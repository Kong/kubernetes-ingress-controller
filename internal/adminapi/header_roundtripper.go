package adminapi

import (
	"net/http"
	"strings"
)

// HeaderRoundTripper injects Headers into requests
// made via RT.
type HeaderRoundTripper struct {
	headers []string
	rt      http.RoundTripper
}

// RoundTrip satisfies the RoundTripper interface.
func (t *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := req.Clone(req.Context())
	for _, s := range t.headers {
		split := strings.SplitN(s, ":", 2)
		if len(split) >= 2 {
			newRequest.Header[split[0]] = append([]string(nil), split[1])
		}
	}
	return t.rt.RoundTrip(newRequest)
}
