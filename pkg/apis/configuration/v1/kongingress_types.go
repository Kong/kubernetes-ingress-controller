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

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:resource:shortName=ki
//+kubebuilder:validation:Optional

// KongIngress is the Schema for the kongingresses API
type KongIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Upstream *KongIngressUpstream `json:"upstream,omitempty"`
	Proxy    *KongIngressService  `json:"proxy,omitempty"`
	Route    *KongIngressRoute    `json:"route,omitempty"`
}

//+kubebuilder:object:root=true

// KongIngressList contains a list of KongIngress
type KongIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongIngress `json:"items"`
}

// KongIngressService contains KongIngress service configuration
//+ It contains the subset of go-kong.kong.Service fields supported by kongstate.Service.overrideByKongIngress
type KongIngressService struct {
	//+kubebuilder:validation:Enum=http;https;grpc;grpcs;tcp;tls;udp
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	//+kubebuilder:validation:Pattern=^/.*$
	Path *string `json:"path,omitempty" yaml:"path,omitempty"`
	//+kubebuilder:validation:Minimum=0
	Retries *int `json:"retries,omitempty" yaml:"retries,omitempty"`
	//+kubebuilder:validation:Minimum=0
	ConnectTimeout *int `json:"connect_timeout,omitempty" yaml:"connect_timeout,omitempty"`
	//+kubebuilder:validation:Minimum=0
	ReadTimeout *int `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty"`
	//+kubebuilder:validation:Minimum=0
	WriteTimeout *int `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty"`
}

// KongIngressRoute contains KongIngress route configuration
//+ It contains the subset of go-kong.kong.Route fields supported by kongstate.Route.overrideByKongIngress
type KongIngressRoute struct {
	Methods                 []*string           `json:"methods,omitempty" yaml:"methods,omitempty"`
	Headers                 map[string][]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Protocols               []*KongProtocol     `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	RegexPriority           *int                `json:"regex_priority,omitempty" yaml:"regex_priority,omitempty"`
	StripPath               *bool               `json:"strip_path,omitempty" yaml:"strip_path,omitempty"`
	PreserveHost            *bool               `json:"preserve_host,omitempty" yaml:"preserve_host,omitempty"`
	HTTPSRedirectStatusCode *int                `json:"https_redirect_status_code,omitempty" yaml:"https_redirect_status_code,omitempty"`
	//+kubebuilder:validation:Enum=v0;v1
	PathHandling      *string   `json:"path_handling,omitempty" yaml:"path_handling,omitempty"`
	SNIs              []*string `json:"snis,omitempty" yaml:"snis,omitempty"`
	RequestBuffering  *bool     `json:"request_buffering,omitempty" yaml:"request_buffering,omitempty"`
	ResponseBuffering *bool     `json:"response_buffering,omitempty" yaml:"response_buffering,omitempty"`
}

// KongIngressUpstream contains KongIngress upstream configuration
//+ It contains the subset of go-kong.kong.Upstream fields supported by kongstate.Upstream.overrideByKongIngress
type KongIngressUpstream struct {
	HostHeader *string `json:"host_header,omitempty" yaml:"host_header,omitempty"`
	//+kubebuilder:validation:Enum=round-robin;consistent-hashing;least-connections
	Algorithm *string `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`
	//+kubebuilder:validation:Minimum=10
	Slots              *int              `json:"slots,omitempty" yaml:"slots,omitempty"`
	Healthchecks       *kong.Healthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
	HashOn             *string           `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`
	HashFallback       *string           `json:"hash_fallback,omitempty" yaml:"hash_fallback,omitempty"`
	HashOnHeader       *string           `json:"hash_on_header,omitempty" yaml:"hash_on_header,omitempty"`
	HashFallbackHeader *string           `json:"hash_fallback_header,omitempty" yaml:"hash_fallback_header,omitempty"`
	HashOnCookie       *string           `json:"hash_on_cookie,omitempty" yaml:"hash_on_cookie,omitempty"`
	HashOnCookiePath   *string           `json:"hash_on_cookie_path,omitempty" yaml:"hash_on_cookie_path,omitempty"`
	// TODO status of this in existing code is unclear. While we supported a raw dump from KongIngress.upstream
	// into the generated upstream, the Kong upstream type only has a certificate ID field here, not a complete
	// certificate/key pair (aka go-kong Certificate). Unclear if db-less or deck have some logic to automagically
	// create the certificate and insert the correct ID into the upstream.
	// We _DID NOT_ show any client cert example at https://docs.konghq.com/kubernetes-ingress-controller/1.3.x/references/custom-resources/#kongingress
	// nor did we have tests for it
	//ClientCertificate  *Certificate `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KongIngress{}, &KongIngressList{})
}
