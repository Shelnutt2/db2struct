# mysql-to-struct

This package produces a golang struct from a mysql table.

It reads details from the INFORMATION_SCHEMA.Columns about the column struct
of the table.

## Usage

```BASH
go get github.com/Shelnutt2/mysql-to-struct
mysql-to-struct --host localhost -d test -t test_table --package myGoPackage --struct testTable -p --user testUser
```
