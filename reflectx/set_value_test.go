package reflectx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/ptr"
)

func TestSetValue(t *testing.T) {
	tt := assert.New(t)

	{
		number := 1
		tt.Nil(SetValue(reflect.ValueOf(&number), reflect.TypeOf(number), 2))
		tt.Equal(2, number)
	}

	{
		number := ptr.Int(1)
		tt.Nil(SetValue(reflect.ValueOf(number), reflect.TypeOf(number), 2))
		tt.NotNil(number)
		tt.Equal(2, *number)
	}

	{
		s := struct {
			Int *int
		}{}
		tt.Nil(SetValue(reflect.Indirect(reflect.ValueOf(&s)).FieldByName("Int"), reflect.TypeOf(s.Int), 2))
		tt.NotNil(s.Int)
		tt.Equal(2, *s.Int)
	}
}
