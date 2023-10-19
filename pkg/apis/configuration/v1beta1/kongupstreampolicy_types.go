package v1beta1

import (
	v1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func init() {
	v1.SchemeBuilder.Register(&KongUpstreamPolicy{}, &KongUpstreamPolicyList{})
}

// KongUpstreamPolicy allows configuring algorithm that should be used for load balancing traffic between Kong
// Upstream's Targets. It also allows configuring health checks for Kong Upstream's Targets.
//
// Its configuration is similar to Kong Upstream object (https://docs.konghq.com/gateway/latest/admin-api/#upstream-object),
// and it is applied to Kong Upstream objects created by the controller.
//
// It can be attached to Services and Gateway API *Routes.
//
// When attached to a Gateway API *Route, it will affect all of Kong Upstreams created for
// the Gateway API *Route. If you want to use different Upstream settings for each of the *Route's individual rules,
// you should create a separate *Route for each of the rules.
//
// When attached to a Service, it will affect all Kong Upstreams created for the Service.
//
// When attached to a Service used in a Gateway API *Route rule with multiple BackendRefs, all of its Services must
// be configured with the same KongUpstreamPolicy (effectively two separate KongUpstreamPolicies with the same
// configuration). Otherwise, the controller will ignore the KongUpstreamPolicy.
//
// When attached to a Service used in a Gateway API *Route that has another KongUpstreamPolicy attached to it,
// the controller will *ignore* the KongUpstreamPolicy attached to the Service.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespace,shortName=kup,categories=kong-ingress-controller
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=gateway.networking.k8s.io/policy=direct
type KongUpstreamPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the configuration of the Kong upstream.
	Spec KongUpstreamPolicySpec `json:"spec,omitempty"`

	// Status represents the current status of the KongUpstreamPolicy resource.
	Status PolicyStatus `json:"status,omitempty"`
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
	// TargetRef identifies an API object to apply policy to.
	TargetRef v1alpha2.PolicyTargetReference `json:"targetRef,omitempty"`

	// Upstream defines configuration to be applied to Kong Upstreams associated with TargetRef.
	Upstream KongUpstreamPolicyConfig `json:"config,omitempty"`
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
	// Algorithm must be set to "consistent-hashing" for this field to have effect.
	HashOn *KongUpstreamHash `json:"hash_on,omitempty"`

	// HasOnFallback defines how to calculate hash for consistent-hashing load balancing algorithm if the primary hash
	// function fails.
	// Algorithm must be set to "consistent-hashing" for this field to have effect.
	HashOnFallback *KongUpstreamHash `json:"hash_on_fallback,omitempty"`

	// Healthchecks defines the health check configurations in Kong.
	Healthchecks *KongUpstreamHealthcheck `json:"healthchecks,omitempty"`
}

// KongUpstreamHash defines how to calculate hash for consistent-hashing load balancing algorithm.
// Only one of the fields must be set.
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

// KongUpstreamHealthcheck represents a health-check config of an Upstream in Kong.
type KongUpstreamHealthcheck struct {
	// Active configures active health check probing.
	Active *KongUpstreamActiveHealthcheck `json:"active,omitempty"`

	// Passive configures passive health check probing.
	Passive *KongUpstreamPassiveHealthcheck `json:"passive,omitempty"`

	// Threshold is the minimum percentage of the upstream’s targets’ weight that must be available for the whole
	// upstream to be considered healthy.
	Threshold *int `json:"threshold,omitempty"`
}

// KongUpstreamActiveHealthcheck configures active health check probing.
type KongUpstreamActiveHealthcheck struct {
	// Type determines how active health checks are collected.
	Type *string `json:"type,omitempty"`

	// Concurrency is the number of targets to check concurrently.
	// +kubebuilder:validation:Minimum=1
	Concurrency *int `json:"concurrency,omitempty"`

	// Healthy configures thresholds and HTTP status codes to mark targets healthy for an upstream.
	Healthy *KongUpstreamHealthcheckHealthy `json:"healthy,omitempty"`

	// Unhealthy configures thresholds and HTTP status codes to mark targets unhealthy for an upstream.
	Unhealthy *KongUpstreamHealthcheckUnhealthy `json:"unhealthy,omitempty"`

	// HTTPPath is the path to use in GET HTTP request to run as a probe.
	// +kubebuilder:validation:Pattern=^/.*$
	HTTPPath *string `json:"http_path,omitempty"`

	// HTTPSSni is the SNI to use in GET HTTPS request to run as a probe.
	HTTPSSni *string `json:"https_sni,omitempty"`

	// HTTPSVerifyCertificate is a boolean value that indicates if the certificate should be verified.
	HTTPSVerifyCertificate *bool `json:"https_verify_certificate,omitempty"`

	// Timeout is the probe timeout in seconds.
	// +kubebuilder:validation:Minimum=0
	Timeout *int `json:"timeout,omitempty"`

	// Headers is a list of HTTP headers to add to the probe request.
	Headers map[string][]string `json:"headers,omitempty"`
}

// KongUpstreamPassiveHealthcheck configures passive checks around
// passive health checks.
type KongUpstreamPassiveHealthcheck struct {
	// Type determines whether to perform passive health checks interpreting HTTP/HTTPS statuses,
	// or just check for TCP connection success.
	// Accepted values are "http", "https", "tcp", "grpc", "grpcs".
	// +kubebuilder:validation:Enum=http;https;tcp;grpc;grpcs
	Type *string `json:"type,omitempty"`

	// Healthy configures thresholds and HTTP status codes to mark targets healthy for an upstream.
	Healthy *KongUpstreamHealthcheckHealthy `json:"healthy,omitempty"`

	// Unhealthy configures thresholds and HTTP status codes to mark targets unhealthy.
	Unhealthy *KongUpstreamHealthcheckUnhealthy `json:"unhealthy,omitempty"`
}

// KongUpstreamHealthcheckHealthy configures thresholds and HTTP status codes to mark targets healthy for an upstream.
type KongUpstreamHealthcheckHealthy struct {
	// HTTPStatuses is a list of HTTP status codes that Kong considers a success.
	HTTPStatuses []int `json:"http_statuses,omitempty"`

	// Interval is the interval between active health checks for an upstream in seconds when in a healthy state.
	// +kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty"`

	// Successes is the number of successes to consider a target healthy.
	// +kubebuilder:validation:Minimum=0
	Successes *int `json:"successes,omitempty"`
}

// KongUpstreamHealthcheckUnhealthy configures thresholds and HTTP status codes to mark targets unhealthy.
type KongUpstreamHealthcheckUnhealthy struct {
	// HTTPFailures is the number of failures to consider a target unhealthy.
	// +kubebuilder:validation:Minimum=0
	HTTPFailures *int `json:"http_failures,omitempty"`

	// HTTPStatuses is a list of HTTP status codes that Kong considers a failure.
	HTTPStatuses []int `json:"http_statuses,omitempty"`

	// TCPFailures is the number of TCP failures in a row to consider a target unhealthy.
	// +kubebuilder:validation:Minimum=0
	TCPFailures *int `json:"tcp_failures,omitempty"`

	// Timeouts is the number of timeouts in a row to consider a target unhealthy.
	// +kubebuilder:validation:Minimum=0
	Timeouts *int `json:"timeouts,omitempty"`

	// Interval is the interval between active health checks for an upstream in seconds when in an unhealthy state.
	// +kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty"`
}
