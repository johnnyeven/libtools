package codegen

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

func NewGenFile(pkgName string, filename string) *GenFile {
	return &GenFile{
		PkgName:  pkgName,
		Filename: filename,
		Importer: &Importer{},
		Buffer:   &bytes.Buffer{},
	}
}

type GenFile struct {
	Filename string
	PkgName  string
	Data     interface{}
	*bytes.Buffer
	*Importer
}

func (f *GenFile) WithData(data interface{}) *GenFile {
	f.Data = data
	return f
}

func (f *GenFile) Block(tpl string) *GenFile {
	f.writeTo(f.Buffer, tpl)
	return f
}

func (f *GenFile) writeTo(writer io.Writer, tpl string) {
	t, parseErr := template.New(f.Filename).Parse(tpl)
	if parseErr != nil {
		panic(fmt.Sprintf("template Prase failded: %s", parseErr.Error()))
	}
	err := t.Execute(writer, f)
	if err != nil {
		panic(fmt.Sprintf("template Execute failded: %s", err.Error()))
	}
}

func (f *GenFile) String() string {
	return fmt.Sprintf(`
package %s
%s
%s
`,
		f.PkgName,
		f.Importer.String(),
		f.Buffer.String(),
	)
}

func (f *GenFile) OutputTo(outputs Outputs) {
	outputs.Add(f.Filename, f.String())
}
