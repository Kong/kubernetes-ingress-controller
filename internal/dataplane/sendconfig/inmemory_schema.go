package sendconfig

import (
	"github.com/kong/deck/file"
)

type ConsumerGroupConsumerRelationship struct {
	ConsumerGroup string `json:"consumer_group"`
	Consumer      string `json:"consumer"`
}

type ConsumerGroupPluginRelationship struct {
	ConsumerGroup string `json:"consumer_group"`
	Plugin        string `json:"plugin"`
}

type DBLessConfig struct {
	file.Content                       `json:",inline"`
	ConsumerGroupConsumerRelationships []ConsumerGroupConsumerRelationship `json:"consumer_group_consumers,omitempty"`
	ConsumerGroupPluginRelationships   []ConsumerGroupPluginRelationship   `json:"consumer_group_plugins,omitempty"`
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
	cleanupConsumerGroups(&dblessConfig)

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

// cleanupConsumerGroups drops consumer groups related fields that are not supported in DBLess schema:
//   - Content.Plugins[].ConsumerGroup
//   - Content.Consumers[].Group,
//   - Content.ConsumerGroups[].Plugins
//   - Content.ConsumerGroups[].Consumers
//
// In place of them, it creates relationships slices:
//   - ConsumerGroupConsumerRelationships
//   - ConsumerGroupPluginRelationships
func cleanupConsumerGroups(dblessConfig *DBLessConfig) {
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
	// them.
	for i := range dblessConfig.Content.ConsumerGroups {
		dblessConfig.Content.ConsumerGroups[i].Consumers = nil
		dblessConfig.Content.ConsumerGroups[i].Plugins = nil
	}
}
