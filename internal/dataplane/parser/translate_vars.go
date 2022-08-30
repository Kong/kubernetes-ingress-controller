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
)

var (
	// MinRegexHeaderKongVersion is the minimum Kong version that supports regex header matches.
	MinRegexHeaderKongVersion = semver.MustParse("2.8.0")

	// MaxHeuristicRegexPathDetectionVersion is the maximum Kong (major) version that detects regular expression paths
	// automatically using a heuristic.
	MaxHeuristicRegexPathDetectionVersion = semver.Version{Major: 2}

	// PluginOrderingVersionCutoff is the Kong version prior to the addition of plugin ordering. Any Kong version <=
	// PluginOrderingVersionCutoff does not support plugin ordering.
	PluginOrderingVersionCutoff = semver.Version{Major: 2}

	// LegacyRegexPathExpression is the regular expression used by Kong <3.0 to determine if a path is a regex
	LegacyRegexPathExpression = regexp.MustCompile(`^\^?[a-zA-Z0-9\.\-_~/%]*\$?$`)
)
