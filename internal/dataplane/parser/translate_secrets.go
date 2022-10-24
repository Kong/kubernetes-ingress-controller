package parser

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

// getCACerts translates CA certificates Secrets to kong.CACertificates. It ensures every certificate's structure and
// validity. In case of violation of any validation rule, a secret gets skipped in a result and error message is logged
// with affected plugins for context.
func (p *Parser) getCACerts() []kong.CACertificate {
	log := p.logger
	caCertSecrets, err := p.storer.ListCACerts()
	if err != nil {
		log.WithError(err).Error("failed to list CA certs")
		return nil
	}

	var caCerts []kong.CACertificate
	for _, certSecret := range caCertSecrets {
		idBytes, ok := certSecret.Data["id"]
		if !ok {
			p.errorsCollector.ParsingError("invalid CA certificate: missing 'id' field in data", certSecret)
			continue
		}
		secretID := string(idBytes)

		caCert, err := toKongCACertificate(certSecret, secretID)
		if err != nil {
			affectedObjects := getPluginsAssociatedWithCACertSecret(secretID, p.storer)
			affectedObjects = append(affectedObjects, certSecret)
			p.errorsCollector.ParsingError(fmt.Sprintf("invalid CA certificate: %s", err), affectedObjects...)
			continue
		}

		caCerts = append(caCerts, caCert)
	}

	return caCerts
}

func toKongCACertificate(certSecret *corev1.Secret, secretID string) (kong.CACertificate, error) {
	caCertbytes, certExists := certSecret.Data["cert"]
	if !certExists {
		return kong.CACertificate{}, errors.New("missing 'cert' field in data")
	}
	pemBlock, _ := pem.Decode(caCertbytes)
	if pemBlock == nil {
		return kong.CACertificate{}, errors.New("invalid PEM block")
	}
	x509Cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return kong.CACertificate{}, errors.New("failed to parse certificate")
	}
	if !x509Cert.IsCA {
		return kong.CACertificate{}, errors.New("certificate is missing the 'CA' basic constraint")
	}
	if time.Now().After(x509Cert.NotAfter) {
		return kong.CACertificate{}, errors.New("expired")
	}

	return kong.CACertificate{
		ID:   kong.String(secretID),
		Cert: kong.String(string(caCertbytes)),
	}, nil
}

func getPluginsAssociatedWithCACertSecret(secretID string, storer store.Storer) []client.Object {
	refersToSecret := func(pluginConfig v1.JSON) bool {
		cfg := struct {
			CACertificates []string `json:"ca_certificates,omitempty"`
		}{}
		err := json.Unmarshal(pluginConfig.Raw, &cfg)
		if err != nil {
			return false
		}

		for _, reference := range cfg.CACertificates {
			if reference == secretID {
				return true
			}
		}
		return false
	}

	var affectedPlugins []client.Object
	for _, p := range storer.ListKongPlugins() {
		if refersToSecret(p.Config) {
			affectedPlugins = append(affectedPlugins, p.DeepCopy())
		}
	}
	for _, p := range storer.ListKongClusterPlugins() {
		if refersToSecret(p.Config) {
			affectedPlugins = append(affectedPlugins, p.DeepCopy())
		}
	}

	return affectedPlugins
}
