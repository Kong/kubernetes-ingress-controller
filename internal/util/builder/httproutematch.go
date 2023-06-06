package builder

import (
	"github.com/samber/lo"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// HTTPRouteMatchBuilder is a builder for gateway api HTTPRouteMatch.
// Primarily used for testing.
// Please note that some methods are not provided yet, as we
// don't need them yet. Feel free to add them as needed.
type HTTPRouteMatchBuilder struct {
	httpRouteMatch gatewayv1beta1.HTTPRouteMatch
}

func NewHTTPRouteMatch() *HTTPRouteMatchBuilder {
	return &HTTPRouteMatchBuilder{
		httpRouteMatch: gatewayv1beta1.HTTPRouteMatch{},
	}
}

func (b *HTTPRouteMatchBuilder) Build() gatewayv1beta1.HTTPRouteMatch {
	return b.httpRouteMatch
}

func (b *HTTPRouteMatchBuilder) WithPathPrefix(pathPrefix string) *HTTPRouteMatchBuilder {
	return b.WithPathType(&pathPrefix, lo.ToPtr(gatewayv1beta1.PathMatchPathPrefix))
}

func (b *HTTPRouteMatchBuilder) WithPathRegex(pathRegexp string) *HTTPRouteMatchBuilder {
	return b.WithPathType(&pathRegexp, lo.ToPtr(gatewayv1beta1.PathMatchRegularExpression))
}

func (b *HTTPRouteMatchBuilder) WithPathExact(pathRegexp string) *HTTPRouteMatchBuilder {
	return b.WithPathType(&pathRegexp, lo.ToPtr(gatewayv1beta1.PathMatchExact))
}

func (b *HTTPRouteMatchBuilder) WithPathType(pathValuePtr *string, pathTypePtr *gatewayv1beta1.PathMatchType) *HTTPRouteMatchBuilder {
	b.httpRouteMatch.Path = &gatewayv1beta1.HTTPPathMatch{
		Type:  pathTypePtr,
		Value: pathValuePtr,
	}
	return b
}

func (b *HTTPRouteMatchBuilder) WithQueryParam(name, value string) *HTTPRouteMatchBuilder {
	b.httpRouteMatch.QueryParams = append(b.httpRouteMatch.QueryParams, gatewayv1beta1.HTTPQueryParamMatch{
		Name:  gatewayv1beta1.HTTPHeaderName(name),
		Value: value,
	})
	return b
}

func (b *HTTPRouteMatchBuilder) WithMethod(method gatewayv1beta1.HTTPMethod) *HTTPRouteMatchBuilder {
	b.httpRouteMatch.Method = &method
	return b
}

func (b *HTTPRouteMatchBuilder) WithHeader(name, value string) *HTTPRouteMatchBuilder {
	b.httpRouteMatch.Headers = append(b.httpRouteMatch.Headers, gatewayv1beta1.HTTPHeaderMatch{
		Name:  gatewayv1beta1.HTTPHeaderName(name),
		Value: value,
	})
	return b
}

func (b *HTTPRouteMatchBuilder) WithHeaderRegex(name, value string) *HTTPRouteMatchBuilder {
	b.httpRouteMatch.Headers = append(b.httpRouteMatch.Headers, gatewayv1beta1.HTTPHeaderMatch{
		Name:  gatewayv1beta1.HTTPHeaderName(name),
		Value: value,
		Type:  lo.ToPtr(gatewayv1beta1.HeaderMatchRegularExpression),
	})
	return b
}
