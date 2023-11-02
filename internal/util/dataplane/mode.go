package dataplane

// IsDBLessMode can be used to detect the proxy mode (db or dbless).
func IsDBLessMode(mode string) bool {
	return mode == "" || mode == "off"
}

// DBBacked returns true if the gateway is DB backed.
// reverse of IsDBLessMode for readability.
func DBBacked(mode string) bool {
	return !IsDBLessMode(mode)
}
