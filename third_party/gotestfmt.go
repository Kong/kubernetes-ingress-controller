//go:build third_party
// +build third_party

package third_party

import (
	_ "github.com/haveyoudebuggedit/gotestfmt/v2"
)

//go:generate go install -modfile go.mod github.com/haveyoudebuggedit/gotestfmt/v2
