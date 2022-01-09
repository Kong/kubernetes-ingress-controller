package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func serviceBackendPortToStr(port networkingv1.ServiceBackendPort) string {
	if port.Name != "" {
		return fmt.Sprintf("pname-%s", port.Name)
	}
	return fmt.Sprintf("pnum-%d", port.Number)
}

func fromIngressV1beta1(log logrus.FieldLogger, ingressList []*networkingv1beta1.Ingress) ingressRules {
	result := newIngressRules()

	var allDefaultBackends []networkingv1beta1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec
		log = log.WithFields(logrus.Fields{
			"ingress_namespace": ingress.Namespace,
			"ingress_name":      ingress.Name,
		})

		if ingressSpec.Backend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1beta1TLS(ingressSpec.TLS, ingress.Namespace)

		for i, rule := range ingressSpec.Rules {
			host := rule.Host
			if rule.HTTP == nil {
				continue
			}
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path

				if strings.Contains(path, "//") {
					log.Errorf("rule skipped: invalid path: '%v'", path)
					continue
				}
				if path == "" {
					path = "/"
				}
				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
						// TODO (#834) Figure out a way to name the routes
						// This is not a stable scheme
						// 1. If a user adds a route in the middle,
						// due to a shift, all the following routes will
						// be PATCHED
						// 2. Is it guaranteed that the order is stable?
						// Meaning, the routes will always appear in the same
						// order?
						Name:              kong.String(fmt.Sprintf("%s.%s.%d%d", ingress.Namespace, ingress.Name, i, j)),
						Paths:             kong.StringSlice(path),
						StripPath:         kong.Bool(false),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						RequestBuffering:  kong.Bool(true),
						ResponseBuffering: kong.Bool(true),
					},
				}
				if host != "" {
					hosts := kong.StringSlice(host)
					r.Hosts = hosts
				}

				serviceName := ingress.Namespace + "." +
					rule.Backend.ServiceName + "." +
					rule.Backend.ServicePort.String()
				service, ok := result.ServiceNameToServices[serviceName]
				if !ok {
					service = kongstate.Service{
						Service: kong.Service{
							Name: kong.String(serviceName),
							Host: kong.String(rule.Backend.ServiceName +
								"." + ingress.Namespace + "." +
								rule.Backend.ServicePort.String() + ".svc"),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend: kongstate.ServiceBackend{
							Name: rule.Backend.ServiceName,
							Port: PortDefFromIntStr(rule.Backend.ServicePort),
						},
					}
				}
				service.Routes = append(service.Routes, r)
				result.ServiceNameToServices[serviceName] = service
			}
		}
	}

	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	// Process the default backend
	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := allDefaultBackends[0].Spec.Backend
		serviceName := allDefaultBackends[0].Namespace + "." +
			defaultBackend.ServiceName + "." +
			defaultBackend.ServicePort.String()
		service, ok := result.ServiceNameToServices[serviceName]
		if !ok {
			service = kongstate.Service{
				Service: kong.Service{
					Name: kong.String(serviceName),
					Host: kong.String(defaultBackend.ServiceName + "." +
						ingress.Namespace + "." +
						defaultBackend.ServicePort.String() + ".svc"),
					Port:           kong.Int(80),
					Protocol:       kong.String("http"),
					ConnectTimeout: kong.Int(60000),
					ReadTimeout:    kong.Int(60000),
					WriteTimeout:   kong.Int(60000),
					Retries:        kong.Int(5),
				},
				Namespace: ingress.Namespace,
				Backend: kongstate.ServiceBackend{
					Name: defaultBackend.ServiceName,
					Port: PortDefFromIntStr(defaultBackend.ServicePort),
				},
			}
		}
		r := kongstate.Route{
			Ingress: util.FromK8sObject(&ingress),
			Route: kong.Route{
				Name:              kong.String(ingress.Namespace + "." + ingress.Name),
				Paths:             kong.StringSlice("/"),
				StripPath:         kong.Bool(false),
				PreserveHost:      kong.Bool(true),
				Protocols:         kong.StringSlice("http", "https"),
				RegexPriority:     kong.Int(0),
				RequestBuffering:  kong.Bool(true),
				ResponseBuffering: kong.Bool(true),
			},
		}
		service.Routes = append(service.Routes, r)
		result.ServiceNameToServices[serviceName] = service
	}

	return result
}

func fromIngressV1(log logrus.FieldLogger, ingressList []*networkingv1.Ingress) ingressRules {
	result := newIngressRules()

	var allDefaultBackends []networkingv1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec
		log = log.WithFields(logrus.Fields{
			"ingress_namespace": ingress.Namespace,
			"ingress_name":      ingress.Name,
		})

		if ingressSpec.DefaultBackend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1TLS(ingressSpec.TLS, ingress.Namespace)

		for i, rule := range ingressSpec.Rules {
			if rule.HTTP == nil {
				continue
			}
			for j, rulePath := range rule.HTTP.Paths {
				if strings.Contains(rulePath.Path, "//") {
					log.Errorf("rule skipped: invalid path: '%v'", rulePath.Path)
					continue
				}

				pathType := networkingv1.PathTypeImplementationSpecific
				if rulePath.PathType != nil {
					pathType = *rulePath.PathType
				}

				paths, err := pathsFromK8s(rulePath.Path, pathType)
				if err != nil {
					log.Errorf("rule skipped: pathsFromK8s: %v", err)
					continue
				}

				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
						// TODO (#834) Figure out a way to name the routes
						// This is not a stable scheme
						// 1. If a user adds a route in the middle,
						// due to a shift, all the following routes will
						// be PATCHED
						// 2. Is it guaranteed that the order is stable?
						// Meaning, the routes will always appear in the same
						// order?
						Name:              kong.String(fmt.Sprintf("%s.%s.%d%d", ingress.Namespace, ingress.Name, i, j)),
						Paths:             paths,
						StripPath:         kong.Bool(false),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(priorityForPath[pathType]),
						RequestBuffering:  kong.Bool(true),
						ResponseBuffering: kong.Bool(true),
					},
				}
				if rule.Host != "" {
					r.Hosts = kong.StringSlice(rule.Host)
				}

				port := PortDefFromServiceBackendPort(&rulePath.Backend.Service.Port)
				serviceName := fmt.Sprintf("%s.%s.%s", ingress.Namespace, rulePath.Backend.Service.Name,
					serviceBackendPortToStr(rulePath.Backend.Service.Port))
				service, ok := result.ServiceNameToServices[serviceName]
				if !ok {
					service = kongstate.Service{
						Service: kong.Service{
							Name: kong.String(serviceName),
							Host: kong.String(fmt.Sprintf("%s.%s.%s.svc", rulePath.Backend.Service.Name, ingress.Namespace,
								port.CanonicalString())),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend: kongstate.ServiceBackend{
							Name: rulePath.Backend.Service.Name,
							Port: port,
						},
					}
				}
				service.Routes = append(service.Routes, r)
				result.ServiceNameToServices[serviceName] = service
			}
		}
	}

	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	// Process the default backend
	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := allDefaultBackends[0].Spec.DefaultBackend
		port := PortDefFromServiceBackendPort(&defaultBackend.Service.Port)
		serviceName := fmt.Sprintf("%s.%s.%s", allDefaultBackends[0].Namespace, defaultBackend.Service.Name,
			port.CanonicalString())
		service, ok := result.ServiceNameToServices[serviceName]
		if !ok {
			service = kongstate.Service{
				Service: kong.Service{
					Name: kong.String(serviceName),
					Host: kong.String(fmt.Sprintf("%s.%s.%d.svc", defaultBackend.Service.Name, ingress.Namespace,
						defaultBackend.Service.Port.Number)),
					Port:           kong.Int(80),
					Protocol:       kong.String("http"),
					ConnectTimeout: kong.Int(60000),
					ReadTimeout:    kong.Int(60000),
					WriteTimeout:   kong.Int(60000),
					Retries:        kong.Int(5),
				},
				Namespace: ingress.Namespace,
				Backend: kongstate.ServiceBackend{
					Name: defaultBackend.Service.Name,
					Port: PortDefFromServiceBackendPort(&defaultBackend.Service.Port),
				},
			}
		}
		r := kongstate.Route{
			Ingress: util.FromK8sObject(&ingress),
			Route: kong.Route{
				Name:              kong.String(ingress.Namespace + "." + ingress.Name),
				Paths:             kong.StringSlice("/"),
				StripPath:         kong.Bool(false),
				PreserveHost:      kong.Bool(true),
				Protocols:         kong.StringSlice("http", "https"),
				RegexPriority:     kong.Int(0),
				RequestBuffering:  kong.Bool(true),
				ResponseBuffering: kong.Bool(true),
			},
		}
		service.Routes = append(service.Routes, r)
		result.ServiceNameToServices[serviceName] = service
	}

	return result
}

// fromHTTPRoutes processes all the HTTPRoute objects present in the cache and translates
// them to Kong Gateway configurations.
func fromHTTPRoutes(log logrus.FieldLogger, httpRouteList []*gatewayv1alpha2.HTTPRoute) ingressRules {
	result := newIngressRules()

	for _, httproute := range httpRouteList {
		// first we grab the spec and gather some metdata about the object
		objectInfo := util.FromK8sObject(httproute)
		spec := httproute.Spec

		// gather the hostnames that will be used (globally) for route matching
		hostnames := make([]*string, 0, len(spec.Hostnames))
		for _, hostname := range spec.Hostnames {
			hostnames = append(hostnames, kong.String(string(hostname)))
		}

		// each rule may represent a different set of backend services that will be accepting
		// traffic, so we make separate routes and Kong services for every present rule.
		for _, rule := range spec.Rules {
			// the HTTPRoute specification upstream specifically defines matches as
			// independent (e.g. each match is an OR with other matches, not an AND).
			// Therefore we treat each match rule as a separate Kong Route, so we iterate through
			// all matches to determine all the routes that will be needed for the services.
			var routes []kongstate.Route
			for matchNumber, match := range rule.Matches {
				// determine the name of the route, identify it as a route that belongs
				// to a Kubernetes HTTPRoute object.
				routeName := kong.String(fmt.Sprintf(
					"httproute.%s.%s.%d",
					httproute.Namespace,
					httproute.Name,
					matchNumber, // TODO: avoid route thrash from re-ordering?
				))

				// TODO: implement query param matches
				if len(match.QueryParams) > 0 {
					errmsg := "query param matches are not yet supported"
					log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %s", errmsg)
					continue
				}

				// TODO: implement regex path matches
				if *match.Path.Type == gatewayv1alpha2.PathMatchRegularExpression {
					errmsg := "regular expression path matches are not yet supported"
					log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %s", errmsg)
					continue
				}

				// build the route object using the method and pathing information
				r := kongstate.Route{
					Ingress: objectInfo,
					Route: kong.Route{
						Name:         routeName,
						Protocols:    kong.StringSlice("http", "https"),
						PreserveHost: kong.Bool(true),
						Hosts:        hostnames,
					},
				}
				log.Debugf("generated route %s for HTTPRoute %s/%s", routeName, httproute.Namespace, httproute.Name)

				// configure path matching information about the route if paths
				// matching was defined.
				if match.Path != nil {
					// determine the path match values
					r.Route.Paths = []*string{match.Path.Value}

					// determine whether path stripping needs to be enabled
					r.Route.StripPath = kong.Bool(match.Path.Type == nil || *match.Path.Type == gatewayv1alpha2.PathMatchPathPrefix)
				}

				// configure method matching information about the route if method
				// matching was defined.
				if match.Method != nil {
					r.Route.Methods = append(r.Route.Methods, kong.String(string(*match.Method)))
				}

				// convert header matching from HTTPRoute to Route format
				headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(match.Headers)
				if err != nil {
					log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %w", err)
					continue
				}
				r.Route.Headers = headers

				// add the route to the list of routes for the service(s)
				routes = append(routes, r)
			}

			// once all routes have been determined based on matching rules
			// we determine the Services they actually route to.
			for _, backendRef := range rule.BackendRefs {
				// determine the namespace for the service, or default to the same namespace
				// as the HTTPRoute object.
				//
				// TODO: need to add validation to restrict namespaces in backendRefs
				namespace := httproute.Namespace
				if backendRef.Namespace != nil {
					namespace = string(*backendRef.Namespace)
				}

				// determine the name of the Service
				serviceName := fmt.Sprintf("%s.%s.%d", namespace, backendRef.Name, *backendRef.Port)

				// determine the Service port
				port := kongstate.PortDef{
					Mode:   kongstate.PortModeByNumber,
					Number: int32(*backendRef.Port),
				}

				// check if the service is already known, and if not create it
				service, ok := result.ServiceNameToServices[serviceName]
				if !ok {
					service = kongstate.Service{
						Service: kong.Service{
							Name:           kong.String(serviceName),
							Host:           kong.String(fmt.Sprintf("%s.%s.%s.svc", backendRef.Name, namespace, port.CanonicalString())),
							Port:           kong.Int(int(*backendRef.Port)),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: httproute.Namespace,
						Backend: kongstate.ServiceBackend{
							Name: string(backendRef.Name),
							Port: port,
						},
					}
					log.Debugf("generated kong service %s for HTTPRoute %s/%s", serviceName, httproute.Namespace, httproute.Name)
				}

				// add all generated routes to this service
				service.Routes = append(service.Routes, routes...)

				// cache the service to avoid duplicates in further loop iterations
				result.ServiceNameToServices[serviceName] = service
			}
		}
	}

	return result
}

func fromTCPIngressV1beta1(log logrus.FieldLogger, tcpIngressList []*configurationv1beta1.TCPIngress) ingressRules {
	result := newIngressRules()

	sort.SliceStable(tcpIngressList, func(i, j int) bool {
		return tcpIngressList[i].CreationTimestamp.Before(
			&tcpIngressList[j].CreationTimestamp)
	})

	for _, ingress := range tcpIngressList {
		ingressSpec := ingress.Spec

		log = log.WithFields(logrus.Fields{
			"tcpingress_namespace": ingress.Namespace,
			"tcpingress_name":      ingress.Name,
		})

		result.SecretNameToSNIs.addFromIngressV1beta1TLS(tcpIngressToNetworkingTLS(ingressSpec.TLS), ingress.Namespace)

		for i, rule := range ingressSpec.Rules {
			if !util.IsValidPort(rule.Port) {
				log.Errorf("invalid TCPIngress: invalid port: %v", rule.Port)
				continue
			}
			r := kongstate.Route{
				Ingress: util.FromK8sObject(ingress),
				Route: kong.Route{
					// TODO (#834) Figure out a way to name the routes
					// This is not a stable scheme
					// 1. If a user adds a route in the middle,
					// due to a shift, all the following routes will
					// be PATCHED
					// 2. Is it guaranteed that the order is stable?
					// Meaning, the routes will always appear in the same
					// order?
					Name:      kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i)),
					Protocols: kong.StringSlice("tcp", "tls"),
					Destinations: []*kong.CIDRPort{
						{
							Port: kong.Int(rule.Port),
						},
					},
				},
			}
			host := rule.Host
			if host != "" {
				r.SNIs = kong.StringSlice(host)
			}
			if rule.Backend.ServiceName == "" {
				log.Errorf("invalid TCPIngress: empty serviceName")
				continue
			}
			if !util.IsValidPort(rule.Backend.ServicePort) {
				log.Errorf("invalid TCPIngress: invalid servicePort: %v", rule.Backend.ServicePort)
				continue
			}

			serviceName := fmt.Sprintf("%s.%s.%d", ingress.Namespace, rule.Backend.ServiceName, rule.Backend.ServicePort)
			service, ok := result.ServiceNameToServices[serviceName]
			if !ok {
				service = kongstate.Service{
					Service: kong.Service{
						Name: kong.String(serviceName),
						Host: kong.String(fmt.Sprintf("%s.%s.%d.svc", rule.Backend.ServiceName, ingress.Namespace,
							rule.Backend.ServicePort)),
						Port:           kong.Int(80),
						Protocol:       kong.String("tcp"),
						ConnectTimeout: kong.Int(60000),
						ReadTimeout:    kong.Int(60000),
						WriteTimeout:   kong.Int(60000),
						Retries:        kong.Int(5),
					},
					Namespace: ingress.Namespace,
					Backend: kongstate.ServiceBackend{
						Name: rule.Backend.ServiceName,
						Port: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
					},
				}
			}
			service.Routes = append(service.Routes, r)
			result.ServiceNameToServices[serviceName] = service
		}
	}

	return result
}

func fromUDPIngressV1beta1(log logrus.FieldLogger, ingressList []*configurationv1beta1.UDPIngress) ingressRules {
	result := newIngressRules()

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec

		log = log.WithFields(logrus.Fields{
			"udpingress_namespace": ingress.Namespace,
			"udpingress_name":      ingress.Name,
		})

		for i, rule := range ingressSpec.Rules {
			// validate the ports and servicenames for the rule
			if !util.IsValidPort(rule.Port) {
				log.Errorf("invalid UDPIngress: invalid port: %d", rule.Port)
				continue
			}
			if rule.Backend.ServiceName == "" {
				log.Errorf("invalid UDPIngress: empty serviceName")
				continue
			}
			if !util.IsValidPort(rule.Backend.ServicePort) {
				log.Errorf("invalid UDPIngress: invalid servicePort: %d", rule.Backend.ServicePort)
				continue
			}

			// generate the kong Route based on the listen port
			route := kongstate.Route{
				Ingress: util.FromK8sObject(ingress),
				Route: kong.Route{
					Name:         kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i) + ".udp"),
					Protocols:    kong.StringSlice("udp"),
					Destinations: []*kong.CIDRPort{{Port: kong.Int(rule.Port)}},
				},
			}

			// generate the kong Service backend for the UDPIngress rules
			host := fmt.Sprintf("%s.%s.%d.svc", rule.Backend.ServiceName, ingress.Namespace, rule.Backend.ServicePort)
			serviceName := fmt.Sprintf("%s.%s.%d.udp", ingress.Namespace, rule.Backend.ServiceName, rule.Backend.ServicePort)
			service, ok := result.ServiceNameToServices[serviceName]
			if !ok {
				service = kongstate.Service{
					Namespace: ingress.Namespace,
					Service: kong.Service{
						Name:     kong.String(serviceName),
						Protocol: kong.String("udp"),
						Host:     kong.String(host),
						Port:     kong.Int(rule.Backend.ServicePort),
					},
					Backend: kongstate.ServiceBackend{
						Name: rule.Backend.ServiceName,
						Port: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
					},
				}
			}
			service.Routes = append(service.Routes, route)
			result.ServiceNameToServices[serviceName] = service
		}
	}

	return result
}

func fromKnativeIngress(log logrus.FieldLogger, ingressList []*knative.Ingress) ingressRules {

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	services := map[string]kongstate.Service{}
	secretToSNIs := newSecretNameToSNIs()

	for _, ingress := range ingressList {
		log = log.WithFields(logrus.Fields{
			"knativeingress_namespace": ingress.Namespace,
			"knativeingress_name":      ingress.Name,
		})

		ingressSpec := ingress.Spec

		secretToSNIs.addFromIngressV1beta1TLS(knativeIngressToNetworkingTLS(ingress.Spec.TLS), ingress.Namespace)

		for i, rule := range ingressSpec.Rules {
			hosts := rule.Hosts
			if rule.HTTP == nil {
				continue
			}
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path

				if path == "" {
					path = "/"
				}
				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
						// TODO (#834) Figure out a way to name the routes
						// This is not a stable scheme
						// 1. If a user adds a route in the middle,
						// due to a shift, all the following routes will
						// be PATCHED
						// 2. Is it guaranteed that the order is stable?
						// Meaning, the routes will always appear in the same
						// order?
						Name:              kong.String(fmt.Sprintf("%s.%s.%d%d", ingress.Namespace, ingress.Name, i, j)),
						Paths:             kong.StringSlice(path),
						StripPath:         kong.Bool(false),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						RequestBuffering:  kong.Bool(true),
						ResponseBuffering: kong.Bool(true),
					},
				}
				r.Hosts = kong.StringSlice(hosts...)

				knativeBackend := knativeSelectSplit(rule.Splits)
				serviceName := fmt.Sprintf("%s.%s.%s", knativeBackend.ServiceNamespace, knativeBackend.ServiceName,
					knativeBackend.ServicePort.String())
				serviceHost := fmt.Sprintf("%s.%s.%s.svc", knativeBackend.ServiceName, knativeBackend.ServiceNamespace,
					knativeBackend.ServicePort.String())
				service, ok := services[serviceName]
				if !ok {

					var headers []string
					for key, value := range knativeBackend.AppendHeaders {
						headers = append(headers, key+":"+value)
					}
					for key, value := range rule.AppendHeaders {
						headers = append(headers, key+":"+value)
					}

					service = kongstate.Service{
						Service: kong.Service{
							Name:           kong.String(serviceName),
							Host:           kong.String(serviceHost),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend: kongstate.ServiceBackend{
							Name: knativeBackend.ServiceName,
							Port: PortDefFromIntStr(knativeBackend.ServicePort),
						},
					}
					if len(headers) > 0 {
						service.Plugins = append(service.Plugins, kong.Plugin{
							Name: kong.String("request-transformer"),
							Config: kong.Configuration{
								"add": map[string]interface{}{
									"headers": headers,
								},
							},
						})
					}
				}
				service.Routes = append(service.Routes, r)
				services[serviceName] = service
			}
		}
	}

	return ingressRules{
		ServiceNameToServices: services,
		SecretNameToSNIs:      secretToSNIs,
	}
}

func knativeSelectSplit(splits []knative.IngressBackendSplit) knative.IngressBackendSplit {
	if len(splits) == 0 {
		return knative.IngressBackendSplit{}
	}
	res := splits[0]
	maxPercentage := splits[0].Percent
	if len(splits) == 1 {
		return res
	}
	for i := 1; i < len(splits); i++ {
		if splits[i].Percent > maxPercentage {
			res = splits[i]
			maxPercentage = res.Percent
		}
	}
	return res
}

func pathsFromK8s(path string, pathType networkingv1.PathType) ([]*string, error) {
	switch pathType {
	case networkingv1.PathTypePrefix:
		base := strings.Trim(path, "/")
		if base == "" {
			return kong.StringSlice("/"), nil
		}
		return kong.StringSlice(
			"/"+base+"$",
			"/"+base+"/",
		), nil
	case networkingv1.PathTypeExact:
		relative := strings.TrimLeft(path, "/")
		return kong.StringSlice("/" + relative + "$"), nil
	case networkingv1.PathTypeImplementationSpecific:
		if path == "" {
			return kong.StringSlice("/"), nil
		}
		return kong.StringSlice(path), nil
	}

	return nil, fmt.Errorf("unknown pathType %v", pathType)
}

var priorityForPath = map[networkingv1.PathType]int{
	networkingv1.PathTypeExact:                  300,
	networkingv1.PathTypePrefix:                 200,
	networkingv1.PathTypeImplementationSpecific: 100,
}

func PortDefFromServiceBackendPort(sbp *networkingv1.ServiceBackendPort) kongstate.PortDef {
	switch {
	case sbp.Name != "":
		return kongstate.PortDef{Mode: kongstate.PortModeByName, Name: sbp.Name}
	case sbp.Number != 0:
		return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: sbp.Number}
	default:
		return kongstate.PortDef{Mode: kongstate.PortModeImplicit}
	}
}

func PortDefFromIntStr(is intstr.IntOrString) kongstate.PortDef {
	if is.Type == intstr.String {
		return kongstate.PortDef{Mode: kongstate.PortModeByName, Name: is.StrVal}
	}
	return kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: is.IntVal}
}

// -----------------------------------------------------------------------------
// Utilities - Gateway APIs
// -----------------------------------------------------------------------------

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayv1alpha2.HTTPHeaderMatch) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		// TODO: implement regex header matching
		if header.Type != nil && *header.Type == gatewayv1alpha2.HeaderMatchRegularExpression {
			return nil, fmt.Errorf("regular expression header matches are not yet supported")
		}

		// split the header values by comma
		values := strings.Split(header.Value, ",")
		convertedHeaders[string(header.Name)] = values
	}

	return convertedHeaders, nil
}
