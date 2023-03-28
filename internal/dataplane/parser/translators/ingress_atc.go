package translators

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	// TODO: copied from kongstate package. should move that to some common place.
	validMethods = regexp.MustCompile(`\A[A-Z]+$`)
	validSNIs    = regexp.MustCompile(`^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*$`)
	validHosts   = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*?(\.\*)?$`)
)

const (
	KongHeaderRegexPrefix = "~*"
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

	regexPrefix := annotations.ExtractRegexPrefix(route.Ingress.Annotations)
	if regexPrefix == "" {
		regexPrefix = ControllerPathRegexPrefix
	}

	for _, httpIngressPath := range m.paths {
		paths := atcPathMatchesFromIngressPath(httpIngressPath, regexPrefix)
		matchRules.Paths = append(matchRules.Paths, paths...)
	}

	matchRules = overrideIngressMatchRulesByAnnotations(matchRules, route.Ingress.Annotations)

	atc.ApplyExpression(&route.Route, matchRules, 1)

	overrideRouteOptionsByAnnotation(&route.Route, route.Ingress.Annotations)

	return route
}

func atcPathMatchesFromIngressPath(httpIngressPath netv1.HTTPIngressPath, regexPrefix string) []atc.MatchRulePath {
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
		}

		// starts with specified regex prefix, translate to regex match on path.
		if strings.HasPrefix(httpIngressPath.Path, regexPrefix) {
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
		}
		return []atc.MatchRulePath{
			{Type: atc.PathMatchPrefix, Path: httpIngressPath.Path},
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

func overrideIngressMatchRulesByAnnotations(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {

	// override protocols.
	matchRules = overrideIngressMatchRuleProtocolsByAnnotation(matchRules, anns)
	// override methods.
	matchRules = overrideIngressMatchRuleMethodsByAnnotation(matchRules, anns)
	// override host aliases.
	matchRules = overrideIngressMatchRuleHostsByAnnotation(matchRules, anns)
	// override SNIs.
	matchRules = overrideIngressMatchRuleSNIsByAnnotation(matchRules, anns)
	// override headers.
	matchRules = overrideIngressMatchRuleHeadersByAnnotation(matchRules, anns)

	return matchRules
}

func overrideIngressMatchRuleProtocolsByAnnotation(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {
	matchProtocols := []string{}
	for _, prot := range annotations.ExtractProtocolNames(anns) {
		if !util.ValidateProtocol(prot) {
			return matchRules
		}
		matchProtocols = append(matchProtocols, prot)
	}
	// TODO: update protocols if force SSL redirection enabled by annotation.

	if len(matchProtocols) > 0 {
		matchRules.Protocols = matchProtocols
	}
	return matchRules
}

func overrideIngressMatchRuleMethodsByAnnotation(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {

	var methods []string
	for _, method := range annotations.ExtractMethods(anns) {
		sanitizedMethod := strings.TrimSpace(strings.ToUpper(method))
		if validMethods.MatchString(sanitizedMethod) {
			methods = append(methods, sanitizedMethod)
		} else {
			return matchRules
		}
	}

	if len(methods) > 0 {
		matchRules.Methods = methods
	}
	return matchRules
}

func overrideIngressMatchRuleHostsByAnnotation(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {
	annHostAliases, exists := annotations.ExtractHostAliases(anns)
	if !exists {
		// the annotation is not set, quit
		return matchRules
	}

	hostMatches := []atc.MatchRuleHost{}
	for _, hostAlias := range annHostAliases {
		sanitizedHost := strings.TrimSpace(hostAlias)
		if !validHosts.MatchString(sanitizedHost) {
			// annotation contains invalid host, return directly without updating match rules.
			return matchRules
		}

		if strings.HasPrefix(sanitizedHost, "*") {
			hostMatches = append(hostMatches, atc.MatchRuleHost{
				Type: atc.HostMatchWildcard,
				Host: strings.TrimPrefix(sanitizedHost, "*"),
			})
		} else {
			hostMatches = append(hostMatches, atc.MatchRuleHost{
				Type: atc.HostMatchExact,
				Host: sanitizedHost,
			})
		}
	}

	if len(hostMatches) > 0 {
		matchRules.Hosts = append(matchRules.Hosts, hostMatches...)
	}
	return matchRules
}

func overrideIngressMatchRuleSNIsByAnnotation(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {
	annSNIs, exists := annotations.ExtractSNIs(anns)
	// this is not a length check because we specifically want to provide a means
	// to set "no SNI criteria", by providing the annotation with an empty string value
	if !exists {
		return matchRules
	}
	var snis []string
	for _, sni := range annSNIs {
		sanitizedSNI := strings.TrimSpace(sni)
		if !validSNIs.MatchString(sanitizedSNI) {
			return matchRules
		}
		snis = append(snis, sanitizedSNI)
	}

	matchRules.SNIs = snis
	return matchRules
}

func overrideIngressMatchRuleHeadersByAnnotation(matchRules *atc.MatchRules, anns map[string]string) *atc.MatchRules {
	headers, exists := annotations.ExtractHeaders(anns)
	if !exists {
		return matchRules
	}

	if matchRules.Headers == nil {
		matchRules.Headers = map[string]atc.MatchRuleHeader{}
	}

	for headerName, values := range headers {
		// should not use header match on "Host" header
		if strings.ToLower(headerName) == "host" {
			continue
		}
		headerMatch := atc.MatchRuleHeader{
			Type: atc.HeaderMatchExact,
		}
		for _, v := range values {
			if strings.HasPrefix(v, KongHeaderRegexPrefix) {
				headerMatch.Type = atc.HeaderMatchRegex
				headerMatch.Values = []string{
					strings.TrimPrefix(v, KongHeaderRegexPrefix),
				}
				break
			}
			headerMatch.Values = append(headerMatch.Values, v)
		}
		if len(headerMatch.Values) > 0 {
			matchRules.Headers[headerName] = headerMatch
		}
	}

	return matchRules
}

func overrideRouteOptionsByAnnotation(r *kong.Route, anns map[string]string) {
	overrideRouteStripPath(r, anns)
	overrideHTTPSRedirectCode(r, anns)
	overrideRequestBuffering(r, anns)
	overrideResponseBuffering(r, anns)
}

func overrideRouteStripPath(r *kong.Route, anns map[string]string) {
	if r == nil {
		return
	}

	stripPathValue := annotations.ExtractStripPath(anns)
	if stripPathValue == "" {
		return
	}

	enabled, err := strconv.ParseBool(strings.ToLower(stripPathValue))
	if err != nil {
		return
	}
	r.StripPath = kong.Bool(enabled)
}

func overrideHTTPSRedirectCode(r *kong.Route, anns map[string]string) {
	if annotations.HasForceSSLRedirectAnnotation(anns) {
		r.HTTPSRedirectStatusCode = kong.Int(302)
	}

	code := annotations.ExtractHTTPSRedirectStatusCode(anns)
	if code == "" {
		return
	}
	statusCode, err := strconv.Atoi(code)
	if err != nil {
		return
	}
	if statusCode != 426 &&
		statusCode != 301 &&
		statusCode != 302 &&
		statusCode != 307 &&
		statusCode != 308 {
		return
	}

	r.HTTPSRedirectStatusCode = kong.Int(statusCode)
}

func overrideRequestBuffering(r *kong.Route, anns map[string]string) {
	annotationValue, ok := annotations.ExtractRequestBuffering(anns)
	if !ok {
		// the annotation is not set, quit
		return
	}

	enabled, err := strconv.ParseBool(strings.ToLower(annotationValue))
	if err != nil {
		// the value provided is not a parseable boolean, quit
		return
	}
	r.RequestBuffering = kong.Bool(enabled)
}

func overrideResponseBuffering(r *kong.Route, anns map[string]string) {
	annotationValue, ok := annotations.ExtractResponseBuffering(anns)
	if !ok {
		// the annotation is not set, quit
		return
	}

	enabled, err := strconv.ParseBool(strings.ToLower(annotationValue))
	if err != nil {
		// the value provided is not a parseable boolean, quit
		return
	}

	r.ResponseBuffering = kong.Bool(enabled)
}
