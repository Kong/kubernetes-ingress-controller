package configuration

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// -----------------------------------------------------------------------------
// Public Types & Functions
// -----------------------------------------------------------------------------

// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted Ingress resources.
const KongIngressFinalizer = "configuration.konghq.com/ingress"

// SetupIngressControllers validates which ingress controllers need to be configured and sets them up with the
// provided controller manager.
func SetupIngressControllers(mgr ctrl.Manager) error {
	netV1Ing := new(netv1.Ingress)
	apiAvailable, err := IsAPIAvailable(mgr, netV1Ing)
	if err != nil {
		return err
	}
	if apiAvailable {
		if err = (&NetV1IngressReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Ingress"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
		}
	} else {
		mgr.GetLogger().Info("API networking.k8s.io/v1/Ingress is not available, skipping controller for this API")
	}

	if !apiAvailable {
		return setupLegacyIngressControllers(mgr)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Private Types & Functions
// -----------------------------------------------------------------------------

// isManaged verifies whether an Ingress resource is managed by Kong controllers by verifying the
// annotations of the object.
// TODO: add these filters to watch options instead!
func isManaged(annotations map[string]string) bool {
	class, ok := annotations["kubernetes.io/ingress.class"]
	if !ok {
		return false
	}
	if class == "kong" {
		return true
	}
	return false
}

// setupLegacyIngressControllers automates controller setup for Ingress resources, but since we
// support some older Ingress versions some decisions need to be made about which controllers run.
// For instance, if networking.k8s.io/v1beta1/Ingress is available, no controller is needed for
// apiextensions.k8s.io/v1beta1/Ingress as the latter will be converted to the former.
func setupLegacyIngressControllers(mgr ctrl.Manager) error {
	// start the networking.k8s.io/v1beta1/Ingress controller (if the API is available)
	netV1Beta1IngAvailable := false
	netV1Beta1Ing := new(netv1beta1.Ingress)
	apiAvailable, err := IsAPIAvailable(mgr, netV1Beta1Ing)
	if err != nil {
		return err
	}
	if apiAvailable {
		if err = (&NetV1Beta1IngressReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("V1Beta1Ingress"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			return err
		}
		netV1Beta1IngAvailable = true
	} else {
		mgr.GetLogger().Info("API networking.k8s.io/v1beta1/Ingress is not available, skipping controller for this API")
	}

	// start the apiextensions.k8s.io/v1beta1/Ingress controller (if the API is available)
	if netV1Beta1IngAvailable {
		mgr.GetLogger().Info("No need to start a controller for apiextensions.k8s.io/v1beta1/Ingress as networking.k8s.io/v1beta1/Ingress covers for this through conversions")
	} else {
		extV1Beta1Ing := new(extv1beta1.Ingress)
		apiAvailable, err = IsAPIAvailable(mgr, extV1Beta1Ing)
		if err != nil {
			return err
		}
		if apiAvailable {
			if err = (&ExtV1Beta1IngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("ExtensionsV1Beta1Ingress"),
				Scheme: mgr.GetScheme(),
			}).SetupWithManager(mgr); err != nil {
				return err
			}
		} else {
			mgr.GetLogger().Info("API apiextensions.k8s.io/v1beta1/Ingress is not available, skipping controller for this API")
		}
	}

	return nil
}
