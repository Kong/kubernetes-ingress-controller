/*
Copyright 2022 Kong, Inc.

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

const (
	IngressClassParametersKind = "IngressClassParameters"
)

//+kubebuilder:object:root=true

// IngressClassParametersList contains a list of IngressClassParameters
type IngressClassParametersList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngressClassParameters `json:"items"`
}

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:resource:categories=kong-ingress-controller

// IngressClassParameters is the Schema for the IngressClassParameters API
type IngressClassParameters struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IngressClassParametersSpec `json:"spec,omitempty"`
}

// IngressClassParametersSpec defines the desired state of IngressClassParameters

type IngressClassParametersSpec struct {
	// Offload load-balancing to kube-proxy or sidecar
	//+kubebuilder:default:=false
	ServiceUpstream bool `json:"serviceUpstream,omitempty"`
}

func init() {
	SchemeBuilder.Register(&IngressClassParameters{}, &IngressClassParametersList{})
}
