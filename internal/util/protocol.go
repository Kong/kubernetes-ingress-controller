package util

import "regexp"

// ValidateProtocol returns a bool of whether string is a valid protocol
func ValidateProtocol(protocol string) bool {
	match := validProtocols.MatchString(protocol)
	return match
}

var validProtocols = regexp.MustCompile(`\Ahttps$|\Ahttp$|\Agrpc$|\Agrpcs|\Atcp|\Atls|\Atls_passthrough$`)
