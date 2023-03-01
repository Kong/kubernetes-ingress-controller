//go:build third_party
// +build third_party

package third_party

import (
	_ "sigs.k8s.io/controller-runtime/tools/setup-envtest"
)

//go:generate go install -modfile go.mod sigs.k8s.io/controller-runtime/tools/setup-envtest
