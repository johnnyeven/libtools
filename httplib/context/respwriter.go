package context

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

type TestResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	Body            bytes.Buffer
	WHeader         http.Header
	HttpStatus      int
	written         bool
	CloseNotifyChan <-chan bool
}

func (w *TestResponseWriter) Header() http.Header {
	if w.WHeader == nil {
		w.WHeader = make(map[string][]string)
	}
	return w.WHeader
}

func (w *TestResponseWriter) WriteString(data string) (int, error) {
	return w.Write([]byte(data))
}

func (w *TestResponseWriter) Write(data []byte) (int, error) {
	w.written = true
	return w.Body.Write(data)
}

func (w *TestResponseWriter) WriteHeader(status int) {
	w.HttpStatus = status
}

func (w *TestResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (w *TestResponseWriter) Flush() {
}

func (w *TestResponseWriter) CloseNotify() <-chan bool {
	return w.CloseNotifyChan
}

func (w *TestResponseWriter) Status() int {
	return w.HttpStatus
}

func (w *TestResponseWriter) Size() int {
	return w.Body.Len()
}

func (w *TestResponseWriter) Written() bool {
	return w.written
}

func (w *TestResponseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.WriteHeader(w.HttpStatus)
	}
}
