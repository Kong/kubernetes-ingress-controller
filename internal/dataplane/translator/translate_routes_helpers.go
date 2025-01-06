package translator

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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
	// Route fulfills client.Object interface (necessary for e.g. GetObjectKind()).
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
			Route:   routeToKongRoute(route, gwPorts, ruleNumber),
		},
	}, nil
}

// routeToKongRoute converts Gateway Route to kong.Route.
func routeToKongRoute[TRoute tTCPorUDPorTLSRoute](
	r TRoute,
	gwPorts []gatewayapi.PortNumber,
	ruleNumber int,
) kong.Route {
	destinations := lo.Map(gwPorts, func(p gatewayapi.PortNumber, _ int) *kong.CIDRPort {
		return &kong.CIDRPort{
			Port: kong.Int(int(p)),
		}
	})

	switch rr := any(r).(type) {
	case *gatewayapi.UDPRoute:
		return udpRouteToKongRoute(rr, destinations, ruleNumber)
	case *gatewayapi.TCPRoute:
		return tcpRouteToKongRoute(rr, destinations, ruleNumber)
	case *gatewayapi.TLSRoute:
		return tlsRouteToKongRoute(rr, ruleNumber)
	default:
		return kong.Route{}
	}
}

func udpRouteToKongRoute(
	r *gatewayapi.UDPRoute,
	destinations []*kong.CIDRPort,
	ruleNumber int,
) kong.Route {
	var ruleName string
	if ruleNumber < len(r.Spec.Rules) {
		ruleName = string(lo.FromPtrOr(r.Spec.Rules[ruleNumber].Name, ""))
	}
	tags := util.GenerateTagsForObject(
		r,
		util.AdditionalTagsK8sNamedRouteRule(ruleName)...,
	)
	return kong.Route{
		Name: kong.String(
			generateRouteName(udpRouteType, r.Namespace, r.Name, ruleNumber),
		),
		Protocols:    kong.StringSlice("udp"),
		Destinations: destinations,
		Tags:         tags,
	}
}

func tcpRouteToKongRoute(
	r *gatewayapi.TCPRoute,
	destinations []*kong.CIDRPort,
	ruleNumber int,
) kong.Route {
	var ruleName string
	if ruleNumber < len(r.Spec.Rules) {
		ruleName = string(lo.FromPtrOr(r.Spec.Rules[ruleNumber].Name, ""))
	}
	tags := util.GenerateTagsForObject(
		r,
		util.AdditionalTagsK8sNamedRouteRule(ruleName)...,
	)
	return kong.Route{
		Name: kong.String(
			generateRouteName(tcpRouteType, r.Namespace, r.Name, ruleNumber),
		),
		Protocols:    kong.StringSlice("tcp"),
		Destinations: destinations,
		Tags:         tags,
	}
}

func tlsRouteToKongRoute(r *gatewayapi.TLSRoute, ruleNumber int) kong.Route {
	hostnames := make([]*string, 0, len(r.Spec.Hostnames))
	for _, hostname := range r.Spec.Hostnames {
		hostnames = append(hostnames, kong.String(string(hostname)))
	}
	var ruleName string
	if ruleNumber < len(r.Spec.Rules) {
		ruleName = string(lo.FromPtrOr(r.Spec.Rules[ruleNumber].Name, ""))
	}
	tags := util.GenerateTagsForObject(
		r,
		util.AdditionalTagsK8sNamedRouteRule(ruleName)...,
	)
	return kong.Route{
		Name: kong.String(
			generateRouteName(tlsRouteType, r.Namespace, r.Name, ruleNumber)),
		Protocols: kong.StringSlice("tls"),
		SNIs:      hostnames,
		Tags:      tags,
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

func (t *Translator) getGatewayListeningPorts(
	routeNamespace string,
	protocol gatewayapi.ProtocolType,
	prs []gatewayapi.ParentReference,
) []gatewayapi.PortNumber {
	var gwPorts []gatewayapi.PortNumber
	for _, pr := range prs {
		// When namespace is explicitly specified in the parentRef,
		// it should be used instead of namespace of the whole Route.
		ns := string(lo.FromPtr(pr.Namespace))
		if ns == "" {
			ns = routeNamespace
		}
		gw, err := t.storer.GetGateway(ns, string(pr.Name))
		if err != nil {
			continue // Skip when attached Gateway is not found.
		}

		// Get explicitly referenced Gateway listening ports by ParentReference configuration.
		// If no sectionName is specified, all ports are used (according to the specification
		// "When unspecified (empty string), this will reference the entire resource." - see
		// https://github.com/kubernetes-sigs/gateway-api/blob/ebe9f31ef27819c3b29f698a3e9b91d279453c59/apis/v1/shared_types.go#L107).
		gwPorts = append(gwPorts, lo.FilterMap(gw.Spec.Listeners, func(l gatewayapi.Listener, _ int) (gatewayapi.PortNumber, bool) {
			if (pr.SectionName == nil || *pr.SectionName == l.Name) && protocol == l.Protocol {
				return l.Port, true
			}
			return 0, false
		})...)
	}
	return gwPorts
}
