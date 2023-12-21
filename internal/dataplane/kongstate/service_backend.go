package kongstate

import (
	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ServiceBackendType is the type of the backend.
type ServiceBackendType string

const (
	// ServiceBackendTypeKongServiceFacade means that the backend is an incubatorv1alpha1.KongServiceFacade.
	ServiceBackendTypeKongServiceFacade ServiceBackendType = "KongServiceFacade"

	// ServiceBackendTypeKubernetesService means that the backend is a Kubernetes Service.
	ServiceBackendTypeKubernetesService ServiceBackendType = "KongService"
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

func NewServiceBackend(
	t ServiceBackendType,
	namespace string,
	name string,
	portDef PortDef,
) ServiceBackend {
	// TODO return error
	return ServiceBackend{
		backendType: t,
		namespace:   namespace,
		name:        name,
		portDef:     portDef,
	}
}

func NewServiceBackendForService(namespace, name string, portDef PortDef) ServiceBackend {
	return ServiceBackend{
		namespace:   namespace,
		name:        name,
		portDef:     portDef,
		backendType: ServiceBackendTypeKubernetesService,
	}
}

func NewServiceBackendForServiceFacade(namespace, name string, portDef PortDef) ServiceBackend {
	return ServiceBackend{
		name:        name,
		namespace:   namespace,
		portDef:     portDef,
		backendType: ServiceBackendTypeKongServiceFacade,
	}
}

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
