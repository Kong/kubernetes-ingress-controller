package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

const (
	// KongUpstreamPolicyAnnotationKey is the key used to attach KongUpstreamPolicy to Services.
	// The value of the annotation is the name of the KongUpstreamPolicy object in the same namespace as the Service.
	KongUpstreamPolicyAnnotationKey = "konghq.com/upstream-policy"
)

func init() {
	SchemeBuilder.Register(&KongUpstreamPolicy{}, &KongUpstreamPolicyList{})
}

// KongUpstreamPolicy allows configuring algorithm that should be used for load balancing traffic between Kong
// Upstream's Targets. It also allows configuring health checks for Kong Upstream's Targets.
//
// Its configuration is similar to Kong Upstream object (https://docs.konghq.com/gateway/latest/admin-api/#upstream-object),
// and it is applied to Kong Upstream objects created by the controller.
//
// It can be attached to Services. To attach it to a Service, it has to be annotated with
// `konghq.com/upstream-policy: <name>`, where `<name>` is the name of the KongUpstreamPolicy
// object in the same namespace as the Service.
//
// When attached to a Service, it will affect all Kong Upstreams created for the Service.
//
// When attached to a Service used in a Gateway API *Route rule with multiple BackendRefs, all of its Services MUST
// be configured with the same KongUpstreamPolicy. Otherwise, the controller will *ignore* the KongUpstreamPolicy.
//
// Note: KongUpstreamPolicy doesn't implement Gateway API's GEP-713 strictly.
// In particular, it doesn't use the TargetRef for attaching to Services and Gateway API *Routes - annotations are
// used instead. This is to allow reusing the same KongUpstreamPolicy for multiple Services and Gateway API *Routes.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kup,categories=kong-ingress-controller
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=gateway.networking.k8s.io/policy=direct
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOn) ? [has(self.spec.hashOn.input), has(self.spec.hashOn.cookie), has(self.spec.hashOn.header), has(self.spec.hashOn.uriCapture), has(self.spec.hashOn.queryArg)].filter(fieldSet, fieldSet == true).size() <= 1 : true", message="Only one of spec.hashOn.(input|cookie|header|uriCapture|queryArg) can be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOn) && has(self.spec.hashOn.cookie) ? has(self.spec.hashOn.cookiePath) : true", message="When spec.hashOn.cookie is set, spec.hashOn.cookiePath is required."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOn) && has(self.spec.hashOn.cookiePath) ? has(self.spec.hashOn.cookie) : true", message="When spec.hashOn.cookiePath is set, spec.hashOn.cookie is required."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOn) ? has(self.spec.algorithm) && self.spec.algorithm == \"consistent-hashing\" : true", message="spec.algorithm must be set to \"consistent-hashing\" when spec.hashOn is set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOnFallback) ? [has(self.spec.hashOnFallback.input), has(self.spec.hashOnFallback.header), has(self.spec.hashOnFallback.uriCapture), has(self.spec.hashOnFallback.queryArg)].filter(fieldSet, fieldSet == true).size() <= 1 : true", message="Only one of spec.hashOnFallback.(input|header|uriCapture|queryArg) can be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOnFallback) ? has(self.spec.algorithm) && self.spec.algorithm == \"consistent-hashing\" : true", message="spec.algorithm must be set to \"consistent-hashing\" when spec.hashOnFallback is set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOnFallback) ? !has(self.spec.hashOnFallback.cookie) : true", message="spec.hashOnFallback.cookie must not be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOnFallback) ? !has(self.spec.hashOnFallback.cookiePath) : true", message="spec.hashOnFallback.cookiePath must not be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.healthchecks) && has(self.spec.healthchecks.passive) && has(self.spec.healthchecks.passive.healthy) ? !has(self.spec.healthchecks.passive.healthy.interval) : true", message="spec.healthchecks.passive.healthy.interval must not be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.healthchecks) && has(self.spec.healthchecks.passive) && has(self.spec.healthchecks.passive.unhealthy) ? !has(self.spec.healthchecks.passive.unhealthy.interval) : true", message="spec.healthchecks.passive.unhealthy.interval must not be set."
// +kubebuilder:validation:XValidation:rule="has(self.spec.hashOn) && has(self.spec.hashOn.cookie) ? !has(self.spec.hashOnFallback) : true", message="spec.hashOnFallback must not be set when spec.hashOn.cookie is set."
type KongUpstreamPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the configuration of the Kong upstream.
	Spec KongUpstreamPolicySpec `json:"spec,omitempty"`

	// Status defines the current state of KongUpstreamPolicy
	Status gatewayv1alpha2.PolicyStatus `json:"status,omitempty"`
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

	// HashOnFallback defines how to calculate hash for consistent-hashing load balancing algorithm if the primary hash
	// function fails.
	// Algorithm must be set to "consistent-hashing" for this field to have effect.
	HashOnFallback *KongUpstreamHash `json:"hashOnFallback,omitempty"`

	// Healthchecks defines the health check configurations in Kong.
	Healthchecks *KongUpstreamHealthcheck `json:"healthchecks,omitempty"`
}

// HashInput is the input for consistent-hashing load balancing algorithm.
// Can be one of: "ip", "consumer", "path".
// +kubebuilder:validation:Enum=ip;consumer;path
type HashInput string

// KongUpstreamHash defines how to calculate hash for consistent-hashing load balancing algorithm.
// Only one of the fields must be set.
type KongUpstreamHash struct {
	// Input allows using one of the predefined inputs (ip, consumer, path).
	// For other parametrized inputs, use one of the fields below.
	Input *HashInput `json:"input,omitempty"`

	// Header is the name of the header to use as hash input.
	Header *string `json:"header,omitempty"`

	// Cookie is the name of the cookie to use as hash input.
	Cookie *string `json:"cookie,omitempty"`

	// CookiePath is cookie path to set in the response headers.
	CookiePath *string `json:"cookiePath,omitempty"`

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

// HTTPStatus is an HTTP status code.
// +kubebuilder:validation:Minimum=100
// +kubebuilder:validation:Maximum=599
type HTTPStatus int

// KongUpstreamHealthcheckHealthy configures thresholds and HTTP status codes to mark targets healthy for an upstream.
type KongUpstreamHealthcheckHealthy struct {
	// HTTPStatuses is a list of HTTP status codes that Kong considers a success.
	HTTPStatuses []HTTPStatus `json:"httpStatuses,omitempty"`

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
	HTTPStatuses []HTTPStatus `json:"httpStatuses,omitempty"`

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
