package manager

import (
	"errors"

	"github.com/google/uuid"
)

// ID is a unique identifier for the Kong Ingress Controller instance.
// It can be an arbitrary string that is unique across all instances.
type ID struct {
	id string
}

// NewRandomID generates a new random manager ID.
func NewRandomID() ID {
	return ID{id: uuid.NewString()}
}

// NewID creates a new manager ID from a string (e.g. a Kubernetes object UID).
func NewID(s string) (ID, error) {
	if s == "" {
		return ID{}, errors.New("manager ID cannot be empty")
	}
	return ID{id: s}, nil
}

func (id ID) String() string {
	return id.id
}
