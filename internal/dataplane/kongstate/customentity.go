package kongstate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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

type entityForeignFieldValue struct {
	fieldName         string
	foreignEntityType kong.EntityType
	foreignEntityID   string
}

// FillCustomEntities fills custom entities in KongState.
func (ks *KongState) FillCustomEntities(
	logger logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	schemaGetter SchemaGetter,
	workspace string,
) {
	entities := s.ListKongCustomEntities()
	if len(entities) == 0 {
		return
	}
	logger = logger.WithName("fillCustomEntities")

	if ks.CustomEntities == nil {
		ks.CustomEntities = map[string]*KongCustomEntityCollection{}
	}
	// Fetch relations between plugins and services/routes/consumers and store the pointer to translated Kong entities.
	// Used for fetching entity referred by a custom entity and fill the ID of referred entity.
	pluginRels := ks.getPluginRelatedEntitiesRef(s, logger)

	for _, entity := range entities {
		// reject the custom entity if its type is in "known" entity types that are already processed.
		if IsKnownEntityType(entity.Spec.EntityType) {
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("cannot use known entity type %s in custom entity", entity.Spec.EntityType),
				entity,
			)
			continue
		}
		// Fetch the entity schema.
		schema, err := ks.fetchEntitySchema(schemaGetter, entity.Spec.EntityType)
		if err != nil {
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("failed to fetch entity schema for entity type %s: %v", entity.Spec.EntityType, err),
				entity,
			)
			continue
		}

		// Fill the "foreign" fields if the entity has such fields referencing services/routes/consumers.
		// First Find out possible foreign field combinations attached to the KCE resource.
		foreignFieldCombinations, err := findCustomEntityForeignFields(logger, s, entity, schema, pluginRels, workspace)
		if err != nil {
			failuresCollector.PushResourceFailure(fmt.Sprintf("failed to find attached foreign entities for custom entity: %v", err), entity)
			continue
		}
		// generate Kong entities from the fields in the KCE itself and attached foreign entities.
		generatedEntities, err := generateCustomEntities(entity, foreignFieldCombinations)
		if err != nil {
			failuresCollector.PushResourceFailure(fmt.Sprintf("failed to generate entities from itself and attach foreign entities: %v", err), entity)
			continue
		}
		for _, generatedEntity := range generatedEntities {
			ks.AddCustomEntity(entity.Spec.EntityType, schema, generatedEntity)
		}
	}

	ks.sortCustomEntities()
}

// AddCustomEntity adds a custom entity into the collection of its type.
func (ks *KongState) AddCustomEntity(entityType string, schema EntitySchema, e CustomEntity) {
	if ks.CustomEntities == nil {
		ks.CustomEntities = map[string]*KongCustomEntityCollection{}
	}
	// Put the entity into the custom collection to store the entities of its type.
	if _, ok := ks.CustomEntities[entityType]; !ok {
		ks.CustomEntities[entityType] = &KongCustomEntityCollection{
			Schema: schema,
		}
	}
	collection := ks.CustomEntities[entityType]
	collection.Entities = append(collection.Entities, e)
}

// CustomEntityTypes returns types of translated custom entities included in the KongState.
func (ks *KongState) CustomEntityTypes() []string {
	return lo.Keys(ks.CustomEntities)
}

// fetchEntitySchema fetches schema of an entity by its type and stores the schema in its custom entity collection
// as a cache to avoid excessive calling of Kong admin APIs.
func (ks *KongState) fetchEntitySchema(schemaGetter SchemaGetter, entityType string) (EntitySchema, error) {
	collection, ok := ks.CustomEntities[entityType]
	if ok {
		return collection.Schema, nil
	}
	// Use `context.Background()` here because `BuildKongConfig` does not provide a context.
	schema, err := schemaGetter.Get(context.Background(), entityType)
	if err != nil {
		return EntitySchema{}, err
	}
	return ExtractEntityFieldDefinitions(schema), nil
}

// sortCustomEntities sorts the custom entities of each type.
// Since there may not be a consistent field to identify an entity, here we sort them by the k8s namespace/name.
func (ks *KongState) sortCustomEntities() {
	for _, collection := range ks.CustomEntities {
		sort.Slice(collection.Entities, func(i, j int) bool {
			e1 := collection.Entities[i]
			e2 := collection.Entities[j]
			// Compare namespace first.
			if e1.K8sKongCustomEntity.Namespace != e2.K8sKongCustomEntity.Namespace {
				return e1.K8sKongCustomEntity.Namespace < e2.K8sKongCustomEntity.Namespace
			}
			// If namespace are the same, compare names.
			if e1.K8sKongCustomEntity.Name != e2.K8sKongCustomEntity.Name {
				return e1.K8sKongCustomEntity.Name < e2.K8sKongCustomEntity.Name
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

func findCustomEntityRelatedPlugin(logger logr.Logger, cacheStore store.Storer, k8sEntity *kongv1alpha1.KongCustomEntity) (string, bool, error) {
	// Find referred entity via the plugin in its spec.parentRef.
	// Then we can fetch the referred service/route/consumer from the reference relations of the plugin.
	parentRef := k8sEntity.Spec.ParentRef
	// Abort if the parentRef is empty or does not refer to a plugin.
	if parentRef == nil ||
		(parentRef.Group == nil || *parentRef.Group != kongv1alpha1.GroupVersion.Group) {
		return "", false, nil
	}
	if parentRef.Kind == nil || (*parentRef.Kind != "KongPlugin" && *parentRef.Kind != "KongClusterPlugin") {
		return "", false, nil
	}

	// Extract the plugin key to get the plugin relations.
	paretRefNamespace := lo.FromPtrOr(parentRef.Namespace, "")
	// if the namespace in parentRef is not same as the namespace of KCE itself, check if the reference is allowed by ReferenceGrant.
	if paretRefNamespace != "" && paretRefNamespace != k8sEntity.Namespace {
		paretRefNamespace, err := extractReferredPluginNamespace(logger, cacheStore, k8sEntity, annotations.NamespacedKongPlugin{
			Namespace: paretRefNamespace,
			Name:      parentRef.Name,
		})
		if err != nil {
			return "", false, err
		}
		return paretRefNamespace + ":" + parentRef.Name, true, nil
	}

	return k8sEntity.Namespace + ":" + parentRef.Name, true, nil
}

func findCustomEntityForeignFields(
	logger logr.Logger,
	cacheStore store.Storer,
	k8sEntity *kongv1alpha1.KongCustomEntity,
	schema EntitySchema,
	pluginRelEntities PluginRelatedEntitiesRefs,
	workspace string,
) ([][]entityForeignFieldValue, error) {
	pluginKey, ok, err := findCustomEntityRelatedPlugin(logger, cacheStore, k8sEntity)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	// Get the relations with other entities of the plugin.
	rels, ok := pluginRelEntities.RelatedEntities[pluginKey]
	if !ok {
		return nil, nil
	}

	var (
		foreignRelations      util.ForeignRelations
		foreignServiceFields  []string
		foreignRouteFields    []string
		foreignConsumerFields []string
	)

	ret := [][]entityForeignFieldValue{}
	for fieldName, field := range schema.Fields {
		if field.Type != EntityFieldTypeForeign {
			continue
		}
		switch field.Reference {
		case string(kong.EntityTypeServices):
			foreignServiceFields = append(foreignServiceFields, fieldName)
			foreignRelations.Service = getServiceIDFromPluginRels(logger, rels, pluginRelEntities.RouteAttachedService, workspace)
		case string(kong.EntityTypeRoutes):
			foreignRouteFields = append(foreignRouteFields, fieldName)
			foreignRelations.Route = lo.FilterMap(rels.Routes, func(r *Route, _ int) (string, bool) {
				if err := r.FillID(workspace); err != nil {
					return "", false
				}
				return *r.ID, true
			})
		case string(kong.EntityTypeConsumers):
			foreignConsumerFields = append(foreignConsumerFields, fieldName)
			foreignRelations.Consumer = lo.FilterMap(rels.Consumers, func(c *Consumer, _ int) (string, bool) {
				if err := c.FillID(workspace); err != nil {
					return "", false
				}
				return *c.ID, true
			})
		} // end of switch
	}

	// TODO: Here we inherited the logic of generating combinations of attached foreign entities for plugins.
	// Actually there are no such case that a custom entity required multiple "foreign" fields in current Kong plugins.
	// So it is still uncertain how to generate foreign field combinations for custom entities.
	for _, combination := range foreignRelations.GetCombinations() {
		foreignFieldValues := []entityForeignFieldValue{}
		for _, fieldName := range foreignServiceFields {
			foreignFieldValues = append(foreignFieldValues, entityForeignFieldValue{
				fieldName:         fieldName,
				foreignEntityType: kong.EntityTypeServices,
				foreignEntityID:   combination.Service,
			})
		}
		for _, fieldName := range foreignRouteFields {
			foreignFieldValues = append(foreignFieldValues, entityForeignFieldValue{
				fieldName:         fieldName,
				foreignEntityType: kong.EntityTypeRoutes,
				foreignEntityID:   combination.Route,
			})
		}
		for _, fieldName := range foreignConsumerFields {
			foreignFieldValues = append(foreignFieldValues, entityForeignFieldValue{
				fieldName:         fieldName,
				foreignEntityType: kong.EntityTypeConsumers,
				foreignEntityID:   combination.Consumer,
			})
		}
		ret = append(ret, foreignFieldValues)
	}

	return ret, nil
}

// generateCustomEntities generates Kong entities from KongCustomEntity resource and combinations of attached foreign entities.
// If the KCE is attached to any foreign entities, it generates one entity per combination of foreign entities.
// If the KCE is not attached, generate one entity for itself.
func generateCustomEntities(
	entity *kongv1alpha1.KongCustomEntity,
	foreignFieldCombinations [][]entityForeignFieldValue,
) ([]CustomEntity, error) {
	copyEntityFields := func() (map[string]any, error) {
		// Unmarshal the fields of the entity to have a fresh copy for each combination as we may modify them.
		fields := map[string]any{}
		if err := json.Unmarshal(entity.Spec.Fields.Raw, &fields); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entity fields: %w", err)
		}
		return fields, nil
	}
	// If there are any foreign fields, generate one entity per each foreign entity combination.
	if len(foreignFieldCombinations) > 0 {
		var customEntities []CustomEntity
		for _, combination := range foreignFieldCombinations {
			entityFields, err := copyEntityFields()
			if err != nil {
				return nil, err
			}
			generatedEntity := CustomEntity{
				K8sKongCustomEntity: entity,
				ForeignEntityIDs:    make(map[kong.EntityType]string),
				Object:              entityFields,
			}
			// Fill the fields referring to foreign entities.
			for _, foreignField := range combination {
				entityFields[foreignField.fieldName] = map[string]any{
					"id": foreignField.foreignEntityID,
				}
				// Save the referred foreign entity IDs for sorting.
				generatedEntity.ForeignEntityIDs[foreignField.foreignEntityType] = foreignField.foreignEntityID
			}
			customEntities = append(customEntities, generatedEntity)
		}
		return customEntities, nil
	}

	// Otherwise (no foreign fields), generate a single entity.
	entityFields, err := copyEntityFields()
	if err != nil {
		return nil, err
	}
	return []CustomEntity{
		{
			K8sKongCustomEntity: entity,
			Object:              entityFields,
		},
	}, nil
}
