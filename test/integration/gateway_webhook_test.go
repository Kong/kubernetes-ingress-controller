//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
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
	const configResourceName = "kong-validations-gateway"
	ensureAdmissionRegistration(
		ctx,
		t,
		ns.Name,
		configResourceName,
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

	t.Log("waiting for webhook service to be connective")
	ensureWebhookServiceIsConnective(ctx, t, configResourceName)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	for _, tt := range []struct {
		name      string
		createdGW gatewayv1.Gateway
		patch     []byte // optional

		wantCreateErr          bool
		wantCreateErrSubstring string

		wantPatchErr          bool
		wantPatchErrSubstring string
	}{
		{
			name: "valid gateway",
			createdGW: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: gatewayv1.ObjectName(unmanagedGatewayClassName),
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Protocol: gatewayv1.HTTPProtocolType,
						Port:     gatewayv1.PortNumber(80),
					}},
				},
			},
			wantCreateErr: false,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, gotCreateErr := gatewayClient.GatewayV1beta1().Gateways(ns.Name).Create(ctx, &tt.createdGW, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Error(t, gotCreateErr)
				require.Contains(t, gotCreateErr.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, gotCreateErr)
			}

			if len(tt.patch) > 0 {
				_, gotUpdateErr := gatewayClient.GatewayV1beta1().Gateways(ns.Name).Patch(ctx, tt.createdGW.Name, k8stypes.MergePatchType, tt.patch, metav1.PatchOptions{})
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
