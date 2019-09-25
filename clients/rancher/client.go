package rancher

import (
	"github.com/johnnyeven/libtools/courier/client"
	"fmt"
	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
)

type ClientRancher struct {
	AccessKey    string `conf:"env"`
	AccessSecret string `conf:"env"`
	client.Client
}

func (ClientRancher) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientRancher); ok {
		cl.Name = "rancher"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientRancher) Init() {
	c.CheckService()
}

func (c ClientRancher) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := status_error.FromError(err)
	if statusErr.Code == int64(status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

func (c ClientRancher) GetServices(stackId string) (resp *GetServicesResponse, err error) {
	resp = &GetServicesResponse{}
	resp.Meta = courier.Metadata{}

	metas := courier.Metadata{}
	metas.Set(httpx.HeaderAuthorization, c.AccessKey, c.AccessSecret)

	err = c.Request(c.Name+".GetServices", "GET", fmt.Sprintf("/v2-beta/stacks/%s/services", stackId), nil, metas).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetServicesResponse struct {
	Meta courier.Metadata
	Body Services
}
