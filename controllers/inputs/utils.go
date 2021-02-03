package inputs

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const configSecretName = "kong-config"

// getOrCreateConfigSecret finds or creates the secret which houses the combined ingress of the cluster for
// eventual parsing and emitting to the Kong Admin API on the proxy instances.
func getOrCreateConfigSecret(ctx context.Context, c client.Client, ns string) (*corev1.Secret, bool, error) {
	secret := new(corev1.Secret)
	if err := c.Get(ctx, types.NamespacedName{Namespace: ns, Name: configSecretName}, secret); err != nil {
		if errors.IsNotFound(err) {
			secret.SetName(configSecretName)
			secret.SetNamespace(ns)
			if err := c.Create(ctx, secret); err != nil {
				return nil, false, err
			}
			return secret, true, nil
		} else {
			return nil, false, err
		}
	}
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	return secret, false, nil
}

// keyFor provides the string key that should be used to store the YAML contents of an object in the
// Kong configuration secret.
func keyFor(obj runtime.Object, nsn types.NamespacedName) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s-%s-%s-%s-%s", gvk.Group, gvk.Version, gvk.Kind, nsn.Namespace, nsn.Name)
}

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
	log.Info("kong configuration secret found", "namespace", nsn.Namespace, "name", configSecretName)

	// get the storage key for this ingress object
	// TODO: patch instead of update for perf
	// TODO: check before overriding
	key := keyFor(obj, nsn)
	secret.Data[key] = cfg
	if err := c.Update(ctx, secret); err != nil {
		if errors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	log.Info("kong configuration patched (TODO: not doing a PATCH yet doing an UPDATE!)", "namespace", nsn.Namespace, "name", configSecretName)
	return ctrl.Result{}, nil
}

// cleanupIngress ensures that a deleted ingress resource is no longer present in the kong configuration secret.
func cleanupIngress(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	// WIP
	return ctrl.Result{}, nil
}
