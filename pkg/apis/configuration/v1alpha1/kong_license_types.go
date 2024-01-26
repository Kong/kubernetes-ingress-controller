package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kl,categories=kong-ingress-controller,path=konglicenses
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// KongLicense stores a Kong enterprise license to apply to managed Kong gateway instances.
type KongLicense struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// RawLicenseString is the raw content of the license in string format.
	RawLicenseString string `json:"rawLicenseString"`
	// Enabled is set to true to let controllers (like KIC) to reconcile it.
	Enabled bool `json:"enabled"`
	// Status is the status of the KongLicense being processed by controllers.
	Status KongLicenseStatus `json:"status,omitempty"`
}

// KongLicenseStatus stores the status of the KongLicense being processesed in each controller that reconciles it.
type KongLicenseStatus struct {
	// +listType=map
	// +listMapKey=controllerName
	KongLicenseControllerStatuses []KongLicenseControllerStatus `json:"controllers,omitempty"`
}

// KongLicenseControllerStatus is the status of owning KongLicense being processed
// identified by the controllerName field.
type KongLicenseControllerStatus struct {
	// ControllerName is an identifier of the controller to reconcile this KongLicense.
	// Should be unique in the list of controller statuses.
	ControllerName string `json:"controllerName"`
	// ControllerRef is the reference of the controller to reconcile this KongLicense.
	// It is usually the name of (KIC/KGO) pod that reconciles it.
	ControllerRef *ControllerReference `json:"controllerRef,omitempty"`
	// Configured is set to true if the controller applied the content of the license on managed Kong gateway.
	Configured bool `json:"configured"`
	// Conditions describe the current conditions of the KongLicense on the controller.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// TODO: copy validation notes from the gateway api package to run the same validation?

// Group refers to a Kubernetes Group. It must either be an empty string or a
// RFC 1123 subdomain.
type Group string

// Kind refers to a kubernetes kind.
type Kind string

// Namespace refers to a Kubernetes namespace.
type Namespace string

// ObjectName refers to the name of a Kubernetes object.
type ObjectName string

type ControllerReference struct {
	// Group is the group of referent.
	// It should be empty if the referent is in "core" group (like pod.)
	Group *Group `json:"group,omitempty"`
	// Kind is the kind of the referent.
	Kind *Kind `json:"kind,omitempty"`
	// Namespace is the namespace of the referent.
	// It should be empty if the referent is cluster scoped.
	Namespace *Namespace `json:"namespace,omitempty"`
	// Name is the name of the referent.
	Name ObjectName `json:"name"`
}

type KongLicensePhase string

// +kubebuilder:object:root=true

// KongLicenseList contains a list of KongLicense.
type KongLicenseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongLicense `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KongLicense{}, &KongLicenseList{})
}
