package enums

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/libtools/courier/enumeration/gen"
	swagger_gen "github.com/johnnyeven/libtools/courier/swagger/gen"
)

var serviceEnumMap = map[string]map[string]enumeration.Enum{}

func RegisterEnumFromExtensions(serviceName string, tpe string, vals interface{}, values interface{}, labels interface{}) {
	if vals, ok := vals.([]interface{}); ok {
		values := values.([]interface{})
		labels := labels.([]interface{})

		enum := enumeration.Enum{}

		for i, v := range vals {
			lable := labels[i].(string)
			if lable == "" {
				lable = values[i].(string)
			}
			o := enumeration.EnumOption{
				Val:   v,
				Value: values[i],
				Label: lable,
			}
			enum = append(enum, o)
		}

		RegisterEnum(serviceName, tpe, enum...)
	}
}

func RegisterEnum(serviceName string, tpe string, options ...enumeration.EnumOption) {
	serviceName = strings.ToLower(codegen.ToUpperCamelCase(serviceName))
	if serviceEnumMap[serviceName] == nil {
		serviceEnumMap[serviceName] = map[string]enumeration.Enum{}
	}
	serviceEnumMap[serviceName][tpe] = options
}

func ToEnums(serviceName string, pkgName string) string {
	serviceName = strings.ToLower(codegen.ToUpperCamelCase(serviceName))
	buf := &bytes.Buffer{}
	imports := &bytes.Buffer{}

	names := make([]string, 0)

	for name := range serviceEnumMap[serviceName] {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		if name == "Bool" {
			continue
		}
		enum := serviceEnumMap[serviceName][name]
		buf.Write(ToEnumDefines(name, enum))
	}

	for _, name := range names {
		if name == "Bool" {
			continue
		}

		enum := serviceEnumMap[serviceName][name]
		e := gen.NewEnum("", pkgName, name, swagger_gen.Enum(enum), false)
		e.WriteAll(buf)
		e.Importer.WriteToImports(imports)
	}

	return `
package ` + pkgName + `

import (
	"errors"
	"bytes"
	"encoding"
	` + imports.String() + `
)

` + buf.String()
}

func ToEnumDefines(name string, enum enumeration.Enum) []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(`
// swagger:enum
type ` + name + ` uint

const (
`)

	buf.WriteString(codegen.ToUpperSnakeCase(name) + `_UNKNOWN ` + name + ` = iota
`)

	sort.Slice(enum, func(i, j int) bool {
		return enum[i].Val.(float64) < enum[j].Val.(float64)
	})

	index := 1
	for _, item := range enum {
		v := int(item.Val.(float64))
		if v > index {
			buf.WriteString(`)

const (
`)
			buf.WriteString(codegen.ToUpperSnakeCase(name) + `__` + item.Value.(string) + fmt.Sprintf(" %s = iota + %d", name, v) + `// ` + item.Label + `
`)
			index = v + 1
			continue
		}
		index++
		buf.WriteString(codegen.ToUpperSnakeCase(name) + `__` + item.Value.(string) + `// ` + item.Label + `
`)
	}

	buf.WriteString(`)`)

	return buf.Bytes()
}
