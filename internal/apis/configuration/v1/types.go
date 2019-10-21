package v1

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
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
// +genclient:noStatus
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
	Config Configuration `json:"config,omitempty"`

	// PluginName is the name of the plugin to which to apply the config
	PluginName string `json:"plugin,omitempty"`

	// RunOn configures the plugin to run on the first or the second or both
	// nodes in case of a service mesh deployment.
	RunOn string `json:"run_on,omitempty"`

	// Protocols configures plugin to run on requests received on specific
	// protocols.
	Protocols []string `json:"protocols,omitempty"`
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

// Configuration contains a plugin configuration
// +k8s:deepcopy-gen=false
type Configuration map[string]interface{}

func init() {
	gob.Register(map[string]interface{}{})
}

// DeepCopyInto deepcopy function, copying the receiver, writing into out. in must be non-nil.
// TODO: change this to be able to use the k8s code generator
func (in *KongPlugin) DeepCopyInto(out *KongPlugin) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Config != nil {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		dec := gob.NewDecoder(&buf)
		err := enc.Encode(in.Config)
		if err != nil {
			glog.Errorf("unexpected error copying configuration: %v", err)
		}
		err = dec.Decode(&out.Config)
		if err != nil {
			glog.Errorf("unexpected error copying configuration: %v", err)
		}
	}
	return
}

// +genclient
// +genclient:noStatus
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

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KongCredential is a top-level type. A client is created for it.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongCredential struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type string `json:"type,omitempty"`

	ConsumerRef string `json:"consumerRef,omitempty"`

	Config Configuration `json:"config,omitempty"`
}

// KongCredentialList is a top-level list type. The client methods for
// lists are automatically created.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KongCredentialList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []KongCredential `json:"items"`
}

func init() {
	gob.Register(map[string]interface{}{})
}

// DeepCopyInto deepcopy function, copying the receiver, writing into out. in must be non-nil.
// TODO: change this to be able to use the k8s code generator
func (in *KongCredential) DeepCopyInto(out *KongCredential) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Config != nil {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		dec := gob.NewDecoder(&buf)
		err := enc.Encode(in.Config)
		if err != nil {
			glog.Errorf("unexpected error copying configuration: %v", err)
		}
		err = dec.Decode(&out.Config)
		if err != nil {
			glog.Errorf("unexpected error copying configuration: %v", err)
		}
	}
	return
}
