//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
)

type routerFlavor string

const (
	traditional           routerFlavor = "traditional"
	traditionalCompatible routerFlavor = "traditional_compatible"
	expressions           routerFlavor = "expressions"
)

func skipTestForRouterFlavors(ctx context.Context, t *testing.T, flavor ...routerFlavor) {
	t.Helper()
	routerFlavor := routerFlavor(eventuallyGetKongRouterFlavor(ctx, t, proxyAdminURL))
	for _, f := range flavor {
		if routerFlavor == f {
			t.Skipf("router flavor: %q for ingress: %q skipping", f, proxyAdminURL)
		}
	}
}

func skipTestForNonKindCluster(t *testing.T) {
	t.Helper()
	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("this test is only available on KIND clusters currently")
	}
}
