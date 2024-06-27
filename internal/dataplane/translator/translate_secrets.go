package translator

import (
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

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// getCACerts translates CA certificates Secrets to kong.CACertificates. It ensures every certificate's structure and
// validity. It skips Secrets that do not contain a valid certificate and reports translation failures for them.
func (t *Translator) getCACerts() []kong.CACertificate {
	caCertSecrets, err := t.storer.ListCACerts()
	if err != nil {
		t.logger.Error(err, "Failed to list CA certs")
		return nil
	}

	caCerts := make([]kong.CACertificate, 0, len(caCertSecrets))
	for _, certSecret := range caCertSecrets {
		idBytes, ok := certSecret.Data["id"]
		if !ok {
			t.registerTranslationFailure("invalid CA certificate: missing 'id' field in data", certSecret)
			continue
		}
		secretID := string(idBytes)

		caCert, err := toKongCACertificate(certSecret, secretID)
		if err != nil {
			relatedObjects := getPluginsAssociatedWithCACertSecret(secretID, t.storer)
			relatedObjects = append(relatedObjects, certSecret.DeepCopy())
			t.registerTranslationFailure(fmt.Sprintf("invalid CA certificate: %s", err), relatedObjects...)
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
