package utils_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	commonv1alpha1 "github.com/kong/kubernetes-configuration/api/common/v1alpha1"

	ctrlutils "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
)

type objectWithCPRefType struct {
	client.Object
	cpRef *commonv1alpha1.ControlPlaneRef
}

func (o *objectWithCPRefType) GetControlPlaneRef() *commonv1alpha1.ControlPlaneRef {
	return o.cpRef
}

func TestGenerateCPReferenceMatchesPredicate(t *testing.T) {
	testCases := []struct {
		name     string
		obj      objectWithCPRefType
		expected bool
	}{
		{
			name: "control plane reference is nil",
			obj: objectWithCPRefType{
				cpRef: nil,
			},
			expected: true,
		},
		{
			name: "control plane reference is set to kic",
			obj: objectWithCPRefType{
				cpRef: &commonv1alpha1.ControlPlaneRef{
					Type: commonv1alpha1.ControlPlaneRefKIC,
				},
			},
			expected: true,
		},
		{
			name: "control plane reference is set to konnect",
			obj: objectWithCPRefType{
				cpRef: &commonv1alpha1.ControlPlaneRef{
					Type:      commonv1alpha1.ControlPlaneRefKonnectID,
					KonnectID: lo.ToPtr("konnect-id"),
				},
			},
			expected: false,
		},
		{
			name: "control plane reference is set to konnect namespaced reference",
			obj: objectWithCPRefType{
				cpRef: &commonv1alpha1.ControlPlaneRef{
					Type: commonv1alpha1.ControlPlaneRefKonnectNamespacedRef,
					KonnectNamespacedRef: &commonv1alpha1.KonnectNamespacedRef{
						Name: "konnect-name",
					},
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pred := ctrlutils.GenerateCPReferenceMatchesPredicate[*objectWithCPRefType]()
			actual := pred.Generic(event.GenericEvent{
				Object: &tc.obj,
			})
			require.Equal(t, tc.expected, actual)
		})
	}
}
