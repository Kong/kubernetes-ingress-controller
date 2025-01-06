package useragent

import (
	"net/http"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

func NewTransport(underlyingTransport http.RoundTripper) http.RoundTripper {
	return &transport{
		underlyingTransport: underlyingTransport,
	}
}

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", metadata.UserAgent())
	if t.underlyingTransport != nil {
		return t.underlyingTransport.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}
