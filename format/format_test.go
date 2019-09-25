package format_test

import (
	"fmt"
	"testing"

	"github.com/johnnyeven/libtools/format"
)

func TestProcess(t *testing.T) {
	result, _ := format.Process("format2_test.go", []byte(`
	package format

	import (
		testing "testing"
		"fmt"

		"github.com/johnnyeven/libtools/format"
		"github.com/johnnyeven/libtools/gin_app"

		"github.com/davecgh/go-spew/spew"
	)

	func Test(t *testing.T) {
		spew.Dump(gin_app.REQUEST_ID_NAME, format.Test)
	}
	`))

	fmt.Println(string(result))
}
