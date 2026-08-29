package diagnostics

import (
	"github.com/kong/go-database-reconciler/pkg/file"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
)

// DumpMeta annotates a config dump.
type DumpMeta struct {
	// Failed indicates the dump was not accepted by the Kong admin API.
	Failed bool
	// Fallback indicates that the dump is a fallback configuration attempted after a failed config update.
	Fallback bool
	// Hash is the configuration hash.
	Hash string
}

// ConfigDump contains a config dump and a flag indicating that the config was not successfully applid.
type ConfigDump struct {
	// Config is the configuration KIC applied or attempted to apply.
	Config file.Content
	// Meta contains information about the status and context of the configuration dump.
	Meta DumpMeta
	// RawResponseBody is the raw Kong Admin API response body from a config apply. It is only available in DB-less mode.
	RawResponseBody []byte
}

// Client contains settings and channels for receiving diagnostic data from the controller's Kong client.
// It encapsulates the channels and exposes methods for sending data to ensure controlled access.
type Client struct {
	// dumpsIncludeSensitive is true if the configuration dump includes sensitive values, such as certificate private
	// keys and credential secrets.
	dumpsIncludeSensitive bool

	// configs is the channel that receives configuration blobs from the configuration update strategy implementation.
	configs chan ConfigDump

	// fallbackCacheMetadata is the channel that receives fallback metadata from the fallback cache generator.
	fallbackCacheMetadata chan fallback.GeneratedCacheMetadata

	// diffs is the channel that receives diff info in DB mode.
	diffs chan ConfigDiff
}

// NewClient creates a new Client with the given configuration.
func NewClient(dumpsIncludeSensitive bool, configBufferSize, fallbackBufferSize, diffsBufferSize int) Client {
	return Client{
		dumpsIncludeSensitive: dumpsIncludeSensitive,
		configs:               make(chan ConfigDump, configBufferSize),
		fallbackCacheMetadata: make(chan fallback.GeneratedCacheMetadata, fallbackBufferSize),
		diffs:                 make(chan ConfigDiff, diffsBufferSize),
	}
}

// DumpsIncludeSensitive returns whether the configuration dump includes sensitive values.
func (c Client) DumpsIncludeSensitive() bool {
	return c.dumpsIncludeSensitive
}

// SendConfig sends a configuration dump to the diagnostics channel.
// It returns false if the channel buffer is full and the send would block.
func (c Client) SendConfig(dump ConfigDump) bool {
	select {
	case c.configs <- dump:
		return true
	default:
		return false
	}
}

// SendFallbackCacheMetadata sends fallback cache metadata to the diagnostics channel.
// It returns false if the channel buffer is full and the send would block.
func (c Client) SendFallbackCacheMetadata(meta fallback.GeneratedCacheMetadata) bool {
	select {
	case c.fallbackCacheMetadata <- meta:
		return true
	default:
		return false
	}
}

// SendDiff sends a configuration diff to the diagnostics channel.
// It returns false if the channel buffer is full and the send would block.
func (c Client) SendDiff(diff ConfigDiff) bool {
	select {
	case c.diffs <- diff:
		return true
	default:
		return false
	}
}

// IsEmpty returns true if the Client has not been initialized (zero value).
func (c Client) IsEmpty() bool {
	return c.configs == nil && c.fallbackCacheMetadata == nil && c.diffs == nil
}

// Configs returns a receive-only channel for reading configuration dumps.
// This is primarily useful for testing and the diagnostics collector.
func (c Client) Configs() <-chan ConfigDump {
	return c.configs
}

// FallbackCacheMetadataCh returns a receive-only channel for reading fallback cache metadata.
// This is primarily useful for testing and the diagnostics collector.
func (c Client) FallbackCacheMetadataCh() <-chan fallback.GeneratedCacheMetadata {
	return c.fallbackCacheMetadata
}

// DiffsCh returns a receive-only channel for reading configuration diffs.
// This is primarily useful for testing and the diagnostics collector.
func (c Client) DiffsCh() <-chan ConfigDiff {
	return c.diffs
}

// AffectedObject is a Kubernetes object associated with diagnostic information.
type AffectedObject struct {
	// UID is the unique identifier of the object.
	UID k8stypes.UID

	// Group is the object's group.
	Group string
	// Kind is the object's Kind.
	Kind string
	// Namespace is the object's Namespace.
	Namespace string
	// Name is the object's Name.
	Name string
}
