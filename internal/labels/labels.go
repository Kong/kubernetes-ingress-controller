package labels

import "github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"

const (
	// LabelPrefix is the string used at the beginning of KIC-specific labels.
	LabelPrefix = annotations.AnnotationPrefix

	// CredentialKey is the key used to indicate a Secret's credential type.
	CredentialKey = "/credential" //nolint:gosec

	// ValidateKey is the key used to indicate a Secret contains plugin configuration.
	ValidateKey = "/validate"

	// CredentialTypeLabel is the label used to indicate a Secret's credential type.
	CredentialTypeLabel = LabelPrefix + CredentialKey

	// ValidateLabel is applied to plugins used for plugin configuration to allow the admission webhook to check
	// updates to them.
	ValidateLabel = LabelPrefix + ValidateKey
)

// ValidateType indicates the type of validation applied to a Secret.
type ValidateType string

const (
	// PluginValidate indicates a labeled Secret's contents require plugin configuration validation.
	PluginValidate ValidateType = "plugin"
)
