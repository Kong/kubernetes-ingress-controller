package parser

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Translate Gateway Routes - Utils
// -----------------------------------------------------------------------------

type tTCPorUDPorTLSRoute interface {
	// Only accept TCP, UDP and TLS routes
	*gatewayv1alpha2.UDPRoute | *gatewayv1alpha2.TCPRoute | *gatewayv1alpha2.TLSRoute
}

type tRoute interface {
	tTCPorUDPorTLSRoute
	// Route fullfills client.Object interface (necessary for e.g. GetObjectKind()).
	client.Object
}

type tRouteRule interface {
	gatewayv1alpha2.UDPRouteRule | gatewayv1alpha2.TCPRouteRule | gatewayv1alpha2.TLSRouteRule
}

// getBackendRefs returns BackendRefs from TCPRouteRule, UDPRouteRule or TLSRouteRule.
func getBackendRefs[TRouteRule tRouteRule](t TRouteRule) ([]gatewayv1alpha2.BackendRef, error) {
	// This is necessary because as of go1.18 (and go1.19) one cannot use common
	// struct fields in generic code.
	//
	// Related golang issue: https://github.com/golang/go/issues/48522
	switch tt := any(t).(type) {
	case gatewayv1alpha2.UDPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return nil, errors.New("UDPRoute rules must include at least one backendRef")
		}
		return tt.BackendRefs, nil
	case gatewayv1alpha2.TCPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return nil, errors.New("TCPRoute rules must include at least one backendRef")
		}
		return tt.BackendRefs, nil
	case gatewayv1alpha2.TLSRouteRule:
		// TLSRoutes don't require BackendRefs.
		return tt.BackendRefs, nil
	}

	// This should never happen because we use type constraints on what types
	// are accepted.
	return nil, nil
}

// generateKongRoutesFromRouteRule converts a Gateway Route (TCP, UDP or TLS) rule
// to one or more Kong Route objects to route traffic to services.
func generateKongRoutesFromRouteRule[T tRoute, TRule tRouteRule](
	route T,
	ruleNumber int,
	rule TRule,
) ([]kongstate.Route, error) {
	backendRefs, err := getBackendRefs(rule)
	if err != nil {
		return []kongstate.Route{}, err
	}

	tags := util.GenerateTagsForObject(route)
	return []kongstate.Route{
		{
			Ingress: util.FromK8sObject(route),
			Route:   routeToKongRoute(route, backendRefs, ruleNumber, tags),
		},
	}, nil
}

// routeToKongRoute converts Gateway Route to kong.Route.
func routeToKongRoute[TRoute tTCPorUDPorTLSRoute](
	r TRoute,
	backendRefs []gatewayv1alpha2.BackendRef,
	ruleNumber int,
	tags []*string,
) kong.Route {
	var kr kong.Route
	switch rr := any(r).(type) {
	case *gatewayv1alpha2.UDPRoute:
		kr = udpRouteToKongRoute(rr, backendRefs, ruleNumber)
	case *gatewayv1alpha2.TCPRoute:
		kr = tcpRouteToKongRoute(rr, backendRefs, ruleNumber)
	case *gatewayv1alpha2.TLSRoute:
		kr = tlsRouteToKongRoute(rr, ruleNumber)
	default:
		kr = kong.Route{}
	}

	kr.Tags = tags
	return kr
}

func udpRouteToKongRoute(
	r *gatewayv1alpha2.UDPRoute,
	backendRefs []gatewayv1alpha2.BackendRef,
	ruleNumber int,
) kong.Route {
	return kong.Route{
		Name: kong.String(
			generateRouteName(udpRouteType, r.Namespace, r.Name, ruleNumber)),
		Protocols:    kong.StringSlice("udp"),
		Destinations: backendRefsToKongCIDRPorts(backendRefs),
	}
}

func tcpRouteToKongRoute(
	r *gatewayv1alpha2.TCPRoute,
	backendRefs []gatewayv1alpha2.BackendRef,
	ruleNumber int,
) kong.Route {
	return kong.Route{
		Name: kong.String(
			generateRouteName(tcpRouteType, r.Namespace, r.Name, ruleNumber)),
		Protocols:    kong.StringSlice("tcp"),
		Destinations: backendRefsToKongCIDRPorts(backendRefs),
	}
}

func backendRefsToKongCIDRPorts(backendRefs []gatewayv1alpha2.BackendRef) []*kong.CIDRPort {
	destinations := make([]*kong.CIDRPort, 0, len(backendRefs))
	for _, backendRef := range backendRefs {
		if backendRef.Port == nil {
			continue // Should we propagate the error?
		}

		destinations = append(destinations,
			&kong.CIDRPort{
				Port: kong.Int(int(*backendRef.Port)),
			},
		)
	}
	return destinations
}

func tlsRouteToKongRoute(r *gatewayv1alpha2.TLSRoute, ruleNumber int) kong.Route {
	hostnames := make([]*string, 0, len(r.Spec.Hostnames))
	for _, hostname := range r.Spec.Hostnames {
		hostnames = append(hostnames, kong.String(string(hostname)))
	}

	return kong.Route{
		Name: kong.String(
			generateRouteName(tlsRouteType, r.Namespace, r.Name, ruleNumber)),
		Protocols: kong.StringSlice("tls"),
		SNIs:      hostnames,
	}
}

type routeType string

const (
	tlsRouteType routeType = "tlsroute"
	tcpRouteType routeType = "tcproute"
	udpRouteType routeType = "udproute"
)

// generateRouteName returns a route name for kong.Route given the provided params.
func generateRouteName(typ routeType, namespace, name string, ruleNumber int) string {
	return fmt.Sprintf(
		"%s.%s.%s.%d.%d",
		typ,
		namespace,
		name,
		ruleNumber,
		0,
	)
}
