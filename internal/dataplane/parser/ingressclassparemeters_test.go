package parser

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
)

func TestGetIngressClassParameters(t *testing.T) {
	var (
		scopeNamespace    = "Namespace"
		testNamespaceName = "test-ns"
		testIcpName       = "test-icp"
	)

	defaultIcpSpec := &configurationv1alpha1.IngressClassParametersSpec{}
	icp := &configurationv1alpha1.IngressClassParameters{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNamespaceName,
			Name:      testIcpName,
		},
		Spec: configurationv1alpha1.IngressClassParametersSpec{
			EnableLegacyRegexDetection: true,
		},
	}

	testCases := []struct {
		name          string
		paramRef      *netv1.IngressClassParametersReference
		parameterSpec *configurationv1alpha1.IngressClassParametersSpec
		hasError      bool
	}{
		{
			name:          "nil-paramref",
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "nil-apigroup",
			paramRef: &netv1.IngressClassParametersReference{
				Kind: "configMap",
				Name: "some-cm",
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "nil-scope",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      configurationv1alpha1.IngressClassParametersKind,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "nil-namespace",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup: &configurationv1alpha1.GroupVersion.Group,
				Kind:     configurationv1alpha1.IngressClassParametersKind,
				Scope:    &scopeNamespace,
				Name:     testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "matched-parameters",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      configurationv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: &icp.Spec,
		},
		{
			name: "unmatched-kind",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      "SomeKind",
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "unmatched-namespace",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      configurationv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: new(string),
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
		{
			name: "unmatched-name",
			paramRef: &netv1.IngressClassParametersReference{
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      configurationv1alpha1.IngressClassParametersKind,
				Scope:     &scopeNamespace,
				Namespace: &testNamespaceName,
				Name:      "another-icp",
			},
			parameterSpec: defaultIcpSpec,
			hasError:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc
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
			assert.NoError(t, err)
			err = cacheStores.Add(ingressClass)
			assert.NoError(t, err)
			err = cacheStores.Add(icp)
			assert.NoError(t, err)
			s := store.New(cacheStores, ingressClass.Name, true, true, true, logrus.New())
			icpSpec, err := getIngressClassParametersOrDefault(s)
			assert.Truef(t, reflect.DeepEqual(*tc.parameterSpec, icpSpec),
				fmt.Sprintf("should get same ingress parameter spec: expected %+v, actual %+v", tc.parameterSpec, icpSpec),
			)
			if tc.hasError {
				assert.NotNil(t, err)
				assert.Truef(t, errors.As(err, &store.ErrNotFound{}), "error should be a store.ErrNotFound, actual %T", err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
