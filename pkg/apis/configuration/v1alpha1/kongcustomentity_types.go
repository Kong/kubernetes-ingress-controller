package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
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
// +kubebuilder:subresource:status

// KongCustomEntity represents a custom entity in Kong.
type KongCustomEntity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec is the specification of the entity.
	Spec KongCustomEntitySpec `json:"spec"`
	// Status is the status of the entity.
	Status KongCustomEntityStatus `json:"status"`
}

// +k8s:deepcopy-gen:

// KongCustomEntitySpec defines the specification of a Kong custom entity.
type KongCustomEntitySpec struct {
	IngressClassName string `json:"ingressClass"`
	// +kubebuilder:validation:Required
	// Type is the type of this custom entity.
	// Should be same as the `Name` of a KongCustomEntityDefinition.
	Type string `json:"type"`
	// Fields is the list of fields in the entity.
	Fields []KongCustomEntityField `json:"fields"`
}

// KongEntityFieldType defines possible type of field in Kong custom entity.
type KongEntityFieldType string

const (
	KongEntityFieldTypeNil     KongEntityFieldType = "nil"
	KongEntityFieldTypeBoolean KongEntityFieldType = "bool"
	KongEntityFieldTypeInteger KongEntityFieldType = "int"
	KongEntityFieldTypeNumber  KongEntityFieldType = "number"
	KongEntityFieldTypeString  KongEntityFieldType = "string"
	KongEntityFieldTypeArray   KongEntityFieldType = "array"
	KongEntityFieldTypeObject  KongEntityFieldType = "object"
)

// +k8s:deepcopy-gen:

// KongCustomEntityField defines one field of Kong custom entity.
type KongCustomEntityField struct {
	// Key is the key of the entity field.
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=nil;bool;int;number;string;array;object
	// Type is the type of the value in the field.
	Type KongEntityFieldType `json:"type"`
	// +kubebuilder:validation:Type=object
	// Value defines the value of this field in JSON format.
	Value     apiextensionsv1.JSON `json:"value,omitempty"`
	ValueFrom *kongv1.ConfigSource `json:"valueFrom,omitempty"`
}

// +k8s:deepcopy-gen:

// KongCustomEntityStatus defines the status of a Kong custom entity.
type KongCustomEntityStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=kongcustomentitydefinitions,shortName=kced,categories=kong-ingress-controller

// KongCustomEntityDefinition represents definition of a custom entity type in Kong.
type KongCustomEntityDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KongCustomEntityDefinitionSpec `json:"spec"`
}

// +k8s:deepcopy-gen:

// KongCustomEntityDefinitionSpec is the specification of KongCustomEntityDefinition
// to define a type of Kong custom entity.
type KongCustomEntityDefinitionSpec struct {
	// +kubebuilder:validation:Required
	// Name is the type name of the entity.
	Name string `json:"name"`
	// AdminAPIName is the name used in admin API paths to CRUD this type of entity.
	// If AdminAPIName is empty, it uses the value of `Name`.
	// For example: Name = "jwt_credentials" and AdminAPIName = "jwts", then we call `/jwts` or `/jwts/<id>` for CRUD.
	AdminAPIName string `json:"adminAPIName,omitempty"`
	// AdminAPINestedName is the name used in the admin API paths to CRUD the entity attached to other entities.
	// If AdminAPINestedName is empty, it is the same as `AdminAPIName`; if they are both empty, it uses the value of `Name`.
	// like Name = "hmacauth_credentials", AdminAPIName = "hmac-auths" and AdminAPINestedName = "hmac-auth"
	// We call `/consumers/*/hmac-auth` or  `/consumers/*/hmac-auth/*` for CRUD.
	AdminAPINestedName string `json:"adminAPINestedName,omitempty"`
	// Dependecies are the entity types which are required by this type.
	// If it is empty, the entity type is a "top level" object that does not dependent on other entities.
	Dependecies []KongEntityForeignKey `json:"dependencies,omitempty"`
	// Schema apiextensions.JSONSchemaDefinitions `json:"schema,omitempty"`
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
