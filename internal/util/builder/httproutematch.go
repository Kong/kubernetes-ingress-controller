package builder

import (
	"k8s.io/utils/pointer"
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
	return b.WithPathType(&pathPrefix, (*gatewayv1beta1.PathMatchType)(pointer.StringPtr(string(gatewayv1beta1.PathMatchPathPrefix))))
}

func (b *HTTPRouteMatchBuilder) WithPathRegex(pathRegexp string) *HTTPRouteMatchBuilder {
	return b.WithPathType(&pathRegexp, (*gatewayv1beta1.PathMatchType)(pointer.StringPtr(string(gatewayv1beta1.PathMatchRegularExpression))))
}

func (b *HTTPRouteMatchBuilder) WithPathExact(pathRegexp string) *HTTPRouteMatchBuilder {
	return b.WithPathType(&pathRegexp, (*gatewayv1beta1.PathMatchType)(pointer.StringPtr(string(gatewayv1beta1.PathMatchExact))))
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
		Name:  name,
		Value: value,
	})
	return b
}
