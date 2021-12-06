package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned"
)

func TestGatewayValidationWebhook(t *testing.T) {
	ns, cleanup := namespace(t)
	defer cleanup()

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("TODO: webhook tests are only supported on KIND based environments right now")
	}

	closer, err := ensureWebhookService()
	assert.NoError(t, err)
	defer closer()

	closer, err = ensureAdmissionRegistration(
		"kong-validations-gateway",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"gateway.networking.k8s.io"},
					APIVersions: []string{"v1alpha2"},
					Resources:   []string{"gateways"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	assert.NoError(t, err, "creating webhook config")
	defer closer()

	waitForWebhookService(t)

	t.Log("creating a gatewayclass to verify gateway validation")
	gatewayc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gatewayClass, err = gatewayc.GatewayV1alpha2().GatewayClasses().Create(ctx, gatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up gatewayclass %s", gatewayClass.Name)
		assert.NoError(t, gatewayc.GatewayV1alpha2().GatewayClasses().Delete(ctx, gatewayClass.Name, metav1.DeleteOptions{}))
	}()

	t.Log("creating an invalid gateway to verify that validation fails")
	gateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
			// the missing annotations here make the gateway invalid
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	_, err = gatewayc.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gateway, metav1.CreateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required annotation")

	t.Log("verifying that we don't validate a gateway that belongs to another controller")
	gateway.Spec.GatewayClassName = gatewayv1alpha2.ObjectName("nonexistentclass")
	gateway, err = gatewayc.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gateway, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that if we update an invalid gateway to be supported by our controller, validation fails")
	gateway.Spec.GatewayClassName = gatewayv1alpha2.ObjectName(gatewayClass.Name)
	_, err = gatewayc.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required annotation")

	defer func() {
		t.Logf("cleaning up gateway %s", gateway.Name)
		if err := gatewayc.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()
}
