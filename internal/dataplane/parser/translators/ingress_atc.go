package translators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// TranslateIngressATC receives a Kubernetes ingress object and from it will
// produce a translated set of kong.Services and expression based kong.Routes
// which will come wrapped in a kongstate.Service object.
func TranslateIngressATC(ingress *netv1.Ingress) []*kongstate.Service {
	index := &ingressTranslationIndex{
		cache: make(map[string]*ingressTranslationMeta),
	}
	index.add(ingress)
	kongStateServices := kongstate.Services(index.translateATC())
	sort.Sort(kongStateServices)
	return kongStateServices
}

func (i *ingressTranslationIndex) translateATC() []*kongstate.Service {
	kongStateServiceCache := make(map[string]*kongstate.Service)
	for _, meta := range i.cache {
		kongServiceName := fmt.Sprintf("%s.%s.%s.%s", meta.parentIngress.GetNamespace(), meta.parentIngress.GetName(), meta.serviceName, meta.servicePort.CanonicalString())
		kongStateService, ok := kongStateServiceCache[kongServiceName]
		if !ok {
			kongStateService = meta.translateIntoKongStateService(kongServiceName, meta.servicePort)
		}

		route := meta.translateIntoKongATCRoutes()
		kongStateService.Routes = append(kongStateService.Routes, *route)

		kongStateServiceCache[kongServiceName] = kongStateService
	}

	kongStateServices := make([]*kongstate.Service, 0, len(kongStateServiceCache))
	for _, kongStateService := range kongStateServiceCache {
		kongStateServices = append(kongStateServices, kongStateService)
	}

	return kongStateServices
}

func (m *ingressTranslationMeta) translateIntoKongATCRoutes() *kongstate.Route {

	routeName := fmt.Sprintf("%s.%s.%s.%s.%s",
		m.parentIngress.GetNamespace(),
		m.parentIngress.GetName(),
		m.serviceName,
		// '_' is not allowed in host, so we use '_' to replace '*' since '*' is not allowed in Kong.
		strings.ReplaceAll(m.ingressHost, "*", "_"),
		m.servicePort.CanonicalString(),
	)
	route := &kongstate.Route{
		Ingress: util.K8sObjectInfo{
			Namespace:   m.parentIngress.GetNamespace(),
			Name:        m.parentIngress.GetName(),
			Annotations: m.parentIngress.GetAnnotations(),
		},
		Route: kong.Route{
			Name:              kong.String(routeName),
			StripPath:         kong.Bool(false),
			PreserveHost:      kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
			ResponseBuffering: kong.Bool(true),
			Tags:              m.ingressTags,
		},
	}

	matchRules := &atc.MatchRules{
		Protocols: []string{"http", "https"},
	}

	if m.ingressHost != "" {
		if strings.HasPrefix(m.ingressHost, "*") {
			matchRules.Hosts = append(matchRules.Hosts, atc.MatchRuleHost{
				Type: atc.HostMatchWildcard,
				Host: strings.TrimPrefix(m.ingressHost, "*"),
			})
		} else {
			matchRules.Hosts = append(matchRules.Hosts, atc.MatchRuleHost{
				Type: atc.HostMatchExact,
				Host: m.ingressHost,
			})
		}
	}

	for _, httpIngressPath := range m.paths {
		paths := ATCPathMatchesFromIngressPath(httpIngressPath, ControllerPathRegexPrefix)
		matchRules.Paths = append(matchRules.Paths, paths...)
	}

	atc.ApplyExpression(&route.Route, matchRules, 1)

	return route
}

func ATCPathMatchesFromIngressPath(httpIngressPath netv1.HTTPIngressPath, regexPrefix string) []atc.MatchRulePath {
	if httpIngressPath.PathType == nil {
		return nil
	}

	switch *httpIngressPath.PathType {
	case netv1.PathTypePrefix:
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			return []atc.MatchRulePath{
				{Type: atc.PathMatchPrefix, Path: "/"},
			}
		} else {
			return []atc.MatchRulePath{
				{Type: atc.PathMatchPrefix, Path: "/" + base + "/"},
				{Type: atc.PathMatchExact, Path: "/" + base},
			}
		}
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		return []atc.MatchRulePath{
			{Type: atc.PathMatchExact, Path: "/" + relative},
		}
	case netv1.PathTypeImplementationSpecific:
		if httpIngressPath.Path == "" {
			return []atc.MatchRulePath{
				{Type: atc.PathMatchPrefix, Path: "/"},
			}
		} else if strings.HasPrefix(httpIngressPath.Path, regexPrefix) {
			// starts with specified regex prefix, translate to regex match on path.
			regexPath := strings.TrimPrefix(httpIngressPath.Path, regexPrefix)
			// The existing regex path matches in Kong traditional/tranditional_compatible
			// applies on prefix of path, but Kong expressions router matches in any place.
			// so we need to prepend a `^` if there is no existing one.
			if !strings.HasPrefix(regexPath, "^") {
				regexPath = "^" + regexPath
			}
			return []atc.MatchRulePath{
				{Type: atc.PathMatchRegex, Path: regexPath},
			}
		} else {
			return []atc.MatchRulePath{
				{Type: atc.PathMatchPrefix, Path: httpIngressPath.Path},
			}
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
}
