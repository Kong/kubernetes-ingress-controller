package sendconfig

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

// CustomEntitiesByType stores all custom entities by types.
// The key is the type of the entity,
// and the corresponding slice stores the sorted list of custom entities with that type.
type CustomEntitiesByType map[string][]custom.Object

// ContentWithHash encapsulates file.Content along with its precalculated hash.
type ContentWithHash struct {
	Content        *file.Content
	CustomEntities CustomEntitiesByType
	Hash           []byte
}

// UpdateStrategy is the way we approach updating data-plane's configuration, depending on its type.
type UpdateStrategy interface {
	// Update applies targetConfig to the DataPlane. When the update is successful, it returns the number of
	// bytes sent to the DataPlane or mo.None when it's impossible to determine the number of bytes sent e.g.
	// for dbmode (deck) strategy.
	Update(ctx context.Context, targetContent ContentWithHash) (mo.Option[int], error)

	// MetricsProtocol returns a string describing the update strategy type to be used in metrics.
	MetricsProtocol() metrics.Protocol

	// Type returns a human-readable debug string representing the UpdateStrategy.
	Type() string
}

type UpdateClient interface {
	IsKonnect() bool
	KonnectControlPlane() string
	AdminAPIClient() *kong.Client
}

type UpdateClientWithBackoff interface {
	UpdateClient
	BackoffStrategy() adminapi.UpdateBackoffStrategy
}

// ResourceError is a Kong configuration error associated with a Kubernetes resource.
type ResourceError struct {
	Name       string
	Namespace  string
	Kind       string
	APIVersion string
	UID        string
	Problems   map[string]string
}

type DefaultUpdateStrategyResolver struct {
	config Config
	logger logr.Logger
}

func NewDefaultUpdateStrategyResolver(config Config, logger logr.Logger) DefaultUpdateStrategyResolver {
	return DefaultUpdateStrategyResolver{
		config: config,
		logger: logger,
	}
}

// ResolveUpdateStrategy returns an UpdateStrategy based on the client and configuration.
// The UpdateStrategy can be either UpdateStrategyDBMode or UpdateStrategyInMemory. Both
// of them implement different ways to populate Kong instances with data-plane configuration.
// If the client implements UpdateClientWithBackoff interface, its strategy will be decorated
// with the backoff strategy it provides.
func (r DefaultUpdateStrategyResolver) ResolveUpdateStrategy(
	client UpdateClient,
	diagnostic *diagnostics.ClientDiagnostic,
) UpdateStrategy {
	updateStrategy := r.resolveUpdateStrategy(client, diagnostic)

	if clientWithBackoff, ok := client.(UpdateClientWithBackoff); ok {
		return NewUpdateStrategyWithBackoff(updateStrategy, clientWithBackoff.BackoffStrategy(), r.logger)
	}

	return updateStrategy
}

func (r DefaultUpdateStrategyResolver) resolveUpdateStrategy(
	client UpdateClient,
	diagnostic *diagnostics.ClientDiagnostic,
) UpdateStrategy {
	adminAPIClient := client.AdminAPIClient()

	// In case the client communicates with Konnect Admin API, we know it has to use DB-mode. There's no need to check
	// config.InMemory that is meant for regular Kong Gateway clients.
	if client.IsKonnect() {
		return NewUpdateStrategyDBModeKonnect(
			adminAPIClient,
			dump.Config{
				KonnectControlPlane: client.KonnectControlPlane(),
			},
			r.config.Version,
			r.config.Concurrency,
			r.logger,
		)
	}

	if !r.config.InMemory {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:     r.config.SkipCACertificates,
				SelectorTags:    r.config.FilterTags,
				IncludeLicenses: true,
			},
			r.config.Version,
			r.config.Concurrency,
			r.logger,
			WithDiagnostic(diagnostic),
		)
	}

	return NewUpdateStrategyInMemory(
		adminAPIClient,
		DefaultContentToDBLessConfigConverter{},
		r.logger,
	)
}
