package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImports(t *testing.T) {
	tt := assert.New(t)
	imports := Importer{}

	tt.Equal(`[]string{"1", "2"}`, imports.Sdump([]string{"1", "2"}))
	tt.Equal(`[]interface {}{"1", nil}`, imports.Sdump([]interface{}{"1", nil}))
	tt.Equal(`map[string]int{"1": 2}`, imports.Sdump(map[string]int{"1": 2}))
}
