package subtranslator

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func getGRPCMatchDefaults() (
	map[gatewayapi.GRPCMethodMatchType]string,
	map[gatewayapi.GRPCMethodMatchType]string,
) {
	// Kong routes derived from a GRPCRoute use a path composed of the match's gRPC service and method
	// If either the service or method is omitted, there is a default regex determined by the match type
	// https://gateway-api.sigs.k8s.io/geps/gep-1016/#matcher-types describes the defaults

	// default path components for the GRPC service
	return map[gatewayapi.GRPCMethodMatchType]string{
			gatewayapi.GRPCMethodMatchType(""):          ".+",
			gatewayapi.GRPCMethodMatchExact:             ".+",
			gatewayapi.GRPCMethodMatchRegularExpression: ".+",
		},
		// default path components for the GRPC method
		map[gatewayapi.GRPCMethodMatchType]string{
			gatewayapi.GRPCMethodMatchType(""):          "",
			gatewayapi.GRPCMethodMatchExact:             "",
			gatewayapi.GRPCMethodMatchRegularExpression: ".+",
		}
}

func GenerateKongRoutesFromGRPCRouteRule(
	grpcroute *gatewayapi.GRPCRoute,
	ruleNumber int,
	storer store.Storer,
) []kongstate.Route {
	if ruleNumber >= len(grpcroute.Spec.Rules) {
		return nil
	}

	routeName := func(namespace string, name string, ruleNumber int, matchNumber int) *string {
		return kong.String(fmt.Sprintf(
			"grpcroute.%s.%s.%d.%d",
			namespace,
			name,
			ruleNumber,
			matchNumber,
		))
	}

	// Gather the K8s object information and hostnames from the GRPCRoute.
	ingressObjectInfo := util.FromK8sObject(grpcroute)
	tags := util.GenerateTagsForObject(grpcroute)
	grpcProtocols := kong.StringSlice("grpc", "grpcs")
	rule := grpcroute.Spec.Rules[ruleNumber]
	// Kong Route expects to have for gRPC, at least one of Hosts, Headers or Paths fields set.
	// For no matches it can be a catch-all or route based on hostnames.
	if len(rule.Matches) == 0 {
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:      routeName(grpcroute.Namespace, grpcroute.Name, ruleNumber, 0),
				Protocols: grpcProtocols,
				Tags:      tags,
			},
		}
		if configuredHostnames := getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute, storer); len(configuredHostnames) > 0 {
			// Match based on hostnames.
			r.Hosts = configuredHostnames
		} else {
			// No hostnames configured, so this is a catch-all.
			// https://docs.konghq.com/gateway/latest/production/configuring-a-grpc-service/#single-grpc-service-and-route
			r.Paths = kong.StringSlice("/")
		}
		return []kongstate.Route{r}
	}

	// Rule matches are configured, hostname may be specified too.
	routes := make([]kongstate.Route, 0, len(rule.Matches))
	for matchNumber, match := range rule.Matches {
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:      routeName(grpcroute.Namespace, grpcroute.Name, ruleNumber, matchNumber),
				Protocols: grpcProtocols,
				Tags:      tags,
				Hosts:     getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute, storer),
			},
		}

		if match.Method != nil {
			serviceMap, methodMap := getGRPCMatchDefaults()
			var method, service string
			matchMethod := match.Method.Method
			matchService := match.Method.Service
			var matchType gatewayapi.GRPCMethodMatchType
			if match.Method.Type == nil {
				matchType = gatewayapi.GRPCMethodMatchExact
			} else {
				matchType = *match.Method.Type
			}
			if matchMethod == nil {
				method = methodMap[matchType]
			} else {
				method = *matchMethod
			}
			if matchService == nil {
				service = serviceMap[matchType]
			} else {
				service = *matchService
			}
			path := kong.String(KongPathRegexPrefix + fmt.Sprintf("/%s/%s", service, method))
			r.Paths = append(r.Paths, path)
		}

		r.Headers = map[string][]string{}
		for _, hmatch := range match.Headers {
			name := string(hmatch.Name)
			r.Headers[name] = append(r.Headers[name], hmatch.Value)
		}

		routes = append(routes, r)
	}
	return routes
}

// -----------------------------------------------------------------------------
// Translate GRPCRoute - Utils
// -----------------------------------------------------------------------------

// getGRPCRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an GRPCRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
// The hostname field is optional. If not specified, the configured value will be obtained from parentRefs.
func getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute *gatewayapi.GRPCRoute, storer store.Storer) []*string {
	hostnames := make([]gatewayapi.Hostname, 0)
	// If no hostnames are specified, we will use the hostname from the Gateway
	// that the GRPCRoute is attached to.
	if len(grpcroute.Spec.Hostnames) == 0 {
		namespace := grpcroute.GetNamespace()

		for _, parentRef := range grpcroute.Spec.ParentRefs {
			// we only care about Gateways
			if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
				continue
			}

			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			name := string(parentRef.Name)

			gateway, err := storer.GetGateway(namespace, name)
			// If we got an error, we will just return nil.
			// Users need to check whether the relevant configurations are correct.
			if err != nil {
				return nil
			}

			if parentRef.SectionName != nil {
				sectionName := string(*parentRef.SectionName)

				for _, listener := range gateway.Spec.Listeners {
					if string(listener.Name) == sectionName {
						if listener.Hostname != nil {
							hostnames = append(hostnames, *listener.Hostname)
						}
					}
				}
			} else {
				for _, listener := range gateway.Spec.Listeners {
					if listener.Hostname != nil {
						hostnames = append(hostnames, *listener.Hostname)
					}
				}
			}
		}
	} else {
		hostnames = grpcroute.Spec.Hostnames
	}

	return lo.Map(hostnames, func(h gatewayapi.Hostname, _ int) *string {
		return lo.ToPtr(string(h))
	})
}
