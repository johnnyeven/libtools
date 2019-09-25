package loaderx

import (
	"go/types"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachFuncResult(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := NewTestProgram(`
	package main

	type String string

	func fnA () String {
		return ""
	}

	func fn () (a interface{}, b String) {
		{
			a = nil
		}
		switch (a) {
		case "a3":
			a = "a3"
			b = "b3"
			return
		}
		if true {
			a = "a0"
			return fnA(), "b1"
		}
		b = "b2"
		return
	}

	func fn2 () (a interface{}, b String) {
		a, b = fn()
		return
	}

	func Fn () (a interface{}, b String) {
		return fn2()
	}
	`)

	q := NewQuery(program, pkgImportPath)
	typeFunc := q.Func("Fn")

	rets := [][]string{}
	ForEachFuncResult(program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
		ret := make([]string, len(resultTypeAndValues))
		for i, r := range resultTypeAndValues {
			if r.IsValue() {
				ret[i], _ = strconv.Unquote(r.Value.String())
			} else {
				ret[i] = r.Type.String()
			}
		}
		rets = append(rets, ret)
	})

	tt.Equal([][]string{
		{"a3", "b3"},
		{pkgImportPath + ".String", "b1"},
		{"untyped nil", "b2"},
	}, rets)
}

func TestMethodOf(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := NewTestProgram(`
	package main

	type SomeType struct {
		F string
	}

	func (a SomeType) Fn() {}

	type SomeTypeAlias = SomeType

	type SomeTypeCompose struct {
		SomeType
	}

	type SomeTypeReDef SomeType


	type SomeTypeDeepCompose struct {
		SomeTypeCompose
	}
	`)

	q := NewQuery(program, pkgImportPath)

	{
		typeName := q.TypeName("SomeType")
		typeFunc := MethodOf(typeName.Type().(*types.Named), "Fn")
		tt.NotNil(typeFunc)
	}

	{
		typeName := q.TypeName("SomeTypeAlias")
		typeFunc := MethodOf(typeName.Type().(*types.Named), "Fn")
		tt.NotNil(typeFunc)
	}

	{
		typeName := q.TypeName("SomeTypeCompose")
		typeFunc := MethodOf(typeName.Type().(*types.Named), "Fn")
		tt.NotNil(typeFunc)
	}

	{
		typeName := q.TypeName("SomeTypeDeepCompose")
		typeFunc := MethodOf(typeName.Type().(*types.Named), "Fn")
		tt.NotNil(typeFunc)
	}

	{
		typeName := q.TypeName("SomeTypeReDef")
		typeFunc := MethodOf(typeName.Type().(*types.Named), "Fn")
		tt.Nil(typeFunc)
	}
}
