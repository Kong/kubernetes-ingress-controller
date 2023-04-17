package versions

import (
	"sync"

	"github.com/blang/semver/v4"
)

var (
	// RegexHeaderVersionCutoff is the Kong version prior to the addition of support for regular expression heade
	// matches.
	RegexHeaderVersionCutoff = semver.Version{Major: 2, Minor: 8}

	// ExplicitRegexPathVersionCutoff is the Kong version prior to the addition of explicit "~" prefixes in regular
	// expression paths.
	ExplicitRegexPathVersionCutoff = semver.Version{Major: 3}

	// PluginOrderingVersionCutoff is the Kong version prior to the addition of plugin ordering.
	PluginOrderingVersionCutoff = semver.Version{Major: 3}

	// MTLSCredentialVersionCutoff is the minimum Kong version that support mTLS credentials. This is a patch version
	// because the original version of the mTLS credential was not compatible with KIC.
	MTLSCredentialVersionCutoff = semver.Version{Major: 2, Minor: 3, Patch: 2}

	// FlattenedErrorCutoff is the Kong version prior to the addition of flattened errors.
	FlattenedErrorCutoff = semver.Version{Major: 3, Minor: 1}
)

var (
	// kongVersion holds the Kong version singleton. If never initialized (during some tests), it defaults to the
	// lowest possible version.
	kongVersion = KongVersion(semver.MustParse("0.0.0"))

	kongVersionOnce sync.Once

	kongVersionLock sync.RWMutex
)

// KongVersion is a Kong version.
type KongVersion semver.Version

// SetKongVersion sets the Kong version. It can only be used once. Repeated calls will not update the Kong
// version.
func SetKongVersion(version semver.Version) {
	kongVersionOnce.Do(func() {
		kongVersionLock.Lock()
		defer kongVersionLock.Unlock()
		kongVersion = KongVersion(version)
	})
}

// GetKongVersion retrieves the Kong version. If the version is not set, it returns the lowest possible version.
func GetKongVersion() KongVersion {
	kongVersionLock.RLock()
	defer kongVersionLock.RUnlock()
	return kongVersion
}

// Full returns a complete Kong version as a semver.Version.
func (v KongVersion) Full() semver.Version {
	return semver.Version(v)
}

// MajorOnly returns a semver.Version with a KongVersion's major version only.
func (v KongVersion) MajorOnly() semver.Version {
	return semver.Version{Major: v.Major}
}

// MajorMinorOnly returns a semver.Version with a KongVersion's major and minor versions only.
func (v KongVersion) MajorMinorOnly() semver.Version {
	return semver.Version{Major: v.Major, Minor: v.Minor}
}

// MajorMinorPatchOnly returns a semver.Version with a KongVersion's major, minor, and patch versions only.
func (v KongVersion) MajorMinorPatchOnly() semver.Version {
	return semver.Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch}
}
