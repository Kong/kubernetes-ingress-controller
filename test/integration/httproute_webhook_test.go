//go:build integration_tests
// +build integration_tests

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
)

func TestHTTPRouteValidationWebhook(t *testing.T) {
	ns, cleanup := namespace(t)
	defer cleanup()

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}

	pathMatchRegex := gatewayv1alpha2.PathMatchRegularExpression

	closer, err := ensureAdmissionRegistration(
		"kong-validations-gateway",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"gateway.networking.k8s.io"},
					APIVersions: []string{"v1alpha2"},
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

	waitForWebhookService(t)

	t.Log("creating a managed gatewayclass")
	gatewayc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gatewaycontroller.ControllerName,
		},
	}
	gatewayClass, err = gatewayc.GatewayV1alpha2().GatewayClasses().Create(ctx, gatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up gatewayclass %s", gatewayClass.Name)
		if err := gatewayc.GatewayV1alpha2().GatewayClasses().Delete(ctx, gatewayClass.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("creating a managed gateway")
	gateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation: "true",
			},
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
	gateway, err = gatewayc.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gateway, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up gateway %s", gateway.Name)
		if err := gatewayc.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("creating an unmanaged gatewayclass")
	require.NoError(t, err)
	unmanagedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: "example.com/gateway-controller",
		},
	}
	unmanagedGatewayClass, err = gatewayc.GatewayV1alpha2().GatewayClasses().Create(ctx, unmanagedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up gatewayclass %s", unmanagedGatewayClass.Name)
		if err := gatewayc.GatewayV1alpha2().GatewayClasses().Delete(ctx, unmanagedGatewayClass.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("creating an unmanaged gateway")
	unmanagedGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(unmanagedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	unmanagedGateway, err = gatewayc.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, unmanagedGateway, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up gateway %s", unmanagedGateway.Name)
		if err := gatewayc.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, unmanagedGateway.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	for _, tt := range []struct {
		name                   string
		route                  *gatewayv1alpha2.HTTPRoute
		wantCreateErr          bool
		wantCreateErrSubstring string
	}{
		{
			name: "a valid httproute linked to a managed gateway passes validation",
			route: &gatewayv1alpha2.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentRef{{
							Namespace: (*gatewayv1alpha2.Namespace)(&gateway.Namespace),
							Name:      gatewayv1alpha2.ObjectName(gateway.Name),
						}},
					},
				},
			},
			wantCreateErr: false,
		},
		{
			name: "an httproute linked to a non-existent gateway fails validation",
			route: &gatewayv1alpha2.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentRef{{
							Namespace: (*gatewayv1alpha2.Namespace)(&gateway.Namespace),
							Name:      gatewayv1alpha2.ObjectName("fake-gateway"),
						}},
					},
				},
			},
			wantCreateErr:          true,
			wantCreateErrSubstring: `Gateway.gateway.networking.k8s.io \"fake-gateway\" not found`,
		},
		{
			name: "an invalid httproute will pass validation if it's not linked to a managed controller (it's not ours)",
			route: &gatewayv1alpha2.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Rules: []gatewayv1alpha2.HTTPRouteRule{{
						Matches: []gatewayv1alpha2.HTTPRouteMatch{{
							Path: &gatewayv1alpha2.HTTPPathMatch{
								Type: &pathMatchRegex, // this route is invalid because we don't support regex path matches (yet)
							},
						}},
					}},
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentRef{{
							Namespace: (*gatewayv1alpha2.Namespace)(&gateway.Namespace),
							Name:      gatewayv1alpha2.ObjectName(unmanagedGateway.Name),
						}},
					},
				},
			},
			wantCreateErr: false, // shouldn't fail because it's considered unmanaged
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gatewayc.GatewayV1alpha2().HTTPRoutes(ns.Name).Create(ctx, tt.route, metav1.CreateOptions{})
			if tt.wantCreateErr {
				require.Contains(t, err.Error(), tt.wantCreateErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
