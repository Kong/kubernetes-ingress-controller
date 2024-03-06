package sendconfig

import (
	"context"
	"errors"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

var (
	ErrAdminAPIUnreachable  = errors.New("Kong Admin API is unreachable")
	ErrAdminAPITagsDisabled = errors.New("Kong Admin API does not support tags")
)

// Config gathers parameters that are needed for sending configuration to Kong Admin APIs.
type Config struct {
	// Currently, this assumes that all underlying clients are using the same version
	// hence this shared field in here.
	Version semver.Version

	// InMemory tells whether a Kong Gateway Admin APIs should be communicated in DB-less mode.
	// It's not relevant for Konnect client.
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

	// ExpressionRoutes indicates whether to use Kong's expression routes.
	ExpressionRoutes bool

	// SanitizeKonnectConfigDumps indicates whether to sanitize Konnect config dumps.
	SanitizeKonnectConfigDumps bool
}

// Init sets up variables that need external calls.
func (c *Config) Init(
	ctx context.Context,
	logger logr.Logger,
	kongClients []*adminapi.Client,
) error {
	if err := tagsFilteringEnabled(ctx, kongClients); err != nil {
		if errors.Is(err, ErrAdminAPITagsDisabled) {
			logger.Error(err, "Tag filtering disabled")
			c.FilterTags = nil
			return nil
		}
		return err
	}

	logger.Info("Tag filtering enabled", "tags", c.FilterTags)
	return nil
}

func tagsFilteringEnabled(ctx context.Context, kongClients []*adminapi.Client) error {
	var errg errgroup.Group
	for _, cl := range kongClients {
		cl := cl
		errg.Go(func() error {
			ok, err := cl.AdminAPIClient().Tags.Exists(ctx)

			if !ok {
				if err == nil {
					return ErrAdminAPITagsDisabled
				}
				return fmt.Errorf("%w: %w", ErrAdminAPIUnreachable, err)
			}
			return nil
		})
	}
	return errg.Wait()
}
