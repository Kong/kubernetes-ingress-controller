package sendconfig

import (
	"github.com/blang/semver/v4"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// Kong Represents a Kong client and connection information
type Kong struct {
	URL        string
	FilterTags []string
	// Headers are injected into every request to Kong's Admin API
	// to help with authorization/authentication.
	Client            *kong.Client
	PluginSchemaStore *util.PluginSchemaStore

	InMemory bool
	// DeprecatedHasTagSupport is not used in KIC 2.x.
	// If the gateway instance does not support tags, pass an empty FilterTags slice instead.
	DeprecatedHasTagSupport bool
	Enterprise              bool

	Version semver.Version

	Concurrency int

	// configuration update
	ConfigDone chan file.Content
}
