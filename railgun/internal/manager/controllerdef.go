package manager

import (
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
)

// -----------------------------------------------------------------------------
// Controller Manager - Controller Definition Interfaces
// -----------------------------------------------------------------------------

// Controller is a Kubernetes controller that can be plugged into Manager.
type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

// AutoHandler decides whether the specific controller shall be enabled (true) or disabled (false).
type AutoHandler func(client.Reader) bool

// ControllerDef is a specification of a Controller that can be conditionally registered with Manager.
type ControllerDef struct {
	IsEnabled   *util.EnablementStatus
	AutoHandler AutoHandler
	Controller  Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// MaybeSetupWithManager runs SetupWithManager on the controller if its EnablementStatus is either "enabled", or "auto"
// and AutoHandler says that it should be enabled.
func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	switch *c.IsEnabled {
	case util.EnablementStatusDisabled:
		return nil

	case util.EnablementStatusAuto:
		if c.AutoHandler == nil {
			return fmt.Errorf("'auto' enablement not supported for controller %q", c.Name())
		}

		if enable := c.AutoHandler(mgr.GetAPIReader()); !enable {
			return nil
		}
		fallthrough

	default: // controller enabled
		return c.Controller.SetupWithManager(mgr)
	}
}

// -----------------------------------------------------------------------------
// Controller Manager - Controller Setup Functions
// -----------------------------------------------------------------------------

func setupControllers(logger logr.Logger, mgr manager.Manager, proxy proxy.Proxy, c *Config) []ControllerDef {
	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			IsEnabled: &c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: &c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: &alwaysEnabled,
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
			IsEnabled: &c.UDPIngressEnabled,
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: &c.TCPIngressEnabled,
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: &c.KongIngressEnabled,
			Controller: &configuration.KongV1KongIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: &c.KongPluginEnabled,
			Controller: &configuration.KongV1KongPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: &c.KongConsumerEnabled,
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
	}

	// ---------------------------------------------------------------------------
	// Dynamic Controller Configurations
	// ---------------------------------------------------------------------------

	// Negotiate Ingress version
	ingressControllers := map[IngressAPI]ControllerDef{
		NetworkingV1: {
			IsEnabled: &c.IngressNetV1Enabled,
			Controller: &configuration.NetV1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		NetworkingV1beta1: {
			IsEnabled: &c.IngressNetV1beta1Enabled,
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		ExtensionsV1beta1: {
			IsEnabled: &c.IngressExtV1beta1Enabled,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
	}

	negotiatedIngressAPI, err := negotiateIngressAPI(c, mgr.GetClient())
	if err == nil {
		controllers = append(controllers, ingressControllers[negotiatedIngressAPI])
	} else {
		logger.Info(`no Ingress controllers enabled or no suitable Ingress version found.
		Disabling Ingress controller`)
	}

	kongClusterPluginGVR := schema.GroupVersionResource{
		Group:    konghqcomv1.SchemeGroupVersion.Group,
		Version:  konghqcomv1.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), kongClusterPluginGVR) == true {
		logger.Info("kongclusterplugins.configuration.konghq.com v1beta1 CRD available on cluster.")
		controller := ControllerDef{
			IsEnabled: &c.KongClusterPluginEnabled,
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		}
		controllers = append(controllers, controller)
	} else {
		logger.Info(`kongclusterplugins.configuration.konghq.com v1beta1 CRD not available on cluster.
		Disabling KongClusterPlugin controller`)
	}

	knativeGVR := schema.GroupVersionResource{
		Group:    knativev1alpha1.SchemeGroupVersion.Group,
		Version:  knativev1alpha1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), knativeGVR) == true {
		logger.Info("ingresses.networking.internal.knative.dev v1alpha1 CRD available on cluster.")
		controller := ControllerDef{
			IsEnabled: &c.KnativeIngressEnabled,
			Controller: &configuration.Knativev1alpha1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		}
		controllers = append(controllers, controller)
	} else {
		logger.Info(`ingresses.networking.internal.knative.dev v1alpha1 CRD not available on cluster.
		Disabling Knative controller`)
		c.KnativeIngressEnabled = util.EnablementStatusDisabled
	}

	return controllers
}
