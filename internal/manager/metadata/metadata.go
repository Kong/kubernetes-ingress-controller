// Package metadata includes metadata variables for logging and reporting.
package metadata

// -----------------------------------------------------------------------------
// Controller Manager - Versioning Information
// -----------------------------------------------------------------------------

// WARNING: moving any of these variables requires changes to both the Makefile
//          and the Dockerfile which modify them during the link step with -X

var (
	// Release returns the release version, generally a semver like v2.0.0.
	Release = "NOT_SET"

	// Repo returns the git repository URL like git@github.com:Kong/kubernetes-ingress-controller.git.
	Repo = "NOT_SET"

	// ProjectURL returns address of project website - git repository like github.com/kong/kubernetes-ingress-controller.
	ProjectURL = "NOT_SET"

	// Commit returns the SHA from the current branch HEAD.
	Commit = "NOT_SET"
)

const (
	Organization = "Kong"
	ProjectName  = "kubernetes-ingress-controller"
)
