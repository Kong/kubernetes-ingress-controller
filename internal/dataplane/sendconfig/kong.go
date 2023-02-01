package sendconfig

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Types
// -----------------------------------------------------------------------------

// Kong Represents a Kong client and connection information.
type Kong struct {
	Clients []*adminapi.Client
	Config  Config
}

type Config struct {
	// Currently, this assumes that all underlying clients are using the same version
	// hence this shared field in here.
	Version  semver.Version
	InMemory bool

	Concurrency int
	FilterTags  []string

	// SkipCACertificates disables CA certificates, to avoid fighting over configuration in multi-workspace
	// environments. See https://github.com/Kong/deck/pull/617
	SkipCACertificates bool

	// DBMode indicates the current database mode of the backend Kong Admin API
	DBMode string

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
	v semver.Version,
	dbMode string,
	concurrency int,
	filterTags []string,
	skipCACertificates bool,
	enableReverseSync bool,
) Kong {
	var (
		errg errgroup.Group
		tags []string
	)

	for _, cl := range kongClients {
		cl := cl
		errg.Go(func() error {
			ok, err := cl.AdminAPIClient().Tags.Exists(ctx)
			if err != nil {
				return fmt.Errorf("Kong Admin API (%s) does not support tags: %w", cl.AdminAPIClient().BaseRootURL(), err)
			}
			if !ok {
				return fmt.Errorf("Kong Admin API (%s) does not support tags", cl.AdminAPIClient().BaseRootURL())
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
		Config: Config{
			Version:            v,
			InMemory:           (dbMode == "off") || (dbMode == ""),
			Concurrency:        concurrency,
			FilterTags:         tags,
			SkipCACertificates: skipCACertificates,
			DBMode:             dbMode,
			EnableReverseSync:  enableReverseSync,
		},
		Clients: kongClients,
	}
}
