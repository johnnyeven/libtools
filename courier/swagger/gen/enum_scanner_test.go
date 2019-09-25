package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/codegen/loaderx"
)

func TestEnumScanner_WithSkip(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	type State int

	const (
		STATE_UNKNOWN State = iota
		STATE__ONE          // one
		STATE__TWO          // two
		STATE__THREE        // three
		_
		STATE__FOUR State = iota + 100 // four
	)
	`)

	scanner := NewEnumScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)
	enums := scanner.Enum(query.TypeName("State"))

	tt.Equal(Enum{
		{
			Label: "four",
			Value: "FOUR",
			Val:   int64(105),
		},
		{
			Value: "ONE",
			Label: "one",
			Val:   int64(1),
		},
		{
			Value: "THREE",
			Label: "three",
			Val:   int64(3),
		},
		{
			Value: "TWO",
			Label: "two",
			Val:   int64(2),
		},
	}, enums)
}

func TestEnumScanner_MultipleBlock(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	type State int

	const (
		STATE_UNKNOWN State = iota
		STATE__ONE          // one
	)

	const (
		STATE__FOUR State = iota + 100 // four
	)
	`)

	scanner := NewEnumScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)
	enums := scanner.Enum(query.TypeName("State"))

	tt.Equal(Enum{
		{
			Label: "four",
			Value: "FOUR",
			Val:   int64(100),
		},
		{
			Label: "one",
			Value: "ONE",
			Val:   int64(1),
		},
	}, enums)
}

func TestEnumScanner_TestEnumValDirectly(t *testing.T) {
	tt := assert.New(t)

	pkgImportPath, program := loaderx.NewTestProgram(`
	package main

	type State string

	const (
		One  State  = "ONE"			// one
		Four State  = "FOUR" 		// four
	)
	`)

	scanner := NewEnumScanner(program)
	query := loaderx.NewQuery(program, pkgImportPath)
	enums := scanner.Enum(query.TypeName("State"))

	tt.Equal(Enum{
		{
			Label: "four",
			Value: "FOUR",
			Val:   "FOUR",
		},
		{
			Label: "one",
			Value: "ONE",
			Val:   "ONE",
		},
	}, enums)
}
