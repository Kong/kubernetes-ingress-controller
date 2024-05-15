package kongstate

import "github.com/kong/go-kong/kong"

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
