package v1

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Required defines the required fields to work between Kubernetes
// and Kong and also defines common field present in Kong entities
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Required struct {
	metav1.TypeMeta   `json:"-"`
	metav1.ObjectMeta `json:"-"`

	ID string `json:"id,omitempty"`

	Tags []string `json:"tags,omitempty"`

	CreatedAt int `json:"created_at,omitempty"`
	UpdatedAt int `json:"updated_at,omitempty"`
}

// RequiredList ...
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RequiredList struct {
	metav1.TypeMeta `json:"-"`
	metav1.ListMeta `json:"-"`

	NextPage string `json:"next"`
	Offset   string `json:"offset"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SNI struct {
	Required `json:",inline"`

	Name        string `json:"name"`
	Certificate string `json:"ssl_certificate_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SNIList struct {
	RequiredList `json:",inline"`

	Items []SNI `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Certificate struct {
	Required `json:",inline"`

	Cert  string   `json:"cert"`
	Key   string   `json:"key"`
	Hosts []string `json:"snis"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CertificateList struct {
	RequiredList `json:",inline"`

	Items []Certificate `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Credential struct {
	Required `json:",inline"`

	Consumer string `json:"consumer_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CredentialList struct {
	RequiredList `json:",inline"`

	Items []Credential `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Service struct {
	Required `json:",inline"`

	Name string `json:"name"`

	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`

	Retries int `json:"retries,omitempty"`

	ConnectTimeout int `json:"connect_timeout,omitempty"`
	ReadTimeout    int `json:"read_timeout,omitempty"`
	WriteTimeout   int `json:"write_timeout,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServiceList struct {
	metav1.TypeMeta `json:"-"`
	metav1.ListMeta `json:"-"`

	Items    []Service `json:"data"`
	NextPage string    `json:"next"`
	Offset   string    `json:"offset"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Upstream struct {
	Required `json:",inline"`

	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Target struct {
	Required `json:",inline"`

	Target   string `json:"target"`
	Weight   int    `json:"weight,omitempty"`
	Upstream string `json:"upstream_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Route struct {
	Required `json:",inline"`

	Protocols []string `json:"protocols"`
	Hosts     []string `json:"hosts"`
	Paths     []string `json:"paths"`
	Methods   []string `json:"methods"`

	PreserveHost bool `json:"preserve_host"`
	StripPath    bool `json:"strip_path"`

	Service InlineService `json:"service"`
}

type InlineService struct {
	ID string `json:"id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RouteList struct {
	RequiredList `json:",inline"`

	Items []Route `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type UpstreamList struct {
	RequiredList `json:",inline"`

	Items []Upstream `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TargetList struct {
	RequiredList `json:",inline"`

	Items []Target `json:"data"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Plugin struct {
	Required `json:",inline"`

	Name string `json:"name"`

	Config   Configuration `json:"config,omitempty"`
	Enabled  bool          `json:"enabled,omitempty"`
	Service  string        `json:"service_id,omitempty"`
	Route    string        `json:"route_id,omitempty"`
	Consumer string        `json:"consumer_id,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PluginList struct {
	RequiredList `json:",inline"`

	Items []Plugin `json:"data"`
}

// DeepCopyInto deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Plugin) DeepCopyInto(out *Plugin) {
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Consumer struct {
	Required `json:",inline"`

	Username string `json:"username,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConsumerList struct {
	RequiredList `json:",inline"`

	Items []Consumer `json:"data"`
}

// Configuration contains a plugin configuration
// +k8s:deepcopy-gen=false
type Configuration map[string]interface{}

// Equal tests for equality between two Configuration types
func (r1 *Route) Equal(r2 *Route) bool {
	if r1 == r2 {
		return true
	}
	if r1 == nil || r2 == nil {
		return false
	}

	if r1.Service.ID != r2.Service.ID {
		return false
	}

	if len(r1.Hosts) != len(r2.Hosts) {
		return false
	}

	for _, r1b := range r1.Hosts {
		found := false
		for _, r2b := range r2.Hosts {
			if r1b == r2b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(r1.Paths) != len(r2.Paths) {
		return false
	}

	for _, r1b := range r1.Paths {
		found := false
		for _, r2b := range r2.Paths {
			if r1b == r2b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
