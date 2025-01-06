//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestGatewayValidationWebhook(t *testing.T) {
	ctx := context.Background()

	ns := helpers.Namespace(ctx, t, env)

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}
	ensureAdmissionRegistration(ctx, t, env.Cluster().Client(), "kong-validations-gateway", ns.Name)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	for _, tt := range []struct {
		name      string
		createdGW gatewayapi.Gateway
		patch     []byte // optional

		wantCreateErr          bool
		wantCreateErrSubstring string

		wantPatchErr          bool
		wantPatchErrSubstring string
	}{
		{
			name: "valid gateway",
			createdGW: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.GatewayClassUnmanagedKey: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: gatewayapi.ObjectName(unmanagedGatewayClassName),
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Protocol: gatewayapi.HTTPProtocolType,
						Port:     gatewayapi.PortNumber(80),
					}},
				},
			},
			wantCreateErr: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, gotCreateErr := gatewayClient.GatewayV1().Gateways(ns.Name).Create(ctx, &tt.createdGW, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Error(t, gotCreateErr)
				require.Contains(t, gotCreateErr.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, gotCreateErr)
			}

			if len(tt.patch) > 0 {
				_, gotUpdateErr := gatewayClient.GatewayV1().Gateways(ns.Name).Patch(ctx, tt.createdGW.Name, k8stypes.MergePatchType, tt.patch, metav1.PatchOptions{})
				if tt.wantPatchErr {
					require.Error(t, gotUpdateErr)
					require.Contains(t, gotUpdateErr.Error(), tt.wantPatchErrSubstring)
				} else {
					require.NoError(t, gotUpdateErr)
				}
			}

			if err := gatewayClient.GatewayV1().Gateways(ns.Name).Delete(ctx, tt.createdGW.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}
		})
	}
}
