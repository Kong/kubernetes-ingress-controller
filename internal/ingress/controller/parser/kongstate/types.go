package kongstate

import (
	"github.com/kong/go-kong/kong"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceBackend struct {
	Name string
	Port intstr.IntOrString
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
