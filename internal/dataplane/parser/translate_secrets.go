package parser

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/sirupsen/logrus"
)

func getCACerts(log logrus.FieldLogger, storer store.Storer) []kong.CACertificate {
	caCertSecrets, err := storer.ListCACerts()
	if err != nil {
		log.WithError(err).Error("failed to list CA certs")
		return nil
	}

	var caCerts []kong.CACertificate
	for _, certSecret := range caCertSecrets {
		secretName := certSecret.Namespace + "/" + certSecret.Name

		idbytes, idExists := certSecret.Data["id"]
		log = log.WithFields(logrus.Fields{
			"secret_name":      secretName,
			"secret_namespace": certSecret.Namespace,
		})
		if !idExists {
			log.Errorf("invalid CA certificate: missing 'id' field in data")
			continue
		}
		secretID := string(idbytes)
		log = logWithAffectedPlugins(log, storer, secretID)

		caCertbytes, certExists := certSecret.Data["cert"]
		if !certExists {
			log.Errorf("invalid CA certificate: missing 'cert' field in data")
			continue
		}

		pemBlock, _ := pem.Decode(caCertbytes)
		if pemBlock == nil {
			log.Errorf("invalid CA certificate: invalid PEM block")
			continue
		}
		x509Cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			log.WithError(err).Errorf("invalid CA certificate: failed to parse certificate")
			continue
		}
		if !x509Cert.IsCA {
			log.WithError(err).Errorf("invalid CA certificate: certificate is missing the 'CA' basic constraint")
			continue
		}
		if time.Now().After(x509Cert.NotAfter) {
			log.WithError(err).Errorf("expired CA certificate")
			continue
		}

		caCerts = append(caCerts, kong.CACertificate{
			ID:   kong.String(secretID),
			Cert: kong.String(string(caCertbytes)),
		})
	}

	return caCerts
}

func logWithAffectedPlugins(log logrus.FieldLogger, storer store.Storer, secretID string) logrus.FieldLogger {
	affectedPlugins := getPluginsAssociatedWithSecret(storer, secretID)
	return log.WithField("affected_plugins", affectedPlugins)
}

func getPluginsAssociatedWithSecret(storer store.Storer, secretID string) []string {
	var affectedPlugins []string

	clusterPlugins, err := storer.ListGlobalKongClusterPlugins()
	if err != nil {
		return nil
	}
	for _, p := range clusterPlugins {
		if pluginConfigRefersToSecret(p.Config, secretID) {
			affectedPlugins = append(affectedPlugins, p.Name)
		}
	}

	plugins, err := storer.ListGlobalKongPlugins()
	if err != nil {
		return affectedPlugins
	}
	for _, p := range plugins {
		if pluginConfigRefersToSecret(p.Config, secretID) {
			affectedPlugins = append(affectedPlugins, fmt.Sprintf("%s/%s", p.Namespace, p.Name))
		}
	}
	return affectedPlugins
}

func pluginConfigRefersToSecret(cfg apiextensionsv1.JSON, secretID string) bool {
	pluginConfig := struct {
		CACertificates []string `json:"ca_certificates,omitempty"`
	}{}

	if err := json.Unmarshal(cfg.Raw, &pluginConfig); err != nil {
		return false
	}

	for _, reference := range pluginConfig.CACertificates {
		if reference == secretID {
			return true
		}
	}
	return false
}
