package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigSource is a wrapper around SecretValueFromSource
type ConfigSource struct {
	metav1.TypeMeta `json:",inline"`
	SecretValue     SecretValueFromSource `json:"secretKeyRef,omitempty"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource
type NamespacedConfigSource struct {
	metav1.TypeMeta `json:",inline"`
	SecretValue     NamespacedSecretValueFromSource `json:"secretKeyRef,omitempty"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:validation:Optional

// SecretValueFromSource represents the source of a secret value
type SecretValueFromSource struct {
	metav1.TypeMeta `json:",inline"`
	// the secret containing the key
	//+kubebuilder:validation:Required
	Secret string `json:"name,omitempty"`
	// the key containing the value
	//+kubebuilder:validation:Required
	Key string `json:"key,omitempty"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:validation:Optional

// NamespacedSecretValueFromSource represents the source of a secret value specifying the secret namespace
type NamespacedSecretValueFromSource struct {
	metav1.TypeMeta `json:",inline"`
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
