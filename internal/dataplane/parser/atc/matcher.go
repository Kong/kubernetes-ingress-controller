package atc

import (
	"fmt"
	"strings"
)

type Matcher interface {
	Expression() string
}

var (
	_ Matcher = &OrMatcher{}
	_ Matcher = &AndMatcher{}
)

type OrMatcher struct {
	subMatchers []Matcher
}

func (m *OrMatcher) Expression() string {
	if m == nil {
		return ""
	}
	expressions := []string{}
	for _, subMathcher := range m.subMatchers {
		exp := subMathcher.Expression()
		if len(exp) == 0 {
			continue
		}
		expressions = append(expressions, exp)
	}

	if len(expressions) == 0 {
		return ""
	}
	if len(expressions) == 1 {
		return expressions[0]
	}

	for i, exp := range expressions {
		expressions[i] = fmt.Sprintf("(%s)", exp)
	}

	return strings.Join(expressions, " || ")
}

func (m *OrMatcher) Or(matcher Matcher) *OrMatcher {
	if matcher == nil {
		return m
	}
	m.subMatchers = append(m.subMatchers, matcher)
	return m
}

func Or(matchers ...Matcher) *OrMatcher {
	ret := &OrMatcher{}
	for _, m := range matchers {
		if m == nil {
			continue
		}
		ret.subMatchers = append(ret.subMatchers, m)
	}
	return ret
}

type AndMatcher struct {
	subMatchers []Matcher
}

func (m *AndMatcher) Expression() string {
	if m == nil {
		return ""
	}
	expressions := []string{}
	for _, subMathcher := range m.subMatchers {
		exp := subMathcher.Expression()
		if len(exp) == 0 {
			continue
		}
		expressions = append(expressions, exp)
	}

	if len(expressions) == 0 {
		return ""
	}
	if len(expressions) == 1 {
		return expressions[0]
	}

	for i, exp := range expressions {
		expressions[i] = fmt.Sprintf("(%s)", exp)
	}

	return strings.Join(expressions, " && ")
}

func (m *AndMatcher) And(matcher Matcher) *AndMatcher {
	if matcher == nil {
		return m
	}
	m.subMatchers = append(m.subMatchers, matcher)
	return m
}

func And(matchers ...Matcher) *AndMatcher {
	ret := &AndMatcher{}
	for _, m := range matchers {
		if m == nil {
			continue
		}
		ret.subMatchers = append(ret.subMatchers, m)
	}
	return ret
}
