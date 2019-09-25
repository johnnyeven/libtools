package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/status_error"
)

type Anonymous string

type SomeData struct {
	String  string  `json:"string" validate:"@string[3,]"`
	Pointer *string `json:"pointer" validate:"@string[3,]"`
}

type SomeReq struct {
	a string
	Anonymous
	P       P        `name:"p" in:"query"`
	Query   string   `name:"query" in:"query" validate:"@string[3,]"`
	Pointer *string  `name:"pointer" in:"query" validate:"@string[3,]"`
	Bytes   []byte   `name:"bytes" in:"query"`
	Data    SomeData `in:"body"`
}

func (v SomeReq) ValidateQuery() string {
	return "hook failed"
}

func TestParameterGroupFromValue(t *testing.T) {
	tt := assert.New(t)

	pg := ParameterGroupFromValue(&SomeReq{})

	tt.Len(pg.Parameters.List(), 4)

	valid, errFields := pg.ValidateNoBodyByHook()
	tt.False(valid)
	tt.Equal(status_error.ErrorFields{
		status_error.NewErrorField("query", "query", "hook failed"),
	}, errFields)
}

type SomeFormData struct {
	a string
	Anonymous
	String  string  `name:"string" validate:"@string[3,]"`
	Pointer *string `name:"pointer" validate:"@string[3,]"`
}

type SomeReqWithFormData struct {
	S        string       `in:"formData" name:"s" validate:"@string[3,]"`
	FormData SomeFormData `in:"formData"`
}

func TestParameterGroupFromWithFormData(t *testing.T) {
	tt := assert.New(t)

	pg := ParameterGroupFromValue(&SomeReqWithFormData{})

	tt.Len(pg.Parameters.List(), 0)
	tt.Len(pg.FormData.List(), 3)
}
