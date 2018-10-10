/*
Copyright 2016 The Kubernetes Authors.

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

package ingress

import (
	"crypto/x509"
	"time"

	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Configuration holds the definition of all the parts required to describe all
// ingresses reachable by the ingress controller (using a filter by namespace)
type Configuration struct {
	// Backends are a list of backends used by all the Ingress rules in the
	// ingress controller. This list includes the default backend
	Backends []*Backend `json:"backends,omitempty"`
	// Servers
	Servers []*Server `json:"servers,omitempty"`
}

// Backend describes one or more remote server/s (endpoints) associated with a service
// +k8s:deepcopy-gen=true
type Backend struct {
	// Name represents an unique apiv1.Service name formatted as <namespace>-<name>-<port>
	Name    string             `json:"name"`
	Service *apiv1.Service     `json:"service,omitempty"`
	Port    intstr.IntOrString `json:"port"`
	// This indicates if the communication protocol between the backend and the endpoint is HTTP or HTTPS
	// Allowing the use of HTTPS
	// The endpoint/s must provide a TLS connection.
	// The certificate used in the endpoint cannot be a self signed certificate
	Secure bool `json:"secure"`
	// Endpoints contains the list of endpoints currently running
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// Endpoint describes a kubernetes endpoint in a backend
// +k8s:deepcopy-gen=true
type Endpoint struct {
	// Address IP address of the endpoint
	Address string `json:"address"`
	// Port number of the TCP port
	Port string `json:"port"`
	// MaxFails returns the number of unsuccessful attempts to communicate
	// allowed before this should be considered dow.
	// Setting 0 indicates that the check is performed by a Kubernetes probe
	MaxFails int `json:"maxFails"`
	// FailTimeout returns the time in seconds during which the specified number
	// of unsuccessful attempts to communicate with the server should happen
	// to consider the endpoint unavailable
	FailTimeout int `json:"failTimeout"`
	// Target returns a reference to the object providing the endpoint
	Target *apiv1.ObjectReference `json:"target,omipempty"`
}

// Server describes a website
type Server struct {
	// Hostname returns the FQDN of the server
	Hostname string `json:"hostname"`

	SSLCert *SSLCert `json:"-"`

	// Locations list of URIs configured in the server.
	Locations []*Location `json:"locations,omitempty"`
	// Alias return the alias of the server name
	Alias string `json:"alias,omitempty"`
}

// Location describes an URI inside a server.
// Also contains additional information about annotations in the Ingress.
//
// In some cases when more than one annotations is defined a particular order in the execution
// is required.
// The chain in the execution order of annotations should be:
// - Whitelist
// - RateLimit
// - BasicDigestAuth
// - ExternalAuth
// - Redirect
type Location struct {
	// Path is an extended POSIX regex as defined by IEEE Std 1003.1,
	// (i.e this follows the egrep/unix syntax, not the perl syntax)
	// matched against the path of an incoming request. Currently it can
	// contain characters disallowed from the conventional "path"
	// part of a URL as defined by RFC 3986. Paths must begin with
	// a '/'. If unspecified, the path defaults to a catch all sending
	// traffic to the backend.
	Path string `json:"path"`
	// IsDefBackend indicates if service specified in the Ingress
	// contains active endpoints or not. Returning true means the location
	// uses the default backend.
	IsDefBackend bool `json:"isDefBackend"`
	// Ingress returns the ingress from which this location was generated
	Ingress *extensions.Ingress `json:"ingress"`
	// Backend describes the name of the backend to use.
	Backend string `json:"backend"`
	// Service describes the referenced services from the ingress
	Service *apiv1.Service `json:"service,omitempty"`
	// Port describes to which port from the service
	Port intstr.IntOrString `json:"port"`
	// Denied returns an error when this location cannot not be allowed
	// Requesting a denied location should return HTTP code 403.
	Denied         error          `json:"denied,omitempty"`
	DefaultBackend *apiv1.Service `json:"defaultBackend,omitempty"`

	KongPlugins map[string]string
}

// SSLCert describes a SSL certificate to be used in a server
type SSLCert struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ID contains the object metadata UID from the Kubenretes Secret
	ID string `json:"id,omitempty"`

	Certificate *x509.Certificate `json:"certificate,omitempty"`

	// CN contains all the common names defined in the SSL certificate
	CN []string `json:"cn"`
	// ExpiresTime contains the expiration of this SSL certificate in timestamp format
	ExpireTime time.Time `json:"expires"`

	Raw RawSSLCert
}

// RawSSLCert represnts cert and key in bytes
type RawSSLCert struct {
	Cert []byte
	Key  []byte
}
