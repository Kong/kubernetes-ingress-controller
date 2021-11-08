package manager

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/proxy"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Controller Definition Interfaces
// -----------------------------------------------------------------------------

// Controller is a Kubernetes controller that can be plugged into Manager.
type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

// AutoHandler decides whether the specific controller shall be enabled (true) or disabled (false).
type AutoHandler func(client.Client) bool

// ControllerDef is a specification of a Controller that can be conditionally registered with Manager.
type ControllerDef struct {
	Enabled     bool
	AutoHandler AutoHandler
	Controller  Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// MaybeSetupWithManager runs SetupWithManager on the controller if it is enabled
// and its AutoHandler (if any) indicates that it can load
func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	if !c.Enabled {
		return nil
	}

	if c.AutoHandler != nil {
		if enable := c.AutoHandler(mgr.GetClient()); !enable {
			return nil
		}
	}
	return c.Controller.SetupWithManager(mgr)
}

// -----------------------------------------------------------------------------
// Controller Manager - Controller Setup Functions
// -----------------------------------------------------------------------------

func setupControllers(mgr manager.Manager, proxy proxy.Proxy, c *Config, featureGates map[string]bool) ([]ControllerDef, error) {
	// Choose the best API version of Ingress to inform which ingress controller to enable.
	var ingressPicker ingressControllerStrategy
	if err := ingressPicker.Initialize(c, mgr.GetClient()); err != nil {
		return nil, fmt.Errorf("ingress version picker failed: %w", err)
	}

	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled:     c.IngressNetV1Enabled,
			AutoHandler: ingressPicker.IsNetV1,
			Controller: &configuration.NetV1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled:     c.IngressNetV1beta1Enabled,
			AutoHandler: ingressPicker.IsNetV1beta1,
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled:     c.IngressExtV1beta1Enabled,
			AutoHandler: ingressPicker.IsExtV1beta1,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			Enabled: true,
			Controller: &configuration.CoreV1SecretReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Secrets"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		// ---------------------------------------------------------------------------
		// Kong API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: c.UDPIngressEnabled,
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled: c.TCPIngressEnabled,
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled: c.KongIngressEnabled,
			Controller: &configuration.KongV1KongIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			Enabled: c.KongPluginEnabled,
			Controller: &configuration.KongV1KongPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			Enabled: c.KongConsumerEnabled,
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			Enabled: c.KongClusterPluginEnabled,
			AutoHandler: crdExistsChecker{GVR: schema.GroupVersionResource{
				Group:    konghqcomv1.SchemeGroupVersion.Group,
				Version:  konghqcomv1.SchemeGroupVersion.Version,
				Resource: "kongclusterplugins",
			}}.CRDExists,
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		// ---------------------------------------------------------------------------
		// Other Controllers
		// ---------------------------------------------------------------------------
		{
			// knative is a special case because it existed before we added feature gates functionality
			// for this controller (only) the existing --enable-controller-knativeingress flag overrides
			// any feature gate configuration. See FEATURE_GATES.md for more information.
			Enabled: featureGates["Knative"] || c.KnativeIngressEnabled,
			AutoHandler: crdExistsChecker{GVR: schema.GroupVersionResource{
				Group:    knativev1alpha1.SchemeGroupVersion.Group,
				Version:  knativev1alpha1.SchemeGroupVersion.Version,
				Resource: "ingresses",
			}}.CRDExists,
			Controller: &configuration.Knativev1alpha1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		// ---------------------------------------------------------------------------
		// GatewayAPI Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: featureGates["Gateway"],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "gateways",
				}}.CRDExists,
			Controller: &gateway.GatewayReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("Gateway"),
				Scheme:          mgr.GetScheme(),
				Proxy:           proxy,
				PublishService:  c.PublishService,
				WatchNamespaces: c.WatchNamespaces,
			},
		},
	}

	return controllers, nil
}

// crdExistsChecker verifies whether the resource type defined by GVR is supported by the k8s apiserver.
type crdExistsChecker struct {
	GVR schema.GroupVersionResource
}

// CRDExists returns true iff the apiserver supports the specified group/version/resource.
func (c crdExistsChecker) CRDExists(r client.Client) bool {
	return ctrlutils.CRDExists(r, c.GVR)
}

// ingressControllerStrategy picks the best Ingress API supported by k8s apiserver.
type ingressControllerStrategy struct {
	chosenVersion IngressAPI
}

// Initialize negotiates the best Ingress API version supported by both KIC and the k8s apiserver.
func (s *ingressControllerStrategy) Initialize(cfg *Config, cl client.Client) error {
	var err error
	s.chosenVersion, err = negotiateIngressAPI(cfg, cl)
	return err
}

// IsExtV1beta1 returns true iff the best supported API version is extensions/v1beta1.
func (s *ingressControllerStrategy) IsExtV1beta1(_ client.Client) bool {
	return s.chosenVersion == ExtensionsV1beta1
}

// IsExtV1beta1 returns true iff the best supported API version is networking.k8s.io/v1beta1.
func (s *ingressControllerStrategy) IsNetV1beta1(_ client.Client) bool {
	return s.chosenVersion == NetworkingV1beta1
}

// IsExtV1beta1 returns true iff the best supported API version is networking.k8s.io/v1.
func (s *ingressControllerStrategy) IsNetV1(_ client.Client) bool {
	return s.chosenVersion == NetworkingV1
}
