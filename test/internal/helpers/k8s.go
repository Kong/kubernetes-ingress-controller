package helpers

import (
	"strings"
	"testing"
)

// LabelValueForTest returns a sanitized test name that can be used as kubernetes
// label value.
func LabelValueForTest(t *testing.T) string {
	s := strings.ReplaceAll(t.Name(), "/", ".")
	// Trim to adhere to k8s label requirements:
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set
	if len(s) > 63 {
		return s[:63]
	}
	return s
}
