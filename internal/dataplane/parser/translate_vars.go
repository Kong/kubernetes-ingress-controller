package parser

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
