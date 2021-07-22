package util

import "github.com/kong/deck/file"

// ConfigDumpDiagnostic contains settings and channels for receiving diagnostic configuration dumps
type ConfigDumpDiagnostic struct {
	DumpsIncludeSensitive bool
	SuccessfulConfigs     chan file.Content
	FailedConfigs         chan file.Content
}
