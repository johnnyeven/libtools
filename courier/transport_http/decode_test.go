package transport_http

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/johnnyeven/libtools/courier"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/transport_http/transform"
	"github.com/johnnyeven/libtools/ptr"
)

type GetItem struct {
	httpx.MethodGet
	ID     int      `name:"id" in:"path"`
	Ping   *string  `name:"ping" in:"query"`
	Header *string  `name:"header" default:"" in:"header"`
	Slice  []string `name:"slice" default:"" in:"query"`
	Body   struct {
		String string `json:"body" default:"" in:"query"`
	} `in:"body"`
}

func (getItem GetItem) MarshalDefaults(v interface{}) {
	v.(*GetItem).Ping = getItem.Ping
}

func (getItem GetItem) Path() string {
	return "/:id"
}

func (getItem GetItem) Output(ctx context.Context) (resp interface{}, err error) {
	return
}

func TestCreateHttpRequestDecoder(t *testing.T) {
	tt := assert.New(t)

	request := GetItem{
		ID:    1,
		Ping:  ptr.String("1"),
		Slice: []string{"1", "2", "3"},
	}
	request.Body.String = "111"

	req, err := transform.NewRequest(request.Method(), request.Path(), request)
	if ok := tt.NoError(err); !ok {
		return
	}
	t.Log(req.URL.String())

	getItem := GetItem{}
	err = MarshalOperator(req, &getItem)
	tt.Nil(err)

	tt.Equal(request, getItem)
}

func TestCreateHttpRequestDecoder_VersionSwitch(t *testing.T) {
	tt := assert.New(t)

	{
		respData := struct {
			courier.WithVersionSwitch
			courier.EmptyOperator
		}{}

		req, _ := transform.NewRequest("GET", "/", nil, courier.MetadataWithVersionSwitch("VERSION"))

		err := MarshalOperator(req, &respData)
		tt.NoError(err)

		tt.Equal(respData.XVersion, "VERSION")
	}

	{
		respData := struct {
			courier.WithVersionSwitch
			courier.EmptyOperator
		}{}

		req, _ := transform.NewRequest("GET", "/", nil, courier.Metadata{
			httpx.HeaderRequestID: []string{courier.ModifyRequestIDWithVersionSwitch("adadasd", "VERSION")},
		})

		err := MarshalOperator(req, &respData)
		tt.NoError(err)

		tt.Equal(respData.XVersion, "VERSION")
	}
}

type PostForm struct {
	httpx.MethodPost
	ID       int `name:"id" in:"path"`
	FormData struct {
		FirstFile  *multipart.FileHeader `name:"firstFile"`
		SecondFile *multipart.FileHeader `name:"secondFile"`
		Data       struct {
			Value string `json:"value"`
		} `name:"data"`
	} `in:"formData,multipart"`
}

func (postForm PostForm) Path() string {
	return "/:id"
}

func (postForm PostForm) Output(ctx context.Context) (resp interface{}, err error) {
	return
}

func TestCreateHttpRequestDecoderWithForm(t *testing.T) {
	tt := assert.New(t)

	request := PostForm{}
	request.ID = 2
	request.FormData.Data.Value = "1111"
	request.FormData.FirstFile, _ = transform.NewFileHeader("firstFile", "SingleFile", []byte("1"))
	request.FormData.SecondFile, _ = transform.NewFileHeader("secondFile", "SecondFile", []byte("2"))

	req, err := transform.NewRequest(request.Method(), request.Path(), request)
	if ok := tt.NoError(err); !ok {
		return
	}
	t.Log(req.URL.String())

	postForm := PostForm{}
	err = MarshalOperator(req, &postForm)
	tt.NoError(err)
	tt.Equal(request, postForm)

	{
		actualFile, _ := postForm.FormData.FirstFile.Open()
		actualBytes, _ := ioutil.ReadAll(actualFile)
		tt.Equal([]byte("1"), actualBytes)
	}

	{
		actualFile, _ := postForm.FormData.SecondFile.Open()
		actualBytes, _ := ioutil.ReadAll(actualFile)
		tt.Equal([]byte("2"), actualBytes)
	}
}

type PostFormURLEncoded struct {
	httpx.MethodPost
	FormData struct {
		String string   `name:"string"`
		Slice  []string `name:"slice"`
		Data   struct {
			Value string `json:"value"`
		} `name:"data"`
	} `in:"formData,urlencoded"`
}

func (PostFormURLEncoded) Path() string {
	return "/"
}

func (PostFormURLEncoded) TransformHttpRequest(req *http.Request) {
	req.Header.Set(httpx.HeaderContentType, "application/x-www-form-urlencoded; param=value")
}

func (postForm PostFormURLEncoded) Output(ctx context.Context) (resp interface{}, err error) {
	return
}

func TestCreateHttpRequestDecoderWithPostFormURLEncoded(t *testing.T) {
	tt := assert.New(t)

	request := PostFormURLEncoded{}
	request.FormData.String = "1111"
	request.FormData.Slice = []string{"1111", "2222"}
	request.FormData.Data.Value = "1"

	req, err := transform.NewRequest(request.Method(), request.Path(), request)
	if ok := tt.NoError(err); !ok {
		return
	}
	t.Log(req.URL.String())

	postForm := PostFormURLEncoded{}
	err = MarshalOperator(req, &postForm)
	tt.NoError(err)

	tt.Equal(request, postForm)
}

func TestCreateHttpRequestDecoderWithPostFormURLEncodedLimitContentType(t *testing.T) {
	tt := assert.New(t)

	request := PostFormURLEncoded{}
	request.FormData.String = "1111"
	request.FormData.Slice = []string{"1111", "2222"}
	request.FormData.Data.Value = "1"

	req, err := transform.NewRequest(request.Method(), request.Path(), request)
	if ok := tt.NoError(err); !ok {
		return
	}

	buf := &bytes.Buffer{}

	io.Copy(buf, req.Body)

	finalRequest, err := http.NewRequest(
		request.Method(),
		request.Path(),
		buf,
	)

	t.Log(req.Header.Get(httpx.HeaderContentType))

	postForm := PostFormURLEncoded{}
	err = MarshalOperator(finalRequest, &postForm)
	tt.NoError(err)

	tt.Equal(request, postForm)
}
