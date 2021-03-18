package configuration

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/configsecret"
)

// hasFinalizer is a helper function to check whether a client.Object
// already has a specific finalizer set.
func hasFinalizer(obj client.Object, finalizer string) bool {
	hasFinalizer := false
	for _, finalizer := range obj.GetFinalizers() {
		if finalizer == finalizer {
			hasFinalizer = true
		}
	}
	return hasFinalizer
}

// isAPIAvailable is a hack to short circuit controllers for APIs which aren't available on the cluster,
// enabling us to keep separate logic and logging for some legacy API versions.
func isAPIAvailable(mgr ctrl.Manager, obj client.Object) (bool, error) {
	if err := mgr.GetAPIReader().Get(context.Background(), client.ObjectKey{Namespace: controllers.DefaultNamespace, Name: "non-existent"}, obj); err != nil {
		if strings.Contains(err.Error(), "no matches for kind") {
			return false, nil
		}
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return true, nil
}

// storeObjUpdates reconciles storing the YAML contents of Ingress resources (which are managed by Kong)
// from multiple versions which remain supported.
func storeObjUpdates(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO need EVENTS here
	// TODO need more status updates

	// if this is an Ingress resource make sure it's managed by Kong
	if obj.GetObjectKind().GroupVersionKind().Kind == "Ingress" {
		if !isManaged(obj.GetAnnotations()) {
			return ctrl.Result{}, nil
		}
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
	key := configsecret.KeyFor(obj, nsn)
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

// cleanupObj ensures that a deleted ingress resource is no longer present in the kong configuration secret.
func cleanupObj(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
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

	key := configsecret.KeyFor(obj, nsn)
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
