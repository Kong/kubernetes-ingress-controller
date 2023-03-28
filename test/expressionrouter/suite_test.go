//go:build expression_router_tests
// +build expression_router_tests

package expressionrouter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

var (
	env environments.Environment
	ctx context.Context
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("INFO: setting up test environment")
	kongbuilder, _, err := helpers.GenerateKongBuilder(ctx)
	kongbuilder.WithProxyEnvVar("router_flavor", "expressions")
	exitOnErr(err)
	kongAddon := kongbuilder.Build()
	builder := environments.NewBuilder().WithAddons(kongAddon)

	fmt.Println("INFO: no existing cluster found, deploying using Kubernetes In Docker (KIND)")

	builder.WithAddons(metallb.New())

	if testenv.ClusterVersion() != "" {
		var err error
		clusterVersion, err := semver.Parse(strings.TrimPrefix(testenv.ClusterVersion(), "v"))
		exitOnErr(err)

		fmt.Printf("INFO: build a new KIND cluster with version %s\n", clusterVersion.String())
		builder.WithKubernetesVersion(clusterVersion)
	}

	fmt.Println("INFO: building test environment")
	env, err = builder.Build(ctx)
	exitOnErr(err)

	cleaner := clusters.NewCleaner(env.Cluster())

	fmt.Println("INFO: waiting for cluster and addons to be ready")
	exitOnErr(<-env.WaitForReady(ctx))

	code = m.Run()

	fmt.Printf("INFO: cleaning up cluster for env %s\n", env.Name())
	if err := cleaner.Cleanup(ctx); err != nil {
		fmt.Printf("ERROR: failed cleaning up the cluster: %v\n", err)
		code = 1
	}

}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
