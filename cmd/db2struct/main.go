package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	goopt "github.com/droundy/goopt"
	"github.com/ericksonjoseph/db2struct"
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
var dbAnnotation = goopt.Flag([]string{"--db"}, []string{}, "Add db annotations (tags)", "")
var gureguTypes = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")
var targetFile = goopt.String([]string{"--target"}, "", "Save file path")

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
	} else if mariadbPassword == nil {
		p := ""
		mariadbPassword = &p
	}

	if *verbose {
		fmt.Println("Connecting to mysql server " + mariadbHost + ":" + strconv.Itoa(*mariadbPort))
	}

	if mariadbDatabase == nil || *mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}
	// If packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "newpackage"
	}

	var tablesSorted []string

	// If no table is specified, process all tables in the schema
	if mariadbTable == nil || *mariadbTable == "" {
		var tErr error
		tablesSorted, tErr = db2struct.GetTablesFromMysqlSchema(*mariadbUser, *mariadbPassword, mariadbHost, *mariadbPort, *mariadbDatabase)
		if tErr != nil {
			fmt.Printf("Error in selecting table data from mysql information schema %s", tErr)
			return
		}
		fmt.Printf("--table flag missing so we will process all tables in the schema i.e. %+v\n", tablesSorted)
	} else {
		tablesSorted = []string{*mariadbTable}
	}

	heading := fmt.Sprintf("// This file was automatically generated. Do not edit.\n\npackage %s\n\nimport (\n\t\"database/sql\"\n\t\"time\"\n\n\t\"github.com/skyhop-tech/go-sky/internal/database\"\n)", *packageName)

	var file *os.File

	if targetFile != nil && *targetFile != "" {
		var err error
		file, err = os.OpenFile(*targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Open File fail: " + err.Error())
			return
		}
	}

	if file != nil {
		length, err := file.WriteString(fmt.Sprintf("%s", heading))
		if err != nil {
			fmt.Println("Save File fail: " + err.Error())
			return
		}
		fmt.Printf("wrote %d bytes\n", length)
	} else {
		fmt.Println(string(heading))
	}

	for _, table := range tablesSorted {

		columnDataTypes, columnsSorted, err := db2struct.GetColumnsFromMysqlTable(*mariadbUser, *mariadbPassword, mariadbHost, *mariadbPort, *mariadbDatabase, table)
		if err != nil {
			fmt.Printf("Error in selecting column data from mysql information schema %s", err)
			return
		}

		var stName string
		// If structName is not set we need to default it to the table name
		if structName == nil || *structName == "" {
			stName = strings.Title(table)
		} else {
			stName = *structName
			// If it is set we need to default it to the table name
			if len(tablesSorted) > 1 {
				stName = fmt.Sprintf("%s%s", *structName, strings.Title(table))
			}
		}

		// Generate struct string based on columnDataTypes
		struc, err := db2struct.Generate(*columnDataTypes, columnsSorted, table, stName, *jsonAnnotation, *gormAnnotation, *dbAnnotation, *gureguTypes)
		if err != nil {
			fmt.Println("Error in creating struct from json: " + err.Error())
			return
		}
		if file != nil {
			length, err := file.WriteString(fmt.Sprintf("%s", struc))
			if err != nil {
				fmt.Println("Save File fail: " + err.Error())
				return
			}
			fmt.Printf("wrote %d bytes\n", length)
		} else {
			fmt.Println(string(struc))
		}
	}
}

func getMariadbPassword(password string) error {
	mariadbPassword = new(string)
	*mariadbPassword = password
	return nil
}
