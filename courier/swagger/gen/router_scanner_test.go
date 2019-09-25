package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

func TestRouterScanner(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"context"
		"github.com/johnnyeven/libtools/courier/httpx"
		"github.com/johnnyeven/libtools/courier"
	)

	type Root struct {
		courier.EmptyOperator
	}

	type Group struct {
		courier.EmptyOperator
	}

	func (g Group) Path() string {
		return "/group"
	}

	type Auth struct{}

	func (auth Auth) Output(ctx context.Context) (result interface{}, err error) {
		return
	}

	type Get struct {
		httpx.MethodGet
	}

	func (get Get) Path() string {
		return "/id"
	}

	func (get Get) Output(ctx context.Context) (result interface{}, err error) {
		return
	}

	var Router = courier.NewRouter(Root{})

	func main() {
		group := courier.NewRouter(Group{})
		group.Register(courier.NewRouter(Get{}))

		Router.Register(group)
		Router.Register(courier.NewRouter(Auth{}, Get{}))
	}
	`)

	query := loaderx.NewQuery(program, pkgImportPath)
	scanner := NewRouterScanner(program)

	router := scanner.Router(query.Var("Router"))
	routes := router.Routes(program)

	tt.Len(routes, 2)
}
