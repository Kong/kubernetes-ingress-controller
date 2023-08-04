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

func (d DefaultContentToDBLessConfigConverter) Convert(content *file.Content) DBLessConfig {
	dblessConfig := DBLessConfig{
		Content: *content,
	}

	// DBLess schema does not support decK's Info section.
	dblessConfig.Content.Info = nil

	// DBLess schema does not support nulls in plugin configs.
	cleanUpNullsInPluginConfigs(&dblessConfig.Content)

	// DBLess schema does not support ConsumerGroup.Consumer and ConsumerGroup.Plugins ...
	for i, cg := range dblessConfig.Content.ConsumerGroups {
		// ... instead, it uses ConsumerGroupConsumerRelationships and ConsumerGroupPluginRelationships.
		for _, consumer := range cg.Consumers {
			dblessConfig.ConsumerGroupConsumerRelationships = append(dblessConfig.ConsumerGroupConsumerRelationships, ConsumerGroupConsumerRelationship{
				ConsumerGroup: *cg.Name,
				Consumer:      *consumer.Username,
			})
		}
		for _, plugin := range cg.Plugins {
			dblessConfig.ConsumerGroupPluginRelationships = append(dblessConfig.ConsumerGroupPluginRelationships, ConsumerGroupPluginRelationship{
				ConsumerGroup: *cg.Name,
				Plugin:        *plugin.Name,
			})
		}

		// ... so we need to remove them from the ConsumerGroup struct.
		content.ConsumerGroups[i].Consumers = nil
		content.ConsumerGroups[i].Plugins = nil
	}

	return dblessConfig
}

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
