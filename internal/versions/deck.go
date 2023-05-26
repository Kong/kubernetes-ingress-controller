package versions

import "github.com/blang/semver/v4"

// DeckFileFormat returns Deck file format based on Kong version.
func DeckFileFormat(kongVersion semver.Version) string {
	if kongVersion.GTE(ExplicitRegexPathVersionCutoff) {
		return "3.0"
	}
	return "1.1"
}
