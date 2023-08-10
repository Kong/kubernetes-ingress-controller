package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

type SelfSignedCertificateOptions struct {
	CommonName string
	DNSNames   []string
}

type SelfSignedCertificateOptionsDecorator func(SelfSignedCertificateOptions) SelfSignedCertificateOptions

func WithCommonName(commonName string) SelfSignedCertificateOptionsDecorator {
	return func(opts SelfSignedCertificateOptions) SelfSignedCertificateOptions {
		opts.CommonName = commonName
		return opts
	}
}

func WithDNSNames(dnsNames ...string) SelfSignedCertificateOptionsDecorator {
	return func(opts SelfSignedCertificateOptions) SelfSignedCertificateOptions {
		opts.DNSNames = append(opts.DNSNames, dnsNames...)
		return opts
	}
}

// MustGenerateSelfSignedCert generates a tls.Certificate struct to be used in TLS client/listener configurations.
func MustGenerateSelfSignedCert(decorators ...SelfSignedCertificateOptionsDecorator) tls.Certificate {
	// Generate a new RSA private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("failed to generate RSA key: %s", err))
	}

	options := SelfSignedCertificateOptions{
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
	if err != nil {
		panic(fmt.Sprintf("failed to create x509 certificate: %s", err))
	}

	// Create a tls.Certificate from the generated private key and certificate.
	certificate := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return certificate
}

// MustGenerateSelfSignedCertPEMFormat generates self-signed certificate
// and returns certificate and key in PEM format.
func MustGenerateSelfSignedCertPEMFormat(decorators ...SelfSignedCertificateOptionsDecorator) (cert []byte, key []byte) {
	tlsCert := MustGenerateSelfSignedCert(decorators...)

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsCert.Certificate[0],
	}

	privateKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		panic("Private Key should be convertible to *rsa.PrivateKey")
	}
	keyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock)
}

var kongSystemServiceCert, kongSystemServiceKey = MustGenerateSelfSignedCertPEMFormat(
	WithCommonName(fmt.Sprintf("*.%s.svc", consts.ControllerNamespace)),
	WithDNSNames(fmt.Sprintf("*.%s.svc", consts.ControllerNamespace)),
)

// GetKongSystemSelfSignedCerts returns the self-signed certificate and key
// with CN=*.<controllerNamespace>.svc and subjectAltName=DNS:*.<controllerNamespace>.svc.
func GetKongSystemSelfSignedCerts() (cert []byte, key []byte) {
	return kongSystemServiceCert, kongSystemServiceKey
}
