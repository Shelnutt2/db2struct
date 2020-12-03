package db2struct

import (
	"fmt"
	"go/format"
	"strconv"
	"strings"
	"unicode"
)

// Constants for return types of golang
const (
	golangByteArray  = "[]byte"
	gureguNullInt    = "int"
	sqlNullInt       = "int64"
	golangInt        = "int"
	golangInt64      = "int64"
	gureguNullFloat  = "float"
	sqlNullFloat     = "float64"
	golangFloat      = "float"
	golangFloat32    = "float32"
	golangFloat64    = "float64"
	gureguNullString = "string"
	sqlNullString    = "string"
	gureguNullTime   = "time.Time"
	golangTime       = "time.Time"
)

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

//Debug level logging
var Debug = false

// Generate Given a Column map with datatypes and a name structName,
// attempts to generate a struct definition
func Generate(columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) ([]byte, error) {
	var dbTypes string
	dbTypes = generateMysqlTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	strImport := `import (
		"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	)`

	src := fmt.Sprintf("package %s\n\n%s\n\ntype %s %s}",
		pkgName,
		strImport,
		structName,
		dbTypes)
	if gormAnnotation == true {
		//tableNameFunc := "// TableName sets the insert table name for this struct type\n" +
		//	//"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
		//	GetStructTile(structName) + "TableName() string {\n" +
		//	"	return \"" + tableName + "\"" +
		//	"}\n"

		var funList string
		funList += GetTableNameFun(structName, tableName)
		funList += GetDeleteFun(structName)
		src = fmt.Sprintf("%s\n%s", src, funList)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

func GetStructTile(structName string) string {
	str := "func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") "
	return str
}

func GetTableNameFun(structName string, tableName string) string {
	str := `
// TableName  返回数据库名字
func (pLock *MisLock) TableName() string {
	return "mis_lock"
}
`
	res := ReplaceFun(str, structName)
	res = strings.ReplaceAll(res, "mis_lock", tableName)
	return res
}

func GetDeleteFun(structName string) string {
	str := `
// Delete 删除函数，根据ID删除数据
func (pLock *MisLock) Delete() error {
	return db.Engine().Delete(pLock).Error
}
`
	return ReplaceFun(str, structName)
}

func ReplaceFun(str string, structName string) string {
	res := str
	m := strings.ToLower(string(structName[0]))
	res = strings.ReplaceAll(res, "MisLock", structName)
	res = strings.ReplaceAll(res, "pLock", m)
	return res
}

// fmtFieldName formats a string as a struct key
//
// Example:
// 	fmtFieldName("foo_id")
// Output: FooID
func fmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	for len(name) > 0 && name[0] == '_' {
		name = name[1:]
	}

	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}
