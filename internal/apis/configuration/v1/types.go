package v1

import (
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

	Upstream *Upstream `json:"upstream,omitempty"`
	Proxy    *Proxy    `json:"proxy,omitempty"`
	Route    *Route    `json:"route,omitempty"`
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

// Route defines optional settings defined in Kong Routes
type Route struct {
	Methods       []string `json:"methods"`
	RegexPriority int      `json:"regex_priority"`
	StripPath     bool     `json:"strip_path"`
	PreserveHost  bool     `json:"preserve_host"`
	Protocols     []string `json:"protocols"`
}

type Upstream struct {
	HashOn             string        `json:"hash_on"`
	HashOnCookie       string        `json:"hash_on_cookie"`
	HashOnCookiePath   string        `json:"hash_on_cookie_path"`
	HashOnHeader       string        `json:"hash_on_header"`
	HashFallback       string        `json:"hash_fallback"`
	HashFallbackHeader string        `json:"hash_fallback_header"`
	Healthchecks       *Healthchecks `json:"healthchecks,omitempty"`
	Slots              int           `json:"slots"`
}

type Proxy struct {
	Protocol       string `json:"protocol"`
	Path           string `json:"path"`
	ConnectTimeout int    `json:"connect_timeout"`
	Retries        int    `json:"retries"`
	ReadTimeout    int    `json:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
}

type Healthchecks struct {
	Active  *ActiveHealthCheck `json:"active,omitempty"`
	Passive *Passive           `json:"passive,omitempty"`
}

type ActiveHealthCheck struct {
	Concurrency int        `json:"concurrency"`
	Healthy     *Healthy   `json:"healthy"`
	HTTPPath    string     `json:"http_path"`
	Timeout     int        `json:"timeout"`
	Unhealthy   *Unhealthy `json:"unhealthy"`
}

type Passive struct {
	Healthy   *Healthy   `json:"healthy,omitempty"`
	Unhealthy *Unhealthy `json:"unhealthy,omitempty"`
}

type Healthy struct {
	HTTPStatuses []int `json:"http_statuses"`
	Interval     int   `json:"interval"`
	Successes    int   `json:"successes"`
}

type Unhealthy struct {
	HTTPFailures int   `json:"http_failures"`
	HTTPStatuses []int `json:"http_statuses"`
	Interval     int   `json:"interval"`
	TCPFailures  int   `json:"tcp_failures"`
	Timeouts     int   `json:"timeouts"`
}
