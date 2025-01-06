//go:build integration_tests

package isolated

import "fmt"

func examplesManifestPath(manifestName string) string {
	return fmt.Sprintf("../../../examples/%s", manifestName)
}
