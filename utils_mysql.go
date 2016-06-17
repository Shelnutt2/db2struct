package db2struct

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// GetColumnsFromMysqlTable Select column details from information schema and return map of map
func GetColumnsFromMysqlTable(mariadbUser string, mariadbPassword string, mariadbHost string, mariadbPort int, mariadbDatabase string, mariadbTable string) (*map[string]map[string]string, error) {

	var err error
	var db *sql.DB
	if mariadbPassword != "" {
		db, err = sql.Open("mysql", mariadbUser+":"+mariadbPassword+"@tcp("+mariadbHost+":"+strconv.Itoa(mariadbPort)+")/"+mariadbDatabase+"?&parseTime=True")
	} else {
		db, err = sql.Open("mysql", mariadbUser+"@tcp("+mariadbHost+":"+strconv.Itoa(mariadbPort)+")/"+mariadbDatabase+"?&parseTime=True")
	}
	defer db.Close()

	// Check for error in db, note this does not check connectivity but does check uri
	if err != nil {
		fmt.Println("Error opening mysql db: " + err.Error())
		return nil, err
	}

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA
	columnDataTypeQuery := "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery, mariadbDatabase, mariadbTable)

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
		var column string
		var dataType string
		var nullable string
		rows.Scan(&column, &dataType, &nullable)

		columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable}
	}

	return &columnDataTypes, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateMysqlTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX
		if gureguTypes == true {
			valueType = mysqlTypeToGureguType(mysqlType["value"], nullable)
		} else {
			valueType = mysqlTypeToGoType(mysqlType["value"], nullable)
		}

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\"", key))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n%s %s",
				fieldName,
				valueType)
		}
	}
	return structure
}

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint":
		if nullable {
			return "sql.NullInt64"
		}
		return "int"
	case "int":
		if nullable {
			return "sql.NullInt64"
		}
		return "int"
	case "bigint":
		if nullable {
			return "sql.NullInt64"
		}
		return "int64"
	case "varchar":
		if nullable {
			return "sql.NullString"
		}
		return "string"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "decimal":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	case "float":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float32"
	case "double":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	}

	return ""
}

// mysqlTypeToGureguType converts the mysql types to go compatible guregu (https://github.com/guregu/null) types
func mysqlTypeToGureguType(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint":
		if nullable {
			return "null.Int"
		}
		return "int"
	case "int":
		if nullable {
			return "null.Int"
		}
		return "int"
	case "bigint":
		if nullable {
			return "null.Int"
		}
		return "int64"
	case "varchar":
		if nullable {
			return "null.String"
		}
		return "string"
	case "datetime":
		if nullable {
			return "null.Time"
		}
		return "time.Time"
	case "date":
		if nullable {
			return "null.Time"
		}
		return "time.Time"
	case "time":
		if nullable {
			return "null.Time"
		}
		return "time.Time"
	case "timestamp":
		if nullable {
			return "null.Time"
		}
		return "time.Time"
	case "decimal":
		if nullable {
			return "null.Float"
		}
		return "float64"
	case "float":
		if nullable {
			return "null.Float"
		}
		return "float32"
	case "double":
		if nullable {
			return "null.Float"
		}
		return "float64"
	}

	return ""
}
