package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/shelnutt2/db2struct"
	goopt "github.com/droundy/goopt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
)

var mariadbHost = os.Getenv("MYSQL_HOST")
var mariadbHostPassed = goopt.String([]string{"-H", "--host"}, "", "Host to check mariadb status of")
var mariadbPort = goopt.Int([]string{"--mysql_port"}, 3306, "Specify a port to connect to")
var mariadbTable = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
var mariadbDatabase = goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
var mariadbPassword *string
var mariadbUser = goopt.String([]string{"-u", "--user"}, "user", "user to connect to database")
var verbose = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
var packageName = goopt.String([]string{"--package"}, "", "name to set for package")
var structName = goopt.String([]string{"--struct"}, "", "name to set for struct")

var jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
var gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")

func init() {
	goopt.OptArg([]string{"-p", "--password"}, "", "Mysql password", getMariadbPassword)
	//goopt.ReqArg([]string{"-u", "--user"}, "user", "user to connect to database", setUser)

	// Setup goopts
	goopt.Description = func() string {
		return "Mariadb http Check"
	}
	goopt.Version = "0.0.2"
	goopt.Summary = "db2struct [-H] [-p] [-v] --package pkgName --struct structName --database databaseName --table tableName"

	//Parse options
	goopt.Parse(nil)

}

func main() {

	// Username is required
	if mariadbUser == nil || *mariadbUser == "user" {
		fmt.Println("Username is required! Add it with --user=name")
		return
	}

	// If a mariadb host is passed use it
	if mariadbHostPassed != nil && *mariadbHostPassed != "" {
		mariadbHost = *mariadbHostPassed
	}

	if mariadbPassword != nil && *mariadbPassword == "" {
		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		stringPass := string(pass)
		mariadbPassword = &stringPass
		if err != nil {
			fmt.Println("Error reading password: " + err.Error())
			return
		}
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

	var err error
	var db *sql.DB
	if mariadbPassword != nil {
		db, err = sql.Open("mysql", *mariadbUser+":"+*mariadbPassword+"@tcp("+mariadbHost+":"+strconv.Itoa(*mariadbPort)+")/"+*mariadbDatabase+"?&parseTime=True")
	} else {
		db, err = sql.Open("mysql", *mariadbUser+"@tcp("+mariadbHost+":"+strconv.Itoa(*mariadbPort)+")/"+*mariadbDatabase+"?&parseTime=True")
	}
	defer db.Close()

	// Check for error in db, note this does not check connectivity but does check uri
	if err != nil {
		fmt.Println("Error opening mysql db: " + err.Error())
		return
	}

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA
	columnDataTypeQuery := "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"

	if *verbose {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery, *mariadbDatabase, *mariadbTable)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return
	}
	if db != nil {
		defer rows.Close()
	}

	for rows.Next() {
		var column string
		var dataType string
		var nullable string
		rows.Scan(&column, &dataType, &nullable)

		columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable}
	}

	// Generate struct string based on columnDataTypes
	struc, err := db2struct.Generate(columnDataTypes, *structName, *packageName, *jsonAnnotation, *gormAnnotation)

	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
	}

	fmt.Println(string(struc))

}

func getMariadbPassword(password string) error {
	mariadbPassword = new(string)
	*mariadbPassword = password
	return nil
}
