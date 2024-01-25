package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
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
	KongLicenseParentStatuses []KongLicenseParentStatus `json:"parents,omitempty"`
}

// KongLicenseParentStatus is the status of owning KongLicense being processed in the controller in ControllerRef field.
type KongLicenseParentStatus struct {
	// ControllerRef is the reference of the "controller" to reconcile this KongLicense.
	// It is usually the name of (KIC/KGO) pod that reconciles it.
	ControllerRef ControllerReference `json:"controllerRef"`
	// Configured is set to true if the controller applied the content of the license on managed Kong gateway.
	Configured bool `json:"configured"`
	// Phase is the phase of the KongLicense being reconciled on the controller present in ControllerRef.
	Phase KongLicensePhase `json:"phase"`
	// Reason is the reason why the KongLicense stays in this phase.
	Reason string `json:"reason"`
	// TODO: add a field to annotate the controller type?
}

type ControllerReference struct {
	// Group is the group of referent.
	// It should be empty if the referent is in "core" group (like pod.)
	Group *gatewayv1.Group `json:"group,omitempty"`
	// Kind is the kind of the referent.
	Kind *gatewayv1.Kind `json:"kind,omitempty"`
	// Namespace is the namespace of the referent.
	// It should be empty if the referent is cluster scoped.
	Namespace *gatewayv1.Namespace `json:"namespace,omitempty"`
	// Name is the name of the referent.
	Name gatewayv1.ObjectName `json:"name"`
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
