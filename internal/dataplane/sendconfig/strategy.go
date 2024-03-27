package sendconfig

import (
	"context"

	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// ContentWithHash encapsulates file.Content along with its precalculated hash.
type ContentWithHash struct {
	Content *file.Content
	Hash    []byte
}

// UpdateStrategy is the way we approach updating data-plane's configuration, depending on its type.
type UpdateStrategy interface {
	// Update applies targetConfig to the data-plane.
	Update(ctx context.Context, targetContent ContentWithHash) (
		err error,
		resourceErrors []ResourceError,
		resourceErrorsParseErr error,
	)

	// MetricsProtocol returns a string describing the update strategy type to be used in metrics.
	MetricsProtocol() metrics.Protocol

	// Type returns a human-readable debug string representing the UpdateStrategy.
	Type() string
}

type UpdateClient interface {
	IsKonnect() bool
	KonnectRuntimeGroup() string
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
	log    logrus.FieldLogger
}

func NewDefaultUpdateStrategyResolver(config Config, log logrus.FieldLogger) DefaultUpdateStrategyResolver {
	return DefaultUpdateStrategyResolver{
		config: config,
		log:    log,
	}
}

// ResolveUpdateStrategy returns an UpdateStrategy based on the client and configuration.
// The UpdateStrategy can be either UpdateStrategyDBMode or UpdateStrategyInMemory. Both
// of them implement different ways to populate Kong instances with data-plane configuration.
// If the client implements UpdateClientWithBackoff interface, its strategy will be decorated
// with the backoff strategy it provides.
func (r DefaultUpdateStrategyResolver) ResolveUpdateStrategy(
	client UpdateClient,
) UpdateStrategy {
	updateStrategy := r.resolveUpdateStrategy(client)

	if clientWithBackoff, ok := client.(UpdateClientWithBackoff); ok {
		return NewUpdateStrategyWithBackoff(updateStrategy, clientWithBackoff.BackoffStrategy(), r.log)
	}

	return updateStrategy
}

func (r DefaultUpdateStrategyResolver) resolveUpdateStrategy(client UpdateClient) UpdateStrategy {
	adminAPIClient := client.AdminAPIClient()

	// In case the client communicates with Konnect Admin API, we know it has to use DB-mode. There's no need to check
	// config.InMemory that is meant for regular Kong Gateway clients.
	if client.IsKonnect() {
		return NewUpdateStrategyDBModeKonnect(
			adminAPIClient,
			dump.Config{
				SkipCACerts: true,
			},
			r.config.Version,
			r.config.Concurrency,
		)
	}

	if !r.config.InMemory {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:  r.config.SkipCACertificates,
				SelectorTags: r.config.FilterTags,
			},
			r.config.Version,
			r.config.Concurrency,
		)
	}

	return NewUpdateStrategyInMemory(
		adminAPIClient,
		DefaultContentToDBLessConfigConverter{},
		r.log,
		r.config.Version,
	)
}
