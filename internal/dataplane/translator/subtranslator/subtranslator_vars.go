package subtranslator

const (
	// KongPathRegexPrefix is the reserved prefix string that instructs Kong 3.0+ to interpret a path as a regex.
	KongPathRegexPrefix = "~"

	// KongHeaderRegexPrefix is a reserved prefix string that Kong uses to determine if it should parse a header value
	// as a regex.
	KongHeaderRegexPrefix = "~*"

	// ControllerPathRegexPrefix is the prefix string used to indicate that the controller should treat a path as a
	// regular expression. The controller replaces this prefix with KongPathRegexPrefix when sending routes to Kong.
	ControllerPathRegexPrefix = "/~"

	// DefaultServiceTimeout indicates the amount of time (by default) for
	// connections, reads and writes to a service over a network should
	// be given before timing out by default.
	DefaultServiceTimeout = 60000

	// DefaultRetries indicates the number of times a connection should be
	// retried by default.
	DefaultRetries = 5

	// DefualtKongServiceProtocol is the default protocol in translated Kong service.
	DefualtKongServiceProtocol = "http"

	// maxKongServiceNameLength is the maximum length of generated Kong service name.
	// if the length of generated Kong service name exceeds the limit, the name will be trimmed.
	maxKongServiceNameLength = (512 - 1)
)
