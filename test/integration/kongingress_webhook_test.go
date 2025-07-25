//go:build integration_tests

package integration

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/net"

	configurationv1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1"
	"github.com/kong/kubernetes-configuration/v2/pkg/clientset"
	"github.com/kong/kubernetes-configuration/v2/pkg/clientset/scheme"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestKongIngressValidationWebhook(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t.Context(), t, expressions)

	ctx := t.Context()
	ns, _ := helpers.Setup(ctx, t, env)

	ensureAdmissionRegistration(ctx, t, env.Cluster().Client(), "kong-validations-kongingress", ns.Name)

	kongClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Run("when deprecated fields are populated warnings are returned", func(t *testing.T) {
		kongIngress := &configurationv1.KongIngress{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "kong-ingress-validation-",
			},
			// Upstream field is deprecated, expecting warning for it.
			Upstream: &configurationv1.KongIngressUpstream{
				HashOn: lo.ToPtr("none"),
			},
		}

		result := kongClient.ConfigurationV1().RESTClient().Post().
			Namespace(ns.Name).
			Resource("kongingresses").
			VersionedParams(&metav1.CreateOptions{}, scheme.ParameterCodec).
			Body(kongIngress).
			Do(ctx)

		assert.NoError(t, result.Error())
		expectedWarnings := []string{
			"configuration.konghq.com/v1 KongIngress is deprecated",
			"'upstream' is DEPRECATED and will be removed in a future version. Use a KongUpstreamPolicy resource instead.",
		}
		receivedWarnings := lo.Map(result.Warnings(), func(item net.WarningHeader, _ int) string {
			return item.Text
		})
		assert.ElementsMatch(t, expectedWarnings, receivedWarnings)
	})
}
