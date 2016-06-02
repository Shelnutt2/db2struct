package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	goopt "github.com/droundy/goopt"
	_ "github.com/go-sql-driver/mysql"
)

var mariadbHost = os.Getenv("MYSQL_HOST")
var mariadbHostPassed = goopt.String([]string{"-H", "--host"}, "", "Host to check mariadb status of")
var mariadbPort = goopt.Int([]string{"-p", "--mysql_port"}, 3306, "Specify a port to connect to")
var mariadbTable = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
var mariadbDatabase = goopt.String([]string{"-d", "--database"}, "", "Database to for connection")
var mysqlUser = os.Getenv("MYSQL_USERNAME")
var mysqlPassword = os.Getenv("MYSQL_PASSWORD")
var verbose = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
var packageName = goopt.String([]string{"--package"}, "", "name to set for package")
var structName = goopt.String([]string{"--struct"}, "", "name to set for struct")

func init() {
	//Parse options
	goopt.Parse(nil)

	// Setup goopts
	goopt.Description = func() string {
		return "Mariadb http Check"
	}
	goopt.Version = "0.0.1"
	goopt.Summary = "mysql-to-struct [-H] [-p] [-v]"

}

func main() {
	if mariadbHostPassed != nil && *mariadbHostPassed != "" {
		mariadbHost = *mariadbHostPassed
	}
	if *verbose {
		fmt.Println("Connecting to mysql server " + mariadbHost + ":" + strconv.Itoa(*mariadbPort))
	}

	if mariadbDatabase == nil || *mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}

	if mariadbTable == nil || *mariadbTable == "" {
		fmt.Println("Table can not be null")
		return
	}

	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@tcp("+mariadbHost+":"+strconv.Itoa(*mariadbPort)+")/"+*mariadbDatabase+"?&parseTime=True")
	defer db.Close()

	if err != nil {
		fmt.Println("Error opening mysql db: " + err.Error())
		return
	}

	columnDataTypes := make(map[string]string)
	columnDataTypeQuery := "SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"

	if *verbose {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery, *mariadbDatabase, *mariadbTable)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		var dataType string
		rows.Scan(&column, &dataType)

		columnDataTypes[column] = dataType
	}

	struc, err := Generate(columnDataTypes, *structName, *packageName)

	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
	}

	fmt.Println(string(struc))

}
