package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress"
	"github.com/kong/kubernetes-ingress-controller/internal/ssl"
	apiv1 "k8s.io/api/core/v1"
)

func (s k8sStore) GetCertFromSecret(secretName string) (*ingress.SSLCert, error) {
	secret, err := s.listers.Secret.ByKey(secretName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving secret %v: %v", secretName, err)
	}
	cert, okcert := secret.Data[apiv1.TLSCertKey]
	key, okkey := secret.Data[apiv1.TLSPrivateKeyKey]

	if !okcert || !okkey {
		return nil, fmt.Errorf("no keypair could be found in %v", secretName)
	}

	cert = []byte(strings.TrimSpace(bytes.NewBuffer(cert).String()))
	key = []byte(strings.TrimSpace(bytes.NewBuffer(key).String()))

	sslCert := &ingress.SSLCert{
		Raw: ingress.RawSSLCert{
			Cert: cert,
			Key:  key,
		},
		ID: fmt.Sprintf("%v", secret.GetUID()),
	}
	sslCert.Namespace = secret.Namespace

	certificate, err := ssl.ParseX509Certificate(cert, key)
	if err != nil {
		return nil, err
	}
	sslCert.Certificate = certificate
	cn := ssl.ParseCommonNamesFromCert(certificate)
	sslCert.CN = cn
	sslCert.ExpireTime = certificate.NotAfter

	return sslCert, nil
}
