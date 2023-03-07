/*
Copyright 2021 Kong, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/kong/go-kong/kong"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:shortName=ki,categories=kong-ingress-controller
// +kubebuilder:validation:Optional

// KongIngress is the Schema for the kongingresses API.
type KongIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Upstream represents a virtual hostname and can be used to loadbalance
	// incoming requests over multiple targets (e.g. Kubernetes `Services` can
	// be a target, OR `Endpoints` can be targets).
	Upstream *KongIngressUpstream `json:"upstream,omitempty"`

	// Proxy defines additional connection options for the routes to be configured in the
	// Kong Gateway, e.g. `connection_timeout`, `retries`, etc.
	Proxy *KongIngressService `json:"proxy,omitempty"`

	// Route define rules to match client requests.
	// Each Route is associated with a Service,
	// and a Service may have multiple Routes associated to it.
	Route *KongIngressRoute `json:"route,omitempty"`
}

// +kubebuilder:object:root=true

// KongIngressList contains a list of KongIngress.
type KongIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongIngress `json:"items"`
}

// KongIngressService contains KongIngress service configuration.
// It contains the subset of go-kong.kong.Service fields supported by kongstate.Service.overrideByKongIngress.
// Deprecated: use Service's annotations instead.
type KongIngressService struct {
	// The protocol used to communicate with the upstream.
	// Deprecated: use Service's "konghq.com/protocol" annotation instead.
	// +kubebuilder:validation:Enum=http;https;grpc;grpcs;tcp;tls;udp
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	// (optional) The path to be used in requests to the upstream server.
	// Deprecated: use Service's "konghq.com/path" annotation instead.
	// +kubebuilder:validation:Pattern=^/.*$
	Path *string `json:"path,omitempty" yaml:"path,omitempty"`

	// The number of retries to execute upon failure to proxy.
	// Deprecated: use Service's "konghq.com/retries" annotation instead.
	// +kubebuilder:validation:Minimum=0
	Retries *int `json:"retries,omitempty" yaml:"retries,omitempty"`

	// The timeout in milliseconds for	establishing a connection to the upstream server.
	// Deprecated: use Service's "konghq.com/connect-timeout" annotation instead.
	// +kubebuilder:validation:Minimum=0
	ConnectTimeout *int `json:"connect_timeout,omitempty" yaml:"connect_timeout,omitempty"`

	// The timeout in milliseconds between two successive read operations
	// for transmitting a request to the upstream server.
	// Deprecated: use Service's "konghq.com/read-timeout" annotation instead.
	// +kubebuilder:validation:Minimum=0
	ReadTimeout *int `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty"`

	// The timeout in milliseconds between two successive write operations
	// for transmitting a request to the upstream server.
	// Deprecated: use Service's "konghq.com/write-timeout" annotation instead.
	// +kubebuilder:validation:Minimum=0
	WriteTimeout *int `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty"`
}

// KongIngressRoute contains KongIngress route configuration.
// It contains the subset of `go-kong.kong.Route` fields supported by `kongstate.Route.overrideByKongIngress`.
// Deprecated: use Ingress' annotations instead.
type KongIngressRoute struct {
	// Methods is a list of HTTP methods that match this Route.
	// Deprecated: use Ingress' "konghq.com/methods" annotation instead.
	Methods []*string `json:"methods,omitempty" yaml:"methods,omitempty"`

	// Headers contains one or more lists of values indexed by header name
	// that will cause this Route to match if present in the request.
	// The Host header cannot be used with this attribute.
	// Deprecated: use Ingress' "konghq.com/headers" annotation instead.
	Headers map[string][]string `json:"headers,omitempty" yaml:"headers,omitempty"`

	// Protocols is an array of the protocols this Route should allow.
	// Deprecated: use Ingress' "konghq.com/protocols" annotation instead.
	Protocols []*KongProtocol `json:"protocols,omitempty" yaml:"protocols,omitempty"`

	// RegexPriority is a number used to choose which route resolves a given request
	// when several routes match it using regexes simultaneously.
	// Deprecated: use Ingress' "konghq.com/regex-priority" annotation instead.
	RegexPriority *int `json:"regex_priority,omitempty" yaml:"regex_priority,omitempty"`

	// StripPath sets When matching a Route via one of the paths
	// strip the matching prefix from the upstream request URL.
	// Deprecated: use Ingress' "konghq.com/strip-path" annotation instead.
	StripPath *bool `json:"strip_path,omitempty" yaml:"strip_path,omitempty"`

	// PreserveHost sets When matching a Route via one of the hosts domain names,
	// use the request Host header in the upstream request headers.
	// If set to false, the upstream Host header will be that of the Service’s host.
	// Deprecated: use Ingress' "konghq.com/preserve-host" annotation instead.
	PreserveHost *bool `json:"preserve_host,omitempty" yaml:"preserve_host,omitempty"`

	// HTTPSRedirectStatusCode is the status code Kong responds with
	// when all properties of a Route match except the protocol.
	// Deprecated: use Ingress' "ingress.kubernetes.io/force-ssl-redirect" or
	// "konghq.com/https-redirect-status-code" annotations instead.
	HTTPSRedirectStatusCode *int `json:"https_redirect_status_code,omitempty" yaml:"https_redirect_status_code,omitempty"`

	// PathHandling controls how the Service path, Route path and requested path
	// are combined when sending a request to the upstream.
	// +kubebuilder:validation:Enum=v0;v1
	// Deprecated: use Ingress' "konghq.com/path-handling" annotation instead.
	PathHandling *string `json:"path_handling,omitempty" yaml:"path_handling,omitempty"`

	// SNIs is a list of SNIs that match this Route when using stream routing.
	// Deprecated: use Ingress' "konghq.com/snis" annotation instead.
	SNIs []*string `json:"snis,omitempty" yaml:"snis,omitempty"`

	// RequestBuffering sets whether to enable request body buffering or not.
	// Deprecated: use Ingress' "konghq.com/request-buffering" annotation instead.
	RequestBuffering *bool `json:"request_buffering,omitempty" yaml:"request_buffering,omitempty"`

	// ResponseBuffering sets whether to enable response body buffering or not.
	// Deprecated: use Ingress' "konghq.com/response-buffering" annotation instead.
	ResponseBuffering *bool `json:"response_buffering,omitempty" yaml:"response_buffering,omitempty"`
}

// KongIngressUpstream contains KongIngress upstream configuration.
// It contains the subset of `go-kong.kong.Upstream` fields supported by `kongstate.Upstream.overrideByKongIngress`.
type KongIngressUpstream struct {
	// HostHeader is The hostname to be used as Host header
	// when proxying requests through Kong.
	HostHeader *string `json:"host_header,omitempty" yaml:"host_header,omitempty"`

	// Algorithm is the load balancing algorithm to use.
	// +kubebuilder:validation:Enum=round-robin;consistent-hashing;least-connections
	Algorithm *string `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`

	// Slots is the number of slots in the load balancer algorithm.
	// +kubebuilder:validation:Minimum=10
	Slots *int `json:"slots,omitempty" yaml:"slots,omitempty"`

	// Healthchecks defines the health check configurations in Kong.
	Healthchecks *kong.Healthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`

	// HashOn defines what to use as hashing input.
	// Accepted values are: "none", "consumer", "ip", "header", "cookie", "path", "query_arg", "uri_capture".
	HashOn *string `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`

	// HashFallback defines What to use as hashing input
	// if the primary hash_on does not return a hash.
	// Accepted values are: "none", "consumer", "ip", "header", "cookie".
	HashFallback *string `json:"hash_fallback,omitempty" yaml:"hash_fallback,omitempty"`

	// HashOnHeader defines the header name to take the value from as hash input.
	// Only required when "hash_on" is set to "header".
	HashOnHeader *string `json:"hash_on_header,omitempty" yaml:"hash_on_header,omitempty"`

	// HashFallbackHeader is the header name to take the value from as hash input.
	// Only required when "hash_fallback" is set to "header".
	HashFallbackHeader *string `json:"hash_fallback_header,omitempty" yaml:"hash_fallback_header,omitempty"`

	// The cookie name to take the value from as hash input.
	// Only required when "hash_on" or "hash_fallback" is set to "cookie".
	HashOnCookie *string `json:"hash_on_cookie,omitempty" yaml:"hash_on_cookie,omitempty"`

	// The cookie path to set in the response headers.
	// Only required when "hash_on" or "hash_fallback" is set to "cookie".
	HashOnCookiePath *string `json:"hash_on_cookie_path,omitempty" yaml:"hash_on_cookie_path,omitempty"`

	// HashOnQueryArg is the query string parameter whose value is the hash input when "hash_on" is set to "query_arg".
	HashOnQueryArg *string `json:"hash_on_query_arg,omitempty" yaml:"hash_on_query_arg,omitempty"`

	// HashFallbackQueryArg is the "hash_fallback" version of HashOnQueryArg.
	HashFallbackQueryArg *string `json:"hash_fallback_query_arg,omitempty" yaml:"hash_fallback_query_arg,omitempty"`

	// HashOnURICapture is the name of the capture group whose value is the hash input when "hash_on" is set to
	// "uri_capture".
	HashOnURICapture *string `json:"hash_on_uri_capture,omitempty" yaml:"hash_on_uri_capture,omitempty"`

	// HashFallbackURICapture is the "hash_fallback" version of HashOnURICapture.
	HashFallbackURICapture *string `json:"hash_fallback_uri_capture,omitempty" yaml:"hash_fallback_uri_capture,omitempty"`

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2075
	// ClientCertificate  *CertificateSecretRef `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KongIngress{}, &KongIngressList{})
}
