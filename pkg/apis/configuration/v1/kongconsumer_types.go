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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:shortName=kc,categories=kong-ingress-controller
// +kubebuilder:validation:Optional
// +kubebuilder:printcolumn:name="Username",type=string,JSONPath=`.username`,description="Username of a Kong Consumer"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"

// KongConsumer is the Schema for the kongconsumers API.
type KongConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Username is a Kong cluster-unique username of the consumer.
	Username string `json:"username,omitempty"`

	// CustomID is a Kong cluster-unique existing ID for the consumer - useful for mapping
	// Kong with users in your existing database.
	CustomID string `json:"custom_id,omitempty"`

	// Credentials are references to secrets containing a credential to be
	// provisioned in Kong.
	Credentials []string `json:"credentials,omitempty"`

	// Status represents the current status of the KongConsumer resource.
	Status KongConsumerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KongConsumerList contains a list of KongConsumer.
type KongConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongConsumer `json:"items"`
}

// KongConsumerStatus represents the current status of the KongConsumer resource.
type KongConsumerStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KongConsumer{}, &KongConsumerList{})
}
