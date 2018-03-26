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

// Configuration contains a plugin configuration
// +k8s:deepcopy-gen=false
type Configuration map[string]interface{}

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
