package adminapi

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var ca *x509.Certificate
var caPEM *bytes.Buffer
var caPrivateKeyPEM *bytes.Buffer
var certPEM *bytes.Buffer
var certPrivateKeyPEM *bytes.Buffer

func TestMakeHTTPClient(t *testing.T) {

	buildTLS(t)

	var opts = HTTPClientOpts{
		true,
		"",
		"",
		caPEM.String(),
		nil,
		"",
		certPEM.String(),
		"",
		certPrivateKeyPEM.String(),
	}

	httpclient, err := MakeHTTPClient(&opts)
	if err != nil {
		t.Errorf("Fail to make the HTTP client - %s", err.Error())
	}

	assert.NotNil(t, httpclient)

	validate(t, httpclient)

}

func TestMakeHTTPClientWithFilePaths(t *testing.T) {

	buildTLS(t)

	caFile, err := ioutil.TempFile("", "ca.crt")
	if err != nil {
		t.Errorf("Fail to create CA certificate file - %s", err.Error())
	}
	caFile.Write(caPEM.Bytes())
	defer os.Remove(caFile.Name())

	certFile, err := ioutil.TempFile("", "cert.crt")
	if err != nil {
		t.Errorf("Fail to create CA certificate file - %s", err.Error())
	}
	certFile.Write(certPEM.Bytes())
	defer os.Remove(caFile.Name())

	certPrivateKeyFile, err := ioutil.TempFile("", "cert.key")
	if err != nil {
		t.Errorf("Fail to create CA certificate file - %s", err.Error())
	}
	certPrivateKeyFile.Write(certPrivateKeyPEM.Bytes())
	defer os.Remove(caFile.Name())

	var opts = HTTPClientOpts{
		true,
		"",
		caFile.Name(),
		"",
		nil,
		certFile.Name(),
		"",
		certPrivateKeyFile.Name(),
		"",
	}

	httpclient, err := MakeHTTPClient(&opts)
	if err != nil {
		t.Errorf("Fail to make the HTTP client - %s", err.Error())
	}

	assert.NotNil(t, httpclient)

	validate(t, httpclient)

}

func buildTLS(t *testing.T) (err error) {

	ca = &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization:  []string{"Kong HQ"},
			Country:       []string{"US"},
			Province:      []string{"California"},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"150 Spear Street, Suite 1600"},
			PostalCode:    []string{"94105"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 1),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Errorf("Fail to generate CA key %s", err.Error())
		return err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		t.Errorf("Fail to generate CA certificate %s", err.Error())
		return err
	}

	caPEM = new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivateKeyPEM = new(bytes.Buffer)
	pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Kong HQ"},
			Country:       []string{"US"},
			Province:      []string{"California"},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"150 Spear Street, Suite 1600"},
			PostalCode:    []string{"94105"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Errorf("Fail to generate ingress key %s", err.Error())
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivateKey.PublicKey, certPrivateKey)
	if err != nil {
		t.Errorf("Fail to generate ingress certificate %s", err.Error())
		return err
	}

	certPEM = new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivateKeyPEM = new(bytes.Buffer)
	pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivateKey),
	})

	return
}

func validate(t *testing.T, httpclient *http.Client) (err error) {

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivateKeyPEM.Bytes())
	if err != nil {
		t.Errorf("Fail to load server certificates %s", err.Error())
		return err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPEM.Bytes())

	serverTLSConf := &tls.Config{
		RootCAs:      certPool,
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAnyClientCert,
		Certificates: []tls.Certificate{serverCert},
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "...")
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()

	response, err := httpclient.Get(server.URL)
	if err != nil {
		t.Errorf("HTTP client failed to issue a GET request %s", err.Error())
		return err
	}

	// verify the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("HTTP client failed to process a GET request %s", err.Error())
		return err
	}

	body := strings.TrimSpace(string(data[:]))
	if body != "..." {
		t.Errorf("Invalid server response")
		return err
	}

	return
}
