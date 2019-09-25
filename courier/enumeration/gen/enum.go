package gen

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
)

func NewEnum(importPath, pkgName, name string, options gen.Enum, hasOffset bool) *Enum {
	return &Enum{
		HasOffset:  hasOffset,
		Name:       name,
		ImportPath: importPath,
		PkgName:    pkgName,
		Options:    options,
		Importer: &codegen.Importer{
			Local: importPath,
		},
	}
}

type Enum struct {
	ImportPath string
	PkgName    string
	Name       string
	Options    gen.Enum
	Importer   *codegen.Importer
	HasOffset  bool
}

func (m *Enum) ConstPrefix() string {
	return codegen.ToUpperSnakeCase(m.Name)
}

func (m *Enum) ConstOffset() string {
	return m.ConstPrefix() + "_OFFSET"
}

func (m *Enum) ConstUnknown() string {
	return m.ConstPrefix() + "_UNKNOWN"
}

func (m *Enum) InvalidError() string {
	return fmt.Sprintf("Invalid%s", m.Name)
}

func (m *Enum) ConstKey(key interface{}) string {
	return fmt.Sprintf("%s__%v", m.ConstPrefix(), key)
}

func (m *Enum) WriteAll(w io.Writer) {
	m.WriteVars(w)
	m.WriteInitFunc(w)
	m.WriteParseXFromString(w)
	m.WriteParseXFromLabelString(w)
	m.WriteEnumDescriptor(w)
	m.WriteStringer(w)
	m.WriteLabeler(w)
	m.WriteTextMarshalerAndUnmarshaler(w)
	m.WriteScannerAndValuer(w)
}

func (m *Enum) String() string {
	buf := &bytes.Buffer{}

	m.WriteAll(buf)

	return fmt.Sprintf(`
	package %s

	%s

	%s
	`,
		m.PkgName,
		m.Importer.String(),
		buf.String(),
	)
}

func (m *Enum) WriteVars(writer io.Writer) {
	io.WriteString(writer, `
var `+m.InvalidError()+` = errors.New("invalid `+m.Name+`")
`)
}

func (m *Enum) WriteInitFunc(w io.Writer) {
	io.WriteString(w, `
func init () {
	`+m.Importer.Use("github.com/johnnyeven/libtools/courier/enumeration.RegisterEnums")+`("`+m.Name+`", map[string]string{
`)

	for _, option := range enumeration.Enum(m.Options) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%v", option.Value))+":"+strconv.Quote(option.Label)+`,
`)
	}

	io.WriteString(w, `
	})
}
`)
}

func (m *Enum) WriteParseXFromString(w io.Writer) {
	io.WriteString(w, `
func `+fmt.Sprintf("Parse%sFromString", m.Name)+`(s string) (`+m.Name+`, error) {
	switch s {`)

	io.WriteString(w, `
		case `+strconv.Quote("")+`:
			return `+m.ConstUnknown()+`, nil`)

	for _, option := range m.Options {
		io.WriteString(w, `
		case `+fmt.Sprintf(`"%v"`, option.Value)+`:
			return `+m.ConstKey(option.Value)+`, nil`)
	}

	io.WriteString(w, `
	}
	return `+m.ConstUnknown()+`, `+m.InvalidError()+`
}
`)
}

func (m *Enum) WriteParseXFromLabelString(w io.Writer) {
	io.WriteString(w, `
func `+fmt.Sprintf("Parse%sFromLabelString", m.Name)+`(s string) (`+m.Name+`, error) {
	switch s {`)

	io.WriteString(w, `
		case `+strconv.Quote("")+`:
			return `+m.ConstUnknown()+`, nil`)

	for _, option := range m.Options {
		io.WriteString(w, `
		case `+strconv.Quote(option.Label)+`:
			return `+m.ConstKey(option.Value)+`, nil`)
	}

	io.WriteString(w, `
	}
	return `+m.ConstUnknown()+`,`+m.InvalidError()+`
}	
`)
}

func (m *Enum) WriteStringer(w io.Writer) {
	io.WriteString(w, `
func (v `+m.Name+`) String() string {
	switch v {`)

	io.WriteString(w, `
		case `+m.ConstUnknown()+`:
			return ""`)

	for _, option := range m.Options {
		io.WriteString(w, `
		case `+m.ConstKey(option.Value)+`:
			return `+fmt.Sprintf(`"%v"`, option.Value))
	}

	io.WriteString(w, `
	}
	return "UNKNOWN"
}
`)
}

func (m *Enum) WriteLabeler(w io.Writer) {
	io.WriteString(w, `
func (v `+m.Name+`) Label() string {
	switch v {`)

	io.WriteString(w, `
		case `+m.ConstUnknown()+`:
			return ""`)

	for _, option := range m.Options {
		io.WriteString(w, `
		case `+m.ConstKey(option.Value)+`:
			return `+strconv.Quote(option.Label))
	}

	io.WriteString(w, `
	}
	return "UNKNOWN"
}
	`)
}

func (m *Enum) WriteEnumDescriptor(w io.Writer) {
	io.WriteString(w, `
func (`+m.Name+`) EnumType() string {
	return "`+m.Name+`"
}

func (`+m.Name+`) Enums() map[int][]string {
	return map[int][]string{
`)

	for _, option := range m.Options {
		io.WriteString(w, fmt.Sprintf(`int(%s): {%s, %s},
`, m.ConstKey(option.Value), strconv.Quote(fmt.Sprintf("%s", option.Value)), strconv.Quote(option.Label)))
	}

	io.WriteString(w, `
	}
}`)
}

func (m *Enum) WriteTextMarshalerAndUnmarshaler(w io.Writer) {
	io.WriteString(w, `
var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*`+m.Name+`)(nil)

func (v `+m.Name+`) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, `+m.InvalidError()+`
	}
	return []byte(str), nil
}

func (v *`+m.Name+`) UnmarshalText(data []byte) (err error) {
	*v, err = Parse`+m.Name+`FromString(string(bytes.ToUpper(data)))
	return
}`)
}

func (m *Enum) WriteScannerAndValuer(w io.Writer) {
	if !m.HasOffset {
		return
	}

	io.WriteString(w, `
var _ interface {
	sql.Scanner
	driver.Valuer
} = (*`+m.Name+`)(nil)

func (v *`+m.Name+`) Scan(src interface{}) error {
	integer, err := `+m.Importer.Use("github.com/johnnyeven/libtools/courier/enumeration.AsInt64")+`(src, `+m.ConstOffset()+`)
	if err != nil {
		return err
	}
	*v = `+m.Name+`(integer - `+m.ConstOffset()+`)
	return nil
}

func (v `+m.Name+`) Value() (driver.Value, error) {
	return int64(v) + `+m.ConstOffset()+`, nil
}
`)
}
