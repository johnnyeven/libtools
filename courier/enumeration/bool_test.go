package enumeration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	tt := assert.New(t)
	{
		bytes, _ := json.Marshal(BOOL__TRUE)
		tt.Equal("true", string(bytes))
	}
	{
		bytes, _ := json.Marshal(BOOL__FALSE)
		tt.Equal("false", string(bytes))
	}
	{
		bytes, _ := json.Marshal(BOOL_UNKNOWN)
		tt.Equal("null", string(bytes))
	}

	{
		var b Bool
		json.Unmarshal([]byte("true"), &b)
		tt.Equal(BOOL__TRUE, b)
	}
	{
		var b Bool
		json.Unmarshal([]byte("false"), &b)
		tt.Equal(BOOL__FALSE, b)
	}
	{
		var b Bool
		json.Unmarshal([]byte("null"), &b)
		tt.Equal(BOOL_UNKNOWN, b)
	}
}
