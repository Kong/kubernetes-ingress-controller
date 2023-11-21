package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
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
	servicesCache := IngressesV1ToKongServices(
		p.featureFlags,
		ingressList,
		icp,
		p.parsedObjectsCollector,
		p.failuresCollector,
	)
	for i := range servicesCache {
		service := servicesCache[i]
		if err := translators.MaybeRewriteURI(&service, p.featureFlags.RewriteURIs); err != nil {
			p.registerTranslationFailure(err.Error(), service.Parent)
			continue
		}

		result.ServiceNameToServices[*service.Name] = service
		result.ServiceNameToParent[*service.Name] = service.Parent
	}

	// Add a default backend if it exists.
	defaultBackendService, ok := getDefaultBackendService(allDefaultBackends, p.featureFlags.ExpressionRoutes)
	if ok {
		// When such service would overwrite an existing service, merge the routes.
		if svc, ok := result.ServiceNameToServices[*defaultBackendService.Name]; ok {
			svc.Routes = append(svc.Routes, defaultBackendService.Routes...)
			defaultBackendService = svc
		}
		result.ServiceNameToServices[*defaultBackendService.Name] = defaultBackendService
		result.ServiceNameToParent[*defaultBackendService.Name] = defaultBackendService.Parent
	}

	return result
}

// KongServicesCache is a cache of Kong Services indexed by their name.
type KongServicesCache map[string]kongstate.Service

// IngressesV1ToKongServices translates IngressV1 object into Kong Service, returns them indexed by name.
// Argument parsedObjectsCollector is used to register all successfully parsed objects. In case of a failure,
// the object is registered in failuresCollector.
func IngressesV1ToKongServices(
	featureFlags FeatureFlags,
	ingresses []*netv1.Ingress,
	icp kongv1alpha1.IngressClassParametersSpec,
	parsedObjectsCollector *ObjectsCollector,
	failuresCollector *failures.ResourceFailuresCollector,
) KongServicesCache {
	if featureFlags.CombinedServiceRoutes {
		return ingressV1ToKongServiceCombinedRoutes(featureFlags, ingresses, icp, parsedObjectsCollector)
	}
	return ingressV1ToKongServiceLegacy(featureFlags, ingresses, icp, parsedObjectsCollector, failuresCollector)
}

// ingressV1ToKongServiceLegacy translates a slice of IngressV1 object into Kong Services.
func ingressV1ToKongServiceCombinedRoutes(
	featureFlags FeatureFlags,
	ingresses []*netv1.Ingress,
	icp kongv1alpha1.IngressClassParametersSpec,
	parsedObjectsCollector *ObjectsCollector,
) KongServicesCache {
	return translators.TranslateIngresses(ingresses, icp, translators.TranslateIngressFeatureFlags{
		RegexPathPrefix:  featureFlags.RegexPathPrefix,
		ExpressionRoutes: featureFlags.ExpressionRoutes,
		CombinedServices: featureFlags.CombinedServices,
	}, parsedObjectsCollector)
}

// ingressV1ToKongServiceLegacy translates a slice IngressV1 object into Kong Services.
func ingressV1ToKongServiceLegacy(
	featureFlags FeatureFlags,
	ingresses []*netv1.Ingress,
	icp kongv1alpha1.IngressClassParametersSpec,
	parsedObjectsCollector *ObjectsCollector,
	failuresCollector *failures.ResourceFailuresCollector,
) KongServicesCache {
	servicesCache := make(KongServicesCache)

	for _, ingress := range ingresses {
		ingressSpec := ingress.Spec
		maybePrependRegexPrefixFn := translators.MaybePrependRegexPrefixForIngressV1Fn(ingress, icp.EnableLegacyRegexDetection && featureFlags.RegexPathPrefix)
		routeName := routeNamer(failuresCollector, ingress)
		for i, rule := range ingressSpec.Rules {
			if rule.HTTP == nil {
				continue
			}
			for j, rulePath := range rule.HTTP.Paths {
				rulePath := rulePath
				pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
				if rulePath.PathType == nil {
					rulePath.PathType = &pathTypeImplementationSpecific
				}

				paths := translators.PathsFromIngressPaths(rulePath, featureFlags.RegexPathPrefix)
				if paths == nil {
					// Registering a failure, but technically it should never happen thanks to Kubernetes API validations.
					failuresCollector.PushResourceFailure(fmt.Sprintf("could not translate Ingress Path %s to Kong paths", rulePath.Path), ingress)
					continue
				}

				for i, path := range paths {
					paths[i] = maybePrependRegexPrefixFn(*path)
				}

				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
						Name:              routeName(i, j),
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
				parsedObjectsCollector.Add(ingress) // Register successfully parsed object.
			}
		}
	}

	return servicesCache
}

// routeNamer returns function for generating a name for a Kong Route based on the Ingress name, namespace, rule index and path index.
// If the name won't be unique, it registers a failure.
func routeNamer(failuresCollector *failures.ResourceFailuresCollector, objIngress client.Object) func(ruleIndex, pathIndex int) *string {
	uniqueRouteNames := make(map[string]struct{})
	return func(ruleIndex, pathIndex int) *string {
		routeName := kong.String(fmt.Sprintf("%s.%s.%d%d", objIngress.GetNamespace(), objIngress.GetName(), ruleIndex, pathIndex))
		if _, conflicting := uniqueRouteNames[*routeName]; conflicting {
			failuresCollector.PushResourceFailure(
				fmt.Sprint(
					"ERROR: Kong route with conflicting name: ", *routeName, " ",
					"use feature gate CombinedRoutes=true ",
					"or update Kong Kubernetes Ingress Controller version to 3.0.0 or above ",
					"(both remediation changes naming schema of Kong routes)",
				),
				objIngress,
			)
		}
		uniqueRouteNames[*routeName] = struct{}{}
		return routeName
	}
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
