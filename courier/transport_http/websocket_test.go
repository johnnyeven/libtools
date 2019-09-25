package transport_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/timelib"
)

func TestPayloadMarshalUnMarshal(t *testing.T) {
	tt := assert.New(t)

	{
		s := "string"
		bytes, err := PayloadMarshal(s)
		tt.Nil(err)
		t.Log(bytes)
		var s2 string
		errForUnmarshal := PayloadUnmarshal(bytes, &s2)
		tt.Nil(errForUnmarshal)
		t.Log(s2)
		tt.Equal(s2, s)
	}

	{
		s := 1
		bytes, err := PayloadMarshal(s)
		tt.Nil(err)
		t.Log(bytes)
		var s2 int
		errForUnmarshal := PayloadUnmarshal(bytes, &s2)
		tt.Nil(errForUnmarshal)
		t.Log(s2)
		tt.Equal(s2, s)
	}

	{
		s := true
		bytes, err := PayloadMarshal(s)
		tt.Nil(err)
		t.Log(bytes)
		var s2 bool
		errForUnmarshal := PayloadUnmarshal(bytes, &s2)
		tt.Nil(errForUnmarshal)
		t.Log(s2)
		tt.Equal(s2, s)
	}

	{
		s := timelib.MySQLTimestamp(time.Now())
		bytes, err := PayloadMarshal(s)
		tt.Nil(err)
		t.Log(bytes)
		var s2 timelib.MySQLTimestamp
		errForUnmarshal := PayloadUnmarshal(bytes, &s2)
		tt.Nil(errForUnmarshal)
		t.Log(s2)
		tt.Equal(s2.String(), s.String())
	}
}
