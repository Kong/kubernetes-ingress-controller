package versions

import (
	"github.com/blang/semver/v4"
)

var (
	// RegexHeaderVersionCutoff is the Kong version prior to the addition of support for regular expression heade
	// matches.
	RegexHeaderVersionCutoff = semver.Version{Major: 2, Minor: 8}

	// ExplicitRegexPathVersionCutoff is the lowest Kong version adding the explicit "~" prefixes in regular expression paths.
	ExplicitRegexPathVersionCutoff = semver.Version{Major: 3, Minor: 0}

	// PluginOrderingVersionCutoff is the Kong version prior to the addition of plugin ordering.
	PluginOrderingVersionCutoff = semver.Version{Major: 3}

	// MTLSCredentialVersionCutoff is the minimum Kong version that support mTLS credentials. This is a patch version
	// because the original version of the mTLS credential was not compatible with KIC.
	MTLSCredentialVersionCutoff = semver.Version{Major: 2, Minor: 3, Patch: 2}

	// FlattenedErrorCutoff is the lowest Kong version with support for flattened errors.
	FlattenedErrorCutoff = semver.Version{Major: 3, Minor: 2}

	// TLSPassthroughCutoff is the lowest Kong version with support for TLS passthrough.
	TLSPassthroughCutoff = semver.Version{Major: 2, Minor: 7}
)
