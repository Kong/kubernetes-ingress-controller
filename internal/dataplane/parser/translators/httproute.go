package translators

import (
	"sort"
	"strconv"
	"strings"

	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// HTTPRouteTranslationMeta is a translation of a single HTTPRoute into metadata
// that can be used to instantiate Kong routes and services.
// Rules from this object should route traffic to BackendRefs from this object.
type HTTPRouteTranslationMeta struct {
	BackendRefs  []gatewayv1beta1.HTTPBackendRef
	Rules        []gatewayv1beta1.HTTPRouteRule
	RulesNumbers []int
}

// TranslateHTTPRoutes translates a list of HTTPRoutes into a list of HTTPRouteTranslationMeta
// objects that can be used to instantiate Kong routes and services.
// The translation is done by grouping the HTTPRoutes by their backendRefs.
// This means that all the rules of a single HTTPRoute will be grouped together
// if they share the same backendRefs.
func TranslateHTTPRoute(route *gatewayv1beta1.HTTPRoute) []*HTTPRouteTranslationMeta {
	index := newHTTPRouteTranslationIndex()
	index.addRoute(route)
	return index.translate()
}

// -----------------------------------------------------------------------------
// HTTPRoute Translation - Private - Index
// -----------------------------------------------------------------------------

// httpRouteTranslationIndex aggregates all rules routing to the same backends group.
type httpRouteTranslationIndex struct {
	backendsRules map[httpBackendRefsKey]*rulesEntry
}

type rulesEntry struct {
	rules       []gatewayv1beta1.HTTPRouteRule
	ruleNumbers []int
}

// newHTTPRouteTranslationIndex creates a new httpRouteTranslationIndex.
func newHTTPRouteTranslationIndex() *httpRouteTranslationIndex {
	return &httpRouteTranslationIndex{
		backendsRules: make(map[httpBackendRefsKey]*rulesEntry),
	}
}

// addRoute an HTTPRoute to the index, grouping the rules by their backendRefs.
func (i *httpRouteTranslationIndex) addRoute(route *gatewayv1beta1.HTTPRoute) {
	for ruleNumber, rule := range route.Spec.Rules {
		backendRefsKey := getHTTPBackendRefsKey(rule.BackendRefs...)
		entry, ok := i.backendsRules[backendRefsKey]
		if !ok {
			entry = &rulesEntry{
				rules:       []gatewayv1beta1.HTTPRouteRule{},
				ruleNumbers: []int{},
			}
			i.backendsRules[backendRefsKey] = entry
		}
		entry.rules = append(entry.rules, rule)
		entry.ruleNumbers = append(entry.ruleNumbers, ruleNumber)
	}
}

// translate the index into a list of HTTPRouteTranslationMeta objects.
func (i *httpRouteTranslationIndex) translate() []*HTTPRouteTranslationMeta {
	translations := make([]*HTTPRouteTranslationMeta, 0)
	for _, rulesEntry := range i.backendsRules {
		// get the backendRefs from any rule, as they are all the same
		backendRefs := rulesEntry.rules[0].BackendRefs

		translations = append(translations, &HTTPRouteTranslationMeta{
			BackendRefs:  backendRefs,
			Rules:        rulesEntry.rules,
			RulesNumbers: rulesEntry.ruleNumbers,
		})
	}

	return translations
}

// -----------------------------------------------------------------------------
// HttpRoute Translation - Private - Metadata
// -----------------------------------------------------------------------------

// httpBackendRefsKey is a key computed from a list of backendRefs.
type httpBackendRefsKey string

// getHTTPBackendRefsKey computes a key from a list of backendRefs.
func getHTTPBackendRefsKey(backendRefs ...gatewayv1beta1.HTTPBackendRef) httpBackendRefsKey {
	backendKeys := make([]string, 0, len(backendRefs))

	// create a list of backend keys
	for _, backendRef := range backendRefs {
		var backendKey strings.Builder

		if backendRef.Group != nil {
			backendKey.WriteString(string(*backendRef.Group))
		}
		backendKey.WriteString(".")

		if backendRef.Kind != nil {
			backendKey.WriteString(string(*backendRef.Kind))
		}
		backendKey.WriteString(".")

		if backendRef.Namespace != nil {
			backendKey.WriteString(string(*backendRef.Namespace))
		}
		backendKey.WriteString(".")

		backendKey.WriteString(string(backendRef.Name))
		backendKey.WriteString(".")

		if backendRef.Port != nil {
			backendKey.WriteString(strconv.Itoa(int(*backendRef.Port)))
		}
		backendKey.WriteString(".")

		if backendRef.Weight != nil {
			backendKey.WriteString(strconv.Itoa(int(*backendRef.Weight)))
		}

		backendKeys = append(backendKeys, backendKey.String())
	}
	sort.Strings(backendKeys)

	return httpBackendRefsKey(strings.Join(backendKeys, ";"))
}
