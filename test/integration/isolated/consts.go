//go:build integration_tests

package isolated

import "fmt"

func manifestPath(manifestName string) string {
	return fmt.Sprintf("../../../examples/%s", manifestName)
}
