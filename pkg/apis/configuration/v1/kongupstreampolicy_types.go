package v1

// REVIEW: do we want to create a separate package for Gateway API extensions?

import (
	"github.com/kong/go-kong/kong"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func init() {
	SchemeBuilder.Register(&KongUpstreamPolicy{}, &KongUpstreamPolicyList{})
}

// KongUpstreamPolicy allows attaching Kong Upstream Policies to Gateway API resources:
//   - HTTPRoute,
//   - TCPRoute,
//   - UDPRoute,
//   - GRPCRoute.
//
// It allows configuring algorithm that should be used for load balancing traffic between Kong Upstream's
// Targets. It also allows configuring health checks for Kong Upstream's Targets.
//
// In the case of Gateway API *Routes, the KongUpstreamPolicy will affect all of Kong Upstreams
// associated with the Gateway API *Route.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespace,shortName=kup,categories=kong-ingress-controller
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type KongUpstreamPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the configuration of the Kong upstream.
	Spec KongUpstreamPolicySpec `json:"spec,omitempty"`

	// Status represents the current status of the KongUpstreamPolicy resource.
	Status KongUpstreamPolicyStatus `json:"status,omitempty"`
}

// KongUpstreamPolicyList contains a list of KongUpstreamPolicy.
// +kubebuilder:object:root=true
type KongUpstreamPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongUpstreamPolicy `json:"items"`
}

// KongUpstreamPolicySpec contains the specification for KongUpstreamPolicy.
type KongUpstreamPolicySpec struct {
	// TargetRef identifies an API object to apply policy to (supported: HTTPRoute, TCPRoute, UDPRoute, GRPCRoute).
	TargetRef v1alpha2.PolicyTargetReference `json:"targetRef,omitempty"`

	// Override defines policy configuration that should override policy configuration attached below the targeted
	// resource in the hierarchy.
	Override *KongUpstreamPolicyConfig `json:"override,omitempty"`

	// Default defines default policy configuration for the targeted resource.
	Default *KongUpstreamPolicyConfig `json:"default,omitempty"`
}

// KongUpstreamPolicyConfig contains the configuration parameters for Kong upstream.
type KongUpstreamPolicyConfig struct {
	// Algorithm is the load balancing algorithm to use.
	// Accepted values are: "round-robin", "consistent-hashing", "least-connections", "latency".
	// +kubebuilder:validation:Enum=round-robin;consistent-hashing;least-connections;latency
	Algorithm *string `json:"algorithm,omitempty"`

	// Slots is the number of slots in the load balancer algorithm.
	// If not set, the default value in Kong for the algorithm is used.
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=65536
	Slots *int `json:"slots,omitempty"`

	// HashOn defines how to calculate hash for consistent-hashing load balancing algorithm.
	HashOn *KongUpstreamHash `json:"hash_on,omitempty"`

	// HasOnFallback defines how to calculate hash for consistent-hashing load balancing algorithm if the primary hash
	// function fails.
	HashOnFallback *KongUpstreamHash `json:"hash_on_fallback,omitempty"`

	// Healthchecks defines the health check configurations in Kong.
	// REVIEW: I think we should not depend on go-kong types here.
	Healthchecks *kong.Healthcheck `json:"healthchecks,omitempty"`

	// HostHeader is the hostname to be used as Host header when proxying requests through Kong.
	// REVIEW: this could be achieved with Gateway API HTTPHeaderFilter, do we need that?
	HostHeader *string `json:"host_header,omitempty"`
}

// KongUpstreamHash defines how to calculate hash for consistent-hashing load balancing algorithm.
// One of the fields must be set.
type KongUpstreamHash struct {
	// Header is the name of the header to use as hash input.
	Header *string `json:"header,omitempty"`

	// Cookie is the name of the cookie to use as hash input.
	Cookie *string `json:"cookie,omitempty"`

	// QueryArg is the name of the query argument to use as hash input.
	QueryArg *string `json:"query_arg,omitempty"`

	// URICapture is the name of the URI capture group to use as hash input.
	URICapture *string `json:"uri_capture,omitempty"`
}

// KongUpstreamPolicyStatus represents the current status of the KongUpstreamPolicy resource.
type KongUpstreamPolicyStatus struct {
	// Conditions describe the current conditions of the ACMEServicePolicy.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
