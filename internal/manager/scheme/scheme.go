package scheme

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
)

// Get returns a scheme aware of all types the manager can interact with.
func Get() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()

	if err := apiextensionsv1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := kongv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := kongv1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := kongv1beta1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := incubatorv1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := gatewayv1alpha2.Install(scheme); err != nil {
		return nil, err
	}

	if err := gatewayv1beta1.Install(scheme); err != nil {
		return nil, err
	}

	if err := gatewayv1.Install(scheme); err != nil {
		return nil, err
	}

	return scheme, nil
}
