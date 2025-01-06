package util

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

// ExtractKongCredentialType returns the credential type of a Secret or an error if no credential type is present.
func ExtractKongCredentialType(secret *corev1.Secret) (string, error) {
	credType, labelOk := secret.Labels[labels.CredentialTypeLabel]
	if !labelOk {
		return "", fmt.Errorf("Secret %s/%s used as credential, but lacks %s label",
			secret.Namespace, secret.Name, labels.CredentialTypeLabel)
	}
	return credType, nil
}
