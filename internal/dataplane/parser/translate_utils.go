package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

// -----------------------------------------------------------------------------
// Translate Utilities - Gateway
// -----------------------------------------------------------------------------

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayv1.HTTPHeaderMatch, kongVersion semver.Version) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		if _, exists := convertedHeaders[string(header.Name)]; exists {
			return nil, fmt.Errorf("multiple header matches for the same header are not allowed: %s",
				string(header.Name))
		}
		if header.Type != nil && *header.Type == gatewayv1.HeaderMatchRegularExpression {
			if kongVersion.LT(versions.RegexHeaderVersionCutoff) {
				return nil, fmt.Errorf("Kong version %s does not support HeaderMatchRegularExpression", kongVersion)
			}
			convertedHeaders[string(header.Name)] = []string{kongHeaderRegexPrefix + header.Value}
		} else if header.Type == nil || *header.Type == gatewayv1.HeaderMatchExact {
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
	from gatewayv1beta1.ReferenceGrantFrom,
	grants []*gatewayv1beta1.ReferenceGrant,
) map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo {
	allowed := make(map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo)
	// loop over all From values in all grants. if we find a match, add all Tos to the list of Tos allowed for the
	// grant namespace. this technically could add duplicate copies of the Tos if there are duplicate Froms (it makes
	// no sense to add them, but it's allowed), but duplicate Tos are harmless (we only care about having at least one
	// matching To when checking if a ReferenceGrant allows a reference)
	for _, grant := range grants {
		for _, otherFrom := range grant.Spec.From {
			if reflect.DeepEqual(from, otherFrom) {
				allowed[gatewayv1.Namespace(grant.ObjectMeta.Namespace)] = append(allowed[gatewayv1.Namespace(grant.ObjectMeta.Namespace)], grant.Spec.To...)
			}
		}
	}

	return allowed
}

// generateKongServiceFromBackendRefWithName translates backendRefs into a Kong service for use with the
// rules generated from a Gateway APIs route. The service name is provided by the caller.
func generateKongServiceFromBackendRefWithName(
	logger logrus.FieldLogger,
	storer store.Storer,
	rules *ingressRules,
	serviceName string,
	route client.Object,
	protocol string,
	backendRefs ...gatewayv1.BackendRef,
) (kongstate.Service, error) {
	objName := fmt.Sprintf("%s %s/%s",
		route.GetObjectKind().GroupVersionKind().String(), route.GetNamespace(), route.GetName())
	grants, err := storer.ListReferenceGrants()
	if err != nil {
		return kongstate.Service{}, fmt.Errorf("could not retrieve ReferenceGrants for %s: %w", objName, err)
	}
	allowed := getPermittedForReferenceGrantFrom(gatewayv1beta1.ReferenceGrantFrom{
		Group:     gatewayv1alpha2.Group(route.GetObjectKind().GroupVersionKind().Group),
		Kind:      gatewayv1alpha2.Kind(route.GetObjectKind().GroupVersionKind().Kind),
		Namespace: gatewayv1alpha2.Namespace(route.GetNamespace()),
	}, grants)

	backends := backendRefsToKongStateBackends(logger, route, backendRefs, allowed)

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
	if len(service.Backends) == 0 && len(backendRefs) != 0 {
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

// generateKongServiceFromBackendRefWithRuleNumber translates backendRefs for rule ruleNumber into a Kong service for use with the
// rules generated from a Gateway APIs route. The service name is computed from route and ruleNumber by the function.
func generateKongServiceFromBackendRefWithRuleNumber(
	logger logrus.FieldLogger,
	storer store.Storer,
	rules *ingressRules,
	route client.Object,
	ruleNumber int,
	protocol string,
	backendRefs ...gatewayv1.BackendRef,
) (kongstate.Service, error) {
	// the service name needs to uniquely identify this service given it's list of
	// one or more backends.
	serviceName := fmt.Sprintf("%s.%d", getUniqueKongServiceNameForObject(route), ruleNumber)

	return generateKongServiceFromBackendRefWithName(
		logger,
		storer,
		rules,
		serviceName,
		route,
		protocol,
		backendRefs...,
	)
}

// maybePrependRegexPrefix takes a path, controller regex prefix, and a legacy heuristic toggle. It returns the path
// with the Kong regex path prefix if it either began with the controller prefix or did not, but matched the legacy
// heuristic, and the heuristic was enabled.
func maybePrependRegexPrefix(path, controllerPrefix string, applyLegacyHeuristic bool) string {
	if strings.HasPrefix(path, controllerPrefix) {
		path = strings.Replace(path, controllerPrefix, translators.KongPathRegexPrefix, 1)
	} else if applyLegacyHeuristic {
		// this regex matches if the path _is not_ considered a regex by Kong 2.x
		if LegacyRegexPathExpression.FindString(path) == "" {
			if !strings.HasPrefix(path, translators.KongPathRegexPrefix) {
				path = translators.KongPathRegexPrefix + path
			}
		}
	}
	return path
}

func applyExpressionToIngressRules(result *ingressRules) {
	for _, svc := range result.ServiceNameToServices {
		for i := range svc.Routes {
			translators.ApplyExpressionToL4KongRoute(&svc.Routes[i])
			svc.Routes[i].Destinations = nil
			svc.Routes[i].SNIs = nil
		}
	}
}
