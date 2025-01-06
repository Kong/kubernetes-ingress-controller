package translator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestGetIngressClassParameters(t *testing.T) {
	var (
		scopeNamespace    = "Namespace"
		testNamespaceName = "test-ns"
		testIcpName       = "test-icp"
	)

	defaultIcpSpec := &kongv1alpha1.IngressClassParametersSpec{}
	icp := &kongv1alpha1.IngressClassParameters{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNamespaceName,
			Name:      testIcpName,
		},
		Spec: kongv1alpha1.IngressClassParametersSpec{
			EnableLegacyRegexDetection: true,
		},
	}

	testCases := []struct {
		name          string
		paramRef      *netv1.IngressClassParametersReference
		parameterSpec *kongv1alpha1.IngressClassParametersSpec
		err           error
	}{
		{
			name:          "nil-paramref",
			parameterSpec: defaultIcpSpec,
			// No error: it's OK to not reference ingress class parameters.
		},
		{
			name: "nil-apigroup",
			paramRef: &netv1.IngressClassParametersReference{
				Kind: "configMap",
				Name: "some-cm",
			},
			parameterSpec: defaultIcpSpec,
			err:           fmt.Errorf("IngressClass nil-apigroup should reference parameters in apiGroup:configuration.konghq.com"),
		},
		{
			name: "nil-scope",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &kongv1alpha1.GroupVersion.Group,
				Kind:      kongv1alpha1.IngressClassParametersKind,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			err:           fmt.Errorf("IngressClass nil-scope should reference namespaced parameters"),
		},
		{
			name: "nil-namespace",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup: &kongv1alpha1.GroupVersion.Group,
				Kind:     kongv1alpha1.IngressClassParametersKind,
				Scope:    &scopeNamespace,
				Name:     testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			err:           fmt.Errorf("IngressClass nil-namespace should reference namespaced parameters"),
		},
		{
			name: "matched-parameters",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &kongv1alpha1.GroupVersion.Group,
				Kind:      kongv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: &icp.Spec,
		},
		{
			name: "unmatched-kind",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &kongv1alpha1.GroupVersion.Group,
				Kind:      "SomeKind",
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			err:           fmt.Errorf("IngressClass unmatched-kind should reference parameters with kind:IngressClassParameters"),
		},
		{
			name: "unmatched-namespace",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &kongv1alpha1.GroupVersion.Group,
				Kind:      kongv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: new(string),
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			err:           store.NotFoundError{Message: "IngressClassParameters test-icp not found"},
		},
		{
			name: "unmatched-name",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &kongv1alpha1.GroupVersion.Group,
				Kind:      kongv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      "another-icp",
			},
			parameterSpec: defaultIcpSpec,
			err:           store.NotFoundError{Message: "IngressClassParameters another-icp not found"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ingressClass := &netv1.IngressClass{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespaceName,
					Name:      tc.name,
				},
				Spec: netv1.IngressClassSpec{
					Parameters: tc.paramRef,
				},
			}
			cacheStores, err := store.NewCacheStoresFromObjs()
			require.NoError(t, err)
			err = cacheStores.Add(ingressClass)
			require.NoError(t, err)
			err = cacheStores.Add(icp)
			require.NoError(t, err)
			s := store.New(cacheStores, ingressClass.Name, zapr.NewLogger(zap.NewNop()))
			icpSpec, err := getIngressClassParametersOrDefault(s)
			assert.Truef(t, reflect.DeepEqual(*tc.parameterSpec, icpSpec),
				fmt.Sprintf("should get same ingress parameter spec: expected %+v, actual %+v", tc.parameterSpec, icpSpec),
			)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())

				if errors.As(tc.err, &store.NotFoundError{}) {
					assert.ErrorAs(t, err, &store.NotFoundError{})
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
