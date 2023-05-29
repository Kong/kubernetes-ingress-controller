package manager_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestIngressControllerConditions(t *testing.T) {
	const kind = "Ingress"
	networkingV1 := schema.GroupVersionResource{
		Group:    netv1.SchemeGroupVersion.Group,
		Version:  netv1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}

	type ingressTestOpts struct {
		enabled      bool
		crdInstalled bool
	}

	testCases := []struct {
		name string

		ingressNetV1      ingressTestOpts
		ingressClassNetV1 ingressTestOpts

		expectIngressNetV1      bool
		expectIngressClassNetV1 bool
		expectError             bool
	}{
		{
			name:                    "netV1_takes_precedence_over_all",
			ingressNetV1:            ingressTestOpts{enabled: true, crdInstalled: true},
			ingressClassNetV1:       ingressTestOpts{enabled: true, crdInstalled: true},
			expectIngressNetV1:      true,
			expectIngressClassNetV1: true,
		},
		{
			name:              "no_crds_installed",
			ingressNetV1:      ingressTestOpts{enabled: true},
			ingressClassNetV1: ingressTestOpts{enabled: true},
			expectError:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			restMapper := meta.NewDefaultRESTMapper(nil)
			if tc.ingressNetV1.crdInstalled || tc.ingressClassNetV1.crdInstalled {
				restMapper.Add(schema.GroupVersionKind{
					Group:   networkingV1.Group,
					Version: networkingV1.Version,
					Kind:    kind,
				}, meta.RESTScopeRoot)
			}

			conditions, err := manager.NewIngressControllersConditions(
				&manager.Config{
					IngressNetV1Enabled:      tc.ingressNetV1.enabled,
					IngressClassNetV1Enabled: tc.ingressClassNetV1.enabled,
				},
				restMapper,
			)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectIngressNetV1, conditions.IngressNetV1Enabled())
			assert.Equal(t, tc.expectIngressClassNetV1, conditions.IngressClassNetV1Enabled())
		})
	}
}

func TestShouldEnableCRDController(t *testing.T) {
	knownGvr := schema.GroupVersionResource{
		Group:    "group",
		Version:  "v1",
		Resource: "resources",
	}
	unknownGVR := schema.GroupVersionResource{
		Group:    "otherGroup",
		Version:  "v1",
		Resource: "resources",
	}

	restMapper := meta.NewDefaultRESTMapper(nil)
	restMapper.Add(schema.GroupVersionKind{
		Group:   knownGvr.Group,
		Version: knownGvr.Version,
		Kind:    "Resource",
	}, meta.RESTScopeRoot)

	testCases := []struct {
		name           string
		gvr            schema.GroupVersionResource
		expectedResult bool
	}{
		{
			name:           "registered_resource",
			gvr:            knownGvr,
			expectedResult: true,
		},
		{
			name:           "not_registered_resource",
			gvr:            unknownGVR,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				manager.ShouldEnableCRDController(tc.gvr, restMapper),
			)
		})
	}
}
