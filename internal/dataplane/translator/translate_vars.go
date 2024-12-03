package translator

import (
	"regexp"
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
)

// LegacyRegexPathExpression is the regular expression used by Kong <3.0 to determine if a path is not a regex.
var LegacyRegexPathExpression = regexp.MustCompile(`^[a-zA-Z0-9\.\-_~/%]*$`)
