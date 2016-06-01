package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ChimeraCoder/gojson"
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

func init() {
	//Parse options
	goopt.Parse(nil)

	// Setup goopts
	goopt.Description = func() string {
		return "Mariadb http Check"
	}
	goopt.Version = "0.0.1"
	goopt.Summary = "mysql-to-struct [-H] [-p]"

}

func main() {
	if mariadbHostPassed != nil && *mariadbHostPassed != "" {
		mariadbHost = *mariadbHostPassed
	}
	fmt.Println("Connecting to mysql server " + mariadbHost + ":" + strconv.Itoa(*mariadbPort))

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

	query := "SELECT * FROM " + *mariadbTable + " limit 1"

	fmt.Println("running: " + query)

	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	finalResult := make(map[string]interface{})
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)

		if err != nil {
			fmt.Println("Could not scan row: " + err.Error())
			continue
		}

		tmpStruct := make(map[string]interface{})

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			tmpStruct[col] = v
		}

		finalResult = tmpStruct
		break
	}

	json, err := json.MarshalIndent(finalResult, "", "\t")

	if err != nil {
		fmt.Println("Error creating json: " + err.Error())
	}

	fmt.Println(string(json))

	struc, err := json2struct.Generate(bytes.NewReader(json), "", "")

	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
	}

	fmt.Println(string(struc))

}
