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

package ssl

import (
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/cert/triple"
)

// generateRSACerts generates a self signed certificate using a self generated ca
func generateRSACerts(host string) (*triple.KeyPair, *triple.KeyPair, error) {
	ca, err := triple.NewCA("self-sign-ca")
	if err != nil {
		return nil, nil, err
	}

	key, err := certutil.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create a server private key: %v", err)
	}

	config := certutil.Config{
		CommonName: host,
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	cert, err := certutil.NewSignedCert(config, key, ca.Cert, ca.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to sign the server certificate: %v", err)
	}

	return &triple.KeyPair{
		Key:  key,
		Cert: cert,
	}, ca, nil
}

func TestAddOrUpdateCertAndKey(t *testing.T) {

	cert, _, err := generateRSACerts("echoheaders")
	if err != nil {
		t.Fatalf("unexpected error creating SSL certificate: %v", err)
	}

	name := fmt.Sprintf("test-%v", time.Now().UnixNano())

	c := certutil.EncodeCertPEM(cert.Cert)
	k := certutil.EncodePrivateKeyPEM(cert.Key)

	ngxCert, err := AddOrUpdateCertAndKey(name, c, k)
	if err != nil {
		t.Fatalf("unexpected error checking SSL certificate: %v", err)
	}

	if len(ngxCert.CN) == 0 {
		t.Fatalf("expected at least one cname but none returned")
	}

	if ngxCert.CN[0] != "echoheaders" {
		t.Fatalf("expected cname echoheaders but %v returned", ngxCert.CN[0])
	}
}

func TestCACert(t *testing.T) {

	cert, _, err := generateRSACerts("echoheaders")
	if err != nil {
		t.Fatalf("unexpected error creating SSL certificate: %v", err)
	}

	name := fmt.Sprintf("test-%v", time.Now().UnixNano())

	c := certutil.EncodeCertPEM(cert.Cert)
	k := certutil.EncodePrivateKeyPEM(cert.Key)

	_, err = AddOrUpdateCertAndKey(name, c, k)
	if err != nil {
		t.Fatalf("unexpected error checking SSL certificate: %v", err)
	}
}

func TestGetFakeSSLCert(t *testing.T) {
	k, c := GetFakeSSLCert()
	if len(k) == 0 {
		t.Fatalf("expected a valid key")
	}
	if len(c) == 0 {
		t.Fatalf("expected a valid certificate")
	}
}
