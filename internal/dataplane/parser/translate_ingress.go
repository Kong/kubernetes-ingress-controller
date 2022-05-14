package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func (p *Parser) ingressRulesFromIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList := p.storer.ListIngressesV1beta1()

	var allDefaultBackends []networkingv1beta1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec
		log := p.logger.WithFields(logrus.Fields{
			"ingress_namespace": ingress.Namespace,
			"ingress_name":      ingress.Name,
		})

		if ingressSpec.Backend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1beta1TLS(ingressSpec.TLS, ingress.Namespace)

		var objectSuccessfullyParsed bool

		// from the ingress rules we're going to create a map of the Kubernetes
		// Service names to the hosts which should be routed, to the ports for the
		// service which should be routed for that host, and then the individual
		// paths for that combination of service, host and port.
		// This allows us to build a single route for every specific service, host
		// and port, rather than a route for each, saving on some space in the size
		// of the dataplane configuration.
		routableServices := make(map[string]map[string]map[int32]map[string]networkingv1beta1.PathType)
		for _, ingressRule := range ingressSpec.Rules {
			// first we check if there are actually any rules present, if there aren't
			// this is a misconfigured rule and we'll simply skip over it.
			if ingressRule.HTTP == nil || len(ingressRule.HTTP.Paths) < 1 {
				continue
			}

			host := ingressRule.Host

			// now that we know this ingress rule actually has paths, we'll map out all
			// of the hosts and paths that belong to the service and port.
			for _, pathRule := range ingressRule.HTTP.Paths {
				// check that the path is actually valid, don't try to route this rule
				// if not.
				path := pathRule.Path
				if strings.Contains(path, "//") {
					log.Errorf("rule skipped: invalid path: '%v'", path)
					continue
				}

				if path == "" {
					path = "/"
				}

				serviceName := pathRule.Backend.ServiceName
				if routableServices[serviceName] == nil {
					routableServices[serviceName] = make(map[string]map[int32]map[string]networkingv1beta1.PathType)
				}

				if routableServices[serviceName][host] == nil {
					routableServices[serviceName][host] = make(map[int32]map[string]networkingv1beta1.PathType)
				}

				port := pathRule.Backend.ServicePort.IntVal
				if routableServices[serviceName][host][port] == nil {
					routableServices[serviceName][host][port] = make(map[string]networkingv1beta1.PathType)
				}

				pathType := networkingv1beta1.PathTypePrefix
				if pathRule.PathType != nil {
					pathType = *pathRule.PathType
				}
				routableServices[serviceName][host][port][path] = pathType
			}
		}

		for serviceName, hostnameRules := range routableServices {
			for host, portRules := range hostnameRules {
				for portNumber, pathsMap := range portRules {
					portDef := PortDefFromIntStr(intstr.FromInt(int(portNumber)))
					kongServiceName := ingress.Namespace + "." + serviceName + "." + portDef.CanonicalString()
					kongStateService, ok := result.ServiceNameToServices[serviceName]
					if !ok {
						kongStateService = kongstate.Service{
							Service: kong.Service{
								Name:           kong.String(kongServiceName),
								Host:           kong.String(fmt.Sprintf("%s.%s.%d.svc", serviceName, ingress.Namespace, portNumber)),
								Port:           kong.Int(DefaultHTTPPort),
								Protocol:       kong.String("http"),
								Path:           kong.String("/"),
								ConnectTimeout: kong.Int(DefaultServiceTimeout),
								ReadTimeout:    kong.Int(DefaultServiceTimeout),
								WriteTimeout:   kong.Int(DefaultServiceTimeout),
								Retries:        kong.Int(DefaultRetries),
							},
							Namespace: ingress.Namespace,
							Backends: []kongstate.ServiceBackend{{
								Name:    serviceName,
								PortDef: portDef,
							}},
						}
					}

					var routeName string
					if host == "" {
						routeName = fmt.Sprintf("%s.%s.%s.%d", ingress.Namespace, ingress.Name, serviceName, portNumber)
					} else {
						routeName = fmt.Sprintf("%s.%s.%s.%s.%d", ingress.Namespace, ingress.Name, serviceName, host, portNumber)
					}

					kongStateRoute := kongstate.Route{
						Ingress: util.FromK8sObject(ingress),
						Route: kong.Route{
							Name:              kong.String(routeName),
							StripPath:         kong.Bool(false),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							RequestBuffering:  kong.Bool(true),
							ResponseBuffering: kong.Bool(true),
						},
					}

					if host != "" {
						kongStateRoute.Route.Hosts = append(kongStateRoute.Route.Hosts, kong.String(host))
					}

					for path, pathType := range pathsMap {
						paths, err := pathsFromK8sLegacy(path, pathType)
						if err != nil {
							log.Errorf("skipping route path %s for ingress %s/%s: %s", path, ingress.Namespace, ingress.Name, err.Error())
							continue
						}
						kongStateRoute.Paths = append(kongStateRoute.Paths, paths...)
					}

					kongStateService.Routes = append(kongStateService.Routes, kongStateRoute)
					result.ServiceNameToServices[serviceName] = kongStateService
					objectSuccessfullyParsed = true
				}
			}
		}

		if objectSuccessfullyParsed {
			p.ReportKubernetesObjectUpdate(ingress)
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
					Port:           kong.Int(DefaultHTTPPort),
					Protocol:       kong.String("http"),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: ingress.Namespace,
				Backends: []kongstate.ServiceBackend{{
					Name:    defaultBackend.ServiceName,
					PortDef: PortDefFromIntStr(defaultBackend.ServicePort),
				}},
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

func (p *Parser) ingressRulesFromIngressV1() ingressRules {
	result := newIngressRules()

	ingressList := p.storer.ListIngressesV1()

	var allDefaultBackends []networkingv1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec
		log := p.logger.WithFields(logrus.Fields{
			"ingress_namespace": ingress.Namespace,
			"ingress_name":      ingress.Name,
		})

		if ingressSpec.DefaultBackend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1TLS(ingressSpec.TLS, ingress.Namespace)

		var objectSuccessfullyParsed bool

		// from the ingress rules we're going to create a map of the Kubernetes
		// Service names to the hosts which should be routed, to the ports for the
		// service which should be routed for that host, and then the individual
		// paths for that combination of service, host and port.
		// This allows us to build a single route for every specific service, host
		// and port, rather than a route for each, saving on some space in the size
		// of the dataplane configuration.
		routableServices := make(map[string]map[string]map[int32]map[string]networkingv1.PathType)
		for _, ingressRule := range ingressSpec.Rules {
			// first we check if there are actually any rules present, if there aren't
			// this is a misconfigured rule and we'll simply skip over it.
			if ingressRule.HTTP == nil || len(ingressRule.HTTP.Paths) < 1 {
				continue
			}

			host := ingressRule.Host

			// now that we know this ingress rule actually has paths, we'll map out all
			// of the hosts and paths that belong to the service and port.
			for _, pathRule := range ingressRule.HTTP.Paths {
				// check that the path is actually valid, don't try to route this rule
				// if not.
				path := pathRule.Path
				if strings.Contains(path, "//") {
					log.Errorf("rule skipped: invalid path: '%v'", path)
					continue
				}

				if path == "" {
					path = "/"
				}

				serviceName := pathRule.Backend.Service.Name
				if routableServices[serviceName] == nil {
					routableServices[serviceName] = make(map[string]map[int32]map[string]networkingv1.PathType)
				}

				if routableServices[serviceName][host] == nil {
					routableServices[serviceName][host] = make(map[int32]map[string]networkingv1.PathType)
				}

				port := pathRule.Backend.Service.Port.Number
				if routableServices[serviceName][host][port] == nil {
					routableServices[serviceName][host][port] = make(map[string]networkingv1.PathType)
				}

				pathType := networkingv1.PathTypePrefix
				if pathRule.PathType != nil {
					pathType = *pathRule.PathType
				}
				routableServices[serviceName][host][port][path] = pathType
			}
		}

		for serviceName, hostnameRules := range routableServices {
			for host, portRules := range hostnameRules {
				for portNumber, pathsMap := range portRules {
					portDef := PortDefFromIntStr(intstr.FromInt(int(portNumber)))
					kongServiceName := ingress.Namespace + "." + serviceName + "." + portDef.CanonicalString()
					kongStateService, ok := result.ServiceNameToServices[serviceName]
					if !ok {
						kongStateService = kongstate.Service{
							Service: kong.Service{
								Name:           kong.String(kongServiceName),
								Host:           kong.String(fmt.Sprintf("%s.%s.%d.svc", serviceName, ingress.Namespace, portNumber)),
								Port:           kong.Int(DefaultHTTPPort),
								Protocol:       kong.String("http"),
								Path:           kong.String("/"),
								ConnectTimeout: kong.Int(DefaultServiceTimeout),
								ReadTimeout:    kong.Int(DefaultServiceTimeout),
								WriteTimeout:   kong.Int(DefaultServiceTimeout),
								Retries:        kong.Int(DefaultRetries),
							},
							Namespace: ingress.Namespace,
							Backends: []kongstate.ServiceBackend{{
								Name:    serviceName,
								PortDef: portDef,
							}},
						}
					}

					var routeName string
					if host == "" {
						routeName = fmt.Sprintf("%s.%s.%s.%d", ingress.Namespace, ingress.Name, serviceName, portNumber)
					} else {
						routeName = fmt.Sprintf("%s.%s.%s.%s.%d", ingress.Namespace, ingress.Name, serviceName, host, portNumber)
					}

					kongStateRoute := kongstate.Route{
						Ingress: util.FromK8sObject(ingress),
						Route: kong.Route{
							Name:              kong.String(routeName),
							StripPath:         kong.Bool(false),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							RequestBuffering:  kong.Bool(true),
							ResponseBuffering: kong.Bool(true),
						},
					}

					if host != "" {
						kongStateRoute.Route.Hosts = append(kongStateRoute.Route.Hosts, kong.String(host))
					}

					for path, pathType := range pathsMap {
						paths, err := pathsFromK8s(path, pathType)
						if err != nil {
							log.Errorf("skipping route path %s for ingress %s/%s: %s", path, ingress.Namespace, ingress.Name, err.Error())
							continue
						}
						kongStateRoute.Paths = append(kongStateRoute.Paths, paths...)
					}

					kongStateService.Routes = append(kongStateService.Routes, kongStateRoute)
					result.ServiceNameToServices[serviceName] = kongStateService
					objectSuccessfullyParsed = true
				}
			}
		}

		if objectSuccessfullyParsed {
			p.ReportKubernetesObjectUpdate(ingress)
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
					Port:           kong.Int(DefaultHTTPPort),
					Protocol:       kong.String("http"),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: ingress.Namespace,
				Backends: []kongstate.ServiceBackend{{
					Name:    defaultBackend.Service.Name,
					PortDef: PortDefFromServiceBackendPort(&defaultBackend.Service.Port),
				}},
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
