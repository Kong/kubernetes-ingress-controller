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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
