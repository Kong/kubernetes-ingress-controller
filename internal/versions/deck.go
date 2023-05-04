package versions

func DeckFileFormat(kongVersion KongVersion) string {
	if kongVersion.MajorMinorOnly().GTE(ExplicitRegexPathVersionCutoff) {
		return "3.0"
	}
	return "1.1"
}
