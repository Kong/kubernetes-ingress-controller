package manager

import (
	"github.com/google/uuid"
)

// InstanceIDProvider provides a unique identifier for a running instance of the manager.
// It should be used by all components that register the instance in any external system.
type InstanceIDProvider struct {
	id uuid.UUID
}

func NewInstanceIDProvider() *InstanceIDProvider {
	return &InstanceIDProvider{
		id: uuid.New(),
	}
}

func (p *InstanceIDProvider) GetID() uuid.UUID {
	return p.id
}
