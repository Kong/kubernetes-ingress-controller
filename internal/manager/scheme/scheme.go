package scheme

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// Get returns the scheme for the manager, enabling all the default schemes and
// those that were enabled via the feature flags.
func Get(fg map[string]bool) (*runtime.Scheme, error) {
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

	if v, ok := fg[featuregates.KnativeFeature]; ok && v {
		if err := knativev1alpha1.AddToScheme(scheme); err != nil {
			return nil, err
		}
	}

	if v, ok := fg[featuregates.GatewayAlphaFeature]; ok && v {
		if err := gatewayv1alpha2.Install(scheme); err != nil {
			return nil, err
		}
	}
	if v, ok := fg[featuregates.GatewayFeature]; ok && v {
		if err := gatewayv1.Install(scheme); err != nil {
			return nil, err
		}
	}

	return scheme, nil
}
