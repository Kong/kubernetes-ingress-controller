package multiinstance

import "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"

// InstanceWithIDAlreadyScheduledError is an error indicating that an instance with the same ID is already scheduled.
type InstanceWithIDAlreadyScheduledError struct {
	id manager.ID
}

func NewInstanceWithIDAlreadyScheduledError(id manager.ID) InstanceWithIDAlreadyScheduledError {
	return InstanceWithIDAlreadyScheduledError{id: id}
}

func (e InstanceWithIDAlreadyScheduledError) Error() string {
	return "instance with ID " + e.id.String() + " already exists"
}

// InstanceNotFoundError is an error indicating that an instance with the given ID was not found in the manager.
// It can indicate that the instance was never scheduled or was stopped.
type InstanceNotFoundError struct {
	id manager.ID
}

func NewInstanceNotFoundError(id manager.ID) InstanceNotFoundError {
	return InstanceNotFoundError{id: id}
}

func (e InstanceNotFoundError) Error() string {
	return "instance with ID " + e.id.String() + " not found"
}
