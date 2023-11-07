package v1

// ConfigSource is a wrapper around SecretValueFromSource.
// +kubebuilder:object:generate=true
type ConfigSource struct {
	// Specifies a name and a key of a secret to refer to. The namespace is implicitly set to the one of referring object.
	SecretValue SecretValueFromSource `json:"secretKeyRef"`
}

// NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource.
// +kubebuilder:object:generate=true
type NamespacedConfigSource struct {
	// Specifies a name, a namespace, and a key of a secret to refer to.
	SecretValue NamespacedSecretValueFromSource `json:"secretKeyRef"`
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
