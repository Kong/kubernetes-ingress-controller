package translator

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func getCertFromSecret(secret *corev1.Secret) (string, string, error) {
	certData, okcert := secret.Data[corev1.TLSCertKey]
	keyData, okkey := secret.Data[corev1.TLSPrivateKeyKey]

	if !okcert || !okkey {
		return "", "", fmt.Errorf("no keypair could be found in"+
			" secret '%v/%v'", secret.Namespace, secret.Name)
	}

	cert := bytes.TrimSpace(certData)
	key := bytes.TrimSpace(keyData)

	if _, err := tls.X509KeyPair(cert, key); err != nil {
		return "", "", fmt.Errorf("parsing TLS key-pair in secret '%v/%v': %w",
			secret.Namespace, secret.Name, err)
	}

	return string(cert), string(key), nil
}

type certWrapper struct {
	identifier        string
	cert              kong.Certificate
	snis              []string
	CreationTimestamp metav1.Time
}

func (t *Translator) getGatewayCerts() []certWrapper {
	logger := t.logger
	s := t.storer
	certs := []certWrapper{}
	gateways, err := s.ListGateways()
	if err != nil {
		logger.Error(err, "Failed to list Gateways")
		return certs
	}
	for _, gateway := range gateways {
		statuses := make(map[gatewayapi.SectionName]gatewayapi.ListenerStatus, len(gateway.Status.Listeners))
		for _, status := range gateway.Status.Listeners {
			statuses[status.Name] = status
		}

		for _, listener := range gateway.Spec.Listeners {
			status, ok := statuses[listener.Name]
			if !ok {
				logger.V(util.DebugLevel).Info("Listener missing status information",
					"gateway", gateway.Name,
					"listener", listener.Name,
					"listener_protocol", listener.Protocol,
					"listener_port", listener.Port,
				)
				continue
			}

			// Check if listener is marked as programmed
			if !util.CheckCondition(
				status.Conditions,
				util.ConditionType(gatewayapi.ListenerConditionProgrammed),
				util.ConditionReason(gatewayapi.ListenerReasonProgrammed),
				metav1.ConditionTrue,
				gateway.Generation,
			) {
				continue
			}

			if listener.TLS != nil {
				if len(listener.TLS.CertificateRefs) > 0 {
					if len(listener.TLS.CertificateRefs) > 1 {
						// TODO support cert_alt and key_alt if there are 2 SecretObjectReferences
						// https://github.com/Kong/kubernetes-ingress-controller/issues/2604
						t.registerTranslationFailure("listener '%s' has more than one certificateRef, it's not supported", gateway)
						continue
					}

					// determine the Secret Namespace
					ref := listener.TLS.CertificateRefs[0]
					namespace := gateway.Namespace
					if ref.Namespace != nil {
						namespace = string(*ref.Namespace)
					}

					// retrieve the Secret and extract the PEM strings
					secret, err := s.GetSecret(namespace, string(ref.Name))
					if err != nil {
						logger.Error(err, "Failed to fetch secret",
							"gateway", gateway.Name,
							"listener", listener.Name,
							"secret_name", string(ref.Name),
							"secret_namespace", namespace,
						)
						continue
					}
					cert, key, err := getCertFromSecret(secret)
					if err != nil {
						t.registerTranslationFailure("failed to construct certificate from secret", secret, gateway)
						continue
					}

					// determine the SNI
					hostname := "*"
					if listener.Hostname != nil {
						hostname = string(*listener.Hostname)
					}

					// create a Kong certificate, wrap it in metadata, and add it to the certs slice
					certs = append(certs, certWrapper{
						identifier: cert + key,
						cert: kong.Certificate{
							ID:   kong.String(string(secret.UID)),
							Cert: kong.String(cert),
							Key:  kong.String(key),
							Tags: util.GenerateTagsForObject(secret),
						},
						CreationTimestamp: secret.CreationTimestamp,
						snis:              []string{hostname},
					})
				}
			}
		}
	}
	return certs
}

func (t *Translator) getCerts(secretsToSNIs SecretNameToSNIs) []certWrapper {
	certs := []certWrapper{}

	for secretKey, SNIs := range secretsToSNIs.secretToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := t.storer.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			t.registerTranslationFailure(fmt.Sprintf("Failed to fetch the secret (%s)", secretKey), SNIs.Parents()...)
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			causingObjects := append(SNIs.Parents(), secret)
			t.registerTranslationFailure("failed to construct certificate from secret", causingObjects...)
			continue
		}
		certs = append(certs, certWrapper{
			identifier: cert + key,
			cert: kong.Certificate{
				ID:   kong.String(string(secret.UID)),
				Cert: kong.String(cert),
				Key:  kong.String(key),
				Tags: util.GenerateTagsForObject(secret),
			},
			CreationTimestamp: secret.CreationTimestamp,
			snis:              SNIs.Hosts(),
		})
	}

	return certs
}

type certIDToMergedCertID map[string]string

type identicalCertIDSet struct {
	mergedCertID string
	certIDs      []string
}

func mergeCerts(logger logr.Logger, certLists ...[]certWrapper) ([]kongstate.Certificate, certIDToMergedCertID) {
	snisSeen := make(map[string]string)
	certsSeen := make(map[string]certWrapper)
	certIDSets := make(map[string]identicalCertIDSet)

	for _, cl := range certLists {
		for _, cw := range cl {
			current, ok := certsSeen[cw.identifier]
			if !ok {
				current = cw
			} else {
				// multiple Secrets that contain identical certificates are collapsed, because we only create one
				// Kong resource for a given cert+key pair. however, because we reuse the Secret ID and creation time
				// for the Kong resource equivalents, the selection of those needs to be deterministic to avoid
				// pointless configuration updates
				if current.CreationTimestamp.After(cw.CreationTimestamp.Time) {
					current.cert.ID = cw.cert.ID
					current.CreationTimestamp = cw.CreationTimestamp
				} else if current.CreationTimestamp.Time.Equal(cw.CreationTimestamp.Time) && (current.cert.ID == nil || *current.cert.ID > *cw.cert.ID) {
					current.cert.ID = cw.cert.ID
					current.CreationTimestamp = cw.CreationTimestamp
				}
				current.snis = append(current.snis, cw.snis...)
			}

			// although we use current in the end, we only warn/exclude on new ones here. SNIs already in the slice
			// have already been vetted by some previous iteration and /are/ in the seen list, but they're in the seen
			// list because the current we retrieved from certsSeen added them
			for _, sni := range cw.snis {
				if seen, ok := snisSeen[sni]; !ok {
					snisSeen[sni] = *current.cert.ID
					current.cert.SNIs = append(current.cert.SNIs, kong.String(sni))
				} else {
					// TODO this should really log information about the requesting Listener or Ingress-like, which is
					// what binds the SNI to a given Secret. Knowing the Secret ID isn't of great use beyond knowing
					// what cert will be served. however, the secretToSNIs input to getCerts does not provide this info
					// https://github.com/Kong/kubernetes-ingress-controller/issues/2605
					logger.Error(nil, "Same SNI requested for multiple certs, can only serve one cert",
						"served_secret_cert", seen,
						"requested_secret_cert", *current.cert.ID,
						"sni", sni)
				}
			}
			certsSeen[current.identifier] = current

			idSet := certIDSets[current.identifier]
			idSet.mergedCertID = *current.cert.ID
			idSet.certIDs = append(certIDSets[current.identifier].certIDs, *cw.cert.ID)
			certIDSets[current.identifier] = idSet

		}
	}
	res := make([]kongstate.Certificate, 0, len(certsSeen))
	for _, cw := range certsSeen {
		sort.SliceStable(cw.cert.SNIs, func(i, j int) bool {
			return strings.Compare(*cw.cert.SNIs[i], *cw.cert.SNIs[j]) < 0
		})
		res = append(res, kongstate.Certificate{
			Certificate: cw.cert,
		})
	}

	idToMergedID := certIDToMergedCertID{}
	for _, idSet := range certIDSets {
		for _, certID := range idSet.certIDs {
			idToMergedID[certID] = idSet.mergedCertID
		}
	}
	return res, idToMergedID
}
