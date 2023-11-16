package v1

// ConfigSource is a wrapper around SecretValueFromSource.
// +kubebuilder:object:generate=true
type ConfigSource struct {
	// Specifies a name and a key of a secret to refer to. The namespace is implicitly set to the one of referring object.
	SecretValue SecretValueFromSource `json:"secretKeyRef"`
}

// ConfigPatch is a JSON patch (RFC6902) to add values from Secret to the generated configuration.
// It is an equivalent of the following patch:
// `{"op": "add", "path": {.Path}, "value": {.ComputedValueFrom}}`.
// +kubebuilder:object:generate=true
type ConfigPatch struct {
	// Path is the JSON-Pointer value (RFC6901) that references a location within the target configuration.
	Path string `json:"path"`
	// ValueFrom is the reference to a key of a secret where the patched value comes from.
	ValueFrom ConfigSource `json:"valueFrom"`
}

// NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource.
// +kubebuilder:object:generate=true
type NamespacedConfigSource struct {
	// Specifies a name, a namespace, and a key of a secret to refer to.
	SecretValue NamespacedSecretValueFromSource `json:"secretKeyRef"`
}

// NamespacedConfigPatch is a JSON patch to add values from secrets to KongClusterPlugin
// to the generated configuration of plugin in Kong.
// +kubebuilder:object:generate=true
type NamespacedConfigPatch struct {
	// Path is the JSON path to add the patch.
	Path string `json:"path"`
	// ValueFrom is the reference to a key of a secret where the patched value comes from.
	ValueFrom NamespacedConfigSource `json:"valueFrom"`
}

// SecretValueFromSource represents the source of a secret value.
// +kubebuilder:object:generate=true
type SecretValueFromSource struct {
	// The secret containing the key.
	Secret string `json:"name"`
	// The key containing the value.
	Key string `json:"key"`
}

// NamespacedSecretValueFromSource represents the source of a secret value specifying the secret namespace.
// +kubebuilder:object:generate=true
type NamespacedSecretValueFromSource struct {
	// The namespace containing the secret.
	Namespace string `json:"namespace"`
	// The secret containing the key.
	Secret string `json:"name"`
	// The key containing the value.
	Key string `json:"key"`
}
