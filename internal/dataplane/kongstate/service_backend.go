package kongstate

import (
	"errors"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ServiceBackendType is the type of the backend.
type ServiceBackendType string

const (
	// ServiceBackendTypeKongServiceFacade means that the backend is an incubatorv1alpha1.KongServiceFacade.
	ServiceBackendTypeKongServiceFacade ServiceBackendType = "KongServiceFacade"

	// ServiceBackendTypeKubernetesService means that the backend is a Kubernetes Service.
	ServiceBackendTypeKubernetesService ServiceBackendType = "KubernetesService"
)

type ServiceBackends []ServiceBackend

// ServiceBackend represents a backend for a Kong Service. It can be a Kubernetes Service or a KongServiceFacade.
type ServiceBackend struct {
	backendType ServiceBackendType
	name        string
	namespace   string
	portDef     PortDef
	weight      *int
}

// NewServiceBackend creates a new ServiceBackend with an arbitrary backend type.
func NewServiceBackend(
	t ServiceBackendType,
	namespace string,
	name string,
	portDef PortDef,
) (ServiceBackend, error) {
	if t == "" {
		return ServiceBackend{}, errors.New("backend type cannot be empty")
	}
	if namespace == "" {
		return ServiceBackend{}, errors.New("namespace cannot be empty")
	}
	if name == "" {
		return ServiceBackend{}, errors.New("name cannot be empty")
	}
	return ServiceBackend{
		backendType: t,
		namespace:   namespace,
		name:        name,
		portDef:     portDef,
	}, nil
}

// NewServiceBackendForService creates a new ServiceBackend for a Kubernetes Service.
func NewServiceBackendForService(namespace, name string, portDef PortDef) (ServiceBackend, error) {
	return NewServiceBackend(
		ServiceBackendTypeKubernetesService,
		namespace,
		name,
		portDef,
	)
}

// NewServiceBackendForServiceFacade creates a new ServiceBackend for a KongServiceFacade.
func NewServiceBackendForServiceFacade(namespace, name string, portDef PortDef) (ServiceBackend, error) {
	return NewServiceBackend(
		ServiceBackendTypeKongServiceFacade,
		namespace,
		name,
		portDef,
	)
}

// SetWeight sets the weight of the backend used for load-balancing.
func (s *ServiceBackend) SetWeight(weight int32) {
	s.weight = lo.ToPtr(int(weight))
}

// Name returns the name of the backend resource (Service or KongServiceFacade).
func (s *ServiceBackend) Name() string {
	return s.name
}

// Namespace returns the namespace of the backend resource (Service or KongServiceFacade).
func (s *ServiceBackend) Namespace() string {
	return s.namespace
}

// PortDef returns the port definition of the backend.
func (s *ServiceBackend) PortDef() PortDef {
	return s.portDef
}

// Weight returns the weight of the backend used for load-balancing.
func (s *ServiceBackend) Weight() mo.Option[int] {
	if s.weight != nil {
		return mo.Some(*s.weight)
	}
	return mo.None[int]()
}

// IsServiceFacade returns true if the backend is a KongServiceFacade. Otherwise, returns false
// what means that the backend is a Kubernetes Service.
func (s *ServiceBackend) IsServiceFacade() bool {
	return s.backendType == ServiceBackendTypeKongServiceFacade
}
