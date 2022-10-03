package manager_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

type ingressTestOpts struct {
	enabled      bool
	crdInstalled bool
}

func TestIngressControllerConditions(t *testing.T) {
	const kind = "Ingress"
	var (
		networkingV1 = schema.GroupVersionResource{
			Group:    netv1.SchemeGroupVersion.Group,
			Version:  netv1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}
		networkingV1beta1 = schema.GroupVersionResource{
			Group:    netv1beta1.SchemeGroupVersion.Group,
			Version:  netv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}
		extensionsV1beta1 = schema.GroupVersionResource{
			Group:    extensionsv1beta1.SchemeGroupVersion.Group,
			Version:  extensionsv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}
	)

	testCases := []struct {
		name string

		ingressNetV1      ingressTestOpts
		ingressClassNetV1 ingressTestOpts
		ingressNetV1beta  ingressTestOpts
		ingressExtV1beta  ingressTestOpts

		expectIngressNetV1      bool
		expectIngressClassNetV1 bool
		expectIngressNetV1beta  bool
		expectIngressExtV1beta  bool
		expectError             bool
	}{
		{
			name:                    "netV1_takes_precedence_over_all",
			ingressNetV1:            ingressTestOpts{true, true},
			ingressClassNetV1:       ingressTestOpts{true, true},
			ingressNetV1beta:        ingressTestOpts{true, true},
			ingressExtV1beta:        ingressTestOpts{true, true},
			expectIngressNetV1:      true,
			expectIngressClassNetV1: true,
		},
		{
			name:                   "netV1beta_wins_when_netV1_crds_not_installed",
			ingressNetV1:           ingressTestOpts{true, false},
			ingressClassNetV1:      ingressTestOpts{true, false},
			ingressNetV1beta:       ingressTestOpts{true, true},
			ingressExtV1beta:       ingressTestOpts{true, true},
			expectIngressNetV1beta: true,
		},
		{
			name:                   "extV1beta_wins_when_netV1_netV1beta_crds_not_installed",
			ingressNetV1:           ingressTestOpts{true, false},
			ingressClassNetV1:      ingressTestOpts{true, false},
			ingressNetV1beta:       ingressTestOpts{true, false},
			ingressExtV1beta:       ingressTestOpts{true, true},
			expectIngressExtV1beta: true,
		},
		{
			name:                   "netV1_not_picked_when_disabled",
			ingressNetV1:           ingressTestOpts{false, true},
			ingressClassNetV1:      ingressTestOpts{true, true},
			ingressNetV1beta:       ingressTestOpts{true, true},
			ingressExtV1beta:       ingressTestOpts{true, true},
			expectIngressNetV1beta: true,
		},
		{
			name:              "no_crds_installed",
			ingressNetV1:      ingressTestOpts{true, false},
			ingressClassNetV1: ingressTestOpts{true, false},
			ingressNetV1beta:  ingressTestOpts{true, false},
			ingressExtV1beta:  ingressTestOpts{true, false},
			expectError:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			restMapper := meta.NewDefaultRESTMapper(nil)
			if tc.ingressNetV1.crdInstalled || tc.ingressClassNetV1.crdInstalled {
				restMapper.Add(schema.GroupVersionKind{
					Group:   networkingV1.Group,
					Version: networkingV1.Version,
					Kind:    kind,
				}, meta.RESTScopeRoot)
			}
			if tc.ingressNetV1beta.crdInstalled {
				restMapper.Add(schema.GroupVersionKind{
					Group:   networkingV1beta1.Group,
					Version: networkingV1beta1.Version,
					Kind:    kind,
				}, meta.RESTScopeRoot)
			}
			if tc.ingressExtV1beta.crdInstalled {
				restMapper.Add(schema.GroupVersionKind{
					Group:   extensionsV1beta1.Group,
					Version: extensionsV1beta1.Version,
					Kind:    kind,
				}, meta.RESTScopeRoot)
			}

			conditions, err := manager.NewIngressControllersConditions(
				&manager.Config{
					IngressNetV1Enabled:      tc.ingressNetV1.enabled,
					IngressClassNetV1Enabled: tc.ingressClassNetV1.enabled,
					IngressNetV1beta1Enabled: tc.ingressNetV1beta.enabled,
					IngressExtV1beta1Enabled: tc.ingressExtV1beta.enabled,
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
			assert.Equal(t, tc.expectIngressNetV1beta, conditions.IngressNetV1beta1Enabled())
			assert.Equal(t, tc.expectIngressExtV1beta, conditions.IngressExtV1beta1Enabled())
		})
	}
}
