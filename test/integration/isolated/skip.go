//go:build integration_tests

package isolated

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func SkipIfRouterNotExpressions(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
	flavor := testenv.KongRouterFlavor()
	// TODO refactor the router flavor consts.
	if flavor != "expressions" {
		t.Skipf("skiping, %q router flavor specified via TEST_KONG_ROUTER_FLAVOR env but %q is required", flavor, "expressions")
	}
	// fg := os.Getenv("KONG_CONTROLLER_FEATURE_GATES")
	// if !strings.Contains(fg, featuregates.ExpressionRoutesFeature) {
	// 	t.Skipf("skiping, 'expressions' router flavor specified but via TEST_KONG_ROUTER_FLAVOR env but %q is required", flavor, "expressions")
	// }
	return ctx
}
