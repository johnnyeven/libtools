package sign

import (
	"fmt"
	"net/http"

	"github.com/johnnyeven/libtools/courier/transport_http"
)

var _ interface {
	http.RoundTripper
} = (*AutoSignTransport)(nil)

func NewAutoSignTransport(secretExchanger SecretExchanger) transport_http.TransportWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		if httpRt, ok := rt.(*http.Transport); ok {
			autoSignTransport := &AutoSignTransport{
				SecretExchanger: secretExchanger,
				Transport:       httpRt,
			}
			return autoSignTransport
		}
		return rt
	}
}

type AutoSignTransport struct {
	SecretExchanger SecretExchanger
	*http.Transport
}

func (t *AutoSignTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.SecretExchanger == nil {
		return nil, fmt.Errorf("missing SecretExchanger")
	}

	query := req.URL.Query()

	sign, _, err := getSign(req, query, t.SecretExchanger)
	if err != nil {
		return nil, err
	}
	if sign != nil {
		query.Set(Sign, string(sign))
		req.URL.RawQuery = query.Encode()
	}

	return t.Transport.RoundTrip(req)
}
