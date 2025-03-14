package kongstate

import (
	"github.com/kong/go-kong/kong"

	configurationv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
)

// ConsumerGroup holds a Kong Consumer.
type ConsumerGroup struct {
	kong.ConsumerGroup

	K8sKongConsumerGroup configurationv1beta1.KongConsumerGroup
}
