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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:validation:Optional
//+kubebuilder:printcolumn:name="Address",type=string,JSONPath=`.status.loadBalancer.ingress[*].ip`,description="Address of the load balancer"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"

// TCPIngress is the Schema for the tcpingresses API
type TCPIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TCPIngressSpec   `json:"spec,omitempty"`
	Status TCPIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TCPIngressList contains a list of TCPIngress
type TCPIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TCPIngress `json:"items"`
}

// TCPIngressSpec defines the desired state of TCPIngress
type TCPIngressSpec struct {
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

// TCPIngressStatus defines the observed state of TCPIngress
type TCPIngressStatus struct {
	// LoadBalancer contains the current status of the load-balancer.
	// +optional
	LoadBalancer corev1.LoadBalancerStatus `json:"loadBalancer,omitempty"`
}

func init() {
	SchemeBuilder.Register(&TCPIngress{}, &TCPIngressList{})
}
