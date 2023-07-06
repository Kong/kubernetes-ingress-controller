package test

import (
	"fmt"
	"os/exec"
	"strings"
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
