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

	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
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

	Name        string            `json:"name"`
	Certificate InlineCertificate `json:"certificate"`
}

type InlineCertificate struct {
	ID string `json:"id"`
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

func NewUpstream(name string) *Upstream {
	return &Upstream{
		Name:         name,
		HashOn:       "none",
		HashFallback: "none",
		Slots:        1000,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Upstream struct {
	Required `json:",inline"`

	Name string `json:"name"`

	HashOn             string        `json:"hash_on,omitempty"`
	HashFallback       string        `json:"hash_fallback,omitempty"`
	HashFallbackHeader string        `json:"hash_fallback_header,omitempty"`
	HashOnHeader       string        `json:"hash_on_header,omitempty"`
	HashOnCookie       string        `json:"hash_on_cookie,omitempty"`
	HashOnCookiePath   string        `json:"hash_on_cookie_path,omitempty"`
	Healthchecks       *Healthchecks `json:"healthchecks,omitempty"`
	Slots              int           `json:"slots,omitempty"`
}

type Healthchecks struct {
	Active  *ActiveHealthCheck `json:"active,omitempty"`
	Passive *Passive           `json:"passive,omitempty"`
}

type ActiveHealthCheck struct {
	Concurrency int        `json:"concurrency,omitempty"`
	Healthy     *Healthy   `json:"healthy,omitempty"`
	HTTPPath    string     `json:"http_path,omitempty"`
	Timeout     int        `json:"timeout,omitempty"`
	Unhealthy   *Unhealthy `json:"unhealthy,omitempty"`
}

type Passive struct {
	Healthy   Healthy    `json:"healthy,omitempty"`
	Unhealthy *Unhealthy `json:"unhealthy,omitempty"`
}

type Healthy struct {
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	Interval     int   `json:"interval,omitempty"`
	Successes    int   `json:"successes,omitempty"`
}

type Unhealthy struct {
	HTTPFailures int   `json:"http_failures,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	Interval     int   `json:"interval,omitempty"`
	TCPFailures  int   `json:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty"`
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

	PreserveHost  bool `json:"preserve_host"`
	StripPath     bool `json:"strip_path"`
	RegexPriority int  `json:"regex_priority"`

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
