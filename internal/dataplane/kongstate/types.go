package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PortMode int

const (
	// PortModeImplicit means that the Ingress does not specify the Kubernetes Service port, and that KIC should expect
	// the Service to have only one port defined.
	PortModeImplicit PortMode = iota
	// PortModeByNumber means that the Ingress specifies the Service port by raw port number.
	PortModeByNumber PortMode = iota
	// PortModeByName means that the Ingress specifies the Service port by its name field.
	PortModeByName PortMode = iota
)

type PortDef struct {
	Mode PortMode

	// Name is the port name as stated in the Kubernetes service. Must be set iff Mode == PortModeName.
	Name string

	// Number is the port number. Must be set iff PortMode == PortModeNumber.
	Number int32
}

const ImplicitPort = "implicitPort"

func (p *PortDef) CanonicalString() string {
	switch p.Mode {
	case PortModeByNumber:
		return fmt.Sprintf("%d", p.Number)
	case PortModeByName:
		return p.Name
	case PortModeImplicit:
		return ImplicitPort
	}
	return ImplicitPort
}

// ServiceBackendType is the type of the backend.
type ServiceBackendType string

const (
	// ServiceBackendTypeKongServiceFacade means that the backend is an incubatorv1alpha1.KongServiceFacade.
	ServiceBackendTypeKongServiceFacade ServiceBackendType = "KongServiceFacade"
)

type ServiceBackend struct {
	// Type is the type of the backend. If left empty, it is assumed to be a Kubernetes Service.
	Type ServiceBackendType

	// Name is the name of the backend resource (Service or KongServiceFacade).
	Name string

	// Namespace is the namespace of the backend resource (Service or KongServiceFacade).
	Namespace string

	// PortDef is the port definition of the backend.
	PortDef PortDef

	// Weight is the weight of the backend used for load-balancing.
	Weight *int32
}

type ServiceBackends []ServiceBackend

// Target is a wrapper around Target object in Kong.
type Target struct {
	kong.Target
}

// Certificate represents the certificate object in Kong.
type Certificate struct {
	kong.Certificate
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *Certificate) SanitizedCopy() *Certificate {
	return &Certificate{
		kong.Certificate{
			ID:        c.ID,
			Cert:      c.Cert,
			Key:       redactedString,
			CreatedAt: c.CreatedAt,
			SNIs:      c.SNIs,
			Tags:      c.Tags,
		},
	}
}

// Plugin represents a plugin Object in Kong.
type Plugin struct {
	kong.Plugin
	K8sParent client.Object
}

func (p Plugin) DeepCopy() Plugin {
	return Plugin{
		Plugin:    *p.Plugin.DeepCopy(),
		K8sParent: p.K8sParent,
	}
}
