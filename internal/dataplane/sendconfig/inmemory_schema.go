package sendconfig

import (
	"github.com/kong/deck/file"
)

// DBLessConfig is the configuration that is sent to Kong's data-plane via its `POST /config` endpoint after being
// marshalled to JSON.
// It uses file.Content as its base schema, but it also includes additional fields that are not part of decK's schema.
type DBLessConfig struct {
	file.Content
	ConsumerGroupConsumerRelationships []ConsumerGroupConsumerRelationship `json:"consumer_group_consumers,omitempty"`
	ConsumerGroupPluginRelationships   []ConsumerGroupPluginRelationship   `json:"consumer_group_plugins,omitempty"`
}

// ConsumerGroupConsumerRelationship is a relationship between a ConsumerGroup and a Consumer.
type ConsumerGroupConsumerRelationship struct {
	ConsumerGroup string `json:"consumer_group"`
	Consumer      string `json:"consumer"`
}

// ConsumerGroupPluginRelationship is a relationship between a ConsumerGroup and a Plugin.
type ConsumerGroupPluginRelationship struct {
	ConsumerGroup string `json:"consumer_group"`
	Plugin        string `json:"plugin"`
}

type DefaultContentToDBLessConfigConverter struct{}

func (DefaultContentToDBLessConfigConverter) Convert(content *file.Content) DBLessConfig {
	dblessConfig := DBLessConfig{
		Content: *content,
	}

	// DBLess schema does not support decK's Info section.
	dblessConfig.Content.Info = nil

	// DBLess schema does not support nulls in plugin configs.
	cleanUpNullsInPluginConfigs(&dblessConfig.Content)

	// DBLess schema does not 1-1 match decK's schema for ConsumerGroups.
	convertConsumerGroups(&dblessConfig)

	return dblessConfig
}

// cleanUpNullsInPluginConfigs removes null values from plugins' configs.
func cleanUpNullsInPluginConfigs(state *file.Content) {
	for _, s := range state.Services {
		for _, p := range s.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
		for _, r := range state.Routes {
			for _, p := range r.Plugins {
				for k, v := range p.Config {
					if v == nil {
						delete(p.Config, k)
					}
				}
			}
		}
	}

	for _, c := range state.Consumers {
		for _, p := range c.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
	}

	for _, p := range state.Plugins {
		for k, v := range p.Config {
			if v == nil {
				delete(p.Config, k)
			}
		}
	}
}

// convertConsumerGroups drops consumer groups related fields that are not supported in DBLess schema:
//   - Content.Plugins[].ConsumerGroup
//   - Content.Consumers[].Group,
//   - Content.ConsumerGroups[].Plugins
//   - Content.ConsumerGroups[].Consumers
//
// In their place it creates relationships slices:
//   - ConsumerGroupConsumerRelationships
//   - ConsumerGroupPluginRelationships
func convertConsumerGroups(dblessConfig *DBLessConfig) {
	// Before any modification to the Plugins slice, we need to make a copy of it to not modify the original slice.
	if len(dblessConfig.Content.Consumers) > 0 {
		copiedConsumers := make([]file.FConsumer, len(dblessConfig.Content.Consumers))
		copy(copiedConsumers, dblessConfig.Content.Consumers)
		dblessConfig.Content.Consumers = copiedConsumers
	}

	// DBLess schema does not support Consumer.Groups field...
	for i, c := range dblessConfig.Content.Consumers {
		// ... therefore we need to convert them to relationships...
		for _, cg := range dblessConfig.Content.Consumers[i].Groups {
			dblessConfig.ConsumerGroupConsumerRelationships = append(dblessConfig.ConsumerGroupConsumerRelationships, ConsumerGroupConsumerRelationship{
				ConsumerGroup: *cg.Name,
				Consumer:      *c.Username,
			})
		}
		// ... and remove them from the Consumer struct.
		dblessConfig.Content.Consumers[i].Groups = nil
	}

	// Before any modification to the Plugins slice, we need to make a copy of it to not modify the original slice.
	if len(dblessConfig.Content.Plugins) > 0 {
		copiedPlugins := make([]file.FPlugin, len(dblessConfig.Content.Plugins))
		copy(copiedPlugins, dblessConfig.Content.Plugins)
		dblessConfig.Content.Plugins = copiedPlugins
	}

	// DBLess schema does not support Consumer.ConsumerGroup field...
	for i, p := range dblessConfig.Content.Plugins {
		// ... therefore we need to convert it to relationships...
		if p.ConsumerGroup != nil {
			dblessConfig.ConsumerGroupPluginRelationships = append(dblessConfig.ConsumerGroupPluginRelationships, ConsumerGroupPluginRelationship{
				ConsumerGroup: *p.ConsumerGroup.Name,
				Plugin:        *p.Name,
			})
		}
		// ... and remove it from the Plugin struct.
		dblessConfig.Content.Plugins[i].ConsumerGroup = nil
	}

	// DBLess schema does not support ConsumerGroups.Consumers and ConsumerGroups.Plugins fields so we need to remove
	// them. We also need to make a copy of the ConsumerGroups slice to not modify the original slice.
	if len(dblessConfig.Content.ConsumerGroups) > 0 {
		copiedConsumerGroups := make([]file.FConsumerGroupObject, len(dblessConfig.Content.ConsumerGroups))
		copy(copiedConsumerGroups, dblessConfig.Content.ConsumerGroups)
		dblessConfig.Content.ConsumerGroups = copiedConsumerGroups
	}
	for i := range dblessConfig.Content.ConsumerGroups {
		dblessConfig.Content.ConsumerGroups[i].Consumers = nil
		dblessConfig.Content.ConsumerGroups[i].Plugins = nil
	}
}
