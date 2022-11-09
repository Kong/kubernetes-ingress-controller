package translators

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Ingress Translation - Public Functions
// -----------------------------------------------------------------------------

// TranslateIngress receives a Kubernetes ingress object and from it will
// produce a translated set of kong.Services and kong.Routes which will come
// wrapped in a kongstate.Service object.
func TranslateIngress(ingress *netv1.Ingress, addRegexPrefix bool) []*kongstate.Service {
	index := &ingressTranslationIndex{
		cache:          make(map[string]*ingressTranslationMeta),
		addRegexPrefix: addRegexPrefix,
	}
	index.add(ingress)
	kongStateServices := kongstate.Services(index.translate())
	sort.Sort(kongStateServices)
	return kongStateServices
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private Consts & Vars
// -----------------------------------------------------------------------------

var defaultHTTPIngressPathType = netv1.PathTypeImplementationSpecific

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
//
// The addRegexPrefix flag indicates if generated regex paths for path type handling include the Kong 3.0+ "~" regular
// expression prefix.
type ingressTranslationIndex struct {
	cache          map[string]*ingressTranslationMeta
	addRegexPrefix bool
}

func (i *ingressTranslationIndex) add(ingress *netv1.Ingress) {
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

			cacheKey := fmt.Sprintf("%s.%s.%s.%s.%d", ingress.Namespace, ingress.Name, ingressRule.Host, serviceName, servicePort)
			meta, ok := i.cache[cacheKey]
			if !ok {
				meta = &ingressTranslationMeta{
					ingressNamespace: ingress.Namespace,
					ingressName:      ingress.Name,
					ingressHost:      ingressRule.Host,
					serviceName:      serviceName,
					servicePort:      servicePort,
					addRegexPrefix:   i.addRegexPrefix,
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
	paths              []netv1.HTTPIngressPath
	addRegexPrefix     bool
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
		Parent: &netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Ingress",
				APIVersion: netv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.ingressName,
				Namespace: m.ingressNamespace,
			},
		},
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
		paths := PathsFromIngressPaths(httpIngressPath, m.addRegexPrefix)
		route.Paths = append(route.Paths, paths...)
	}

	return route
}

// -----------------------------------------------------------------------------
// Ingress Translation - Private - Helper Functions
// -----------------------------------------------------------------------------

// TODO this is exported because most of the parser translate functions are still in the parser package. if/when we
// refactor to move them here, this should become private.

// PathsFromIngressPaths takes a path and Ingress path type and returns a set of Kong route paths that satisfy that path
// type. It optionally adds the Kong 3.x regex path prefix for path types that require a regex path. It rejects
// unknown path types with an error.
func PathsFromIngressPaths(httpIngressPath netv1.HTTPIngressPath, addRegexPrefix bool) []*string {
	routePaths := []string{}
	routeRegexPaths := []string{}
	if httpIngressPath.PathType == nil {
		return nil
	}

	switch *httpIngressPath.PathType {
	case netv1.PathTypePrefix:
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, "/"+base+"/")
			routeRegexPaths = append(routeRegexPaths, "/"+base+"$")
		}
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		routeRegexPaths = append(routeRegexPaths, "/"+relative+"$")
	case netv1.PathTypeImplementationSpecific:
		if httpIngressPath.Path == "" {
			routePaths = append(routePaths, "/")
		} else {
			routePaths = append(routePaths, httpIngressPath.Path)
		}
	default:
		// the default case here is mostly to provide a home for this comment: we explicitly do not handle unknown
		// PathTypes, and leave it up to the callers if they want to handle empty responses. barring spec changes,
		// however, this should not be a concern: Kubernetes rejects any Ingress with an unknown PathType already, so
		// none should ever end up here. prior versions of this function returned an error in this case, but it
		// should be unnecessary in practice and not returning one simplifies the call chain above (this would be the
		// only part of translation that can error)
		return nil
	}

	if addRegexPrefix {
		for i, orig := range routeRegexPaths {
			routeRegexPaths[i] = KongPathRegexPrefix + orig
		}
	}
	routePaths = append(routePaths, routeRegexPaths...)
	return kong.StringSlice(routePaths...)
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
