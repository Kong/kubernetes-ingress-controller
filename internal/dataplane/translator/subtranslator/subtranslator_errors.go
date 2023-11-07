package subtranslator

import "errors"

var (
	ErrRouteValidationNoRules                          = errors.New("no rules provided")
	ErrRouteValidationQueryParamMatchesUnsupported     = errors.New("query param matches are not yet supported")
	ErrRouteValidationNoMatchRulesOrHostnamesSpecified = errors.New("no match rules or hostnames specified")
	ErrRotueValidationRuleNoBackendRef                 = errors.New("no backendRefs in rule")
)
