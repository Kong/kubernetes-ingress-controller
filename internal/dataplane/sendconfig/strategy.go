package sendconfig

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"

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
	// Update applies targetConfig to the data-plane.
	Update(ctx context.Context, targetContent ContentWithHash) error

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
			// The DB mode update strategy is used for both DB mode gateways and Konnect-integrated controllers. In the
			// Konnect case, we don't actually want to collect diffs, and don't actually provide a diagnostic when setting
			// it up, so we only collect and send diffs if we're talking to a gateway.
			//
			// TODO maybe this is wrong? I'm not sure if we actually support (or if not, explicitly prohibit)
			// configuring a controller to use both DB mode and talk to Konnect, or if we only support DB-less when using
			// Konnect. If those are mutually exclusive, maybe we can just collect diffs for Konnect mode? If they're
			// not mutually exclusive, trying to do diagnostics diff updates for both the updates would have both attempt
			// to store diffs. This is... maybe okay. They should be identical, but that's a load-bearing "should": we know
			// Konnect can sometimes differ in what it accepts versus the gateway, and we have some Konnect configuration
			// (consumer exclude, sensitive value mask) where they're _definitely_ different. That same configuration could
			// make the diff confusing even if it's DB mode only, since it doesn't reflect what we're sending to the gateway
			// in some cases.
			nil,
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
			diagnostic,
			r.logger,
		)
	}

	return NewUpdateStrategyInMemory(
		adminAPIClient,
		DefaultContentToDBLessConfigConverter{},
		r.logger,
	)
}
