package client

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/johnnyeven/libtools/courier"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/status_error"
)

type IpInfo struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

func TestClient_DoRequest(t *testing.T) {
	tt := assert.New(t)

	ipInfoClient := Client{
		Service: "test",
		Mode:    "http",
		Host:    "ip-api.com",
		Timeout: 100 * time.Second,
	}

	{
		ipInfo := IpInfo{}
		err := ipInfoClient.
			Request("id", "GET", "/json", nil).
			Do().
			Into(&ipInfo)

		if err == nil {
			tt.Equal("China", ipInfo.Country)
			tt.Equal("CN", ipInfo.CountryCode)
		} else {
			spew.Dump(err.(*status_error.StatusError))
		}
	}
	{
		data, _ := json.Marshal(IpInfo{
			Country:     "USA",
			CountryCode: "US",
		})

		ipInfo := IpInfo{}

		mockData := MetadataWithMocker(Mock(ipInfoClient.Service).For("id", MockData{
			Data: data,
		}))

		err := ipInfoClient.Request("id", "GET", "/json", nil, mockData).Do().Into(&ipInfo)

		tt.Nil(err)
		tt.Equal("USA", ipInfo.Country)
		tt.Equal("US", ipInfo.CountryCode)
	}
}

func TestClient_DoRequestWithMetadata(t *testing.T) {
	tt := assert.New(t)

	ipInfoClient := Client{
		Service: "test",
		Mode:    "http",
		Host:    "ip-api.com",
		Timeout: 100 * time.Second,
	}

	ipInfo := IpInfo{}

	err := ipInfoClient.Request("id", "GET", "/json", nil, courier.MetadataWithVersionSwitch("VERSION")).
		Do().
		Into(&ipInfo)

	tt.Nil(err)

	spew.Dump(ipInfo)
}
