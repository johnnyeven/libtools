package codegen

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type ImportPkg struct {
	*build.Package
	Alias string
}

func (importPkg *ImportPkg) GetID() string {
	if importPkg.Alias != "" {
		return importPkg.Alias
	}
	return importPkg.Name
}

func (importPkg *ImportPkg) String() string {
	if importPkg.Alias != "" {
		return importPkg.Alias + " " + strconv.Quote(importPkg.ImportPath) + "\n"
	}
	return strconv.Quote(importPkg.ImportPath)
}

type Importer struct {
	Local string
	pkgs  map[string]*ImportPkg
}

func getPkgImportPathAndExpose(s string) (pkgImportPath string, expose string) {
	idxSlash := strings.LastIndex(s, "/")
	idxDot := strings.LastIndex(s, ".")
	if idxDot > idxSlash {
		return s[0:idxDot], s[idxDot+1:]
	}
	return s, ""
}

func (importer *Importer) ExposeVar(name string) string {
	return ToUpperCamelCase(name)
}

func (importer *Importer) Var(name string) string {
	return ToLowerCamelCase(name)
}

func (importer *Importer) PureUse(importPath string, subPkgs ...string) string {
	pkgPath, expose := getPkgImportPathAndExpose(strings.Join(append([]string{importPath}, subPkgs...), "/"))

	importPkg := importer.Import(pkgPath, false)

	if expose != "" {
		if pkgPath == importer.Local {
			return expose
		}
		return fmt.Sprintf("%s.%s", importPkg.GetID(), expose)
	}

	return importPkg.GetID()
}

// use and alias
func (importer *Importer) Use(importPath string, subPkgs ...string) string {
	pkgPath, expose := getPkgImportPathAndExpose(strings.Join(append([]string{importPath}, subPkgs...), "/"))

	importPkg := importer.Import(pkgPath, true)

	if expose != "" {
		if pkgPath == importer.Local {
			return expose
		}
		return fmt.Sprintf("%s.%s", importPkg.GetID(), expose)
	}

	return importPkg.GetID()
}

func (importer *Importer) Import(importPath string, alias bool) *ImportPkg {
	importPath = DeVendor(importPath)
	if importer.pkgs == nil {
		importer.pkgs = map[string]*ImportPkg{}
	}

	importPkg, exists := importer.pkgs[importPath]
	if !exists {
		pkg, err := build.Import(importPath, "", build.ImportComment)
		if err != nil {
			panic(err)
		}
		importPkg = &ImportPkg{
			Package: pkg,
		}
		if alias {
			importPkg.Alias = ToLowerSnakeCase(importPath)
		}
		importer.pkgs[importPath] = importPkg
	}

	return importPkg
}

func DeVendor(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func (importer *Importer) WriteToImports(w io.Writer) {
	if len(importer.pkgs) > 0 {
		for _, importPkg := range importer.pkgs {
			io.WriteString(w, importPkg.String()+"\n")
		}
	}
}

func (importer *Importer) String() string {
	buf := &bytes.Buffer{}
	if len(importer.pkgs) > 0 {
		buf.WriteString("import (\n")
		importer.WriteToImports(buf)
		buf.WriteString(")")
	}
	return buf.String()
}

func (importer *Importer) Type(tpe reflect.Type) string {
	if tpe.PkgPath() != "" {
		return importer.Use(fmt.Sprintf("%s.%s", tpe.PkgPath(), tpe.Name()))
	}

	switch tpe.Kind() {
	case reflect.Slice:
		return fmt.Sprintf("[]%s", importer.Type(tpe.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", importer.Type(tpe.Key()), importer.Type(tpe.Elem()))
	default:
		return tpe.String()
	}
}

func (importer *Importer) Sdump(v interface{}) string {
	tpe := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		parts := make([]string, 0)
		isMulti := rv.Len() > 1
		for _, key := range rv.MapKeys() {
			s := importer.Sdump(key.Interface()) + ": " + importer.Sdump(rv.MapIndex(key).Interface())
			if isMulti {
				parts = append(parts, s+",\n")
			} else {
				parts = append(parts, s)
			}
		}
		sort.Strings(parts)

		if isMulti {
			f := `%s{
				%s
			}`
			return fmt.Sprintf(f, importer.Type(tpe), strings.Join(parts, ""))
		}
		f := "%s{%s}"
		return fmt.Sprintf(f, importer.Type(tpe), strings.Join(parts, ", "))
	case reflect.Slice:
		buf := new(bytes.Buffer)
		for i := 0; i < rv.Len(); i++ {
			s := importer.Sdump(rv.Index(i).Interface())
			if i == 0 {
				buf.WriteString(s)
			} else {
				buf.WriteString(", " + s)
			}
		}
		return fmt.Sprintf("%s{%s}", importer.Type(tpe), buf.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%d", v)
	case reflect.Bool:
		return strconv.FormatBool(v.(bool))
	case reflect.Float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case reflect.Invalid:
		return "nil"
	case reflect.String:
		return strconv.Quote(v.(string))
	}
	return ""
}
