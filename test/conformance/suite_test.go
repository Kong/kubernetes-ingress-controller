//go:build conformance_tests
// +build conformance_tests

package conformance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
)

var (
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")
	ingressClass    = "kong-conformance-tests"

	env environments.Environment
	ctx context.Context
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

		// For some reason we've ended up using kind v0.15.0 (which by default deploys k8s v1.25)
		// even though that
		// * our CI runners use ubuntu-latest (which at the time of writing this comment was ubuntu20.04)
		//   which uses kind v0.14 https://github.com/actions/runner-images/blob/main/images/linux/Ubuntu2004-Readme.md
		// * when used with ktf (which at the time of writing this comment was set to v0.19.0) we
		//   should use kind v0.14 since that what ktf has set as dependency
		//
		// With all that said, we still managed to get kind v0.15 on our CI
		// https://github.com/Kong/kubernetes-ingress-controller/runs/8211490522?check_suite_focus=true#step:5:6
		// which causes issues down the line (metallb manifests using PSP which is not available
		// in k8s v1.25+).
		builder.WithKubernetesVersion(semver.Version{
			Major: 1,
			Minor: 24,
			Patch: 4,
		})
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
