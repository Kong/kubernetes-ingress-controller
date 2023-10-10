package versions

import (
	"github.com/blang/semver/v4"
)

var (
	// KICv3VersionCutoff is the lowest version version of Kong Gateway supported by KIC >=v3.0.0.
	KICv3VersionCutoff = semver.Version{Major: 3, Minor: 4, Patch: 1}

	// ExplicitRegexPathVersionCutoff is the lowest Kong version requiring the explicit "~" prefixes in regular expression paths.
	ExplicitRegexPathVersionCutoff = semver.Version{Major: 3, Minor: 0}
)
