package atc

import "net/http"

type Matcher interface {
	Expression() string
	Matches(*http.Request) bool
}

var _ Matcher = &OrMatcher{}
var _ Matcher = &AndMatcher{}

type OrMatcher struct {
	subMatchers []Matcher
}

func (m *OrMatcher) Expression() string {
	if len(m.subMatchers) == 0 {
		return ""
	}
	if len(m.subMatchers) == 1 {
		return m.subMatchers[0].Expression()
	}

	ret := ""
	for i, subMathcher := range m.subMatchers {
		exp := "( " + subMathcher.Expression() + " )"
		if i != len(m.subMatchers)-1 {
			exp = exp + " || "
		}
		ret = ret + exp
	}
	return ret
}

func (m *OrMatcher) Matches(req *http.Request) bool {
	for _, subMatcher := range m.subMatchers {
		if subMatcher.Matches(req) {
			return true
		}
	}
	return false
}

func (m *OrMatcher) Or(matcher Matcher) *OrMatcher {
	m.subMatchers = append(m.subMatchers, matcher)
	return m
}

func Or(matchers ...Matcher) *OrMatcher {
	return &OrMatcher{
		subMatchers: matchers,
	}
}

type AndMatcher struct {
	subMatchers []Matcher
}

func (m *AndMatcher) Expression() string {
	if len(m.subMatchers) == 0 {
		return ""
	}
	if len(m.subMatchers) == 1 {
		return m.subMatchers[0].Expression()
	}

	ret := ""
	for i, subMathcher := range m.subMatchers {
		exp := " (" + subMathcher.Expression() + ") "
		if i != len(m.subMatchers)-1 {
			exp = exp + " && "
		}
		ret = ret + exp
	}
	return ret
}

func (m *AndMatcher) Matches(req *http.Request) bool {
	for _, subMatcher := range m.subMatchers {
		if !subMatcher.Matches(req) {
			return false
		}
	}

	return true
}

func (m *AndMatcher) And(matcher Matcher) *AndMatcher {
	m.subMatchers = append(m.subMatchers, matcher)
	return m
}

func And(matchers ...Matcher) *AndMatcher {
	return &AndMatcher{
		subMatchers: matchers,
	}
}
