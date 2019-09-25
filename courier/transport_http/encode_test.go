package transport_http

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/courier/transport_http/transform"
)

func NewTestRespWriter() *TestRespWriter {
	return &TestRespWriter{
		H: http.Header{},
	}
}

type TestRespWriter struct {
	bytes.Buffer
	Status int
	H      http.Header
}

func (rw *TestRespWriter) Header() http.Header {
	return rw.H
}

func (rw *TestRespWriter) WriteHeader(code int) {
	rw.Status = code
}

func TestEncodeHttpResponse_NoContent(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()
	req, _ := transform.NewRequest("GET", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, nil)
	tt.NoError(err)
	tt.Equal(http.StatusNoContent, rw.Status)
}

func TestEncodeHttpResponse_WithMetaAndContentType(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()
	req, _ := transform.NewRequest("GET", "/", nil)

	file := NewFile("123.txt", "text/plain")
	file.Write([]byte("123123"))

	err := encodeHttpResponse(context.Background(), rw, req, file)
	tt.NoError(err)
	tt.Equal(http.StatusOK, rw.Status)
	tt.Equal([]byte("123123"), rw.Bytes())
	tt.Equal("attachment; filename=123.txt", rw.Header().Get("Content-Disposition"))
	tt.Equal("text/plain;charset=utf-8", rw.Header().Get("Content-Type"))
}

func TestEncodeHttpResponse_WithStatus(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()
	req, _ := transform.NewRequest("GET", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, status_error.InvalidStruct.StatusError())
	tt.NoError(err)
	tt.Equal(status_error.InvalidStruct.StatusError().Status(), rw.Status)
}

func TestEncodeHttpResponse_SomeJSONForGet(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()

	respData := struct {
		A string `json:"a"`
		B string `json:"b"`
	}{
		A: "a",
		B: "b",
	}

	req, _ := transform.NewRequest("GET", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, respData)
	tt.NoError(err)
	tt.Equal(http.StatusOK, rw.Status)
	tt.Equal(`{"a":"a","b":"b"}`, rw.String())
}

type XMLData struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (XMLData) ContentType() string {
	return httpx.MIMEXML
}

func TestEncodeHttpResponse_SomeXMLForGet(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()

	respData := XMLData{
		A: "a",
		B: "b",
	}

	req, _ := transform.NewRequest("GET", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, respData)
	tt.NoError(err)
	tt.Equal(http.StatusOK, rw.Status)
	tt.Equal(`<XMLData><A>a</A><B>b</B></XMLData>`, rw.String())
}

func TestEncodeHttpResponse_SomeJSONForPOST(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()

	respData := struct {
		A string `json:"a"`
		B string `json:"b"`
	}{
		A: "a",
		B: "b",
	}

	req, _ := transform.NewRequest("POST", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, respData)
	tt.NoError(err)
	tt.Equal(http.StatusCreated, rw.Status)
	tt.Equal(`{"a":"a","b":"b"}`, rw.String())
}

func TestEncodeHttpResponse_ByteDirectly(t *testing.T) {
	tt := assert.New(t)
	rw := NewTestRespWriter()

	req, _ := transform.NewRequest("POST", "/", nil)

	err := encodeHttpResponse(context.Background(), rw, req, []byte("123"))
	tt.NoError(err)
	tt.Equal(http.StatusCreated, rw.Status)
	tt.Equal("123", rw.String())
}
