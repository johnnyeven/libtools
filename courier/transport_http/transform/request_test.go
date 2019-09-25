package transform

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"encoding/json"

	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/ptr"
)

func MarshalParametersWithPath(req *http.Request, path string, v interface{}) error {
	return MarshalParameters(ParameterGroupFromReflectValue(reflect.ValueOf(v)), &ParameterValuesGetterWithPath{
		ParameterValuesGetter: NewParameterValuesGetter(req),
		Path: path,
	})
}

type ParameterValuesGetterWithPath struct {
	*ParameterValuesGetter
	Path   string
	params httprouter.Params
}

func (getter *ParameterValuesGetterWithPath) Initial() (err error) {
	getter.params, err = GetParams(getter.Path, getter.Request.URL.Path)
	if err != nil {
		return
	}
	return getter.ParameterValuesGetter.Initial()
}

func GetParams(path string, url string) (params httprouter.Params, err error) {
	pathArr := strings.Split(httprouter.CleanPath(path), "/")
	urlArr := strings.Split(httprouter.CleanPath(url), "/")

	if len(pathArr) != len(urlArr) {
		return nil, fmt.Errorf("url %s is not match path %s", url, path)
	}

	for i, p := range pathArr {
		if strings.HasPrefix(p, ":") {
			params = append(params, httprouter.Param{
				Key:   p[1:],
				Value: urlArr[i],
			})
		}
	}

	return params, nil
}

func (getter *ParameterValuesGetterWithPath) Param(name string) string {
	return getter.params.ByName(name)
}

type MarshalInt int

func (marshalInt *MarshalInt) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		s = string(data)
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*marshalInt = MarshalInt(i)
	return nil
}

func (marshalInt MarshalInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, marshalInt)), nil
}

type StructBody struct {
	A int    `json:"a" xml:"a"`
	C uint   `json:"c" xml:"c"`
	D bool   `json:"d" xml:"d"`
	E string `json:"e" xml:"e"`
}

func TestTransToHttpReq(t *testing.T) {
	tt := assert.New(t)

	type CommonReq struct {
		A int              `name:"a" in:"header"`
		B enumeration.Bool `name:"b" in:"header"`
		C uint             `name:"c" in:"header"`
		D bool             `name:"d" in:"header"`
	}

	type SuccessStruct struct {
		CommonReq
		E string     `name:"e" in:"header"`
		Q MarshalInt `name:"q" in:"header"`

		F int     `name:"f" in:"path"`
		G uint64  `name:"g" in:"path" default:"0"`
		H float64 `name:"h" in:"path"`
		I *bool   `name:"i" in:"path"`
		J string  `name:"j" in:"path"`

		K    int        `name:"k" in:"query"`
		L    uint       `name:"l" in:"query"`
		N    bool       `name:"n" in:"query"`
		O    string     `name:"o" in:"query"`
		P    MarshalInt `name:"p" in:"query"`
		Ids  []int      `name:"ids" in:"query"`
		Body StructBody `in:"body"`
	}

	ss := &SuccessStruct{
		CommonReq: CommonReq{
			A: -1,
			B: enumeration.BOOL__FALSE,
			C: 1,
			D: true,
		},
		E: "a",
		Q: 11,

		F: -1,
		G: 1,
		H: 1.1,
		I: ptr.Bool(false),
		J: "b",

		K:   -1,
		L:   1,
		N:   true,
		O:   "c",
		P:   11,
		Ids: []int{1, 2},

		Body: StructBody{
			A: -1,
			C: 1,
			D: true,
			E: "d",
		},
	}

	req, err := NewRequest("POST", "http://127.0.0.1/:f/:g/:h/:i/:j", ss)
	tt.Nil(err)

	tt.Equal("-1", req.Header.Get("a"))
	tt.Equal("false", req.Header.Get("b"))
	tt.Equal("1", req.Header.Get("c"))
	tt.Equal("true", req.Header.Get("d"))
	tt.Equal("a", req.Header.Get("e"))
	tt.Equal("11", req.Header.Get("q"))
	tt.Equal("ids=1%2C2&k=-1&l=1&n=true&o=c&p=11", req.URL.Query().Encode())
	tt.Equal("/-1/1/1.1/false/b", req.URL.Path)

	bodyBytes, err := CloneRequestBody(req)
	tt.NoError(err)

	reqBody := &StructBody{}
	json.Unmarshal(bodyBytes, reqBody)
	tt.Equal(ss.Body, *reqBody)

	{
		ssForReceive := &SuccessStruct{}
		err := MarshalParametersWithPath(req, "/:f/:g/:h/:i/:j", ssForReceive)
		tt.NoError(err)
		tt.Equal(ss, ssForReceive)
	}
}

func TestTransToHttpReqWithXML(t *testing.T) {
	tt := assert.New(t)

	type SuccessStruct struct {
		Body struct {
			XMLName xml.Name `xml:"person" json:"-"`
			StructBody
		} `in:"body" fmt:"xml"`
	}

	ss := &SuccessStruct{}
	ss.Body.StructBody = StructBody{
		A: -1,
		C: 1,
		D: true,
		E: "d",
	}

	req, err := NewRequest("POST", "http://127.0.0.1", ss)
	tt.Nil(err)

	bodyBytes, err := CloneRequestBody(req)
	tt.NoError(err)

	reqBody := &struct {
		XMLName xml.Name `xml:"person" json:"-"`
		StructBody
	}{}

	xml.Unmarshal(bodyBytes, reqBody)
	tt.Equal(ss.Body.StructBody, reqBody.StructBody)
}

func TestTransToHttpReqWithBodyError(t *testing.T) {
	tt := assert.New(t)
	type SuccessStruct struct {
		Body StructBody `in:"body"`
	}

	ss := &SuccessStruct{
		Body: StructBody{
			A: -1,
			C: 1,
			D: true,
			E: "d",
		},
	}

	req, err := NewRequest("POST", "http://127.0.0.1", ss)
	tt.Nil(err)

	bodyBytes, err := ioutil.ReadAll(req.Body)
	tt.NoError(err)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0:9]))

	reqBody := &StructBody{}
	json.Unmarshal(bodyBytes, reqBody)
	tt.Equal(ss.Body, *reqBody)

	{
		ssForReceive := &SuccessStruct{}
		err := MarshalParametersWithPath(req, "/", ssForReceive)
		tt.Error(err)
		tt.Equal(int64(status_error.InvalidBodyStruct), status_error.FromError(err).Code)
	}
}

func TestTransToHttpReqWithoutParams(t *testing.T) {
	tt := assert.New(t)
	req, err := NewRequest("POST", "http://127.0.0.1", nil)
	tt.Nil(err)
	tt.Equal("http://127.0.0.1/", req.URL.String())
}

type FormData struct {
	A []int `name:"a"`
	C uint  `name:"c" `
}

func TestTransToHttpReqWithCookie(t *testing.T) {
	tt := assert.New(t)

	type Req struct {
		Cookie string   `in:"cookie" name:"cookie"`
		Slice  []string `in:"cookie" name:"slice"`
	}

	request := &Req{
		Cookie: "cookie",
		Slice:  []string{"1", "2"},
	}

	req, err := NewRequest("POST", "http://127.0.0.1", request)
	tt.Nil(err)
	tt.NotNil(req.Header.Get("Cookie"))

	{
		requestForReceive := &Req{}
		err := MarshalParametersWithPath(req, "/", requestForReceive)
		tt.NoError(err)
		tt.Equal(request, requestForReceive)
	}
}

func TestTransToHttpReqWithFormData(t *testing.T) {
	tt := assert.New(t)

	type Req struct {
		FormData `in:"formData"`
	}

	request := Req{
		FormData: FormData{
			A: []int{
				-1, 1,
			},
			C: 1,
		},
	}

	req, err := NewRequest("POST", "http://127.0.0.1", request)
	tt.Nil(err)
	tt.Nil(req.ParseForm())

	tt.Equal([]string{"-1", "1"}, req.Form["a"])
	tt.Equal([]string{"1"}, req.Form["c"])
}

type FormDataMultipart struct {
	Bytes []byte     `name:"bytes"`
	A     []int      `name:"a"`
	C     uint       `name:"c" `
	Data  StructBody `name:"data"`

	File  *multipart.FileHeader   `name:"file"`
	Files []*multipart.FileHeader `name:"files"`
}

func TestTransToHttpReqWithMultipart(t *testing.T) {
	tt := assert.New(t)

	type Req struct {
		FormDataMultipart `in:"formData,multipart"`
	}

	request := &Req{
		FormDataMultipart: FormDataMultipart{
			A:     []int{-1, 1},
			C:     1,
			Bytes: []byte("bytes"),
			Data: StructBody{
				A: -1,
				C: 1,
				D: true,
				E: "d",
			},
		},
	}

	fileHeader, err := NewFileHeader("file", "test.txt", []byte("test test"))
	tt.NoError(err)
	request.File = fileHeader

	fileHeader0, err := NewFileHeader("files", "test0.txt", []byte("test0 test0"))
	tt.NoError(err)
	request.Files = append(request.Files, fileHeader0)

	fileHeader1, err := NewFileHeader("files", "test1.txt", []byte("test1 test1"))
	tt.NoError(err)
	request.Files = append(request.Files, fileHeader1)

	req, err := NewRequest("POST", "http://127.0.0.1", request)
	tt.NoError(err)

	{
		requestForReceive := &Req{}
		err := MarshalParametersWithPath(req, "/", requestForReceive)
		tt.NoError(err)
		tt.Equal(request, requestForReceive)
	}

	tt.Equal("1", req.FormValue("c"))
	tt.Equal([]string{"-1", "1"}, req.Form["a"])
	tt.Equal("bytes", req.FormValue("bytes"))
	tt.Equal(`{"a":-1,"c":1,"d":true,"e":"d"}`, req.FormValue("data"))

	_, _, errForRead := req.FormFile("files")
	tt.NoError(errForRead)
	tt.Equal(fileHeader, req.MultipartForm.File["file"][0])
	tt.Equal([]*multipart.FileHeader{fileHeader0, fileHeader1}, req.MultipartForm.File["files"])
}

func TestTransToHttpReqWithMultipartForReadFailed(t *testing.T) {
	tt := assert.New(t)

	{
		type Req struct {
			FormData struct {
				File *multipart.FileHeader `name:"file"`
			} `in:"formData,multipart"`
		}

		request := &Req{}

		fileHeader, err := NewFileHeader("file", "test.txt", nil)
		tt.NoError(err)
		request.FormData.File = fileHeader

		req, err := NewRequest("POST", "http://127.0.0.1", request)
		tt.NoError(err)

		// destroy body
		bodyBytes, err := ioutil.ReadAll(req.Body)
		tt.NoError(err)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0 : len(bodyBytes)-5]))

		{
			requestForReceive := &Req{}
			err := MarshalParametersWithPath(req, "/", requestForReceive)
			tt.Error(status_error.ReadFormFileFailed, status_error.FromError(err).Code)
		}
	}

	{
		type Req struct {
			FormData struct {
				Files []*multipart.FileHeader `name:"files"`
			} `in:"formData,multipart"`
		}

		request := &Req{}

		fileHeader, err := NewFileHeader("files", "test.txt", nil)
		tt.NoError(err)
		request.FormData.Files = append(request.FormData.Files, fileHeader)

		req, err := NewRequest("POST", "http://127.0.0.1", request)
		tt.NoError(err)

		// destroy body
		bodyBytes, err := ioutil.ReadAll(req.Body)
		tt.NoError(err)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0 : len(bodyBytes)-5]))

		{
			requestForReceive := &Req{}
			err := MarshalParametersWithPath(req, "/", requestForReceive)
			tt.Error(err)
			tt.Equal(int64(status_error.ReadFormFileFailed), status_error.FromError(err).Code)
		}
	}
}
