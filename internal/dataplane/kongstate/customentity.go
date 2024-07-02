package kongstate

import (
	"context"
	"sort"

	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

// EntityFieldType represents type of a Kong entity field.
// possible field types include boolean, integer, number, string, array, set, map, record, json, foreign.
type EntityFieldType string

// These types and field properties are defined upstream in the Kong DAO:
// https://github.com/Kong/kong/blob/3.6.1/kong/db/schema/init.lua#L131-L143
// https://docs.konghq.com/gateway/latest/plugin-development/custom-entities/#define-a-schema

const (
	EntityFieldTypeBoolean EntityFieldType = "boolean"
	EntityFieldTypeInteger EntityFieldType = "integer"
	EntityFieldTypeNumber  EntityFieldType = "number"
	EntityFieldTypeString  EntityFieldType = "string"
	EntityFieldTypeSet     EntityFieldType = "set"
	EntityFieldTypeArray   EntityFieldType = "array"
	EntityFieldTypeMap     EntityFieldType = "map"
	EntityFieldTypeRecord  EntityFieldType = "record"
	EntityFieldTypeJSON    EntityFieldType = "json"
	// EntityFieldTypeForeign means that this field refers to another entity by the key (typically ID).
	EntityFieldTypeForeign EntityFieldType = "foreign"
)

type EntityField struct {
	// Name is the name of the field.
	Name string `json:"name"`
	// Type stands for the type of the field.
	Type EntityFieldType `json:"type"`
	// Required is true means that the field must present in the entity.
	Required bool `json:"required,omitempty"`
	// Auto is true means that the field is automatically generated when it is created in Kong gateway.
	Auto bool `json:"auto,omitempty"`
	// UUID is true means that the field is in UUID format.
	UUID bool `json:"uuid,omitempty"`
	// Default is the default value of the field when it is not given.
	Default interface{} `json:"default,omitempty"`
	// Reference is the type referring entity when the field is "foreign" to refer to another entity.
	Reference string `json:"reference,omitempty"`
	// Other attributes in field metadata that do not affect validation and translation are omitted.
}

// EntitySchema is the schema of an entity.
type EntitySchema struct {
	Fields map[string]EntityField
}

// ExtractEntityFieldDefinitions extracts the fields in response of retrieving entity schema from Kong gateway
// and fill the definition of each field in the `Fields` map of returning value.
func ExtractEntityFieldDefinitions(schema kong.Schema) EntitySchema {
	retSchema := EntitySchema{
		Fields: make(map[string]EntityField),
	}

	fieldList, ok := schema["fields"].([]interface{})
	if !ok {
		return retSchema
	}
	for _, item := range fieldList {
		fieldDef, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		for fieldName, fieldAttributes := range fieldDef {
			fieldAttributesMap, ok := fieldAttributes.(map[string]interface{})
			if !ok {
				continue
			}

			fieldTypeStr, _ := fieldAttributesMap["type"].(string)
			fieldAuto, _ := fieldAttributesMap["auto"].(bool)
			fieldIsUUID, _ := fieldAttributesMap["uuid"].(bool)
			fieldRequired, _ := fieldAttributesMap["required"].(bool)
			fieldReference, _ := fieldAttributesMap["reference"].(string)

			f := EntityField{
				Name:      fieldName,
				Type:      EntityFieldType(fieldTypeStr),
				Auto:      fieldAuto,
				UUID:      fieldIsUUID,
				Required:  fieldRequired,
				Default:   fieldAttributesMap["default"],
				Reference: fieldReference,
			}
			retSchema.Fields[fieldName] = f
		}
	}
	return retSchema
}

// IsKnownEntityType returns true if the entities of the type are "standard" and processed elsewhere in KIC.
func IsKnownEntityType(entityType string) bool {
	switch entityType {
	case
		// Types of standard Kong entities that are processed elsewhere in KIC.
		// So the entities cannot be specified via KongCustomEntity types.
		string(kong.EntityTypeServices),
		string(kong.EntityTypeRoutes),
		string(kong.EntityTypeUpstreams),
		string(kong.EntityTypeTargets),
		string(kong.EntityTypeConsumers),
		string(kong.EntityTypeConsumerGroups),
		string(kong.EntityTypePlugins):
		return true
	default:
		return false
	}
}

// KongCustomEntityCollection is a collection of custom Kong entities with the same type.
type KongCustomEntityCollection struct {
	// Schema is the Schema of the entity.
	Schema EntitySchema `json:"-"`
	// Entities is the list of entities in the collection.
	Entities []CustomEntity
}

// CustomEntity saves content of a Kong custom entity with the pointer to the k8s resource translating to it.
type CustomEntity struct {
	custom.Object
	// K8sKongCustomEntity refers to the KongCustomEntity resource that translate to it.
	K8sKongCustomEntity *kongv1alpha1.KongCustomEntity
	// ForeignEntityIDs stores the IDs of the foreign Kong entities attached to the entity.
	ForeignEntityIDs map[kong.EntityType]string
}

// SchemaGetter is the interface to fetch the schema of a Kong entity by its type.
// Used for fetching schema of custom entity for filling "foreign" field referring to other entities.
type SchemaGetter interface {
	Get(ctx context.Context, entityType string) (kong.Schema, error)
}

// sortCustomEntities sorts the custom entities of each type.
// Since there may not be a consistent field to identify an entity, here we sort them by the k8s namespace/name.
func (ks *KongState) sortCustomEntities() {
	for _, collection := range ks.CustomEntities {
		sort.Slice(collection.Entities, func(i, j int) bool {
			e1 := collection.Entities[i]
			e2 := collection.Entities[j]
			// Compare namespace first.
			if e1.K8sKongCustomEntity.Namespace < e2.K8sKongCustomEntity.Namespace {
				return true
			}
			if e1.K8sKongCustomEntity.Namespace > e2.K8sKongCustomEntity.Namespace {
				return false
			}
			// If namespace are the same, compare names.
			if e1.K8sKongCustomEntity.Name < e2.K8sKongCustomEntity.Name {
				return true
			}
			if e1.K8sKongCustomEntity.Name > e2.K8sKongCustomEntity.Name {
				return false
			}
			// Namespace and name are all the same.
			// This means the two entities are generated from the same KCE resource but attached to different foreign entities.
			// So we need to compare foreign entities.
			if e1.ForeignEntityIDs != nil && e2.ForeignEntityIDs != nil {
				// Compare IDs of attached entities in services, routes, consumers order.
				foreignEntityTypeList := []kong.EntityType{
					kong.EntityTypeServices,
					kong.EntityTypeRoutes,
					kong.EntityTypeConsumers,
				}
				for _, t := range foreignEntityTypeList {
					if e1.ForeignEntityIDs[t] != e2.ForeignEntityIDs[t] {
						return e1.ForeignEntityIDs[t] < e2.ForeignEntityIDs[t]
					}
				}
			}
			// Should not reach here when k8s namespace/names are the same, and foreign entities are also the same.
			// This means we generated two Kong entities from one KCE (and attached to the same foreign entities if any).
			return true
		})
	}
}

type entityForeignFieldValue struct {
	fieldName         string
	foreignEntityType kong.EntityType
	foreignEntityID   string
}
