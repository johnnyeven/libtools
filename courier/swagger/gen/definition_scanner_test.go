package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

func TestDefinitionScanner(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"time"
		"github.com/johnnyeven/libtools/courier/enumeration"
	)

	type String string
	type Int int
	type Bool bool
	type SomeBool = enumeration.Bool

	type SliceString []string
	type SliceNamed  []String

	type Struct struct {
		// name
		Name      *string ^^json:"name" validate:"@string[0,)"^^
		Id        **string ^^json:"id,omitempty"^^
		Enum      Enum ^^json:"enum"^^
	}

	// swagger:strfmt date-time
	// 日期
	type Time time.Time

	// swagger:enum
	type Enum int

	const (
		ENUM__ONE Enum = iota + 1 // one
		ENUM__TWO  // two
	)
	`)

	scanner := NewDefinitionScanner(program)

	cases := []struct {
		typeName string
		schema   interface{}
	}{
		{
			typeName: "String",
			schema: map[string]interface{}{
				"type": "string",
			},
		},
		{
			typeName: "Int",
			schema: map[string]interface{}{
				"type":   "integer",
				"format": "int64",
			},
		},
		{
			typeName: "Bool",
			schema: map[string]interface{}{
				"type": "boolean",
			},
		},
		{
			typeName: "SomeBool",
			schema: map[string]interface{}{
				"type": "boolean",
			},
		},
		{
			typeName: "SliceNamed",
			schema: map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"$ref": "#/components/schemas/String",
				},
			},
		},
		{
			typeName: "SliceString",
			schema: map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		{
			typeName: "Struct",
			schema: map[string]interface{}{
				"type":     "object",
				"required": []string{"name", "enum"},
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						XField:        "Name",
						XTagJSON:      "name",
						"description": "name",
						"minLength":   0,
						"type":        "string",
						XPointer:      1,
						XTagValidate:  "@string[0,)",
					},
					"id": map[string]interface{}{
						XField:   "Id",
						XTagJSON: "id,omitempty",
						XPointer: 2,
						"type":   "string",
					},
					"enum": map[string]interface{}{
						"allOf": []map[string]interface{}{
							{
								"$ref": "#/components/schemas/Enum",
							},
							{},
						},
						XField:   "Enum",
						XTagJSON: "enum",
					},
				},
			},
		},
		{
			typeName: "Time",
			schema: map[string]interface{}{
				"description": "日期",
				"type":        "string",
				"format":      "date-time",
			},
		},
		{
			typeName: "Enum",
			schema: map[string]interface{}{
				"type":      "string",
				"enum":      []interface{}{"ONE", "TWO"},
				XEnumVals:   []interface{}{1, 2},
				XEnumLabels: []string{"one", "two"},
				XEnumValues: []string{"ONE", "TWO"},
			},
		},
	}

	query := loaderx.NewQuery(program, pkgImportPath)

	for _, c := range cases {
		schema := scanner.Def(query.TypeName(c.typeName))
		tt.Equal(ToMap(c.schema), ToMap(schema))
	}
}

func TestDefinitionScannerWithError(t *testing.T) {
	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	type PartA struct {
		Name      string ^^json:"name" validate:"@string[0,)"^^
	}

	type PartB struct {
		Name      string ^^json:"name" validate:"@string[0,)"^^
	}

	type Composed struct {
		PartA
		PartB
	}
	`)

	scanner := NewDefinitionScanner(program)

	query := loaderx.NewQuery(program, pkgImportPath)

	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()

	scanner.Def(query.TypeName("Composed"))
}
