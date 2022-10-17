package parser

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func getCACerts(log logrus.FieldLogger, storer store.Storer, plugins []kongstate.Plugin) []kong.CACertificate {
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
		log = logWithAffectedPlugins(log, plugins, secretID)

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

func logWithAffectedPlugins(log logrus.FieldLogger, plugins []kongstate.Plugin, secretID string) logrus.FieldLogger {
	affectedPlugins := getPluginsAssociatedWithCACertSecret(plugins, secretID)
	return log.WithField("affected_plugins", affectedPlugins)
}

func getPluginsAssociatedWithCACertSecret(plugins []kongstate.Plugin, secretID string) []string {
	refersToSecret := func(pluginConfig map[string]interface{}) bool {
		caCertReferences, ok := pluginConfig["ca_certificates"].([]string)
		if !ok {
			return false
		}

		for _, reference := range caCertReferences {
			if reference == secretID {
				return true
			}
		}
		return false
	}

	var affectedPlugins []string
	for _, p := range plugins {
		if refersToSecret(p.Config) && p.Name != nil {
			affectedPlugins = append(affectedPlugins, *p.Name)
		}
	}

	return affectedPlugins
}
