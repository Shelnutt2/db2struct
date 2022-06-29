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
	case "date", "datetime", "time", "timestamp":
		if nullable {
			if gureguTypes {
				return gureguNullTime
			}
			return golangNullTime
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
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
