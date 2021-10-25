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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+genclient
//+genclient:nonNamespaced
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster,shortName=kcp
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:validation:Optional
//+kubebuilder:printcolumn:name="Plugin-Type",type=string,JSONPath=`.plugin`,description="Name of the plugin"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"
//+kubebuilder:printcolumn:name="Disabled",type=boolean,JSONPath=`.disabled`,description="Indicates if the plugin is disabled",priority=1
//+kubebuilder:printcolumn:name="Config",type=string,JSONPath=`.config`,description="Configuration of the plugin",priority=1

// KongClusterPlugin is the Schema for the kongclusterplugins API
type KongClusterPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ConsumerRef is a reference to a particular consumer
	ConsumerRef string `json:"consumerRef,omitempty"`

	// Disabled set if the plugin is disabled or not
	Disabled bool `json:"disabled,omitempty"`

	// Config contains the plugin configuration.
	//+kubebuilder:validation:Type=object
	Config apiextensionsv1.JSON `json:"config,omitempty"`

	// ConfigFrom references a secret containing the plugin configuration.
	ConfigFrom *NamespacedConfigSource `json:"configFrom,omitempty"`

	// PluginName is the name of the plugin to which to apply the config
	//+kubebuilder:validation:Required
	PluginName string `json:"plugin,omitempty"`

	// RunOn configures the plugin to run on the first or the second or both
	// nodes in case of a service mesh deployment.
	//+kubebuilder:validation:Enum:=first;second;all
	RunOn string `json:"run_on,omitempty"`

	// Protocols configures plugin to run on requests received on specific
	// protocols.
	Protocols []KongProtocol `json:"protocols,omitempty"`
}

//+kubebuilder:object:root=true

// KongClusterPluginList contains a list of KongClusterPlugin
type KongClusterPluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongClusterPlugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KongClusterPlugin{}, &KongClusterPluginList{})
}
