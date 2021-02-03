package inputs

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ConfigSecretName = "kong-config"

// getOrCreateConfigSecret finds or creates the secret which houses the combined configurations of the cluster
// for eventual parsing and emitting to the Kong Admin API on the proxy instances.
func getOrCreateConfigSecret(ctx context.Context, c client.Client, ns string) (*corev1.Secret, bool, error) {
	secret := new(corev1.Secret)
	if err := c.Get(ctx, types.NamespacedName{Namespace: ns, Name: ConfigSecretName}, secret); err != nil {
		if errors.IsNotFound(err) {
			secret.SetName(ConfigSecretName)
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
