package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

func TestOperatorScanner(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"fmt"
		"net/http"
		"github.com/johnnyeven/libtools/courier/status_error"
	)

	type With struct {
		Query string  ^^name:"query" in:"query"^^
	}

	// swagger:enum
	type Status uint8

	const (
		STATUS_UNKNOWN Status = iota
		STATUS__A    // A
		STATUS__B    // B
	)

	// Op
	//go:generate echo
	type Op struct {
		With
		Param string  ^^name:"param" in:"path"^^
		Status Status  ^^name:"status,omitempty" in:"query" validate:"@string{A}"^^
		FormData struct {
			Name string ^^name:"status"^^
		}  ^^in:"formData"^^
	}

	func (op Op) Output() (resp interface{}, err error) {
		if op.Query == "" {
			return nil, status_error.InvalidStruct
		}
		resp = &Resp{
			Name: "s",
		}

		call := func() int {
			return 1
		}
		fmt.Println(call())

		return
	}

	type Resp struct {
		Name string ^^json:"name"^^
	}

	func (resp Resp) ContentType() string {
		return "application/json"
	}

	func (resp Resp) Status() int {
		return http.StatusOK
	}
	`)

	scanner := NewOperatorScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)

	operator := scanner.Operator(query.TypeName("Op"))

	tt.Equal("Op", operator.Summary)

	tt.Equal(ToMap(map[string]map[string]interface{}{
		"query": {
			"name":     "query",
			"in":       "query",
			"required": true,
			XField:     "Query",
			"schema": map[string]interface{}{
				"type": "string",
			},
		},
		"status": {
			"name": "status",
			"in":   "query",
			"schema": map[string]interface{}{
				"allOf": []map[string]interface{}{
					{
						"$ref": "#/components/schemas/Status",
					},
					{
						"enum":       []interface{}{"A"},
						XTagValidate: "@string{A}",
					},
				},
			},
			XField: "Status",
		},
		"param": {
			"name":     "param",
			"in":       "path",
			"required": true,
			XField:     "Param",
			"schema": map[string]interface{}{
				"type": "string",
			},
		},
	}), ToMap(operator.NonBodyParameters))
}

func TestOperatorScannerWithFile(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"github.com/johnnyeven/libtools/courier/transport_http"
	)

	// Op
	type Op struct {
	}

	// @success content-type text/plain
	func (op Op) Output() (resp interface{}, err error) {
		file := transport_http.NewFile("1.txt", "text/plain")
		resp = file
		return
	}

	`)

	scanner := NewOperatorScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)

	operator := scanner.Operator(query.TypeName("Op"))

	_, ok := operator.SuccessResponse.WithContent.Content["text/plain"]
	tt.True(ok)
}

func TestOperatorScannerForWebSocket(t *testing.T) {
	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"github.com/johnnyeven/libtools/courier/transport_http"
	)

	type WS struct {
	}

	func (ws WS) Output() (resp interface{}, err error) {
		listeners := transport_http.Listeners{}
		listeners.On(Ping{}, func(v interface{}, ws *transport_http.WSClient) error {
			p := v.(*Ping)
			return ws.Send(Pong{
				MS: p.MS + 1,
			})
		})
		resp = listeners
		return
	}

	type Ping struct {
		MS int ^^json:"ms"^^
	}

	type Pong struct {
		MS int ^^json:"ms"^^
	}
	`)

	scanner := NewOperatorScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)

	scanner.Operator(query.TypeName("WS"))
}
