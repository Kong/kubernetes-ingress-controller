package inputs

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/kong/railgun/controllers"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// -----------------------------------------------------------------------------
// Public Types & Functions
// -----------------------------------------------------------------------------

// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted Ingress resources.
const KongIngressFinalizer = "networking.konghq.com/ingress"

// SetupIngressControllers validates which ingress controllers need to be configured and sets them up with the
// provided controller manager.
func SetupIngressControllers(mgr ctrl.Manager) error {
	netV1Ing := new(netv1.Ingress)
	apiAvailable, err := isAPIAvailable(mgr, netV1Ing)
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
	class, ok := annotations["kubernetes.io/ingress-class"]
	if !ok {
		return false
	}
	if class == "kong" {
		return true
	}
	return false
}

// storeIngressUpdates reconciles storing the YAML contents of Ingress resources (which are managed by Kong)
// from multiple versions which remain supported.
func storeIngressUpdates(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO need EVENTS here
	// TODO need more status updates

	anns := obj.GetAnnotations()

	// ensure this Ingress is managed by KONG
	if !isManaged(anns) {
		return ctrl.Result{}, nil
	}

	// marshal to YAML for later storage
	cfg, err := yaml.Marshal(obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// get the configuration secret namespace
	secretNamespace := os.Getenv(controllers.CtrlNamespaceEnv)
	if secretNamespace == "" {
		return ctrl.Result{}, fmt.Errorf("kong can not be configured because the required %s env var is not present", controllers.CtrlNamespaceEnv)
	}

	// get the configuration secret
	secret, created, err := getOrCreateConfigSecret(ctx, c, secretNamespace)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("kong configuration secret was created elsewhere retrying", "namespace", nsn.Namespace, "ingress", nsn.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if created {
		log.Info("kong configuration did not exist, was created successfully", "namespace", nsn.Namespace, "ingress", nsn.Name)
		return ctrl.Result{Requeue: true}, nil
	}
	log.Info("kong configuration secret found", "namespace", nsn.Namespace, "name", controllers.ConfigSecretName)

	// before we store configuration data for this Ingress object, ensure that it has our finalizer set
	if !hasFinalizer(obj, KongIngressFinalizer) {
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, KongIngressFinalizer))
		if err := c.Update(ctx, obj); err != nil { // TODO: patch here instead of update
			return ctrl.Result{}, err
		}
	}

	// get the storage key for this ingress object and update it
	// TODO: check before overriding
	key := keyFor(obj, nsn)
	secret.Data[key] = cfg
	if err := c.Update(ctx, secret); err != nil { // TODO: patch here instead of update for perf
		if errors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	log.Info("kong configuration patched", "namespace", nsn.Namespace, "name", controllers.ConfigSecretName)
	return ctrl.Result{}, nil
}

// cleanupIngress ensures that a deleted ingress resource is no longer present in the kong configuration secret.
func cleanupIngress(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO need EVENTS here
	// TODO need more status updates

	// get the configuration secret namespace
	secretNamespace := os.Getenv(controllers.CtrlNamespaceEnv)
	if secretNamespace == "" {
		return ctrl.Result{}, fmt.Errorf("kong can not be configured because the required %s env var is not present", controllers.CtrlNamespaceEnv)
	}

	// grab the configuration secret from the API
	secret := new(corev1.Secret)
	if err := c.Get(ctx, types.NamespacedName{Namespace: secretNamespace, Name: controllers.ConfigSecretName}, secret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	key := keyFor(obj, nsn)
	if _, ok := secret.Data[key]; ok {
		delete(secret.Data, key)
		if err := c.Update(ctx, secret); err != nil { // TODO: patch here instead of update
			return ctrl.Result{}, err
		}
		log.Info("kong ingress record removed from kong configuration", "ingress", obj.GetName(), "config", secret.GetName())
		return ctrl.Result{Requeue: true}, nil
	}

	if hasFinalizer(obj, KongIngressFinalizer) {
		log.Info("kong ingress finalizer needs to be removed from ingress resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		finalizers := []string{}
		for _, finalizer := range obj.GetFinalizers() {
			if finalizer != KongIngressFinalizer {
				finalizers = append(finalizers, finalizer)
			}
		}
		obj.SetFinalizers(finalizers)
		if err := c.Update(ctx, obj); err != nil { // TODO: patch here instead of update
			return ctrl.Result{}, err
		}
		log.Info("the kong ingress finalizer was removed from an ingress resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// setupLegacyIngressControllers automates controller setup for Ingress resources, but since we
// support some older Ingress versions some decisions need to be made about which controllers run.
// For instance, if networking.k8s.io/v1beta1/Ingress is available, no controller is needed for
// apiextensions.k8s.io/v1beta1/Ingress as the latter will be converted to the former.
func setupLegacyIngressControllers(mgr ctrl.Manager) error {
	// start the networking.k8s.io/v1beta1/Ingress controller (if the API is available)
	netV1Beta1IngAvailable := false
	netV1Beta1Ing := new(netv1beta1.Ingress)
	apiAvailable, err := isAPIAvailable(mgr, netV1Beta1Ing)
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
		apiAvailable, err = isAPIAvailable(mgr, extV1Beta1Ing)
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
