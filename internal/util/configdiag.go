package util

import "github.com/kong/go-database-reconciler/pkg/file"

// ConfigDump contains a config dump and a flag indicating that the config was not successfully applid.
type ConfigDump struct {
	Config file.Content
	Failed bool
	Raw    []byte
}

// ConfigDumpDiagnostic contains settings and channels for receiving diagnostic configuration dumps.
type ConfigDumpDiagnostic struct {
	DumpsIncludeSensitive bool
	Configs               chan ConfigDump
}
