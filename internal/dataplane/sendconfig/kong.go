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

	// Currently, this assumes that all underlying clients are using the same version
	// hence this shared field in here.
	Version  semver.Version
	InMemory bool

	Concurrency int
	FilterTags  []string
}

// New creates new Kong client that is responsible for sending configurations
// to Kong instance(s) through Admin API.
func New(
	ctx context.Context,
	logger logr.Logger,
	kongClients []*adminapi.Client,
	v semver.Version,
	dbMode string,
	concurrency int,
	filterTags []string,
) Kong {
	var (
		errg errgroup.Group
		tags []string
	)

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
		logger.Error(err, "tag filtering disabled")
	} else {
		logger.Info("tag filtering enabled", "tags", filterTags)
		tags = filterTags
	}

	return Kong{
		InMemory:    (dbMode == "off") || (dbMode == ""),
		Version:     v,
		FilterTags:  tags,
		Concurrency: concurrency,
		Clients: lo.Map(kongClients, func(client *adminapi.Client, index int) ClientWithPluginStore {
			return ClientWithPluginStore{
				Client:            client,
				PluginSchemaStore: util.NewPluginSchemaStore(client.Client),
			}
		}),
	}
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
