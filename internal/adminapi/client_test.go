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
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeHTTPClientWithTLSOpts(t *testing.T) {
	var caPEM *bytes.Buffer
	var certPEM *bytes.Buffer
	var certPrivateKeyPEM *bytes.Buffer
	var err error

	caPEM, certPEM, certPrivateKeyPEM, err = buildTLS(t)
	if err != nil {
		t.Errorf("Fail to build TLS certificates - %s", err.Error())
	}

	opts := HTTPClientOpts{
		TLSSkipVerify:     true,
		TLSServerName:     "",
		CACertPath:        "",
		CACert:            caPEM.String(),
		Headers:           nil,
		TLSClientCertPath: "",
		TLSClientCert:     certPEM.String(),
		TLSClientKeyPath:  "",
		TLSClientKey:      certPrivateKeyPEM.String(),
	}

	httpclient, err := MakeHTTPClient(&opts)
	require.NoError(t, err)

	assert.NotNil(t, httpclient)

	err = validate(t, httpclient, caPEM, certPEM, certPrivateKeyPEM)
	require.NoError(t, err)
}

func TestMakeHTTPClientWithTLSOptsAndFilePaths(t *testing.T) {
	var caPEM *bytes.Buffer
	var certPEM *bytes.Buffer
	var certPrivateKeyPEM *bytes.Buffer
	var err error

	caPEM, certPEM, certPrivateKeyPEM, err = buildTLS(t)
	if err != nil {
		t.Errorf("Fail to build TLS certificates - %s", err.Error())
	}

	caFile, err := os.CreateTemp(os.TempDir(), "ca.crt")
	require.NoError(t, err)
	writtenBytes, err := caFile.Write(caPEM.Bytes())
	require.NoError(t, err)
	require.Equal(t, caPEM.Len(), writtenBytes)
	defer os.Remove(caFile.Name())

	certFile, err := os.CreateTemp(os.TempDir(), "cert.crt")
	require.NoError(t, err)
	writtenBytes, err = certFile.Write(certPEM.Bytes())
	require.NoError(t, err)
	require.Equal(t, certPEM.Len(), writtenBytes)
	defer os.Remove(caFile.Name())

	certPrivateKeyFile, err := os.CreateTemp(os.TempDir(), "cert.key")
	require.NoError(t, err)
	writtenBytes, err = certPrivateKeyFile.Write(certPrivateKeyPEM.Bytes())
	require.NoError(t, err)
	require.Equal(t, certPrivateKeyPEM.Len(), writtenBytes)
	defer os.Remove(caFile.Name())

	opts := HTTPClientOpts{
		TLSSkipVerify:     true,
		TLSServerName:     "",
		CACertPath:        caFile.Name(),
		CACert:            "",
		Headers:           nil,
		TLSClientCertPath: certFile.Name(),
		TLSClientCert:     "",
		TLSClientKeyPath:  certPrivateKeyFile.Name(),
		TLSClientKey:      "",
	}

	httpclient, err := MakeHTTPClient(&opts)
	require.NoError(t, err)

	assert.NotNil(t, httpclient)

	err = validate(t, httpclient, caPEM, certPEM, certPrivateKeyPEM)
	require.NoError(t, err)
}

func buildTLS(t *testing.T) (caPEM *bytes.Buffer, certPEM *bytes.Buffer, certPrivateKeyPEM *bytes.Buffer, err error) {
	var ca *x509.Certificate
	var caPrivateKeyPEM *bytes.Buffer

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
		return nil, nil, nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		t.Errorf("Fail to generate CA certificate %s", err.Error())
		return nil, nil, nil, err
	}

	caPEM = new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		t.Errorf("Fail to encode CA certificate %s", err.Error())
		return nil, nil, nil, err
	}

	caPrivateKeyPEM = new(bytes.Buffer)
	err = pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})
	if err != nil {
		t.Errorf("Fail to encode CA key %s", err.Error())
		return nil, nil, nil, err
	}

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
		return nil, nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivateKey.PublicKey, certPrivateKey)
	if err != nil {
		t.Errorf("Fail to generate ingress certificate %s", err.Error())
		return nil, nil, nil, err
	}

	certPEM = new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		t.Errorf("Fail to encode certificate %s", err.Error())
		return nil, nil, nil, err
	}

	certPrivateKeyPEM = new(bytes.Buffer)
	err = pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivateKey),
	})
	if err != nil {
		t.Errorf("Fail to encode key %s", err.Error())
		return nil, nil, nil, err
	}

	return caPEM, certPEM, certPrivateKeyPEM, nil
}

func validate(t *testing.T,
	httpclient *http.Client,
	caPEM *bytes.Buffer,
	certPEM *bytes.Buffer,
	certPrivateKeyPEM *bytes.Buffer,
) (err error) {
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
		MinVersion:   tls.VersionTLS12,
	}

	successMessage := "connection successful"
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, successMessage)
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()

	response, err := httpclient.Get(server.URL)
	if err != nil {
		t.Errorf("HTTP client failed to issue a GET request %s", err.Error())
		return err
	}
	defer response.Body.Close()

	// verify the response
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("HTTP client failed to process a GET request %s", err.Error())
		return err
	}

	body := strings.TrimSpace(string(data[:]))
	if body != successMessage {
		t.Errorf("Invalid server response")
		return err
	}

	return nil
}
