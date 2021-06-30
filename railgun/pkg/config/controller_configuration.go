package config

import (
	"github.com/go-logr/logr"
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

	knativeGVR := schema.GroupVersionResource{
		Group:    knativev1alpha1.SchemeGroupVersion.Group,
		Version:  knativev1alpha1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if controllerConfig.KnativeIngressEnabled != util.EnablementStatusDisabled {
		if ctrlutils.CRDExists(mgr.GetClient(), knativeGVR) {
			log.Info("ingresses.networking.internal.knative.dev v1alpha1 CRD available on cluster.")
			controllerConfig.KnativeIngressEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(`ingresses.networking.internal.knative.dev v1alpha1 CRD not available on cluster.
				 Disabling Knative controller`)
			controllerConfig.KnativeIngressEnabled = util.EnablementStatusDisabled
		}
	}

	kongClusterPluginGVR := schema.GroupVersionResource{
		Group:    konghqcomv1.SchemeGroupVersion.Group,
		Version:  konghqcomv1.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	if controllerConfig.KongClusterPluginEnabled != util.EnablementStatusDisabled {
		if ctrlutils.CRDExists(mgr.GetClient(), kongClusterPluginGVR) {
			log.Info("kongclusterplugins.configuration.konghq.com v1beta1 CRD available on cluster.")
			controllerConfig.KongClusterPluginEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(`kongclusterplugins.configuration.konghq.com v1beta1 CRD not available on cluster.
				 Disabling KongClusterPlugin controller`)
			controllerConfig.KongClusterPluginEnabled = util.EnablementStatusDisabled
		}
	}

	return &controllerConfig
}
