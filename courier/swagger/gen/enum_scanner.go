package gen

import (
	"go/ast"
	"go/constant"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/morlay/oas"
	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/enumeration"
)

func NewEnumScanner(program *loader.Program) *EnumScanner {
	return &EnumScanner{
		program: program,
	}
}

type EnumScanner struct {
	program *loader.Program
	Enums   map[*types.TypeName]Enum
}

func (scanner *EnumScanner) HasOffset(typeName *types.TypeName) bool {
	pkgInfo := scanner.program.Package(typeName.Pkg().Path())
	if pkgInfo == nil {
		return false
	}
	for _, def := range pkgInfo.Defs {
		if typeConst, ok := def.(*types.Const); ok {
			if typeConst.Name() == codegen.ToUpperSnakeCase(typeName.Name())+"_OFFSET" {
				return true
			}
		}
	}
	return false
}

func (scanner *EnumScanner) Enum(typeName *types.TypeName) Enum {
	if enumOptions, ok := scanner.Enums[typeName]; ok {
		return enumOptions.Sort()
	}

	pkgInfo := scanner.program.Package(typeName.Pkg().Path())
	if pkgInfo == nil {
		return nil
	}

	typeNameString := typeName.Name()

	for ident, def := range pkgInfo.Defs {
		if typeConst, ok := def.(*types.Const); ok {
			if typeConst.Type() == typeName.Type() {
				name := typeConst.Name()

				if name != "_" {
					val := typeConst.Val()
					label := strings.TrimSpace(ident.Obj.Decl.(*ast.ValueSpec).Comment.Text())

					if strings.HasPrefix(name, codegen.ToUpperSnakeCase(typeNameString)) {
						var values = strings.SplitN(name, "__", 2)
						if len(values) == 2 {
							scanner.addEnum(typeName, values[1], getConstVal(val), label)
						}
					} else {
						v := getConstVal(val)
						scanner.addEnum(typeName, v, v, label)
					}
				}
			}
		}
	}

	return scanner.Enums[typeName].Sort()
}

func (scanner *EnumScanner) addEnum(typeName *types.TypeName, value interface{}, val interface{}, label string) {
	if scanner.Enums == nil {
		scanner.Enums = map[*types.TypeName]Enum{}
	}
	scanner.Enums[typeName] = append(scanner.Enums[typeName], enumeration.EnumOption{
		Value: value,
		Val:   val,
		Label: label,
	})
}

type Enum enumeration.Enum

func (enum Enum) Sort() Enum {
	sort.Slice(enum, func(i, j int) bool {
		switch enum[i].Value.(type) {
		case string:
			return enum[i].Value.(string) < enum[j].Value.(string)
		case int64:
			return enum[i].Value.(int64) < enum[j].Value.(int64)
		case float64:
			return enum[i].Value.(float64) < enum[j].Value.(float64)
		}
		return true
	})
	return enum
}

func (enum Enum) Labels() (labels []string) {
	for _, e := range enum {
		labels = append(labels, e.Label)
	}
	return
}

func (enum Enum) Vals() (vals []interface{}) {
	for _, e := range enum {
		vals = append(vals, e.Val)
	}
	return
}

func (enum Enum) Values() (values []interface{}) {
	for _, e := range enum {
		values = append(values, e.Value)
	}
	return
}

func (enum Enum) ToSchema() *oas.Schema {
	values := enum.Values()

	// nullable bool
	if len(enum) == 2 && reflect.DeepEqual(values, []string{"FALSE", "TRUE"}) {
		return oas.Boolean()
	}

	typeName, _ := getSchemaTypeFromBasicType(reflect.TypeOf(values[0]).Name())

	s := oas.NewSchema(typeName, "").WithValidation(&oas.SchemaValidation{
		Enum: values,
	})
	s.AddExtension(XEnumLabels, enum.Labels())
	s.AddExtension(XEnumVals, enum.Vals())
	s.AddExtension(XEnumValues, values)
	return s
}

func getConstVal(constVal constant.Value) interface{} {
	switch constVal.Kind() {
	case constant.String:
		stringVal, _ := strconv.Unquote(constVal.String())
		return stringVal
	case constant.Int:
		intVal, _ := strconv.ParseInt(constVal.String(), 10, 64)
		return intVal
	case constant.Float:
		floatVal, _ := strconv.ParseFloat(constVal.String(), 10)
		return floatVal
	}
	return nil
}
