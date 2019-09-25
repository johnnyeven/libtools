package codegen

import (
	"testing"
)

type GenDemo struct {
}

func (demo *GenDemo) Load(cwd string) {
}

func (demo *GenDemo) Pick() {
}

func (demo *GenDemo) Output(cwd string) Outputs {
	return Outputs{
		"./doc_test.go": `
		package codegen

		func Test(t *testing.T) {
		}
		`,
	}
}

func TestGenerate(t *testing.T) {
	genDemo := GenDemo{}
	Generate(&genDemo)
}
