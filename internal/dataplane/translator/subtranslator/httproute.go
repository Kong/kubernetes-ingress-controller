package subtranslator

import (
	"encoding/json"
	"fmt"
	pathlib "path"
	"sort"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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

// TranslateHTTPRoutesToKongstateServices translates a set of HTTPRoutes to kongstate services,
// and collect the translation errors in the process of translating.
func TranslateHTTPRoutesToKongstateServices(
	logger logr.Logger,
	storer store.Storer,
	routes []*gatewayapi.HTTPRoute,
	combinedServicesFromDifferentHTTPRoutes bool,
) (map[string]kongstate.Service, map[k8stypes.NamespacedName][]error) {
	serviceNameToRules := groupRulesFromHTTPRoutesByKongServiceName(routes, combinedServicesFromDifferentHTTPRoutes)

	kongstateServiceCache := map[string]kongstate.Service{}
	routeTranslationErrors := map[k8stypes.NamespacedName][]error{}
	for serviceName, rulesMeta := range serviceNameToRules {
		if len(rulesMeta) == 0 {
			continue
		}
		service, err := translateHTTPRouteRulesMetaToKongstateService(logger, storer, serviceName, rulesMeta)
		if err != nil {
			httpRoutes := extractUniqueHTTPRoutes(rulesMeta)
			for _, route := range httpRoutes {
				nn := k8stypes.NamespacedName{
					Namespace: route.Namespace,
					Name:      route.Name,
				}
				routeTranslationErrors[nn] = append(routeTranslationErrors[nn], err)
			}
			continue
		}
		kongstateServiceCache[serviceName] = service
	}

	return kongstateServiceCache, routeTranslationErrors
}

// groupRulesFromHTTPRoutesByKongServiceName groups rules from HTTPRoutes into groups of rules, where each rule group
// is translated to a Kong service, and assign translated Kong state service name to the rule groups.
// If combinedServicesFromDifferentHTTPRoutes is true, it groups rules from different HTTPRoutes in the same namespace.
// Otherwise, it groups rules in the same HTTPRoute.
func groupRulesFromHTTPRoutesByKongServiceName(
	routes []*gatewayapi.HTTPRoute,
	combinedServicesFromDifferentHTTPRoutes bool,
) map[string][]httpRouteRuleMeta {
	// If we enabled combined service from different HTTPRoutes, we group rules sharing the same backendRefs
	// from all HTTPRoutes in the same namespace, and generate the Kong service name by backends.
	if combinedServicesFromDifferentHTTPRoutes {
		backendCache := newHTTPRouteBackendCache()
		for _, route := range routes {
			for ruleNumber, rule := range route.Spec.Rules {
				ruleMeta := httpRouteRuleMeta{
					Rule:        rule,
					RuleNumber:  ruleNumber,
					parentRoute: route,
				}
				backendCache.addRule(ruleMeta)
			}
		}
		return backendCache.backends
	}

	// Otherwise, we still group rules in the same HTTPRoute sharing the same backends,
	// and generate the service name by the namespace and name of HTTPRoute and index of the first rule.
	serviceNameToRules := map[string][]httpRouteRuleMeta{}
	for _, route := range routes {
		rulesMetaFromRoute := make([]httpRouteRuleMeta, 0, len(route.Spec.Rules))
		for ruleNumber, rule := range route.Spec.Rules {
			ruleMeta := httpRouteRuleMeta{
				Rule:        rule,
				RuleNumber:  ruleNumber,
				parentRoute: route,
			}
			rulesMetaFromRoute = append(rulesMetaFromRoute, ruleMeta)
		}
		// Since the grouping keeps the order of rules in each group, the first rule in each group should be the rule
		// with the smallest index in the original HTTPRoute in the same group.
		for _, rulesMeta := range groupRulesByBackendRefs(rulesMetaFromRoute) {
			if len(rulesMeta) == 0 {
				continue
			}
			serviceName := fmt.Sprintf("httproute.%s.%s.%d", route.Namespace, route.Name, rulesMeta[0].RuleNumber)
			serviceNameToRules[serviceName] = rulesMeta
		}
	}
	return serviceNameToRules
}

func httpBackendRefsToBackendRefs(httpBackendRef []gatewayapi.HTTPBackendRef, parentRoute *gatewayapi.HTTPRoute) []gatewayapi.BackendRef {
	backendRefs := make([]gatewayapi.BackendRef, 0, len(httpBackendRef))

	for _, hRef := range httpBackendRef {
		backendRef := hRef.BackendRef
		if backendRef.BackendObjectReference.Group == nil {
			backendRef.BackendObjectReference.Group = lo.ToPtr(gatewayapi.Group(""))
		}
		if backendRef.BackendObjectReference.Namespace == nil {
			backendRef.BackendObjectReference.Namespace = lo.ToPtr(gatewayapi.Namespace(parentRoute.Namespace))
		}
		backendRefs = append(backendRefs, backendRef)
	}
	return backendRefs
}

// translateHTTPRouteRulesMetaToKongstateService translates a set of rules sharing the same backends to a `kongstate.Service`.
// The rules should be grouped before calling and the service name should be specified.
func translateHTTPRouteRulesMetaToKongstateService(
	logger logr.Logger,
	storer store.Storer,
	serviceName string,
	rulesMeta []httpRouteRuleMeta,
) (kongstate.Service, error) {
	// Fill in the common fields of the kongstate.Service.
	service := kongstate.Service{
		Service: kong.Service{
			Name:           kong.String(serviceName),
			Host:           kong.String(serviceName),
			Protocol:       kong.String(DefualtKongServiceProtocol),
			ConnectTimeout: kong.Int(DefaultServiceTimeout),
			ReadTimeout:    kong.Int(DefaultServiceTimeout),
			WriteTimeout:   kong.Int(DefaultServiceTimeout),
			Retries:        kong.Int(DefaultRetries),
		},
	}

	// Extract reference grants for checking if reference of backends in another namespace is allowed.
	// Since all RuleMeta grouped here are from HTTPRoutes in the same namespace, it is OK to use the parent HTTPRoute of the first rule as the representative.
	firstHTTPRoute := rulesMeta[0].parentRoute
	grants, err := storer.ListReferenceGrants()
	if err != nil {
		return kongstate.Service{}, fmt.Errorf("could not retrieve ReferenceGrants: %w", err)
	}

	grantFrom := gatewayapi.ReferenceGrantFrom{
		Group:     gatewayapi.Group(firstHTTPRoute.GetObjectKind().GroupVersionKind().Group),
		Kind:      gatewayapi.Kind(firstHTTPRoute.GetObjectKind().GroupVersionKind().Kind),
		Namespace: gatewayapi.Namespace(firstHTTPRoute.GetNamespace()),
	}
	allowed := gatewayapi.GetPermittedForReferenceGrantFrom(
		logger,
		grantFrom,
		grants,
	)

	// generate backends from the backendRefs. The rules should share the same backend ref despite that we do not support filters per backend yet.
	backendRefs := httpBackendRefsToBackendRefs(rulesMeta[0].Rule.BackendRefs, firstHTTPRoute)
	kongstateBackends := backendRefsToKongStateBackends(
		logger,
		storer,
		firstHTTPRoute,
		backendRefs,
		allowed,
	)
	service.Backends = kongstateBackends
	// set the metadata of the kong state service to point to the first HTTPRoute.
	service.Namespace = firstHTTPRoute.GetNamespace()
	service.Parent = firstHTTPRoute

	// In the context of the gateway API conformance tests, if there is no service for the backend,
	// the response must have a status code of 500. Since The default behavior of Kong is returning 503
	// if there is no backend for a service, we inject a plugin that terminates all requests with 500
	// as status code.
	if len(service.Backends) == 0 && len(backendRefs) != 0 {
		if service.Plugins == nil {
			service.Plugins = make([]kong.Plugin, 0)
		}
		service.Plugins = append(service.Plugins, kong.Plugin{
			Name: kong.String("request-termination"),
			Config: kong.Configuration{
				"status_code": 500,
				"message":     "no existing backendRef provided",
			},
		})
	}

	// applyTimeoutToServiceFromHTTPRouteRule applies timeouts from HTTPRoute to the service.
	// If the service is translated from multiple rules, the timeout from the last rule which has specific timeout will be applied to the service in the loop.
	for _, ruleMeta := range rulesMeta {
		applyTimeoutToServiceFromHTTPRouteRule(&service, ruleMeta.Rule)
	}

	routes, err := translateHTTPRouteRulesMetaToKongstateRoutes(rulesMeta)
	if err != nil {
		return kongstate.Service{}, err
	}
	// preallocate slice for the routes.
	service.Routes = make([]kongstate.Route, 0, len(routes))
	service.Routes = append(service.Routes, routes...)

	return service, nil
}

// applyTimeoutToServiceFromHTTPRouteRule applies timeout on the translated Kong service from the timeout settings in the rule.
func applyTimeoutToServiceFromHTTPRouteRule(svc *kongstate.Service, rule gatewayapi.HTTPRouteRule) {
	if rule.Timeouts == nil {
		return
	}
	backendRequestTimeout := DefaultServiceTimeout
	if rule.Timeouts != nil && rule.Timeouts.BackendRequest != nil {
		duration, err := time.ParseDuration(string(*rule.Timeouts.BackendRequest))
		// We ignore the error here because the rule.Timeouts.BackendRequest is validated
		// to be a strict subset of Golang time.ParseDuration so it should never happen
		if err != nil {
			return
		}
		backendRequestTimeout = int(duration.Milliseconds())
	}
	// if the backendRequestTimeout is the same as the default timeout, we don't need to apply it to the service.
	if backendRequestTimeout == DefaultServiceTimeout {
		return
	}
	svc.Service.ReadTimeout = kong.Int(backendRequestTimeout)
	svc.Service.ConnectTimeout = kong.Int(backendRequestTimeout)
	svc.Service.WriteTimeout = kong.Int(backendRequestTimeout)
}

// getHTTPRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an HTTPRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayapi.HTTPRoute) []*string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayapi.Hostname, _ int) *string {
		return kong.String(string(h))
	})
}

// translateHTTPRouteRulesMetaToKongstateRoutes translate the matches and filters under the rules sharing the same backends
// to list of kongstate.Route.
func translateHTTPRouteRulesMetaToKongstateRoutes(
	rulesMeta []httpRouteRuleMeta,
) ([]kongstate.Route, error) {
	rulesGroupedByFilter := groupRulesByFilter(rulesMeta)
	routes := make([]kongstate.Route, 0)

	for _, rulesWithSameFilter := range rulesGroupedByFilter {
		if len(rulesWithSameFilter) == 0 {
			continue
		}
		filters := rulesWithSameFilter[0].Rule.Filters
		// Group the matches for each rule. Then aggregate the matches eligible for the consolidation
		// into a single match group.
		matchGroups := make(map[string]httpRouteMatchMetaList)
		for _, ruleMeta := range rulesWithSameFilter {
			ruleMatchGroups := groupSliceByKeyFn(ruleMeta.matches(), httpRouteMatchMeta.getKey)
			for matchGroupKey, matchGroup := range ruleMatchGroups {
				matchGroups[matchGroupKey] = append(matchGroups[matchGroupKey], matchGroup...)
			}
			// There is a rule that has no matches, which means we should generate a route to match only hostnames for the rule
			// or a catch-all match when hostname is empty.
			if len(ruleMatchGroups) == 0 {
				parentRoute := ruleMeta.parentRoute
				routeName := fmt.Sprintf("httproute.%s.%s.%d.0",
					parentRoute.Namespace, parentRoute.Name, ruleMeta.RuleNumber)
				objectInfo := util.FromK8sObject(parentRoute)
				tags := util.GenerateTagsForObject(parentRoute)
				hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(parentRoute)
				routesFromRule, err := GenerateKongRoutesFromHTTPRouteMatches(
					routeName,
					[]gatewayapi.HTTPRouteMatch{},
					filters,
					objectInfo,
					hostnames,
					tags,
				)
				if err != nil {
					return nil, err
				}
				routes = append(routes, routesFromRule...)
			}
		}

		for _, matchGroup := range matchGroups {
			// Since we have grouped matches by their parent routes, all the matches in the same group are from the same HTTPRoute.
			parentRoute := matchGroup[0].parentRoute
			objectInfo := util.FromK8sObject(parentRoute)
			tags := util.GenerateTagsForObject(parentRoute)
			routeName := translateToKongRouteName(matchGroup, parentRoute.GetNamespace(), parentRoute.GetName())
			matches := matchGroup.httpRouteMatches()
			// Since the grouped matches here are from the same HTTPRoute, it is OK to use the hostnames from the first HTTPRoute.
			hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(parentRoute)

			routesFromMatchGroup, err := GenerateKongRoutesFromHTTPRouteMatches(
				routeName,
				matches,
				filters,
				objectInfo,
				hostnames,
				tags,
			)
			if err != nil {
				return nil, err
			}
			routes = append(routes, routesFromMatchGroup...)
		}
	}
	// Sort the routes by name to produce a stable order.
	sort.Slice(routes, func(i, j int) bool {
		return *routes[i].Name < *routes[j].Name
	})
	return routes, nil
}

// extractUniqueHTTPRoutes extracts unique HTTPRoutes in a grouped list of HTTPRouteRuleMeta.
func extractUniqueHTTPRoutes(rulesMeta []httpRouteRuleMeta) []*gatewayapi.HTTPRoute {
	routes := lo.Map(rulesMeta, func(m httpRouteRuleMeta, _ int) *gatewayapi.HTTPRoute {
		return m.parentRoute
	})
	uniqueRoutes := lo.UniqBy(routes, func(route *gatewayapi.HTTPRoute) k8stypes.NamespacedName {
		return k8stypes.NamespacedName{
			Namespace: route.Namespace,
			Name:      route.Name,
		}
	})
	return uniqueRoutes
}

func translateToKongRouteName(matchesMeta httpRouteMatchMetaList, namespace string, name string) string {
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
		namespace,
		name,
		firstRuleNumber,
		firstMatchNumber,
	)
}

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
	Rule        gatewayapi.HTTPRouteRule
	RuleNumber  int
	parentRoute *gatewayapi.HTTPRoute
}

// getHTTPBackendRefsKey computes a key from a list of backendRefs.
// The order of backedRefs is not important.
func (m httpRouteRuleMeta) getHTTPBackendRefsKey() string {
	return getSortedItemsString(m.Rule.BackendRefs)
}

// getFiltersKey computes a key from a list of filters.
// The order of the filters is not important.
func (m httpRouteRuleMeta) getFiltersKey() string {
	return getSortedItemsString(m.Rule.Filters)
}

// getKongServiceNameByBackendRefs generates service name based on rule's backendRefs and the namespace of the parent HTTPRoute
// to group rules with same backends and from HTTPRoutes in the same namespace to the same Kong service.
// Grouping by namespace of parent HTTPRoute is required, because HTTPRoute from different namespaces may have different reference grants
// to backends in different namespaces, thus the same backend may have different validity in HTTPRoutes from different namespaces.
// The Kong service name is composed by sorted identifier of each backend of the rule, including:
// - namespace (filled from parent HTTPRoute if not given in the rule's backend ref, for sorting)
// - name
// - port number (if exists)
// - weight (if exists)
// Theoretically this can produce a name with at most 5389 characters:
//   - For a backend, maxLen(namespace)=63, maxLen(name) = 253, maxLen(port) = 5, maxLen(weight) = 7,
//     so max length of a single backend is 63+1+253+1+5+1+7=331.
//   - A rule can contain at most 16 backends, so the part from backends can have at most 331*16+15=5311 chars.
//   - A namespace can have at most 63 chars, so the max length of prefix before backends is 9+1+63+1+3+1=78.
//   - So the maximum possible total length is 5311+78=5389 characters.
func (m httpRouteRuleMeta) getKongServiceNameByBackendRefs() string {
	backendRefs := m.Rule.BackendRefs

	backendNames := make([]string, 0, len(backendRefs))
	for _, backendRef := range backendRefs {
		var backendNamespace string
		if backendRef.Namespace == nil {
			backendNamespace = m.parentRoute.Namespace
		} else {
			backendNamespace = string(*backendRef.Namespace)
		}

		backendFullName := backendNamespace + "." + string(backendRef.Name)
		// We only support `Service`s as backends of `HTTPRoute` for now,
		// so we do not need to check for groups and kinds for backendRefs.

		if backendRef.Port != nil {
			backendFullName = backendFullName + "." + strconv.Itoa(int(*backendRef.Port))
		}

		if backendRef.Weight != nil {
			backendFullName = backendFullName + "." + strconv.Itoa(int(*backendRef.Weight))
		}
		backendNames = append(backendNames, backendFullName)
	}
	// Sort the backend names to guarantee that the same set of backend {namespace, name, port, weight}
	// generates the same name.
	sort.Strings(backendNames)
	return fmt.Sprintf("httproute.%s.svc.%s", m.parentRoute.Namespace, strings.Join(backendNames, "_"))
}

func (m *httpRouteRuleMeta) matches() httpRouteMatchMetaList {
	matches := make([]httpRouteMatchMeta, 0, len(m.Rule.Matches))

	for matchNumber, match := range m.Rule.Matches {
		matches = append(matches, httpRouteMatchMeta{
			Match:       &match,
			RuleNumber:  m.RuleNumber,
			MatchNumber: matchNumber,
			parentRoute: m.parentRoute,
		})
	}

	return matches
}

type httpRouteMatchMeta struct {
	Match       *gatewayapi.HTTPRouteMatch
	RuleNumber  int
	MatchNumber int
	parentRoute *gatewayapi.HTTPRoute
}

// getKey computes a key from an HTTPRouteMatch. Two HTTPRouteMatches will generate the same key if their
// parent HTTPRoute, methods, headers, and query parameters are identical.
// HTTPRouteMatches with the same key can be combined into a single Kong route.
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
		Namespace string
		Name      string
		Method    *gatewayapi.HTTPMethod
		Headers   []gatewayapi.HTTPHeaderMatch
		Query     []gatewayapi.HTTPQueryParamMatch
	}{
		Namespace: m.parentRoute.Namespace,
		Name:      m.parentRoute.Name,
		Method:    m.Match.Method,
		Headers:   headers,
		Query:     queryParams,
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

// -----------------------------------------------------------------------------
// HttpRoute Translation - Private - Backend Cache
// -----------------------------------------------------------------------------

type HTTPRouteBackendCache struct {
	backends map[string][]httpRouteRuleMeta
}

func newHTTPRouteBackendCache() *HTTPRouteBackendCache {
	return &HTTPRouteBackendCache{
		backends: make(map[string][]httpRouteRuleMeta),
	}
}

func (c *HTTPRouteBackendCache) addRule(r httpRouteRuleMeta) {
	kongServiceName := r.getKongServiceNameByBackendRefs()
	c.backends[kongServiceName] = append(c.backends[kongServiceName], r)
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

// GenerateKongRoutesFromHTTPRouteMatches converts an HTTPRouteMatches to a slice of Kong Route objects with traditional routes.
// This function assumes that the HTTPRouteMatches share the query params, headers and methods.
func GenerateKongRoutesFromHTTPRouteMatches(
	routeName string,
	matches []gatewayapi.HTTPRouteMatch,
	filters []gatewayapi.HTTPRouteFilter,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []*string,
	tags []*string,
) ([]kongstate.Route, error) {
	if len(matches) == 0 {
		// it's acceptable for an HTTPRoute to have no matches in the rulesets,
		// but only backends as long as there are hostnames. In this case, we
		// match all traffic based on the hostname and leave all other routing
		// options default.
		// for rules with no hostnames, we generate a "catch-all" route for it.
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:         kong.String(routeName),
				Protocols:    kong.StringSlice("http", "https"),
				PreserveHost: kong.Bool(true),
				Tags:         tags,
			},
		}
		r.Hosts = append(r.Hosts, hostnames...)

		return []kongstate.Route{r}, nil
	}

	r := generateKongstateHTTPRoute(routeName, ingressObjectInfo, hostnames)
	r.Tags = tags

	// convert header matching from HTTPRoute to Route format
	headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(matches[0].Headers)
	if err != nil {
		return []kongstate.Route{}, err
	}
	if len(headers) > 0 {
		r.Route.Headers = headers
	}

	// stripPath needs to be disabled by default to be conformant with the Gateway API
	r.StripPath = kong.Bool(false)

	// Check if the route has a RequestRedirect or URLRewrite with non-nil ReplacePrefixMatch - if it does, we need to
	// generate a route for each match as the path is used to modify routes and generate plugins.
	hasRedirectFilter := lo.ContainsBy(filters, func(filter gatewayapi.HTTPRouteFilter) bool {
		return filter.Type == gatewayapi.HTTPRouteFilterRequestRedirect
	})

	routes, err := getRoutesFromMatches(matches, &r, filters, tags, hasRedirectFilter)
	if err != nil {
		return nil, err
	}

	var path string
	if hasURLRewriteWithReplacePrefixMatchFilter := lo.ContainsBy(filters, func(filter gatewayapi.HTTPRouteFilter) bool {
		return filter.Type == gatewayapi.HTTPRouteFilterURLRewrite &&
			filter.URLRewrite.Path != nil &&
			filter.URLRewrite.Path.Type == gatewayapi.PrefixMatchHTTPPathModifier &&
			filter.URLRewrite.Path.ReplacePrefixMatch != nil
	}); hasURLRewriteWithReplacePrefixMatchFilter {
		// In the case of URLRewrite with non-nil ReplacePrefixMatch, we rely on a CEL validation rule that disallows
		// rules with multiple matches if the URLRewrite filter is present. We can be certain that if the filter is
		// present, there is at most only one match. Based on that, we can determine the path from the first match.
		// See: https://github.com/kubernetes-sigs/gateway-api/blob/29e68bffffb9af568e35545305d78d0001a1a0f7/apis/v1/httproute_types.go#L131
		if len(matches) > 0 && matches[0].Path != nil && matches[0].Path.Value != nil {
			path = *matches[0].Path.Value
		}
	}

	// If the redirect filter has not been set, we still need to set the route plugins.
	if !hasRedirectFilter {
		if err := SetRoutePlugins(&r, filters, path, tags, false); err != nil {
			return nil, err
		}
		routes = []kongstate.Route{r}
	}

	return routes, nil
}

func generateKongstateHTTPRoute(routeName string, ingressObjectInfo util.K8sObjectInfo, hostnames []*string) kongstate.Route {
	// build the route object using the method and pathing information
	r := kongstate.Route{
		Ingress: ingressObjectInfo,
		Route: kong.Route{
			Name:         kong.String(routeName),
			Protocols:    kong.StringSlice("http", "https"),
			PreserveHost: kong.Bool(true),
			// metadata tags aren't added here, they're added by the caller
		},
	}

	// attach any hostnames associated with the httproute
	if len(hostnames) > 0 {
		r.Hosts = hostnames
	}

	return r
}

// convertGatewayMatchHeadersToKongRouteMatchHeaders takes an input list of Gateway APIs HTTPHeaderMatch
// and converts these header matching rules to the format expected by go-kong.
func convertGatewayMatchHeadersToKongRouteMatchHeaders(headers []gatewayapi.HTTPHeaderMatch) (map[string][]string, error) {
	// iterate through each provided header match checking for invalid
	// options and otherwise converting to kong type format.
	convertedHeaders := make(map[string][]string)
	for _, header := range headers {
		if _, exists := convertedHeaders[string(header.Name)]; exists {
			return nil, fmt.Errorf("multiple header matches for the same header are not allowed: %s",
				string(header.Name))
		}
		switch {
		case header.Type != nil && *header.Type == gatewayapi.HeaderMatchRegularExpression:
			convertedHeaders[string(header.Name)] = []string{KongHeaderRegexPrefix + header.Value}
		case header.Type == nil || *header.Type == gatewayapi.HeaderMatchExact:
			convertedHeaders[string(header.Name)] = []string{header.Value}
		default:
			return nil, fmt.Errorf("unknown/unsupported header match type: %s", string(*header.Type))
		}
	}

	return convertedHeaders, nil
}

// getRoutesFromMatches converts all the httpRoute matches to the proper set of kong routes.
func getRoutesFromMatches(
	matches []gatewayapi.HTTPRouteMatch,
	route *kongstate.Route,
	filters []gatewayapi.HTTPRouteFilter,
	tags []*string,
	hasRedirectFilter bool,
) ([]kongstate.Route, error) {
	seenMethods := make(map[string]struct{})
	routes := make([]kongstate.Route, 0)

	for _, match := range matches {
		// if the rule specifies the redirectFilter, we cannot put all the paths under the same route,
		// as the kong plugin needs to know the exact path to use to perform redirection.
		if hasRedirectFilter {
			matchRoute := route
			// configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value
			if match.Path != nil {
				paths := generateKongRoutePathFromHTTPRouteMatch(match)
				for _, p := range paths {
					matchRoute.Route.Paths = append(matchRoute.Route.Paths, kong.String(p))
				}
			}

			// configure method matching information about the route if method
			// matching was defined.
			if match.Method != nil {
				method := string(*match.Method)
				if _, ok := seenMethods[method]; !ok {
					matchRoute.Route.Methods = append(matchRoute.Route.Methods, kong.String(string(*match.Method)))
					seenMethods[method] = struct{}{}
				}
			}
			path := ""
			if match.Path.Value != nil {
				path = *match.Path.Value
			}

			// generate kong plugins from rule.filters
			if err := SetRoutePlugins(matchRoute, filters, path, tags, false); err != nil {
				return nil, err
			}

			routes = append(routes, *route)
		} else {
			// Configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value.
			if match.Path != nil {
				for _, path := range generateKongRoutePathFromHTTPRouteMatch(match) {
					route.Route.Paths = append(route.Route.Paths, kong.String(path))
				}
			}

			if match.Method != nil {
				method := string(*match.Method)
				if _, ok := seenMethods[method]; !ok {
					route.Route.Methods = append(route.Route.Methods, kong.String(string(*match.Method)))
					seenMethods[method] = struct{}{}
				}
			}
		}
	}
	return routes, nil
}

func generateKongRoutePathFromHTTPRouteMatch(match gatewayapi.HTTPRouteMatch) []string {
	switch *match.Path.Type {
	case gatewayapi.PathMatchExact:
		return []string{KongPathRegexPrefix + *match.Path.Value + "$"}

	case gatewayapi.PathMatchPathPrefix:
		paths := make([]string, 0, 2)
		path := *match.Path.Value
		paths = append(paths, fmt.Sprintf("%s%s$", KongPathRegexPrefix, path))
		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%s/", path)
		}
		return append(paths, path)

	case gatewayapi.PathMatchRegularExpression:
		return []string{KongPathRegexPrefix + *match.Path.Value}
	}

	return []string{""} // unreachable code
}

// SetRoutePlugins converts HTTPRouteFilter into Kong plugins. The plugins are set into the given kongstate.Route.
// The plugins can be set in two different ways:
// - Direct conversion from the respective HTTPRouteFilter.
// - ExtensionRef to plugins annotation from the ExtensionRef filter.
func SetRoutePlugins(
	route *kongstate.Route,
	filters []gatewayapi.HTTPRouteFilter,
	path string,
	tags []*string,
	expressionsRouterEnabled bool,
) error {
	generatedPlugins, err := generatePluginsFromHTTPRouteFilters(filters, path, tags, expressionsRouterEnabled)
	if err != nil {
		return err
	}
	route.Plugins = append(route.Plugins, generatedPlugins.Plugins...)
	for _, modifier := range generatedPlugins.KongRouteModifiers {
		modifier(route)
	}
	return nil
}

// httpRouteFiltersOriginatedPlugins is a set of Kong plugins generated from HTTPRoute filters along with other
// metadata that is required to apply these plugins to a Kong Route.
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
// It returns httpRouteFiltersOriginatedPlugins which contains:
// - generated plugins
// - Kong Route modifiers that need to be applied to the Kong Route
// - PluginsAnnotation that is generated from the ExtensionRef filter.
func generatePluginsFromHTTPRouteFilters(
	filters []gatewayapi.HTTPRouteFilter,
	path string,
	tags []*string,
	expressionsRouterEnabled bool,
) (httpRouteFiltersOriginatedPlugins, error) {
	if len(filters) == 0 {
		return httpRouteFiltersOriginatedPlugins{}, nil
	}
	var (
		transformerPlugins          []transformerPlugin
		kongPlugins                 []kong.Plugin
		pluginNamesFromExtensionRef []string
		kongRouteModifiers          []kongRouteModifier
	)

	for _, filter := range filters {
		switch filter.Type {
		case gatewayapi.HTTPRouteFilterRequestHeaderModifier:
			transformerPlugins = append(transformerPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))

		case gatewayapi.HTTPRouteFilterRequestRedirect:
			kongPlugin, transformerPlugin := generateRequestRedirectKongPlugin(filter.RequestRedirect, path)
			kongPlugins = append(kongPlugins, kongPlugin)
			transformerPlugins = append(transformerPlugins, transformerPlugin)

		case gatewayapi.HTTPRouteFilterResponseHeaderModifier:
			transformerPlugins = append(transformerPlugins, generateResponseHeaderModifierKongPlugin(filter.ResponseHeaderModifier))

		case gatewayapi.HTTPRouteFilterExtensionRef:
			plugin, err := generateExtensionRefKongPlugin(filter.ExtensionRef)
			if err != nil {
				return httpRouteFiltersOriginatedPlugins{}, err
			}
			pluginNamesFromExtensionRef = append(pluginNamesFromExtensionRef, plugin)

		case gatewayapi.HTTPRouteFilterURLRewrite:
			plugins, routeModifiers, err := generateRequestTransformerForURLRewrite(filter.URLRewrite, path, expressionsRouterEnabled)
			if err != nil {
				return httpRouteFiltersOriginatedPlugins{}, err
			}
			transformerPlugins = append(transformerPlugins, plugins...)
			kongRouteModifiers = append(kongRouteModifiers, routeModifiers...)

		case gatewayapi.HTTPRouteFilterRequestMirror:
			// not supported
			return httpRouteFiltersOriginatedPlugins{}, fmt.Errorf("httpFilter %s unsupported", filter.Type)
		}
	}

	// It's possible the above loop generates multiple transformerPlugins of the same type, so we need to merge them.
	// It can happen for example when both RequestHeaderModifier and HTTPRouteFilterURLRewrite filters are present.
	transformerPlugins, err := mergePluginsOfTheSameType(transformerPlugins)
	if err != nil {
		return httpRouteFiltersOriginatedPlugins{}, fmt.Errorf("failed to merge transformerPlugins of the same type: %w", err)
	}
	kongPlugins = append(kongPlugins, transformerPluginsToKongPlugins(transformerPlugins)...)

	for i := range kongPlugins {
		// This plugin is derived from an HTTPRoute filter, not a KongPlugin, so we apply tags indicating that
		// HTTPRoute as the parent Kubernetes resource for these generated transformerPlugins.
		kongPlugins[i].Tags = tags
	}

	if len(pluginNamesFromExtensionRef) > 0 {
		kongRouteModifiers = append(kongRouteModifiers, generateKongRouteModifierFromExtensionRef(pluginNamesFromExtensionRef))
	}

	return httpRouteFiltersOriginatedPlugins{
		Plugins:            kongPlugins,
		KongRouteModifiers: kongRouteModifiers,
	}, nil
}

// mergePluginsOfTheSameType merges plugins of the same type into a single plugin with merged configurations.
func mergePluginsOfTheSameType(plugins []transformerPlugin) ([]transformerPlugin, error) {
	pluginsByName := lo.GroupBy(plugins, func(p transformerPlugin) transformerPluginType {
		return p.Type // Name is effectively a plugin type.
	})
	for pluginName, plugins := range pluginsByName {
		// If we produced multiple plugins of the same type, we need to merge their configurations now.
		if len(plugins) > 1 {
			mergedPlugin := transformerPlugin{}
			for _, plugin := range plugins {
				if err := mergo.Merge(&mergedPlugin, plugin, mergo.WithAppendSlice); err != nil {
					// Should never happen as we're passing the same type of objects.
					return nil, fmt.Errorf("failed to merge %q plugin configurations: %w", pluginName, err)
				}
			}
			pluginsByName[pluginName] = []transformerPlugin{mergedPlugin}
		}
	}

	// Sort the plugins by name to ensure that the order is deterministic.
	mergedPlugins := lo.Flatten(lo.Values(pluginsByName))
	sort.Slice(mergedPlugins, func(i, j int) bool {
		return mergedPlugins[i].Type < mergedPlugins[j].Type
	})
	return mergedPlugins, nil
}

func transformerPluginsToKongPlugins(plugins []transformerPlugin) []kong.Plugin {
	kongPlugins := make([]kong.Plugin, 0, len(plugins))
	for _, plugin := range plugins {
		kongPlugins = append(kongPlugins, transformerPluginToKongPlugin(plugin))
	}
	return kongPlugins
}

func transformerPluginToKongPlugin(plugin transformerPlugin) kong.Plugin {
	res := kong.Plugin{
		Name:   kong.String(string(plugin.Type)),
		Config: kong.Configuration{},
	}
	if !cmp.Equal(plugin.Replace, TransformerPluginReplaceConfig{}) {
		res.Config["replace"] = plugin.Replace
	}
	if !cmp.Equal(plugin.Add, TransformerPluginConfig{}) {
		res.Config["add"] = plugin.Add
	}
	if !cmp.Equal(plugin.Append, TransformerPluginConfig{}) {
		res.Config["append"] = plugin.Append
	}
	if !cmp.Equal(plugin.Remove, TransformerPluginConfig{}) {
		res.Config["remove"] = plugin.Remove
	}
	return res
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
func generateRequestRedirectKongPlugin(modifier *gatewayapi.HTTPRequestRedirectFilter, path string) (kong.Plugin, transformerPlugin) {
	requestTerminationPlugin := kong.Plugin{
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

	transformerPlugin := transformerPlugin{
		Type: transformerPluginTypeResponse,
		Add: TransformerPluginConfig{
			Headers: []string{locationHeader},
		},
	}

	return requestTerminationPlugin, transformerPlugin
}

func generateExtensionRefKongPlugin(modifier *gatewayapi.LocalObjectReference) (string, error) {
	if modifier.Group != "configuration.konghq.com" || modifier.Kind != "KongPlugin" {
		return "", fmt.Errorf("plugin %s/%s unsupported", modifier.Group, modifier.Kind)
	}
	return string(modifier.Name), nil
}

// generateRequestHeaderModifierKongPlugin converts a gatewayapi.HTTPRequestHeaderFilter into a
// kong.Plugin of type request-transformer.
func generateRequestHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter) transformerPlugin {
	return generateHeaderModifierKongPlugin(modifier, transformerPluginTypeRequest)
}

// generateResponseHeaderModifierKongPlugin converts a gatewayapi.HTTPResponseHeaderFilter into a
// kong.Plugin of type response-transformer.
func generateResponseHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter) transformerPlugin {
	return generateHeaderModifierKongPlugin(modifier, transformerPluginTypeResponse)
}

type transformerPluginType string

const (
	transformerPluginTypeRequest  transformerPluginType = "request-transformer"
	transformerPluginTypeResponse transformerPluginType = "response-transformer"
)

// transformerPlugin is a configuration for request-transformer and response-transformer plugins.
type transformerPlugin struct {
	// Type is the type of the transformer plugin (request-transformer or response-transformer).
	Type transformerPluginType `json:"-"`

	Replace TransformerPluginReplaceConfig `json:"replace,omitempty"`
	Add     TransformerPluginConfig        `json:"add,omitempty"`
	Append  TransformerPluginConfig        `json:"append,omitempty"`
	Remove  TransformerPluginConfig        `json:"remove,omitempty"`
}

// TransformerPluginConfig is a configuration for request-transformer and response-transformer plugins'
// "add", "append", and "remove" fields.
type TransformerPluginConfig struct {
	Headers []string `json:"headers,omitempty"`
}

// TransformerPluginReplaceConfig is a configuration for request-transformer and response-transformer plugins'
// "replace" field.
type TransformerPluginReplaceConfig struct {
	Headers []string `json:"headers,omitempty"`
	URI     string   `json:"uri,omitempty"`
}

func generateHeaderModifierKongPlugin(modifier *gatewayapi.HTTPHeaderFilter, pluginType transformerPluginType) transformerPlugin {
	plugin := transformerPlugin{
		Type: pluginType,
	}

	// modifier.Set is converted to a pair composed of "replace" and "add"
	if modifier.Set != nil {
		setModifiers := make([]string, 0, len(modifier.Set))
		for _, s := range modifier.Set {
			setModifiers = append(setModifiers, kongHeaderFormatter(s))
		}
		plugin.Replace = TransformerPluginReplaceConfig{
			Headers: setModifiers,
		}
		plugin.Add = TransformerPluginConfig{
			Headers: setModifiers,
		}
	}

	// modifier.Add is converted to "append"
	if modifier.Add != nil {
		appendModifiers := make([]string, 0, len(modifier.Add))
		for _, a := range modifier.Add {
			appendModifiers = append(appendModifiers, kongHeaderFormatter(a))
		}
		plugin.Append = TransformerPluginConfig{
			Headers: appendModifiers,
		}
	}

	if modifier.Remove != nil {
		plugin.Remove = TransformerPluginConfig{
			Headers: modifier.Remove,
		}
	}

	return plugin
}

func kongHeaderFormatter(header gatewayapi.HTTPHeader) string {
	return fmt.Sprintf("%s:%s", header.Name, header.Value)
}

func generateRequestTransformerForURLRewrite(
	filter *gatewayapi.HTTPURLRewriteFilter,
	path string,
	expressionsRouterEnabled bool,
) ([]transformerPlugin, []kongRouteModifier, error) {
	if filter == nil {
		return nil, nil, fmt.Errorf("%s is not provided", gatewayapi.HTTPRouteFilterURLRewrite)
	}
	if filter.Path == nil && filter.Hostname == nil {
		return nil, nil, fmt.Errorf("%s missing Path and Hostname", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	var (
		plugins        []transformerPlugin
		routeModifiers []kongRouteModifier
	)
	if filter.Path != nil {
		plugin, modifier, err := generateRequestTransformerForURLRewritePath(filter, path, expressionsRouterEnabled)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate request-transformer plugin for Path: %w", err)
		}
		plugins = append(plugins, plugin)
		if modifier != nil {
			routeModifiers = append(routeModifiers, modifier)
		}
	}
	if filter.Hostname != nil {
		plugins = append(plugins, generateRequestTransformerForURLRewriteHostname(filter))
	}

	return plugins, routeModifiers, nil
}

func generateRequestTransformerForURLRewritePath(
	filter *gatewayapi.HTTPURLRewriteFilter,
	path string,
	expressionsRouterEnabled bool,
) (transformerPlugin, kongRouteModifier, error) {
	switch filter.Path.Type {
	case gatewayapi.FullPathHTTPPathModifier:
		plugin, err := generateRequestTransformerForURLRewriteFullPath(filter)
		if err != nil {
			return transformerPlugin{}, nil, fmt.Errorf("failed to generate request-transformer plugin for %s: %w", gatewayapi.HTTPRouteFilterURLRewrite, err)
		}
		return plugin, nil, nil

	case gatewayapi.PrefixMatchHTTPPathModifier:
		plugin, routeModifier := generateRequestTransformerForURLRewritePrefixMatch(filter, path, expressionsRouterEnabled)
		return plugin, routeModifier, nil
	default:
		return transformerPlugin{}, nil, fmt.Errorf("unsupported path type %s for %s", filter.Path.Type, gatewayapi.HTTPRouteFilterURLRewrite)
	}
}

func generateRequestTransformerForURLRewriteHostname(
	filter *gatewayapi.HTTPURLRewriteFilter,
) transformerPlugin {
	return transformerPlugin{
		Type: transformerPluginTypeRequest,
		Replace: TransformerPluginReplaceConfig{
			Headers: []string{
				fmt.Sprintf("host:%s", string(*filter.Hostname)),
			},
		},
		Add: TransformerPluginConfig{
			Headers: []string{
				fmt.Sprintf("host:%s", string(*filter.Hostname)),
			},
		},
	}
}

// generateRequestTransformerForURLRewriteFullPath generates a request-transformer plugin for the URLRewrite filter
// with a FullPathHTTPPathModifier.
func generateRequestTransformerForURLRewriteFullPath(filter *gatewayapi.HTTPURLRewriteFilter) (transformerPlugin, error) {
	if filter.Path.ReplaceFullPath == nil {
		return transformerPlugin{}, fmt.Errorf("%s missing ReplaceFullPath", gatewayapi.HTTPRouteFilterURLRewrite)
	}

	return transformerPlugin{
		Type: transformerPluginTypeRequest,
		Replace: TransformerPluginReplaceConfig{
			URI: *filter.Path.ReplaceFullPath,
		},
	}, nil
}

// generateRequestTransformerForURLRewritePrefixMatch generates a request-transformer plugin and a route modifier
// for the URLRewrite filter with a PrefixMatchHTTPPathModifier.
func generateRequestTransformerForURLRewritePrefixMatch(
	filter *gatewayapi.HTTPURLRewriteFilter,
	path string,
	expressionsRouterEnabled bool,
) (transformerPlugin, kongRouteModifier) {
	// Normalize the path before passing it down.
	path = normalizePath(path)

	return transformerPlugin{
		Type: transformerPluginTypeRequest,
		Replace: TransformerPluginReplaceConfig{
			URI: generateRequestTransformerReplaceURIForURLRewritePrefixMatch(
				filter.Path.ReplacePrefixMatch,
				path,
			),
		},
	}, generateKongRouteModifierForURLRewritePrefixMatch(path, expressionsRouterEnabled)
}

// generateRequestTransformerReplaceURIForURLRewritePrefixMatch generates the replacement URI for the request-transformer
// plugin for the URLRewrite filter with a PrefixMatchHTTPPathModifier.
func generateRequestTransformerReplaceURIForURLRewritePrefixMatch(
	replacePrefixMatch *string,
	path string,
) string {
	if replacePrefixMatch == nil {
		// If no ReplacePrefixMatch is provided, we need to use a slash as the replacement URI.
		return "/"
	}

	// If the ReplacePrefixMatch is provided, we need to use it as the replacement URI.
	// Trim the trailing slash from the ReplacePrefixMatch to avoid double slashes in the final URI.
	*replacePrefixMatch = strings.TrimSuffix(*replacePrefixMatch, "/")
	pathIsRoot := isPathRoot(path)

	// In the case of an empty replacePrefixMatch, we need to make sure that the path will always start with a slash,
	// even if we have no capture group from the incoming request's URI.
	if *replacePrefixMatch == "" {
		// If path is "/", we need to add a slash before URI captures because the capture group won't include
		// the leading slash.
		if pathIsRoot {
			// The below is a Lua ternary operator that checks if the captured group is nil, and if so, replaces it with
			// a slash. Otherwise, it appends the captured group to a slash.
			return `$(uri_captures[1] == nil and "/" or "/" .. uri_captures[1])`
		}

		// Otherwise, we do not need to add a leading slash before URI captures.
		// The below is a Lua ternary operator that checks if the captured group is nil, and if so, replaces it with
		// a slash. Otherwise, it returns the captured group (in this case the captured group will always have a
		// leading slash).
		return `$(uri_captures[1] == nil and "/" or uri_captures[1])`
	}

	// Otherwise, we concatenate the replacement URI with the captured group.
	// If path is "/", we need to add a slash before URI captures because the capture group won't include the
	// leading slash.
	if pathIsRoot {
		// The below Lua ternary operator checks if the captured group is nil, and if so, replaces it with
		// an empty string (as we already know replacePrefixMatch is not empty so the resulting path will always have
		// a leading slash). Otherwise, it appends the captured group to a slash (as the captured group won't
		// have the leading slash).
		return fmt.Sprintf(`%s$(uri_captures[1] == nil and "" or "/" .. uri_captures[1])`, *replacePrefixMatch)
	}
	// Simply concatenate the replacement URI with the captured group as the captured group will always have a
	// leading slash.
	return fmt.Sprintf(`%s$(uri_captures[1])`, *replacePrefixMatch)
}

// generateKongRouteModifierForURLRewritePrefixMatch generates a Kong route modifier for the URLRewrite filter with a
// PrefixMatchHTTPPathModifier.
// One of the paths will match the exact path, and the other will match subpaths with a capture group.
// The capture group can be used in the request-transformer plugin to replace the path.
// Please note that depending on the path, the capture group will:
// - Exclude the leading slash if the path is root.
// - Include the leading slash otherwise.
func generateKongRouteModifierForURLRewritePrefixMatch(path string, expressionsRouterEnabled bool) func(route *kongstate.Route) {
	pathIsRoot := isPathRoot(path)

	// If expressions router is enabled, we need to set the expression on the Kong Route.
	if expressionsRouterEnabled {
		return func(route *kongstate.Route) {
			exactPrefixPredicate := atc.NewPredicateHTTPPath(atc.OpEqual, path)
			subpathsPredicate := func() atc.Predicate {
				if pathIsRoot {
					// If the path is "/", we don't capture the slash as Kong Route's path has to begin with a slash.
					// If we captured the slash, we'd generate "(/.*)", and it'd be rejected by Kong.
					return atc.NewPredicateHTTPPath(atc.OpRegexMatch, "^/(.*)")
				}
				// If the path is not "/", i.e. it has a prefix, we capture the slash to make it possible to
				// route "/prefix" to "/replacement" and "/prefix/" to "/replacement/" correctly.
				return atc.NewPredicateHTTPPath(atc.OpRegexMatch, fmt.Sprintf("^%s(/.*)", path))
			}()
			route.Route.Expression = lo.ToPtr(atc.Or(exactPrefixPredicate, subpathsPredicate).Expression())
		}
	}

	// Otherwise, we set the Kong Route's paths.
	return func(route *kongstate.Route) {
		paths := make([]*string, 0, 2)
		// The first path matches the exact path.
		paths = append(paths, lo.ToPtr(fmt.Sprintf("%s%s$", KongPathRegexPrefix, path)))
		// The second path matches subpaths, including a single slash.
		if pathIsRoot {
			// If the path is "/", we don't capture the slash as Kong Route's path has to begin with a slash.
			// If we captured the slash, we'd generate "(/.*)", and it'd be rejected by Kong.
			paths = append(paths, lo.ToPtr(fmt.Sprintf("%s/(.*)", KongPathRegexPrefix)))
		} else {
			// If the path is not "/", i.e. it has a prefix, we capture the slash to make it possible to
			// route "/prefix" to "/replacement" and "/prefix/" to "/replacement/" correctly.
			paths = append(paths, lo.ToPtr(
				fmt.Sprintf("%s%s(/.*)", KongPathRegexPrefix, path)),
			)
		}
		route.Paths = paths
	}
}

// normalizePath normalizes the path by making it:
// - a slash if it's empty or "/",
// - a path without trailing slash otherwise.
func normalizePath(path string) string {
	if path == "/" || path == "" {
		return "/"
	}
	return strings.TrimSuffix(path, "/")
}

func isPathRoot(path string) bool {
	return path == "/"
}
