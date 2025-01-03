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

	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

type certificateOptions struct {
	CommonName        string
	DNSNames          []string
	CATrue            bool
	Expired           bool
	Usage             x509.KeyUsage
	MaxPathLen        int
	ParentCertificate *tls.Certificate
}

type SelfSignedCertificateOption func(certificateOptions) certificateOptions

func WithCommonName(commonName string) SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.CommonName = commonName
		return opts
	}
}

func WithDNSNames(dnsNames ...string) SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.DNSNames = append(opts.DNSNames, dnsNames...)
		return opts
	}
}

// WithCATrue allows to use returned certificate to sign other certificates (uses BasicConstraints extension).
func WithCATrue() SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.CATrue = true
		return opts
	}
}

// WithAlreadyExpired allows to generate an already expired certificate.
func WithAlreadyExpired() SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.Expired = true
		return opts
	}
}

// WithMaxPathLen sets the MaxPathLen constraint in the certificate.
func WithMaxPathLen(maxLen int) SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.MaxPathLen = maxLen
		return opts
	}
}

// WithParent allows to sign the certificate with a parent certificate.
func WithParent(parent tls.Certificate) SelfSignedCertificateOption {
	return func(opts certificateOptions) certificateOptions {
		opts.ParentCertificate = &parent
		return opts
	}
}

// MustGenerateCert generates a tls.Certificate struct to be used in TLS client/listener configurations.
// If no parent certificate is passed using WithParent option, the certificate is self-signed thus returned cert can be
// used as CA for it.
func MustGenerateCert(opts ...SelfSignedCertificateOption) tls.Certificate {
	// Generate a new RSA private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate RSA key: %s", err))
	}

	options := certificateOptions{}
	for _, opt := range opts {
		options = opt(options)
	}

	notBefore := time.Now()
	notAfter := notBefore.AddDate(1, 0, 0)
	if options.Expired {
		notBefore = notBefore.AddDate(-2, 0, 0)
		notAfter = notAfter.AddDate(-2, 0, 0)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(130), nil))
	if err != nil {
		panic(fmt.Sprintf("Failed to generate serial number: %s", err))
	}

	// Create a self-signed X.509 certificate.
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"Kong HQ"},
			Country:       []string{"US"},
			Province:      []string{"California"},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"150 Spear Street, Suite 1600"},
			PostalCode:    []string{"94105"},
			CommonName:    options.CommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		DNSNames:              options.DNSNames,
		BasicConstraintsValid: true,
		IsCA:                  options.CATrue,
		KeyUsage:              options.Usage,
		MaxPathLen:            options.MaxPathLen,
	}

	var (
		// If ParentCertificate is not provided, create a self-signed certificate.
		parent     = template
		signingKey = privateKey
	)
	if options.ParentCertificate != nil {
		parent = lo.Must(x509.ParseCertificate(options.ParentCertificate.Certificate[0]))
		signingKey = options.ParentCertificate.PrivateKey.(*rsa.PrivateKey)
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, signingKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create x509 certificate: %s", err))
	}

	// Create a tls.Certificate from the generated private key and certificate.
	certificate := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return certificate
}

// MustGenerateCertPEMFormat generates a certificate and returns certificate and key in PEM format.
// If no parent certificate is passed using WithParent option, the certificate is self-signed thus returned cert can be
// used as CA for it.
func MustGenerateCertPEMFormat(opts ...SelfSignedCertificateOption) (cert []byte, key []byte) {
	return CertToPEMFormat(MustGenerateCert(opts...))
}

// CertToPEMFormat converts a tls.Certificate to PEM format.
func CertToPEMFormat(tlsCert tls.Certificate) (cert []byte, key []byte) {
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

var kongSystemServiceCert, kongSystemServiceKey = MustGenerateCertPEMFormat(
	WithCommonName(fmt.Sprintf("*.%s.svc", consts.ControllerNamespace)),
	WithDNSNames(fmt.Sprintf("*.%s.svc", consts.ControllerNamespace)),
)

// GetKongSystemSelfSignedCerts returns the self-signed certificate and key
// with CN=*.<controllerNamespace>.svc and subjectAltName=DNS:*.<controllerNamespace>.svc.
func GetKongSystemSelfSignedCerts() (cert []byte, key []byte) {
	return kongSystemServiceCert, kongSystemServiceKey
}
