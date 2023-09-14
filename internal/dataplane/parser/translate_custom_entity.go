package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
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
			p.logger.Errorf("failed to compose Kong custom entity from resource %s/%s, error %v", e.Namespace, e.Name, err)
			continue
		}
		p.logger.Debugf("kong entity type %s, object %+v", def.Spec.Name, kongEntity.Object())
	}
}

func (p *Parser) renderCustomEntity(e *kongv1alpha1.KongCustomEntity, defSpec kongv1alpha1.KongCustomEntityDefinitionSpec) (custom.Entity, error) {
	kongEntity := custom.NewEntityObject(custom.Type(defSpec.Name))
	o := custom.Object{}

	var rawPatches []string
	for _, patch := range e.Spec.Patches {
		if patch.ConfigSource != nil {
			secretName := patch.ConfigSource.SecretValue.Secret
			secretKey := patch.ConfigSource.SecretValue.Key
			secert, err := p.storer.GetSecret(e.Namespace, secretName)
			if err != nil {
				return nil, err
			}
			data, ok := secert.Data[secretKey]
			if !ok {
				return nil, fmt.Errorf("key %s not found in secret %s", secretKey, secretName)
			}

			// REVIEW: JSON patch needs the value to be QUOTED (like "key":"\"value\"" format).
			// Should we define the config source as key + type + secret reference format?
			rawPatches = append(rawPatches, fmt.Sprintf(`{"op": "add", "path": "%s", "value": %s}`, patch.Path, data))
		}
	}
	patches := fmt.Sprintf("[%s]", strings.Join(rawPatches, ","))

	patch, err := jsonpatch.DecodePatch([]byte(patches))
	if err != nil {
		return nil, fmt.Errorf("failed to decode patches from %s: %w", patches, err)
	}

	patchedConfig, err := patch.Apply(e.Spec.Fields.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patches: %w", err)
	}
	err = json.Unmarshal(patchedConfig, &o)
	if err != nil {
		p.logger.Debug("patched config:", string(patchedConfig))
		return nil, fmt.Errorf("failed to unmarshal entity: %w", err)
	}

	kongEntity.SetObject(o)
	return kongEntity, nil
}
