package versions

import (
	"github.com/blang/semver/v4"
)

// KICv3VersionCutoff is the lowest version version of Kong Gateway supported by KIC >=v3.0.0.
var KICv3VersionCutoff = semver.Version{Major: 3, Minor: 4, Patch: 1}

// KongRedirectPluginCutoff is the lowest version of Kong Gateway that supports `redirect` plugin.
var KongRedirectPluginCutoff = semver.Version{Major: 3, Minor: 9, Patch: 0}

// KongStickySessionsCutoff is the lowest version of Kong Gateway that supports `sticky_sessions` loadbalancing algorithm.
var KongStickySessionsCutoff = semver.Version{Major: 3, Minor: 11, Patch: 0}

// DeckFileFormatVersion is the version of the decK file format used by KIC everywhere.
const DeckFileFormatVersion = "3.0"
