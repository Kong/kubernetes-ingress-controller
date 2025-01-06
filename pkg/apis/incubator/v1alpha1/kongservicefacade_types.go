package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KongServiceFacadeKind is the string value representing the KongServiceFacade kind in Kubernetes.
	KongServiceFacadeKind = "KongServiceFacade"
)

func init() {
	SchemeBuilder.Register(&KongServiceFacade{}, &KongServiceFacadeList{})
}

// KongServiceFacade allows creating separate Kong Services for a single Kubernetes
// Service. It can be used as Kubernetes Ingress' backend (via its path's `backend.resource`
// field). It's designed to enable creating two "virtual" Services in Kong that will point
// to the same Kubernetes Service, but will have different configuration (e.g. different
// set of plugins, different load balancing algorithm, etc.).
//
// KongServiceFacade requires `kubernetes.io/ingress.class` annotation with a value
// matching the ingressClass of the Kong Ingress Controller (`kong` by default) to be reconciled.
//
// +genclient
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories=kong-ingress-controller
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type KongServiceFacade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KongServiceFacadeSpec   `json:"spec"`
	Status            KongServiceFacadeStatus `json:"status,omitempty"`
}

// KongServiceFacadeList contains a list of KongServiceFacade.
// +kubebuilder:object:root=true
type KongServiceFacadeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongServiceFacade `json:"items"`
}

// KongServiceFacadeSpec defines the desired state of KongServiceFacade.
type KongServiceFacadeSpec struct {
	// Backend is a reference to a Kubernetes Service that is used as a backend
	// for this Kong Service Facade.
	// +kubebuilder:validation:Required
	Backend KongServiceFacadeBackend `json:"backendRef"`
}

// KongServiceFacadeBackend is a reference to a Kubernetes Service
// that is used as a backend for a Kong Service Facade.
type KongServiceFacadeBackend struct {
	// Name is the name of the referenced Kubernetes Service.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Port is the port of the referenced Kubernetes Service.
	// +kubebuilder:validation:Required
	Port int32 `json:"port"`
}

// KongServiceFacadeStatus defines the observed state of KongServiceFacade.
type KongServiceFacadeStatus struct {
	// Conditions describe the current conditions of the KongServiceFacade.
	//
	// Known condition types are:
	//
	// * "Programmed"
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={{type: "Programmed", status: "Unknown", reason:"Pending", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
