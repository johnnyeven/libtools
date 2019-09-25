package transform

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/reflectx"
)

func TestParameterMap(t *testing.T) {
	tt := assert.New(t)

	pm := ParameterMap{}

	type SomeReq struct {
		Query   string  `name:"query" in:"query" validate:"@string[3,]"`
		Pointer *string `name:"pointer" in:"query" validate:"@string[3,]"`
	}

	req := &SomeReq{}
	tpe := reflectx.IndirectType(reflect.TypeOf(req))
	rv := reflect.Indirect(reflect.ValueOf(req))

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := rv.Field(i)

		tagIn, _, tagInFlags := GetTagIn(&field)
		pm.Add(NewParameterMeta(&field, fieldValue, tagIn, tagInFlags))
	}

	tt.Equal(2, len(pm.List()))
	tt.Equal(2, pm.Len())

	{
		p, ok := pm.Get("Query")
		tt.True(ok)
		tt.Equal("query", p.In)
	}
}
