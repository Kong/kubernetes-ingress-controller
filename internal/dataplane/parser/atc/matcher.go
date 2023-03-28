package atc

import (
	"fmt"
	"strings"
)

type Matcher interface {
	// Expression returns a string representation of the Matcher.
	Expression() string

	// IsEmpty() returns a boolean indicating if the Matcher is empty. It is true if the Matcher is an empty struct,
	// if the Matcher has zero subMatchers, or if a single-predicate Matcher has no value.
	IsEmpty() bool
}

var (
	_ Matcher = &OrMatcher{}
	_ Matcher = &AndMatcher{}
)

type OrMatcher struct {
	subMatchers []Matcher
}

func (m *OrMatcher) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.subMatchers) == 0
}

func (m *OrMatcher) Expression() string {
	if m == nil {
		return ""
	}
	if m.IsEmpty() {
		return ""
	}
	if len(m.subMatchers) == 1 {
		return m.subMatchers[0].Expression()
	}

	var grouped []string

	for _, m := range m.subMatchers {
		grouped = append(grouped, fmt.Sprintf("(%s)", m.Expression()))
	}

	return strings.Join(grouped, " || ")
}

func (m *OrMatcher) Or(matcher Matcher) *OrMatcher {
	if matcher == nil {
		return m
	}
	if !matcher.IsEmpty() {
		m.subMatchers = append(m.subMatchers, matcher)
	}
	return m
}

func Or(matchers ...Matcher) *OrMatcher {
	actual := []Matcher{}
	for _, m := range matchers {
		if m != nil {
			if !m.IsEmpty() {
				actual = append(actual, m)
			}
		}
	}
	return &OrMatcher{
		subMatchers: actual,
	}
}

type AndMatcher struct {
	subMatchers []Matcher
}

func (m *AndMatcher) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.subMatchers) == 0
}

func (m *AndMatcher) Expression() string {
	if m == nil {
		return ""
	}
	if m.IsEmpty() {
		return ""
	}
	if len(m.subMatchers) == 1 {
		return m.subMatchers[0].Expression()
	}

	var grouped []string

	for _, m := range m.subMatchers {
		grouped = append(grouped, fmt.Sprintf("(%s)", m.Expression()))
	}

	return strings.Join(grouped, " && ")
}

func (m *AndMatcher) And(matcher Matcher) *AndMatcher {
	if matcher == nil {
		return m
	}
	if !matcher.IsEmpty() {
		m.subMatchers = append(m.subMatchers, matcher)
	}
	return m
}

func And(matchers ...Matcher) *AndMatcher {
	actual := []Matcher{}
	for _, m := range matchers {
		if m != nil {
			if !m.IsEmpty() {
				actual = append(actual, m)
			}
		}
	}
	return &AndMatcher{
		subMatchers: actual,
	}
}
