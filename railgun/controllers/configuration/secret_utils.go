package configuration

import (
	"bytes"
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/configsecret"
)

// -----------------------------------------------------------------------------
// Secret Utils - Storage
// -----------------------------------------------------------------------------

// storeIngressObj reconciles storing the YAML contents of Ingress resources (which are managed by Kong)
// from multiple versions which remain supported.
func storeIngressObj(ctx context.Context, c client.Client, log logr.Logger, targetNsn, ingressNsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO need EVENTS here
	// TODO need more status updates
	// TODO: (shane) I want to refactor this into several smaller functions
	// TODO: collapse nsn + obj, this is redudant as obj includes nsn
	// ^ follow up for these items is in: https://github.com/Kong/kubernetes-ingress-controller/issues/1094

	// if this is an Ingress resource make sure it's managed by KIC
	// BUG: this takes only the kind into account, not the API group.
	if obj.GetObjectKind().GroupVersionKind().Kind == "Ingress" {
		if !isManaged(obj.GetAnnotations()) {
			return ctrl.Result{}, nil
		}
	}

	// get the configuration secret
	secret, created, err := getOrCreateConfigSecret(ctx, c, targetNsn)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("kong configuration secret was created elsewhere retrying", "namespace", ingressNsn.Namespace, "ingress", ingressNsn.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if created {
		log.Info("kong configuration did not exist, was created successfully", "namespace", ingressNsn.Namespace, "ingress", ingressNsn.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	// before we store configuration data for this Ingress object, ensure that it has our finalizer set
	if !hasFinalizer(obj, KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", ingressNsn.Namespace, ingressNsn.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, KongIngressFinalizer))
		if err := c.Update(ctx, obj); err != nil { // TODO: patch here instead of update
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// store the ingress record
	if err := storeRuntimeObject(ctx, c, secret, obj, ingressNsn); err != nil {
		if errors.IsConflict(err) {
			log.Error(err, "object updated while reconcilation was running, retrying", ingressNsn.Namespace, ingressNsn.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	log.Info("kong secret configuration successfully patched patched", "namespace", ingressNsn.Namespace, "name", ingressNsn.Name)
	return ctrl.Result{}, nil
}

// isRuntimeObjectSame indicates whether a runtime.Object you intend to store in the configuration secret is the same as what's already stored.
// This can be used to decide whether or not an update needs to be performed on the configuration secret.
func isRuntimeObjectSame(secret *corev1.Secret, obj runtime.Object, nsn types.NamespacedName) (bool, error) {
	// marshal to YAML to check contents
	cfg, err := yaml.Marshal(obj)
	if err != nil {
		return false, err
	}

	// check if there's any existing object
	key := configsecret.KeyFor(obj, nsn)
	foundCFG, ok := secret.Data[key]
	return ok && bytes.Equal(foundCFG, cfg), nil
}

// storeRuntimeObject stores a runtime.Object in the configuration secret. Callers should re-queue after this completes successfully.
func storeRuntimeObject(ctx context.Context, c client.Client, secret *corev1.Secret, obj runtime.Object, nsn types.NamespacedName) error {
	// marshal to YAML for storage
	cfg, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	// patch the secret with the runtime.Object contents
	key := configsecret.KeyFor(obj, nsn)
	secret.Data[key] = cfg

	return c.Update(ctx, secret) // TODO: patch here instead of update for perf
}

// cleanupObj ensures that a deleted ingress resource is no longer present in the kong configuration secret.
func cleanupObj(ctx context.Context, c client.Client, log logr.Logger, targetNsn, ingressNsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// TODO need EVENTS here
	// TODO need more status updates
	// TODO: (shane) I want to refactor this into several smaller functions
	// ^ follow up for these items is in: https://github.com/Kong/kubernetes-ingress-controller/issues/1094

	// grab the configuration secret from the API
	secret := new(corev1.Secret)
	if err := c.Get(ctx, targetNsn, secret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	key := configsecret.KeyFor(obj, ingressNsn)
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
