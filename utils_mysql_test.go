package db2struct

import (
	"testing"

	_ "github.com/go-sql-driver/mysql" // Initialize mysql driver
	. "github.com/smartystreets/goconvey/convey"
)

const testMariadbUsername = "root"
const testMariadbPassword = ""
const testMariadbHost = "127.0.0.1"
const testMariadbPort = 3306
const testMariadbDatabase = "test"

func TestGetColumnsFromMysqlTable(t *testing.T) {
	var testTable = "all_data_types"
	columMap, _, err := GetColumnsFromMysqlTable(testMariadbUsername, testMariadbPassword, testMariadbHost, testMariadbPort, testMariadbDatabase, testTable)
	Convey("Should be able to connect to test database and create columnMap", t, func() {
		So(err, ShouldBeNil)
		So(columMap, ShouldNotBeNil)
		So(*columMap, ShouldNotBeEmpty)
	})

	columMap, _, err = GetColumnsFromMysqlTable(testMariadbUsername, testMariadbPassword, "doesnotexists", testMariadbPort, testMariadbDatabase, testTable)
	Convey("Should get an error connecting to test database", t, func() {
		So(err, ShouldNotBeNil)
		So(columMap, ShouldBeNil)
	})
}
