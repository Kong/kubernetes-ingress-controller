package scheme

import (
	"sync/atomic"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1alpha3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
)

var (
	kicScheme       = runtime.NewScheme()
	addToSchemeDone atomic.Bool
)

// REVIEW: Is it better to move initialization of scheme elsewhere,
// and change the `Get` to only return the scheme (or directly use the scheme)?

// Get returns a scheme aware of all types the manager can interact with.
func Get() (*runtime.Scheme, error) {
	// Only add schemes to the runtime scheme when adding is not done.
	if !addToSchemeDone.Load() {
		if err := apiextensionsv1.AddToScheme(kicScheme); err != nil {
			return nil, err
		}

		if err := clientgoscheme.AddToScheme(kicScheme); err != nil {
			return nil, err
		}

		if err := kongv1.AddToScheme(kicScheme); err != nil {
			return nil, err
		}
		if err := kongv1alpha1.AddToScheme(kicScheme); err != nil {
			return nil, err
		}
		if err := kongv1beta1.AddToScheme(kicScheme); err != nil {
			return nil, err
		}
		if err := incubatorv1alpha1.AddToScheme(kicScheme); err != nil {
			return nil, err
		}

		if err := gatewayv1alpha2.Install(kicScheme); err != nil {
			return nil, err
		}
		if err := gatewayv1alpha3.Install(kicScheme); err != nil {
			return nil, err
		}
		if err := gatewayv1beta1.Install(kicScheme); err != nil {
			return nil, err
		}
		if err := gatewayv1.Install(kicScheme); err != nil {
			return nil, err
		}
	}
	addToSchemeDone.Store(true)
	// REVIEW: is it safe to return the same scheme everywhere it is used?
	return kicScheme, nil
}
