package gateway

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestGetBackendTLSPolicyAncestors(t *testing.T) {
	tests := []struct {
		name        string
		policy      gatewayapi.BackendTLSPolicy
		objects     []client.Object
		expected    []k8stypes.NamespacedName
		expectError bool
	}{
		{
			name: "target ref not a service",
			policy: gatewayapi.BackendTLSPolicy{
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: gatewayapi.Group("other-group"),
								Kind:  gatewayapi.Kind("other-kind"),
								Name:  "example-service",
							},
						},
					},
				},
			},
		},
		{
			name: "valid target ref, httproute with resolved refs",
			policy: gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
				},
			},
			objects: []client.Object{
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Group: lo.ToPtr(gatewayapi.Group("core")),
												Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
												Name:  "example-service",
											},
										},
									},
								},
							},
						},
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
									Name:  "example-gateway",
								},
							},
						},
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ControllerName: GetControllerName(),
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
										Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:      "example-gateway",
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
									},
									Conditions: []metav1.Condition{
										{
											Type:   string(gatewayapi.RouteConditionResolvedRefs),
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "example-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: "example-gateway-class",
					},
				},
			},
			expected: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "example-gateway",
				},
			},
		},
		{
			name: "valid target ref, httproute without resolved refs",
			policy: gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
				},
			},
			objects: []client.Object{
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Group: lo.ToPtr(gatewayapi.Group("core")),
												Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
												Name:  "example-service",
											},
										},
									},
								},
							},
						},
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
									Name:  "example-gateway",
								},
							},
						},
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ControllerName: GetControllerName(),
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
										Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:      "example-gateway",
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
									},
									Conditions: []metav1.Condition{
										{
											Type:   string(gatewayapi.RouteConditionResolvedRefs),
											Status: metav1.ConditionFalse,
										},
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "example-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: "example-gateway-class",
					},
				},
			},
		},
		{
			name: "multiple gateways belonging to multiple controllers",
			policy: gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
				},
			},
			objects: []client.Object{
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Group: lo.ToPtr(gatewayapi.Group("core")),
												Kind:  lo.ToPtr(gatewayapi.Kind("Service")),
												Name:  "example-service",
											},
										},
									},
								},
							},
						},
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
									Name:  "other-gateway",
								},
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
									Name:  "example-gateway",
								},
							},
						},
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ControllerName: GetControllerName(),
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
										Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:      "example-gateway",
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
									},
									Conditions: []metav1.Condition{
										{
											Type:   string(gatewayapi.RouteConditionResolvedRefs),
											Status: metav1.ConditionTrue,
										},
									},
								},
								{
									ControllerName: "other-controller",
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
										Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:      "other-gateway",
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
									},
									Conditions: []metav1.Condition{
										{
											Type:   string(gatewayapi.RouteConditionResolvedRefs),
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "example-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: "example-gateway-class",
					},
				},
			},
			expected: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "example-gateway",
				},
			},
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, gatewayapi.InstallV1(scheme))
	require.NoError(t, gatewayapi.InstallV1alpha3(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.objects...).
				WithIndex(
					&gatewayapi.HTTPRoute{},
					httpRouteBackendRefIndexKey,
					indexHTTPRouteOnBackendRef,
				).
				Build()

			r := &BackendTLSPolicyReconciler{
				Client: cl,
			}

			gateways, err := r.getBackendTLSPolicyAncestors(context.Background(), tt.policy)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for _, gateway := range gateways {
					key := client.ObjectKeyFromObject(&gateway)
					assert.Contains(t, tt.expected, key)
				}
			}
		})
	}
}

func TestSortGateways(t *testing.T) {
	tests := []struct {
		name              string
		gateways          []gatewayapi.Gateway
		expected          []gatewayapi.Gateway
		existingAncestors []gatewayapi.PolicyAncestorStatus
		policyNamespace   string
	}{
		{
			name: "different namespaces, no existing ancestors",
			gateways: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
			},
			expected: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
			},
		},
		{
			name: "same namespace, no existing ancestors",
			gateways: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
			},
			expected: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
			},
		},
		{
			name: "multiple combinations, no existing ancestors",
			gateways: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
			},
			expected: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
			},
		},
		{
			name: "multiple combinations, with existing ancestors",
			gateways: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-2",
					},
				},
			},
			expected: []gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-2",
						Name:      "gateway-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "namespace-1",
						Name:      "gateway-2",
					},
				},
			},
			existingAncestors: []gatewayapi.PolicyAncestorStatus{
				{
					AncestorRef: gatewayapi.ParentReference{
						Namespace: lo.ToPtr(gatewayapi.Namespace("namespace-2")),
						Name:      "gateway-2",
					},
				},
				{
					AncestorRef: gatewayapi.ParentReference{
						Name: "gateway-1",
					},
				},
			},
			policyNamespace: "namespace-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortGateways(tt.gateways, tt.existingAncestors, tt.policyNamespace)
			assert.Equal(t, tt.expected, tt.gateways)
		})
	}
}

func TestValidateBackendTLSPolicy(t *testing.T) {
	tests := []struct {
		name        string
		policy      *gatewayapi.BackendTLSPolicy
		objects     []client.Object
		expected    *metav1.Condition
		expectError bool
	}{
		{
			name: "valid policy with no conflicts",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:   string(gatewayapi.PolicyConditionAccepted),
				Status: metav1.ConditionTrue,
				Reason: string(gatewayapi.PolicyConditionAccepted),
			},
		},
		{
			name: "policy with conflicting target refs",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
				},
			},
			objects: []client.Object{
				&gatewayapi.BackendTLSPolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "conflicting-policy",
						Namespace: "default",
					},
					Spec: gatewayapi.BackendTLSPolicySpec{
						TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
							{
								LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
									Group: "core",
									Kind:  "Service",
									Name:  "example-service",
								},
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonConflicted),
				Message: "Multiple BackendTLSPolicies target the same service",
			},
		},
		{
			name: "policy with unsupported CACertificateRefs",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						CACertificateRefs: []gatewayapi.LocalObjectReference{
							{
								Group: "core",
								Kind:  "Secret",
								Name:  "example-secret",
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "CACertificateRefs must reference ConfigMaps in the core group",
			},
		},
		{
			name: "policy with unsupported SubjectAltNames",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						SubjectAltNames: []gatewayapi.SubjectAltName{
							{
								Type:     gatewayapi.HostnameSubjectAltNameType,
								Hostname: "example.com",
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "SubjectAltNames feature is not currently supported",
			},
		},
		{
			name: "policy with unsupported WellKnownCACertificates",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						WellKnownCACertificates: lo.ToPtr(gatewayapi.WellKnownCACertificatesType("well-known-ca")),
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "WellKnownCACertificates feature is not currently supported",
			},
		},
		{
			name: "policy with both unsupported features",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						SubjectAltNames: []gatewayapi.SubjectAltName{
							{
								Type:     gatewayapi.HostnameSubjectAltNameType,
								Hostname: "example.com",
							},
						},
						WellKnownCACertificates: lo.ToPtr(gatewayapi.WellKnownCACertificatesType("well-known-ca")),
						CACertificateRefs: []gatewayapi.LocalObjectReference{
							{
								Group: "core",
								Kind:  "Secret",
								Name:  "example-secret",
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "CACertificateRefs must reference ConfigMaps in the core group - SubjectAltNames feature is not currently supported - WellKnownCACertificates feature is not currently supported",
			},
		},
		{
			name: "valid policy referencing not existing CACert (ConfigMap)",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						CACertificateRefs: []gatewayapi.LocalObjectReference{
							{
								Group: "",
								Kind:  "ConfigMap",
								Name:  gatewayapi.ObjectName("example-configmap"),
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "failed getting ConfigMap default/example-configmap set as CACertificateRef: configmaps \"example-configmap\" not found",
			},
		},
		{
			name: "valid policy referencing not existing CACert (ConfigMap, group core)",
			policy: &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-policy",
					Namespace: "default",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  "example-service",
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						CACertificateRefs: []gatewayapi.LocalObjectReference{
							{
								Group: "core",
								Kind:  "ConfigMap",
								Name:  gatewayapi.ObjectName("example-configmap"),
							},
						},
					},
				},
			},
			expected: &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonInvalid),
				Message: "failed getting ConfigMap default/example-configmap set as CACertificateRef: configmaps \"example-configmap\" not found",
			},
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, clientgoscheme.AddToScheme(scheme))
	require.NoError(t, gatewayapi.InstallV1(scheme))
	require.NoError(t, gatewayapi.InstallV1alpha3(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects := []client.Object{
				tt.policy,
			}
			objects = append(objects, tt.objects...)
			cl := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				WithIndex(
					&gatewayapi.BackendTLSPolicy{},
					backendTLSPolicyTargetRefIndexKey,
					indexBackendTLSPolicyOnTargetRef,
				).
				WithIndex(
					&gatewayapi.BackendTLSPolicy{},
					backendTLSPolicyValidationCARefIndexKey,
					indexBackendTLSPolicyOnValidationCACertificateRef,
				).
				Build()

			r := &BackendTLSPolicyReconciler{
				Client: cl,
			}

			condition, err := r.validateBackendTLSPolicy(context.Background(), *tt.policy)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, condition)
			}
		})
	}
}
