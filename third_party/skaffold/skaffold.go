//go:build third_party

// skaffold is extracted to its own package to avoid issues with its stale dependencies.
package skaffold

import (
	_ "github.com/GoogleContainerTools/skaffold/v2/cmd/skaffold"
)

//go:generate go install -modfile go.mod github.com/GoogleContainerTools/skaffold/v2/cmd/skaffold
