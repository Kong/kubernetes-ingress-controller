//go:build third_party
// +build third_party

package third_party

import (
	_ "github.com/go-task/task/v3/cmd/task"
)

//go:generate go install -modfile go.mod github.com/go-task/task/v3/cmd/task
