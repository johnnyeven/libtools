package client_id

import (
	"fmt"

	github_com_johnnyeven_libtools_courier "github.com/johnnyeven/libtools/courier"
	github_com_johnnyeven_libtools_courier_client "github.com/johnnyeven/libtools/courier/client"
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
)

type ClientIDInterface interface {
	GetNewId(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetNewIdResponse, err error)
}

type ClientID struct {
	github_com_johnnyeven_libtools_courier_client.Client
}

func (ClientID) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientID); ok {
		cl.Name = "id"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientID) Init() {
	c.CheckService()
}

func (c ClientID) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := github_com_johnnyeven_libtools_courier_status_error.FromError(err)
	if statusErr.Code == int64(github_com_johnnyeven_libtools_courier_status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

func (c ClientID) GetNewId(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetNewIdResponse, err error) {
	resp = &GetNewIdResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".GetNewId", "GET", "/id/v0/id", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetNewIdResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body UniqueID
}

func (c ClientID) Swagger(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *SwaggerResponse, err error) {
	resp = &SwaggerResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".Swagger", "GET", "/id", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type SwaggerResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body JSONBytes
}
