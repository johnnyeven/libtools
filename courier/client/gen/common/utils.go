package common

import (
	"regexp"
	"strings"

	"github.com/johnnyeven/libtools/codegen"
)

func RequestOf(id string) string {
	return id + "Request"
}

func ResponseOf(id string) string {
	return id + "Response"
}

func PathFromSwaggerPath(str string) string {
	r := regexp.MustCompile(`/\{([^/\\}]+)\}`)
	result := r.ReplaceAllString(str, "/:$1")
	return result
}

func RefName(str string) string {
	parts := strings.Split(str, "/")
	return parts[len(parts)-1]
}

func BasicType(schemaType string, format string, importer *codegen.Importer) string {
	switch format {
	case "binary":
		return importer.Use("mime/multipart.FileHeader")
	case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		return format
	case "float":
		return "float32"
	case "double":
		return "float64"
	default:
		switch schemaType {
		case "boolean":
			return "bool"
		default:
			return "string"
		}
	}
}
