package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/status_error"
)

func TestParameterErrors(t *testing.T) {
	tt := assert.New(t)

	parameterErrors := ParameterErrors{}

	parameterErrors.Merge(nil)
	tt.Nil(parameterErrors.StatusError)

	parameterErrors.Merge(status_error.InvalidField)
	tt.Equal(int64(status_error.InvalidField), parameterErrors.StatusError.Code)

	parameterErrors.Merge(status_error.InvalidField.StatusError().WithErrorField("query", "query", "error"))
	tt.Equal(status_error.ErrorFields{
		{
			In:    "query",
			Field: "query",
			Msg:   "error",
		},
	}, parameterErrors.StatusError.ErrorFields)

}
