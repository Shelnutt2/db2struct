package db2struct

import (
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"strings"
	"unicode"
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

// Generate Given a Column map with datatypes and a name structName,
// attempts to generate a struct definition
func Generate(columnTypes map[string]map[string]string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) ([]byte, error) {
	src := fmt.Sprintf("package %s\ntype %s %s}",
		pkgName,
		structName,
		generateTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes))
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX
		if gureguTypes == true {
			valueType = mysqlTypeToGureguType(mysqlType["value"], nullable)
		} else {
			valueType = mysqlTypeToGoType(mysqlType["value"], nullable)
		}

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\"", key))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n%s %s",
				fieldName,
				valueType)
		}
	}
	return structure
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

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint":
		if nullable {
			return "sql.NullInt64"
		}
		return "int"
	case "int":
		if nullable {
			return "sql.NullInt64"
		}
		return "int"
	case "bigint":
		if nullable {
			return "sql.NullInt64"
		}
		return "int64"
	case "varchar":
		if nullable {
			return "sql.NullString"
		}
		return "string"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "decimal":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	case "float":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float32"
	case "double":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	}

	return ""
}

// mysqlTypeToGureguType converts the mysql types to go compatible guregu (https://github.com/guregu/null) types
func mysqlTypeToGureguType(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint":
		if nullable {
			return "null.Int"
		}
		return "int"
	case "int":
		if nullable {
			return "null.Int"
		}
		return "int"
	case "bigint":
		if nullable {
			return "null.Int"
		}
		return "int64"
	case "varchar":
		if nullable {
			return "null.String"
		}
		return "string"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "decimal":
		if nullable {
			return "null.Float"
		}
		return "float64"
	case "float":
		if nullable {
			return "null.Float"
		}
		return "float32"
	case "double":
		if nullable {
			return "null.Float"
		}
		return "float64"
	}

	return ""
}
