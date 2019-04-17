package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	apiv1 "k8s.io/api/core/v1"
)

func (s k8sStore) GetCertFromSecret(secretName string) (*utils.RawSSLCert, error) {
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

	return &utils.RawSSLCert{
		Cert: cert,
		Key:  key,
	}, nil
}
