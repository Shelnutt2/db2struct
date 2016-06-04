package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerate(t *testing.T) {
	expectedStruct :=
		`package test

type testStruct struct {
	StringColumn string
}
`

	columnMap := map[string]map[string]string{"stringColumn": {"nullable": "NO", "value": "varchar"}}
	bytes, err := Generate(columnMap, "testStruct", "test")

	Convey("Should be able to generate map from string column", t, func() {
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, expectedStruct)
	})
}
