package db2struct

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMysqlStringGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	NullStringColumn sql.NullString
	StringColumn     string
}
`

	columnMap := map[string]map[string]string{
		"stringColumn":     {"nullable": "NO", "value": "varchar"},
		"nullStringColumn": {"nullable": "YES", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlDateGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	DateColumn      time.Time
	DateTimeColumn  time.Time
	TimeColumn      time.Time
	TimeStampColumn time.Time
}
`

	columnMap := map[string]map[string]string{
		"DateColumn":      {"nullable": "NO", "value": "date"},
		"DateTimeColumn":  {"nullable": "NO", "value": "datetime"},
		"TimeColumn":      {"nullable": "NO", "value": "time"},
		"TimeStampColumn": {"nullable": "NO", "value": "timestamp"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlFloatGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	DecimalColumn     float64
	DoubleColumn      float64
	FloatColumn       float32
	NullDecimalColumn sql.NullFloat64
	NullDoubleColumn  sql.NullFloat64
	NullFloatColumn   sql.NullFloat64
}
`

	columnMap := map[string]map[string]string{
		"floatColumn":       {"nullable": "NO", "value": "float"},
		"nullFloatColumn":   {"nullable": "YES", "value": "float"},
		"doubleColumn":      {"nullable": "NO", "value": "double"},
		"nullDoubleColumn":  {"nullable": "YES", "value": "double"},
		"decimalColumn":     {"nullable": "NO", "value": "decimal"},
		"nullDecimalColumn": {"nullable": "YES", "value": "decimal"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlIntGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	BigIntColumn     int64
	IntColumn        int
	NullBigIntColumn sql.NullInt64
	NullIntColumn    sql.NullInt64
}
`

	columnMap := map[string]map[string]string{
		"intColumn":        {"nullable": "NO", "value": "int"},
		"nullIntColumn":    {"nullable": "YES", "value": "int"},
		"bigIntColumn":     {"nullable": "NO", "value": "bigint"},
		"nullBigIntColumn": {"nullable": "YES", "value": "bigint"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlJSONStringGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	NullStringColumn sql.NullString ` + "`json:\"nullStringColumn\"`" + `
	StringColumn     string         ` + "`json:\"stringColumn\"`" + `
}
`

	columnMap := map[string]map[string]string{
		"stringColumn":     {"nullable": "NO", "value": "varchar"},
		"nullStringColumn": {"nullable": "YES", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", true, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlGormStringGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	NullStringColumn sql.NullString ` + "`gorm:\"column:nullStringColumn\"`" + `
	StringColumn     string         ` + "`gorm:\"column:stringColumn\"`" + `
}
`

	columnMap := map[string]map[string]string{
		"stringColumn":     {"nullable": "NO", "value": "varchar"},
		"nullStringColumn": {"nullable": "YES", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, true)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlStringWithIntGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	OneStringColumn string
}
`

	columnMap := map[string]map[string]string{
		"1stringColumn": {"nullable": "NO", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlStringWithUnderscoresGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	StringColumn string
}
`

	columnMap := map[string]map[string]string{
		"string_Column": {"nullable": "NO", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlStringWithCommonInitialismGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	API string
}
`

	columnMap := map[string]map[string]string{
		"API": {"nullable": "NO", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "testStruct", "test", false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}
