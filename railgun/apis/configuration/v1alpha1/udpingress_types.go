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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// UDPIngress is the Schema for the udpingresses API
type UDPIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UDPIngressSpec   `json:"spec,omitempty"`
	Status UDPIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UDPIngressList contains a list of UDPIngress
type UDPIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UDPIngress `json:"items"`
}

// UDPIngressSpec defines the desired state of UDPIngress
type UDPIngressSpec struct {
	// Host indicates where to send the UDP datagrams
	Host string `json:"host,required" yaml:"host,required"`

	// ListenPort indicates the Kong proxy port which will accept the ingress datagrams
	ListenPort int `json:"listenPort,required" yaml:"listenPort,required"`

	// TargetPort indicates the backend Host port which kong will proxy the UDP datagrams to
	TargetPort int `json:"targetPort,required" yaml:"targetPort,required"`
}

// UDPIngressStatus defines the observed state of UDPIngress
type UDPIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

func init() {
	SchemeBuilder.Register(&UDPIngress{}, &UDPIngressList{})
}
