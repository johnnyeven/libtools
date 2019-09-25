package validatetpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidateChar(t *testing.T) {
	tt := assert.New(t)
	{
		valid, _ := NewValidateChar(0, 2)("中文")
		tt.True(valid)
	}
	{
		valid, msg := NewValidateChar(0, 2)("中文中文")
		tt.False(valid)
		tt.Equal("字符串字数不在[0, 2]范围内,当前长度: 4", msg)
	}
}
