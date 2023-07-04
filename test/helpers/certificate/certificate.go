package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type SelfSignedCeritificateOptions struct {
	CommonName string
	DNSNames   []string
}

type SelfSignedCeritificateOptionsDecorator func(SelfSignedCeritificateOptions) SelfSignedCeritificateOptions

func WithCommonName(commonName string) SelfSignedCeritificateOptionsDecorator {
	return func(opts SelfSignedCeritificateOptions) SelfSignedCeritificateOptions {
		opts.CommonName = commonName
		return opts
	}
}

func WithDNSNames(dnsNames ...string) SelfSignedCeritificateOptionsDecorator {
	return func(opts SelfSignedCeritificateOptions) SelfSignedCeritificateOptions {
		opts.DNSNames = append(opts.DNSNames, dnsNames...)
		return opts
	}
}

// GenerateSelfSignedCert generates a tls.Certificate struct to be used in TLS client/listener configurations.
func GenerateSelfSignedCert(t *testing.T, decorators ...SelfSignedCeritificateOptionsDecorator) tls.Certificate {
	t.Helper()
	// Generate a new RSA private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate RSA key")

	options := SelfSignedCeritificateOptions{
		CommonName: "",
		DNSNames:   []string{},
	}

	for _, decorator := range decorators {
		options = decorator(options)
	}

	// Create a self-signed X.509 certificate.
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Kong HQ"},
			Country:       []string{"US"},
			Province:      []string{"California"},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"150 Spear Street, Suite 1600"},
			PostalCode:    []string{"94105"},
			CommonName:    options.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		DNSNames:              options.DNSNames,
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err, "failed to create x509 certificate")

	// Create a tls.Certificate from the generated private key and certificate.
	certificate := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return certificate
}

// GenerateSelfSignedCertPEMFormat generates self-signed certificate
// and returns certificate and key in PEM format.
func GenerateSelfSignedCertPEMFormat(t *testing.T, decorators ...SelfSignedCeritificateOptionsDecorator) (cert []byte, key []byte) {
	t.Helper()
	tlsCert := GenerateSelfSignedCert(t, decorators...)

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsCert.Certificate[0],
	}

	privateKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	require.True(t, ok, "Private Key should be convertible to *rsa.PrivateKey")
	keyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock)
}
