package parser

import "errors"

var (
	errRouteValidationNoRules                          = errors.New("no rules provided")
	errRouteValidationMissingBackendRefs               = errors.New("missing backendRef in rule")
	errRouteValidationQueryParamMatchesUnsupported     = errors.New("query param matches are not yet supported")
	errRouteValidationNoMatchRulesOrHostnamesSpecified = errors.New("no match rules or hostnames specified")
)
