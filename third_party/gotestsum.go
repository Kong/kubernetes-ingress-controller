//go:build third_party
// +build third_party

package third_party

import (
	_ "gotest.tools/gotestsum"
)

//go:generate go install -modfile go.mod gotest.tools/gotestsum
