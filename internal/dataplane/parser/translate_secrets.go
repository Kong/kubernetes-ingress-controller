package parser

import (
	"context"
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
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// getCACerts translates CA certificates Secrets to kong.CACertificates. It ensures every certificate's structure and
// validity. It skips Secrets that do not contain a valid certificate and reports translation failures for them.
func (p *Parser) getCACerts(ctx context.Context) []kong.CACertificate {
	caCertSecrets, err := p.storer.ListCACerts(ctx)
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
			relatedObjects := getPluginsAssociatedWithCACertSecret(ctx, secretID, p.storer)
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
		Tags: util.GenerateTagsForObject(certSecret),
	}, nil
}

func getPluginsAssociatedWithCACertSecret(ctx context.Context, secretID string, storer store.Storer) []client.Object {
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
	for _, p := range storer.ListKongPlugins(ctx) {
		if refersToSecret(p.Config) {
			affectedPlugins = append(affectedPlugins, p.DeepCopy())
		}
	}
	for _, p := range storer.ListKongClusterPlugins(ctx) {
		if refersToSecret(p.Config) {
			affectedPlugins = append(affectedPlugins, p.DeepCopy())
		}
	}

	return affectedPlugins
}
