package translator

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func (t *Translator) ingressRulesFromIngressV1() ingressRules {
	result := newIngressRules()

	ingressList := t.storer.ListIngressesV1()
	icp, err := getIngressClassParametersOrDefault(t.storer)
	if err != nil {
		if !errors.As(err, &store.NotFoundError{}) {
			// anything else is unexpected
			t.logger.Error(err, "Could not retrieve IngressClassParameters, using defaults")
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
	servicesCache := subtranslator.TranslateIngresses(
		ingressList,
		icp,
		subtranslator.TranslateIngressFeatureFlags{
			ExpressionRoutes:  t.featureFlags.ExpressionRoutes,
			KongServiceFacade: t.featureFlags.KongServiceFacade,
		},
		t.translatedObjectsCollector,
		t.failuresCollector,
		t.storer,
	)
	for i := range servicesCache {
		service := servicesCache[i]
		if err := subtranslator.MaybeRewriteURI(&service, t.featureFlags.RewriteURIs); err != nil {
			t.registerTranslationFailure(err.Error(), service.Parent)
			continue
		}

		result.ServiceNameToServices[*service.Name] = service
		result.ServiceNameToParent[*service.Name] = service.Parent
	}

	// Add a default backend if it exists.
	defaultBackendService, ok := getDefaultBackendService(t.storer, t.failuresCollector, allDefaultBackends, t.featureFlags)
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

// getDefaultBackendService picks the oldest Ingress with a DefaultBackend defined and returns a Kong Service for it.
func getDefaultBackendService(
	storer store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	allDefaultBackends []netv1.Ingress,
	features FeatureFlags,
) (kongstate.Service, bool) {
	// Sort the default backends by creation timestamp, so that the oldest one is picked.
	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := ingress.Spec.DefaultBackend
		route := translateIngressDefaultBackendRoute(&ingress, util.GenerateTagsForObject(&ingress), features.ExpressionRoutes)

		// If the default backend is defined as an arbitrary resource, then we need handle it differently.
		if resource := defaultBackend.Resource; resource != nil {
			return translateIngressDefaultBackendResource(
				resource,
				ingress,
				route,
				storer,
				failuresCollector,
				features,
			)
		}

		// Otherwise, the default backend is defined as a Kubernetes Service.
		return translateIngressDefaultBackendService(ingress, route)
	}

	return kongstate.Service{}, false
}

func translateIngressDefaultBackendResource(
	resource *corev1.TypedLocalObjectReference,
	ingress netv1.Ingress,
	route *kongstate.Route,
	storer store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	features FeatureFlags,
) (kongstate.Service, bool) {
	if !subtranslator.IsKongServiceFacade(resource) {
		gk := resource.Kind
		if resource.APIGroup != nil {
			gk = *resource.APIGroup + "/" + gk
		}
		failuresCollector.PushResourceFailure(fmt.Sprintf("default backend: unsupported resource type %s", gk), &ingress)
		return kongstate.Service{}, false
	}
	if !features.KongServiceFacade {
		failuresCollector.PushResourceFailure(
			fmt.Sprintf("default backend: KongServiceFacade is not enabled, please set the %q feature gate to 'true' to enable it", featuregates.KongServiceFacade),
			&ingress,
		)
		return kongstate.Service{}, false
	}
	facade, err := storer.GetKongServiceFacade(ingress.Namespace, resource.Name)
	if err != nil {
		failuresCollector.PushResourceFailure(
			fmt.Sprintf("default backend: KongServiceFacade %q could not be fetched: %s", resource.Name, err),
			&ingress,
		)
		return kongstate.Service{}, false
	}

	serviceName := fmt.Sprintf("%s.%s.svc.facade", ingress.Namespace, resource.Name)
	return kongstate.Service{
		Service: kong.Service{
			Name:           kong.String(serviceName),
			Host:           kong.String(serviceName),
			Port:           kong.Int(DefaultHTTPPort),
			Protocol:       kong.String("http"),
			ConnectTimeout: kong.Int(DefaultServiceTimeout),
			ReadTimeout:    kong.Int(DefaultServiceTimeout),
			WriteTimeout:   kong.Int(DefaultServiceTimeout),
			Retries:        kong.Int(DefaultRetries),
			// We do not populate Service's Tags field here because it would get overridden anyway later in the
			// Translator pipeline (see ingressRules.generateKongServiceTags).
		},
		Namespace: ingress.Namespace,
		Backends: []kongstate.ServiceBackend{{
			Type:      kongstate.ServiceBackendTypeKongServiceFacade,
			Name:      resource.Name,
			Namespace: ingress.Namespace,
			PortDef:   subtranslator.PortDefFromPortNumber(facade.Spec.Backend.Port),
		}},
		Parent: facade,
		Routes: []kongstate.Route{*route},
	}, true
}

func translateIngressDefaultBackendService(ingress netv1.Ingress, route *kongstate.Route) (kongstate.Service, bool) {
	defaultBackend := ingress.Spec.DefaultBackend
	port := subtranslator.PortDefFromServiceBackendPort(&defaultBackend.Service.Port)
	serviceName := fmt.Sprintf(
		"%s.%s.%s",
		ingress.Namespace,
		defaultBackend.Service.Name,
		port.CanonicalString(),
	)
	return kongstate.Service{
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
			// We do not populate Service's Tags field here because it would get overridden anyway later in the
			// Translator pipeline (see ingressRules.generateKongServiceTags).
		},
		Namespace: ingress.Namespace,
		Backends: []kongstate.ServiceBackend{{
			Name:    defaultBackend.Service.Name,
			PortDef: port,
		}},
		Parent: &ingress,
		Routes: []kongstate.Route{*route},
	}, true
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
		atc.ApplyExpression(&r.Route, catchAllMatcher, subtranslator.IngressDefaultBackendPriority)
	} else {
		r.Route.Paths = kong.StringSlice("/")
		r.Route.Protocols = kong.StringSlice("http", "https")
		r.Route.RegexPriority = kong.Int(0)
	}
	return r
}
