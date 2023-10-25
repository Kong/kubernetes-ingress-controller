package parser

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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

// checkBackendRefs checks if BackendRefs are properly configured for
// TCPRouteRule, UDPRouteRule or TLSRouteRule.
func checkBackendRefs[TRouteRule tRouteRule](t TRouteRule) error {
	// This is necessary because as of go1.18 (and go1.19) one cannot use common
	// struct fields in generic code.
	//
	// Related golang issue: https://github.com/golang/go/issues/48522
	switch tt := any(t).(type) {
	case gatewayapi.UDPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return errors.New("UDPRoute rules must include at least one backendRef")
		}
	case gatewayapi.TCPRouteRule:
		if len(tt.BackendRefs) == 0 {
			return errors.New("TCPRoute rules must include at least one backendRef")
		}
	case gatewayapi.TLSRouteRule:
		// TLSRoutes don't require BackendRefs.
	}
	return nil
}

// generateKongRoutesFromRouteRule converts a Gateway Route (TCP, UDP or TLS) rule
// to one or more Kong Route objects to route traffic to services.
func generateKongRoutesFromRouteRule[T tRoute, TRule tRouteRule](
	route T,
	gwPorts []gatewayapi.PortNumber,
	ruleNumber int,
	rule TRule,
) ([]kongstate.Route, error) {
	if err := checkBackendRefs(rule); err != nil {
		return []kongstate.Route{}, err
	}
	return []kongstate.Route{
		{
			Ingress: util.FromK8sObject(route),
			Route:   routeToKongRoute(route, gwPorts, ruleNumber, util.GenerateTagsForObject(route)),
		},
	}, nil
}

// routeToKongRoute converts Gateway Route to kong.Route.
func routeToKongRoute[TRoute tTCPorUDPorTLSRoute](
	r TRoute,
	gwPorts []gatewayapi.PortNumber,
	ruleNumber int,
	tags []*string,
) kong.Route {
	destinations := lo.Map(gwPorts, func(p gatewayapi.PortNumber, _ int) *kong.CIDRPort {
		return &kong.CIDRPort{
			Port: kong.Int(int(p)),
		}
	})

	var kr kong.Route
	switch rr := any(r).(type) {
	case *gatewayapi.UDPRoute:
		kr = udpRouteToKongRoute(rr, destinations, ruleNumber)
	case *gatewayapi.TCPRoute:
		kr = tcpRouteToKongRoute(rr, destinations, ruleNumber)
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
	destinations []*kong.CIDRPort,
	ruleNumber int,
) kong.Route {
	return kong.Route{
		Name: kong.String(
			generateRouteName(udpRouteType, r.Namespace, r.Name, ruleNumber),
		),
		Protocols:    kong.StringSlice("udp"),
		Destinations: destinations,
	}
}

func tcpRouteToKongRoute(
	r *gatewayapi.TCPRoute,
	destinations []*kong.CIDRPort,
	ruleNumber int,
) kong.Route {
	return kong.Route{
		Name: kong.String(
			generateRouteName(tcpRouteType, r.Namespace, r.Name, ruleNumber),
		),
		Protocols:    kong.StringSlice("tcp"),
		Destinations: destinations,
	}
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
