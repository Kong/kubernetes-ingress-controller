package labels

import "github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"

const (
	// LabelPrefix is the string used at the beginning of KIC-specific labels.
	LabelPrefix = annotations.AnnotationPrefix

	// CredentialKey is the key used to indicate a Secret's credential type.
	CredentialKey = "/credential" //nolint:gosec

	// CredentialTypeLabel is the label used to indicate a Secret's credential type.
	CredentialTypeLabel = LabelPrefix + CredentialKey
)
