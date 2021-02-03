package inputs

import (
	"context"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted Ingress resources.
const KongIngressFinalizer = "networking.konghq.com/ingress"

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

	// get the configuration secret
	secret, created, err := getOrCreateConfigSecret(ctx, c, nsn.Namespace)
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
	log.Info("kong configuration secret found", "namespace", nsn.Namespace, "name", ConfigSecretName)

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

	log.Info("kong configuration patched (TODO: not doing a PATCH yet doing an UPDATE!)", "namespace", nsn.Namespace, "name", ConfigSecretName)
	return ctrl.Result{}, nil
}

// cleanupIngress ensures that a deleted ingress resource is no longer present in the kong configuration secret.
func cleanupIngress(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO remove secret entry

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
