package manager

import (
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

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
	IsEnabled   bool
	AutoHandler AutoHandler
	Controller  Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// -----------------------------------------------------------------------------
// Controller Manager - Controller Setup Functions
// -----------------------------------------------------------------------------

// setupControllers mostly just provides a list of controllers which are flagged as enabled
// or disabled based on user configuration, however wherever special handling might be needed
// (particularly when trying to accomodate for older API versions on older K8s clusters) this
// function will take care of configuring for the latest API version of a requested controller.
func setupControllers(logger logr.Logger, mgr manager.Manager, proxy proxy.Proxy, c *Config) ([]ControllerDef, error) {
	// When it comes to upstream Ingress API you only ever need one of them and that will always be the latest version,
	// but we technically have different controllers for the two legacy options which may need to be enabled for some
	// older versions of Kubernetes. if the caller enabled Ingress support, the following call simply determines the
	// latest version available on the cluster and enables that.
	extv1beta1Enabled, netv1beta1Enabled, netv1Enabled, err := determineIngressesAvailable(logger, mgr, c)
	if err != nil {
		return nil, err
	}

	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Networking API Controllers
		// ---------------------------------------------------------------------------
		{
			IsEnabled: netv1Enabled,
			Controller: &configuration.NetV1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: netv1beta1Enabled,
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: extv1beta1Enabled,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			IsEnabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: true,
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
			IsEnabled: c.UDPIngressEnabled,
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: c.TCPIngressEnabled,
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: c.KongIngressEnabled,
			Controller: &configuration.KongV1KongIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: c.KongPluginEnabled,
			Controller: &configuration.KongV1KongPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme: mgr.GetScheme(),
				Proxy:  proxy,
			},
		},
		{
			IsEnabled: c.KongConsumerEnabled,
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: c.KongClusterPluginEnabled,
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
		// ---------------------------------------------------------------------------
		// 3rd Party API Controllers
		// ---------------------------------------------------------------------------
		{
			IsEnabled: c.KnativeIngressEnabled,
			Controller: &configuration.Knativev1alpha1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            proxy,
				IngressClassName: c.IngressClassName,
			},
		},
	}

	return controllers, nil
}

// determineIngressesAvailable will check the cluster for available Ingress APIs from newest to oldest
// and will indicate which controllers should be enabled based on the latest available version.
func determineIngressesAvailable(logger logr.Logger, mgr manager.Manager, c *Config) (extv1beta1Enabled, netv1beta1Enabled, netv1Enabled bool, err error) {
	if !c.IngressEnabled {
		return
	}
	logger.Info("ingress controller enabled: determining latest available ingress api on cluster")

	netv1GVR := schema.GroupVersionResource{
		Group:    netv1.SchemeGroupVersion.Group,
		Version:  netv1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), netv1GVR) {
		logger.Info(fmt.Sprintf("%s was found and is available", netv1GVR.String()))
		netv1Enabled = true
		return
	}

	netv1beta1GVR := schema.GroupVersionResource{
		Group:    netv1beta1.SchemeGroupVersion.Group,
		Version:  netv1beta1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), netv1beta1GVR) {
		logger.Info(fmt.Sprintf("%s was found and is available", netv1beta1GVR.String()))
		netv1beta1Enabled = true
		return
	}

	extv1beta1GVR := schema.GroupVersionResource{
		Group:    extv1beta1.SchemeGroupVersion.Group,
		Version:  extv1beta1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), extv1beta1GVR) {
		logger.Info(fmt.Sprintf("%s was found and is available", extv1beta1GVR.String()))
		extv1beta1Enabled = true
		return
	}

	err = fmt.Errorf("ingress was enabled, but the cluster has no available ingress apis")
	return
}
