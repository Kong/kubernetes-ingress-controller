package util

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4853 remove field handling when no longer supported.

// CredentialTypeSource indicates the source of credential type information (or lack thereof) in a Secret.
type CredentialTypeSource int

const (
	// CredentialTypeAbsent indicates that no credential information is present in a Secret.
	CredentialTypeAbsent CredentialTypeSource = iota
	// CredentialTypeFromLabel indicates that a Secret's credential type was determined from a label.
	CredentialTypeFromLabel
	// CredentialTypeFromField indicates that a Secret's credential type was determined from a data field.
	CredentialTypeFromField
)

// ExtractKongCredentialType returns the credential type of a Secret and a code indicating whether the credential type
// was obtained from a label, field, or not at all. Labels take precedence over fields if both are present.
func ExtractKongCredentialType(secret *corev1.Secret) (string, CredentialTypeSource) {
	credType, labelOk := secret.Labels[labels.CredentialTypeLabel]
	if !labelOk {
		// if no label, fall back to the deprecated field
		credBytes, fieldOk := secret.Data["kongCredType"]
		if !fieldOk {
			return "", CredentialTypeAbsent
		}
		return string(credBytes), CredentialTypeFromField
	}
	return credType, CredentialTypeFromLabel
}
