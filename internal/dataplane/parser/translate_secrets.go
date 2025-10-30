package parser

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// getCACerts translates CA certificates Secrets to kong.CACertificates. It ensures every certificate's structure and
// validity. It skips Secrets that do not contain a valid certificate and reports translation failures for them.
func (p *Parser) getCACerts() []kong.CACertificate {
	caCertSecrets, err := p.storer.ListCACerts()
	if err != nil {
		p.logger.WithError(err).Error("failed to list CA certs")
		return nil
	}

	var caCerts []kong.CACertificate
	for _, certSecret := range caCertSecrets {
		idBytes, ok := certSecret.Data["id"]
		if !ok {
			p.registerTranslationFailure("invalid CA certificate: missing 'id' field in data", certSecret)
			continue
		}
		secretID := string(idBytes)

		caCert, err := toKongCACertificate(certSecret, secretID)
		if err != nil {
			relatedObjects := getPluginsAssociatedWithCACertSecret(secretID, p.storer)
			relatedObjects = append(relatedObjects, certSecret.DeepCopy())
			p.registerTranslationFailure(fmt.Sprintf("invalid CA certificate: %s", err), relatedObjects...)
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

	rest := bytes.TrimSpace(caCertbytes)
	pemBlocks := make([][]byte, 0, 1)

	for len(rest) > 0 {
		block, r := pem.Decode(rest)
		if block == nil {
			return kong.CACertificate{}, errors.New("invalid PEM block")
		}
		if block.Type != "CERTIFICATE" {
			return kong.CACertificate{}, errors.New("invalid PEM block type")
		}
		pemBlocks = append(pemBlocks, block.Bytes)
		rest = bytes.TrimSpace(r)
	}

	if len(pemBlocks) == 0 {
		return kong.CACertificate{}, errors.New("invalid PEM block")
	}
	if len(pemBlocks) > 1 {
		return kong.CACertificate{}, errors.New("multiple PEM certificates found")
	}

	x509Cert, err := x509.ParseCertificate(pemBlocks[0])
	if err != nil {
		return kong.CACertificate{}, errors.New("failed to parse certificate")
	}
	if !x509Cert.IsCA {
		return kong.CACertificate{}, errors.New("certificate is missing the 'CA' basic constraint")
	}
	if time.Now().After(x509Cert.NotAfter) {
		return kong.CACertificate{}, errors.New("expired")
	}

	// Re-encode a single clean CERTIFICATE PEM block to ensure only one is sent to Kong.
	singlePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: pemBlocks[0]})

	return kong.CACertificate{
		ID:   kong.String(secretID),
		Cert: kong.String(string(singlePEM)),
		Tags: util.GenerateTagsForObject(certSecret),
	}, nil
}

func getPluginsAssociatedWithCACertSecret(secretID string, storer store.Storer) []client.Object {
	refersToSecret := func(pluginConfig apiextensionsv1.JSON) bool {
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
