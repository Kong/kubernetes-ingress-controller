package atc

type HostMatchType int

const (
	HostMatchExact HostMatchType = iota
	HostMatchWildcard
)

type MatchRuleHost struct {
	Type HostMatchType
	Host string
}

type PathMatchType int

const (
	PathMatchExact PathMatchType = iota
	PathMatchPrefix
	PathMatchRegex
)

type MatchRulePath struct {
	Type PathMatchType
	Path string
}

type HeaderMatchType int

const (
	HeaderMatchExact HeaderMatchType = iota
	HeaderMatchRegex
)

type MatchRuleHeader struct {
	Type   HeaderMatchType
	Values []string
}

type MatchRules struct {
	Protocols []string
	Methods   []string
	Hosts     []MatchRuleHost
	Paths     []MatchRulePath
	Headers   map[string]MatchRuleHeader
	SNIs      []string
}

func (r *MatchRules) Merge(other *MatchRules, overrideHeaders bool) *MatchRules {
	if r == nil {
		return other
	}
	if other == nil {
		return r
	}

	r.Protocols = append(r.Protocols, other.Protocols...)
	r.Methods = append(r.Methods, other.Methods...)
	r.Hosts = append(r.Hosts, other.Hosts...)
	r.Paths = append(r.Paths, other.Paths...)
	r.SNIs = append(r.SNIs, other.SNIs...)

	if len(other.Headers) == 0 {
		return r
	}

	if r.Headers == nil {
		r.Headers = make(map[string]MatchRuleHeader, len(other.Headers))
	}

	for key, headerMatch := range other.Headers {
		_, ok := r.Headers[key]
		if !ok || overrideHeaders {
			r.Headers[key] = headerMatch
		}
	}
	return r
}

func (r *MatchRules) IsEmpty() bool {
	if r == nil {
		return true
	}
	return len(r.Protocols) == 0 &&
		len(r.Hosts) == 0 &&
		len(r.Methods) == 0 &&
		len(r.Paths) == 0 &&
		len(r.Headers) == 0 &&
		len(r.SNIs) == 0
}

func (r *MatchRules) Expression() string {
	fieldMatchers := []Matcher{}

	if len(r.Protocols) > 0 {
		predicates := []Matcher{}
		for _, protocol := range r.Protocols {
			predicates = append(predicates, NewPredicateNetProtocol(OpEqual, protocol))
		}
		fieldMatchers = append(fieldMatchers, Or(predicates...))
	}

	if len(r.Methods) > 0 {
		predicates := []Matcher{}
		for _, method := range r.Methods {
			predicates = append(predicates, NewPredicateHTTPMethod(OpEqual, method))
		}
		fieldMatchers = append(fieldMatchers, Or(predicates...))
	}

	if len(r.Hosts) > 0 {
		predicates := []Matcher{}
		for _, hostMatch := range r.Hosts {
			switch hostMatch.Type {
			case HostMatchExact:
				predicates = append(predicates, NewPrediacteHTTPHost(OpEqual, hostMatch.Host))
			case HostMatchWildcard:
				predicates = append(predicates, NewPrediacteHTTPHost(OpSuffixMatch, hostMatch.Host))
			}
		}
		fieldMatchers = append(fieldMatchers, Or(predicates...))
	}

	if len(r.Paths) > 0 {
		predicates := []Matcher{}
		for _, pathMatch := range r.Paths {
			switch pathMatch.Type {
			case PathMatchExact:
				predicates = append(predicates, NewPredicateHTTPPath(OpEqual, pathMatch.Path))
			case PathMatchPrefix:
				predicates = append(predicates, NewPredicateHTTPPath(OpPrefixMatch, pathMatch.Path))
			case PathMatchRegex:
				predicates = append(predicates, NewPredicateHTTPPath(OpRegexMatch, pathMatch.Path))
			}
		}
		fieldMatchers = append(fieldMatchers, Or(predicates...))
	}

	for key, headerMatch := range r.Headers {
		switch headerMatch.Type {
		case HeaderMatchExact:
			predicates := []Matcher{}
			for _, v := range headerMatch.Values {
				predicates = append(predicates, NewPredicateHTTPHeader(key, OpEqual, v))
			}
			fieldMatchers = append(fieldMatchers, Or(predicates...))
		case HeaderMatchRegex:
			if len(headerMatch.Values) == 1 {
				v := headerMatch.Values[0]
				fieldMatchers = append(fieldMatchers, NewPredicateHTTPHeader(key, OpRegexMatch, v))
			}
		}
	}

	if len(r.SNIs) > 0 {
		predicates := []Matcher{}
		for _, sni := range r.SNIs {
			predicates = append(predicates, NewPredicateTLSSNI(OpEqual, sni))
		}
		fieldMatchers = append(fieldMatchers, Or(predicates...))
	}

	return And(fieldMatchers...).Expression()
}
