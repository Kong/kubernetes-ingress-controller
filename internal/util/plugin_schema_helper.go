package util

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
)

// PluginSchemaStore retrives a schema of a Plugin from Kong.
type PluginSchemaStore struct {
	client  *kong.Client
	schemas map[string]map[string]interface{}
}

// NewPluginSchemaStore creates a PluginSchemaStore.
func NewPluginSchemaStore(client *kong.Client) *PluginSchemaStore {
	return &PluginSchemaStore{
		client:  client,
		schemas: make(map[string]map[string]interface{}),
	}
}

// Schema retrives schema of a plugin.
// A cache is used to save the responses and subsequent queries are served from
// the cache.
func (p *PluginSchemaStore) Schema(ctx context.Context, pluginName string) (map[string]interface{}, error) {
	if pluginName == "" {
		return nil, fmt.Errorf("pluginName can not be empty")
	}

	// lookup in cache
	if schema, ok := p.schemas[pluginName]; ok {
		return schema, nil
	}

	// not present in cache, lookup
	schema, err := p.client.Plugins.GetFullSchema(ctx, &pluginName)
	if err != nil {
		return nil, err
	}
	p.schemas[pluginName] = schema
	return schema, nil
}
