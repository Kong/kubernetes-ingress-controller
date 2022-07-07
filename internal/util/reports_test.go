package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

type TLSPair struct {
	Key, Cert string
}

var reportTestTLSCert = TLSPair{
	Cert: `-----BEGIN CERTIFICATE-----
MIIC2DCCAcACCQC32eFOsWpKojANBgkqhkiG9w0BAQsFADAuMRcwFQYDVQQDDA5z
ZWN1cmUtZm9vLWJhcjETMBEGA1UECgwKa29uZ2hxLm9yZzAeFw0xODEyMTgyMTI4
MDBaFw0xOTEyMTgyMTI4MDBaMC4xFzAVBgNVBAMMDnNlY3VyZS1mb28tYmFyMRMw
EQYDVQQKDAprb25naHEub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqhl/HSwV6PbMv+cMFU9X+HuM7QbNNPh39GKa4pkxzFgiAnuuJ4jw9V/bzsEy
S+ZIyjzo+QKB1LzmgdcX4vkdI22BjxUd9HPHdZxtv3XilbNmSk9UOl2Hh1fORJoS
7YH+VbvVwiz5lo7qKRepbg/jcKkbs6AUE0YWFygtDLTvhP2qkphQkxZ0m8qroW91
CWgI73Ar6U2W/YQBRI3+LwtsKo0p2ASDijvqxElQBgBIiyGIr0RZc5pkCJ1eQdDB
2F6XaMfpeEyBj0MxypNL4S9HHfchOt55J1KOzYnUPkQnSoxp6oEjef4Q/ZCj5BRL
EGZnTb3tbwzHZCxGtgl9KqO9pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAKQ5BX
kkBL+alERL31hsOgWgRiUMw+sPDtRS96ozUlPtVvAg9XFdpY4ldtWkxFcmBnhKzp
UewjrHkf9rR16NISwUTjlGIwaJu/ACQrY15v+r301Crq2DV+GjiUJFVuT495dp/l
0LZbt2Sh/uD+r3UNTcJpJ7jb1V0UP7FWXFj8oafsoFSgmxAPjpKQySTC54JK4AYb
QSnWu1nQLyohnrB9qLZhe2+jOQZnkKuCcWJQ5njvU6SxT3SOKE5XaOZCezEQ6IVL
U47YCCXsq+7wKWXBhKl4H2Ztk6x3HOC56l0noXWezsMfrou/kjwGuuViGnrjqelS
WQ7uVeNCUBY+l+qY
-----END CERTIFICATE-----`,
	Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCqGX8dLBXo9sy/
5wwVT1f4e4ztBs00+Hf0YprimTHMWCICe64niPD1X9vOwTJL5kjKPOj5AoHUvOaB
1xfi+R0jbYGPFR30c8d1nG2/deKVs2ZKT1Q6XYeHV85EmhLtgf5Vu9XCLPmWjuop
F6luD+NwqRuzoBQTRhYXKC0MtO+E/aqSmFCTFnSbyquhb3UJaAjvcCvpTZb9hAFE
jf4vC2wqjSnYBIOKO+rESVAGAEiLIYivRFlzmmQInV5B0MHYXpdox+l4TIGPQzHK
k0vhL0cd9yE63nknUo7NidQ+RCdKjGnqgSN5/hD9kKPkFEsQZmdNve1vDMdkLEa2
CX0qo72lAgMBAAECggEADxMTYNJ3Xp4Ap0EioQDXGv5YDul7ZiZe+xmCAHLzJtjo
qq+rT3WjZRuJr1kPzAosiT+8pdTDDMdw5jDZvRO2sV0TDksgzHk2RAYI897OpdWw
SwWcwU9oo2X0sb+1zbang5GR8BNsSxt/RQUDzu05itJx0gltvgeIDaVR2L5wO6ja
USa8OVuj/92XtIIve9OtyK9jAzgR6LQOTFrCCEv89/vmy5Bykv4Uz8s8swZmTs3v
XJmAmruHGuSLMfXk8lBRp/gVyNTi3uMsdph5AJbVKnra5TZLguEozZKbLdNUYk0p
+aAc7rxDcH2sPqa/7DwRvei9dvd5oB3VJlxGVgC8AQKBgQDfznRSSKAD15hoSDzt
cKNyhLgWAL+MD0jhHKUy3x+Z9OCvf0DVnmru5HfQKq5UfT0t8VTRPGKmOtAMD4cf
LYjIurvMvpVzQGSJfhtHQuULZTh3dfsM7xivMqSV+9txklMAakM7vGQlOQxhrScM
21Mp5LWDU6+e2pFCrQPop0IPkQKBgQDCkVE+dou2yFuJx3uytCH1yKPSy9tkdhQH
dGF12B5dq8MZZozAz5P9YN/COa9WjsNKDqWbEgLEksEQUq4t8SBjHnSV/D3x7rEF
qgwii0GETYxax6gms8nueIqWZQf+0NbX7Gc5mTqeVb7v3TrhsKr0VNMFRXXQwP2E
M/pxJq8q1QKBgQC3rH7oXLP+Ez0AMHDYSL3LKULOw/RvpMeh/9lQA6+ysTaIsP3r
kuSdhCEUVULXEiVYhBug0FcBp3jAvSmem8cLPb0Mjkim2mzoLfeDJ1JEZODPoaLU
fZEbj4tlj9oLvhOiXpMo/jaOGeCgdPN8aK86zXlt+wtBao0WVFnF4SalEQKBgQC1
uLfi2SGgs/0a8B/ORoO5ZY3s4c2lRMtsMvyb7iBeaIAuByPLKZUVABe89deXxnsL
fiaacPX41wBO2IoqCp2vNdC6DP9mKQNZQPtYgCvPAAbo+rVIgH9HpXn7AZ24FyGy
RfAbUcv3+in9KelGxZTF4zu8HqXtNXMSuOFeMT1FiQKBgF0R+IFDGHhD4nudAQvo
hncXsgyzK6QUzak6HmFji/CMZ6EU9q6A67JkiEWrYoKqIAKZ2Og8+Eucr/rDdGWc
kqlmLPBJAJeUsP/9KidBjTE5mIbn/2n089VPMBvnlt2xIcuB6+zrf2NjvlcZEyKS
Gn+T2uCyOP4a1DTUoPyoNJXo
-----END PRIVATE KEY-----`,
}

func TestMain(m *testing.M) {
	reportsHost = "localhost"
	pingInterval = 1
	tlsConf = tls.Config{InsecureSkipVerify: true} //nolint:gosec
	os.Exit(m.Run())
}

func TestReporterOnce(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
		FeatureGates: map[string]bool{
			"Knative": false,
			"Gateway": true,
		},
	}
	reporter := Reporter{
		Info:   info,
		Logger: logr.Discard(),
	}
	reqs := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := getTLSListener()
	assert.Nil(err)
	defer listener.Close()
	go runTestTLSServer(ctx, t, listener, reqs)

	reporter.once()
	got := make(map[string]string)
	for _, line := range strings.Split(reporter.serializedInfo, ";") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		k, v := parts[0], parts[1]
		got[k] = v
	}

	want := map[string]string{
		"v":               "kic.version",
		"k8sv":            "k8s.version",
		"kv":              "kong.version",
		"db":              "off",
		"id":              "6acb7447-eedf-4815-a193-d714c5108f7b",
		"hn":              "example.local",
		"feature-knative": "false",
		"feature-gateway": "true",
	}
	assert.Equal(want, got)
}

func TestReporterSendStart(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logr.Discard(),
	}

	reqs := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := getTLSListener()
	assert.Nil(err)
	defer listener.Close()
	go runTestTLSServer(ctx, t, listener, reqs)

	reporter.once()

	reporter.sendStart()

	serialized := "<14>signal=kic-start;uptime=0;v=kic.version;" +
		"k8sv=k8s.version;kv=kong.version;db=off;" +
		"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;"
	received, ok := <-reqs
	assert.True(ok)
	short := string(bytes.Trim(received, "\x00"))
	assert.Equal(len(serialized), len(short))
	assert.Equal(serialized, short)
}

func TestReporterSendPing(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logr.Discard(),
	}

	reqs := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := getTLSListener()
	assert.Nil(err)
	defer listener.Close()
	go runTestTLSServer(ctx, t, listener, reqs)

	reporter.once()

	reporter.sendPing(42)

	serialized := "<14>signal=kic-ping;uptime=42;v=kic.version;" +
		"k8sv=k8s.version;kv=kong.version;db=off;" +
		"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;"
	received, ok := <-reqs
	assert.True(ok)
	short := string(bytes.Trim(received, "\x00"))
	assert.Equal(len(serialized), len(short))
	assert.Equal(serialized, short)
}

func TestReporterRun(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logr.Discard(),
	}

	reqs := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := getTLSListener()
	assert.Nil(err)
	defer listener.Close()
	go runTestTLSServer(ctx, t, listener, reqs)

	reporter.once()
	done := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		reporter.Run(done)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serializedContent := []string{
			"<14>signal=kic-start;uptime=0;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
			"<14>signal=kic-ping;uptime=1;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
			"<14>signal=kic-ping;uptime=2;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
		}
		for _, expect := range serializedContent {
			received, ok := <-reqs
			assert.True(ok)
			short := string(bytes.Trim(received, "\x00"))
			assert.Equal(len(expect), len(short))
			assert.Equal(expect, short)
		}
		close(done)
	}()
	wg.Wait()
}

// getTLSListener builds a TLS listener using the test certificates
func getTLSListener() (net.Listener, error) {
	testCertificate, err := tls.X509KeyPair([]byte(reportTestTLSCert.Cert), []byte(reportTestTLSCert.Key))
	if err != nil {
		return nil, err
	}
	// Most TLS configurations in this project use TLS 1.2 for FIPS mode compatibility
	// This does not because it causes an unexplained test failure where the test clients do report sending data
	// but the test server receives only an EOF/connection closed without any data. This only occurs in GitHub Actions
	// As this test server is not shipped as part of the product, allowing 1.3 here does not affect FIPS compatibility
	conf := &tls.Config{
		Certificates: []tls.Certificate{
			testCertificate,
		},
		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS12,
	}
	listen, err := tls.Listen("tcp", net.JoinHostPort(reportsHost, strconv.FormatUint(uint64(reportsPort), 10)), conf)
	if err != nil {
		return nil, err
	}
	return listen, nil
}

// runTLSServer creates a new test TLS server for the reporting system. It accepts connections using the provided
// listener and sends all requests it receives over the reqs channel
func runTestTLSServer(ctx context.Context, t *testing.T, listen net.Listener, reqs chan []byte) {
	defer close(reqs)
	for {
		select {
		case <-ctx.Done():
			listen.Close()
			return
		default:
			conn, err := listen.Accept()
			if err != nil {
				// we expect "use of closed network connection" when the test ends, since it will be blocked on accept
				if errors.Is(err, net.ErrClosed) {
					return
				}
				t.Logf("could not accept TLS connection: %v", err)
				return
			}
			go handleConnection(t, reqs, conn)
		}
	}
}

func handleConnection(t *testing.T, reqs chan []byte, conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		close(reqs)
		t.Logf("could not read from connection: %s", err)
	}
	reqs <- buffer
}
