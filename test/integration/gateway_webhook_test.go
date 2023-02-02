//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestGatewayValidationWebhook(t *testing.T) {
	ctx := context.Background()

	ns := helpers.Namespace(ctx, t, env)

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}

	closer, err := ensureAdmissionRegistration(ctx,
		"kong-validations-gateway",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"gateway.networking.k8s.io"},
					APIVersions: []string{"v1beta1"},
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

	t.Log("waiting for webhook service to be connective")
	err = waitForWebhookServiceConnective(ctx, "kong-validations-gateway")
	require.NoError(t, err)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	for _, tt := range []struct {
		name      string
		createdGW gatewayv1beta1.Gateway
		patch     []byte // optional

		wantCreateErr          bool
		wantCreateErrSubstring string

		wantPatchErr          bool
		wantPatchErrSubstring string
	}{
		{
			name: "valid gateway",
			createdGW: gatewayv1beta1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1beta1.GatewaySpec{
					GatewayClassName: gatewayv1beta1.ObjectName(unmanagedGatewayClassName),
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Protocol: gatewayv1beta1.HTTPProtocolType,
						Port:     gatewayv1beta1.PortNumber(80),
					}},
				},
			},
			wantCreateErr: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, gotCreateErr := gatewayClient.GatewayV1beta1().Gateways(ns.Name).Create(ctx, &tt.createdGW, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Error(t, gotCreateErr)
				require.Contains(t, gotCreateErr.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, gotCreateErr)
			}

			if len(tt.patch) > 0 {
				_, gotUpdateErr := gatewayClient.GatewayV1beta1().Gateways(ns.Name).Patch(ctx, tt.createdGW.Name, types.MergePatchType, tt.patch, metav1.PatchOptions{})
				if tt.wantPatchErr {
					require.Error(t, gotUpdateErr)
					require.Contains(t, gotUpdateErr.Error(), tt.wantPatchErrSubstring)
				} else {
					require.NoError(t, gotUpdateErr)
				}
			}

			if err := gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, tt.createdGW.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}
		})
	}
}
