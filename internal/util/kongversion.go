package util

import (
	"sync"

	"github.com/blang/semver/v4"
)

var (
	kongVersion     = semver.MustParse("0.0.0")
	kongVersionOnce sync.Once
)

// SetKongVersion sets the Kong version. It can only be used once. Repeated calls will not update the Kong
// version
func SetKongVersion(version semver.Version) {
	kongVersionOnce.Do(func() {
		kongVersion = version
	})
}

// GetKongVersion retrieves the Kong version. If the version is not set, it returns the lowest possible version
func GetKongVersion() semver.Version {
	return kongVersion
}
