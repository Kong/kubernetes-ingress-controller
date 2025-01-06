package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=kl,categories=kong-ingress-controller,path=konglicenses
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.enabled`,description="Enabled to configure on Kong gateway instances"

// KongLicense stores a Kong enterprise license to apply to managed Kong gateway instances.
type KongLicense struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// RawLicenseString is a string with the raw content of the license.
	RawLicenseString string `json:"rawLicenseString"`
	// Enabled is set to true to let controllers (like KIC or KGO) to reconcile it.
	// Default value is true to apply the license by default.
	// +kubebuilder:default=true
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
	// Conditions describe the current conditions of the KongLicense on the controller.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={{type: "Programmed", status: "Unknown", reason:"Pending", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// Group refers to a Kubernetes Group. It must either be an empty string or a
// RFC 1123 subdomain.
// +kubebuilder:validation:MaxLength=253
// +kubebuilder:validation:Pattern=`^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
type Group string

// Kind refers to a Kubernetes kind.
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
type Kind string

// Namespace refers to a Kubernetes namespace. It must be a RFC 1123 label.
// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
type Namespace string

// ObjectName refers to the name of a Kubernetes object.
// Object names can have a variety of forms, including RFC1123 subdomains,
// RFC 1123 labels, or RFC 1035 labels.
//
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=253
type ObjectName string

type ControllerReference struct {
	// Group is the group of referent.
	// It should be empty if the referent is in "core" group (like pod).
	Group *Group `json:"group,omitempty"`
	// Kind is the kind of the referent.
	// By default the nil kind means kind Pod.
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
