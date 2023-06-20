//go:build third_party

package third_party

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

//go:generate go install -modfile go.mod github.com/golangci/golangci-lint/cmd/golangci-lint
