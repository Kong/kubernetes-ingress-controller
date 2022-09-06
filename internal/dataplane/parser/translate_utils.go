package parser

import (
	"fmt"
	"reflect"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// kongHeaderRegexPrefix is a reserved prefix string that Kong uses to determine if it should parse a header value
// as a regex.
const kongHeaderRegexPrefix = "~*"

// MinRegexHeaderKongVersion is the minimum Kong version that supports regex header matches.
var MinRegexHeaderKongVersion = semver.MustParse("2.8.0")

// -----------------------------------------------------------------------------
// Translate Utilities - Gateway
// -----------------------------------------------------------------------------

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayv1beta1.HTTPHeaderMatch) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		if _, exists := convertedHeaders[string(header.Name)]; exists {
			return nil, fmt.Errorf("multiple header matches for the same header are not allowed: %s",
				string(header.Name))
		}
		if header.Type != nil && *header.Type == gatewayv1beta1.HeaderMatchRegularExpression {
			if util.GetKongVersion().LT(MinRegexHeaderKongVersion) {
				return nil, fmt.Errorf("Kong version %s does not support HeaderMatchRegularExpression",
					util.GetKongVersion().String())
			}
			convertedHeaders[string(header.Name)] = []string{kongHeaderRegexPrefix + header.Value}
		} else if header.Type == nil || *header.Type == gatewayv1beta1.HeaderMatchExact {
			convertedHeaders[string(header.Name)] = []string{header.Value}
		} else {
			return nil, fmt.Errorf("unknown/unsupported header match type: %s", string(*header.Type))
		}
	}

	return convertedHeaders, nil
}

// getPermittedForReferenceGrantFrom takes a ReferenceGrant From (a namespace, group, and kind) and returns a map
// from a namespace to a slice of ReferenceGrant Tos. When a To is included in the slice, the key namespace has a
// ReferenceGrant with those Tos and the input From.
func getPermittedForReferenceGrantFrom(
	from gatewayv1alpha2.ReferenceGrantFrom,
	grants []*gatewayv1alpha2.ReferenceGrant,
) map[gatewayv1beta1.Namespace][]gatewayv1alpha2.ReferenceGrantTo {
	allowed := make(map[gatewayv1beta1.Namespace][]gatewayv1alpha2.ReferenceGrantTo)
	// loop over all From values in all grants. if we find a match, add all Tos to the list of Tos allowed for the
	// grant namespace. this technically could add duplicate copies of the Tos if there are duplicate Froms (it makes
	// no sense to add them, but it's allowed), but duplicate Tos are harmless (we only care about having at least one
	// matching To when checking if a ReferenceGrant allows a reference)
	for _, grant := range grants {
		for _, otherFrom := range grant.Spec.From {
			if reflect.DeepEqual(from, otherFrom) {
				allowed[gatewayv1beta1.Namespace(grant.ObjectMeta.Namespace)] = append(allowed[gatewayv1beta1.Namespace(grant.ObjectMeta.Namespace)], grant.Spec.To...)
			}
		}
	}

	return allowed
}

// generateKongServiceFromBackendRef translates backendRefs for rule ruleNumber into a Kong service for use with the
// rules generated from a Gateway APIs route.
func generateKongServiceFromBackendRef[
	T types.BackendRefT,
](
	logger logrus.FieldLogger,
	storer store.Storer,
	rules *ingressRules,
	route client.Object,
	ruleNumber int,
	protocol string,
	backendRefs ...T,
) (kongstate.Service, error) {
	objName := fmt.Sprintf("%s %s/%s",
		route.GetObjectKind().GroupVersionKind().String(), route.GetNamespace(), route.GetName())
	if len(backendRefs) == 0 {
		return kongstate.Service{}, fmt.Errorf("no backendRefs present for %s, cannot build Kong service", objName)
	}

	grants, err := storer.ListReferenceGrants()
	if err != nil {
		return kongstate.Service{}, fmt.Errorf("could not retrieve ReferenceGrants for %s: %w", objName, err)
	}
	allowed := getPermittedForReferenceGrantFrom(gatewayv1alpha2.ReferenceGrantFrom{
		Group:     gatewayv1alpha2.Group(route.GetObjectKind().GroupVersionKind().Group),
		Kind:      gatewayv1alpha2.Kind(route.GetObjectKind().GroupVersionKind().Kind),
		Namespace: gatewayv1alpha2.Namespace(route.GetNamespace()),
	}, grants)

	backends := backendRefsToKongStateBackends(logger, route, backendRefs, allowed)

	// the service name needs to uniquely identify this service given it's list of
	// one or more backends.
	serviceName := fmt.Sprintf("%s.%d", getUniqueKongServiceNameForObject(route), ruleNumber)

	// the service host needs to be a resolvable name due to legacy logic so we'll
	// use the anchor backendRef as the basis for the name
	serviceHost := serviceName

	// check if the service is already known, and if not create it
	service, ok := rules.ServiceNameToServices[serviceName]
	if !ok {
		service = kongstate.Service{
			Service: kong.Service{
				Name:           kong.String(serviceName),
				Host:           kong.String(serviceHost),
				Protocol:       kong.String(protocol),
				ConnectTimeout: kong.Int(DefaultServiceTimeout),
				ReadTimeout:    kong.Int(DefaultServiceTimeout),
				WriteTimeout:   kong.Int(DefaultServiceTimeout),
				Retries:        kong.Int(DefaultRetries),
			},
			Namespace: route.GetNamespace(),
			Backends:  backends,
			Parent:    route,
		}
	}

	// In the context of the gateway API conformance tests, if there is no service for the backend,
	// the response must have a status code of 500. Since The default behavior of Kong is returning 503
	// if there is no backend for a service, we inject a plugin that terminates all requests with 500
	// as status code
	if len(service.Backends) == 0 {
		if service.Plugins == nil {
			service.Plugins = make([]kong.Plugin, 0)
		}
		service.Plugins = append(service.Plugins, kong.Plugin{
			Name: kong.String("request-termination"),
			Config: kong.Configuration{
				"status_code": 500,
				"message":     "no existing backendRef provided",
			},
		})
	}

	return service, nil
}
