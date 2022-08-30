//go:build third_party
// +build third_party

package third_party

import (
	_ "sigs.k8s.io/kustomize/kustomize/v4"
)

//go:generate go install -modfile go.mod sigs.k8s.io/kustomize/kustomize/v4
