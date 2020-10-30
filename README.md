# db2struct 
[![Build Status](https://travis-ci.org/Shelnutt2/db2struct.svg?branch=master)](https://travis-ci.org/Shelnutt2/db2struct) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Shelnutt2/db2struct)](https://goreportcard.com/report/github.com/Shelnutt2/db2struct) 
[![Coverage Status](https://coveralls.io/repos/github/Shelnutt2/db2struct/badge.svg?branch=1-add-coveralls-support)](https://coveralls.io/github/Shelnutt2/db2struct?branch=1-add-coveralls-support) 
[![GoDoc](https://godoc.org/github.com/Shelnutt2/db2struct?status.svg)](https://godoc.org/github.com/Shelnutt2/db2struct)  

The db2struct package produces a usable golang struct from a given database table for use in a .go file.

By reading details from the database about the column structure, db2struct generates a go compatible struct type
with the required column names, data types, and annotations.

Generated datatypes include support for nullable columns [sql.NullX types](https://golang.org/pkg/database/sql/#NullBool) or [guregu null.X types](https://github.com/guregu/null)
and the expected basic built in go types.

Db2Struct is based/inspired by the work of ChimeraCoder's gojson package
[gojson](https://github.com/ChimeraCoder/gojson)



## Usage

```BASH
go get github.com/Shelnutt2/db2struct/cmd/db2struct
db2struct --host localhost -d test -t test_table --package myGoPackage --struct testTable -p --user testUser
```

## Example

MySQL table named users with four columns: id (int), user_name (varchar(255)), number_of_logins (int(11),nullable), and LAST_NAME (varchar(255), nullable)  

Example below uses guregu's null package, but without the option it procuded the sql.NullInt64 and so on.
```BASH
db2struct --host localhost -d example.com -t users --package example --struct user -p --user exampleUser --guregu --gorm
```

Output:
```GOLANG

package example

type User struct {
  ID              int   `gorm:"column:id"`
  UserName        string `gorm:"column:user_name"`
  NumberOfLogins  null.Int `gorm:"column:number_of_logins"`
  LastName        null.String `gorm:"column:LAST_NAME"`
}
```

## Supported Databases

Currently Supported
-   MariaDB
-   MySQL

Planned Support
-   PostgreSQL
-   Oracle
-   Microsoft SQL Server

### MariaDB/MySQL

Structures are created by querying the INFORMATION_SCHEMA.Columns table and then formatting the types, column names,
and metadata to create a usable go compatible struct type.

#### Supported Datatypes

Currently only a limited number of MariaDB/MySQL datatypes are supported. Initial support includes:
-   tinyint (sql.NullInt64 or null.Int)
-   int      (sql.NullInt64 or null.Int)
-   smallint      (sql.NullInt64 or null.Int)
-   mediumint      (sql.NullInt64 or null.Int)
-   bigint (sql.NullInt64 or null.Int)
-   decimal (sql.NullFloat64 or null.Float)
-   float (sql.NullFloat64 or null.Float)
-   double (sql.NullFloat64 or null.Float)
-   datetime (null.Time)
-   time  (null.Time)
-   date (null.Time)
-   timestamp (null.Time)
-   var (sql.String or null.String)
-   enum (sql.String or null.String)
-   varchar (sql.String or null.String)
-   longtext (sql.String or null.String)
-   mediumtext (sql.String or null.String)
-   text (sql.String or null.String)
-   tinytext (sql.String or null.String)
-   binary
-   blob
-   longblob
-   mediumblob
-   varbinary
-   json
