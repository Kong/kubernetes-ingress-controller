package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
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

	// Collect all default backends and TLS SNIs.
	var allDefaultBackends []netv1.Ingress
	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec
		if ingressSpec.DefaultBackend != nil {
			allDefaultBackends = append(allDefaultBackends, *ingress)
		}
		result.SecretNameToSNIs.addFromIngressV1TLS(ingressSpec.TLS, ingress)
	}

	// Translate Ingress objects into Kong Services.
	for _, service := range p.ingressesV1ToKongServices(ingressList, icp) {
		result.ServiceNameToServices[*service.Name] = service
		result.ServiceNameToParent[*service.Name] = service.Parent
	}

	// Add a default backend if it exists.
	defaultBackendService, ok := getDefaultBackendService(allDefaultBackends, p.featureFlags.ExpressionRoutes)
	if ok {
		result.ServiceNameToServices[*defaultBackendService.Name] = defaultBackendService
		result.ServiceNameToParent[*defaultBackendService.Name] = defaultBackendService.Parent
	}

	return result
}

// ingressV1ToKongServicesCache is a cache of Kong Services indexed by their name.
type kongServicesCache map[string]kongstate.Service

// ingressesV1ToKongServices translates IngressV1 object into Kong Service. It inserts the Kong Service into the passed servicesCache.
// Returns true if the passed servicesCache was updated.
func (p *Parser) ingressesV1ToKongServices(
	ingresses []*netv1.Ingress,
	icp v1alpha1.IngressClassParametersSpec,
) kongServicesCache {
	if p.featureFlags.CombinedServiceRoutes {
		return p.ingressV1ToKongServiceCombinedRoutes(ingresses, icp)
	}
	return p.ingressV1ToKongServiceLegacy(ingresses, icp)
}

// ingressV1ToKongServiceLegacy translates a slice of IngressV1 object into Kong Services.
func (p *Parser) ingressV1ToKongServiceCombinedRoutes(
	ingresses []*netv1.Ingress,
	icp v1alpha1.IngressClassParametersSpec,
) kongServicesCache {
	return translators.TranslateIngresses(ingresses, icp, translators.TranslateIngressFeatureFlags{
		RegexPathPrefix:  p.featureFlags.RegexPathPrefix,
		ExpressionRoutes: p.featureFlags.ExpressionRoutes,
		CombinedServices: p.featureFlags.CombinedServices,
	}, p.parsedObjectsCollector)
}

// ingressV1ToKongServiceLegacy translates a slice IngressV1 object into Kong Services.
func (p *Parser) ingressV1ToKongServiceLegacy(ingresses []*netv1.Ingress, icp v1alpha1.IngressClassParametersSpec) kongServicesCache {
	servicesCache := make(kongServicesCache)

	for _, ingress := range ingresses {
		ingressSpec := ingress.Spec
		maybePrependRegexPrefixFn := translators.MaybePrependRegexPrefixForIngressV1Fn(ingress, icp.EnableLegacyRegexDetection && p.featureFlags.RegexPathPrefix)
		for i, rule := range ingressSpec.Rules {
			if rule.HTTP == nil {
				continue
			}
			for j, rulePath := range rule.HTTP.Paths {
				pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
				if rulePath.PathType == nil {
					rulePath.PathType = &pathTypeImplementationSpecific
				}

				paths := translators.PathsFromIngressPaths(rulePath, p.featureFlags.RegexPathPrefix)
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
				p.registerSuccessfullyParsedObject(ingress)
			}
		}
	}

	return servicesCache
}

// getDefaultBackendService picks the oldest Ingress with a DefaultBackend defined and returns a Kong Service for it.
func getDefaultBackendService(allDefaultBackends []netv1.Ingress, expressionRoutes bool) (kongstate.Service, bool) {
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
		r := translateIngressDefaultBackendRoute(&ingress, util.GenerateTagsForObject(&ingress), expressionRoutes)
		service.Routes = append(service.Routes, *r)
		return service, true
	}

	return kongstate.Service{}, false
}

func translateIngressDefaultBackendRoute(ingress *netv1.Ingress, tags []*string, expressionRoutes bool) *kongstate.Route {
	r := &kongstate.Route{
		Ingress: util.FromK8sObject(ingress),
		Route: kong.Route{
			Name:              kong.String(ingress.Namespace + "." + ingress.Name),
			StripPath:         kong.Bool(false),
			PreserveHost:      kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ResponseBuffering: kong.Bool(true),
			Tags:              tags,
		},
		ExpressionRoutes: expressionRoutes,
	}

	if expressionRoutes {
		catchAllMatcher := atc.And(
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/"),
			atc.Or(atc.NewPredicateNetProtocol(atc.OpEqual, "http"), atc.NewPredicateNetProtocol(atc.OpEqual, "https")),
		)
		atc.ApplyExpression(&r.Route, catchAllMatcher, translators.IngressDefaultBackendPriority)
	} else {
		r.Route.Paths = kong.StringSlice("/")
		r.Route.Protocols = kong.StringSlice("http", "https")
		r.Route.RegexPriority = kong.Int(0)
	}
	return r
}
