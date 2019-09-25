package v3

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/morlay/oas"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/client/gen/common"
	"github.com/johnnyeven/libtools/courier/client/gen/enums"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
	"github.com/johnnyeven/libtools/godash"
)

func ToTypes(serviceName string, pkgName string, openAPI *oas.OpenAPI) string {
	buf := &bytes.Buffer{}

	importer := &codegen.Importer{}

	definitionNames := make([]string, 0)
	for name := range openAPI.Components.Schemas {
		definitionNames = append(definitionNames, name)
	}
	sort.Strings(definitionNames)

	for _, name := range definitionNames {
		s := openAPI.Components.Schemas[name]
		tpe, alias := NewTypeGenerator(serviceName, importer).Type(s)
		op := " "
		if alias {
			op = " = "
		}
		buf.WriteString(`
type ` + name + op + tpe + `
`)
	}

	return fmt.Sprintf(`
	package %s

	%s

	%s
	`,
		pkgName,
		importer.String(),
		buf.String(),
	)
}

func NewTypeGenerator(serviceName string, importer *codegen.Importer) *TypeGenerator {
	return &TypeGenerator{
		ServiceName: serviceName,
		Importer:    importer,
	}
}

type TypeGenerator struct {
	ServiceName string
	Importer    *codegen.Importer
}

func (g *TypeGenerator) PrefixType(tpe string) string {
	return codegen.ToUpperCamelCase(g.ServiceName) + tpe
}

func (g *TypeGenerator) Type(schema *oas.Schema) (string, bool) {
	pointer := ""
	if schema.Extensions[gen.XPointer] != nil {
		pointer = strings.Repeat("*", int(schema.Extensions[gen.XPointer].(float64)))
	}
	tpe, alias := g.TypeIndirect(schema)
	return pointer + tpe, alias
}

func (g *TypeGenerator) TypeIndirect(schema *oas.Schema) (string, bool) {
	if schema == nil {
		return "interface{}", false
	}

	if schema.Ref != "" {
		return common.RefName(schema.Ref), true
	}

	if schema.Extensions[gen.XNamed] != nil {
		if schema.Type == "array" && schema.Items != nil && schema.Items.Format == "uint64" {
			return g.Importer.Use("github.com/johnnyeven/libtools/httplib.Uint64List"), true
		}

		typeFullName := fmt.Sprint(schema.Extensions[gen.XNamed])
		isInCommonLib := strings.Contains(typeFullName, "johnnyeven/libtools")

		if schema.Type == "string" || schema.Type == "boolean" || schema.Enum != nil || isInCommonLib {

			pkgImportName, typeName := loaderx.GetPkgImportPathAndExpose(typeFullName)

			if schema.Extensions[gen.XEnumVals] != nil {
				typeName = g.PrefixType(typeName)
				enums.RegisterEnumFromExtensions(g.ServiceName, typeName, schema.Extensions[gen.XEnumVals], schema.Extensions[gen.XEnumValues], schema.Extensions[gen.XEnumLabels])
				return typeName, true
			}

			if schema.Type == "boolean" {
				typeName = "Bool"
				pkgImportName = "github.com/johnnyeven/libtools/courier/enumeration"
				isInCommonLib = true
			}

			return g.Importer.Use(fmt.Sprintf("%s.%s", pkgImportName, typeName)), true
		}
	}

	if len(schema.AllOf) > 0 {
		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		for _, subSchema := range schema.AllOf {
			if subSchema.Ref != "" {
				field := common.NewField(common.RefName(subSchema.Ref))
				buf.WriteString(field.String())
			}

			if subSchema.Properties != nil {
				g.WriteFields(buf, subSchema)
			}
		}

		buf.WriteString(`}`)
		return buf.String(), false
	}

	if schema.Type == "object" {
		if schema.AdditionalProperties != nil {
			tpe, _ := g.Type(schema.AdditionalProperties.Schema)
			return fmt.Sprintf("map[string]%s", tpe), false
		}

		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		g.WriteFields(buf, schema)

		buf.WriteString(`}`)
		return buf.String(), false
	}

	if schema.Type == "array" {
		if schema.Items != nil {
			tpe, _ := g.Type(schema.Items)
			return fmt.Sprintf("[]%s", tpe), false
		}
	}

	return common.BasicType(string(schema.Type), schema.Format, g.Importer), false
}

func (g *TypeGenerator) WriteFields(w io.Writer, schema *oas.Schema) {
	if schema.Properties == nil {
		return
	}

	fieldNames := make([]string, 0)
	for fieldName := range schema.Properties {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	for _, fieldName := range fieldNames {
		propSchema := mayComposedFieldSchema(schema.Properties[fieldName])

		io.WriteString(w, g.FieldFrom(fieldName, propSchema, schema.Required...).String())
	}
}

func (g *TypeGenerator) FieldFrom(name string, propSchema *oas.Schema, requiredFields ...string) *common.Field {
	isRequired := godash.StringIncludes(requiredFields, name)

	fieldName := name
	if propSchema.Extensions[gen.XField] != nil {
		fieldName = propSchema.Extensions[gen.XField].(string)
	}

	field := common.NewField(fieldName)
	field.Comment = propSchema.Description

	if propSchema.Extensions[gen.XTagJSON] != nil {
		tagName := fmt.Sprintf("%s", propSchema.Extensions[gen.XTagJSON])
		flags := make([]string, 0)
		if !isRequired && !strings.Contains(tagName, "omitempty") {
			flags = append(flags, "omitempty")
		}
		field.AddTag("json", tagName, flags...)
	}

	if propSchema.Extensions[gen.XTagName] != nil {
		tagName := fmt.Sprintf("%s", propSchema.Extensions[gen.XTagName])
		flags := make([]string, 0)
		if !isRequired && !strings.Contains(tagName, "omitempty") {
			flags = append(flags, "omitempty")
		}
		field.AddTag("name", tagName, flags...)
	}

	if propSchema.Extensions[gen.XTagXML] != nil {
		field.AddTag("xml", fmt.Sprintf("%s", propSchema.Extensions[gen.XTagXML]))
	}

	if propSchema.Extensions[gen.XTagStyle] != nil {
		field.AddTag("style", fmt.Sprintf("%s", propSchema.Extensions[gen.XTagStyle]))
	}

	if propSchema.Extensions[gen.XTagFmt] != nil {
		field.AddTag("fmt", fmt.Sprintf("%s", propSchema.Extensions[gen.XTagFmt]))
	}

	if propSchema.Extensions[gen.XTagValidate] != nil {
		field.AddTag("validate", fmt.Sprintf("%s", propSchema.Extensions[gen.XTagValidate]))
	}

	if propSchema.Default != nil {
		field.AddTag("default", fmt.Sprintf("%v", propSchema.Default))
	}

	field.Type, _ = g.Type(propSchema)

	return field
}

func mayComposedFieldSchema(schema *oas.Schema) *oas.Schema {
	// for named field
	if schema.AllOf != nil && len(schema.AllOf) == 2 && schema.AllOf[len(schema.AllOf)-1].Type == "" {
		nextSchema := &oas.Schema{
			Reference:    schema.AllOf[0].Reference,
			SchemaObject: schema.AllOf[1].SchemaObject,
		}

		for k, v := range schema.AllOf[1].SpecExtensions.Extensions {
			nextSchema.AddExtension(k, v)
		}

		for k, v := range schema.SpecExtensions.Extensions {
			nextSchema.AddExtension(k, v)
		}

		return nextSchema
	}

	return schema
}
