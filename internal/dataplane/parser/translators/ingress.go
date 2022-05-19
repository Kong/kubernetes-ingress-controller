package translators

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kong/go-kong/kong"
	networkingv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Ingress Translation - Public Functions
// -----------------------------------------------------------------------------

// TranslateIngress receives a Kubernetes ingress object and from it will
// produce a translated set of kong.Services and kong.Routes which will come
// wrapped in a kongstate.Service object.
func TranslateIngress(ingress *networkingv1.Ingress) []*kongstate.Service {
	index := &ingressTranslationIndex{cache: make(map[string]*ingressTranslationMeta)}
	index.add(ingress)
	kongStateServices := kongstate.Services(index.translate())
	sort.Sort(kongStateServices)
	return kongStateServices
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private Consts & Vars
// -----------------------------------------------------------------------------

var defaultHTTPIngressPathType = networkingv1.PathTypePrefix

const (
	defaultHTTPPort = 80
	defaultRetries  = 5

	// defaultServiceTimeout indicates the amount of time by default that we wait
	// for connections to an underlying Kubernetes service to complete in the
	// data-plane. The current value is based on a historical default that started
	// in version 0 of the ingress controller.
	defaultServiceTimeout = time.Second * 60
)

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Index
// -----------------------------------------------------------------------------

// ingressTranslationIndex is a de-duplicating index of the contents of a
// Kubernetes Ingress resource, where the key is a combination of data from
// that resource which makes it unique for the purpose of translating it into
// kong.Services and kong.Routes and the value is a combination of various
// metadata needed to configure and name those kong.Services and kong.Routes
// plus the URL paths which make those routes actionable. This index is used
// to enable compiling a minimal set of kong.Routes when translating into
// Kong resources for each rule in the ingress spec, where the combination of:
//
// - ingress.Namespace
// - ingress.Name
// - host for the Ingress rule
// - Kubernetes Service for the Ingress rule
// - the port for the Kubernetes Service
//
// are unique. For ingress spec rules which are not unique along those
// data-points, a separate kong.Service and separate kong.Routes will be created
// for each unique combination.
type ingressTranslationIndex struct {
	cache map[string]*ingressTranslationMeta
}

func (i *ingressTranslationIndex) add(ingress *networkingv1.Ingress) {
	for _, ingressRule := range ingress.Spec.Rules {
		if ingressRule.HTTP == nil || len(ingressRule.HTTP.Paths) < 1 {
			continue
		}

		for _, httpIngressPath := range ingressRule.HTTP.Paths {
			httpIngressPath.Path = flattenMultipleSlashes(httpIngressPath.Path)

			if httpIngressPath.Path == "" {
				httpIngressPath.Path = "/"
			}

			if httpIngressPath.PathType == nil {
				httpIngressPath.PathType = &defaultHTTPIngressPathType
			}

			serviceName := httpIngressPath.Backend.Service.Name
			servicePort := httpIngressPath.Backend.Service.Port.Number

			cacheKey := fmt.Sprintf("%s%s%s%s%d", ingress.Namespace, ingress.Name, ingressRule.Host, serviceName, servicePort)
			meta, ok := i.cache[cacheKey]
			if !ok {
				meta = &ingressTranslationMeta{
					ingressNamespace: ingress.Namespace,
					ingressName:      ingress.Name,
					ingressHost:      ingressRule.Host,
					serviceName:      serviceName,
					servicePort:      servicePort,
				}
			}

			meta.paths = append(meta.paths, httpIngressPath)
			meta.ingressAnnotations = ingress.Annotations
			i.cache[cacheKey] = meta
		}
	}
}

func (i *ingressTranslationIndex) translate() []*kongstate.Service {
	kongStateServiceCache := make(map[string]*kongstate.Service)
	for _, meta := range i.cache {
		portDef := kongstate.PortDef{
			Number: meta.servicePort,
			Mode:   kongstate.PortModeByNumber,
		}

		kongServiceName := fmt.Sprintf("%s.%s.%s.%s", meta.ingressNamespace, meta.ingressName, meta.serviceName, portDef.CanonicalString())
		kongStateService, ok := kongStateServiceCache[kongServiceName]
		if !ok {
			kongStateService = meta.translateIntoKongStateService(kongServiceName, portDef)
		}

		route := meta.translateIntoKongRoutes()
		kongStateService.Routes = append(kongStateService.Routes, *route)

		kongStateServiceCache[kongServiceName] = kongStateService
	}

	kongStateServices := make([]*kongstate.Service, 0, len(kongStateServiceCache))
	for _, kongStateService := range kongStateServiceCache {
		kongStateServices = append(kongStateServices, kongStateService)
	}

	return kongStateServices
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Metadata
// -----------------------------------------------------------------------------

type ingressTranslationMeta struct {
	ingressAnnotations map[string]string
	ingressNamespace   string
	ingressName        string
	ingressHost        string
	serviceName        string
	servicePort        int32
	paths              []networkingv1.HTTPIngressPath
}

func (m *ingressTranslationMeta) translateIntoKongStateService(kongServiceName string, portDef kongstate.PortDef) *kongstate.Service {
	return &kongstate.Service{
		Namespace: m.ingressNamespace,
		Service: kong.Service{
			Name:           kong.String(kongServiceName),
			Host:           kong.String(fmt.Sprintf("%s.%s.%d.svc", m.serviceName, m.ingressNamespace, portDef.Number)),
			Port:           kong.Int(defaultHTTPPort),
			Protocol:       kong.String("http"),
			Path:           kong.String("/"),
			ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
			ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
			WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
			Retries:        kong.Int(defaultRetries),
		},
		Backends: []kongstate.ServiceBackend{{
			Name:      m.serviceName,
			Namespace: m.ingressNamespace,
			PortDef:   portDef,
		}},
	}
}

func (m *ingressTranslationMeta) translateIntoKongRoutes() *kongstate.Route {
	routeName := fmt.Sprintf("%s.%s.%s.%s.%d", m.ingressNamespace, m.ingressName, m.serviceName, m.ingressHost, m.servicePort)
	route := &kongstate.Route{
		Ingress: util.K8sObjectInfo{
			Namespace:   m.ingressNamespace,
			Name:        m.ingressName,
			Annotations: m.ingressAnnotations,
		},
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

	if m.ingressHost != "" {
		route.Route.Hosts = append(route.Route.Hosts, kong.String(m.ingressHost))
	}

	for _, httpIngressPath := range m.paths {
		paths := pathsFromIngressPaths(httpIngressPath)
		route.Paths = append(route.Paths, paths...)
	}

	return route
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Helper Functions
// -----------------------------------------------------------------------------

func pathsFromIngressPaths(httpIngressPath networkingv1.HTTPIngressPath) []*string {
	switch *httpIngressPath.PathType { //nolint:exhaustive
	case networkingv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		if httpIngressPath.Path == "" {
			return kong.StringSlice("/")
		}
		return kong.StringSlice("/" + relative + "$")
	case networkingv1.PathTypeImplementationSpecific:
		return kong.StringSlice(httpIngressPath.Path)
	default:
		// path type is prefix
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			return kong.StringSlice("/")
		}
		return kong.StringSlice(
			"/"+base+"$",
			"/"+base+"/",
		)
	}
}

func flattenMultipleSlashes(path string) string {
	var out []rune
	in := []rune(path)
	for i := 0; i < len(in); i++ {
		c := in[i]
		if c == '/' {
			for i < len(in)-1 && in[i+1] == '/' {
				i++
			}
		}
		out = append(out, c)
	}
	return string(out)
}
