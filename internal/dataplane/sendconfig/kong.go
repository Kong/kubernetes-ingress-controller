package sendconfig

import (
	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Types
// -----------------------------------------------------------------------------

// Kong Represents a Kong client and connection information.
type Kong struct {
	URL string
	// If the gateway instance does not support tags, pass an empty FilterTags slice instead.
	FilterTags []string
	// Headers are injected into every request to Kong's Admin API
	// to help with authorization/authentication.
	Client            *kong.Client
	PluginSchemaStore *util.PluginSchemaStore
	InMemory          bool
	Version           semver.Version
	Concurrency       int
}
