package labels

import "github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"

const (
	// LabelPrefix is the string used at the beginning of KIC-specific labels.
	LabelPrefix = annotations.AnnotationPrefix

	// CredentialKey is the key used to indicate a Secret's credential type.
	CredentialKey = "/credential" //nolint:gosec

	// PluginConfigKey is the key used to indicate a Secret contains plugin configuration.
	PluginConfigKey = "/plugin-config"

	// CredentialTypeLabel is the label used to indicate a Secret's credential type.
	CredentialTypeLabel = LabelPrefix + CredentialKey

	// PluginConfigLabel is applied to plugins used for plugin configuration to allow the admission webhook to check
	// updates to them.
	PluginConfigLabel = LabelPrefix + PluginConfigKey
)
