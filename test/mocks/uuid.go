package mocks

import "github.com/kong/kubernetes-ingress-controller/v3/internal/util"

var _ = util.UUIDGenerator(&StaticUUIDGenerator{})

type StaticUUIDGenerator struct {
	UUID string
}

func (s StaticUUIDGenerator) NewString() string {
	return s.UUID
}
