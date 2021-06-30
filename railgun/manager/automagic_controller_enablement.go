package manager

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

// autoCTRL will automatically enable or disable controllers based on whether they've been
// explicitly disabled (e.g. util.EnablementStatusDisabled) and whether or not the API for
// that controller is available and loaded in the Kubernetes API.
func autoCTRL(mgr manager.Manager, log logr.Logger, c *config.Config) {
	knativeGVR := schema.GroupVersionResource{
		Group:    knativev1alpha1.SchemeGroupVersion.Group,
		Version:  knativev1alpha1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	if c.KnativeIngressEnabled != util.EnablementStatusDisabled {
		if ctrlutils.CRDExists(mgr.GetClient(), knativeGVR) {
			log.Info("ingresses.networking.internal.knative.dev v1alpha1 CRD available on cluster.")
			c.KnativeIngressEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(`ingresses.networking.internal.knative.dev v1alpha1 CRD not available on cluster.
				 Disabling Knative controller`)
			c.KnativeIngressEnabled = util.EnablementStatusDisabled
		}
	}

	kongClusterPluginGVR := schema.GroupVersionResource{
		Group:    konghqcomv1.SchemeGroupVersion.Group,
		Version:  konghqcomv1.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	if c.KongClusterPluginEnabled != util.EnablementStatusDisabled {
		if ctrlutils.CRDExists(mgr.GetClient(), kongClusterPluginGVR) {
			log.Info("kongclusterplugins.configuration.konghq.com v1beta1 CRD available on cluster.")
			c.KongClusterPluginEnabled = util.EnablementStatusEnabled
		} else {
			log.Info(`kongclusterplugins.configuration.konghq.com v1beta1 CRD not available on cluster.
				 Disabling KongClusterPlugin controller`)
			c.KongClusterPluginEnabled = util.EnablementStatusDisabled
		}
	}

	return
}
