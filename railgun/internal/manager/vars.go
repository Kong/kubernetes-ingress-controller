package manager

import "github.com/kong/kubernetes-ingress-controller/pkg/util"

// -----------------------------------------------------------------------------
// Controller Manager - Versioning Information
// -----------------------------------------------------------------------------

var (
	// Release returns the release version
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "../Dockerfile.railgun" for invocation via "go build".
	Release = "UNKNOWN"

	// Repo returns the git repository URL
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "../Dockerfile.railgun" for invocation via "go build".
	Repo = "UNKNOWN"

	// Commit returns the short sha from git
	// NOTE: the value of this is set at compile time using the -X flag for go tool link.
	//       See: "go doc cmd/link" for details, and "../Dockerfile.railgun" for invocation via "go build".
	Commit = "UNKNOWN"
)

// -----------------------------------------------------------------------------
// Controller Manager - Configuration Vars
// -----------------------------------------------------------------------------

// alwaysEnabled is a convenience alias to indicate that a component is not configurable
// and is always enabled in any configuration.
var alwaysEnabled = util.EnablementStatusEnabled
