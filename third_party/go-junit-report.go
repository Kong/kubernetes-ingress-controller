//go:build third_party
// +build third_party

package third_party

import (
	_ "github.com/jstemmer/go-junit-report/v2"
)

//go:generate go install -modfile go.mod github.com/jstemmer/go-junit-report/v2
