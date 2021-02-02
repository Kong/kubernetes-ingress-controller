package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const configSecretName = "kong-config"

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

func keyFor(obj runtime.Object, nsn types.NamespacedName) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s-%s-%s-%s-%s", gvk.Group, gvk.Version, gvk.Kind, nsn.Namespace, nsn.Name)
}

// TODO: add these filters to watch options instead!
func isManaged(metadata metav1.ObjectMeta) bool {
	class, ok := metadata.Annotations["kubernetes.io/ingress-class"]
	if !ok {
		return false
	}
	if class == "kong" {
		return true
	}
	return false
}
