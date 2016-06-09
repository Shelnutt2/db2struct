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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

func TestMysqlIntGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	BigIntColumn      int64
	IntColumn         int
	NullBigIntColumn  sql.NullInt64
	NullIntColumn     sql.NullInt64
	NullTinyIntColumn sql.NullInt64
	TinyIntColumn     int
}
`

	columnMap := map[string]map[string]string{
		"intColumn":         {"nullable": "NO", "value": "int"},
		"nullIntColumn":     {"nullable": "YES", "value": "int"},
		"tinyIntColumn":     {"nullable": "NO", "value": "tinyint"},
		"nullTinyIntColumn": {"nullable": "YES", "value": "tinyint"},
		"bigIntColumn":      {"nullable": "NO", "value": "bigint"},
		"nullBigIntColumn":  {"nullable": "YES", "value": "bigint"},
	}
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", true, false, false)

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

// TableName sets the insert table name for this struct type
func (t *testStruct) TableName() string {
	return "test_table"
}
`

	columnMap := map[string]map[string]string{
		"stringColumn":     {"nullable": "NO", "value": "varchar"},
		"nullStringColumn": {"nullable": "YES", "value": "varchar"},
	}
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, true, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

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
	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, false)

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}

// TestMysqlTypeToGureguType generates the struct and outputs nullable columns as guregu null types
func TestMysqlTypeToGureguType(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	BigInt    int64
	Date      null.Time
	DateTime  null.Time
	Decimal   null.Float
	Double    float64
	Float     null.Float
	Int       null.Int
	Time      time.Time
	TimeStamp null.Time
	TinyInt   int
	VarChar   null.String
}
`

	columnMap := map[string]map[string]string{
		"VarChar":   {"nullable": "YES", "value": "varchar"},
		"TinyInt":   {"nullable": "NO", "value": "tinyint"},
		"Int":       {"nullable": "YES", "value": "int"},
		"BigInt":    {"nullable": "NO", "value": "bigint"},
		"Decimal":   {"nullable": "YES", "value": "decimal"},
		"Float":     {"nullable": "YES", "value": "float"},
		"Double":    {"nullable": "NO", "value": "double"},
		"DateTime":  {"nullable": "YES", "value": "datetime"},
		"Time":      {"nullable": "NO", "value": "time"},
		"Date":      {"nullable": "YES", "value": "date"},
		"TimeStamp": {"nullable": "YES", "value": "timestamp"},
	}

	bytes, err := Generate(columnMap, "test_table", "testStruct", "test", false, false, true)

	Convey("Should be able to generate map for guregu types", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}
