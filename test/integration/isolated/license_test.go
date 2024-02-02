package isolated

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
	"github.com/stretchr/testify/require"
)

func TestKongLicense(t *testing.T) {
	f := features.
		New("essentials").
		WithLabel(testlabels.Kind, testlabels.KindKongLicense).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("prepare clients", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)

			kongClients, err := clientset.NewForConfig(cluster.Config())
			require.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, kongClients)

			gatewayClient, err := gatewayclient.NewForConfig(cluster.Config())
			require.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			return ctx
		})
	tenv.Test(t, f.Feature())
}
