package translator

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// -----------------------------------------------------------------------------
// Translate Utilities - Gateway
// -----------------------------------------------------------------------------

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayapi.HTTPHeaderMatch) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		if _, exists := convertedHeaders[string(header.Name)]; exists {
			return nil, fmt.Errorf("multiple header matches for the same header are not allowed: %s",
				string(header.Name))
		}
		switch {
		case header.Type != nil && *header.Type == gatewayapi.HeaderMatchRegularExpression:
			convertedHeaders[string(header.Name)] = []string{kongHeaderRegexPrefix + header.Value}
		case header.Type == nil || *header.Type == gatewayapi.HeaderMatchExact:
			convertedHeaders[string(header.Name)] = []string{header.Value}
		default:
			return nil, fmt.Errorf("unknown/unsupported header match type: %s", string(*header.Type))
		}
	}

	return convertedHeaders, nil
}

// generateKongServiceFromBackendRefWithName translates backendRefs into a Kong service for use with the
// rules generated from a Gateway APIs route. The service name is provided by the caller.
func generateKongServiceFromBackendRefWithName(
	logger logr.Logger,
	storer store.Storer,
	rules *ingressRules,
	serviceName string,
	route client.Object,
	protocol string,
	backendRefs ...gatewayapi.BackendRef,
) (kongstate.Service, error) {
	objName := fmt.Sprintf("%s %s/%s",
		route.GetObjectKind().GroupVersionKind().String(), route.GetNamespace(), route.GetName())
	grants, err := storer.ListReferenceGrants()
	if err != nil {
		return kongstate.Service{}, fmt.Errorf("could not retrieve ReferenceGrants for %s: %w", objName, err)
	}
	allowed := gatewayapi.GetPermittedForReferenceGrantFrom(
		logger,
		gatewayapi.ReferenceGrantFrom{
			Group:     gatewayapi.Group(route.GetObjectKind().GroupVersionKind().Group),
			Kind:      gatewayapi.Kind(route.GetObjectKind().GroupVersionKind().Kind),
			Namespace: gatewayapi.Namespace(route.GetNamespace()),
		},
		grants,
	)

	backends := backendRefsToKongStateBackends(logger, storer, route, backendRefs, allowed)

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
	logger logr.Logger,
	storer store.Storer,
	rules *ingressRules,
	route client.Object,
	ruleNumber int,
	protocol string,
	backendRefs ...gatewayapi.BackendRef,
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

func applyExpressionToIngressRules(result *ingressRules) {
	for _, svc := range result.ServiceNameToServices {
		for i := range svc.Routes {
			subtranslator.ApplyExpressionToL4KongRoute(&svc.Routes[i])
			svc.Routes[i].Destinations = nil
			svc.Routes[i].SNIs = nil
		}
	}
}
