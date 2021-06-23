package v1beta1

// UDPIngressRule represents a rule to apply against incoming requests
// wherein no Host matching is available for request routing, only the port
// is used to match requests.
type UDPIngressRule struct {
	// Port indicates the port for the Kong proxy to accept incoming traffic
	// on, which will then be routed to the service Backend.
	Port int `json:"port"`

	// Backend defines the Kubernetes service which accepts traffic from the
	// listening Port defined above.
	Backend IngressBackend `json:"backend"`
}

// IngressRule represents a rule to apply against incoming requests.
// Matching is performed based on an (optional) SNI and port.
type IngressRule struct {
	// Host is the fully qualified domain name of a network host, as defined
	// by RFC 3986.
	// If a Host is specified, the protocol must be TLS over TCP.
	// A plain-text TCP request cannot be routed based on Host. It can only
	// be routed based on Port.
	// +optional
	Host string `json:"host,omitempty"`

	// Port is the port on which to accept TCP or TLS over TCP sessions and
	// route. It is a required field. If a Host is not specified, the requested
	// are routed based only on Port.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port,omitempty"`

	// Backend defines the referenced service endpoint to which the traffic
	// will be forwarded to.
	Backend IngressBackend `json:"backend"`
}

// IngressBackend describes all endpoints for a given service and port.
type IngressBackend struct {
	// Specifies the name of the referenced service.
	ServiceName string `json:"serviceName"`

	// Specifies the port of the referenced service.
	ServicePort int `json:"servicePort"`
}
