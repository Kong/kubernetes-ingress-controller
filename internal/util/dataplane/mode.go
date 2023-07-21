package dataplane

// IsDBLessMode can be used to detect the proxy mode (db or dbless).
func IsDBLessMode(mode string) bool {
	return mode == "" || mode == "off"
}
