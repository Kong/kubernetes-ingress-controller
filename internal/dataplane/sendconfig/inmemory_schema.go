package sendconfig

import (
	"github.com/kong/go-database-reconciler/pkg/file"
)

// DBLessConfig is the configuration that is sent to Kong's data-plane via its `POST /config` endpoint after being
// marshalled to JSON.
// It uses file.Content as its base schema, but it also includes additional fields that are not part of decK's schema.
type DBLessConfig struct {
	file.Content
	ConsumerGroupConsumerRelationships []ConsumerGroupConsumerRelationship `json:"consumer_group_consumers,omitempty"`
}

// ConsumerGroupConsumerRelationship is a relationship between a ConsumerGroup and a Consumer.
type ConsumerGroupConsumerRelationship struct {
	ConsumerGroup string `json:"consumer_group"`
	Consumer      string `json:"consumer"`
}

type DefaultContentToDBLessConfigConverter struct{}

func (DefaultContentToDBLessConfigConverter) Convert(content *file.Content) DBLessConfig {
	dblessConfig := DBLessConfig{
		Content: *content,
	}

	// DBLess schema does not support decK's Info section.
	dblessConfig.Info = nil

	// DBLess schema does not 1-1 match decK's schema for ConsumerGroups.
	convertConsumerGroups(&dblessConfig)

	return dblessConfig
}

// convertConsumerGroups drops consumer groups related fields that are not supported in DBLess schema:
//   - Content.Consumers[].Groups,
//   - Content.ConsumerGroups[].Plugins
//   - Content.ConsumerGroups[].Consumers
//
// In their place it creates relationships slices:
//   - ConsumerGroupConsumerRelationships
func convertConsumerGroups(dblessConfig *DBLessConfig) {
	// DBLess schema does not support Consumer.Groups field...
	for i, c := range dblessConfig.Consumers {
		// ... therefore we need to convert them to relationships...
		for _, cg := range dblessConfig.Content.Consumers[i].Groups {
			dblessConfig.ConsumerGroupConsumerRelationships = append(dblessConfig.ConsumerGroupConsumerRelationships, ConsumerGroupConsumerRelationship{
				// ... by using FriendlyName() that ensures returning ID if Name is nil...
				ConsumerGroup: cg.FriendlyName(),
				Consumer:      c.FriendlyName(),
			})
		}
		// ... and remove them from the Consumer struct.
		dblessConfig.Content.Consumers[i].Groups = nil
	}
	// DBLess schema does not support ConsumerGroups.Consumers and ConsumerGroups.Plugins fields so we need to remove
	// them.
	for i := range dblessConfig.ConsumerGroups {
		dblessConfig.Content.ConsumerGroups[i].Consumers = nil
		dblessConfig.Content.ConsumerGroups[i].Plugins = nil
	}
}
