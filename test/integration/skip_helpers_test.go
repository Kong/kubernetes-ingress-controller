//go:build integration_tests

package integration

import (
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
)

type routerFlavor string

const (
	traditional           routerFlavor = "traditional"
	traditionalCompatible routerFlavor = "traditional_compatible"
	expressions           routerFlavor = "expressions"
)

func skipTestForRouterFlavor(t *testing.T, flavor ...routerFlavor) {
	t.Helper()
	routerFlavor := routerFlavor(eventuallyGetKongRouterFlavor(t, proxyAdminURL))
	for _, f := range flavor {
		if routerFlavor == f {
			t.Skipf("router flavor:%q skipping", f)
		}
	}
}

// Expression router is not supported for some objects and features.
// For example, KongIngress is not supported by intention;
// TCPRoute is not supported because Kong (< 3.4) does not support expression router on stream proxy.
// When the test case depends on the object or feature not supported, we skip it if expression router is used.
func skipTestForExpressionRouter(t *testing.T) {
	t.Helper()
	skipTestForRouterFlavor(t, expressions)
}

func skipTestForNonKindCluster(t *testing.T) {
	t.Helper()
	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("this test is only available on KIND clusters currently")
	}
}
