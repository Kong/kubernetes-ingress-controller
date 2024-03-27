package builder

import (
	"github.com/samber/lo"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteFilterBuilder is a builder for gateway api HTTPRouteMatch.
// Primarily used for testing.
type HTTPRouteFilterBuilder struct {
	httpRouteFilter gatewayv1.HTTPRouteFilter
}

func (b *HTTPRouteFilterBuilder) Build() gatewayv1.HTTPRouteFilter {
	return b.httpRouteFilter
}

// NewHTTPRouteRequestRedirectFilter builds a request redirect HTTPRoute filter.
func NewHTTPRouteRequestRedirectFilter() *HTTPRouteFilterBuilder {
	filter := gatewayv1.HTTPRouteFilter{
		Type:            gatewayv1.HTTPRouteFilterRequestRedirect,
		RequestRedirect: &gatewayv1.HTTPRequestRedirectFilter{},
	}
	return &HTTPRouteFilterBuilder{httpRouteFilter: filter}
}

// WithRequestRedirectScheme sets scheme of request redirect filter.
func (b *HTTPRouteFilterBuilder) WithRequestRedirectScheme(scheme string) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestRedirect ||
		b.httpRouteFilter.RequestRedirect == nil {
		return b
	}

	b.httpRouteFilter.RequestRedirect.Scheme = lo.ToPtr(scheme)
	return b
}

// WithRequestRedirectHost sets host of request redirect filter.
func (b *HTTPRouteFilterBuilder) WithRequestRedirectHost(host string) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestRedirect ||
		b.httpRouteFilter.RequestRedirect == nil {
		return b
	}

	preciseHost := (gatewayv1.PreciseHostname)(host)
	b.httpRouteFilter.RequestRedirect.Hostname = lo.ToPtr(preciseHost)
	return b
}

// WithRequestRedirectStatusCode sets status code of response in request redirect filter.
func (b *HTTPRouteFilterBuilder) WithRequestRedirectStatusCode(code int) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestRedirect ||
		b.httpRouteFilter.RequestRedirect == nil {
		return b
	}

	b.httpRouteFilter.RequestRedirect.StatusCode = lo.ToPtr(code)
	return b
}

// NewHTTPRouteRequestHeaderModifierFilter builds a request header modifier HTTPRoute filter.
func NewHTTPRouteRequestHeaderModifierFilter() *HTTPRouteFilterBuilder {
	filter := gatewayv1.HTTPRouteFilter{
		Type:                  gatewayv1.HTTPRouteFilterRequestHeaderModifier,
		RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{},
	}
	return &HTTPRouteFilterBuilder{httpRouteFilter: filter}
}

func (b *HTTPRouteFilterBuilder) WithRequestHeaderAdd(headers []gatewayv1.HTTPHeader) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestHeaderModifier ||
		b.httpRouteFilter.RequestHeaderModifier == nil {
		return b
	}
	b.httpRouteFilter.RequestHeaderModifier.Add = headers
	return b
}

func (b *HTTPRouteFilterBuilder) WithRequestHeaderSet(headers []gatewayv1.HTTPHeader) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestHeaderModifier ||
		b.httpRouteFilter.RequestHeaderModifier == nil {
		return b
	}
	b.httpRouteFilter.RequestHeaderModifier.Set = headers
	return b
}

func (b *HTTPRouteFilterBuilder) WithRequestHeaderRemove(headerNames []string) *HTTPRouteFilterBuilder {
	if b.httpRouteFilter.Type != gatewayv1.HTTPRouteFilterRequestHeaderModifier ||
		b.httpRouteFilter.RequestHeaderModifier == nil {
		return b
	}
	b.httpRouteFilter.RequestHeaderModifier.Remove = headerNames
	return b
}
