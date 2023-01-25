package test

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	// KongSystemServiceCert is a testing TLS certificate with SAN *.kong-system.svc.
	//
	// created with:
	//
	// openssl req \
	//   -x509 \
	//   -newkey rsa:4096 \
	//   -sha256 \
	//   -days 3560 \
	//   -nodes \
	//   -keyout tls.key \
	//   -out tls.crt \
	//   -subj '/CN=*.kong-system.svc' \
	//   -extensions san \
	//   -config <( \
	//     echo '[req]'; \
	//     echo 'distinguished_name=req'; \
	//     echo '[san]'; \
	//     echo 'subjectAltName=DNS:*.kong-system.svc')
	KongSystemServiceCert = `-----BEGIN CERTIFICATE-----
MIIE5jCCAs6gAwIBAgIUCAQsA6yH5M6/LgmSg/89y4NPB8UwDQYJKoZIhvcNAQEL
BQAwHDEaMBgGA1UEAwwRKi5rb25nLXN5c3RlbS5zdmMwHhcNMjIwMTA1MTY0MjE4
WhcNMzExMDA1MTY0MjE4WjAcMRowGAYDVQQDDBEqLmtvbmctc3lzdGVtLnN2YzCC
AiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL6fCBkdOlOvmGoXxW7ZsMQD
nENJ3BRCeco4Jezb0SJYG/epXH3meXtih0E5trF1LWGVOfGPgS/32NCAGnPVLerc
hQxAVL/prh7+/PjYm+2WvxDgoLqtI+iOda7ouXPnl54Qtp0rB906CEyhO+PH+lzm
DoKHRPGNPCVG2AOb0xhHWlouKLufNExyQ4+ny/qpSCcepOKu3ME1HJpiI8kMkwrj
asw/gw8LCxZwAdXP2y6U5cvynr7B9lGAgsqHifF1xHM1W7kqn1+E7EsnAH93gRtq
uRtapCX3FjH831aRCGKCkJRgfrMgszX0eoM0BkCWlarLd4Za8zQvZAiCf3z2uEvn
M/7myyyVNeGlzAQOnoerTq8lKBa+GAvHG6DhQ2S+aqiGKqV+yZRtBnqEhcb7LCIi
CpSFnq7AZHYg3mahMYbnw8iXUUdnNO70ONnq9OWuzAOQSWfUQuhY7BEQquNFFN3l
0XohYbLB64W+ZCBQ63AXoktIXPeOdRl9RzywnCCcyZEaiuOdYDYos54AStLCVM7u
cnv1e/63tyhdzmmD53U0KwSOjCxB7rpA0vFCVE1BncHDpupJQ65j2c9bFgV4ehQ3
jJb4uZ2glENFaBqKjoZtu1uCsuEbH0vIy27nXrrsFu30NFHt7sVU4+OCR3T/acOJ
+QKOygBYaXpEokXKbXI5AgMBAAGjIDAeMBwGA1UdEQQVMBOCESoua29uZy1zeXN0
ZW0uc3ZjMA0GCSqGSIb3DQEBCwUAA4ICAQA1jxc3q7vzmG/8NxPrCzlXKvcoTHTJ
ysVglzfvLTbJ1oKmLzVyvX43QWDD8rhyJr7W3agsLwHTdzs6hHlSl+DPWYK02FvA
x8/CTmycGcySpj1NCO9JIif+h01mt7bC65dP/tcWBPAxWq9vQ4PzL/s7/7RXgCxm
9qmY5Ri+spmKcMzU8ekG+IDZfCbJ5qPNdnCmTd1G2tBAeDkYeg+iS7rcnfIEVppd
0nIBhbEuMU6P7WXzUb7+nFTPN7R46PZ/u+lifkCPt2DIniJgI38W7M8Aspu513gU
n6L2rjxfwv1ZJbpTTwjws4fki54i9FhZJcPPWRV2t1zcT3k2nMH8OQXXjbo49uEc
jW3BvufX1zc+NXnN+5HIZtxNksPglY0MHpdFoD5GMHTrauEANGpoYG6vDGwA8NQd
8Z5l+IV0nGD0I7Jz+Q83bZtCezk9s+RlhO6TpTNM1O3yTkOWctCCVJQQH3wKSrkW
YXT9Zzp47E4sI/MKK4DCl4hPsZCdIsBjreBgawBcEb+xczWS2euFuAc+MHdenDS2
ftDKUoUKvhFhetW5VsTVGESWP26LhmQ6oaz6eJPC26skKq+w2rAeZ2EU8fhBwzCX
vv2SWWk6xtzxMyD2w3IXPlUCzpkTwkT6tC7qf5RYA/y+Qr81pbbp6q3OEQ51BaKI
nyC7Cdnomxtx6A==
-----END CERTIFICATE-----`
	// KongSystemServiceKey is the private key for the testing certificate kongSystemServiceCert.
	KongSystemServiceKey = `-----BEGIN PRIVATE KEY-----
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
ctFsgXhf5+tDgbBZpcuTMpd3KnaDUYg=
-----END PRIVATE KEY-----`
)

// This hack is tracked in https://github.com/Kong/kubernetes-ingress-controller/issues/1613:
//
// The test process (`go test github.com/Kong/kubernetes-ingress-controller/test/integration/...`) serves the webhook
// endpoints to be consumed by the apiserver (so that the tests can apply a ValidatingWebhookConfiguration and test
// those validations).
//
// In order to make that possible, we needed to allow the apiserver (that gets spun up by the test harness) to access
// the system under test (which runs as a part of the `go test` process).
// Below, we're making an audacious assumption that the host running the `go test` process is either:
//
// - a direct Docker host on the default bridge, and that the apiserver is running within a context
// (such as KIND running on that same docker bridge), from which 172.17.0.1 routes to the host OR
// - a Colima host, and that the apiserver is running within a docker container hosted by Colima
// from which 192.168.5.2 routes to the host (https://github.com/abiosoft/colima/issues/220)
//
// This works if the test runs against a KIND cluster, and does not work against cloud providers (like GKE).

var AdmissionWebhookListenHost = admissionWebhookListenHost()

const (
	AdmissionWebhookListenPort = 49023

	colimaHostAddress                 = "192.168.5.2"
	defaultDockerBridgeNetworkGateway = "172.17.0.1"
)

func admissionWebhookListenHost() string {
	if isColimaHost() {
		return colimaHostAddress
	}

	return defaultDockerBridgeNetworkGateway
}

func isColimaHost() bool {
	cmd := exec.Command("docker", "info", "--format", "{{.Name}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to run %q command %s\n", cmd.String(), err)
		fmt.Println(string(out))
		return false
	}

	return strings.Contains(string(out), "colima")
}
