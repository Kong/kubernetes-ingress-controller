package subtranslator

import (
	"encoding/json"
	"fmt"
	pathlib "path"
	"sort"
	"strings"

	"dario.cat/mergo"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/samber/lo"
)

// KongServiceTranslation is a translation of a single HTTPRoute into metadata
// that can be used to instantiate Kong routes and services.
// Routes from this object should route traffic to BackendRefs from this object.
type KongServiceTranslation struct {
	Name        string
	BackendRefs []gatewayapi.HTTPBackendRef
	KongRoutes  []KongRouteTranslation
}

// KongRouteTranslation is a translation of a single HTTPRoute rule into metadata
// that can be used to instantiate Kong routes.
type KongRouteTranslation struct {
	Name    string
	Matches []gatewayapi.HTTPRouteMatch
	Filters []gatewayapi.HTTPRouteFilter
}

// TranslateHTTPRoute translates a list of HTTPRoutes into a list of HTTPRouteTranslationMeta
// objects that can be used to instantiate Kong routes and services.
// The translation is done by grouping the HTTPRoutes by their backendRefs.
// This means that all the rules of a single HTTPRoute will be grouped together
// if they share the same backendRefs.
func TranslateHTTPRoute(route *gatewayapi.HTTPRoute) []*KongServiceTranslation {
	index := httpRouteTranslationIndex{}
	index.setRoute(route)
	return index.translate()
}

// -----------------------------------------------------------------------------
// HTTPRoute Translation - Private - Index
// -----------------------------------------------------------------------------

// httpRouteTranslationIndex aggregates all rules routing to the same backends group.
type httpRouteTranslationIndex struct {
	httpRoute *gatewayapi.HTTPRoute
	rulesMeta []httpRouteRuleMeta
}

func (i *httpRouteTranslationIndex) setRoute(route *gatewayapi.HTTPRoute) {
	i.httpRoute = route
	i.extractRulesMeta(route)
}

func (i *httpRouteTranslationIndex) extractRulesMeta(route *gatewayapi.HTTPRoute) {
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
	// This should never happen, as we validate for the number of matches in the translator,
	// but just in case anything changes in the future to avoid panics.
	firstRuleInGroup := -1
	if len(rulesMeta) > 0 {
		// Rules are guaranteed to retain their order, so we can use the first one.
		firstRuleInGroup = rulesMeta[0].RuleNumber
	}
	return fmt.Sprintf(
		"httproute.%s.%s.%d",
		i.httpRoute.Namespace,
		i.httpRoute.Name,
		firstRuleInGroup,
	)
}

func (i *httpRouteTranslationIndex) translateToKongServiceBackends(rulesMeta []httpRouteRuleMeta) []gatewayapi.HTTPBackendRef {
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
	var filters []gatewayapi.HTTPRouteFilter
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
	// This should never happen, as we validate for the number of matches in the translator,
	// but just in case anything changes in the future to avoid panics.
	firstRuleNumber := -1
	firstMatchNumber := -1
	if len(matchesMeta) > 0 {
		// Matches are guaranteed to retain their order, so we can use the first one.
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
	Rule       gatewayapi.HTTPRouteRule
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
	Match       *gatewayapi.HTTPRouteMatch
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
	headers := make([]gatewayapi.HTTPHeaderMatch, 0, len(m.Match.Headers))
	for _, header := range m.Match.Headers {
		name := strings.ToLower(string(header.Name))
		header.Name = gatewayapi.HTTPHeaderName(strings.ToLower(string(header.Name)))
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
	queryParams := make([]gatewayapi.HTTPQueryParamMatch, 0, len(m.Match.QueryParams))
	for _, queryParam := range m.Match.QueryParams {
		if _, ok := seenQueryParams[string(queryParam.Name)]; ok {
			continue
		}
		seenQueryParams[string(queryParam.Name)] = struct{}{}
		queryParams = append(queryParams, queryParam)
	}

	keySource := struct {
		Method  *gatewayapi.HTTPMethod
		Headers []gatewayapi.HTTPHeaderMatch
		Query   []gatewayapi.HTTPQueryParamMatch
	}{
		m.Match.Method,
		headers,
		queryParams,
	}

	return mustMarshalJSON(keySource)
}

type httpRouteMatchMetaList []httpRouteMatchMeta

func (l httpRouteMatchMetaList) httpRouteMatches() []gatewayapi.HTTPRouteMatch {
	matches := make([]gatewayapi.HTTPRouteMatch, 0, len(l))
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

type SetRoutePluginsArgs struct {
	Route   *kongstate.Route
	Filters []gatewayapi.HTTPRouteFilter
}

// SetRoutePlugins converts HTTPRouteFilter into Kong plugins. The plugins are set into the given kongstate.Route.
// The plugins can be set in two different ways:
// - Direct conversion from the respective HTTPRouteFilter.
// - ExtensionRef to plugins annotation from the ExtensionRef filter.
func SetRoutePlugins(route *kongstate.Route, filters []gatewayapi.HTTPRouteFilter, path string, tags []*string) error {
	generatedPlugins, err := generatePluginsFromHTTPRouteFilters(filters, path, tags)
	if err != nil {
		return err
	}
	route.Plugins = append(route.Plugins, generatedPlugins.Plugins...)
	for _, modifier := range generatedPlugins.KongRouteModifiers {
		modifier(route)
	}
	return nil
}

type httpRouteFiltersOriginatedPlugins struct {
	// Plugins is a list of Kong plugins generated from HTTPRoute filters to be applied to a Kong Route.
	Plugins []kong.Plugin

	// KongRouteModifiers is a list of functions that will be used to modify a Kong Route.
	// These can be used to set additional properties on the Kong Route or changing its properties if required by
	// the filters (e.g. setting different paths).
	KongRouteModifiers []kongRouteModifier

	// PluginsAnnotation is a `konghq.com/plugins: <plugins>` annotation value generated from the ExtensionRef filter.
	PluginsAnnotation string
}

// kongRouteModifier is a function that modifies a Kong route.
type kongRouteModifier func(*kongstate.Route)

// generatePluginsFromHTTPRouteFilters converts HTTPRouteFilter into Kong plugins.
// path is the parameter to be used by the redirect plugin, to perform redirection.
// It returns two values:
// - A set of plugins generated by the conversion of all the provided filters, excluding  ExtensionRefs.
// - A plugins annotation value, generated by the ExtensionRef filter.
func generatePluginsFromHTTPRouteFilters(filters []gatewayapi.HTTPRouteFilter, path string, tags []*string) (httpRouteFiltersOriginatedPlugins, error) {
	if len(filters) == 0 {
		return httpRouteFiltersOriginatedPlugins{}, nil
	}
	var (
		kongPlugins                 []kong.Plugin
		pluginNamesFromExtensionRef []string
		kongRouteModifiers          []kongRouteModifier
	)

	for _, filter := range filters {
		switch filter.Type {
		case gatewayapi.HTTPRouteFilterRequestHeaderModifier:
			kongPlugins = append(kongPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))

		case gatewayapi.HTTPRouteFilterRequestRedirect:
			kongPlugins = append(kongPlugins, generateRequestRedirectKongPlugin(filter.RequestRedirect, path)...)

		case gatewayapi.HTTPRouteFilterResponseHeaderModifier:
			kongPlugins = append(kongPlugins, generateResponseHeaderModifierKongPlugin(filter.ResponseHeaderModifier))

		case gatewayapi.HTTPRouteFilterExtensionRef:
			plugin, err := generateExtensionRefKongPlugin(filter.ExtensionRef)
			if err != nil {
				return httpRouteFiltersOriginatedPlugins{}, err
			}
			pluginNamesFromExtensionRef = append(pluginNamesFromExtensionRef, plugin)

		case gatewayapi.HTTPRouteFilterURLRewrite:
			plugin, routeModifier, err := generateRequestTransformerForURLRewrite(filter.URLRewrite, path)
			if err != nil {
				return httpRouteFiltersOriginatedPlugins{}, err
			}
			kongPlugins = append(kongPlugins, plugin)
			if routeModifier != nil {
				kongRouteModifiers = append(kongRouteModifiers, routeModifier)
			}

		case gatewayapi.HTTPRouteFilterRequestMirror:
			// not supported
			return httpRouteFiltersOriginatedPlugins{}, fmt.Errorf("httpFilter %s unsupported", filter.Type)
		}
	}
	for _, p := range kongPlugins {
		// This plugin is derived from an HTTPRoute filter, not a KongPlugin, so we apply tags indicating that
		// HTTPRoute as the parent Kubernetes resource for these generated plugins.
		p.Tags = tags
	}

	// It's possible the above loop generates multiple plugins of the same type, so we need to merge them.
	// It can happen for example when both RequestHeaderModifier and HTTPRouteFilterURLRewrite are present.
	kongPlugins, err := mergePluginsOfTheSameType(kongPlugins)
	if err != nil {
		return httpRouteFiltersOriginatedPlugins{}, fmt.Errorf("failed to merge plugins of the same type: %w", err)
	}

	if len(pluginNamesFromExtensionRef) > 0 {
		kongRouteModifiers = append(kongRouteModifiers, generateKongRouteModifierFromExtensionRef(pluginNamesFromExtensionRef))
	}

	return httpRouteFiltersOriginatedPlugins{
		Plugins:            kongPlugins,
		KongRouteModifiers: kongRouteModifiers,
	}, nil
}

func mergePluginsOfTheSameType(plugins []kong.Plugin) ([]kong.Plugin, error) {
	pluginsByName := lo.GroupBy(plugins, func(p kong.Plugin) string {
		return *p.Name // Name is effectively a plugin type.
	})
	for pluginName, plugins := range pluginsByName {
		// If we produced multiple plugins of the same type, we need to merge their configurations now.
		if len(plugins) > 1 {
			mergedPlugin := *plugins[0].DeepCopy()
			for _, plugin := range plugins {
				if err := mergo.Merge(&mergedPlugin.Config, plugin.Config); err != nil {
					// Should never happen as we're passing the same type of objects.
					return nil, fmt.Errorf("failed to merge %q plugin configurations: %w", pluginName, err)
				}
			}
			pluginsByName[pluginName] = []kong.Plugin{mergedPlugin}
		}
	}
	return lo.Flatten(lo.Values(pluginsByName)), nil
}

func generateKongRouteModifierFromExtensionRef(pluginNamesFromExtensionRef []string) kongRouteModifier {
	return func(route *kongstate.Route) {
		if route.Ingress.Annotations == nil {
			route.Ingress.Annotations = make(map[string]string)
		}
		annotationValue := strings.Join(pluginNamesFromExtensionRef, ",")
		const pluginAnnotationKey = annotations.AnnotationPrefix + annotations.PluginsKey
		if _, ok := route.Ingress.Annotations[pluginAnnotationKey]; !ok {
			route.Ingress.Annotations[pluginAnnotationKey] = annotationValue
		} else {
			route.Ingress.Annotations[pluginAnnotationKey] = fmt.Sprintf("%s,%s",
				route.Ingress.Annotations[pluginAnnotationKey],
				annotationValue)
		}
	}
}

// generateRequestRedirectKongPlugin generates configurations of plugins to satisfy the specification
// of request redirect filter.
func generateRequestRedirectKongPlugin(modifier *gatewayapi.HTTPRequestRedirectFilter, path string) []kong.Plugin {
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
		modifier.Path.Type == gatewayapi.FullPathHTTPPathModifier &&
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

func generateExtensionRefKongPlugin(modifier *gatewayapi.LocalObjectReference) (string, error) {
	if modifier.Group != "configuration.konghq.com" || modifier.Kind != "KongPlugin" {
		return "", fmt.Errorf("plugin %s/%s unsupported", modifier.Group, modifier.Kind)
	}
	return string(modifier.Name), nil
}

// generateRequestHeaderModifierKongPlugin converts a gatewayapi.HTTPRequestHeaderFilter into a
// kong.Plugin of type request-transformer.
func generateRequestHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter) kong.Plugin {
	return generateHeaderModifierKongPlugin(modifier, "request-transformer")
}

// generateResponseHeaderModifierKongPlugin converts a gatewayapi.HTTPResponseHeaderFilter into a
// kong.Plugin of type response-transformer.
func generateResponseHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter) kong.Plugin {
	return generateHeaderModifierKongPlugin(modifier, "response-transformer")
}

func generateHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter, pluginName string) kong.Plugin {
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

func kongHeaderFormatter(header gatewayapi.HTTPHeader) string {
	return fmt.Sprintf("%s:%s", header.Name, header.Value)
}

func generateRequestTransformerForURLRewrite(filter *gatewayapi.HTTPURLRewriteFilter, path string) (kong.Plugin, kongRouteModifier, error) {
	if filter == nil {
		return kong.Plugin{}, nil, fmt.Errorf("%s is not provided", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	if filter.Path == nil && filter.Hostname == nil {
		return kong.Plugin{}, nil, fmt.Errorf("%s missing Path and Hostname", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	if filter.Path != nil {
		switch filter.Path.Type {
		case gatewayapi.FullPathHTTPPathModifier:
			plugin, err := generateRequestTransformerForURLRewriteFullPath(filter)
			if err != nil {
				return kong.Plugin{}, nil, fmt.Errorf("failed to generate request-transformer plugin for %s: %w", gatewayapi.HTTPRouteFilterURLRewrite, err)
			}
			return plugin, nil, nil

		case gatewayapi.PrefixMatchHTTPPathModifier:
			plugin, routeModifier := generateRequestTransformerForURLRewritePrefixMatch(filter, path)
			return plugin, routeModifier, nil
		}
	}

	// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3685
	if filter.Hostname != nil {
		return kong.Plugin{}, nil, fmt.Errorf("unsupported hostname replace for %s", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	return kong.Plugin{}, nil, fmt.Errorf("invalid %s config", gatewayapi.HTTPRouteFilterURLRewrite)
}

// generateRequestTransformerForURLRewriteFullPath generates a request-transformer plugin for the URLRewrite filter
// with a FullPathHTTPPathModifier.
func generateRequestTransformerForURLRewriteFullPath(filter *gatewayapi.HTTPURLRewriteFilter) (kong.Plugin, error) {
	if filter.Path.ReplaceFullPath == nil {
		return kong.Plugin{}, fmt.Errorf("%s missing ReplaceFullPath", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	return kong.Plugin{
		Name: kong.String("request-transformer"),
		Config: kong.Configuration{
			"replace": map[string]string{
				"uri": *filter.Path.ReplaceFullPath,
			},
		},
	}, nil
}

// generateRequestTransformerForURLRewritePrefixMatch generates a request-transformer plugin and a route modifier
// for the URLRewrite filter with a PrefixMatchHTTPPathModifier.
func generateRequestTransformerForURLRewritePrefixMatch(filter *gatewayapi.HTTPURLRewriteFilter, path string) (kong.Plugin, kongRouteModifier) {
	// By default, we replace the URI with a slash.
	replaceURI := "/"

	// If the ReplacePrefixMatch is provided, we need to use it as the replacement URI.
	if replacePrefixMatch := filter.Path.ReplacePrefixMatch; replacePrefixMatch != nil {
		// Trim the trailing slash from the ReplacePrefixMatch to avoid double slashes in the final URI.
		replaceURI = strings.TrimSuffix(*replacePrefixMatch, "/")
		if replaceURI == "" {
			// In the case of an empty ReplacePrefixMatch, we need to make sure that the path will always start with a slash,
			// even if we have no capture group captured from the incoming request's URI.
			// The below is a Lua ternary operator that checks if the captured group is nil, and if so, replaces it with a slash.
			replaceURI = `$(uri_captures[1] == nil and "/" or uri_captures[1])`
		} else {
			// Otherwise, we concatenate the replacement URI with the captured group.
			replaceURI = fmt.Sprintf(`%s$(uri_captures[1])`, replaceURI)
		}
	}
	plugin := kong.Plugin{
		Name: kong.String("request-transformer"),
		Config: kong.Configuration{
			"replace": map[string]string{
				"uri": replaceURI,
			},
		},
	}

	// We need to change the paths in the Kong route to a regex with a capture group that can be used
	// in the request-transformer plugin.
	routeModifier := func(route *kongstate.Route) {
		paths := make([]*string, 0, 2)
		// The first path matches the exact path.
		paths = append(paths, lo.ToPtr(fmt.Sprintf("%s%s$", KongPathRegexPrefix, path)))
		// The second path matches subpaths, including a single slash.
		paths = append(paths, lo.ToPtr(
			fmt.Sprintf("%s%s(/.*)", KongPathRegexPrefix, strings.TrimSuffix(path, "/"))),
		)
		route.Paths = paths
	}

	return plugin, routeModifier
}
