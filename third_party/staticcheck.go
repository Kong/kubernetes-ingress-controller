//go:build third_party

package third_party

import (
	_ "honnef.co/go/tools/cmd/staticcheck"
)

//go:generate go install -modfile go.mod honnef.co/go/tools/cmd/staticcheck
