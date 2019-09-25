package context

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"encoding/json"
)

func GetTestContextBody(c *gin.Context) []byte {
	if rwriter, ok := c.Writer.(*TestResponseWriter); ok {
		return rwriter.Body.Bytes()
	}
	panic("context's writer is not TestResponseWriter")
}

func GetTestContextStatus(c *gin.Context) int {
	return int(reflect.Indirect(reflect.ValueOf(c)).FieldByName("writermem").FieldByName("status").Int())
}

func UnmarshalJSONTestContextBody(c *gin.Context, resp interface{}) {
	bodyBytes := GetTestContextBody(c)
	err := json.Unmarshal(bodyBytes, resp)
	if err != nil {
		panic(fmt.Sprintf("unmashal context's body failed[err:%s]", err.Error()))
	}
}

type RequestOpts struct {
	FormData map[string]string
	Query    map[string]string
	Params   map[string]string
	Headers  map[string]string
	Context  map[string]interface{}
	Body     interface{}
}

func (opts RequestOpts) Merge(reqOptions RequestOpts) RequestOpts {
	for k, v := range reqOptions.FormData {
		opts.FormData[k] = v
	}
	for k, v := range reqOptions.Query {
		opts.Query[k] = v
	}
	for k, v := range reqOptions.Params {
		opts.Params[k] = v
	}
	for k, v := range reqOptions.Headers {
		opts.Headers[k] = v
	}
	for k, v := range reqOptions.Context {
		opts.Context[k] = v
	}
	if reqOptions.Body != nil {
		opts.Body = reqOptions.Body
	}
	return opts
}

func toRawQuery(queryMap map[string]string) string {
	queryString := ""
	i := 0

	for k, v := range queryMap {
		keyValue := k + "=" + v
		if i == 0 {
			queryString += "?" + keyValue
		} else {
			queryString += "&" + keyValue
		}
		i++
	}

	return queryString
}

func WithRequest(opts RequestOpts) *gin.Context {
	c := gin.Context{}
	rv := reflect.ValueOf(opts.Body)
	var request *http.Request

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		request, _ = http.NewRequest("GET", "/"+toRawQuery(opts.Query), nil)
	} else {
		bodyBytes, _ := json.Marshal(opts.Body)
		request, _ = http.NewRequest("GET", "/"+toRawQuery(opts.Query), bytes.NewBuffer(bodyBytes))
	}

	for k, v := range opts.Headers {
		request.Header.Add(k, v)
	}

	for k, v := range opts.FormData {
		request.PostForm.Add(k, v)
	}

	c.Request = request
	rwriter := TestResponseWriter{}
	rwriter.Body = *bytes.NewBuffer([]byte{})
	c.Writer = &rwriter

	params := gin.Params{}
	for k, v := range opts.Params {
		params = append(params, gin.Param{
			Key:   k,
			Value: v,
		})
	}
	c.Params = params

	for k, v := range opts.Context {
		c.Set(k, v)
	}

	return &c
}
