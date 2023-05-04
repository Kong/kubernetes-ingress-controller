package translators

import (
	"fmt"
	"regexp"
	"sort"
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
	hosts = append(hosts, hostAliases...)
	if len(hosts) > 0 {
		hostMatcher := hostMatcherFromIngressHosts(hosts)
		routeMatcher.And(hostMatcher)
	}

	// translate paths.
	pathMatchers := make([]atc.Matcher, 0, len(m.paths))
	pathRegexPrefix := annotations.ExtractRegexPrefix(ingressAnnotations)
	if pathRegexPrefix == "" {
		pathRegexPrefix = ControllerPathRegexPrefix
	}
	for _, path := range m.paths {
		pathMatchers = append(pathMatchers, pathMatcherFromIngressPath(path, pathRegexPrefix))
	}
	routeMatcher.And(atc.Or(pathMatchers...))

	// translate protocols.
	protocols := []string{"http", "https"}
	annonationProtocols := annotations.ExtractProtocolNames(ingressAnnotations)
	if len(annonationProtocols) > 0 {
		protocols = annonationProtocols
	}
	protocolMatcher := protocolMatcherFromProtocols(protocols)
	routeMatcher.And(protocolMatcher)

	// translate headers.
	headers, exist := annotations.ExtractHeaders(ingressAnnotations)
	if len(headers) > 0 && exist {
		headerMatcher := headerMatcherFromHeaders(headers)
		routeMatcher.And(headerMatcher)
	}

	// translate methods.
	methods := annotations.ExtractMethods(ingressAnnotations)
	if len(methods) > 0 {
		methodMatcher := methodMatcherFromMethods(methods)
		routeMatcher.And(methodMatcher)
	}

	// translate SNIs.
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
			// wildcard match on hosts (like *.foo.com), genreate a suffix match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpSuffixMatch, strings.TrimPrefix(host, "*")))
		} else {
			// exact match on hosts, generate an exact match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpEqual, host))
		}
	}
	return atc.Or(matchers...)
}

// pathMatcherFromIngressPath translate ingress path into matcher to match the path.
func pathMatcherFromIngressPath(httpIngressPath netv1.HTTPIngressPath, regexPathPrefix string) atc.Matcher {
	switch *httpIngressPath.PathType {
	// Prefix paths.
	case netv1.PathTypePrefix:
		base := strings.Trim(httpIngressPath.Path, "/")
		if base == "" {
			// empty string in prefix path matches prefix "/".
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		return atc.Or(
			// otherwise, match /<path>/* or /<path>.
			atc.NewPredicateHTTPPath(atc.OpEqual, "/"+base),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/"+base+"/"),
		)
	// Exact paths.
	case netv1.PathTypeExact:
		relative := strings.TrimLeft(httpIngressPath.Path, "/")
		return atc.NewPredicateHTTPPath(atc.OpEqual, "/"+relative)
	// Implementation Specific match. treat it as regex match if it begins with a regex prefix (/~ by default),
	// otherwise generate a prefix match.
	case netv1.PathTypeImplementationSpecific:
		// empty path. matches prefix "/" to match any path.
		if httpIngressPath.Path == "" {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		// regex match.
		if regexPathPrefix != "" && strings.HasPrefix(httpIngressPath.Path, regexPathPrefix) {
			regex := strings.TrimPrefix(httpIngressPath.Path, regexPathPrefix)
			if regex == "" {
				regex = "^/"
			}
			// regex match matches a prefix of the whole path, so we need to add a line start annotation in the regex.
			if !strings.HasPrefix(regex, "^") {
				regex = "^" + regex
			}
			return atc.NewPredicateHTTPPath(atc.OpRegexMatch, regex)
		}
		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, httpIngressPath.Path)
	}

	return nil
}

// protocolMatcherFromProtocols gernerates matchers from protocols.
func protocolMatcherFromProtocols(protocols []string) atc.Matcher {
	matchers := []atc.Matcher{}
	for _, protocol := range protocols {
		if !util.ValidateProtocol(protocol) {
			continue
		}
		matchers = append(matchers, atc.NewPredicateNetProtocol(atc.OpEqual, protocol))
	}
	return atc.Or(matchers...)
}

// headerMatcherFromHeaders generates matcher to match headers in HTTP requests.
func headerMatcherFromHeaders(headers map[string][]string) atc.Matcher {
	matchers := make([]atc.Matcher, 0, len(headers))

	// To make a stable result from the same annotations, sort the header names first.
	headerNames := make([]string, 0, len(headers))
	for headerName := range headers {
		headerNames = append(headerNames, headerName)
	}
	sort.Strings(headerNames)

	for _, headerName := range headerNames {
		// extract values.
		values := headers[headerName]
		if len(values) == 0 {
			continue
		}

		// transfer header name to lowercase and replace "-" with "_", which is used in expressions of kong routes.
		headerName = strings.ReplaceAll(strings.ToLower(headerName), "-", "_")
		// header "Host" should be skipped, they are processed in "http.host".
		if headerName == "host" {
			continue
		}

		// values for the same headers ar "or"ed to match any of the values.
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
	// matchers from different headers are "and"ed to match al rules to match headers.
	return atc.And(matchers...)
}

// methodMatcherFromMethods generates matcher to match http methods.
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

// sniMatcherFromSNIs generates matchers to match TLS SNIs.
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
