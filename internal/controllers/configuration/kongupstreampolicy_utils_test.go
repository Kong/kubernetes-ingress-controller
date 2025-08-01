package configuration

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakectrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	configurationv1beta1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/v2/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/scheme"
)

func TestEnforceKongUpstreamPolicyStatus(t *testing.T) {
	const (
		policyName        = "test-policy"
		anotherPolicyName = "another-test-policy"
		testNamespace     = "test"
	)

	now := metav1.Now()

	testCases := []struct {
		name                             string
		kongUpstreamPolicy               configurationv1beta1.KongUpstreamPolicy
		inputObjects                     []client.Object
		objectsConfiguredInDataPlane     bool
		expectedKongUpstreamPolicyStatus gatewayapi.PolicyStatus
		updated                          bool
	}{
		{
			name: "3 services referencing the same policy, all accepted. Status update.",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(10 * time.Second),
						},
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-3",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(5 * time.Second),
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
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: "svc-3",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			objectsConfiguredInDataPlane: true,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					// Order of ancestors should be svc-1 - svc-3 - svc-2 for the order of creationTimestamp.
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("core")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-3"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
				},
			},
			updated: true,
		},
		{
			name: "2 services referencing the same policy, all accepted. No status update.",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
								{
									Type:   string(gatewayapi.GatewayConditionProgrammed),
									Status: metav1.ConditionTrue,
									Reason: string(gatewayapi.GatewayReasonProgrammed),
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
								{
									Type:   string(gatewayapi.GatewayConditionProgrammed),
									Status: metav1.ConditionTrue,
									Reason: string(gatewayapi.GatewayReasonProgrammed),
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
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
			objectsConfiguredInDataPlane: true,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
				},
			},
			updated: false,
		},
		{
			name: "2 services in the same httproute rule referencing different policies, conflict",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: anotherPolicyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(10 * time.Second),
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
			objectsConfiguredInDataPlane: true,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionFalse,
								Reason: string(gatewayapi.GatewayReasonPending),
							},
						},
					},
				},
			},
			updated: true,
		},
		{
			name: "2 services referencing different policies in different http route rules, accepted",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: anotherPolicyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(10 * time.Second),
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
								},
							},
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
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
			objectsConfiguredInDataPlane: true,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
				},
			},
			updated: true,
		},
		{
			name: "service and kong service facade referencing the same policy, accepted",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade-1",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(10 * time.Second),
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
									builder.NewHTTPBackendRef("svc-1").Build(),
								},
							},
						},
					},
				},
			},
			objectsConfiguredInDataPlane: true,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group(incubatorv1alpha1.GroupVersion.Group)),
							Kind:      lo.ToPtr(gatewayapi.Kind(incubatorv1alpha1.KongServiceFacadeKind)),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-facade-1"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.GatewayReasonProgrammed),
							},
						},
					},
				},
			},
			updated: true,
		},
		{
			name: "service and kong service facade not configured in data plane, programmed=false",
			kongUpstreamPolicy: configurationv1beta1.KongUpstreamPolicy{
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
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: now,
					},
				},
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade-1",
						Namespace: testNamespace,
						Annotations: map[string]string{
							configurationv1beta1.KongUpstreamPolicyAnnotationKey: policyName,
						},
						CreationTimestamp: metav1.Time{
							Time: now.Add(10 * time.Second),
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
									builder.NewHTTPBackendRef("svc-1").Build(),
								},
							},
						},
					},
				},
			},
			objectsConfiguredInDataPlane: false,
			expectedKongUpstreamPolicyStatus: gatewayapi.PolicyStatus{
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
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionFalse,
								Reason: string(gatewayapi.GatewayReasonPending),
							},
						},
					},
					{
						AncestorRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group(incubatorv1alpha1.GroupVersion.Group)),
							Kind:      lo.ToPtr(gatewayapi.Kind(incubatorv1alpha1.KongServiceFacadeKind)),
							Namespace: lo.ToPtr(gatewayapi.Namespace(testNamespace)),
							Name:      gatewayapi.ObjectName("svc-facade-1"),
						},
						ControllerName: gatewaycontroller.GetControllerName(),
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.PolicyConditionAccepted),
								Status: metav1.ConditionTrue,
								Reason: string(gatewayapi.PolicyReasonAccepted),
							},
							{
								Type:   string(gatewayapi.GatewayConditionProgrammed),
								Status: metav1.ConditionFalse,
								Reason: string(gatewayapi.GatewayReasonPending),
							},
						},
					},
				},
			},
			updated: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.inputObjects = append(tc.inputObjects, &tc.kongUpstreamPolicy)
			fakeClient := fakectrlruntimeclient.
				NewClientBuilder().
				WithScheme(scheme.Get()).
				WithObjects(tc.inputObjects...).
				WithStatusSubresource(tc.inputObjects...).
				WithIndex(&corev1.Service{}, upstreamPolicyIndexKey, indexServicesOnUpstreamPolicyAnnotation).
				WithIndex(&gatewayapi.HTTPRoute{}, routeBackendRefServiceNameIndexKey, indexRoutesOnBackendRefServiceName).
				WithIndex(&incubatorv1alpha1.KongServiceFacade{}, upstreamPolicyIndexKey, indexServiceFacadesOnUpstreamPolicyAnnotation).
				Build()

			reconciler := KongUpstreamPolicyReconciler{
				Client:                   fakeClient,
				DataplaneClient:          DataPlaneStatusClientMock{ObjectsConfigured: tc.objectsConfiguredInDataPlane},
				KongServiceFacadeEnabled: true,
				HTTPRouteEnabled:         true,
			}

			updated, err := reconciler.enforceKongUpstreamPolicyStatus(t.Context(), &tc.kongUpstreamPolicy)
			assert.NoError(t, err)
			assert.Equal(t, tc.updated, updated)
			newPolicy := &configurationv1beta1.KongUpstreamPolicy{}
			assert.NoError(t, fakeClient.Get(t.Context(), k8stypes.NamespacedName{
				Namespace: tc.kongUpstreamPolicy.Namespace,
				Name:      tc.kongUpstreamPolicy.Name,
			}, newPolicy))
			ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
			assert.Empty(t, cmp.Diff(tc.expectedKongUpstreamPolicyStatus, newPolicy.Status, ignoreLastTransitionTime))
		})
	}
}

type DataPlaneStatusClientMock struct {
	controllers.DataPlane
	ObjectsConfigured bool
}

func (d DataPlaneStatusClientMock) KubernetesObjectIsConfigured(client.Object) bool {
	return d.ObjectsConfigured
}

func TestHttpRouteHasUpstreamPolicyConflictedBackendRefsWithService(t *testing.T) {
	testCases := []struct {
		name                   string
		httpRoute              gatewayapi.HTTPRoute
		upstreamPolicyServices servicesSet
		serviceRef             serviceKey
		expected               bool
	}{
		{
			name:       "service not referenced by the http route",
			serviceRef: "default/svc-not-referenced",
			upstreamPolicyServices: servicesSet{
				"default/svc-not-referenced": {},
			},
			httpRoute: gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpRoute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-1").Build(),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name:       "service referenced by the http route alone in a rule",
			serviceRef: "default/svc-1",
			upstreamPolicyServices: servicesSet{
				"default/svc-1": {},
			},
			httpRoute: gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpRoute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-1").Build(),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name:       "service referenced by the http route alone in a rule while there's another rule with a service not using the same policy",
			serviceRef: "default/svc-1",
			upstreamPolicyServices: servicesSet{
				"default/svc-1": {},
			},
			httpRoute: gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpRoute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-1").Build(),
							},
						},
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-2").Build(),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name:       "service referenced by the http route in a rule together with another service using the same policy",
			serviceRef: "default/svc-1",
			upstreamPolicyServices: servicesSet{
				"default/svc-1": {},
				"default/svc-2": {},
			},
			httpRoute: gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpRoute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-1").Build(),
								builder.NewHTTPBackendRef("svc-2").Build(),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name:       "service referenced by the http route in a rule together with another service not using the same policy",
			serviceRef: "default/svc-1",
			upstreamPolicyServices: servicesSet{
				"default/svc-1": {},
			},
			httpRoute: gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpRoute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("svc-1").Build(),
								builder.NewHTTPBackendRef("svc-2").Build(),
							},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t,
				tc.expected,
				httpRouteHasUpstreamPolicyConflictedBackendRefsWithService(tc.httpRoute, tc.upstreamPolicyServices, tc.serviceRef),
			)
		})
	}
}

func TestBuildPolicyStatus(t *testing.T) {
	acceptedCondition := metav1.Condition{
		Type:   string(gatewayapi.PolicyConditionAccepted),
		Status: metav1.ConditionTrue,
		Reason: string(gatewayapi.PolicyReasonAccepted),
	}
	programmedCondition := metav1.Condition{
		Type:   string(gatewayapi.GatewayConditionProgrammed),
		Status: metav1.ConditionTrue,
		Reason: string(gatewayapi.GatewayReasonProgrammed),
	}

	serviceStatus := func(name string, creationTimestamp time.Time) ancestorStatus {
		return ancestorStatus{
			namespacedName:      k8stypes.NamespacedName{Namespace: "default", Name: name},
			ancestorKind:        upstreamPolicyAncestorKindService,
			acceptedCondition:   acceptedCondition,
			programmedCondition: programmedCondition,
			creationTimestamp:   metav1.NewTime(creationTimestamp),
		}
	}
	serviceFacadeStatus := func(name string, creationTimestamp time.Time) ancestorStatus {
		return ancestorStatus{
			namespacedName:      k8stypes.NamespacedName{Namespace: "default", Name: name},
			ancestorKind:        upstreamPolicyAncestorKindKongServiceFacade,
			acceptedCondition:   acceptedCondition,
			programmedCondition: programmedCondition,
			creationTimestamp:   metav1.NewTime(creationTimestamp),
		}
	}

	serviceExpectedPolicyAncestorStatus := func(name string) gatewayapi.PolicyAncestorStatus {
		return gatewayapi.PolicyAncestorStatus{
			AncestorRef: gatewayapi.ParentReference{
				Group:     lo.ToPtr(gatewayapi.Group("core")),
				Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
				Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
				Name:      gatewayapi.ObjectName(name),
			},
			ControllerName: gatewaycontroller.GetControllerName(),
			Conditions: []metav1.Condition{
				acceptedCondition,
				programmedCondition,
			},
		}
	}
	serviceFacadeExpectedPolicyAncestorStatus := func(name string) gatewayapi.PolicyAncestorStatus {
		return gatewayapi.PolicyAncestorStatus{
			AncestorRef: gatewayapi.ParentReference{
				Group:     lo.ToPtr(gatewayapi.Group(incubatorv1alpha1.GroupVersion.Group)),
				Kind:      lo.ToPtr(gatewayapi.Kind(incubatorv1alpha1.KongServiceFacadeKind)),
				Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
				Name:      gatewayapi.ObjectName(name),
			},
			ControllerName: gatewaycontroller.GetControllerName(),
			Conditions: []metav1.Condition{
				acceptedCondition,
				programmedCondition,
			},
		}
	}

	now := time.Now()
	testCases := []struct {
		name              string
		ancestorsStatuses []ancestorStatus
		expected          gatewayapi.PolicyStatus
	}{
		{
			name: "all ordered by creation timestamp (oldest first)",
			ancestorsStatuses: []ancestorStatus{
				serviceStatus("svc-1", now.Add(4*time.Second)),
				serviceStatus("svc-2", now.Add(3*time.Second)),
				serviceFacadeStatus("svc-facade-1", now.Add(2*time.Second)),
				serviceFacadeStatus("svc-facade-2", now.Add(1*time.Second)),
			},
			expected: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					serviceFacadeExpectedPolicyAncestorStatus("svc-facade-2"),
					serviceFacadeExpectedPolicyAncestorStatus("svc-facade-1"),
					serviceExpectedPolicyAncestorStatus("svc-2"),
					serviceExpectedPolicyAncestorStatus("svc-1"),
				},
			},
		},
		{
			name: "ordered by kind, namespace and name if their creationTimestamp is the same",
			ancestorsStatuses: []ancestorStatus{
				serviceStatus("svc-1", now),
				serviceStatus("svc-2", now),
				serviceFacadeStatus("svc-facade-1", now),
				serviceFacadeStatus("svc-facade-2", now),
			},
			expected: gatewayapi.PolicyStatus{
				Ancestors: []gatewayapi.PolicyAncestorStatus{
					serviceFacadeExpectedPolicyAncestorStatus("svc-facade-1"),
					serviceFacadeExpectedPolicyAncestorStatus("svc-facade-2"),
					serviceExpectedPolicyAncestorStatus("svc-1"),
					serviceExpectedPolicyAncestorStatus("svc-2"),
				},
			},
		},
		{
			name: "more ancestors than allowed - keeps only maxNAncestors oldest ones",
			ancestorsStatuses: func() []ancestorStatus {
				var ancestors []ancestorStatus
				for i := range maxNAncestors + 2 {
					ancestors = append(ancestors, serviceStatus(fmt.Sprintf("svc-%d", i), now.Add(time.Duration(i)*time.Second)))
				}
				return ancestors
			}(),
			expected: gatewayapi.PolicyStatus{
				Ancestors: func() []gatewayapi.PolicyAncestorStatus {
					var ancestors []gatewayapi.PolicyAncestorStatus
					for i := range maxNAncestors {
						ancestors = append(ancestors, serviceExpectedPolicyAncestorStatus(fmt.Sprintf("svc-%d", i)))
					}
					return ancestors
				}(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &KongUpstreamPolicyReconciler{Log: logr.Discard()}
			policyStatus, err := r.buildPolicyStatus(k8stypes.NamespacedName{Namespace: "default", Name: "test-policy"}, tc.ancestorsStatuses)
			require.NoError(t, err)
			require.Equal(t, tc.expected, policyStatus)
		})
	}
}

func TestIsSupportedHTTPRouteBackendRef(t *testing.T) {
	testCases := []struct {
		name       string
		backendRef gatewayapi.BackendObjectReference
		expected   bool
	}{
		{
			name: "core service backend ref",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("core")),
				Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
			},
			expected: true,
		},
		{
			name: "core service with nil group",
			backendRef: gatewayapi.BackendObjectReference{
				Kind: lo.ToPtr(gatewayapi.Kind("Service")),
			},
			expected: true,
		},
		{
			name: "core service with nil kind",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("core")),
			},
			expected: true,
		},
		{
			name:       "core service with nil group and kind",
			backendRef: gatewayapi.BackendObjectReference{},
			expected:   true,
		},
		{
			name: "core group unsupported kind",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("core")),
				Kind:  lo.ToPtr(gatewayapi.Kind("UnsupportedKind")),
			},
			expected: false,
		},
		{
			name: "service unsupported group",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("unsupported")),
				Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
			},
			expected: false,
		},
		{
			name: "unsupported group",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("unsupported")),
			},
			expected: false,
		},
		{
			name: "unsupported kind",
			backendRef: gatewayapi.BackendObjectReference{
				Kind: lo.ToPtr(gatewayapi.Kind("UnsupportedKind")),
			},
			expected: false,
		},
		{
			name: "empty group",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("")),
				Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
			},
			expected: true,
		},
		{
			name: "empty group with nil kind",
			backendRef: gatewayapi.BackendObjectReference{
				Group: lo.ToPtr(gatewayapi.Group("")),
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, isSupportedHTTPRouteBackendRef(gatewayapi.BackendRef{BackendObjectReference: tc.backendRef}))
		})
	}
}
