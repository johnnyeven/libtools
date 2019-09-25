package sign

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/libtools/courier/status_error"
)

func TestAutoSignTransport(t *testing.T) {
	tt := assert.New(t)

	ipInfoClient := client.Client{
		Service: "test",
		Mode:    "http",
		Host:    "ip-api.com",
		Timeout: 100 * time.Second,
		WrapTransport: NewAutoSignTransport(func(key string) (string, error) {
			return "111", nil
		}),
	}

	p := SignParams{
		AccessKey:  "123",
		RandString: "1123",
	}

	ipInfo := IpInfo{}
	err := ipInfoClient.
		Request("id", "GET", "/json", p).
		Do().
		Into(&ipInfo)

	if err == nil {
		tt.Equal("China", ipInfo.Country)
		tt.Equal("CN", ipInfo.CountryCode)
	} else {
		spew.Dump(err.(*status_error.StatusError))
	}
}

type IpInfo struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}
