package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// KongCustomEntityList is a list of kong custom entities.
type KongCustomEntityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongCustomEntity `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kongcustomentities,shortName=kce,categories=kong-ingress-controller

// KongCustomEntity represents a custom entity in Kong.
type KongCustomEntity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:Required
	// Type is the type of this custom entity.
	// Should be same as the `Name` of a KongCustomEntityDefinition.
	Type string `json:"type"`
	// Fields is the list of fields in the entity.
	Fields apiextensionsv1.JSON `json:"fields"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=kongcustomentitydefinitions,shortName=kced,categories=kong-ingress-controller

// KongCustomEntityDefinition represents definition of a custom entity type in Kong.
type KongCustomEntityDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:Required
	// Name is the type name of the entity.
	Name string `json:"name"`
	// AdminAPIName is the name used in admin API paths to CRUD this type of entity.
	// If AdminAPIName is empty, it is seen as same as `Name`.
	// For example: Name = "jwt_credentials" and AdminAPIName = "jwts", then we call `/jwts` or `/jwts/<id>` for CRUD.
	AdminAPIName string `json:"adminAPIName,omitempty"`
	// Dependecies are the entity types which are required by this type.
	// If it is empty, the entity type is a "top level" object that does not dependent on other entities.
	Dependecies []KongEntityForeignKey `json:"dependencies,omitempty"`
}

// +k8s:deepcopy-gen:

// KongEntityForeignKey represents a foreign key constraint of Kong entity.
type KongEntityForeignKey struct {
	// +kubebuilder:validation:Required
	// Type is the type of the dependent entity in the foreign key constraint.
	Type string `json:"type"`
	// +kubebuilder:validation:Required
	// PrimaryKey is the primary key to identify the foreign dependency, like "id" in service.
	PrimaryKey string `json:"primaryKey"`
	// AlternativeKeys are other fields that could identify the foreign dependency, like "name" in service.
	AlternativeKeys []string `json:"alternativeKeys,omitempty"`
}

// +kubebuilder:object:root=true

// KongCustomEntityDefinitionList is the list of KongCustomEntityDefinitions.
type KongCustomEntityDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KongCustomEntityDefinition `json:"items"`
}
