package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizations(t *testing.T) {
	tt := assert.New(t)

	auths := Authorizations{}

	auths.Add("Bearer", "xxxxx")
	auths.Add("WechatBearer", "yyyyy")

	t.Log(auths.String())

	tt.Equal(auths, ParseAuthorization(auths.String()))
	tt.Equal("xxxxx", auths.Get("bearer"))
	tt.Equal("yyyyy", auths.Get("WechatBearer"))
}
