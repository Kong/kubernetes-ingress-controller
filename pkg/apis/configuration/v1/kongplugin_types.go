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
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:shortName=kp,categories=kong-ingress-controller
// +kubebuilder:validation:Optional
// +kubebuilder:printcolumn:name="Plugin-Type",type=string,JSONPath=`.plugin`,description="Name of the plugin"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"
// +kubebuilder:printcolumn:name="Disabled",type=boolean,JSONPath=`.disabled`,description="Indicates if the plugin is disabled",priority=1
// +kubebuilder:printcolumn:name="Config",type=string,JSONPath=`.config`,description="Configuration of the plugin",priority=1

// KongPlugin is the Schema for the kongplugins API.
type KongPlugin struct {
	metav1.TypeMeta `json:",inline"`
	// Setting a `global` label to `true` will apply the plugin to every request proxied by the Kong.
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
	// Only one of `config` or `configFrom` may be used in a KongPlugin, not both at once.
	// +kubebuilder:validation:Type=object
	Config apiextensionsv1.JSON `json:"config,omitempty"`

	// ConfigFrom references a secret containing the plugin configuration.
	// This should be used when the plugin configuration contains sensitive information,
	// such as AWS credentials in the Lambda plugin or the client secret in the OIDC plugin.
	// Only one of `config` or `configFrom` may be used in a KongPlugin, not both at once.
	ConfigFrom *ConfigSource `json:"configFrom,omitempty"`

	// PluginName is the name of the plugin to which to apply the config.
	// +kubebuilder:validation:Required
	PluginName string `json:"plugin,omitempty"`

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
}

// +kubebuilder:object:root=true

// KongPluginList contains a list of KongPlugin.
type KongPluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongPlugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KongPlugin{}, &KongPluginList{})
}
