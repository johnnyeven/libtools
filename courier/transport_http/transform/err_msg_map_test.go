package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/status_error"
)

func TestErrMsgMap(t *testing.T) {
	tt := assert.New(t)

	errMsgMap := ErrMsgMap{}

	tt.Nil(errMsgMap.ErrorFieldsIn("query", ""))

	errMsgMap.Set("key", "error")

	tt.Equal(ErrMsgMap{
		"key": "error",
	}, errMsgMap)

	errMsgMap2 := ErrMsgMap{}
	errMsgMap2.Set("key2", "error")

	tt.Equal(ErrMsgMap{
		"key2": "error",
	}, errMsgMap2)

	newErrMsgMap := errMsgMap.Merge(errMsgMap2).Merge(nil)

	tt.Equal(ErrMsgMap{
		"key":  "error",
		"key2": "error",
	}, newErrMsgMap)

	tt.Equal(status_error.ErrorFields{
		{
			Field: "key",
			In:    "query",
			Msg:   "error",
		},
		{
			Field: "key2",
			In:    "query",
			Msg:   "error",
		},
	}, newErrMsgMap.ErrorFieldsIn("query", ""))
}
