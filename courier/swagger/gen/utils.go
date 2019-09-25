package gen

import (
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier"
)

func getTagNameAndFlags(tagValue string) (name string, flags map[string]bool) {
	values := strings.SplitN(tagValue, ",", -1)
	for _, flag := range values[1:] {
		if flags == nil {
			flags = map[string]bool{}
		}
		flags[flag] = true
	}
	name = values[0]
	return
}

func docOfTypeName(obj types.Object, program *loader.Program) string {
	pkgInfo := program.Package(obj.Pkg().Path())
	for ident, def := range pkgInfo.Defs {
		if def == obj {
			return loaderx.CommentsOf(program.Fset, ident, pkgInfo.Files...)
		}
	}
	return ""
}

func CommentValuesSplit(s string) (values []string) {
	parts := strings.Split(s, " ")
	for _, p := range parts {
		if p != "" {
			values = append(values, p)
		}
	}
	return
}

func ParseSuccessMetadata(doc string) courier.Metadata {
	metadata := courier.Metadata{}

	matches := regexp.MustCompile(`@success ([^\n]+)`).FindAllStringSubmatch(doc, -1)

	for _, subMatch := range matches {
		if len(subMatch) == 2 {
			parts := CommentValuesSplit(subMatch[1])
			metadata.Set(strings.TrimSpace(parts[0]), parts[1:]...)
		}
	}

	return metadata
}

func reflectTypeString(tpe reflect.Type) string {
	for tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}
	return fmt.Sprintf("%s.%s", loaderx.ResolvePkgImport(tpe.PkgPath()), tpe.Name())
}
