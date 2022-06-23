//go:build integration_tests
// +build integration_tests

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

func TestGatewayValidationWebhook(t *testing.T) {
	ns, cleanup := namespace(t)
	defer cleanup()

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}

	closer, err := ensureAdmissionRegistration(
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
	defer func() {
		assert.NoError(t, closer())
	}()

	waitForWebhookService(t)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	for _, tt := range []struct {
		name      string
		createdGW gatewayv1alpha2.Gateway
		patch     []byte // optional

		wantCreateErr          bool
		wantCreateErrSubstring string

		wantPatchErr          bool
		wantPatchErrSubstring string
	}{
		{
			name: "valid gateway",
			createdGW: gatewayv1alpha2.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation: "true",
					},
				},
				Spec: gatewayv1alpha2.GatewaySpec{
					GatewayClassName: gatewayv1alpha2.ObjectName(managedGatewayClassName),
					Listeners: []gatewayv1alpha2.Listener{{
						Name:     "http",
						Protocol: gatewayv1alpha2.HTTPProtocolType,
						Port:     gatewayv1alpha2.PortNumber(80),
					}},
				},
			},
			wantCreateErr: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, gotCreateErr := gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, &tt.createdGW, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Error(t, gotCreateErr)
				require.Contains(t, gotCreateErr.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, gotCreateErr)
			}

			if len(tt.patch) > 0 {
				_, gotUpdateErr := gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Patch(ctx, tt.createdGW.Name, types.MergePatchType, tt.patch, metav1.PatchOptions{})
				if tt.wantPatchErr {
					require.Error(t, gotUpdateErr)
					require.Contains(t, gotUpdateErr.Error(), tt.wantPatchErrSubstring)
				} else {
					require.NoError(t, gotUpdateErr)
				}
			}

			if err := gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, tt.createdGW.Name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		})
	}
}
