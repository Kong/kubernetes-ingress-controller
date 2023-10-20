package v1beta1

import (
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	v1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// KongUpstreamPolicyCondition is the condition type for KongUpstreamPolicy.
type KongUpstreamPolicyCondition string

// KongUpstreamPolicyStatus is the status type for KongUpstreamPolicy conditions.
type KongUpstreamPolicyStatus string

const (
	// KongUpstreamPolicyConditionAccepted describes the status of a KongUpstreamPolicy with respect to its ancestor.
	KongUpstreamPolicyConditionAccepted KongUpstreamPolicyCondition = KongUpstreamPolicyCondition(v1alpha2.PolicyConditionAccepted)

	// KongUpstreamPolicyStatusAccepted means that the policy is successfully attached to the ancestor.
	KongUpstreamPolicyStatusAccepted KongUpstreamPolicyStatus = KongUpstreamPolicyStatus(v1alpha2.PolicyReasonAccepted)

	// KongUpstreamPolicyStatusConflicted means that the policy couldn't be attached to the ancestor because of a conflict.
	// The conflict might be e.g. attaching KongUpstreamPolicy to a Service and a Gateway API *Route that uses the Service.
	KongUpstreamPolicyStatusConflicted KongUpstreamPolicyStatus = KongUpstreamPolicyStatus(v1alpha2.PolicyReasonConflicted)
)

const (
	// KongUpstreamPolicyAnnotationKey is the key used to attach KongUpstreamPolicy to Services and Gateway API *Routes.
	KongUpstreamPolicyAnnotationKey = annotations.AnnotationPrefix + "/upstream-policy"
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
// It can be attached to Services and Gateway API *Routes. To attach it to an object, the object must be annotated with
// `konghq.com/upstream-policy: <name>`, where `<name>` is the name of the KongUpstreamPolicy object in the same namespace
// as the object.
//
// If attached to multiple objects (ancestors), the controller will populate the status of the KongUpstreamPolicy with
// a separate status entry for each of the ancestors.
//
// When attached to a Gateway API *Route, it will affect all of Kong Upstreams created for
// the Gateway API *Route. If you want to use different Upstream settings for each of the *Route's individual rules,
// you should create a separate *Route for each of the rules.
//
// When attached to a Service, it will affect all Kong Upstreams created for the Service.
//
// When attached to a Service used in a Gateway API *Route rule with multiple BackendRefs, all of its Services MUST
// be configured with the same KongUpstreamPolicy. Otherwise, the controller will *ignore* the KongUpstreamPolicy.
//
// When attached to a Service used in a Gateway API *Route that has another KongUpstreamPolicy attached to it,
// the controller will *ignore* the KongUpstreamPolicy attached to the Service.
//
// Note: KongUpstreamPolicy doesn't implement Gateway API's GEP-713 strictly.
// In particular, it doesn't use the TargetRef for attaching to Services and Gateway API *Routes - annotations are
// used instead. This is to allow reusing the same KongUpstreamPolicy for multiple Services and Gateway API *Routes.
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
	HashOn *KongUpstreamHash `json:"hashOn,omitempty"`

	// HasOnFallback defines how to calculate hash for consistent-hashing load balancing algorithm if the primary hash
	// function fails.
	// Algorithm must be set to "consistent-hashing" for this field to have effect.
	HashOnFallback *KongUpstreamHash `json:"hashOnFallback,omitempty"`

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
	QueryArg *string `json:"queryArg,omitempty"`

	// URICapture is the name of the URI capture group to use as hash input.
	URICapture *string `json:"uriCapture,omitempty"`
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
	// Type determines whether to perform active health checks using HTTP or HTTPS, or just attempt a TCP connection.
	// Accepted values are "http", "https", "tcp", "grpc", "grpcs".
	// +kubebuilder:validation:Enum=http;https;tcp;grpc;grpcs
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
	HTTPPath *string `json:"httpPath,omitempty"`

	// HTTPSSNI is the SNI to use in GET HTTPS request to run as a probe.
	HTTPSSNI *string `json:"httpsSni,omitempty"`

	// HTTPSVerifyCertificate is a boolean value that indicates if the certificate should be verified.
	HTTPSVerifyCertificate *bool `json:"httpsVerifyCertificate,omitempty"`

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
	HTTPStatuses []int `json:"httpStatuses,omitempty"`

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
	HTTPFailures *int `json:"httpFailures,omitempty"`

	// HTTPStatuses is a list of HTTP status codes that Kong considers a failure.
	HTTPStatuses []int `json:"httpStatuses,omitempty"`

	// TCPFailures is the number of TCP failures in a row to consider a target unhealthy.
	// +kubebuilder:validation:Minimum=0
	TCPFailures *int `json:"tcpFailures,omitempty"`

	// Timeouts is the number of timeouts in a row to consider a target unhealthy.
	// +kubebuilder:validation:Minimum=0
	Timeouts *int `json:"timeouts,omitempty"`

	// Interval is the interval between active health checks for an upstream in seconds when in an unhealthy state.
	// +kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty"`
}
