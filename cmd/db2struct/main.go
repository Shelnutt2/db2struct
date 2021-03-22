package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Shelnutt2/db2struct"
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
var projectName = goopt.String([]string{"--project"}, "", "name to set projectName")
var tableNote = goopt.String([]string{"--note"}, "", "table note 表-中文名")

var jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
var gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
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

	log.SetFlags(log.Lshortfile | log.LstdFlags)

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
			log.Println("Error reading password: " + err.Error())
			return
		}
	} else if mariadbPassword == nil {
		p := ""
		mariadbPassword = &p
	}

	if *verbose {
		log.Println("Connecting to mysql server " + mariadbHost + ":" + strconv.Itoa(*mariadbPort))
	}

	if mariadbDatabase == nil || *mariadbDatabase == "" {
		log.Println("Database can not be null")
		return
	}

	if mariadbTable == nil || *mariadbTable == "" {
		log.Println("Table can not be null")
		return
	}

	// If structName is not set we need to default it
	if structName == nil || *structName == "" {
		*structName = "newstruct"
	}
	// If packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "newpackage"
	}

	dp := db2struct.DBParam{
		MariadbUser:     *mariadbUser,
		MariadbPassword: *mariadbPassword,
		MariadbHost:     mariadbHost,
		MariadbPort:     *mariadbPort,
		MariadbDatabase: *mariadbDatabase,
		MariadbTable:    *mariadbTable,
	}

	tp := db2struct.TableParam{
		TableName:      *mariadbTable,
		StructName:     *structName,
		PkgName:        *packageName,
		JsonAnnotation: *jsonAnnotation,
		GormAnnotation: *gormAnnotation,
		GureguTypes:    *gureguTypes,
		TableNote:      *tableNote,
		ProjectName:    *projectName,
	}

	StartCreate(&tp, &dp)

}

func getMariadbPassword(password string) error {
	mariadbPassword = new(string)
	*mariadbPassword = password
	return nil
}

func StartCreate(tp *db2struct.TableParam, dp *db2struct.DBParam) error {

	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(dp)
	if err != nil {
		log.Println("Error in selecting column data information from mysql information schema")
		return err
	}

	// Generate struct string based on columnDataTypes
	struc, err := db2struct.Generate(*columnDataTypes, tp)

	if err != nil {
		log.Println("Error in creating struct from json: " + err.Error())
		return err
	}
	var saveFile string
	if targetFile != nil && *targetFile != "" {
		saveFile = *targetFile + "/model/" + tp.TableName + ".go"
		length, err := db2struct.AutoSaveFile(saveFile, string(struc))
		if err != nil {
			log.Println("open err: ", err)
			return err
		}
		log.Printf("wrote %d bytes\n", length)
	}
	fmt.Println(string(struc))
	fmt.Println("")
	fmt.Println("======================================")
	fmt.Println("save model  file to ", saveFile)
	fmt.Println("")

	return nil
}
