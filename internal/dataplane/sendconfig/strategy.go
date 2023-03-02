package sendconfig

import (
	"context"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// UpdateStrategy is the way we approach updating data-plane's configuration, depending on its type.
type UpdateStrategy interface {
	// Update applies targetConfig to the data-plane.
	Update(ctx context.Context, targetContent *file.Content) (
		err error,
		resourceErrors []ResourceError,
		resourceErrorsParseErr error,
	)

	// MetricsProtocol returns a string describing the update strategy type to be used in metrics.
	MetricsProtocol() metrics.Protocol
}

type UpdateClient interface {
	IsKonnect() bool
	KonnectRuntimeGroup() string
	AdminAPIClient() *kong.Client
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
func (r DefaultUpdateStrategyResolver) ResolveUpdateStrategy(
	client UpdateClient,
) UpdateStrategy {
	adminAPIClient := client.AdminAPIClient()

	// In case the client communicates with Konnect Admin API, we know it has to use DB-mode. There's no need to check
	// config.InMemory that is meant for regular Kong Gateway clients.
	if client.IsKonnect() {
		return NewUpdateStrategyDBMode(
			adminAPIClient,
			dump.Config{
				SkipCACerts:         true,
				KonnectRuntimeGroup: client.KonnectRuntimeGroup(),
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

	return NewUpdateStrategyInMemory(adminAPIClient, r.log)
}
