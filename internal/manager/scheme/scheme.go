package scheme

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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

var kicScheme *runtime.Scheme = initScheme()

func initScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	utilruntime.Must(apiextensionsv1.AddToScheme(s))
	utilruntime.Must(clientgoscheme.AddToScheme(s))

	utilruntime.Must(kongv1.AddToScheme(s))
	utilruntime.Must(kongv1alpha1.AddToScheme(s))
	utilruntime.Must(kongv1beta1.AddToScheme(s))
	utilruntime.Must(incubatorv1alpha1.AddToScheme(s))

	utilruntime.Must(gatewayv1.Install(s))
	utilruntime.Must(gatewayv1beta1.Install(s))
	utilruntime.Must(gatewayv1alpha2.Install(s))
	utilruntime.Must(gatewayv1alpha3.Install(s))

	return s
}

// Get returns a scheme aware of all types the manager can interact with.
func Get() *runtime.Scheme {
	return kicScheme
}
