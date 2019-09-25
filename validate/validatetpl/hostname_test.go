package validatetpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHostname(t *testing.T) {
	tt := assert.New(t)

	{
		got, _ := ValidateHostname("134.255.255.1")
		tt.True(got)
	}

	{
		got, _ := ValidateHostname("localhost")
		tt.True(got)
	}

	{
		got, _ := ValidateHostname("service-etc-message.service-etc-message.rancher.internal")
		tt.True(got)
	}

	{
		got, _ := ValidateHostname("baidu.com")
		tt.True(got)
	}

	{
		got, errStr := ValidateHostname("134.259.255.1")
		tt.False(got)
		tt.Equal(InvalidHostnameValue, errStr)
	}

	{
		got, errStr := ValidateHostname("134.255.255.1/")
		tt.False(got)
		tt.Equal(InvalidHostnameValue, errStr)
	}

	{
		got, errStr := ValidateHostname(123123)
		tt.False(got)
		tt.Equal(InvalidHostnameType, errStr)
	}
}
