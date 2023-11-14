//go:build integration_tests

package isolated

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func SkipIfRouterNotExpressions(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
	flavor := testenv.KongRouterFlavor()
	if flavor != string(dpconf.RouterFlavorExpressions) {
		t.Skipf("skiping, %q router flavor specified via TEST_KONG_ROUTER_FLAVOR env but %q is required", flavor, "expressions")
	}
	return ctx
}
