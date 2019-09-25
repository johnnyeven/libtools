package testify

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/courier/transport_http/transform"
)

func NewContext(method string, url string, req interface{}) (ctx *gin.Context) {
	ctx = &gin.Context{}
	if req == nil {
		request, err := transform.NewHttpRequestFromParameterGroup(method, url, nil)
		if err != nil {
			panic(err)
		}
		ctx.Request = request
	}

	group := transform.ParameterGroupFromValue(req)
	request, err := transform.NewHttpRequestFromParameterGroup(method, url, group)
	if err != nil {
		panic(err)
	}
	ctx.Request = request

	for _, p := range group.Context {
		ctx.Set(p.Name, p.Value.Interface())
	}

	params, err := ParseParams(url, request.URL.Path)
	if err != nil {
		panic(err)
	}

	logrus.Infof("Gin Test Context %s %s", method, request.URL.String())

	ctx.Params = params

	rwiter := &TestResponseWriter{}
	rwiter.Body = &bytes.Buffer{}
	ctx.Writer = rwiter

	return

}

func ParseParams(path string, url string) (params gin.Params, err error) {
	pathArr := strings.Split(path, "/")
	urlArr := strings.Split(url, "/")

	if len(pathArr) != len(urlArr) {
		return nil, fmt.Errorf("url %s is not match path %s", url, path)
	}

	for i, p := range pathArr {
		if strings.HasPrefix(p, ":") {
			params = append(params, gin.Param{
				Key:   p[1:],
				Value: urlArr[i],
			})
		}
	}

	return params, nil
}
