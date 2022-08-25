package parser

import (
	"regexp"

	"github.com/blang/semver/v4"
)

// -----------------------------------------------------------------------------
// Translation - Vars & Constants
// -----------------------------------------------------------------------------

const (
	// DefaultServiceTimeout indicates the amount of time (by default) for
	// connections, reads and writes to a service over a network should
	// be given before timing out by default.
	DefaultServiceTimeout = 60000

	// DefaultRetries indicates the number of times a connection should be
	// retried by default.
	DefaultRetries = 5

	// DefaultHTTPPort is the network port that should be assumed by default
	// for HTTP traffic to services.
	DefaultHTTPPort = 80

	// kongHeaderRegexPrefix is a reserved prefix string that Kong uses to determine if it should parse a header value
	// as a regex.
	kongHeaderRegexPrefix = "~*"

	// kongPathRegexPrefix is the reserved prefix string that instructs Kong 3.0+ to interpret a path as a regex.
	kongPathRegexPrefix = "~"
)

var (
	// MinRegexHeaderKongVersion is the minimum Kong version that supports regex header matches.
	MinRegexHeaderKongVersion = semver.MustParse("2.8.0")

	// MinExplicitPathRegexKongVersion is the minimum Kong version that requires explicit indication of regex paths.
	MinExplicitPathRegexKongVersion = semver.MustParse("3.0.0")

	// LegacyRegexPathExpression is the regular expression used by Kong <3.0 to determine if a path is a regex
	LegacyRegexPathExpression = regexp.MustCompile(`^[a-zA-Z0-9\.\-_~/%]*$`)
)
