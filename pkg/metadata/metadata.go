// Package metadata includes metadata variables for logging and reporting.
package metadata

import (
	"strings"

	sdkkonnectmetadata "github.com/Kong/sdk-konnect-go/pkg/metadata"
)

func init() {
	// NOTE: We do it this way because speakeasy does not provide a way to set the
	// user-agent for the SDK instance.
	sdkkonnectmetadata.SetUserAgent(UserAgent())
}

// -----------------------------------------------------------------------------
// Controller Manager - Versioning Information
// -----------------------------------------------------------------------------

// WARNING: moving any of these variables requires changes to both the Makefile
//          and the Dockerfile which modify them during the link step with -X

var (
	// Release returns the release version, generally a semver like v2.0.0.
	Release = NotSet

	// Repo returns the git repository URL like git@github.com:Kong/kubernetes-ingress-controller.git.
	Repo = NotSet

	// ProjectURL returns address of project website - git repository like github.com/kong/kubernetes-ingress-controller.
	ProjectURL = NotSet

	// Commit returns the SHA from the current branch HEAD.
	Commit = NotSet

	// ProjectName is the name of repository, the last part of the URL. Thus it's basically name of the application.
	ProjectName = projectNameFromRepo(Repo)
)

const (
	Organization = "Kong"
	NotSet       = "NOT_SET"
)

func projectNameFromRepo(repo string) string {
	parts := strings.Split(repo, "/")
	projectName := strings.TrimSpace(strings.TrimSuffix(parts[len(parts)-1], ".git"))
	if projectName == "" {
		return NotSet
	}
	return projectName
}

func UserAgent() string {
	return "kong-ingress-controller/" + Release
}
