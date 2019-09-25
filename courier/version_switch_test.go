package courier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersionSwitch(t *testing.T) {
	tt := assert.New(t)

	{
		others, v, exist := ParseVersionSwitch("aadfafdads@x-version(VERSION_A)")
		tt.Equal("VERSION_A", v)
		tt.Equal("aadfafdads", others)
		tt.True(exist)
	}

	{
		others, _, exist := ParseVersionSwitch("aadfafdads")
		tt.Equal("aadfafdads", others)
		tt.False(exist)
	}
}
