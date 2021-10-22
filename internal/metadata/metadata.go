// Package metadata includes metadata variables for logging and reporting.
package metadata

// -----------------------------------------------------------------------------
// Controller Manager - Versioning Information
// -----------------------------------------------------------------------------

// WARNING: moving any of these variables requires changes to both the Makefile
//          and the Dockerfile which modify them during the link step with -X

var (
	// Release returns the release version, generally a semver like v2.0.0.
	Release = "UNKNOWN"

	// Repo returns the git repository URL
	Repo = "UNKNOWN"

	// Commit returns the SHA from the current branch HEAD
	Commit = "UNKNOWN"
)
