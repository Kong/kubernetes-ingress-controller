package util

import "github.com/google/uuid"

// RandomName creates a unique name prepended with a predefined prefix.
// It can be used to populate .Name field of Kubernetes objects.
func RandomName(prefix string) string {
	return prefix + uuid.NewString()
}
