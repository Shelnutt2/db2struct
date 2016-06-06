# db2struct [![Build Status](https://travis-ci.org/Shelnutt2/db2struct.svg?branch=master)](https://travis-ci.org/Shelnutt2/db2struct) [![Coverage Status](https://coveralls.io/repos/github/Shelnutt2/db2struct/badge.svg?branch=1-add-coveralls-support)](https://coveralls.io/github/Shelnutt2/db2struct?branch=1-add-coveralls-support)

This package produces a golang struct from a db table.

It reads details from the database about the column structure.


This is based on the work by ChimeraCoder with
[gojson](https://github.com/ChimeraCoder/gojson)

## Usage

```BASH
go get github.com/Shelnutt2/db2struct/db2struct
db2struct --host localhost -d test -t test_table --package myGoPackage --struct testTable -p --user testUser
```
## Supported Database

Right now Only Mariadb/Mysql is supported, long term plans are to support
postgres and others.

### Mariadb

Structures are created by querying the INFORMATION_SCHEMA.Columns and returning details.


#### Supported Datatypes

Currently only a small portion of mariadb datatypes are supported.

Were applicable sql.Null versions are also supported

-   int
-   bigint
-   decimal
-   float
-   double
-   datetime
-   time
-   date
-   timestamp
