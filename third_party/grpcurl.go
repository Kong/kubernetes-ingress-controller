//go:build third_party
// +build third_party

package third_party

import (
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
)

//go:generate go install -modfile go.mod github.com/fullstorydev/grpcurl/cmd/grpcurl
