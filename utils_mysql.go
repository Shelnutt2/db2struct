package db2struct

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func connect(mariadbUser string, mariadbPassword string, mariadbHost string, mariadbPort int, mariadbDatabase string) (*sql.DB, error) {
	if mariadbPassword != "" {
		return sql.Open("mysql", mariadbUser+":"+mariadbPassword+"@tcp("+mariadbHost+":"+strconv.Itoa(mariadbPort)+")/"+mariadbDatabase+"?&parseTime=True")
	}
	return sql.Open("mysql", mariadbUser+"@tcp("+mariadbHost+":"+strconv.Itoa(mariadbPort)+")/"+mariadbDatabase+"?&parseTime=True")
}

// GetTablesFromMysqlSchema Select table details from information schema
func GetTablesFromMysqlSchema(mariadbUser string, mariadbPassword string, mariadbHost string, mariadbPort int, mariadbDatabase string) ([]string, error) {

	db, err := connect(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase)
	// Check for error in db, note this does not check connectivity but does check uri
	if err != nil {
		fmt.Println("Error opening mysql db: " + err.Error())
		return nil, err
	}
	defer db.Close()

	tableNamesSorted := []string{}

	// Select table data from INFORMATION_SCHEMA
	tableDataTypeQuery := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? order by TABLE_NAME asc"

	if Debug {
		fmt.Println("running: " + tableDataTypeQuery)
	}

	rows, err := db.Query(tableDataTypeQuery, mariadbDatabase)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)

		tableNamesSorted = append(tableNamesSorted, tableName)
	}

	return tableNamesSorted, err
}

// GetColumnsFromMysqlTable Select column details from information schema and return map of map
func GetColumnsFromMysqlTable(mariadbUser string, mariadbPassword string, mariadbHost string, mariadbPort int, mariadbDatabase string, mariadbTable string) (*map[string]map[string]string, []string, error) {

	db, err := connect(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase)
	// Check for error in db, note this does not check connectivity but does check uri
	if err != nil {
		fmt.Println("Error opening mysql db: " + err.Error())
		return nil, nil, err
	}
	defer db.Close()

	columnNamesSorted := []string{}

	// Store column as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select column data from INFORMATION_SCHEMA
	columnDataTypeQuery := "SELECT COLUMN_TYPE, COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ? order by ordinal_position asc"

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery, mariadbDatabase, mariadbTable)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var columnType string
		var column string
		var columnKey string
		var dataType string
		var nullable string
		var comment string
		rows.Scan(&columnType, &column, &columnKey, &dataType, &nullable, &comment)

		columnDataTypes[column] = map[string]string{"columnType": columnType, "value": dataType, "nullable": nullable, "primary": columnKey, "comment": comment}
		columnNamesSorted = append(columnNamesSorted, column)
	}

	return &columnDataTypes, columnNamesSorted, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateMysqlTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = mysqlTypeToGoType(mysqlType["value"], nullable, gureguTypes, mysqlType["columnType"])

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if dbAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("db:\"%s\"", key))
		}

		if len(annotations) > 0 {
			// add colulmn comment
			comment := mysqlType["comment"]
			structure += fmt.Sprintf("\n%s %s `%s`", fieldName, valueType, strings.Join(annotations, " "))
			if comment != "" {
				structure += fmt.Sprintf("  // %s", comment)
			}
		} else {
			structure += fmt.Sprintf("\n%s %s", fieldName, valueType)
		}
	}
	return structure
}

func isSigned(columnType string) bool {
	return !strings.Contains(columnType, "unsigned")
}

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool, columnType string) string {
	signed := isSigned(columnType)
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt // Unsiged will fit in this 64 bit number
		}
		return golangInt
	case "bigint":
		// Until they support this https://go-review.googlesource.com/c/go/+/344410/
		if !signed {
			if nullable {
				return golangNullUint64
			}
			return golangUint64
		}
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "date", "time", "datetime", "timestamp":
		if nullable {
			if gureguTypes {
				return gureguNullTime
			}
			return golangNullTime
		}
		return golangTime
	// case "time":
	// if nullable {
	// return skyhopNullTime
	// }
	// return skyhopTime
	case "float", "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "point":
		if nullable {
			return skyhopNullPoint
		}
		return skyhopPoint
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		// This assumes that any binary(16) is a uuid
		if columnType == "binary(16)" {
			if nullable {
				return skyhopNullBinaryUUID
			}
			return skyhopBinaryUUID
		}
		if nullable {
			return golangNullByteArray
		}
		return golangByteArray
	}
	return ""
}

// Generate go struct entries for a map[string]interface{} structure
func generateServiceTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = dbTypeToPrimitiveType(mysqlType["value"], nullable, gureguTypes, mysqlType["columnType"])

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if dbAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("db:\"%s\"", key))
		}

		if len(annotations) > 0 {
			// add colulmn comment
			comment := mysqlType["comment"]
			structure += fmt.Sprintf("\n%s %s `%s`", fieldName, valueType, strings.Join(annotations, " "))
			if comment != "" {
				structure += fmt.Sprintf("  // %s", comment)
			}
		} else {
			structure += fmt.Sprintf("\n%s %s", fieldName, valueType)
		}
	}
	return structure
}

// Generate go struct entries for a map[string]interface{} structure
func generateServiceConversion(obj map[string]map[string]string, columnsSorted []string) string {
	var structure string

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		// Get the corresponding go value type for this mysql type
		conv := mysqlTypeConversion(mysqlType["value"], nullable, false, mysqlType["columnType"])
		fieldName := fmtFieldName(stringifyFirstChar(key))

		o := fmt.Sprintf("obj.%s", fieldName)
		c := fmt.Sprintf(conv, o)
		structure += fmt.Sprintf("\n%s: %s,", fieldName, c)
	}
	return structure
}

func generateSqlConversion(obj map[string]map[string]string, columnsSorted []string) string {
	var structure string

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		// Get the corresponding go value type for this mysql type
		conv := objTypeConversion(mysqlType["value"], nullable, false, mysqlType["columnType"])
		fieldName := fmtFieldName(stringifyFirstChar(key))

		if fieldName == "Id" {
			structure += fmt.Sprintf("\n%s: helper.NewUUIDIfBlank(obj.Id),", fieldName)
		} else {
			o := fmt.Sprintf("obj.%s", fieldName)
			c := fmt.Sprintf(conv, o)
			structure += fmt.Sprintf("\n%s: %s,", fieldName, c)
		}
	}
	return structure
}

func dbTypeToPrimitiveType(mysqlType string, nullable bool, gureguTypes bool, columnType string) string {
	signed := isSigned(columnType)
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			return golangNullInt
		}
		return golangInt
	case "bigint":
		// Until they support this https://go-review.googlesource.com/c/go/+/344410/
		if !signed {
			if nullable {
				return golangNullUint64
			}
			return golangUint64
		}
		if nullable {
			return golangNullInt64
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			return golangNullString
		}
		return "string"
	case "date", "time", "datetime", "timestamp":
		if nullable {
			return golangNullTime
		}
		return golangTime
	case "decimal", "double", "float":
		if nullable {
			return golangNullFloat64
		}
		return golangFloat64
	// case "float":
	// 	if nullable {
	// 		return golangNullFloat32
	// 	}
	// 	return golangFloat32
	case "point":
		if nullable {
			return skyhopNullPoint
		}
		return skyhopPoint
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		// This assumes that any binary(16) is a uuid
		if columnType == "binary(16)" {
			if nullable {
				return golangNullString
			}
			return "string"
		}
		if nullable {
			return golangNullByteArray
		}
		return golangByteArray
	}
	return ""
}

// mysqlTypeConversion converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeConversion(mysqlType string, nullable bool, gureguTypes bool, columnType string) string {
	signed := isSigned(columnType)
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			return "helper.SqlNullInt32ToNullableInt32(%v)"
		}
	case "bigint":
		if !signed {
			if nullable {
				return "helper.SqlNullInt64ToNullableInt64(%v)"
			}
		}
		if nullable {
			return "helper.SqlNullInt32ToNullableInt32(%v)"
		}
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			return "helper.SqlNullStringToNullableString(%v)"
		}
	case "decimal", "double", "float":
		if nullable {
			return "helper.SqlNullFloat64ToNullableFloat64(%v)"
		}
	// case "float":
	// 	if nullable {
	// 		return "helper.SqlNullFloat32ToNullableFloat32(%v)"
	// 	}
	// case "date", "time", "datetime", "timestamp":
	// 	if nullable {
	// 		return golangNullTime
	// 	}
	// case "point":
	// 	if nullable {
	// 		return skyhopNullPoint
	// 	}
	// 	return skyhopPoint
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		// This assumes that any binary(16) is a uuid
		if columnType == "binary(16)" {
			if nullable {
				return "helper.NullableUUIDToString(%v)"
			}
			return "helper.UUIDToString(%v)"
		}
		// if nullable {
		// 	return golangNullByteArray
		// }
		// return golangByteArray
	}

	return "%v"
}

// mysqlTypeConversion converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func objTypeConversion(mysqlType string, nullable bool, gureguTypes bool, columnType string) string {
	signed := isSigned(columnType)
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			return "helper.NullableInt32ToSqlNullInt32(%v)"
		}
	case "bigint":
		if !signed {
			if nullable {
				return "helper.NullableInt64ToSqlNullInt64(%v)"
			}
		}
		if nullable {
			return "helper.NullableInt32ToSqlNullInt32(%v)"
		}
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			return "helper.NullableStringToSqlNullString(%v)"
		}
	case "decimal", "double", "float":
		if nullable {
			return "helper.NullableFloat64ToSqlNullFloat64(%v)"
		}
	// case "float":
	// 	if nullable {
	// 		return "helper.SqlNullFloat32ToNullableFloat32(%v)"
	// 	}
	// case "date", "time", "datetime", "timestamp":
	// 	if nullable {
	// 		return golangNullTime
	// 	}
	// case "point":
	// 	if nullable {
	// 		return skyhopNullPoint
	// 	}
	// 	return skyhopPoint
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		// This assumes that any binary(16) is a uuid
		if columnType == "binary(16)" {
			if nullable {
				return "helper.StringToNullableUUID(%v)"
			}
			return "helper.StringToUUID(%v)"
		}
		// if nullable {
		// 	return golangNullByteArray
		// }
		// return golangByteArray
	}

	return "%v"
}

// UnwrapNullableString(w *wrapperspb.StringValue)
// UnwrapString(w *wrapperspb.StringValue)
// WrapString(s string)
// WrapNullableString(s *string)
// WrapNullableInt64(n *int64)
// UnwrapNullableInt64(w *wrapperspb.Int64Value)
// StringToNullableUUID(s *string)
// NullableUUIDToString(id *database.BinaryUUID)
// WrapTimeAsNullableString(t *time.Time)
// NullableStringToSqlNullString(s *string)
// SqlNullStringToNullableString(s sql.NullString)
// NullableInt64ToSqlNullInt64(n *int64)
// SqlNullInt64ToNullableInt64(n sql.NullInt64)
// NullableFloat64ToSqlNullInt64(n)
// SqlNullFloat64ToNullableInt64(n)
