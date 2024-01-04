package util

import "regexp"

// ValidateProtocol returns a bool of whether string is a valid protocol.
func ValidateProtocol(protocol string) bool {
	if protocol == "" {
		return true
	}
	match := validProtocols.MatchString(protocol)
	return match
}

var validProtocols = regexp.MustCompile(`\Ahttps$|\Ahttp$|\Agrpc$|\Agrpcs|\Aws$|\Awss|\Atcp|\Atls|\Atls_passthrough$`)
