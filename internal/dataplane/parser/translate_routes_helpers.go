package parser

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/gatewayapi"
)

// -----------------------------------------------------------------------------
// Translate Gateway Routes - Utils
// -----------------------------------------------------------------------------

type tTCPorUDPorTLSRoute interface {
	// Only accept TCP, UDP and TLS routes
	*gatewayapi.UDPRoute | *gatewayapi.TCPRoute | *gatewayapi.TLSRoute
}

type tRoute interface {
	tTCPorUDPorTLSRoute
	// Route fullfills client.Object interface (necessary for e.g. GetObjectKind()).
	client.Object
}

type tRouteRule interface {
	gatewayapi.UDPRouteRule | gatewayapi.TCPRouteRule | gatewayapi.TLSRouteRule
}

// getBackendRefs returns BackendRefs from TCPRouteRule, UDPRouteRule or TLSRouteRule.
func getBackendRefs[TRouteRule tRouteRule](t TRouteRule) ([]gatewayapi.BackendRef, error) {
	// This is necessary because as of go1.18 (and go1.19) one cannot use common
	// struct fields in generic code.
	//
	// Related golang issue: https://github.com/golang/go/issues/48522
	switch tt := any(t).(type) {
	case gatewayapi.UDPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return nil, errors.New("UDPRoute rules must include at least one backendRef")
		}
		return tt.BackendRefs, nil
	case gatewayapi.TCPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return nil, errors.New("TCPRoute rules must include at least one backendRef")
		}
		return tt.BackendRefs, nil
	case gatewayapi.TLSRouteRule:
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
	backendRefs []gatewayapi.BackendRef,
	ruleNumber int,
	tags []*string,
) kong.Route {
	var kr kong.Route
	switch rr := any(r).(type) {
	case *gatewayapi.UDPRoute:
		kr = udpRouteToKongRoute(rr, backendRefs, ruleNumber)
	case *gatewayapi.TCPRoute:
		kr = tcpRouteToKongRoute(rr, backendRefs, ruleNumber)
	case *gatewayapi.TLSRoute:
		kr = tlsRouteToKongRoute(rr, ruleNumber)
	default:
		kr = kong.Route{}
	}

	kr.Tags = tags
	return kr
}

func udpRouteToKongRoute(
	r *gatewayapi.UDPRoute,
	backendRefs []gatewayapi.BackendRef,
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
	r *gatewayapi.TCPRoute,
	backendRefs []gatewayapi.BackendRef,
	ruleNumber int,
) kong.Route {
	return kong.Route{
		Name: kong.String(
			generateRouteName(tcpRouteType, r.Namespace, r.Name, ruleNumber)),
		Protocols:    kong.StringSlice("tcp"),
		Destinations: backendRefsToKongCIDRPorts(backendRefs),
	}
}

func backendRefsToKongCIDRPorts(backendRefs []gatewayapi.BackendRef) []*kong.CIDRPort {
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

func tlsRouteToKongRoute(r *gatewayapi.TLSRoute, ruleNumber int) kong.Route {
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
