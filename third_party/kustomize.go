//go:build third_party

package third_party

import (
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)

//go:generate go install -modfile go.mod sigs.k8s.io/kustomize/kustomize/v5
