package parser

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

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
		idBytes, ok := certSecret.Data["id"]
		if !ok {
			log.Errorf("skipping synchronisation, invalid CA certificate: missing 'id' field in data")
			continue
		}
		secretID := string(idBytes)

		caCert, err := toKongCACertificate(certSecret, secretID)
		if err != nil {
			logWithAffectedPlugins(log, plugins, secretID).WithFields(logrus.Fields{
				"secret_name":      certSecret.Name,
				"secret_namespace": certSecret.Namespace,
			}).WithError(err).Error("skipping synchronisation, invalid CA certificate")
			continue
		}

		caCerts = append(caCerts, caCert)
	}

	return caCerts
}

func toKongCACertificate(certSecret *v1.Secret, secretID string) (kong.CACertificate, error) {
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
