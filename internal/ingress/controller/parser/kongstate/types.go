package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type PortMode int

const (
	// PortModeImplicit means that the Ingress does not specify the Kubernetes Service port, and that KIC should expect
	// the Service to have only one port defined.
	PortModeImplicit PortMode = iota
	// PortModeNumber means that the Ingress specifies the Service port by raw port number.
	PortModeByNumber PortMode = iota
	// PortModeNumber means that the Ingress specifies the Service port by its name field.
	PortModeByName PortMode = iota
)

type PortDef struct {
	Mode PortMode

	// Name is the port name as stated in the Kubernetes service. Must be set iff Mode == PortModeName.
	Name string
	// Number is the port number. Must be set iff PortMode == PortModeNumber.
	Number int32
}

func PortDefFromIntStr(is intstr.IntOrString) PortDef {
	if is.Type == intstr.String {
		return PortDef{Mode: PortModeByName, Name: is.StrVal}
	}
	return PortDef{Mode: PortModeByNumber, Number: is.IntVal}
}

func (p *PortDef) CanonicalString() string {
	switch p.Mode {
	case PortModeByNumber:
		return fmt.Sprintf("%d", p.Number)
	case PortModeByName:
		return p.Name
	}
	return "implicitPort"
}

type ServiceBackend struct {
	Name string
	Port PortDef
}

// Target is a wrapper around Target object in Kong.
type Target struct {
	kong.Target
}

// Certificate represents the certificate object in Kong.
type Certificate struct {
	kong.Certificate
}

// Plugin represetns a plugin Object in Kong.
type Plugin struct {
	kong.Plugin
}
