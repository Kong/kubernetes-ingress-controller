/*
Copyright 2017 The Kubernetes Authors.

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

// Equal tests for equality between two Configuration types
func (c1 *Configuration) Equal(c2 *Configuration) bool {
	if c1 == c2 {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}

	if len(c1.Backends) != len(c2.Backends) {
		return false
	}

	for _, c1b := range c1.Backends {
		found := false
		for _, c2b := range c2.Backends {
			if c1b.Equal(c2b) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(c1.Servers) != len(c2.Servers) {
		return false
	}

	// Servers are sorted
	for idx, c1s := range c1.Servers {
		if !c1s.Equal(c2.Servers[idx]) {
			return false
		}
	}

	return true
}

// Equal tests for equality between two Backend types
func (b1 *Backend) Equal(b2 *Backend) bool {
	if b1 == b2 {
		return true
	}
	if b1 == nil || b2 == nil {
		return false
	}
	if b1.Name != b2.Name {
		return false
	}

	if b1.Service != b2.Service {
		if b1.Service == nil || b2.Service == nil {
			return false
		}
		if b1.Service.GetNamespace() != b2.Service.GetNamespace() {
			return false
		}
		if b1.Service.GetName() != b2.Service.GetName() {
			return false
		}
		if b1.Service.GetResourceVersion() != b2.Service.GetResourceVersion() {
			return false
		}
	}

	if b1.Port != b2.Port {
		return false
	}
	if b1.Secure != b2.Secure {
		return false
	}
	if len(b1.Endpoints) != len(b2.Endpoints) {
		return false
	}

	for _, udp1 := range b1.Endpoints {
		found := false
		for _, udp2 := range b2.Endpoints {
			if (&udp1).Equal(&udp2) {
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

// Equal checks the equality against an Endpoint
func (e1 *Endpoint) Equal(e2 *Endpoint) bool {
	if e1 == e2 {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Address != e2.Address {
		return false
	}
	if e1.Port != e2.Port {
		return false
	}
	if e1.MaxFails != e2.MaxFails {
		return false
	}
	if e1.FailTimeout != e2.FailTimeout {
		return false
	}

	if e1.Target != e2.Target {
		if e1.Target == nil || e2.Target == nil {
			return false
		}
		if e1.Target.UID != e2.Target.UID {
			return false
		}
		if e1.Target.ResourceVersion != e2.Target.ResourceVersion {
			return false
		}
	}

	return true
}

// Equal tests for equality between two Server types
func (s1 *Server) Equal(s2 *Server) bool {
	if s1 == s2 {
		return true
	}
	if s1 == nil || s2 == nil {
		return false
	}
	if s1.Hostname != s2.Hostname {
		return false
	}
	if s1.Alias != s2.Alias {
		return false
	}

	if len(s1.Locations) != len(s2.Locations) {
		return false
	}

	// Location are sorted
	for idx, s1l := range s1.Locations {
		if !s1l.Equal(s2.Locations[idx]) {
			return false
		}
	}

	return true
}

// Equal tests for equality between two Location types
func (l1 *Location) Equal(l2 *Location) bool {
	if l1 == l2 {
		return true
	}
	if l1 == nil || l2 == nil {
		return false
	}
	if l1.Path != l2.Path {
		return false
	}
	if l1.IsDefBackend != l2.IsDefBackend {
		return false
	}
	if l1.Backend != l2.Backend {
		return false
	}

	if l1.Service != l2.Service {
		if l1.Service == nil || l2.Service == nil {
			return false
		}
		if l1.Service.GetNamespace() != l2.Service.GetNamespace() {
			return false
		}
		if l1.Service.GetName() != l2.Service.GetName() {
			return false
		}
		if l1.Service.GetResourceVersion() != l2.Service.GetResourceVersion() {
			return false
		}
	}

	if l1.Port.StrVal != l2.Port.StrVal {
		return false
	}
	if l1.Denied != l2.Denied {
		return false
	}

	return true
}

// Equal tests for equality between two L4Backend types
func (s1 *SSLCert) Equal(s2 *SSLCert) bool {
	if s1 == s2 {
		return true
	}
	if s1 == nil || s2 == nil {
		return false
	}
	if !s1.ExpireTime.Equal(s2.ExpireTime) {
		return false
	}

	for _, cn1 := range s1.CN {
		found := false
		for _, cn2 := range s2.CN {
			if cn1 == cn2 {
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
