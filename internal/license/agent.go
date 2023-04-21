package license

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// NewLicenseAgent creates a new license agent that retrieves a license from the given url once every given period.
func NewLicenseAgent(
	period time.Duration,
	url string,
	konnectAPIClient *konnect.LicenseAPIClient,
	logger logr.Logger,
) *Agent {
	return &Agent{
		logger:           logger,
		upstreamURL:      url,
		ticker:           time.NewTicker(period),
		mutex:            sync.RWMutex{},
		konnectAPIClient: konnectAPIClient,
	}
}

// Agent handles retrieving a Kong license and providing it to other KIC subsystems.
type Agent struct {
	license          konnect.LicenseItem
	logger           logr.Logger
	upstreamURL      string
	ticker           *time.Ticker
	mutex            sync.RWMutex
	konnectAPIClient *konnect.LicenseAPIClient
}

// NeedLeaderElection indicates if the Agent requires leadership to run. It always returns true.
func (a *Agent) NeedLeaderElection() bool {
	return true
}

// Start starts the Agent. It attempts to pull an initial license from upstream, and failing that, pulls it from local
// cache. If both fail, startup fails.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.V(util.DebugLevel).Info("starting license agent")
	updateTimeout, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()
	err := a.UpdateLicense(updateTimeout)
	if err != nil {
		a.logger.Error(err, "could not retrieve license from upstream")
		err := a.UpdateLicenseFromCache(ctx)
		if err != nil {
			a.logger.Error(err, "could not retrieve license from local cache")
		}
	}
	go a.Run(ctx)
	return nil
}

// Run updates the license on a regular interval until the context is cancelled.
func (a *Agent) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("context done, shutting down license agent")
			a.ticker.Stop()
			return
		case <-a.ticker.C:
			a.logger.V(util.DebugLevel).Info("retrieving license from external service")
			updateTimeout, cancel := context.WithTimeout(ctx, time.Minute*5)
			defer cancel()
			if err := a.UpdateLicense(updateTimeout); err != nil {
				a.logger.Error(err, "could not update license")
			}
		}
	}
}

// UpdateLicense retrievs a license from an outside system. If it successfully retrieves a license, it updates the in-memory
// and persistent license caches.
func (a *Agent) UpdateLicense(ctx context.Context) error {
	// TODO this is an array because it's a Kong entity collection, even though we only expect to have
	// exactly one license. this is manageable, but a bit messy
	licenses, err := a.konnectAPIClient.List(ctx, 0)
	if err != nil {
		return fmt.Errorf("could not retrieve license: %w", err)
	}
	if len(licenses.Items) == 0 {
		return fmt.Errorf("received empty license response")
	}
	license := licenses.Items[0]
	if license.UpdatedAt > a.license.UpdatedAt {
		a.logger.V(util.DebugLevel).Info("retrieved license has later expiration than current license, updating license cache")
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.license = *license

		err = persistLicense(license.License)
		if err != nil {
			a.logger.Error(err, "could not store license in Secret")
		}
	}
	return nil
}

// UpdateLicenseFromCache retrieves a license from a local cache.
func (a *Agent) UpdateLicenseFromCache(_ context.Context) error {
	// TODO make this not a stub https://github.com/Kong/kubernetes-ingress-controller/issues/3923
	return fmt.Errorf("not implemented")
}

// GetLicense returns the agent's current license as a go-kong License struct. It omits the origin timestamps,
// as Kong will auto-populate these when adding the license to its config database.
func (a *Agent) GetLicense() kong.License {
	a.mutex.RLock()
	a.logger.V(util.DebugLevel).Info("retrieving license from cache")
	defer a.mutex.RUnlock()
	return kong.License{
		ID:      kong.String(a.license.ID),
		Payload: kong.String(a.license.License),
	}
}

// PersistLicense saves the current license to a Secret.
func persistLicense(_ string) error {
	// TODO make this not a stub https://github.com/Kong/kubernetes-ingress-controller/issues/3923
	return nil
}
