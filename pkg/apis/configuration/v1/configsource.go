package v1

// ConfigSource is a wrapper around SecretValueFromSource
//+kubebuilder:object:generate=true
type ConfigSource struct {
	SecretValue SecretValueFromSource `json:"secretKeyRef,omitempty"`
}

// NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource
//+kubebuilder:object:generate=true
type NamespacedConfigSource struct {
	SecretValue NamespacedSecretValueFromSource `json:"secretKeyRef,omitempty"`
}

// SecretValueFromSource represents the source of a secret value
//+kubebuilder:object:generate=true
type SecretValueFromSource struct {
	// the secret containing the key
	//+kubebuilder:validation:Required
	Secret string `json:"name,omitempty"`
	// the key containing the value
	//+kubebuilder:validation:Required
	Key string `json:"key,omitempty"`
}

// NamespacedSecretValueFromSource represents the source of a secret value specifying the secret namespace
//+kubebuilder:object:generate=true
type NamespacedSecretValueFromSource struct {
	// The namespace containing the secret
	//+kubebuilder:validation:Required
	Namespace string `json:"namespace,omitempty"`
	// the secret containing the key
	//+kubebuilder:validation:Required
	Secret string `json:"name,omitempty"`
	// the key containing the value
	//+kubebuilder:validation:Required
	Key string `json:"key,omitempty"`
}
