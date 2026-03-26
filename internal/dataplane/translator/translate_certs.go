package translator

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
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

// getCertAlgorithm returns the public key algorithm of a PEM-encoded certificate.
func getCertAlgorithm(certPEM string) (x509.PublicKeyAlgorithm, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return x509.UnknownPublicKeyAlgorithm, fmt.Errorf("failed to decode certificate PEM block")
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return x509.UnknownPublicKeyAlgorithm, fmt.Errorf("failed to parse certificate: %w", err)
	}
	return parsed.PublicKeyAlgorithm, nil
}

// verifyCertSANsMatch returns an error if the two PEM-encoded certificates do not share
// the same Subject CN and DNS SANs.
func verifyCertSANsMatch(cert1PEM, cert2PEM string) error {
	cn1, dns1, err := parseCertCNAndDNSSANs(cert1PEM)
	if err != nil {
		return err
	}
	cn2, dns2, err := parseCertCNAndDNSSANs(cert2PEM)
	if err != nil {
		return err
	}
	if cn1 != cn2 {
		return fmt.Errorf("CN mismatch: %q != %q", cn1, cn2)
	}
	sorted1 := slices.Clone(dns1)
	sorted2 := slices.Clone(dns2)
	slices.Sort(sorted1)
	slices.Sort(sorted2)
	if !slices.Equal(sorted1, sorted2) {
		return fmt.Errorf("DNS SAN mismatch: %v != %v", sorted1, sorted2)
	}
	return nil
}

func parseCertCNAndDNSSANs(certPEM string) (string, []string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", nil, fmt.Errorf("failed to decode certificate PEM block")
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	return parsed.Subject.CommonName, parsed.DNSNames, nil
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
		gwc, err := s.GetGatewayClass(string(gateway.Spec.GatewayClassName))
		if err != nil {
			logger.Error(err, "Failed to get GatewayClass for Gateway, skipping", "gateway", gateway.Name, "gateway_class", gateway.Spec.GatewayClassName)
			continue
		}

		// Skip the gateway when the gateway's GatewayClass is not controlled by the KIC instance.
		if gwc.Spec.ControllerName != gatewayapi.GatewayController(t.gatewayControllerName) {
			continue
		}

		statuses := lo.SliceToMap(gateway.Status.Listeners, func(status gatewayapi.ListenerStatus) (gatewayapi.SectionName, gatewayapi.ListenerStatus) {
			return status.Name, status
		})

		for _, listener := range gateway.Spec.Listeners {
			status, ok := statuses[listener.Name]
			if !ok {
				logger.V(logging.DebugLevel).Info("Listener missing status information",
					"gateway", gateway.Name,
					"listener", listener.Name,
					"listener_protocol", listener.Protocol,
					"listener_port", listener.Port,
				)
				continue
			}

			// Check if listener is marked as programmed when the gateway's GatewayClass has the "Unmanaged" annotation.
			// If the GatewayClass does not have the annotation, the gateway is considered to be managed by other components (for example Kong Operator),
			// so we do not check the "Programmed" condition before extracting the certificate from the listener
			// to prevent unexpected deletion of certificates when the instance is managed by Kong Operator.
			if annotations.ExtractUnmanagedGatewayClassMode(gwc.Annotations) != "" &&
				!util.CheckCondition(
					status.Conditions,
					util.ConditionType(gatewayapi.ListenerConditionProgrammed),
					util.ConditionReason(gatewayapi.ListenerReasonProgrammed),
					metav1.ConditionTrue,
					gateway.Generation,
				) {
				continue
			}

			if listener.TLS != nil {
				numRefs := len(listener.TLS.CertificateRefs)
				if numRefs > 2 {
					t.registerTranslationFailure("listener '%s' has more than two certificateRefs, at most two are supported", gateway)
					continue
				}

				if numRefs == 2 {
					// fetch both secrets
					ref0, ref1 := listener.TLS.CertificateRefs[0], listener.TLS.CertificateRefs[1]
					ns0, ns1 := gateway.Namespace, gateway.Namespace
					if ref0.Namespace != nil {
						ns0 = string(*ref0.Namespace)
					}
					if ref1.Namespace != nil {
						ns1 = string(*ref1.Namespace)
					}

					secret0, err := s.GetSecret(ns0, string(ref0.Name))
					if err != nil {
						logger.Error(err, "Failed to fetch secret",
							"gateway", gateway.Name,
							"listener", listener.Name,
							"secret_name", string(ref0.Name),
							"secret_namespace", ns0,
						)
						continue
					}
					secret1, err := s.GetSecret(ns1, string(ref1.Name))
					if err != nil {
						logger.Error(err, "Failed to fetch secret",
							"gateway", gateway.Name,
							"listener", listener.Name,
							"secret_name", string(ref1.Name),
							"secret_namespace", ns1,
						)
						continue
					}

					cert0, key0, err := getCertFromSecret(secret0)
					if err != nil {
						t.registerTranslationFailure("failed to construct certificate from secret", secret0, gateway)
						continue
					}
					cert1, key1, err := getCertFromSecret(secret1)
					if err != nil {
						t.registerTranslationFailure("failed to construct certificate from secret", secret1, gateway)
						continue
					}

					// verify each certificate uses a supported algorithm (RSA or ECDSA)
					algo0, err := getCertAlgorithm(cert0)
					if err != nil {
						t.registerTranslationFailure("failed to detect certificate algorithm", secret0, gateway)
						continue
					}
					algo1, err := getCertAlgorithm(cert1)
					if err != nil {
						t.registerTranslationFailure("failed to detect certificate algorithm", secret1, gateway)
						continue
					}
					if (algo0 != x509.RSA && algo0 != x509.ECDSA) || (algo1 != x509.RSA && algo1 != x509.ECDSA) {
						t.registerTranslationFailure("listener '%s' has certificateRef with unsupported algorithm; only RSA and ECDSA are supported", gateway)
						continue
					}
					// verify the two certificates use different algorithms so Kong can select by client support
					if algo0 == algo1 {
						t.registerTranslationFailure("listener '%s' has two certificateRefs with the same algorithm; one must be RSA and one ECDSA", gateway)
						continue
					}

					// verify both certificates cover the same CN and DNS SANs so SNI selection remains unambiguous
					if err := verifyCertSANsMatch(cert0, cert1); err != nil {
						t.registerTranslationFailure("listener '%s' has certificateRefs with mismatched CN/SANs: "+err.Error(), gateway)
						continue
					}

					// Kong stores RSA in cert/key and ECDSA in cert_alt/key_alt
					var rsaCert, rsaKey, ecdsaCert, ecdsaKey string
					var primarySecret *corev1.Secret
					if algo0 == x509.RSA {
						rsaCert, rsaKey = cert0, key0
						ecdsaCert, ecdsaKey = cert1, key1
						primarySecret = secret0
					} else {
						rsaCert, rsaKey = cert1, key1
						ecdsaCert, ecdsaKey = cert0, key0
						primarySecret = secret1
					}

					// determine the SNI
					hostname := "*"
					if listener.Hostname != nil {
						hostname = string(*listener.Hostname)
					}

					certs = append(certs, certWrapper{
						identifier: rsaCert + rsaKey + ecdsaCert + ecdsaKey,
						cert: kong.Certificate{
							ID:      kong.String(string(primarySecret.UID)),
							Cert:    kong.String(rsaCert),
							Key:     kong.String(rsaKey),
							CertAlt: kong.String(ecdsaCert),
							KeyAlt:  kong.String(ecdsaKey),
							Tags:    util.GenerateTagsForObject(primarySecret),
						},
						CreationTimestamp: primarySecret.CreationTimestamp,
						snis:              []string{hostname},
					})
				} else if numRefs == 1 {
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
			}

			// although we use current in the end, we only warn/exclude on new ones here. SNIs already in the slice
			// have already been vetted by some previous iteration and /are/ in the seen list, but they're in the seen
			// list because the current we retrieved from certsSeen added them
			for _, sni := range cw.snis {
				if seen, ok := snisSeen[sni]; !ok {
					snisSeen[sni] = *current.cert.ID
					current.cert.SNIs = append(current.cert.SNIs, kong.String(sni))
				} else if seen != *current.cert.ID {
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
			idSet.certIDs = append(idSet.certIDs, *cw.cert.ID)
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
