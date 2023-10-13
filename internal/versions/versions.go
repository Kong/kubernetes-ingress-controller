package versions

import (
	"github.com/blang/semver/v4"
)

// KICv3VersionCutoff is the lowest version version of Kong Gateway supported by KIC >=v3.0.0.
var KICv3VersionCutoff = semver.Version{Major: 3, Minor: 4, Patch: 1}

// DeckFileFormatVersion is the version of the decK file format used by KIC everywhere.
const DeckFileFormatVersion = "3.0"
