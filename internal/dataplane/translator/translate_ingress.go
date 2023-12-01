package translator

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
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
	defaultBackendService, ok := getDefaultBackendService(allDefaultBackends, t.featureFlags.ExpressionRoutes)
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
func getDefaultBackendService(allDefaultBackends []netv1.Ingress, expressionRoutes bool) (kongstate.Service, bool) {
	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := allDefaultBackends[0].Spec.DefaultBackend
		port := subtranslator.PortDefFromServiceBackendPort(&defaultBackend.Service.Port)
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
		atc.ApplyExpression(&r.Route, catchAllMatcher, subtranslator.IngressDefaultBackendPriority)
	} else {
		r.Route.Paths = kong.StringSlice("/")
		r.Route.Protocols = kong.StringSlice("http", "https")
		r.Route.RegexPriority = kong.Int(0)
	}
	return r
}
