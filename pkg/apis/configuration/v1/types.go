package v1

import (
	"encoding/gob"

	"github.com/kong/go-kong/kong"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KongIngress is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongIngress struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Upstream *kong.Upstream `json:"upstream,omitempty"`
	Proxy    *kong.Service  `json:"proxy,omitempty"`
	Route    *kong.Route    `json:"route,omitempty"`
}

// KongIngressList is a top-level list type. The client methods for
// lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongIngressList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []KongIngress `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KongClusterPlugin is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongClusterPlugin struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ConsumerRef is a reference to a particular consumer
	ConsumerRef string `json:"consumerRef,omitempty"`

	// Disabled set if the plugin is disabled or not
	Disabled bool `json:"disabled,omitempty"`

	// Config contains the plugin configuration.
	Config apiextensionsv1.JSON `json:"config,omitempty"`

	// ConfigFrom references a secret containing the plugin configuration.
	ConfigFrom NamespacedConfigSource `json:"configFrom,omitempty"`

	// PluginName is the name of the plugin to which to apply the config
	PluginName string `json:"plugin,omitempty"`

	// RunOn configures the plugin to run on the first or the second or both
	// nodes in case of a service mesh deployment.
	RunOn string `json:"run_on,omitempty"`

	// Protocols configures plugin to run on requests received on specific
	// protocols.
	Protocols []string `json:"protocols,omitempty"`
}

// KongClusterPluginList is a top-level list type. The client methods for lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongClusterPluginList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []KongClusterPlugin `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KongPlugin is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongPlugin struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ConsumerRef is a reference to a particular consumer
	ConsumerRef string `json:"consumerRef,omitempty"`

	// Disabled set if the plugin is disabled or not
	Disabled bool `json:"disabled,omitempty"`

	// Config contains the plugin configuration.
	Config apiextensionsv1.JSON `json:"config,omitempty"`

	// ConfigFrom references a secret containing the plugin configuration.
	ConfigFrom ConfigSource `json:"configFrom,omitempty"`

	// PluginName is the name of the plugin to which to apply the config
	PluginName string `json:"plugin,omitempty"`

	// RunOn configures the plugin to run on the first or the second or both
	// nodes in case of a service mesh deployment.
	RunOn string `json:"run_on,omitempty"`

	// Protocols configures plugin to run on requests received on specific
	// protocols.
	Protocols []string `json:"protocols,omitempty"`
}

// ConfigSource is a wrapper around SecretValueFromSource
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConfigSource struct {
	metav1.TypeMeta `json:",inline"`
	SecretValue     SecretValueFromSource `json:"secretKeyRef,omitempty"`
}

// NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NamespacedConfigSource struct {
	metav1.TypeMeta `json:",inline"`
	SecretValue     NamespacedSecretValueFromSource `json:"secretKeyRef,omitempty"`
}

// SecretValueFromSource represents the source of a secret value
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SecretValueFromSource struct {
	metav1.TypeMeta `json:",inline"`
	// the secret containing the key
	Secret string `json:"name,omitempty"`
	// the key containing the value
	Key string `json:"key,omitempty"`
}

// NamespacedSecretValueFromSource represents the source of a secret value,
// specifying the secret namespace
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NamespacedSecretValueFromSource struct {
	metav1.TypeMeta `json:",inline"`
	// The namespace containing the secret
	Namespace string `json:"namespace,omitempty"`
	// the secret containing the key
	Secret string `json:"name,omitempty"`
	// the key containing the value
	Key string `json:"key,omitempty"`
}

// KongPluginList is a top-level list type. The client methods for lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongPluginList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []KongPlugin `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KongConsumer is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongConsumer struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Username unique username of the consumer.
	Username string `json:"username,omitempty"`

	// CustomID existing unique ID for the consumer - useful for mapping
	// Kong with users in your existing database
	CustomID string `json:"custom_id,omitempty"`

	// Credentials are references to secrets containing a credential to be
	// provisioned in Kong.
	Credentials []string `json:"credentials,omitempty"`
}

// KongConsumerList is a top-level list type. The client methods for
// lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []KongConsumer `json:"items"`
}

func init() {
	gob.Register(map[string]interface{}{})
}
