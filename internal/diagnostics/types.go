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

// ClientDiagnostic contains settings and channels for receiving diagnostic data from the controller's Kong client.
type ClientDiagnostic struct {
	// DumpsIncludeSensitive is true if the configuration dump includes sensitive values, such as certificate private
	// keys and credential secrets.
	DumpsIncludeSensitive bool

	// Configs is the channel that receives configuration blobs from the configuration update strategy implementation.
	Configs chan ConfigDump

	// FallbackCacheMetadata is the channel that receives fallback metadata from the fallback cache generator.
	FallbackCacheMetadata chan fallback.GeneratedCacheMetadata

	// Diffs is the channel that receives diff info in DB mode.
	Diffs chan ConfigDiff
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
