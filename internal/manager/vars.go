package manager

// -----------------------------------------------------------------------------
// Controller Manager - Versioning Information
// -----------------------------------------------------------------------------

var (
	// Release returns the release version
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "Dockerfile" for invocation via "go build".
	Release = "UNKNOWN"

	// Repo returns the git repository URL
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "Dockerfile" for invocation via "go build".
	Repo = "UNKNOWN"

	// Commit returns the short sha from git
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "Dockerfile" for invocation via "go build".
	Commit = "UNKNOWN"
)

// -----------------------------------------------------------------------------
// Controller Manager - Configuration Vars
// -----------------------------------------------------------------------------

// neverDisabled is a convenience alias to indicate that a controller cannot be disabled
var neverDisabled = false
