package util

const (
	// minPort is the minimum networking port number.
	minPort = 1

	// maxPort is the maximum networking port number.
	maxPort = 65535
)

// IsValidPort is a convenience function to validate whether or not
// a given integer is a valid networking port number.
func IsValidPort(port int) bool {
	if port >= minPort && port <= maxPort {
		return true
	}
	return false
}
