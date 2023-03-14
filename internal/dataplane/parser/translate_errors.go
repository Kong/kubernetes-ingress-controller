package parser

import "errors"

var (
	errRouteValidationNoRules                          = errors.New("no rules provided")
	errRouteValidationQueryParamMatchesUnsupported     = errors.New("query param matches are not yet supported")
	errRouteValidationNoMatchRulesOrHostnamesSpecified = errors.New("no match rules or hostnames specified")
)
