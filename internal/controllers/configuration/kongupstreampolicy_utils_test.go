package configuration

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakectrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

func TestEnforceKongUpstreamPolicyStatus(t *testing.T) {
	const (
		policyName        = "test-policy"
		anotherPolicyName = "another-test-policy"
		testNamespace     = "test"
	)

	testCases := []struct {
		name                             string
		kongupstreamPolicy               kongv1beta1.KongUpstreamPolicy
		inputObjects                     []client.Object
		expectedkongUpstreamPolicyStatus gatewayapi.PolicyStatus
		updated                          bool
	}{
		{
			name: "2 services referencing the same policy, all accepted. Status update.",
			kongupstreamPolicy: kongv1beta1.KongUpstreamPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      policyName,
					Namespace: testNamespace,
				},
			},
			inputObjects: []client.Object{
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Now(),
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: metav1.Now().Add(10 * time.Second),
						},
					},
				},
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httpRoute",
						Namespace: testNamespace,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-1",
											},
										},
									},
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedkongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-1"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
						},
					},
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-2"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
						},
					},
				},
			},
			updated: true,
		},
		{
			name: "2 services referencing the same policy, all accepted. No status update.",
			kongupstreamPolicy: kongv1beta1.KongUpstreamPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      policyName,
					Namespace: testNamespace,
				},
				Status: gatewayapi.PolicyStatus{
					Ancestors: []gatewayapi.PolicyAncestorStatus{
						{
							AncestorRef: gatewayapi.ParentReference{
								Group:     lo.ToPtr(gatewayapi.Group("core")),
								Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
								Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
								Name:      gatewayapi.ObjectName("svc-1"),
							},
							ControllerName: gatewaycontroller.GetControllerName(),
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.PolicyConditionAccepted),
									Status: metav1.ConditionTrue,
									Reason: string(gatewayapi.PolicyReasonAccepted),
								},
							},
						},
						{
							AncestorRef: gatewayapi.ParentReference{
								Group:     lo.ToPtr(gatewayapi.Group("core")),
								Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
								Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
								Name:      gatewayapi.ObjectName("svc-2"),
							},
							ControllerName: gatewaycontroller.GetControllerName(),
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.PolicyConditionAccepted),
									Status: metav1.ConditionTrue,
									Reason: string(gatewayapi.PolicyReasonAccepted),
								},
							},
						},
					},
				},
			},
			inputObjects: []client.Object{
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Now(),
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: metav1.Now().Add(10 * time.Second),
						},
					},
				},
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httpRoute",
						Namespace: testNamespace,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-1",
											},
										},
									},
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedkongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-1"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
						},
					},
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-2"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
						},
					},
				},
			},
			updated: false,
		},
		{
			name: "2 services referencing different policies, policy with conflict",
			kongupstreamPolicy: kongv1beta1.KongUpstreamPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      policyName,
					Namespace: testNamespace,
				},
			},
			inputObjects: []client.Object{
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Now(),
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: anotherPolicyName,
						},
						CreationTimestamp: metav1.Time{
							Time: metav1.Now().Add(10 * time.Second),
						},
					},
				},
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httpRoute",
						Namespace: testNamespace,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-1",
											},
										},
									},
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedkongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-1"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionFalse,
								Reason: string(gatewayapi.PolicyReasonConflicted),
							},
						},
					},
				},
			},
			updated: true,
		},
	}

	assert.NoError(t, kongv1beta1.AddToScheme(scheme.Scheme))
	assert.NoError(t, gatewayv1.AddToScheme(scheme.Scheme))

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			objectsToAdd := append(tc.inputObjects, &tc.kongupstreamPolicy)
			fakeClient := fakectrlruntimeclient.
				NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithObjects(objectsToAdd...).
				WithStatusSubresource(objectsToAdd...).
				WithIndex(&corev1.Service{}, upstreamPolicyIndexKey, indexServicesOnUpstreamPolicyAnnotation).
				WithIndex(&gatewayapi.HTTPRoute{}, routeBackendRefServiceNameIndexKey, indexRoutesOnBackendRefServiceName).
				Build()

			reconciler := KongUpstreamPolicyReconciler{
				Client: fakeClient,
			}

			updated, err := reconciler.enforceKongUpstreamPolicyStatus(context.TODO(), &tc.kongupstreamPolicy)
			assert.NoError(t, err)
			assert.Equal(t, tc.updated, updated)
			newPolicy := &kongv1beta1.KongUpstreamPolicy{}
			assert.NoError(t, fakeClient.Get(context.TODO(), k8stypes.NamespacedName{
				Namespace: tc.kongupstreamPolicy.Namespace,
				Name:      tc.kongupstreamPolicy.Name,
			}, newPolicy))
			assert.Len(t, newPolicy.Status.Ancestors, len(tc.expectedkongUpstreamPolicyStatus.Ancestors))
			for i, got := range newPolicy.Status.Ancestors {
				expected := tc.expectedkongUpstreamPolicyStatus.Ancestors[i]
				assert.Equal(t, expected.ControllerName, got.ControllerName)
				assert.Equal(t, expected.AncestorRef, got.AncestorRef)
				assert.Len(t, got.Conditions, len(expected.Conditions))
				for j, gotCond := range got.Conditions {
					// we cannot assert the whole condition is equal because of the LastTransitionTime
					expectedCond := expected.Conditions[j]
					assert.Equal(t, expectedCond.Type, gotCond.Type)
					assert.Equal(t, expectedCond.Status, gotCond.Status)
					assert.Equal(t, expectedCond.Reason, gotCond.Reason)
					assert.Equal(t, expectedCond.Message, gotCond.Message)
				}
			}
		})
	}
}
