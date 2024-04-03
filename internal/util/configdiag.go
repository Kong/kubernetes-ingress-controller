package util

import "github.com/kong/go-database-reconciler/pkg/file"

// ConfigDump contains a config dump and a flag indicating that the config was not successfully applid.
type ConfigDump struct {
	// Config is the configuration KIC applied or attempted to apply.
	Config file.Content
	// Failed is true if the configuration apply failed.
	Failed bool
	// RawResponseBody is the raw Kong Admin API response body from a config apply. It is only available in DB-less mode.
	RawResponseBody []byte
}

// ConfigDumpDiagnostic contains settings and channels for receiving diagnostic configuration dumps.
type ConfigDumpDiagnostic struct {
	// DumpsIncludeSensitive is true if the configuration dump includes sensitive values, such as certificate private
	// keys and credential secrets.
	DumpsIncludeSensitive bool
	// Configs is the channel that receives configuration blobs from the configuration update strategy implementation.
	Configs chan ConfigDump
}
