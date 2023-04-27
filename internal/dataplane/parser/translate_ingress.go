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
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
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
						Tags:              util.GenerateTagsForObject(ingress),
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
				result.ServiceNameToParent[serviceName] = ingress
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
					Tags:           util.GenerateTagsForObject(result.ServiceNameToParent[serviceName]),
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
				Tags:              util.GenerateTagsForObject(result.ServiceNameToParent[serviceName]),
			},
		}
		service.Routes = append(service.Routes, r)
		result.ServiceNameToServices[serviceName] = service
		result.ServiceNameToParent[serviceName] = &ingress
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

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	servicesCache := make(kongServicesCache)
	var allDefaultBackends []netv1.Ingress
	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec

		if ingressSpec.DefaultBackend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}

		result.SecretNameToSNIs.addFromIngressV1TLS(ingressSpec.TLS, ingress)

		if wasServicesCacheUpdated := p.ingressV1ToKongService(ingress, icp, servicesCache); wasServicesCacheUpdated {
			p.ReportKubernetesObjectUpdate(ingress)
		}
	}

	for _, service := range servicesCache {
		result.ServiceNameToServices[*service.Name] = service
		result.ServiceNameToParent[*service.Name] = service.Parent
	}

	// Add a default backend if it exists.
	defaultBackendService, ok := getDefaultBackendService(allDefaultBackends)
	if ok {
		result.ServiceNameToServices[*defaultBackendService.Name] = defaultBackendService
		result.ServiceNameToParent[*defaultBackendService.Name] = defaultBackendService.Parent
	}

	return result
}

// ingressV1ToKongServicesCache is a cache of Kong Services indexed by their name.
type kongServicesCache map[string]kongstate.Service

// ingressV1ToKongService translates IngressV1 object into Kong Service. It inserts the Kong Service into the passed servicesCache.
// Returns true if the passed servicesCache was updated.
func (p *Parser) ingressV1ToKongService(
	ingress *netv1.Ingress,
	icp v1alpha1.IngressClassParametersSpec,
	servicesCache kongServicesCache,
) bool {
	if p.featureEnabledCombinedServiceRoutes {
		return p.ingressV1ToKongServiceCombinedRoutes(ingress, icp, servicesCache)
	}

	return p.ingressV1ToKongServiceLegacy(ingress, icp, servicesCache)
}

// ingressV1ToKongServiceLegacy translates IngressV1 object into Kong Service. It inserts the Kong Service into the passed servicesCache.
// Returns true if the passed servicesCache was updated. It is used when CombinedRoutes feature flag is enabled.
func (p *Parser) ingressV1ToKongServiceCombinedRoutes(
	ingress *netv1.Ingress,
	icp v1alpha1.IngressClassParametersSpec,
	servicesCache kongServicesCache,
) bool {
	wasServicesCacheUpdated := false

	regexPrefix := translators.ControllerPathRegexPrefix
	if prefix, ok := ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.RegexPrefixKey]; ok {
		regexPrefix = prefix
	}
	for _, kongStateService := range translators.TranslateIngress(ingress, p.flagEnabledRegexPathPrefix) {
		for _, route := range kongStateService.Routes {
			for i, path := range route.Paths {
				newPath := translators.MaybePrependRegexPrefix(*path, regexPrefix, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
				route.Paths[i] = &newPath
			}
		}

		servicesCache[*kongStateService.Service.Name] = *kongStateService
		wasServicesCacheUpdated = true
	}

	return wasServicesCacheUpdated
}

// ingressV1ToKongServiceLegacy translates IngressV1 object into Kong Service. It inserts the Kong Service into the passed servicesCache.
// Returns true if the passed servicesCache was updated. It is used when the CombinedRoutes feature flag is disabled.
func (p *Parser) ingressV1ToKongServiceLegacy(
	ingress *netv1.Ingress,
	icp v1alpha1.IngressClassParametersSpec,
	servicesCache kongServicesCache,
) bool {
	wasServicesCacheUpdated := false

	ingressSpec := ingress.Spec
	maybePrependRegexPrefixFn := translators.MaybePrependRegexPrefixForIngressV1Fn(ingress, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
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
				paths[i] = maybePrependRegexPrefixFn(*path)
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
					Tags:              util.GenerateTagsForObject(ingress),
				},
			}
			if rule.Host != "" {
				r.Hosts = kong.StringSlice(rule.Host)
			}

			port := translators.PortDefFromServiceBackendPort(&rulePath.Backend.Service.Port)
			serviceName := fmt.Sprintf(
				"%s.%s.%s",
				ingress.Namespace,
				rulePath.Backend.Service.Name,
				serviceBackendPortToStr(rulePath.Backend.Service.Port),
			)
			service, ok := servicesCache[serviceName]
			if !ok {
				service = kongstate.Service{
					Service: kong.Service{
						Name: kong.String(serviceName),
						Host: kong.String(fmt.Sprintf(
							"%s.%s.%s.svc",
							rulePath.Backend.Service.Name,
							ingress.Namespace,
							port.CanonicalString(),
						)),
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
			servicesCache[serviceName] = service
			wasServicesCacheUpdated = true
		}
	}

	return wasServicesCacheUpdated
}

func getDefaultBackendService(allDefaultBackends []netv1.Ingress) (kongstate.Service, bool) {
	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := allDefaultBackends[0].Spec.DefaultBackend
		port := translators.PortDefFromServiceBackendPort(&defaultBackend.Service.Port)
		serviceName := fmt.Sprintf(
			"%s.%s.%s",
			allDefaultBackends[0].Namespace,
			defaultBackend.Service.Name,
			port.CanonicalString(),
		)
		service := kongstate.Service{
			Service: kong.Service{
				Name: kong.String(serviceName),
				Host: kong.String(fmt.Sprintf(
					"%s.%s.%s.svc",
					defaultBackend.Service.Name,
					ingress.Namespace,
					port.CanonicalString(),
				)),
				Port:           kong.Int(DefaultHTTPPort),
				Protocol:       kong.String("http"),
				ConnectTimeout: kong.Int(DefaultServiceTimeout),
				ReadTimeout:    kong.Int(DefaultServiceTimeout),
				WriteTimeout:   kong.Int(DefaultServiceTimeout),
				Retries:        kong.Int(DefaultRetries),
				Tags:           util.GenerateTagsForObject(&ingress),
			},
			Namespace: ingress.Namespace,
			Backends: []kongstate.ServiceBackend{{
				Name:    defaultBackend.Service.Name,
				PortDef: port,
			}},
			Parent: &ingress,
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
				Tags:              util.GenerateTagsForObject(&ingress),
			},
		}
		service.Routes = append(service.Routes, r)
		return service, true
	}

	return kongstate.Service{}, false
}
