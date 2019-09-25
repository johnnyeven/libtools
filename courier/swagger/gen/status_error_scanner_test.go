package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

func TestStatusErrorScanner(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	import (
		"github.com/johnnyeven/libtools/courier/status_error"
		"net/http"
	)

	const (
		// 内部未明确定义的错误
		UnknownError status_error.StatusErrorCode = http.StatusInternalServerError*1e6 + 1 + iota
		// 内部用于接收参数时非法的结构
		InvalidStructError
	)

	// @httpError(500000003,ReadFailed,"Read 调用时发生错误","",false);
	// doc
	func Err2() error {
		return UnknownError
	}

	func Err3() error {
		return status_error.DemoErr()
	}

	func Err() (err error) {
		if true {
			err = Err3()
			return
		}
		return Err2()
	}
	`)

	scanner := NewStatusErrorScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)
	statusErrors := scanner.StatusErrorsInFunc(query.Func("Err"))

	tt.Len(statusErrors, 3)
}
