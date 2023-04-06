package parser

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrlclientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/fake"
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
				APIGroup:  &configurationv1alpha1.GroupVersion.Group,
				Kind:      configurationv1alpha1.IngressClassParametersKind,
				Namespace: &testNamespaceName,
				Name:      testIcpName,
			},
			parameterSpec: defaultIcpSpec,
			err:           fmt.Errorf("IngressClass nil-scope should reference namespaced parameters"),
		},

		// TODO(pmalek): these need to be fixed because now getting an ingressclass
		// is namespace sensitive, i.e. store.GetIngressClassName() should take into
		// account that IngressClass can also have namespace specified.
		// ref: https://kubernetes.io/docs/concepts/services-networking/ingress/#ingressclass-scope
		//
		// {
		// 	name: "matched-parameters",
		// 	paramRef: &netv1.IngressClassParametersReference{
		// 		APIGroup:  &configurationv1alpha1.GroupVersion.Group,
		// 		Kind:      configurationv1alpha1.IngressClassParametersKind,
		// 		Scope:     &scopeNamespace,
		// 		Namespace: &testNamespaceName,
		// 		Name:      testIcpName,
		// 	},
		// 	parameterSpec: &icp.Spec,
		// },
		// {
		// 	name: "unmatched-kind",
		// 	paramRef: &netv1.IngressClassParametersReference{
		// 		APIGroup:  &configurationv1alpha1.GroupVersion.Group,
		// 		Kind:      "SomeKind",
		// 		Scope:     &scopeNamespace,
		// 		Namespace: &testNamespaceName,
		// 		Name:      testIcpName,
		// 	},
		// 	parameterSpec: defaultIcpSpec,
		// 	err:           fmt.Errorf("IngressClass unmatched-kind should reference parameters with kind:IngressClassParameters"),
		// },
		// {
		// 	name: "unmatched-namespace",
		// 	paramRef: &netv1.IngressClassParametersReference{
		// 		APIGroup:  &configurationv1alpha1.GroupVersion.Group,
		// 		Kind:      configurationv1alpha1.IngressClassParametersKind,
		// 		Scope:     &scopeNamespace,
		// 		Namespace: new(string),
		// 		Name:      testIcpName,
		// 	},
		// 	parameterSpec: defaultIcpSpec,
		// 	err:           store.ErrNotFound{Message: "IngressClassParameters test-icp not found"},
		// },
		// {
		// 	name: "unmatched-name",
		// 	paramRef: &netv1.IngressClassParametersReference{
		// 		APIGroup:  &configurationv1alpha1.GroupVersion.Group,
		// 		Kind:      configurationv1alpha1.IngressClassParametersKind,
		// 		Scope:     &scopeNamespace,
		// 		Namespace: &testNamespaceName,
		// 		Name:      "another-icp",
		// 	},
		// 	parameterSpec: defaultIcpSpec,
		// 	err:           store.ErrNotFound{Message: "IngressClassParameters another-icp not found"},
		// },
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ingressClass := &netv1.IngressClass{
				TypeMeta: metav1.TypeMeta{
					Kind:       "IngressClass",
					APIVersion: "networking.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: tc.name,
				},
				Spec: netv1.IngressClassSpec{
					Parameters: tc.paramRef,
				},
			}
			if tc.paramRef != nil && tc.paramRef.Scope != nil && *tc.paramRef.Scope == scopeNamespace {
				ingressClass.Namespace = testNamespaceName
			}

			require.NoError(t, fake.AddToScheme(scheme.Scheme))
			client := ctrlclientfake.NewClientBuilder().
				WithObjects(
					ingressClass,
					icp,
				).
				Build()
			s := store.New(client, ingressClass.Name, logrus.New())
			icpSpec, err := getIngressClassParametersOrDefault(context.TODO(), s)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())

				if errors.As(tc.err, &store.ErrNotFound{}) {
					assert.ErrorAs(t, err, &store.ErrNotFound{})
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, *tc.parameterSpec, icpSpec)
			}
		})
	}
}
