package sendconfig

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Types
// -----------------------------------------------------------------------------

// Kong Represents a Kong client and connection information.
type Kong struct {
	Clients []ClientWithPluginStore
	Config  Config
}

// Config gathers parameters that are needed for sending configuration to Kong Admin API.
type Config struct {
	// Currently, this assumes that all underlying clients are using the same version
	// hence this shared field in here.
	Version semver.Version

	// InMemory
	InMemory bool

	// Concurrency defines how many concurrent goroutines should be used when syncing configuration in DB-mode.
	Concurrency int

	// FilterTags are tags used to manage and filter entities in Kong.
	FilterTags []string

	// SkipCACertificates disables CA certificates, to avoid fighting over configuration in multi-workspace
	// environments. See https://github.com/Kong/deck/pull/617
	SkipCACertificates bool

	// EnableReverseSync indicates that reverse sync should be enabled for
	// updates to the data-plane.
	EnableReverseSync bool
}

// New creates new Kong client that is responsible for sending configurations
// to Kong instance(s) through Admin API.
func New(
	ctx context.Context,
	logger logr.Logger,
	kongClients []*adminapi.Client,
	cfg Config,
) Kong {
	if err := tagsFilteringEnabled(ctx, kongClients); err != nil {
		logger.Error(err, "tag filtering disabled")
		cfg.FilterTags = nil
	} else {
		logger.Info("tag filtering enabled", "tags", cfg.FilterTags)
	}

	return Kong{
		Config: cfg,
		Clients: lo.Map(kongClients, func(client *adminapi.Client, index int) ClientWithPluginStore {
			return ClientWithPluginStore{
				Client:            client,
				PluginSchemaStore: util.NewPluginSchemaStore(client.Client),
			}
		}),
	}
}

func tagsFilteringEnabled(ctx context.Context, kongClients []*adminapi.Client) error {
	var errg errgroup.Group
	for _, cl := range kongClients {
		cl := cl
		errg.Go(func() error {
			ok, err := cl.Tags.Exists(ctx)
			if err != nil {
				return fmt.Errorf("Kong Admin API (%s) does not support tags: %w", cl.BaseRootURL(), err)
			}
			if !ok {
				return fmt.Errorf("Kong Admin API (%s) does not support tags", cl.BaseRootURL())
			}
			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		return err
	}

	return nil
}

type ClientWithPluginStore struct {
	*adminapi.Client
	*util.PluginSchemaStore
	// lastConfigSHA is a checksum of the last successful update to the data-plane
	lastConfigSHA []byte
}

func (c *ClientWithPluginStore) SetLastConfigSHA(s []byte) {
	c.lastConfigSHA = s
}

func (c *ClientWithPluginStore) LastConfigSHA() []byte {
	return c.lastConfigSHA
}
