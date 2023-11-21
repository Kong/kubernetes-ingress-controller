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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=kcp,categories=kong-ingress-controller
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Plugin-Type",type=string,JSONPath=`.plugin`,description="Name of the plugin"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"
// +kubebuilder:printcolumn:name="Disabled",type=boolean,JSONPath=`.disabled`,description="Indicates if the plugin is disabled",priority=1
// +kubebuilder:printcolumn:name="Config",type=string,JSONPath=`.config`,description="Configuration of the plugin",priority=1
// +kubebuilder:printcolumn:name="Programmed",type=string,JSONPath=`.status.conditions[?(@.type=="Programmed")].status`
// +kubebuilder:validation:XValidation:rule="!(has(self.config) && has(self.configFrom))", message="Using both config and configFrom fields is not allowed."
// +kubebuilder:validation:XValidation:rule="!(has(self.configFrom) && has(self.configPatches))", message="Using both configFrom and configPatches fields is not allowed."
// +kubebuilder:validation:XValidation:rule="self.plugin == oldSelf.plugin", message="The plugin field is immutable"

// KongClusterPlugin is the Schema for the kongclusterplugins API.
type KongClusterPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ConsumerRef is a reference to a particular consumer.
	ConsumerRef string `json:"consumerRef,omitempty"`

	// Disabled set if the plugin is disabled or not.
	Disabled bool `json:"disabled,omitempty"`

	// Config contains the plugin configuration. It's a list of keys and values
	// required to configure the plugin.
	// Please read the documentation of the plugin being configured to set values
	// in here. For any plugin in Kong, anything that goes in the `config` JSON
	// key in the Admin API request, goes into this property.
	// Only one of `config` or `configFrom` may be used in a KongClusterPlugin, not both at once.
	// +kubebuilder:validation:Type=object
	Config apiextensionsv1.JSON `json:"config,omitempty"`

	// ConfigFrom references a secret containing the plugin configuration.
	// This should be used when the plugin configuration contains sensitive information,
	// such as AWS credentials in the Lambda plugin or the client secret in the OIDC plugin.
	// Only one of `config` or `configFrom` may be used in a KongClusterPlugin, not both at once.
	ConfigFrom *NamespacedConfigSource `json:"configFrom,omitempty"`

	// ConfigPatches represents JSON patches to the configuration of the plugin.
	// Each item means a JSON patch to add something in the configuration,
	// where path is specified in `path` and value is in `valueFrom` referencing
	// a key in a secret.
	// When Config is specified, patches will be applied to the configuration in Config.
	// Otherwise, patches will be applied to an empty object.
	ConfigPatches []NamespacedConfigPatch `json:"configPatches,omitempty"`

	// PluginName is the name of the plugin to which to apply the config.
	// +kubebuilder:validation:Required
	PluginName string `json:"plugin"`

	// RunOn configures the plugin to run on the first or the second or both
	// nodes in case of a service mesh deployment.
	// +kubebuilder:validation:Enum:=first;second;all
	RunOn string `json:"run_on,omitempty"`

	// Protocols configures plugin to run on requests received on specific
	// protocols.
	Protocols []KongProtocol `json:"protocols,omitempty"`

	// Ordering overrides the normal plugin execution order. It's only available on Kong Enterprise.
	// `<phase>` is a request processing phase (for example, `access` or `body_filter`) and
	// `<plugin>` is the name of the plugin that will run before or after the KongPlugin.
	// For example, a KongPlugin with `plugin: rate-limiting` and `before.access: ["key-auth"]`
	// will create a rate limiting plugin that limits requests _before_ they are authenticated.
	Ordering *kong.PluginOrdering `json:"ordering,omitempty"`

	// InstanceName is an optional custom name to identify an instance of the plugin. This is useful when running the
	// same plugin in multiple contexts, for example, on multiple services.
	InstanceName string `json:"instance_name,omitempty"`

	// Status represents the current status of the KongClusterPlugin resource.
	Status KongClusterPluginStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KongClusterPluginList contains a list of KongClusterPlugin.
type KongClusterPluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongClusterPlugin `json:"items"`
}

// KongClusterPluginStatus represents the current status of the KongClusterPlugin resource.
type KongClusterPluginStatus struct {
	// Conditions describe the current conditions of the KongClusterPluginStatus.
	//
	// Known condition types are:
	//
	// * "Programmed"
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={{type: "Programmed", status: "Unknown", reason:"Pending", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KongClusterPlugin{}, &KongClusterPluginList{})
}
