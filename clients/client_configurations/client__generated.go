package client_configurations

import (
	"fmt"

	github_com_johnnyeven_libtools_courier "github.com/johnnyeven/libtools/courier"
	github_com_johnnyeven_libtools_courier_client "github.com/johnnyeven/libtools/courier/client"
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
)

type ClientConfigurationsInterface interface {
	BatchCreateConfig(req BatchCreateConfigRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *BatchCreateConfigResponse, err error)
	CreateConfig(req CreateConfigRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *CreateConfigResponse, err error)
	GetConfigurations(req GetConfigurationsRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetConfigurationsResponse, err error)
}

type ClientConfigurations struct {
	github_com_johnnyeven_libtools_courier_client.Client
}

func (ClientConfigurations) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientConfigurations); ok {
		cl.Name = "configurations"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientConfigurations) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := github_com_johnnyeven_libtools_courier_status_error.FromError(err)
	if statusErr.Code == int64(github_com_johnnyeven_libtools_courier_status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

type BatchCreateConfigRequest struct {
	//
	Body []CreateConfigurationBody `fmt:"json" in:"body"`
}

func (c ClientConfigurations) BatchCreateConfig(req BatchCreateConfigRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *BatchCreateConfigResponse, err error) {
	resp = &BatchCreateConfigResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".BatchCreateConfig", "POST", "/configurations/v0/configurations/0/batch", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type BatchCreateConfigResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body []byte
}

type CreateConfigRequest struct {
	//
	Body CreateConfigurationBody `fmt:"json" in:"body"`
}

func (c ClientConfigurations) CreateConfig(req CreateConfigRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *CreateConfigResponse, err error) {
	resp = &CreateConfigResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".CreateConfig", "POST", "/configurations/v0/configurations", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type CreateConfigResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body []byte
}

type GetConfigurationsRequest struct {
	// 分页偏移
	// 默认为 0
	Offset int32 `in:"query" name:"offset,omitempty"`
	// StackID
	StackID uint64 `in:"query" name:"stackID"`
	// 分页大小
	// 默认为 10，-1 为查询所有
	Size int32 `default:"10" in:"query" name:"size,omitempty"`
}

func (c ClientConfigurations) GetConfigurations(req GetConfigurationsRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetConfigurationsResponse, err error) {
	resp = &GetConfigurationsResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".GetConfigurations", "GET", "/configurations/v0/configurations", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetConfigurationsResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body GetConfigurationResult
}

func (c ClientConfigurations) Swagger(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *SwaggerResponse, err error) {
	resp = &SwaggerResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".Swagger", "GET", "/configurations", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type SwaggerResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body JSONBytes
}
