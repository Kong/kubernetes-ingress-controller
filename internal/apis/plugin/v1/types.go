package v1

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
