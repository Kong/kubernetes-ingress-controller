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

func (m *ingressTranslationMeta) translateIntoKongExpressionRoute() *kongstate.Route {
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
		hostMatcher := hostMatcherFromHosts(hosts)
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

	priority := calculateExpressionRoutePriority(m.paths, pathRegexPrefix, m.ingressHost, ingressAnnotations)
	atc.ApplyExpression(&route.Route, routeMatcher, priority)
	return route
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

// calculateExpressionRoutePriority calculates the priority of the generated route.
// It basically follows the calculating method used in Kong's traditional compatible router:
//   - First, sort by number of fields specified (methods, hosts, headers, paths, snis)
//   - Then, if both routes have hosts, the routes having only plain(non-wildcard) hosts has
//     higher priority than routes having at least one wildcard host
//   - Then, sort by numer of different headers (maximum header count = 255).
//   - Then, paths with regex match has higher priority than prefix match.
//   - Then, sort by regex_priority field (not supported in KIC with expression routes).
//   - At last, sort by maximum length of paths in the route.
func calculateExpressionRoutePriority(
	paths []netv1.HTTPIngressPath,
	regexPathPrefix string,
	ingressHost string,
	ingressAnnotations map[string]string,
) int {
	matchFields := 0
	plainHost := false
	headerCount := 0
	maxPathLength := 0
	hasRegexPath := false

	// add 1 to matchFields if path is non-empty.
	if len(paths) > 0 {
		matchFields++
	}
	for _, path := range paths {
		// TODO: in non-combined translator, we assigned different regex priorities to different type of path:
		// exact = 300, prefix = 200, implementationSpecific = 100.
		// We may add additional weights on path types.
		if path.PathType == nil {
			continue
		}
		// since we will generate regex path for exact and prefix matches in traditional routes,
		// we should consider them as using regex path.
		var pathLength int
		switch *path.PathType {
		case netv1.PathTypeExact:
			hasRegexPath = true
			// trim all leading '/'s and calculate the length by '/'+ trimmed path.
			// for example: len("/foo") = 4; len("/foo/") = 5; len("//foo") = 4.
			relative := strings.TrimLeft(path.Path, "/")
			pathLength = 1 + len(relative)
		case netv1.PathTypePrefix:
			hasRegexPath = true
			// trim all leading and trailing '/' s and calculate the length by '/' + trimmed path + '/' (except that len("/")=1)
			// for example: len("/foo") = 5; len("/foo/") = 5; len("//foo/bar") = len("/foo/bar//") = 9.
			base := strings.Trim(path.Path, "/")
			if base == "" {
				pathLength = 1
			} else {
				pathLength = 2 + len(base)
			}
		case netv1.PathTypeImplementationSpecific:
			if regexPathPrefix != "" && strings.HasPrefix(path.Path, regexPathPrefix) {
				// regex path match, calculate the length by length of regex part.
				hasRegexPath = true
				regex := strings.TrimPrefix(path.Path, path.Path)
				pathLength = len(regex)
			} else {
				// non-regex path.
				pathLength = len(path.Path)
			}
		}
		if pathLength > maxPathLength {
			maxPathLength = pathLength
		}
	}

	hosts := []string{}
	if ingressHost != "" {
		hosts = append(hosts, ingressHost)
	}
	hostAliases, _ := annotations.ExtractHostAliases(ingressAnnotations)
	hosts = append(hosts, hostAliases...)
	if len(hosts) > 0 {
		matchFields++
		plainHost = true
	}

	for _, host := range hosts {
		if strings.HasPrefix(host, "*") {
			plainHost = false
			break
		}
	}

	headers, exist := annotations.ExtractHeaders(ingressAnnotations)
	if exist {
		matchFields++
		headerCount = len(headers)
	}

	methods := annotations.ExtractMethods(ingressAnnotations)
	if len(methods) > 0 {
		matchFields++
	}

	snis, exist := annotations.ExtractSNIs(ingressAnnotations)
	if exist && len(snis) > 0 {
		matchFields++
	}

	// combine priority into numbers. route.priority in admin API could only use the lowest 52 bits
	// because the numbers in JSON is parsed into double precision floating number by Kong.
	const (
		// lowest 16 bits (0~15) for max path length
		// regexPathShiftBits uses the 16th bit for marking if regex match on path exists.
		regexPathShiftBits = 16
		// bits 17~31 are preserved.

		// plainHostShiftBits uses the 32nd bit for marking if ALL hosts are non-wildcard.
		plainHostShiftBits = 32
		// headerNumberShiftBits makes bits 33~40 used for number of headers.
		headerNumberShiftBits = 33
		// matchFieldsShiftBits uses bits 41 and over (41~43 since there are at most 5 fields).
		matchFieldsShiftBits = 41

		headerNumberLimit = 255
		pathLengthLimit   = (1 << 16) - 1
	)

	var priority int
	// add max path length.
	if maxPathLength > pathLengthLimit {
		maxPathLength = pathLengthLimit
	}
	priority += maxPathLength
	// add regex path mark.
	if hasRegexPath {
		priority += (1 << regexPathShiftBits)
	}
	// add plain host mark.
	if plainHost {
		priority += (1 << plainHostShiftBits)
	}
	if headerCount > headerNumberLimit {
		headerCount = headerNumberLimit
	}
	priority += (headerCount << headerNumberShiftBits)
	priority += (matchFields << matchFieldsShiftBits)

	return priority
}
