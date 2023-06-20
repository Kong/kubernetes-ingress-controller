//go:build third_party

package skaffold

import (
	_ "github.com/GoogleContainerTools/skaffold/v2/cmd/skaffold"
)

//go:generate go install -modfile go.mod github.com/GoogleContainerTools/skaffold/v2/cmd/skaffold
