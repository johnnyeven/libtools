package transport_http

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetParams(t *testing.T) {
	tt := assert.New(t)

	{
		_, err := GetParams("/:a/:b/:c", "/a/c")
		tt.NotNil(err)
	}

	{
		params, err := GetParams("/:a/:b/:c", "/a/b/c")
		tt.Nil(err)
		tt.Equal(httprouter.Params{
			{
				Key:   "a",
				Value: "a",
			},
			{
				Key:   "b",
				Value: "b",
			},
			{
				Key:   "c",
				Value: "c",
			},
		}, params)
	}
}

func TestGetClientIP(t *testing.T) {
}
