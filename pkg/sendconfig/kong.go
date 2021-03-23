package sendconfig

import (
	"github.com/blang/semver"
	"github.com/kong/go-kong/kong"
)

// Kong Represents a Kong client and connection information
type Kong struct {
	URL        string
	FilterTags []string
	// Headers are injected into every request to Kong's Admin API
	// to help with authorization/authentication.
	Client *kong.Client

	InMemory      bool
	HasTagSupport bool
	Enterprise    bool

	Version semver.Version

	Concurrency int
}
