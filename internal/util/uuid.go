package util

import "github.com/google/uuid"

// UUIDGenerator is an interface to generate UUIDs.
type UUIDGenerator interface {
	NewString() string
}

// DefaultUUIDGenerator is the default implementation of UUIDGenerator.
type DefaultUUIDGenerator struct{}

func (DefaultUUIDGenerator) NewString() string {
	return uuid.NewString()
}
