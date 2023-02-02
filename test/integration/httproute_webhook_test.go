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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestHTTPRouteValidationWebhook(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}

	pathMatchRegex := gatewayv1beta1.PathMatchRegularExpression

	closer, err := ensureAdmissionRegistration(ctx,
		"kong-validations-gateway",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"gateway.networking.k8s.io"},
					APIVersions: []string{"v1beta1"},
					Resources:   []string{"httproutes"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	assert.NoError(t, err, "creating webhook config")
	defer func() {
		assert.NoError(t, closer())
	}()

	t.Log("creating a gateway client ")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("creating a managed gateway")
	managedGateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName, func(g *gatewayv1beta1.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(managedGateway)

	t.Log("creating an unmanaged gatewayclass")
	unmanagedGatewayClass, err := DeployGatewayClass(ctx, gatewayClient, uuid.NewString(), func(gc *gatewayv1beta1.GatewayClass) {
		gc.Spec.ControllerName = unsupportedControllerName
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGatewayClass)

	t.Log("creating an unmanaged gateway")
	unmanagedGateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClass.Name, func(g *gatewayv1beta1.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGateway)

	t.Log("waiting for webhook service to be connective")
	err = waitForWebhookServiceConnective(ctx, "kong-validations-gateway")
	require.NoError(t, err)

	for _, tt := range []struct {
		name                   string
		route                  *gatewayv1beta1.HTTPRoute
		wantCreateErr          bool
		wantCreateErrSubstring string
	}{
		{
			name: "a valid httproute linked to a managed gateway passes validation",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Namespace: (*gatewayv1beta1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1beta1.ObjectName(managedGateway.Name),
						}},
					},
				},
			},
			wantCreateErr: false,
		},
		{
			name: "an httproute linked to a non-existent gateway fails validation",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Namespace: (*gatewayv1beta1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
				},
			},
			wantCreateErr:          true,
			wantCreateErrSubstring: `Gateway.gateway.networking.k8s.io \"fake-gateway\" not found`,
		},
		{
			name: "an invalid httproute will pass validation if it's not linked to a managed controller (it's not ours)",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type: &pathMatchRegex, // this route is invalid because we don't support regex path matches (yet)
							},
						}},
					}},
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Namespace: (*gatewayv1beta1.Namespace)(&unmanagedGateway.Namespace),
							Name:      gatewayv1beta1.ObjectName(unmanagedGateway.Name),
						}},
					},
				},
			},
			wantCreateErr: false, // shouldn't fail because it's considered unmanaged
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Create(ctx, tt.route, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Contains(t, err.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
