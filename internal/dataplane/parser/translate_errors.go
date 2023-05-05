package parser

import "errors"

// REVIEW: use error variables in translator package instead?

var (
	errRouteValidationNoRules                          = errors.New("no rules provided")
	errRouteValidationQueryParamMatchesUnsupported     = errors.New("query param matches are not yet supported")
	errRouteValidationNoMatchRulesOrHostnamesSpecified = errors.New("no match rules or hostnames specified")
)
