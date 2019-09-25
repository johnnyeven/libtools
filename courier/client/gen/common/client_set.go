package common

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/johnnyeven/libtools/codegen"
)

func NewClientSet(baseClient, name string) *ClientSet {
	return &ClientSet{
		BaseClient: baseClient,
		Name:       codegen.ToLowerLinkCase(name),
		PkgName:    codegen.ToLowerSnakeCase("Client-" + name),
		ClientName: codegen.ToUpperCamelCase("Client-" + name),
		Ops:        map[string]Op{},
		Importer:   &codegen.Importer{},
	}
}

type ClientSet struct {
	BaseClient string
	Name       string
	PkgName    string
	ClientName string
	Ops        map[string]Op
	Importer   *codegen.Importer
}

func (c *ClientSet) AddOp(op Op) {
	if op != nil {
		c.Ops[op.ID()] = op
	}
}

func (c *ClientSet) WriteAll(w io.Writer) {
	c.WriteTypeInterface(w)
	c.WriteTypeClient(w)
	c.WriteOperations(w)
}
func (c *ClientSet) WriteTypeInterface(w io.Writer) {
	keys := make([]string, 0)
	for key := range c.Ops {
		if key == "Swagger" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	io.WriteString(w, `type `+c.ClientName+`Interface interface {
`)

	for _, key := range keys {
		op := c.Ops[key]

		reqVar := ""
		reqType := ""
		reqTypeInParams := ""

		if op.HasRequest() {
			reqVar = "req"
			reqType = RequestOf(op.ID())
			reqTypeInParams = reqType + ", "
		}

		interfaceMethod := op.ID() + `(` + reqVar + ` ` + reqTypeInParams + `metas... ` + c.Importer.Use("github.com/johnnyeven/libtools/courier.Metadata") + `) (resp *` + ResponseOf(op.ID()) + `, err error)
`

		io.WriteString(w, interfaceMethod)
	}

	io.WriteString(w, `}
`)
}

func (c *ClientSet) WriteOperations(w io.Writer) {
	keys := make([]string, 0)
	for key := range c.Ops {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		op := c.Ops[key]

		reqVar := ""
		reqType := ""
		reqTypeInParams := ""
		reqVarInUse := "nil"

		if op.HasRequest() {
			reqVar = "req"
			reqVarInUse = reqVar
			reqType = RequestOf(op.ID())
			reqTypeInParams = reqType + ", "

			io.WriteString(w, `
type `+reqType+" ")
			op.WriteReqType(w, c.Importer)
		}

		interfaceMethod := op.ID() + `(` + reqVar + ` ` + reqTypeInParams + `metas... ` + c.Importer.Use("github.com/johnnyeven/libtools/courier.Metadata") + `) (resp *` + ResponseOf(op.ID()) + `, err error)`

		io.WriteString(w, `
func (c `+c.ClientName+`) `+interfaceMethod+` {
	resp = &`+ResponseOf(op.ID())+`{}
	resp.Meta = `+c.Importer.Use("github.com/johnnyeven/libtools/courier.Metadata")+`{}

	err = c.Request(c.Name + ".`+op.ID()+`", "`+op.Method()+`", "`+op.Path()+`", `+reqVarInUse+`, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}
`)

		io.WriteString(w, `
type `+ResponseOf(op.ID())+`  struct {
	Meta `+c.Importer.Use("github.com/johnnyeven/libtools/courier.Metadata")+`
	Body `)
		op.WriteRespBodyType(w, c.Importer)
		io.WriteString(w, `
}
`)

	}
}

func (c *ClientSet) WriteTypeClient(w io.Writer) {
	io.WriteString(w, `
type `+c.ClientName+` struct {
	`+c.Importer.Use(c.BaseClient)+`
}

func (`+c.ClientName+`) MarshalDefaults(v interface{}) {
	if cl, ok := v.(* `+c.ClientName+`); ok {
		cl.Name = "`+c.Name+`"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c  `+c.ClientName+`) Init() {
	c.CheckService()
}

func (c  `+c.ClientName+`) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := `+c.Importer.Use("github.com/johnnyeven/libtools/courier/status_error.FromError")+`(err)
	if statusErr.Code == int64(`+c.Importer.Use("github.com/johnnyeven/libtools/courier/status_error.RequestTimeout")+`) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}
`)
}

func (c *ClientSet) String() string {
	buf := &bytes.Buffer{}

	c.WriteAll(buf)

	return fmt.Sprintf(`
	package %s

	%s

	%s
	`,
		c.PkgName,
		c.Importer.String(),
		buf.String(),
	)
}
