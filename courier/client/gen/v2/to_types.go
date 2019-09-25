package v2

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/client/gen/common"
	"github.com/johnnyeven/libtools/courier/client/gen/enums"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
	"github.com/johnnyeven/libtools/godash"
)

func ToTypes(serviceName string, pkgName string, swagger *spec.Swagger) string {
	buf := &bytes.Buffer{}

	importer := &codegen.Importer{}

	definitionNames := make([]string, 0)
	for name := range swagger.Definitions {
		definitionNames = append(definitionNames, name)
	}
	sort.Strings(definitionNames)

	for _, name := range definitionNames {
		s := swagger.Definitions[name]
		buf.WriteString(`
type ` + name + ` ` + NewTypeGenerator(serviceName, importer).Type(&s) + `
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

func (g *TypeGenerator) Type(schema *spec.Schema) string {
	pointer := ""
	if schema.Extensions[gen.XPointer] != nil {
		pointer = strings.Repeat("*", schema.Extensions[gen.XPointer].(int))
	}
	return pointer + g.TypeIndirect(schema)
}

func (g *TypeGenerator) TypeIndirect(schema *spec.Schema) string {
	if schema == nil {
		return "interface{}"
	}

	if schema.Ref.String() != "" {
		return common.RefName(schema.Ref.String())
	}

	if schema.Extensions[gen.XNamed] != nil {
		if schema.Type.Contains("array") && schema.Items != nil && schema.Items.Schema.Format == "uint64" {
			return g.Importer.Use("github.com/johnnyeven/libtools/httplib.Uint64List")
		}

		typeFullName := fmt.Sprint(schema.Extensions[gen.XNamed])

		if schema.Type.Contains("string") ||
			schema.Type.Contains("boolean") ||
			schema.Enum != nil ||
			strings.Contains(typeFullName, "johnnyeven/libtools") {

			pkgImportName, typeName := loaderx.GetPkgImportPathAndExpose(typeFullName)

			if schema.Extensions[gen.XEnumVals] != nil {
				typeName = g.PrefixType(typeName)
				enums.RegisterEnumFromExtensions(g.ServiceName, typeName, schema.Extensions[gen.XEnumVals], schema.Extensions[gen.XEnumValues], schema.Extensions[gen.XEnumLabels])
				return typeName
			}

			if schema.Type.Contains("boolean") {
				typeName = "Bool"
				pkgImportName = "github.com/johnnyeven/libtools/courier/enumeration"
			}

			return g.Importer.Use(fmt.Sprintf("%s.%s", pkgImportName, typeName))
		}
	}

	if len(schema.AllOf) > 0 {
		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		for _, subSchema := range schema.AllOf {
			if subSchema.Ref.String() != "" {
				field := common.NewField(common.RefName(subSchema.Ref.String()))
				buf.WriteString(field.String())
			}

			if subSchema.Properties != nil {
				g.WriteFields(buf, &subSchema)
			}
		}

		buf.WriteString(`}`)
		return buf.String()
	}

	if schema.Type.Contains("object") {
		if schema.AdditionalProperties != nil {
			return fmt.Sprintf("map[string]%s", g.Type(schema.AdditionalProperties.Schema))
		}

		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		g.WriteFields(buf, schema)

		buf.WriteString(`}`)
		return buf.String()
	}

	if schema.Type.Contains("array") {
		if schema.Items != nil && schema.Items.Schema != nil {
			return fmt.Sprintf("[]%s", g.Type(schema.Items.Schema))
		}
	}

	schemaType := "string"
	if len(schema.Type) > 0 {
		schemaType = schema.Type[0]
	}

	return common.BasicType(schemaType, schema.Format, g.Importer)
}

func (g *TypeGenerator) WriteFields(w io.Writer, schema *spec.Schema) {
	if schema.Properties == nil {
		return
	}

	fieldNames := make([]string, 0)
	for fieldName := range schema.Properties {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	for _, fieldName := range fieldNames {
		propSchema := schema.Properties[fieldName]
		io.WriteString(w, g.FieldFrom(fieldName, &propSchema, schema.Required...).String())
	}
}

func (g *TypeGenerator) FieldFrom(name string, propSchema *spec.Schema, requiredFields ...string) *common.Field {
	isRequired := godash.StringIncludes(requiredFields, name)

	fieldName := name
	if propSchema.Extensions[gen.XField] != nil {
		fieldName = propSchema.Extensions[gen.XField].(string)
	}

	field := common.NewField(fieldName)
	field.Comment = propSchema.Description

	field.AddTag("json", name)

	if propSchema.Type.Contains("string") && propSchema.Format == "uint64" {
		field.AddTag("json", name+",string")
	}

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

	if propSchema.Extensions[gen.XTagValidate] != nil {
		field.AddTag("validate", fmt.Sprintf("%s", propSchema.Extensions[gen.XTagValidate]))
	}

	if propSchema.Default != nil {
		field.AddTag("default", fmt.Sprintf("%v", propSchema.Default))
	}

	field.Type = g.Type(propSchema)

	return field
}
