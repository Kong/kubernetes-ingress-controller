package translators

const (
	// KongPathRegexPrefix is the reserved prefix string that instructs Kong 3.0+ to interpret a path as a regex.
	KongPathRegexPrefix = "~"

	// ControllerPathRegePrefix is the prefix string used to indicate that the controller should treat a path as a
	// regular expression. The controller replaces this prefix with KongPathRegexPrefix when sending routes to Kong.
	ControllerPathRegexPrefix = "/~"
)
