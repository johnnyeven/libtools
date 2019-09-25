package transport_http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/johnnyeven/libtools/courier/httpx"
)

func setCORS(headers *http.Header, req *http.Request) {
	referer, err := url.Parse(req.Referer())
	credentials := "false"
	origin := "*"
	if err == nil {
		credentials = "true"
		origin = fmt.Sprintf("%s://%s:%s", referer.Scheme, referer.Hostname(), referer.Port())
	}
	headers.Set("Access-Control-Allow-Credentials", credentials)
	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Access-Control-Allow-Methods", strings.Join([]string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodHead,
		http.MethodDelete,
		http.MethodPatch,
	}, ","))
	headers.Set("Access-Control-Allow-Headers", strings.Join([]string{
		"Origin",
		httpx.HeaderContentType,
		"Content-Length",
		"Authorization",
		"AppToken",
		"AccessKey",
	}, ","))
	headers.Set("Access-Control-Max-Age", strconv.FormatInt(int64(12*time.Hour/time.Second), 10))
}
