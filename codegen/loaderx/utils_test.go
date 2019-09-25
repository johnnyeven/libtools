package loaderx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPkgImportPathAndExpose(t *testing.T) {
	tt := assert.New(t)

	cases := []struct {
		p string
		e string
		s string
	}{
		{
			"a",
			"B",
			"a.B",
		},
		{
			"a.b.c.d/c",
			"B",
			"a.b.c.d/c.B",
		},
	}

	for _, caseItem := range cases {
		p, e := GetPkgImportPathAndExpose(caseItem.s)
		tt.Equal(caseItem.p, p)
		tt.Equal(caseItem.e, e)
	}
}
