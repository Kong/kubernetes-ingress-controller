package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KongCustomEntityKind = "KongCustomEntity"
)

type KongEntityScope string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kce,categories=kong-ingress-controller,path=kongcustomentities,singular=kongcustomentity
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Entity Type",type=string,JSONPath=`.type`,description="type of the Kong entity"
// +kubebuilder:validation:XValidation:rule="self.scope == oldSelf.scope",message="The scope field is immutable"
// +kubebuilder:validation:XValidation:rule="self.type == oldSelf.type",message="The type field is immutable"
// +kubebuilder:validation:XValidation:rule="!(self.type in ['services','routes','upstreams','targets','plugins','consumers','consumer_groups'])",message="The type field cannot be known Kong entity types"
// +kubebuilder:validation:XValidation:rule="!(self.scope == 'independent' && has(self.parentRef)) && !(self.scope == 'attached' && !has(self.parentRef))",message="attached KongCustomEntity must have parentRef; independent KongCustomEntity must not have parentRef"
// REVIEW: put all fields other than "status" under the "spec"?

// KongCustomEntity defines a "custom" Kong entity that KIC cannot support the entity type directly.
type KongCustomEntity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:Enum=independent;attached
	// REVIEW: We can get the schema of entities from Kong gateway admin APIs including fields requiring foreign reference.
	// Is the "Scope" field really required?

	// Scope means whether the entity can be specified independently.
	// "independent" means that the entity does not depend on other entities;
	// "attached" means that the entity has "foreign" fields referring to other entities.
	Scope KongEntityScope `json:"scope"`

	// EntityType is the type of the Kong entity. The type is used in generating declarative configuration.
	EntityType string `json:"type"`
	// Fields defines the fields of the Kong entity itself.
	Fields apiextensionsv1.JSON `json:"fields"`
	// ControllerName specifies the controller that should reconcile it, like ingress class.
	ControllerName string `json:"controllerName"`

	// ParentRef references the kubernetes resource it attached to when its scope is "attached".
	// Currently only KongPlugin/KongClusterPlugin allowed. This will make the custom entity to be attached
	// to the entity(service/route/consumer) where the plugin is attached.
	ParentRef *ObjectReference `json:"parentRef,omitempty"`
	// Status stores the reconciling status of the resource.
	Status KongCustomEntityStatus `json:"status,omitempty"`
}

// REVIEW:
// - Should we define dedicated type aliases for each field (like gateway API does, define types "Group","Kind","Namespace","ObjectName")?
// - Should we define the optional fields to pointer type?
// - Should we preset a "default" value of Group/Kind when they are not present like KongPlugin?

// ObjectReference defines reference of a kubernetes object.
type ObjectReference struct {
	Group *string `json:"group,omitempty"`
	Kind  *string `json:"kind,omitempty"`
	// Empty namespace means the same namespace of the owning object.
	Namespace *string `json:"namespace,omitempty"`
	Name      string  `json:"name"`
}

type KongCustomEntityStatus struct {
	// Conditions describe the current conditions of the KongCustomEntityStatus.
	//
	// Known condition types are:
	//
	// * "Programmed"
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={{type: "Programmed", status: "Unknown", reason:"Pending", message:"Waiting for controller", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	Conditions []metav1.Condition `json:"conditions"`
}

// +kubebuilder:object:root=true

// KongCustomEntityList contains a list of KongCustomEntity.
type KongCustomEntityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongCustomEntity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KongCustomEntity{}, &KongCustomEntityList{})
}
