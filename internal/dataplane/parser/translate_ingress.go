package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func serviceBackendPortToStr(port netv1.ServiceBackendPort) string {
	if port.Name != "" {
		return fmt.Sprintf("pname-%s", port.Name)
	}
	return fmt.Sprintf("pnum-%d", port.Number)
}

var priorityForPath = map[netv1.PathType]int{
	netv1.PathTypeExact:                  300,
	netv1.PathTypePrefix:                 200,
	netv1.PathTypeImplementationSpecific: 100,
}

func (p *Parser) ingressRulesFromIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList := p.storer.ListIngressesV1beta1()
	icp, err := getIngressClassParametersOrDefault(p.storer)
	if err != nil {
		if !errors.As(err, &store.ErrNotFound{}) {
			// anything else is unexpected
			p.logger.Errorf("could not find IngressClassParameters, using defaults: %s", err)
		}
	}

	var allDefaultBackends []netv1beta1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		regexPrefix := translators.ControllerPathRegexPrefix
		if prefix, ok := ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.RegexPrefixKey]; ok {
			regexPrefix = prefix
		}
		ingressSpec := ingress.Spec

		if ingressSpec.Backend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1TLS(v1beta1toV1TLS(ingressSpec.TLS), ingress)

		var objectSuccessfullyParsed bool
		for i, rule := range ingressSpec.Rules {
			host := rule.Host
			if rule.HTTP == nil {
				continue
			}
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path
				path = maybePrependRegexPrefix(path, regexPrefix, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
				if path == "" {
					path = "/"
				}
				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
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
							Name:    rule.Backend.ServiceName,
							PortDef: translators.PortDefFromIntStr(rule.Backend.ServicePort),
						}},
						Parent: ingress,
					}
				}
				service.Routes = append(service.Routes, r)
				result.ServiceNameToServices[serviceName] = service
				objectSuccessfullyParsed = true
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
					PortDef: translators.PortDefFromIntStr(defaultBackend.ServicePort),
				}},
				Parent: &ingress,
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
	icp, err := getIngressClassParametersOrDefault(p.storer)
	if err != nil {
		if !errors.As(err, &store.ErrNotFound{}) {
			// anything else is unexpected
			p.logger.Errorf("could not find IngressClassParameters, using defaults: %s", err)
		}
	}

	var allDefaultBackends []netv1.Ingress
	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		regexPrefix := translators.ControllerPathRegexPrefix
		if prefix, ok := ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.RegexPrefixKey]; ok {
			regexPrefix = prefix
		}
		ingressSpec := ingress.Spec

		if ingressSpec.DefaultBackend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1TLS(ingressSpec.TLS, ingress)

		var objectSuccessfullyParsed bool

		if p.featureEnabledCombinedServiceRoutes {
			for _, kongStateService := range translators.TranslateIngress(ingress, p.flagEnabledRegexPathPrefix) {
				for _, route := range kongStateService.Routes {
					for i, path := range route.Paths {
						newPath := maybePrependRegexPrefix(*path, regexPrefix, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
						route.Paths[i] = &newPath
					}
				}
				result.ServiceNameToServices[*kongStateService.Service.Name] = *kongStateService
				objectSuccessfullyParsed = true
			}
		} else {
			for i, rule := range ingressSpec.Rules {
				if rule.HTTP == nil {
					continue
				}
				for j, rulePath := range rule.HTTP.Paths {
					pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
					if rulePath.PathType == nil {
						rulePath.PathType = &pathTypeImplementationSpecific
					}

					paths := translators.PathsFromIngressPaths(rulePath, p.flagEnabledRegexPathPrefix)
					if paths == nil {
						// registering a failure, but technically it should never happen thanks to Kubernetes API validations
						p.registerTranslationFailure(
							fmt.Sprintf("could not translate Ingress Path %s to Kong paths", rulePath.Path), ingress,
						)
						continue
					}

					for i, path := range paths {
						newPath := maybePrependRegexPrefix(*path, regexPrefix, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
						paths[i] = &newPath
					}

					r := kongstate.Route{
						Ingress: util.FromK8sObject(ingress),
						Route: kong.Route{
							Name:              kong.String(fmt.Sprintf("%s.%s.%d%d", ingress.Namespace, ingress.Name, i, j)),
							Paths:             paths,
							StripPath:         kong.Bool(false),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(priorityForPath[*rulePath.PathType]),
							RequestBuffering:  kong.Bool(true),
							ResponseBuffering: kong.Bool(true),
						},
					}
					if rule.Host != "" {
						r.Hosts = kong.StringSlice(rule.Host)
					}

					port := translators.PortDefFromServiceBackendPort(&rulePath.Backend.Service.Port)
					serviceName := fmt.Sprintf("%s.%s.%s", ingress.Namespace, rulePath.Backend.Service.Name,
						serviceBackendPortToStr(rulePath.Backend.Service.Port))
					service, ok := result.ServiceNameToServices[serviceName]
					if !ok {
						service = kongstate.Service{
							Service: kong.Service{
								Name: kong.String(serviceName),
								Host: kong.String(fmt.Sprintf("%s.%s.%s.svc", rulePath.Backend.Service.Name, ingress.Namespace,
									port.CanonicalString())),
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
								Name:    rulePath.Backend.Service.Name,
								PortDef: port,
							}},
							Parent: ingress,
						}
					}
					service.Routes = append(service.Routes, r)
					result.ServiceNameToServices[serviceName] = service
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
		port := translators.PortDefFromServiceBackendPort(&defaultBackend.Service.Port)
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
					PortDef: translators.PortDefFromServiceBackendPort(&defaultBackend.Service.Port),
				}},
				Parent: &ingress,
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
