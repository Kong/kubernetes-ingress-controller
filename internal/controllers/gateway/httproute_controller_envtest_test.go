//go:build envtest
// +build envtest

package gateway_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/envtest"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers"
)

func init() {
	if err := gatewayv1beta1.Install(scheme.Scheme); err != nil {
		panic(err)
	}
}

func TestHTTPRouteReconcilerProperlyReactsToReferenceGrant(t *testing.T) {
	t.Parallel()

	const (
		waitDuration = 5 * time.Second
		tickDuration = 100 * time.Millisecond
	)

	cfg := envtest.Setup(t, scheme.Scheme)
	var client ctrlclient.Client
	{
		var err error
		client, err = ctrlclient.New(cfg, ctrlclient.Options{
			Scheme: scheme.Scheme,
		})
		require.NoError(t, err)
	}

	// In tests below we use a deferred cancel to stop the manager and not wait
	// for its timeout.

	testcases := []struct {
		name       string
		reconciler *gateway.HTTPRouteReconciler
	}{
		{
			name: "with ReferenceGrant enabled",
			reconciler: &gateway.HTTPRouteReconciler{
				Client:               client,
				EnableReferenceGrant: true,
			},
		},
		{
			name: "with ReferenceGrant disabled",
			reconciler: &gateway.HTTPRouteReconciler{
				Client:               client,
				EnableReferenceGrant: false,
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ns := envtest.CreateNamespace(ctx, t, client)
			nsRoute := envtest.CreateNamespace(ctx, t, client)

			svc := corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "backend-1",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "http",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			}
			require.NoError(t, client.Create(ctx, &svc))

			tc.reconciler.DataplaneClient = gateway.DataplaneMock{}
			envtest.StartReconciler(ctx, t, client.Scheme(), cfg, tc.reconciler, nil)

			gwc := gatewayv1beta1.GatewayClass{
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: gateway.GetControllerName(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						"konghq.com/gatewayclass-unmanaged": "placeholder",
					},
				},
			}
			require.NoError(t, client.Create(ctx, &gwc))
			t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

			gw := gatewayv1beta1.Gateway{
				Spec: gatewayv1beta1.GatewaySpec{
					GatewayClassName: gatewayv1beta1.ObjectName(gwc.Name),
					Listeners: []gatewayv1beta1.Listener{
						{
							Name:     gatewayv1beta1.SectionName("http"),
							Port:     gatewayv1beta1.PortNumber(80),
							Protocol: gatewayv1beta1.HTTPProtocolType,
							AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
								Namespaces: &gatewayv1beta1.RouteNamespaces{
									From: lo.ToPtr(gatewayv1beta1.NamespacesFromAll),
								},
							},
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      uuid.NewString(),
				},
			}
			require.NoError(t, client.Create(ctx, &gw))

			gwOld := gw.DeepCopy()
			gw.Status = gatewayv1beta1.GatewayStatus{
				Addresses: []gatewayv1beta1.GatewayAddress{
					{
						Type:  lo.ToPtr(gatewayv1beta1.IPAddressType),
						Value: "10.0.0.1",
					},
				},
				Conditions: []metav1.Condition{
					{
						Type:               "Ready",
						Status:             metav1.ConditionTrue,
						Reason:             "Ready",
						LastTransitionTime: metav1.Now(),
						ObservedGeneration: gw.Generation,
					},
					{
						Type:               "Accepted",
						Status:             metav1.ConditionTrue,
						Reason:             "Accepted",
						LastTransitionTime: metav1.Now(),
						ObservedGeneration: gw.Generation,
					},
					{
						Type:               "Programmed",
						Status:             metav1.ConditionTrue,
						Reason:             "Programmed",
						LastTransitionTime: metav1.Now(),
						ObservedGeneration: gw.Generation,
					},
				},
				Listeners: []gatewayv1beta1.ListenerStatus{
					{
						Name: gatewayv1beta1.SectionName("http"),
						Conditions: []metav1.Condition{
							{
								Type:               "Accepted",
								Status:             metav1.ConditionTrue,
								Reason:             "Accepted",
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               "Ready",
								Status:             metav1.ConditionTrue,
								Reason:             "Ready",
								LastTransitionTime: metav1.Now(),
							},
						},
						SupportedKinds: []gatewayv1beta1.RouteGroupKind{
							{
								Group: lo.ToPtr(gatewayv1beta1.Group(gatewayv1beta1.GroupVersion.Group)),
								Kind:  "HTTPRoute",
							},
						},
					},
				},
			}
			require.NoError(t, client.Status().Patch(ctx, &gw, ctrlclient.MergeFrom(gwOld)))

			route := gatewayv1beta1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HTTPRoute",
					APIVersion: "v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: nsRoute.Name,
					Name:      uuid.NewString(),
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name:      gatewayv1beta1.ObjectName(gw.Name),
							Namespace: lo.ToPtr(gatewayv1beta1.Namespace(ns.Name)),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						BackendRefs: builder.NewHTTPBackendRef("backend-1").WithNamespace(ns.Name).ToSlice(),
					}},
				},
			}
			require.NoError(t, client.Create(ctx, &route))

			nn := types.NamespacedName{
				Namespace: route.GetNamespace(),
				Name:      route.GetName(),
			}

			t.Logf("verifying that HTTPRoute has ResolvedRefs set to Status False and Reason RefNotPermitted")
			if !assert.Eventually(t,
				helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
					metav1.Condition{
						Type:   "ResolvedRefs",
						Status: "False",
						Reason: "RefNotPermitted",
					},
				),
				waitDuration, tickDuration,
			) {
				t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
			}

			rg := gatewayv1beta1.ReferenceGrant{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      uuid.NewString(),
				},
				Spec: gatewayv1beta1.ReferenceGrantSpec{
					From: []gatewayv1beta1.ReferenceGrantFrom{
						{
							Group:     gatewayv1beta1.Group(gatewayv1beta1.GroupVersion.Group),
							Kind:      "HTTPRoute",
							Namespace: gatewayv1beta1.Namespace(nsRoute.Name),
						},
					},
					To: []gatewayv1beta1.ReferenceGrantTo{
						{
							Group: "",
							Kind:  "Service",
						},
					},
				},
			}
			require.NoError(t, client.Create(ctx, &rg))
			if tc.reconciler.EnableReferenceGrant {
				t.Logf("verifying that HTTPRoute gets accepted by HTTPRouteReconciler after relevant ReferenceGrant gets created")
				if !assert.Eventually(t,
					helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
						metav1.Condition{
							Type:   "ResolvedRefs",
							Status: "True",
							Reason: "ResolvedRefs",
						},
						metav1.Condition{
							Type:   "Accepted",
							Status: "True",
							Reason: "Accepted",
						},
						// Programmed condition requires a bit more work with mocks.
						// It's set only when KubernetesObjectReports are enabled in the underlying
						// dataplane client and then it relies on what's returned by
						// dataplane client in KubernetesObjectConfigurationStatus().
						// This can be done but it's not the main focus of this test.
						// Related: https://github.com/Kong/kubernetes-ingress-controller/issues/3793
					),
					waitDuration, tickDuration,
				) {
					t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
				}
			} else {
				t.Logf("verifying that HTTPRoute's status doesn't change after relevant ReferenceGrant gets created")

				if !assert.Eventually(t,
					helpers.HTTPRouteEventuallyNotContainsConditions(ctx, t, client, nn,
						metav1.Condition{
							Type:   "ResolvedRefs",
							Status: "True",
							Reason: "ResolvedRefs",
						},
						metav1.Condition{
							Type:   "Accepted",
							Status: "True",
							Reason: "Accepted",
						},
					),
					waitDuration, tickDuration,
				) {
					t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
				}
			}

			require.NoError(t, client.Delete(ctx, &rg))
			t.Logf("verifying that HTTPRoute gets its ResolvedRefs condition to Status False and Reason RefNotPermitted when relevant ReferenceGrant gets deleted")

			if !assert.Eventually(t,
				helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
					metav1.Condition{
						Type:   "ResolvedRefs",
						Status: "False",
						Reason: "RefNotPermitted",
					},
				),
				waitDuration, tickDuration,
			) {
				t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
			}
		})
	}
}

func printHTTPRoutesConditions(ctx context.Context, client ctrlclient.Client, nn types.NamespacedName) string {
	var route gateway.HTTPRoute
	err := client.Get(ctx, ctrlclient.ObjectKey{Namespace: nn.Namespace, Name: nn.Name}, &route)
	if err != nil {
		return fmt.Sprintf("failed to get HTTPRoute %s/%s when trying to print its conditions", nn.Namespace, nn.Name)
	}

	if len(route.Status.Parents) == 0 {
		return fmt.Sprintf("HTTPRoute %s/%s has no parents in Status", nn.Namespace, nn.Name)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("HTTPRoute %s/%s has the following Parents in Status:", nn.Namespace, nn.Name))
	for _, p := range route.Status.Parents {
		if p.ParentRef.Namespace != nil {
			_, _ = sb.WriteString(fmt.Sprintf("\nParent %s/%s: ", *p.ParentRef.Namespace, string(p.ParentRef.Name)))
		} else {
			_, _ = sb.WriteString(fmt.Sprintf("\nParent %s: ", string(p.ParentRef.Name)))
		}
		for _, c := range p.Conditions {
			s := fmt.Sprintf(
				"\n\tcondition: Type:%s, Status:%s, Reason:%s, ObservedGeneration:%d",
				c.Type, c.Status, c.Reason, c.ObservedGeneration,
			)
			_, _ = sb.WriteString(s)
		}
		_ = sb.WriteByte('\n')
	}
	return sb.String()
}
