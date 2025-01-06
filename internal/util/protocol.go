package util

// ValidateProtocol returns true if the provided protocol is valid.
func ValidateProtocol(protocol string) bool {
	switch protocol {
	case "", "http", "https", "grpc", "grpcs", "ws", "wss", "tls", "tcp", "tls_passthrough":
		return true
	default:
		return false
	}
}
