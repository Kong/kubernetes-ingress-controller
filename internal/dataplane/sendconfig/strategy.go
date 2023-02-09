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

// ResolveUpdateStrategy returns an UpdateStrategy based on the client and configuration.
// The UpdateStrategy can be either UpdateStrategyDBMode or UpdateStrategyInMemory. Both
// of them implement different ways to populate Kong instances with data-plane configuration.
func ResolveUpdateStrategy(
	client UpdateClient,
	config Config,
) UpdateStrategy {
	adminAPIClient := client.AdminAPIClient()

	// In case the client communicates with Konnect Admin API, we know it has to use DB-mode. There's no need to check
	// config.InMemory that is meant for regular Kong Gateway clients.
	if client.IsKonnect() {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:         config.SkipCACertificates,
				SelectorTags:        config.FilterTags,
				KonnectRuntimeGroup: client.KonnectRuntimeGroup(),
			},
			config.Version,
			config.Concurrency,
		)
	}

	if !config.InMemory {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:  config.SkipCACertificates,
				SelectorTags: config.FilterTags,
			},
			config.Version,
			config.Concurrency,
		)
	}

	return NewUpdateStrategyInMemory(adminAPIClient)
}
