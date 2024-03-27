package translators

import (
	"encoding/json"
	"fmt"
	pathlib "path"
	"sort"
	"strings"

	"github.com/kong/go-kong/kong"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// KongServiceTranslation is a translation of a single HTTPRoute into metadata
// that can be used to instantiate Kong routes and services.
// Routes from this object should route traffic to BackendRefs from this object.
type KongServiceTranslation struct {
	Name        string
	BackendRefs []gatewayv1.HTTPBackendRef
	KongRoutes  []KongRouteTranslation
}

// KongRouteTranslation is a translation of a single HTTPRoute rule into metadata
// that can be used to instantiate Kong routes.
type KongRouteTranslation struct {
	Name    string
	Matches []gatewayv1.HTTPRouteMatch
	Filters []gatewayv1.HTTPRouteFilter
}

// TranslateHTTPRoute translates a list of HTTPRoutes into a list of HTTPRouteTranslationMeta
// objects that can be used to instantiate Kong routes and services.
// The translation is done by grouping the HTTPRoutes by their backendRefs.
// This means that all the rules of a single HTTPRoute will be grouped together
// if they share the same backendRefs.
func TranslateHTTPRoute(route *gatewayv1.HTTPRoute) []*KongServiceTranslation {
	index := httpRouteTranslationIndex{}
	index.setRoute(route)
	return index.translate()
}

// -----------------------------------------------------------------------------
// HTTPRoute Translation - Private - Index
// -----------------------------------------------------------------------------

// httpRouteTranslationIndex aggregates all rules routing to the same backends group.
type httpRouteTranslationIndex struct {
	httpRoute *gatewayv1.HTTPRoute
	rulesMeta []httpRouteRuleMeta
}

func (i *httpRouteTranslationIndex) setRoute(route *gatewayv1.HTTPRoute) {
	i.httpRoute = route
	i.extractRulesMeta(route)
}

func (i *httpRouteTranslationIndex) extractRulesMeta(route *gatewayv1.HTTPRoute) {
	i.rulesMeta = make([]httpRouteRuleMeta, 0, len(route.Spec.Rules))

	for ruleNumber, rule := range route.Spec.Rules {
		i.rulesMeta = append(i.rulesMeta, httpRouteRuleMeta{
			RuleNumber: ruleNumber,
			Rule:       rule,
		})
	}
}

func (i *httpRouteTranslationIndex) translate() []*KongServiceTranslation {
	rulesGroupedByBackendRed := groupRulesByBackendRefs(i.rulesMeta)
	translations := make([]*KongServiceTranslation, 0, len(rulesGroupedByBackendRed))

	for _, rulesByBackends := range rulesGroupedByBackendRed {
		// each backend refs group is a separate Kong service, not eligible for consolidation
		kongServiceTranslation := i.translateToKongService(rulesByBackends)
		i.translateToKongServiceRoutes(kongServiceTranslation, rulesByBackends)
		translations = append(translations, kongServiceTranslation)
	}

	return translations
}

func (i *httpRouteTranslationIndex) translateToKongService(rulesMeta []httpRouteRuleMeta) *KongServiceTranslation {
	return &KongServiceTranslation{
		Name:        i.translateToKongServiceName(rulesMeta),
		BackendRefs: i.translateToKongServiceBackends(rulesMeta),
		KongRoutes:  nil,
	}
}

func (i *httpRouteTranslationIndex) translateToKongServiceName(rulesMeta []httpRouteRuleMeta) string {
	// this should never happen, as we validate for the number of matches in the parser,
	// but just in case anything changes in the future to avoid panics
	firstRuleInGroup := -1
	if len(rulesMeta) > 0 {
		// rules are guaranteed to retain their order, so we can use the first one
		firstRuleInGroup = rulesMeta[0].RuleNumber
	}
	return fmt.Sprintf(
		"httproute.%s.%s.%d",
		i.httpRoute.Namespace,
		i.httpRoute.Name,
		firstRuleInGroup,
	)
}

func (i *httpRouteTranslationIndex) translateToKongServiceBackends(rulesMeta []httpRouteRuleMeta) []gatewayv1.HTTPBackendRef {
	if len(rulesMeta) == 0 {
		return nil
	}
	// get the backendRefs and filters from any rule, as they are all the same,
	// because the the rules are processed in groups wit the same backendRefs and filters.
	return rulesMeta[0].Rule.BackendRefs
}

func (i *httpRouteTranslationIndex) translateToKongServiceRoutes(s *KongServiceTranslation, rulesMeta []httpRouteRuleMeta) {
	for _, rulesByFilter := range groupRulesByFilter(rulesMeta) {
		// each filter group must be a separate Kong route, not eligible for consolidation
		// thus need to be translated separately
		routesByFilter := i.translateToKongRoutes(rulesByFilter)
		s.KongRoutes = append(s.KongRoutes, routesByFilter...)
	}

	// Sort the routes by name to ensure that the order is deterministic.
	sort.Slice(s.KongRoutes, func(i, j int) bool {
		return s.KongRoutes[i].Name < s.KongRoutes[j].Name
	})
}

func (i *httpRouteTranslationIndex) translateToKongRoutes(rulesMeta []httpRouteRuleMeta) []KongRouteTranslation {
	// All the rules in the group have the same backendRefs and filters.
	var filters []gatewayv1.HTTPRouteFilter
	if len(rulesMeta) > 0 {
		filters = rulesMeta[0].Rule.Filters
	}

	// Group the matches for each rule. Then aggregate the matches eligible for the consolidation
	// into a single match group.
	matchGroups := make(map[string]httpRouteMatchMetaList)
	for _, ruleMeta := range rulesMeta {
		ruleMatchGroups := groupSliceByKeyFn(ruleMeta.matches(), httpRouteMatchMeta.getKey)
		for matchGroupKey, matchGroup := range ruleMatchGroups {
			matchGroups[matchGroupKey] = append(matchGroups[matchGroupKey], matchGroup...)
		}
	}

	// Then, for each group, create a KongRoute with the multiple paths.
	kongRoutes := make([]KongRouteTranslation, 0, len(matchGroups))
	for _, matchGroup := range matchGroups {
		kongRouteName := i.translateToKongRouteName(matchGroup)

		kongRoutes = append(kongRoutes, KongRouteTranslation{
			Name:    kongRouteName,
			Matches: matchGroup.httpRouteMatches(),
			Filters: filters,
		})
	}

	// No matches means a catch-all route based on the hostname
	if len(matchGroups) == 0 {
		kongRouteName := fmt.Sprintf("httproute.%s.%s.0.0", i.httpRoute.Namespace, i.httpRoute.Name)
		kongRoutes = append(kongRoutes, KongRouteTranslation{
			Name:    kongRouteName,
			Filters: filters,
		})
	}

	return kongRoutes
}

func (i *httpRouteTranslationIndex) translateToKongRouteName(matchesMeta httpRouteMatchMetaList) string {
	// this should never happen, as we validate for the number of matches in the parser,
	// but just in case anything changes in the future to avoid panics
	firstRuleNumber := -1
	firstMatchNumber := -1
	if len(matchesMeta) > 0 {
		// matches are guaranteed to retain their order, so we can use the first one
		firstRuleNumber = matchesMeta[0].RuleNumber
		firstMatchNumber = matchesMeta[0].MatchNumber
	}

	return fmt.Sprintf(
		"httproute.%s.%s.%d.%d",
		i.httpRoute.Namespace,
		i.httpRoute.Name,
		firstRuleNumber,
		firstMatchNumber,
	)
}

// groupRulesByBackendRefs groups the rules by their backendRefs.
// The backendRefs are grouped by their key function.
// The elements in the groups have the order of the original slice, but the groups themselves are not ordered.
func groupRulesByBackendRefs(ruleEntries []httpRouteRuleMeta) map[string][]httpRouteRuleMeta {
	return groupSliceByKeyFn(ruleEntries, httpRouteRuleMeta.getHTTPBackendRefsKey)
}

// groupRulesByFilter groups the rules by their filters.
// The filters are grouped by deep equality, key being the numeric index.
// The elements in the groups have the order of the original slice, but the groups themselves are not ordered.
func groupRulesByFilter(ruleEntries []httpRouteRuleMeta) map[string][]httpRouteRuleMeta {
	return groupSliceByKeyFn(ruleEntries, httpRouteRuleMeta.getFiltersKey)
}

// groupSliceByKeyFn groups a slice by a key function. The elements in the groups
// have the order of the original slice, but the groups themselves are not ordered.
func groupSliceByKeyFn[T any](vals []T, keyFn func(val T) string) map[string][]T {
	groups := make(map[string][]T, len(vals))
	for _, val := range vals {
		key := keyFn(val)
		groups[key] = append(groups[key], val)
	}

	return groups
}

// -----------------------------------------------------------------------------
// HttpRoute Translation - Private - Metadata
// -----------------------------------------------------------------------------

type httpRouteRuleMeta struct {
	Rule       gatewayv1.HTTPRouteRule
	RuleNumber int
}

// getFiltersKey computes a key from a list of filters.
// The order of the filters is not important.
func (m httpRouteRuleMeta) getFiltersKey() string {
	return getSortedItemsString(m.Rule.Filters)
}

// getHTTPBackendRefsKey computes a key from a list of backendRefs.
// The order of backedRefs is not important.
func (m httpRouteRuleMeta) getHTTPBackendRefsKey() string {
	return getSortedItemsString(m.Rule.BackendRefs)
}

func (m *httpRouteRuleMeta) matches() httpRouteMatchMetaList {
	matches := make([]httpRouteMatchMeta, 0, len(m.Rule.Matches))

	for matchNumber, match := range m.Rule.Matches {
		match := match
		matches = append(matches, httpRouteMatchMeta{
			Match:       &match,
			RuleNumber:  m.RuleNumber,
			MatchNumber: matchNumber,
		})
	}

	return matches
}

type httpRouteMatchMeta struct {
	Match       *gatewayv1.HTTPRouteMatch
	RuleNumber  int
	MatchNumber int
}

// getKey computes a key from an HTTPRouteMatch. Two HTTPRouteMatches will generate the same key if their
// methods, headers, and query parameters are identical. HTTPRouteMatches with the same key can be
// combined into a single Kong route.
func (m httpRouteMatchMeta) getKey() string {
	// Per the HTTPHeader definition at https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io%2fv1beta1.HTTPHeader
	//
	// Name is the name of the HTTP Header to be matched. Name matching MUST be
	// case insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
	//
	// If multiple entries specify equivalent header names, only the first
	// entry with an equivalent name MUST be considered for a match. Subsequent
	// entries with an equivalent header name MUST be ignored. Due to the
	// case-insensitivity of header names, "foo" and "Foo" are considered
	// equivalent.
	//
	seenHeaders := make(map[string]struct{})
	headers := make([]gatewayv1.HTTPHeaderMatch, 0, len(m.Match.Headers))
	for _, header := range m.Match.Headers {
		name := strings.ToLower(string(header.Name))
		header.Name = gatewayv1.HTTPHeaderName(strings.ToLower(string(header.Name)))
		if _, ok := seenHeaders[name]; ok {
			continue
		}
		seenHeaders[name] = struct{}{}
		headers = append(headers, header)
	}
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].Name < headers[j].Name
	})

	// According to the spec of HTTPQueryParamMatch:
	//
	// If multiple entries specify equivalent query param names, only the first
	// entry with an equivalent name MUST be considered for a match. Subsequent
	// entries with an equivalent query param name MUST be ignored.
	//
	seenQueryParams := make(map[string]struct{})
	queryParams := make([]gatewayv1.HTTPQueryParamMatch, 0, len(m.Match.QueryParams))
	for _, queryParam := range m.Match.QueryParams {
		if _, ok := seenQueryParams[string(queryParam.Name)]; ok {
			continue
		}
		seenQueryParams[string(queryParam.Name)] = struct{}{}
		queryParams = append(queryParams, queryParam)
	}

	keySource := struct {
		Method  *gatewayv1.HTTPMethod
		Headers []gatewayv1.HTTPHeaderMatch
		Query   []gatewayv1.HTTPQueryParamMatch
	}{
		m.Match.Method,
		headers,
		queryParams,
	}

	return mustMarshalJSON(keySource)
}

type httpRouteMatchMetaList []httpRouteMatchMeta

func (l httpRouteMatchMetaList) httpRouteMatches() []gatewayv1.HTTPRouteMatch {
	matches := make([]gatewayv1.HTTPRouteMatch, 0, len(l))
	for _, matchMeta := range l {
		matches = append(matches, *matchMeta.Match)
	}
	return matches
}

// getSortedItemsString returns a string representation of a list of items,
// sorted by their string representation. The items are required to be
// to be JSON marshalable.
func getSortedItemsString[T any](items []T) string {
	keys := make([]string, 0, len(items))

	for _, item := range items {
		keys = append(keys, mustMarshalJSON(item))
	}
	sort.Strings(keys)

	return strings.Join(keys, ";")
}

// mustMarshalJSON marshals the given object to JSON or returns an empty string if it fails.
// It's a developer error if the object cannot be marshaled to JSON.
// The functions using this, should be calling it with a struct that is known to be marshalable.
func mustMarshalJSON[T any](val T) string {
	key, err := json.Marshal(val)
	if err != nil {
		return ""
	}
	return string(key)
}

// GeneratePluginsFromHTTPRouteFilters converts HTTPRouteFilter into Kong plugins.
// path is the parameter to be used by the redirect plugin, to perform redirection.
func GeneratePluginsFromHTTPRouteFilters(filters []gatewayv1.HTTPRouteFilter, path string, tags []*string) []kong.Plugin {
	kongPlugins := make([]kong.Plugin, 0)
	if len(filters) == 0 {
		return kongPlugins
	}

	for _, filter := range filters {
		switch filter.Type {
		case gatewayv1.HTTPRouteFilterRequestHeaderModifier:
			kongPlugins = append(kongPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))

		case gatewayv1.HTTPRouteFilterRequestRedirect:
			kongPlugins = append(kongPlugins, generateRequestRedirectKongPlugin(filter.RequestRedirect, path)...)

		case gatewayv1.HTTPRouteFilterResponseHeaderModifier:
			kongPlugins = append(kongPlugins, generateResponseHeaderModifierKongPlugin(filter.ResponseHeaderModifier))

		case gatewayv1.HTTPRouteFilterExtensionRef,
			gatewayv1.HTTPRouteFilterRequestMirror,
			gatewayv1.HTTPRouteFilterURLRewrite:
			// not supported
		}
	}
	for _, p := range kongPlugins {
		// This plugin is derived from an HTTPRoute filter, not a KongPlugin, so we apply tags indicating that
		// HTTPRoute as the parent Kubernetes resource for these generated plugins.
		p.Tags = tags
	}

	return kongPlugins
}

// generateRequestRedirectKongPlugin generates configurations of plugins to satisfy the specification
// of request redirect filter.
func generateRequestRedirectKongPlugin(modifier *gatewayv1.HTTPRequestRedirectFilter, path string) []kong.Plugin {
	plugins := make([]kong.Plugin, 2)
	plugins[0] = kong.Plugin{
		Name: kong.String("request-termination"),
		Config: kong.Configuration{
			"status_code": modifier.StatusCode,
		},
	}

	var locationHeader string
	scheme := "http"
	port := 80

	if modifier.Scheme != nil {
		scheme = *modifier.Scheme
	}
	if modifier.Port != nil {
		port = int(*modifier.Port)
	}
	if modifier.Path != nil &&
		modifier.Path.Type == gatewayv1.FullPathHTTPPathModifier &&
		modifier.Path.ReplaceFullPath != nil {
		// only ReplaceFullPath currently supported
		path = *modifier.Path.ReplaceFullPath
	}
	if modifier.Hostname != nil {
		locationHeader = fmt.Sprintf("Location: %s://%s", scheme, pathlib.Join(fmt.Sprintf("%s:%d", *modifier.Hostname, port), path))
	} else {
		locationHeader = fmt.Sprintf("Location: %s", path)
	}

	plugins[1] = kong.Plugin{
		Name: kong.String("response-transformer"),
		Config: kong.Configuration{
			"add": map[string][]string{
				"headers": {locationHeader},
			},
		},
	}

	return plugins
}

// generateRequestHeaderModifierKongPlugin converts a gatewayv1.HTTPRequestHeaderFilter into a
// kong.Plugin of type request-transformer.
func generateRequestHeaderModifierKongPlugin(modifier *gatewayv1.HTTPHeaderFilter) kong.Plugin {
	return generateHeaderModifierKongPlugin(modifier, "request-transformer")
}

// generateResponseHeaderModifierKongPlugin converts a gatewayv1.HTTPResponseHeaderFilter into a
// kong.Plugin of type response-transformer.
func generateResponseHeaderModifierKongPlugin(modifier *gatewayv1.HTTPHeaderFilter) kong.Plugin {
	return generateHeaderModifierKongPlugin(modifier, "response-transformer")
}

func generateHeaderModifierKongPlugin(modifier *gatewayv1.HTTPHeaderFilter, pluginName string) kong.Plugin {
	plugin := kong.Plugin{
		Name:   kong.String(pluginName),
		Config: make(kong.Configuration),
	}

	// modifier.Set is converted to a pair composed of "replace" and "add"
	if modifier.Set != nil {
		setModifiers := make([]string, 0, len(modifier.Set))
		for _, s := range modifier.Set {
			setModifiers = append(setModifiers, kongHeaderFormatter(s))
		}
		plugin.Config["replace"] = map[string][]string{
			"headers": setModifiers,
		}
		plugin.Config["add"] = map[string][]string{
			"headers": setModifiers,
		}
	}

	// modifier.Add is converted to "append"
	if modifier.Add != nil {
		appendModifiers := make([]string, 0, len(modifier.Add))
		for _, a := range modifier.Add {
			appendModifiers = append(appendModifiers, kongHeaderFormatter(a))
		}
		plugin.Config["append"] = map[string][]string{
			"headers": appendModifiers,
		}
	}

	if modifier.Remove != nil {
		plugin.Config["remove"] = map[string][]string{
			"headers": modifier.Remove,
		}
	}

	return plugin
}

func kongHeaderFormatter(header gatewayv1.HTTPHeader) string {
	return fmt.Sprintf("%s:%s", header.Name, header.Value)
}
