package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateScanner(t *testing.T) {
	tt := assert.New(t)

	{
		var v State
		err := v.Scan(-3)
		tt.NoError(err)
		tt.Equal(STATE__ACTIVE, v)
		dv, _ := v.Value()
		tt.Equal(int64(-3), dv)
	}

	{
		var v State
		err := v.Scan("-3")
		tt.NoError(err)
		tt.Equal(STATE__ACTIVE, v)
		dv, _ := v.Value()
		tt.Equal(int64(-3), dv)
	}

	{
		var v State
		err := v.Scan("")
		tt.NoError(err)
		tt.Equal(STATE_UNKNOWN, v)
	}

	{
		var v State
		err := v.Scan(nil)
		tt.NoError(err)
		tt.Equal(STATE_UNKNOWN, v)
	}
}
