package diagnostics

import (
	"github.com/kong/go-database-reconciler/pkg/file"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

// DumpMeta annotates a config dump.
type DumpMeta struct {
	// Failed indicates the dump was not accepted by the Kong admin API.
	Failed bool
	// Fallback indicates that the dump is a fallback configuration attempted after a failed config update.
	Fallback bool
	// AffectedObjects are objects excluded from the fallback configuration.
	AffectedObjects []AffectedObject
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

type configDumpResponse struct {
	ConfigHash string       `json:"hash"`
	Config     file.Content `json:"config"`
}

type fallbackResponse struct {
	FallbackObjects []FallbackDiagnostic `json:"objects"`
}

// FallbackDiagnosticCollection is a set of fallback object diagnostics associated with a config hash.
type FallbackDiagnosticCollection struct {
	Objects []FallbackDiagnostic
}

// FallbackDiagnostic are fallback objects.
type FallbackDiagnostic struct {
	// GroupKind is the object group and kind.
	GroupKind string `json:"resource"`
	// Namespace is the object namespace.
	Namespace string `json:"namespace"`
	// Namespace is the object name.
	Name string `json:"name"`
	// ID is the object UID.
	ID string `json:"id"`
	// Status is the object's fallback status.
	Status string `json:"status"`
	// CausingObject is the object that triggered this
	CausingObject string `json:"cause"`
}

// ConfigDumpDiagnostic contains settings and channels for receiving diagnostic configuration dumps.
type ConfigDumpDiagnostic struct {
	// DumpsIncludeSensitive is true if the configuration dump includes sensitive values, such as certificate private
	// keys and credential secrets.
	DumpsIncludeSensitive bool
	// Configs is the channel that receives configuration blobs from the configuration update strategy implementation.
	Configs   chan ConfigDump
	Fallbacks chan FallbackDiagnosticCollection
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
