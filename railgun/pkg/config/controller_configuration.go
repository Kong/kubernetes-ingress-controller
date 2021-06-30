package config

import (
	"fmt"

	"github.com/go-logr/logr"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
)

// -----------------------------------------------------------------------------
// Controller Configurations
// -----------------------------------------------------------------------------

// ControllerConfig collects the enablement status of controllers which indicates
// to the controller manager which and which not to run during initialization.
type ControllerConfig struct {
	IngressExtV1beta1Enabled util.EnablementStatus
	IngressNetV1beta1Enabled util.EnablementStatus
	IngressNetV1Enabled      util.EnablementStatus
	UDPIngressEnabled        util.EnablementStatus
	TCPIngressEnabled        util.EnablementStatus
	KongIngressEnabled       util.EnablementStatus
	KnativeIngressEnabled    util.EnablementStatus
	KongClusterPluginEnabled util.EnablementStatus
	KongPluginEnabled        util.EnablementStatus
	KongConsumerEnabled      util.EnablementStatus
	ServiceEnabled           util.EnablementStatus
}

// NewControllerConfigFromManagerConfig will provide a new *ControllerConfig based on a provided *Config.
// Some controllers will be automatically enabled or disabled based on several factors:
//
// - if they've been explicity disabled (e.g. util.EnablementStatusDisabled) they remain disabled
// - if they are configured for auto resolution, and the relevant API is available on the cluster, they are enabled
// - if they are explicitly enabled BUT the API is not available, they are disabled and a warning is logged.
//
// this function is intended for use by the controller manager.Setup() in order to make decisions dynamically
// about controllers so that the controller manager doesn't crash due to missing APIs for controllers.
func NewControllerConfigFromManagerConfig(mgr manager.Manager, log logr.Logger, originalConfig *Config) *ControllerConfig {
	controllerConfig := originalConfig.ControllerConfig
	mgrc := mgr.GetClient()

	// -------------------------------------------------------------------------
	// Kubernetes Core APIs
	// -------------------------------------------------------------------------

	if controllerConfig.IngressNetV1Enabled != util.EnablementStatusDisabled {
		gvr := schema.GroupVersionResource{
			Group:    networkingv1.SchemeGroupVersion.Group,
			Version:  networkingv1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}

		if ctrlutils.CRDExists(mgrc, gvr) {
			log.Info(fmt.Sprintf("%s CRD available on cluster", gvr.String()))
			controllerConfig.IngressNetV1Enabled = util.EnablementStatusEnabled
		} else {
			log.Info(fmt.Sprintf("%s CRD NOT available on cluster: disabling controller", gvr.String()))
			controllerConfig.IngressNetV1Enabled = util.EnablementStatusDisabled
		}
	}

	if controllerConfig.IngressNetV1beta1Enabled != util.EnablementStatusDisabled {
		gvr := schema.GroupVersionResource{
			Group:    networkingv1beta1.SchemeGroupVersion.Group,
			Version:  networkingv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}

		if ctrlutils.CRDExists(mgrc, gvr) {
			log.Info(fmt.Sprintf("%s CRD available on cluster", gvr.String()))
			controllerConfig.IngressNetV1beta1Enabled = util.EnablementStatusEnabled
		} else {
			log.Info(fmt.Sprintf("%s CRD NOT available on cluster: disabling controller", gvr.String()))
			controllerConfig.IngressNetV1beta1Enabled = util.EnablementStatusDisabled
		}
	}

	if controllerConfig.IngressExtV1beta1Enabled != util.EnablementStatusDisabled {
		gvr := schema.GroupVersionResource{
			Group:    extensionsv1beta1.SchemeGroupVersion.Group,
			Version:  extensionsv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}

		if ctrlutils.CRDExists(mgrc, gvr) {
			log.Info(fmt.Sprintf("%s CRD available on cluster", gvr.String()))
			controllerConfig.IngressExtV1beta1Enabled = util.EnablementStatusEnabled
		} else {
			log.Info(fmt.Sprintf("%s CRD NOT available on cluster: disabling controller", gvr.String()))
			controllerConfig.IngressExtV1beta1Enabled = util.EnablementStatusDisabled
		}
	}

	// -------------------------------------------------------------------------
	// Kong APIs
	// -------------------------------------------------------------------------

	if controllerConfig.KongClusterPluginEnabled != util.EnablementStatusDisabled {
		gvr := schema.GroupVersionResource{
			Group:    konghqcomv1.SchemeGroupVersion.Group,
			Version:  konghqcomv1.SchemeGroupVersion.Version,
			Resource: "kongclusterplugins",
		}

		if ctrlutils.CRDExists(mgrc, gvr) {
			log.Info(fmt.Sprintf("%s CRD available on cluster", gvr.String()))
			controllerConfig.KongClusterPluginEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(fmt.Sprintf("%s CRD NOT available on cluster: disabling controller", gvr.String()))
			controllerConfig.KongClusterPluginEnabled = util.EnablementStatusDisabled
		}
	}

	// -------------------------------------------------------------------------
	// 3rd Party APIs
	// -------------------------------------------------------------------------

	if controllerConfig.KnativeIngressEnabled != util.EnablementStatusDisabled {
		gvr := schema.GroupVersionResource{
			Group:    knativev1alpha1.SchemeGroupVersion.Group,
			Version:  knativev1alpha1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}

		if ctrlutils.CRDExists(mgrc, gvr) {
			log.Info(fmt.Sprintf("%s CRD available on cluster", gvr.String()))
			controllerConfig.KnativeIngressEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(fmt.Sprintf("%s CRD NOT available on cluster: disabling controller", gvr.String()))
			controllerConfig.KnativeIngressEnabled = util.EnablementStatusDisabled
		}
	}

	return &controllerConfig
}
