package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const configSecretName = "kong-config"

func getOrCreateConfigSecret(ctx context.Context, c client.Client, nsn types.NamespacedName) (*corev1.Secret, bool, error) {
	secret := new(corev1.Secret)
	if err := c.Get(ctx, nsn, secret); err != nil {
		if errors.IsNotFound(err) {
			secret.SetName(configSecretName)
			secret.SetNamespace(nsn.Namespace)
			if err := c.Create(ctx, secret); err != nil {
				return nil, true, err
			}
		} else {
			return nil, false, err
		}
	}
	return secret, true, nil
}
