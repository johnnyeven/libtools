package rancher

import (
	"testing"
	"github.com/johnnyeven/libtools/courier/client"
)

var cli *ClientRancher

func init() {
	cli = &ClientRancher{
		AccessKey:    "D88A9E1985EF5A62F0FA",
		AccessSecret: "hub81J4KSPwHSByFxzhnVRCNN6SCnNS9wkaDE9N2",
		Client: client.Client{
			Name: "rancher",
			Host: "rancher.profzone.net",
			Port: 38080,
		},
	}
	cli.MarshalDefaults(cli)
}

func TestGetServices(t *testing.T) {
	resp, err := cli.GetServices("1st25")
	if err != nil {
		t.Error(err)
	}

	t.Log(resp.Body)
}
