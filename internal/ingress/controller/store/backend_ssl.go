/*
Copyright 2015 The Kubernetes Authors.

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

package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/glog"

	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress"
	"github.com/kong/kubernetes-ingress-controller/internal/net/ssl"
)

// syncSecret keeps in sync Secrets used by Ingress rules with the files on
// disk to allow copy of the content of the secret to disk to be used
// by external processes.
func (s k8sStore) syncSecret(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	glog.V(3).Infof("starting syncing of secret %v", key)

	// TODO: getPemCertificate should not write to disk to avoid unnecessary overhead
	cert, err := s.getPemCertificate(key)
	if err != nil {
		glog.Warningf("error obtaining PEM from secret %v: %v", key, err)
		return
	}

	// create certificates and add or update the item in the store
	cur, err := s.GetLocalSSLCert(key)
	if err == nil {
		if cur.Equal(cert) {
			// no need to update
			return
		}
		glog.Infof("updating secret %v in the local store", key)
		s.sslStore.Update(key, cert)
		// this update must trigger an update
		// (like an update event from a change in Ingress)
		s.sendDummyEvent()
		return
	}

	glog.Infof("adding secret %v to the local store", key)
	s.sslStore.Add(key, cert)
	// this update must trigger an update
	// (like an update event from a change in Ingress)
	s.sendDummyEvent()
}

// getPemCertificate receives a secret, and creates a ingress.SSLCert as return.
// It parses the secret and verifies if it's a keypair
func (s k8sStore) getPemCertificate(secretName string) (*ingress.SSLCert, error) {
	secret, err := s.listers.Secret.ByKey(secretName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving secret %v: %v", secretName, err)
	}

	cert, okcert := secret.Data[apiv1.TLSCertKey]
	key, okkey := secret.Data[apiv1.TLSPrivateKeyKey]

	// namespace/secretName -> namespace-secretName
	nsSecName := strings.Replace(secretName, "/", "-", -1)

	var sslCert *ingress.SSLCert
	if okcert && okkey {
		if cert == nil {
			return nil, fmt.Errorf("secret %v has no 'tls.crt'", secretName)
		}
		if key == nil {
			return nil, fmt.Errorf("secret %v has no 'tls.key'", secretName)
		}
		if secret.Data["ca.crt"] != nil {
			glog.Warningf("found 'ca.crt' for secret %v, which is not supported by Kong Ingress Controller", secretName)
		}

		sc := bytes.NewBuffer(cert).String()
		sc = strings.TrimSpace(sc)

		sk := bytes.NewBuffer(key).String()
		sk = strings.TrimSpace(sk)

		sslCert, err = ssl.AddOrUpdateCertAndKey(nsSecName, []byte(sc), []byte(sk))
		if err != nil {
			return nil, fmt.Errorf("unexpected error creating pem file: %v", err)
		}

		glog.V(3).Infof("found 'tls.crt' and 'tls.key', configuring %v as a TLS Secret (CN: %v)", secretName, sslCert.CN)
	} else {
		return nil, fmt.Errorf("no keypair could be found in %v", secretName)
	}

	sslCert.Name = secret.Name
	sslCert.Namespace = secret.Namespace
	sslCert.ID = fmt.Sprintf("%v", secret.GetUID())

	return sslCert, nil
}

// sendDummyEvent sends a dummy event to trigger an update
// This is used in when a secret change
func (s *k8sStore) sendDummyEvent() {
	s.updateCh.In() <- Event{
		Type: UpdateEvent,
		Obj: &extensions.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy",
				Namespace: "dummy",
			},
		},
	}
}
