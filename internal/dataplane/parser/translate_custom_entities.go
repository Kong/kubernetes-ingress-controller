package parser

import (
	"encoding/json"
	"fmt"

	"github.com/kong/go-kong/kong/custom"
	"github.com/samber/lo"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
)

func (p *Parser) renderCustomEntities() {
	customEntities := p.storer.ListKongCustomEntities()
	customEntityDefinitions := p.storer.ListKongCustomEntityDefinitions()
	entityTypeMap := lo.SliceToMap(customEntityDefinitions, func(d *kongv1alpha1.KongCustomEntityDefinition) (
		string, *kongv1alpha1.KongCustomEntityDefinition,
	) {
		return d.Name, d
	})

	for _, e := range customEntities {
		entityType := e.Spec.Type
		def, ok := entityTypeMap[entityType]
		if !ok {
			p.logger.Errorf("unknown type of custom entity: %s", entityType)
			continue
		}
		p.logger.Debugf("Kong custom entity type %s, from resource %s/%s", def.Spec.Name, e.Namespace, e.Name)
		kongEntity, err := p.renderCustomEntity(e, def.Spec)
		if err != nil {
			p.logger.Errorf("failed to compose Kong custom entity from resource %s/%s", e.Namespace, e.Name)
		}
		p.logger.Debugf("kong entity type %s, object %+v", def.Spec.Name, kongEntity.Object())
		_ = kongEntity
	}
}

func (p *Parser) renderCustomEntity(e *kongv1alpha1.KongCustomEntity, defSpec kongv1alpha1.KongCustomEntityDefinitionSpec) (custom.Entity, error) {
	kongEntity := custom.NewEntityObject(custom.Type(defSpec.Name))
	o := custom.Object{}
	for _, field := range e.Spec.Fields {
		switch field.Type {
		// allow explicit null values.
		case kongv1alpha1.KongEntityFieldTypeNil:
			o[field.Key] = nil
		case kongv1alpha1.KongEntityFieldTypeBoolean:
			var b bool
			// TODO: visit secrets if ValuFrom is not nil. Ditto the following.
			err := json.Unmarshal(field.Value.Raw, &b)
			if err != nil {
				return nil, err
			}
			o[field.Key] = b
		case kongv1alpha1.KongEntityFieldTypeInteger:
			var i int
			err := json.Unmarshal(field.Value.Raw, &i)
			if err != nil {
				return nil, err
			}
			o[field.Key] = i
		case kongv1alpha1.KongEntityFieldTypeNumber:
			var f float64
			err := json.Unmarshal(field.Value.Raw, &f)
			if err != nil {
				return nil, err
			}
			o[field.Key] = f
		case kongv1alpha1.KongEntityFieldTypeString:
			var s string
			err := json.Unmarshal(field.Value.Raw, &s)
			if err != nil {
				return nil, err
			}
			o[field.Key] = s
		case kongv1alpha1.KongEntityFieldTypeArray:
			var a []any
			err := json.Unmarshal(field.Value.Raw, &a)
			if err != nil {
				return nil, err
			}
			o[field.Key] = a
		case kongv1alpha1.KongEntityFieldTypeObject:
			var m map[string]any
			err := json.Unmarshal(field.Value.Raw, &m)
			if err != nil {
				return nil, err
			}
			o[field.Key] = m
		default:
			return nil, fmt.Errorf("unknown field type %s in field %s", field.Type, field.Key)
		}
	}
	kongEntity.SetObject(o)
	return kongEntity, nil
}
