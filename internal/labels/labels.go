package labels

import "github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"

const (
	// LabelPrefix is the string used at the beginning of KIC-specific labels.
	LabelPrefix = annotations.AnnotationPrefix

	// CredentialKey is the key used to indicate a Secret's credential type.
	CredentialKey = "/credential" //nolint:gosec
)
