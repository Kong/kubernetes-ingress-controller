package translator

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// getCACerts translates CA certificates Secrets to kong.CACertificates. It ensures every certificate's structure and
// validity. It skips Secrets that do not contain a valid certificate and reports translation failures for them.
func (t *Translator) getCACerts() []kong.CACertificate {
	caCertSecrets, caCertConfigMaps, err := t.storer.ListCACerts()
	if err != nil {
		t.logger.Error(err, "Failed to list CA certs")
		return nil
	}

	// intermediateT stores the object and its data for easier handling.
	// intermediateT.data can come either from secret.Data (map[string]byte) or configmap.Data (map[string]string).
	type intermediateT struct {
		obj  client.Object
		data map[string]string
	}

	caCertsData := lo.Map(caCertSecrets, func(secret *corev1.Secret, _ int) intermediateT {
		return intermediateT{
			obj: secret.DeepCopy(),
			data: lo.MapValues(secret.Data, func(v []byte, _ string) string {
				return string(v)
			}),
		}
	})
	caCertsData = append(caCertsData, lo.Map(caCertConfigMaps, func(configmap *corev1.ConfigMap, _ int) intermediateT {
		return intermediateT{obj: configmap.DeepCopy(), data: configmap.Data}
	})...)

	caCerts := make([]kong.CACertificate, 0, len(caCertSecrets)+len(caCertConfigMaps))
	for _, caCertData := range caCertsData {
		secretID, ok := caCertData.data["id"]
		if !ok {
			t.registerTranslationFailure("invalid secret CA certificate: missing 'id' field in data", caCertData.obj)
			continue
		}

		// Allow the certificate key to be named either "cert" or "ca.crt".
		caCertStr, certExists := caCertData.data["cert"]
		if !certExists {
			caCertStr, certExists = caCertData.data["ca.crt"]
			if !certExists {
				relatedObjects := getPluginsAssociatedWithCACertSecret(secretID, t.storer)
				relatedObjects = append(relatedObjects, caCertData.obj)
				t.registerTranslationFailure(fmt.Sprintf(`invalid secret CA certificate %s/%s, neither "cert" nor "ca.crt" key exist`, caCertData.obj.GetNamespace(), caCertData.obj.GetName()), relatedObjects...)
				continue
			}
		}

		caCert, err := toKongCACertificate([]byte(caCertStr), caCertData.obj, secretID)
		if err != nil {
			relatedObjects := getPluginsAssociatedWithCACertSecret(secretID, t.storer)
			relatedObjects = append(relatedObjects, caCertData.obj)
			t.registerTranslationFailure(fmt.Sprintf("invalid secret CA certificate: %s", err), relatedObjects...)
			continue
		}

		caCerts = append(caCerts, caCert)
	}

	return caCerts
}

func toKongCACertificate(caCertBytes []byte, object client.Object, secretID string) (kong.CACertificate, error) {
	pemBlock, _ := pem.Decode(caCertBytes)
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
		Cert: kong.String(string(caCertBytes)),
		Tags: util.GenerateTagsForObject(object),
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
