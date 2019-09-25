package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitToWords(tt *testing.T) {
	t := assert.New(tt)

	t.Equal([]string{}, SplitToWords(""))
	t.Equal([]string{"lowercase"}, SplitToWords("lowercase"))
	t.Equal([]string{"Class"}, SplitToWords("Class"))
	t.Equal([]string{"My", "Class"}, SplitToWords("MyClass"))
	t.Equal([]string{"My", "C"}, SplitToWords("MyC"))
	t.Equal([]string{"HTML"}, SplitToWords("HTML"))
	t.Equal([]string{"PDF", "Loader"}, SplitToWords("PDFLoader"))
	t.Equal([]string{"A", "String"}, SplitToWords("AString"))
	t.Equal([]string{"Simple", "XML", "Parser"}, SplitToWords("SimpleXMLParser"))
	t.Equal([]string{"vim", "RPC", "Plugin"}, SplitToWords("vimRPCPlugin"))
	t.Equal([]string{"GL", "11", "Version"}, SplitToWords("GL11Version"))
	t.Equal([]string{"99", "Bottles"}, SplitToWords("99Bottles"))
	t.Equal([]string{"May", "5"}, SplitToWords("May5"))
	t.Equal([]string{"BFG", "9000"}, SplitToWords("BFG9000"))
	t.Equal([]string{"Böse", "Überraschung"}, SplitToWords("BöseÜberraschung"))
	t.Equal([]string{"Two", "spaces"}, SplitToWords("Two  spaces"))
	t.Equal([]string{"BadUTF8\xe2\xe2\xa1"}, SplitToWords("BadUTF8\xe2\xe2\xa1"))
	t.Equal([]string{"snake", "case"}, SplitToWords("snake_case"))
	t.Equal([]string{"snake", "case"}, SplitToWords("snake_ case"))
}

func TestToUpperCamelCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("SnakeCase", ToUpperCamelCase("snake_case"))
	t.Equal("IDCase", ToUpperCamelCase("id_case"))
}

func TestToLowerCamelCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("snakeCase", ToLowerCamelCase("snake_case"))
	t.Equal("idCase", ToLowerCamelCase("id_case"))
}

func TestToUpperSnakeCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("SNAKE_CASE", ToUpperSnakeCase("snakeCase"))
	t.Equal("ID_CASE", ToUpperSnakeCase("idCase"))
}

func TestToLowerSnakeCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("snake_case", ToLowerSnakeCase("snakeCase"))
	t.Equal("id_case", ToLowerSnakeCase("idCase"))
	t.Equal("i7_case", ToLowerSnakeCase("i7Case"))
}
