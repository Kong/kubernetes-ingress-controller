package sendconfig

import (
	"context"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
)

// UpdateStrategy is the way we approach updating data-plane's configuration, depending on its type.
type UpdateStrategy interface {
	// Update applies targetConfig to the data-plane.
	Update(ctx context.Context, targetContent *file.Content) error

	// MetricsProtocol returns a string describing the update strategy type to be used in metrics.
	MetricsProtocol() string
}

type UpdateClient interface {
	IsKonnect() bool
	KonnectRuntimeGroup() string
	AdminAPIClient() *kong.Client
}

func ResolveUpdateStrategy(
	client UpdateClient,
	config Config,
) UpdateStrategy {
	adminAPIClient := client.AdminAPIClient()

	if !config.InMemory || client.IsKonnect() {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:         config.SkipCACertificates,
				SelectorTags:        config.FilterTags,
				KonnectRuntimeGroup: getKonnectRuntimeGroup(client),
			},
			config.Version,
			config.Concurrency,
		)
	}

	return NewUpdateStrategyInMemory(adminAPIClient)
}

// getKonnectRuntimeGroup returns Konnect's Runtime Group UUID if the client is Konnect-compatible.
// Otherwise, it returns an empty string.
func getKonnectRuntimeGroup(client UpdateClient) string {
	if client.IsKonnect() {
		return client.KonnectRuntimeGroup()
	}
	return ""
}
