//go:build conformance_tests
// +build conformance_tests

package conformance

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"

	testutils "github.com/kong/kubernetes-ingress-controller/v2/test/internal/util"
)

var (
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")
	ingressClass    = "kong-conformance-tests"

	env   environments.Environment
	ctx   context.Context
	admin *url.URL
)

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	kongAddon := kong.NewBuilder().WithControllerDisabled().WithProxyAdminServiceTypeLoadBalancer().Build()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)
	useExistingClusterIfPresent(builder)

	var err error
	env, err = builder.Build(ctx)
	exitOnErr(err)
	defer func() {
		if existingCluster == "" {
			exitOnErr(env.Cleanup(ctx))
		}
	}()

	fmt.Println("INFO: waiting for cluster and addons to be ready")
	exitOnErr(<-env.WaitForReady(ctx))

	fmt.Println("INFO: deploying CRDs")
	exitOnErr(testutils.DeployCRDsForCluster(ctx, env.Cluster()))

	fmt.Println("INFO: gathering the kong proxy admin URL")
	admin, err = kongAddon.ProxyAdminURL(ctx, env.Cluster())
	exitOnErr(err)

	fmt.Println("INFO: starting the controller manager")
	args := []string{
		fmt.Sprintf("--ingress-class=%s", ingressClass),
		fmt.Sprintf("--admission-webhook-cert=%s", testutils.KongSystemServiceCert),
		fmt.Sprintf("--admission-webhook-key=%s", testutils.KongSystemServiceKey),
		fmt.Sprintf("--admission-webhook-listen=%s:%d", testutils.AdmissionWebhookListenHost, testutils.AdmissionWebhookListenPort),
		"--profiling",
		"--dump-config",
		"--log-level=trace",
		"--debug-log-reduce-redundancy",
		"--feature-gates=Gateway=true",
		fmt.Sprintf("--kong-admin-url=%s", admin.String()),
	}
	exitOnErr(testutils.DeployControllerManagerForCluster(ctx, env.Cluster(), args...))

	code := m.Run()

	if existingCluster == "" {
		exitOnErr(env.Cleanup(ctx))
	}

	os.Exit(code)
}

func useExistingClusterIfPresent(builder *environments.Builder) {
	if existingCluster != "" {
		parts := strings.Split(existingCluster, ":")
		if len(parts) != 2 {
			exitOnErr(fmt.Errorf("%s is not a valid value for KONG_TEST_CLUSTER", existingCluster))
		}
		if parts[0] != "kind" {
			exitOnErr(fmt.Errorf("%s is not a supported cluster type for this test suite yet", parts[0]))
		}
		cluster, err := kind.NewFromExisting(parts[1])
		exitOnErr(err)
		fmt.Printf("INFO: using existing kind cluster for test (name: %s)\n", parts[1])
		builder.WithExistingCluster(cluster)
	} else {
		fmt.Println("INFO: creating new kind cluster for conformance tests")
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
