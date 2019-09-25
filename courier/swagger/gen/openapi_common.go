package gen

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/morlay/oas"
)

var (
	rxValidate = regexp.MustCompile(`^@([^\[\(\{]+)([\[\(\{])([^\}^\]^\)]+)([\}\]\)])$`)
)

func BindValidateFromValidateTagString(schema *oas.Schema, validate string) {
	commonValidations := getCommonValidations(validate)
	bindSchemaWithValidations(schema, commonValidations)
	schema.AddExtension(XTagValidate, validate)
}

func getCommonValidations(validateTag string) (commonValidations oas.SchemaValidation) {
	var matched = rxValidate.FindAllStringSubmatch(validateTag, -1)

	// "@int[1,2]", "int", "[", "1,2", "]"
	if len(matched) > 0 && len(matched[0]) == 5 {
		tpe := matched[0][1]
		startBracket := matched[0][2]
		endBracket := matched[0][4]
		values := strings.Split(matched[0][3], ",")

		if startBracket != "{" && endBracket != "}" {
			switch tpe {
			case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
				if len(values) > 0 {
					if val, err := strconv.ParseFloat(values[0], 64); err == nil {
						commonValidations.Minimum = &val
						commonValidations.ExclusiveMinimum = (startBracket == "(")
					}
				}
				if len(values) > 1 {
					if val, err := strconv.ParseFloat(values[1], 64); err == nil {
						commonValidations.Maximum = &val
						commonValidations.ExclusiveMaximum = (endBracket == ")")
					}
				}
			case "string":
				if len(values) > 0 {
					if val, err := strconv.ParseInt(values[0], 10, 64); err == nil {
						commonValidations.MinLength = &val
					}
				}
				if len(values) > 1 {
					if val, err := strconv.ParseInt(values[1], 10, 64); err == nil {
						commonValidations.MaxLength = &val
					}
				}
			}
		} else {
			enums := make([]interface{}, 0)

			for _, value := range values {
				if tpe != "string" {
					if val, err := strconv.ParseInt(value, 10, 64); err == nil {
						enums = append(enums, val)
					}
				} else {
					enums = append(enums, value)
				}
			}

			commonValidations.Enum = enums
		}
	}

	return
}

func bindSchemaWithValidations(schema *oas.Schema, schemaValidation oas.SchemaValidation) {
	schema.SchemaValidation = schemaValidation

	// for partial enum
	if len(schema.Enum) != 0 && len(schemaValidation.Enum) != 0 {
		var enums []interface{}

		for _, enumValueOrIndex := range schemaValidation.Enum {
			switch reflect.TypeOf(enumValueOrIndex).Name() {
			case "string":
				if enumContainsValue(schema.Enum, enumValueOrIndex) {
					enums = append(enums, enumValueOrIndex)
				} else if enumValueOrIndex != "" {
					panic(fmt.Errorf("%s is not value of %s", enumValueOrIndex, schema.Enum))
				}
			default:
				if idx, ok := enumValueOrIndex.(int); ok {
					if schema.Enum[idx] != nil {
						enums = append(enums, schema.Enum[idx])
					} else if idx != 0 {
						panic(fmt.Errorf("%s is out-range of  %s", enumValueOrIndex, schema.Enum))
					}
				}

			}
		}

		schema.Enum = enums
	}
}

func enumContainsValue(enum []interface{}, value interface{}) bool {
	var isContains = false

	for _, enumValue := range enum {
		if enumValue == value {
			isContains = true
		}
	}

	return isContains
}
