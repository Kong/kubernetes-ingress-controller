package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hbagdi/go-kong/kong"
	"github.com/tidwall/gjson"
)

// FIXME
// decK will release this official API soon, use that and remove this code.

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
func (p *PluginSchemaStore) Schema(pluginName string) (map[string]interface{}, error) {
	if pluginName == "" {
		return nil, fmt.Errorf("pluginName can not be empty")
	}

	// lookup in cache
	if schema, ok := p.schemas[pluginName]; ok {
		return schema, nil
	}

	// not present in cache, lookup
	req, err := p.client.NewRequest("GET", "/plugins/schema/"+pluginName,
		nil, nil)
	if err != nil {
		return nil, err
	}
	schema := make(map[string]interface{})
	_, err = p.client.Do(context.TODO(), req, &schema)
	if err != nil {
		return nil, err
	}
	p.schemas[pluginName] = schema
	return schema, nil
}

func fill(schema map[string]interface{},
	config kong.Configuration) (kong.Configuration, error) {
	jsonb, err := json.Marshal(&schema)
	if err != nil {
		return nil, err
	}
	// Get all in the schema
	value := gjson.ParseBytes((jsonb))
	return fillRecord(value, config)
}

func fillRecord(schema gjson.Result, config kong.Configuration) (kong.Configuration, error) {
	if config == nil {
		return nil, nil
	}
	res := config.DeepCopy()
	value := schema.Get("fields")

	value.ForEach(func(key, value gjson.Result) bool {
		// get the key name
		ms := value.Map()
		fname := ""
		for k := range ms {
			fname = k
			break
		}
		ftype := value.Get(fname + ".type")
		if ftype.String() == "record" {
			subConfig := config[fname]
			if subConfig == nil {
				subConfig = make(map[string]interface{})
			}
			newSubConfig, err := fillRecord(value.Get(fname), subConfig.(map[string]interface{}))
			if err != nil {
				panic(err)
			}
			res[fname] = map[string]interface{}(newSubConfig)
			return true
		}
		// check if key is already set in the config
		if _, ok := config[fname]; ok {
			// yes, don't set it
			return true
		}
		// no, set it
		value = value.Get(fname + ".default")
		if value.Exists() {
			res[fname] = value.Value()
		} else {
			// if no default exists, set an explicit nil
			res[fname] = nil
		}
		return true
	})

	return res, nil
}
