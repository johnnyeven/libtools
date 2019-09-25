package gen

import (
	"os"
	"testing"

	"github.com/johnnyeven/libtools/codegen"
)

func init() {
	os.Chdir("./test")
}

func TestGen(t *testing.T) {
	clientGenerator := SqlFuncGenerator{}
	clientGenerator.StructName = "User"
	clientGenerator.Database = "DBTest"
	clientGenerator.WithTableInterfaces = true
	codegen.Generate(&clientGenerator)
}
