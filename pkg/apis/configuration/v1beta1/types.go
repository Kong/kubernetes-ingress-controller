package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TCPIngress is very similar to (and heavily borrows from) Ingress resource
// in the networking.v1beta1 group but for TCP or L4 routing.
// TCPIngress is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TCPIngress struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the TCPIngress.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec IngressSpec `json:"spec,omitempty"`

	// Status is the current state of the Ingress.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Status IngressStatus `json:"status,omitempty"`
}

// TCPIngressList is a top-level list type. The client methods for
// lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TCPIngressList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []TCPIngress `json:"items"`
}

// IngressSpec describes the TCPIngress the user wishes to exist.
type IngressSpec struct {
	// A list of rules used to configure the Ingress.
	Rules []IngressRule `json:"rules,omitempty"`
	// TLS configuration. This is similar to the `tls` section in the
	// Ingress resource in networking.v1beta1 group.
	// The mapping of SNIs to TLS cert-key pair defined here will be
	// used for HTTP Ingress rules as well. Once can define the mapping in
	// this resource or the original Ingress resource, both have the same
	// effect.
	// +optional
	TLS []IngressTLS `json:"tls,omitempty"`
}

// IngressTLS describes the transport layer security.
type IngressTLS struct {
	// Hosts are a list of hosts included in the TLS certificate. The values in
	// this list must match the name/s used in the tlsSecret. Defaults to the
	// wildcard host setting for the loadbalancer controller fulfilling this
	// Ingress, if left unspecified.
	// +optional
	Hosts []string `json:"hosts,omitempty"`
	// SecretName is the name of the secret used to terminate SSL traffic.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

// IngressStatus describe the current state of TLSIngress.
type IngressStatus struct {
	// LoadBalancer contains the current status of the load-balancer.
	// +optional
	LoadBalancer corev1.LoadBalancerStatus `json:"loadBalancer,omitempty"`
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
