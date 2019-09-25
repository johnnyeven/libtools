package client

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestInsecureTLSTransport(t *testing.T) {
	rootCA := []byte(`-----BEGIN CERTIFICATE-----
MIIDfzCCAmegAwIBAgIEfZ0WpjANBgkqhkiG9w0BAQsFADBvMQswCQYDVQQGEwJj
bjELMAkGA1UECBMCemoxCzAJBgNVBAcTAmh6MQ8wDQYDVQQKEwZ6anBvcnQxEDAO
BgNVBAsTB2dhdGV3YXkxIzAhBgNVBAMTGm9wZW5hcGktdGVzdC56anBvcnQuZ292
LmNuMCAXDTE2MTIwODA3MjgwMloYDzIxMTYxMTE0MDcyODAyWjBvMQswCQYDVQQG
EwJjbjELMAkGA1UECBMCemoxCzAJBgNVBAcTAmh6MQ8wDQYDVQQKEwZ6anBvcnQx
EDAOBgNVBAsTB2dhdGV3YXkxIzAhBgNVBAMTGm9wZW5hcGktdGVzdC56anBvcnQu
Z292LmNuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw7wf0u/SjO9z
JRziwkbZ3BDXmNgIrsEbAZKRDWpUP8FLcq46aUdG32q4ra3hNoa+rVX8otgqtjX2
q/PswQ77PjAfzfiyKh7ROvVI90CiNtHxvNe207uEjiJ64xixvDprs6l6YgJHWReM
qyAsRT+7RSfxjze8RsJEHipM8zjUUDcaTjkJf+Ce8TDyv7RGM+AV3UGBrKqazCEt
oxkh9NPdTNRpBaKlb3j81kV7T9OVtdhBK4gVdVpziDpp0Iu9KnjtS+/NTl/NYVwn
7XSW1N3B3i+6Ckphwt4U0JCwbD6PXi6ggHHg/kwx8vx4wHMT5Xf6FaZONXeR9bmJ
ciCqhD+JOQIDAQABoyEwHzAdBgNVHQ4EFgQUbygOTIL8VDKsgs/sofORVKWx7kYw
DQYJKoZIhvcNAQELBQADggEBAKq9OfvfBaCIF5ES915lL3ifBTMZsX2x6EZULAoy
0jmkOlORHHDjSms5Kk5z+o+8CP1jNUFZXh3zKYju5b3oswMQ89LVp1M1J9BqrUhi
CYcUjfkVZk3iHxsxNDwKA3NIbC8E02AwHsRja+WesSo6AGnhqi9XUf0cVbx/RGJ/
nkCeyoHMVD5sy9N4putVsXsrwsLnhNHegUy6PYUa3yK32g30MEntEzGEZFxijkSa
EJ8q0fMGMxYCNaxfjXH5FzBH18HuM+i0Z23KwFuwPJ1YSI5PC5nfDyzh7HiYlHDW
OB1tH8y71EHV0kt9vyUcL1Q+SYyqtinveR8XZyTjvmxZCQE=
-----END CERTIFICATE-----`)

	zjClient := Client{
		Service:       "zjport",
		Mode:          "https",
		Host:          "openapi-test.zjport.gov.cn",
		Port:          8553,
		Timeout:       100 * time.Second,
		WrapTransport: NewInsecureTLSTransport(rootCA),
	}

	tt := assert.New(t)

	req := zjClient.Request("receive", "GET", "/gateway/receive", nil)
	data := make([]byte, 0)
	err := req.Do().Into(&data)

	tt.NoError(err)

	spew.Dump(string(data))
}
