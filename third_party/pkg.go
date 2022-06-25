package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "k8s.io/code-generator/cmd/client-gen"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kustomize/kustomize/v4"
)
