package envtest

import (
	"testing"

	"github.com/stretchr/testify/require"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
)

type SchemeOption func(t *testing.T, s *k8sruntime.Scheme)

// WithGatewayAPI registers the Gateway API types with the scheme.
func WithGatewayAPI(t *testing.T, s *k8sruntime.Scheme) {
	require.NoError(t, gatewayv1.Install(s))
	require.NoError(t, gatewayv1beta1.Install(s))
	require.NoError(t, gatewayv1alpha2.Install(s))
}

// WithKong registers the Kong types with the scheme.
func WithKong(t *testing.T, s *k8sruntime.Scheme) {
	require.NoError(t, kongv1.AddToScheme(s))
	require.NoError(t, kongv1beta1.AddToScheme(s))
	require.NoError(t, kongv1alpha1.AddToScheme(s))
	require.NoError(t, incubatorv1alpha1.AddToScheme(s))
}

// Scheme returns a new scheme with the default Kubernetes types registered.
// It accepts optional SchemeOptions to register additional types.
func Scheme(t *testing.T, opts ...SchemeOption) *k8sruntime.Scheme {
	s := k8sruntime.NewScheme()
	require.NoError(t, k8sscheme.AddToScheme(s))

	for _, opt := range opts {
		opt(t, s)
	}

	return s
}
