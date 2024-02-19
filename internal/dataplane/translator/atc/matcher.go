package atc

import (
	"fmt"
	"strings"
)

// Matcher is a sub-expression within a Kong router expression. It can be a single predicate expression, a
// group of predicates joined by logical operators, or a recursive combination of either of the previous components.
type Matcher interface {
	// Expression returns a string representation of the Matcher that could be a valid Kong route expression.
	Expression() string

	// IsEmpty() returns a boolean indicating if the Matcher is empty. It is true if the Matcher is an empty struct,
	// if the Matcher has zero subMatchers, or if a single-predicate Matcher has no value.
	IsEmpty() bool
}

var (
	_ Matcher = &OrMatcher{}
	_ Matcher = &AndMatcher{}
	_ Matcher = &NotMatcher{}
)

// OrMatcher is a group of Matchers joined by logical ORs.
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
	if m == nil || m.IsEmpty() {
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

// Or appends an additional Matcher to an existing OrMatcher. If the given Matcher is empty, it returns the original
// OrMatcher.
func (m *OrMatcher) Or(matcher Matcher) *OrMatcher {
	if matcher != nil && !matcher.IsEmpty() {
		m.subMatchers = append(m.subMatchers, matcher)
	}
	return m
}

// Or constructs an OrMatcher from a list of Matchers. If any of the given Matchers is empty, Or skips adding it.
func Or(matchers ...Matcher) *OrMatcher {
	actual := []Matcher{}
	for _, m := range matchers {
		if m != nil && !m.IsEmpty() {
			actual = append(actual, m)
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
	if m == nil || m.IsEmpty() {
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

// And appends an additional Matcher to an existing AndMatcher. If the given Matcher is empty, it returns the original
// AndMatcher.
func (m *AndMatcher) And(matcher Matcher) *AndMatcher {
	if matcher != nil && !matcher.IsEmpty() {
		m.subMatchers = append(m.subMatchers, matcher)
	}
	return m
}

// And constructs an AndMatcher from a list of Matchers. If any of the given Matchers is empty, And skips adding it.
func And(matchers ...Matcher) *AndMatcher {
	actual := []Matcher{}
	for _, m := range matchers {
		if m != nil && !m.IsEmpty() {
			actual = append(actual, m)
		}
	}
	return &AndMatcher{
		subMatchers: actual,
	}
}

// NotMatcher is a matcher which negates the internal submatcher.
type NotMatcher struct {
	subMatcher Matcher
}

// Not returns a matcher that negates the matcher in the parameter.
func Not(m Matcher) Matcher {
	return &NotMatcher{
		subMatcher: m,
	}
}

func (m *NotMatcher) IsEmpty() bool {
	return m == nil || m.subMatcher.IsEmpty()
}

func (m *NotMatcher) Expression() string {
	if m == nil || m.IsEmpty() {
		return ""
	}
	return fmt.Sprintf("!(%s)", m.subMatcher.Expression())
}
