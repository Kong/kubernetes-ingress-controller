//go:build third_party
// +build third_party

package third_party

import (
	_ "k8s.io/code-generator/cmd/client-gen"
)

//go:generate go install -modfile go.mod k8s.io/code-generator/cmd/client-gen
