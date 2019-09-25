package client

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/transport_grpc"
	"github.com/johnnyeven/libtools/courier/transport_http"
	"github.com/johnnyeven/libtools/env"
	"github.com/johnnyeven/libtools/log/context"
	"github.com/johnnyeven/libtools/servicex"
)

const depToolStackName = "dep-tools"

type Client struct {
	Name string
	// used in service
	Service       string
	Group         string
	Version       string
	Host          string `conf:"env,upstream" validate:"@hostname"`
	Mode          string
	Port          int
	Timeout       time.Duration
	WrapTransport transport_http.TransportWrapper
}

func (c Client) DockerDefaults() conf.DockerDefaults {
	if c.Name == "rancher" {
		return conf.DockerDefaults{
			// todo make switch in docker or expose
			"Host": conf.RancherInternal(depToolStackName, c.Name),
			"Port": 38080,
		}
	}
	return conf.DockerDefaults{
		// todo make switch in docker or expose
		"Host": conf.RancherInternal(c.Group, "service-"+c.Name),
		"Port": 80,
	}
}

func (Client) MarshalDefaults(v interface{}) {
	if client, ok := v.(*Client); ok {
		if client.Service == "" {
			client.Service = os.Getenv("PROJECT_NAME")
		}
		if client.Group == "" {
			client.Group = os.Getenv(servicex.EnvVarKeyProjectGroup)
		}
		if client.Version == "" {
			client.Version = os.Getenv("PROJECT_REF")
		}
		if client.Mode == "" {
			client.Mode = "http"
		}
		if client.Host == "" {
			client.Host = fmt.Sprintf("service-%s.staging.g7pay.net", client.Name)
		}
		if client.Port == 0 {
			client.Port = 80
		}
		if client.Timeout == 0 {
			client.Timeout = 5 * time.Second
		}
	}
}

func (c Client) GetBaseURL(protocol string) (url string) {
	url = c.Host
	if protocol != "" {
		url = fmt.Sprintf("%s://%s", protocol, c.Host)
	}
	if c.Port > 0 {
		url = fmt.Sprintf("%s:%d", url, c.Port)
	}
	return
}

func (c *Client) Request(id, httpMethod, uri string, req interface{}, metas ...courier.Metadata) IRequest {
	requestID := context.GetLogID()
	metadata := courier.MetadataMerge(metas...)

	if !env.IsOnline() {
		if requestIDInMeta := metadata.Get(httpx.HeaderRequestID); requestIDInMeta != "" {
			requestID = requestIDInMeta
		}
		mocker, err := ParseMockID(c.Service, requestID)
		if err == nil {
			if m, exists := mocker.Mocks[id]; exists {
				return &MockRequest{
					MockData: m,
				}
			}
		}
	}

	if metadata.Has(courier.VersionSwitchKey) {
		requestID = courier.ModifyRequestIDWithVersionSwitch(requestID, metadata.Get(courier.VersionSwitchKey))
	} else {
		if _, v, exists := courier.ParseVersionSwitch(requestID); exists {
			metadata.Set(courier.VersionSwitchKey, v)
		}
	}

	metadata.Add(httpx.HeaderRequestID, requestID)
	metadata.Add(httpx.HeaderUserAgent, c.Service+" "+c.Version)

	switch strings.ToLower(c.Mode) {
	case "grpc":
		serverName, method := parseID(id)
		return &transport_grpc.GRPCRequest{
			BaseURL:    c.GetBaseURL(""),
			ServerName: serverName,
			Method:     method,
			Timeout:    c.Timeout,
			RequestID:  requestID,
			Req:        req,
			Metadata:   metadata,
		}
	default:
		return &transport_http.HttpRequest{
			BaseURL:       c.GetBaseURL(c.Mode),
			Method:        httpMethod,
			URI:           uri,
			ID:            id,
			Timeout:       c.Timeout,
			WrapTransport: c.WrapTransport,
			Req:           req,
			Metadata:      metadata,
		}
	}
}

func parseID(id string) (serverName string, method string) {
	values := strings.Split(id, ".")
	if len(values) == 2 {
		serverName = strings.ToLower(strings.Replace(values[0], "Client", "", -1))
		method = values[1]
	}
	return
}
