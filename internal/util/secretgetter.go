package util

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretGetterFromK8s is a SecretGetter that reads secrets from Kubernetes API.
type SecretGetterFromK8s struct {
	Reader client.Reader
}

// GetSecret reads a core v1 Secret from Kubernetes API.
func (s *SecretGetterFromK8s) GetSecret(namespace string, name string) (*corev1.Secret, error) {
	var res corev1.Secret
	err := s.Reader.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, &res)
	return &res, err
}
