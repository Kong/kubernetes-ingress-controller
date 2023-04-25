package translators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	headerAnnotationRegexPrefix = "~*"
	// TODO: list all available matchers instead.
	validMethods = regexp.MustCompile(`\A[A-Z]+$`)

	// hostnames are complicated. shamelessly cribbed from https://stackoverflow.com/a/18494710
	// TODO if the Kong core adds support for wildcard SNI route match criteria, this should change.
	validSNIs  = regexp.MustCompile(`^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*$`)
	validHosts = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*))*?(\.\*)?$`)
)

var (
	NormalIngressExpressionPriority = 1
	IngressDefaultBackendPriority   = 0
)

func (m *ingressTranslationMeta) translateIntoKongExpressionRoutes() *kongstate.Route {
	ingressHost := m.ingressHost
	if strings.Contains(ingressHost, "*") {
		// '_' is not allowed in host, so we use '_' to replace '*' since '*' is not allowed in Kong.
		ingressHost = strings.ReplaceAll(ingressHost, "*", "_")
	}

	routeName := fmt.Sprintf("%s.%s.%s.%s.%s", m.parentIngress.GetNamespace(), m.parentIngress.GetName(), m.serviceName, ingressHost, m.servicePort.CanonicalString())
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
		ExpressionRoutes: true,
	}

	ingressAnnotations := m.parentIngress.GetAnnotations()

	routeMatcher := atc.And()
	// translate hosts.
	hosts := []string{}
	if m.ingressHost != "" {
		hosts = append(hosts, m.ingressHost)
	}
	hostAliases, _ := annotations.ExtractHostAliases(ingressAnnotations)
	for _, hostAlias := range hostAliases {
		hosts = append(hosts, hostAlias)
	}
	hostMatcher := hostMatcherFromIngressHosts(hosts)
	routeMatcher.And(hostMatcher)

	// translate paths.
	pathMatchers := make([]atc.Matcher, 0, len(m.paths))
	for _, path := range m.paths {
		pathMatchers = append(pathMatchers, pathMatcherFromIngressPath(path, ControllerPathRegexPrefix))
	}
	routeMatcher.And(atc.Or(pathMatchers...))

	// default protocols
	protocols := []string{"http", "https"}
	annonationProtocols := annotations.ExtractProtocolNames(ingressAnnotations)
	if len(annonationProtocols) > 0 {
		protocols = annonationProtocols
	}
	protocolMatcher := protocolMatcherFromProtocols(protocols)
	routeMatcher.And(protocolMatcher)

	headers, exist := annotations.ExtractHeaders(ingressAnnotations)
	if len(headers) > 0 && exist {
		headerMatcher := headerMatcherFromHeaders(headers)
		routeMatcher.And(headerMatcher)
	}

	methods := annotations.ExtractMethods(ingressAnnotations)
	if len(methods) > 0 {
		methodMatcher := methodMatcherFromMethods(methods)
		routeMatcher.And(methodMatcher)
	}

	snis, exist := annotations.ExtractSNIs(ingressAnnotations)
	if exist && len(snis) > 0 {
		sniMatcher := sniMatcherFromSNIs(snis)
		routeMatcher.And(sniMatcher)
	}

	atc.ApplyExpression(&route.Route, routeMatcher, NormalIngressExpressionPriority)
	return route
}

// hostMatcherFromIngressHosts translates hosts in IngressHost format to ATC matcher that matches any of them.
func hostMatcherFromIngressHosts(hosts []string) atc.Matcher {
	matchers := make([]atc.Matcher, 0, len(hosts))
	for _, host := range hosts {
		if !validHosts.MatchString(host) {
			continue
		}

		if strings.HasPrefix(host, "*") {
			// wildcard match on hosts, genreate a prefix match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpPrefixMatch, strings.TrimPrefix(host, "*")))
		} else {
			// exact match on hosts, generate an exact match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpEqual, host))
		}
	}
	return atc.Or(matchers...)
}

func pathMatcherFromIngressPath(httpIngressPath netv1.HTTPIngressPath, regexPathPrefix string) atc.Matcher {
	switch *httpIngressPath.PathType {
	case netv1.PathTypePrefix:
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		return atc.Or(
			atc.NewPredicateHTTPPath(atc.OpEqual, "/"+base+"/"),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/"+base),
		)
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		return atc.NewPredicateHTTPPath(atc.OpEqual, "/"+relative)
	case netv1.PathTypeImplementationSpecific:
		if httpIngressPath.Path == "" {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		if regexPathPrefix != "" && strings.HasPrefix(httpIngressPath.Path, regexPathPrefix) {
			regex := strings.TrimPrefix(httpIngressPath.Path, regexPathPrefix)
			if !strings.HasPrefix(regex, "^") {
				regex = "^" + regex
			}
			return atc.NewPredicateHTTPPath(atc.OpRegexMatch, regex)
		}

		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, httpIngressPath.Path)

	}

	return nil
}

func protocolMatcherFromProtocols(protocols []string) atc.Matcher {
	if len(protocols) == 0 {
		return nil
	}
	matchers := []atc.Matcher{}
	for _, protocol := range protocols {
		if !util.ValidateProtocol(protocol) {
			continue
		}
		matchers = append(matchers, atc.NewPredicateNetProtocol(atc.OpEqual, protocol))
	}
	return atc.Or(matchers...)
}

func headerMatcherFromHeaders(headers map[string][]string) atc.Matcher {

	matchers := make([]atc.Matcher, 0, len(headers))
	for headerName, values := range headers {
		// transfer header name to lowercase and replace "-" with "_"
		headerName = strings.ReplaceAll(strings.ToLower(headerName), "-", "_")
		// header "Host" should be skipped, they are processed in "http.host".
		if headerName == "host" {
			continue
		}
		if len(values) == 0 {
			continue
		}

		singleHeaderMatcher := atc.Or()
		for _, val := range values {
			// generate a predicate using regex match if value starts with the special prefix "~*".
			if strings.HasPrefix(val, headerAnnotationRegexPrefix) {
				regex := strings.TrimPrefix(val, headerAnnotationRegexPrefix)
				singleHeaderMatcher.Or(atc.NewPredicateHTTPHeader(headerName, atc.OpRegexMatch, regex))
			} else {
				// otherwise, genereate a predicate using exact match.
				singleHeaderMatcher.Or(atc.NewPredicateHTTPHeader(headerName, atc.OpEqual, val))
			}
		}
		matchers = append(matchers, singleHeaderMatcher)
	}

	return atc.And(matchers...)
}

func methodMatcherFromMethods(methods []string) atc.Matcher {
	matchers := make([]atc.Matcher, 0, len(methods))
	for _, method := range methods {
		if !validMethods.MatchString(method) {
			continue
		}
		matchers = append(matchers, atc.NewPredicateHTTPMethod(atc.OpEqual, method))
	}
	return atc.Or(matchers...)
}

func sniMatcherFromSNIs(snis []string) atc.Matcher {
	matchers := make([]atc.Matcher, 0, len(snis))
	for _, sni := range snis {
		if !validSNIs.MatchString(sni) {
			continue
		}
		matchers = append(matchers, atc.NewPredicateTLSSNI(atc.OpEqual, sni))
	}
	return atc.Or(matchers...)
}
