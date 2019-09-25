package gen

import (
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"strings"

	"github.com/morlay/oas"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

const (
	XNamed       = `x-go-named`
	XField       = `x-go-name`
	XTagJSON     = `x-go-json`
	XTagName     = `x-tag-name`
	XTagXML      = `x-tag-xml`
	XTagStyle    = `x-tag-style`
	XTagFmt      = `x-tag-fmt`
	XTagValidate = `x-go-validate`
	XPointer     = "x-pointer"
	XEnumValues  = `x-enum-values`
	XEnumLabels  = `x-enum-labels`
	XEnumVals    = `x-enum-vals`
)

func NewDefinitionScanner(program *loader.Program) *DefinitionScanner {
	return &DefinitionScanner{
		EnumScanner: NewEnumScanner(program),
		program:     program,
	}
}

type DefinitionScanner struct {
	EnumScanner *EnumScanner
	program     *loader.Program
	definitions map[*types.TypeName]*oas.Schema
}

func (scanner *DefinitionScanner) BindSchemas(openapi *oas.OpenAPI) {
	for typeName, schema := range scanner.definitions {
		schema.AddExtension(XNamed, fmt.Sprintf("%s.%s", typeName.Pkg().Path(), typeName.Name()))
		defKey := toDefID(typeName.Type().String())
		if _, exists := openapi.Components.Schemas[defKey]; exists {
			logrus.Panicf("`%s` already used by %s", defKey, typeName.String())
		} else {
			openapi.AddSchema(toDefID(typeName.Type().String()), schema)
		}
	}
	return
}

func (scanner *DefinitionScanner) getSchemaByTypeString(typeString string) *oas.Schema {
	pkgImportPath, _ := loaderx.GetPkgImportPathAndExpose(typeString)
	pkgInfo := scanner.program.Package(loaderx.ResolvePkgImport(pkgImportPath))
	for _, def := range pkgInfo.Defs {
		if typeName, ok := def.(*types.TypeName); ok {
			if typeName.Type().String() == typeString {
				return scanner.getSchemaByType(typeName.Type())
			}
		}
	}
	return nil
}

func (scanner *DefinitionScanner) Def(typeName *types.TypeName) *oas.Schema {
	if s, ok := scanner.definitions[typeName]; ok {
		return s
	}

	if typeName.IsAlias() {
		typeName = typeName.Type().(*types.Named).Obj()
	}

	doc := docOfTypeName(typeName.Type().(*types.Named).Obj(), scanner.program)

	if doc, fmtName := ParseStrfmt(doc); fmtName != "" {
		return scanner.addDef(typeName, oas.NewSchema(oas.TypeString, fmtName).WithDesc(doc))
	}

	// todo
	if typeName.Name() == "Time" {
		return scanner.addDef(typeName, oas.DateTime().WithDesc(doc))
	}

	doc, hasEnum := ParseEnum(doc)
	if hasEnum {
		enum := scanner.EnumScanner.Enum(typeName)
		if len(enum) == 2 {
			values := enum.Values()
			if values[0] == "FALSE" && values[1] == "TRUE" {
				return scanner.addDef(typeName, oas.Boolean())
			}
		}
		if enum == nil {
			panic(fmt.Errorf("missing enum option but annotated by swagger:enum"))
		}
		return scanner.addDef(typeName, enum.ToSchema().WithDesc(doc))
	}

	return scanner.addDef(typeName, scanner.getSchemaByType(typeName.Type().Underlying()).WithDesc(doc))
}

func (scanner *DefinitionScanner) addDef(typeName *types.TypeName, schema *oas.Schema) *oas.Schema {
	if scanner.definitions == nil {
		scanner.definitions = map[*types.TypeName]*oas.Schema{}
	}
	scanner.definitions[typeName] = schema
	return schema
}

func (scanner *DefinitionScanner) getSchemaByType(tpe types.Type) *oas.Schema {
	switch tpe.(type) {
	case *types.Interface:
		return &oas.Schema{
			SchemaObject: oas.SchemaObject{
				Type: oas.TypeObject,
				AdditionalProperties: &oas.SchemaOrBool{
					Allows: true,
				},
			},
		}
	case *types.Named:
		named := tpe.(*types.Named)
		if named.String() == "mime/multipart.FileHeader" {
			return oas.Binary()
		}
		scanner.Def(named.Obj())
		return oas.RefSchema(fmt.Sprintf("#/components/schemas/%s", toDefID(named.String())))
	case *types.Basic:
		typeName, format := getSchemaTypeFromBasicType(tpe.(*types.Basic).Name())
		if typeName != "" {
			return oas.NewSchema(typeName, format)
		}
	case *types.Pointer:
		count := 0
		pointer := tpe.(*types.Pointer)
		elem := pointer.Elem()
		for pointer != nil {
			elem = pointer.Elem()
			pointer, _ = pointer.Elem().(*types.Pointer)
			count++
		}
		s := scanner.getSchemaByType(elem)
		markPointer(s, count)
		return s
	case *types.Map:
		keySchema := scanner.getSchemaByType(tpe.(*types.Map).Key())
		if keySchema != nil && len(keySchema.Type) > 0 && keySchema.Type != "string" {
			panic(fmt.Errorf("only support map[string]interface{}"))
		}
		return oas.MapOf(scanner.getSchemaByType(tpe.(*types.Map).Elem()))
	case *types.Slice:
		return oas.ItemsOf(scanner.getSchemaByType(tpe.(*types.Slice).Elem()))
	case *types.Array:
		typArray := tpe.(*types.Array)
		length := typArray.Len()
		return oas.ItemsOf(scanner.getSchemaByType(typArray.Elem())).WithValidation(&oas.SchemaValidation{
			MaxItems: &length,
			MinItems: &length,
		})
	case *types.Struct:
		var structType = tpe.(*types.Struct)

		err := StructFieldUniqueChecker{}.Check(structType, false)
		if err != nil {
			panic(fmt.Errorf("type %s: %s", tpe, err))
		}

		var structSchema = oas.ObjectOf(nil)
		var schemas []*oas.Schema

		for i := 0; i < structType.NumFields(); i++ {
			field := structType.Field(i)

			if !field.Exported() {
				continue
			}

			structFieldType := field.Type()
			structFieldTags := reflect.StructTag(structType.Tag(i))
			jsonTagValue := structFieldTags.Get("json")

			name, flags := getTagNameAndFlags(jsonTagValue)
			if name == "-" {
				continue
			}

			if name == "" && field.Anonymous() {
				s := scanner.getSchemaByType(structFieldType)
				if s != nil {
					schemas = append(schemas, s)
				}
				continue
			}

			if name == "" {
				name = field.Name()
			}

			defaultValue, hasDefault := structFieldTags.Lookup("default")
			validate, hasValidate := structFieldTags.Lookup("validate")

			required := true
			if hasOmitempty, ok := flags["omitempty"]; ok {
				required = !hasOmitempty
			} else {
				// todo don't use non-default as required
				required = !hasDefault
			}

			propSchema := scanner.getSchemaByType(structFieldType)

			if flags != nil && flags["string"] {
				propSchema.Type = oas.TypeString
			}

			if defaultValue != "" {
				propSchema.Default = defaultValue
			}

			if hasValidate {
				BindValidateFromValidateTagString(propSchema, validate)
			}

			propSchema = propSchema.WithDesc(docOfTypeName(field, scanner.program))
			propSchema.AddExtension(XField, field.Name())

			if nameValue, hasName := structFieldTags.Lookup("name"); hasName {
				propSchema.AddExtension(XTagName, nameValue)
			}

			if styleValue, hasStyle := structFieldTags.Lookup("style"); hasStyle {
				propSchema.AddExtension(XTagStyle, styleValue)
			}

			if fmtValue, hasFmt := structFieldTags.Lookup("fmt"); hasFmt {
				propSchema.AddExtension(XTagFmt, fmtValue)
			}

			if xmlValue, hasXML := structFieldTags.Lookup("xml"); hasXML {
				propSchema.AddExtension(XTagXML, xmlValue)
			}

			if jsonTagValue != "" {
				propSchema.AddExtension(XTagJSON, jsonTagValue)
			}

			if propSchema.Ref != "" {
				composedSchema := oas.AllOf(
					propSchema,
					&oas.Schema{
						SchemaObject: propSchema.SchemaObject,
					},
				)
				composedSchema.SpecExtensions = propSchema.SpecExtensions
				structSchema.SetProperty(name, composedSchema, required)
			} else {
				structSchema.SetProperty(name, propSchema, required)
			}

		}

		if len(schemas) > 0 {
			return oas.AllOf(append(schemas, structSchema)...)
		}
		return structSchema
	}
	return nil
}

type StructFieldUniqueChecker map[string]*types.Var

func (checker StructFieldUniqueChecker) Check(structType *types.Struct, anonymous bool) error {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}
		if field.Anonymous() {
			if named, ok := field.Type().(*types.Named); ok {
				if st, ok := named.Underlying().(*types.Struct); ok {
					if err := checker.Check(st, true); err != nil {
						return err
					}
				}
			}
			continue
		}
		if anonymous {
			if _, ok := checker[field.Name()]; ok {
				return fmt.Errorf("%s.%s already defined in other anonymous field", structType.String(), field.Name())
			}
			checker[field.Name()] = field
		}
	}
	return nil
}

type VendorExtensible interface {
	AddExtension(key string, value interface{})
}

func markPointer(vendorExtensible VendorExtensible, count int) {
	vendorExtensible.AddExtension(XPointer, count)
}

func toDefID(s string) string {
	_, expose := loaderx.GetPkgImportPathAndExpose(s)
	return expose
}

var (
	rxEnum   = regexp.MustCompile(`swagger:enum`)
	rxStrFmt = regexp.MustCompile(`swagger:strfmt\s+(\S+)([\s\S]+)?$`)
)

func ParseEnum(doc string) (string, bool) {
	if rxEnum.MatchString(doc) {
		return strings.TrimSpace(strings.Replace(doc, "swagger:enum", "", -1)), true
	}
	return doc, false
}

func ParseStrfmt(doc string) (string, string) {
	matched := rxStrFmt.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

func getSchemaTypeFromBasicType(basicTypeName string) (tpe oas.Type, format string) {
	switch basicTypeName {
	case "bool":
		return "boolean", ""
	case "byte":
		return "integer", "uint8"
	case "error":
		return "string", ""
	case "float32":
		return "number", "float"
	case "float64":
		return "number", "double"
	case "int":
		return "integer", "int64"
	case "int8":
		return "integer", "int8"
	case "int16":
		return "integer", "int16"
	case "int32":
		return "integer", "int32"
	case "int64":
		return "integer", "int64"
	case "rune":
		return "integer", "int32"
	case "string":
		return "string", ""
	case "uint":
		return "integer", "uint64"
	case "uint16":
		return "integer", "uint16"
	case "uint32":
		return "integer", "uint32"
	case "uint64":
		return "integer", "uint64"
	case "uint8":
		return "integer", "uint8"
	case "uintptr":
		return "integer", "uint64"
	default:
		panic(fmt.Errorf("unsupported type %q", basicTypeName))
	}
}
