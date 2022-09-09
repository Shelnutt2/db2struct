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
	golangNullInt     = "*int"
	golangNullInt32   = "*int32"
	golangNullInt64   = "*int64"
	golangNullFloat   = "*float"
	golangNullFloat32 = "*float32"
	golangNullFloat64 = "*float64"
	golangNullString  = "*string"

	golangByteArray      = "[]byte"
	golangNullByteArray  = "*[]byte"
	golangUint64         = "uint64"
	golangNullUint64     = "*uint64"
	gureguNullInt        = "null.Int"
	sqlNullInt           = "sql.NullInt64"
	sqlNullInt32         = "sql.NullInt32"
	golangInt            = "int"
	golangInt32          = "int32"
	golangInt64          = "int64"
	gureguNullFloat      = "null.Float"
	sqlNullFloat         = "sql.NullFloat64"
	golangFloat          = "float"
	golangFloat32        = "float32"
	golangFloat64        = "float64"
	gureguNullString     = "null.String"
	sqlNullString        = "sql.NullString"
	gureguNullTime       = "null.Time"
	golangTime           = "time.Time"
	golangNullTime       = "*time.Time"
	skyhopBinaryUUID     = "database.BinaryUUID"
	skyhopNullBinaryUUID = "*database.BinaryUUID"
	skyhopPoint          = "database.Point"
	skyhopNullPoint      = "*database.Point"
	// skyhopTime           = "database.Time"
	// skyhopNullTime       = "*database.Time"
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
	//"ID":    true,
	"IP":   true,
	"JSON": true,
	"LHS":  true,
	"QPS":  true,
	"RAM":  true,
	"RHS":  true,
	"RPC":  true,
	"SLA":  true,
	"SMTP": true,
	"SSH":  true,
	"TLS":  true,
	"TTL":  true,
	"UI":   true,
	"UID":  true,
	"UUID": true,
	"URI":  true,
	"URL":  true,
	"UTF8": true,
	"VM":   true,
	"XML":  true,
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
func Generate(columnTypes map[string]map[string]string, columnsSorted []string, tableName string, structName string, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) ([]byte, error) {
	var dbTypes string
	dbTypes = generateMysqlTypes(columnTypes, columnsSorted, 0, jsonAnnotation, gormAnnotation, dbAnnotation, gureguTypes)
	src := fmt.Sprintf("\n\ntype %s %s\n}",
		structName,
		dbTypes)
	if gormAnnotation == true || dbAnnotation == true {
		tableNameFunc := "// TableName sets the insert table name for this struct type\n" +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

func GenerateServiceObject(columnTypes map[string]map[string]string, columnsSorted []string, internalDir, tableName, structName string, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) ([]byte, error) {

	title := structNameToTitle(structName)
	dbTypes := generateServiceTypes(columnTypes, columnsSorted, 0, jsonAnnotation, gormAnnotation, dbAnnotation, gureguTypes)
	strct := fmt.Sprintf("\n\ntype %s %s\n}", title, dbTypes)
	heading := "// This file was automatically generated. Do not edit.\n\n" +
		"package object\n\n" +
		"import (\n" +
		"	\"time\"\n" +
		"\n" +
		"	\"github.com/skyhop-tech/go-sky/cmd/manifest-api/internal/database\"\n" +
		"	\"github.com/skyhop-tech/go-sky/internal/helper\"\n" +
		")\n\n"
	from := "\n" +
		"func " + title + "FromModel(obj *database." + structName + ") *" + title + " {\n" +
		"	return &" + title + "{\n" +
		"		%v\n" +
		"	}\n" +
		"}\n"
	fromcnv := generateServiceConversion(columnTypes, columnsSorted)
	from = fmt.Sprintf(from, fromcnv)
	to := "\n" +
		"func New" + title + "Model(obj *" + title + ") *database." + structName + " {\n" +
		"	return &database." + structName + "{\n" +
		"		%v\n" +
		"	}\n" +
		"}\n"
	tocnv := generateSqlConversion(columnTypes, columnsSorted)
	to = fmt.Sprintf(to, tocnv)

	src := fmt.Sprintf("%s\n%s\n%s\n%s\n", heading, strct, from, to)
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

func GenerateServiceObjectText(columnTypes map[string]map[string]string, columnsSorted []string, internalDir, tableName string, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) ([]byte, error) {

	structName := strings.Title(tableName)
	title := structNameToTitle(structName)
	dbTypes := generateServiceTypes(columnTypes, columnsSorted, 0, jsonAnnotation, gormAnnotation, dbAnnotation, gureguTypes)
	strct := fmt.Sprintf("\n\ntype %s %s\n}", title, dbTypes)
	from := "\n" +
		"func " + title + "FromModel(obj *database." + structName + ") *" + title + " {\n" +
		"	return &" + title + "{\n" +
		"		%v\n" +
		"	}\n" +
		"}\n"
	fromcnv := generateServiceConversion(columnTypes, columnsSorted)
	from = fmt.Sprintf(from, fromcnv)
	to := "\n" +
		"func New" + title + "Model(obj *" + title + ") *database." + structName + " {\n" +
		"	return &database." + structName + "{\n" +
		"		%v\n" +
		"	}\n" +
		"}\n"
	tocnv := generateSqlConversion(columnTypes, columnsSorted)
	to = fmt.Sprintf(to, tocnv)

	src := fmt.Sprintf("\n%s\n%s\n%s\n", strct, from, to)
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

func structNameToTitle(name string) string {
	title := name
	for {
		ind := strings.Index(title, "_")
		if ind < 0 {
			break
		}
		title = title[:ind] + strings.Title(title[ind+1:])
	}
	return title
}

// Generate Given a Column map with datatypes and a name structName,
// attempts to generate a struct definition
func GenerateCRUD(columnTypes map[string]map[string]string, columnsSorted []string, tableName string, structName string, jsonAnnotation bool, gormAnnotation bool, dbAnnotation bool, gureguTypes bool) ([]byte, error) {
	var src string
	structNameTitle := structName
	tableNameTitle := tableName
	for {
		ind := strings.Index(structNameTitle, "_")
		if ind < 0 {
			break
		}
		structNameTitle = structNameTitle[:ind] + strings.Title(structNameTitle[ind+1:])
	}
	for {
		ind := strings.Index(tableNameTitle, "_")
		if ind < 0 {
			break
		}
		tableNameTitle = tableNameTitle[:ind] + strings.Title(tableNameTitle[ind+1:])
	}
	tableName = strings.Title(tableName)
	heading := "// This file was automatically generated. Do not edit.\n\n" +
		"package database\n\n" +
		"import (\n" +
		"	\"context\"\n" +
		"	\"fmt\"\n" +
		"	\"time\"\n" +
		"\n" +
		"	\"github.com/pkg/errors\"\n" +
		"	\"github.com/skyhop-tech/go-sky/internal/database\"\n" +
		")\n\n" +
		"var (\n" +
		"	" + tableNameTitle + "TableName = (&" + tableName + "{}).TableName()\n" +
		"	" + tableNameTitle + "TableColumns []string\n" +
		"	" + tableNameTitle + "TableColumnsEscaped []string\n" +
		")\n\n" +
		"func init() {\n" +
		"	" + tableNameTitle + "TableColumns = database.Columns(&" + tableName + "{})\n" +
		"	var escaped []string\n" +
		"	for _, c := range " + tableNameTitle + "TableColumns {\n" +
		"		escaped = append(escaped, fmt.Sprintf(\"`%v`\", c))\n" +
		"	}\n" +
		"	" + tableNameTitle + "TableColumnsEscaped = escaped\n" +
		"}\n"
	createFunc := "// Create" + structNameTitle + "s - insert many\n" +
		"func (c *Client) Create" + structNameTitle + "s(ctx context.Context, input []*" + structName + ") error {\n" +
		"	var insertData []any\n" +
		"	for _, v := range input {\n" +
		"		insertData = append(insertData, v)\n" +
		"	}\n\n" +
		"	escapedColumns, data, err := database.PrepareBulkInsert(" + tableNameTitle + "TableColumns, insertData)\n" +
		"	if err != nil {\n" +
		"		return errors.Wrap(err, \"prepare bulk insert\")\n" +
		"	}\n\n" +
		"	_, err = database.InsertMany(ctx, c.logger, c.db, " + tableNameTitle + "TableName, escapedColumns, data)\n" +
		"	if err != nil {\n" +
		"		return errors.Wrap(err, \"create " + strings.ToLower(structName) + "\")\n" +
		"	}\n\n" +
		"	return nil\n" +
		"}\n"
	listFunc := "// Create" + structNameTitle + "s - insert many\n" +
		"func (c *Client) List" + structNameTitle + "s(ctx context.Context, filters map[string]interface{}, pageSize, pageToken int32) ([]*" + structName + ", error) {" +
		"	var results []*" + structName + "\n\n" +
		"	err := database.GetManyByFilters(ctx, c.logger, c.db, " + tableNameTitle + "TableColumnsEscaped, " + tableNameTitle + "TableName, filters, pageSize, pageToken, func() interface{} {\n" +
		"		return &results\n" +
		"	})\n\n" +
		"	if err != nil {\n" +
		"		return nil, errors.Wrap(err, \"list " + strings.ToLower(structName) + "\")\n" +
		"	}\n\n" +
		"	return results, nil\n" +
		"}\n"
	updateFunc := "" +
		"func (c *Client) Update" + structNameTitle + "(ctx context.Context, input *" + structName + ") error {\n" +
		"	// update value for field Updated_At\n" +
		"	now := time.Now()\n" +
		"	input.UpdatedAt = &now\n" +
		"	data := database.ToMap(input)\n\n" +
		"	// delete Created At/By fields so we don't override it\n" +
		"	delete(data, \"created_at\")\n" +
		"	delete(data, \"created_by\")\n\n" +
		"	result, err := database.UpdateByField(ctx, c.logger, c.db, " + tableNameTitle + "TableName, \"id\", input.Id, data)\n" +
		"	if err != nil {\n" +
		"		return errors.Wrap(err, \"update " + strings.ToLower(structName) + " by id\")\n" +
		"	}\n\n" +
		"	updated, err := result.RowsAffected()\n" +
		"	if err != nil || updated == 0 {\n" +
		"		return errors.Wrap(err, fmt.Sprintf(\"unable to update " + strings.ToLower(structName) + " by id %v\", input.Id))\n" +
		"	}\n\n" +
		"	return nil\n" +
		"}"

	src = fmt.Sprintf("%v\n%s\n%s\n%s\n%s", heading, src, createFunc, listFunc, updateFunc)
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
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
