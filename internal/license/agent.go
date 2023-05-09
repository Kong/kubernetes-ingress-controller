package license

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// PollingInterval is the interval at which the license agent will poll for license updates.
	PollingInterval = time.Hour * 12

	// PollingTimeout is the timeout for retrieving a license from upstream.
	PollingTimeout = time.Minute * 5
)

type UpstreamClient interface {
	List(ctx context.Context, pageNumber int) (*license.ListLicenseResponse, error)
}

// NewAgent creates a new license agent that retrieves a license from the given url once every given period.
func NewAgent(
	konnectAPIClient UpstreamClient,
	logger logr.Logger,
) *Agent {
	return &Agent{
		logger:           logger,
		ticker:           time.NewTicker(PollingInterval),
		mutex:            sync.RWMutex{},
		konnectAPIClient: konnectAPIClient,
	}
}

// Agent handles retrieving a Kong license and providing it to other KIC subsystems.
type Agent struct {
	logger           logr.Logger
	ticker           *time.Ticker
	mutex            sync.RWMutex
	konnectAPIClient UpstreamClient

	// license is the current license retrieved from upstream.
	license license.Item
}

// NeedLeaderElection indicates if the Agent requires leadership to run. It always returns true.
func (a *Agent) NeedLeaderElection() bool {
	return true
}

// Start starts the Agent. It attempts to pull an initial license from upstream, and then polls for updates on a
// regular interval defined by PollingInterval.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.V(util.DebugLevel).Info("starting license agent")

	err := a.updateLicense(ctx)
	if err != nil {
		a.logger.Error(err, "could not retrieve license from upstream")
	}

	return a.run(ctx)
}

// GetLicense returns the agent's current license as a go-kong License struct. It omits the origin timestamps,
// as Kong will auto-populate these when adding the license to its config database.
func (a *Agent) GetLicense() kong.License {
	a.logger.V(util.DebugLevel).Info("retrieving license from cache")
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return kong.License{
		ID:      kong.String(a.license.ID),
		Payload: kong.String(a.license.License),
	}
}

// run updates the license on a regular interval until the context is cancelled.
func (a *Agent) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("context done, shutting down license agent")
			a.ticker.Stop()
			return ctx.Err()
		case <-a.ticker.C:
			a.logger.V(util.DebugLevel).Info("retrieving license from external service")
			if err := a.updateLicense(ctx); err != nil {
				a.logger.Error(err, "could not update license")
			}
		}
	}
}

// updateLicense retrievs a license from an outside system. If it successfully retrieves a license, it updates the
// in-memory license cache.
func (a *Agent) updateLicense(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, PollingTimeout)
	defer cancel()

	// This is an array because it's a Kong entity collection, even though we only expect to have exactly one license.
	licenses, err := a.konnectAPIClient.List(ctx, 0)
	if err != nil {
		return fmt.Errorf("could not retrieve license: %w", err)
	}
	if len(licenses.Items) == 0 {
		return fmt.Errorf("received empty license response")
	}
	license := licenses.Items[0]
	if license.UpdatedAt > a.license.UpdatedAt {
		a.logger.V(util.InfoLevel).Info("updating license cache",
			"old_updated_at", time.Unix(int64(a.license.UpdatedAt), 0).String(),
			"new_updated_at", time.Unix(int64(license.UpdatedAt), 0).String(),
		)
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.license = *license
	} else {
		a.logger.V(util.DebugLevel).Info("license cache is up to date")
	}

	return nil
}
