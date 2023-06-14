package license

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// DefaultPollingPeriod is the period at which the license agent will poll for license updates by default.
	DefaultPollingPeriod = time.Hour * 12

	// DefaultInitialPollingPeriod is the period at which the license agent will poll for a license until it retrieves
	// one.
	DefaultInitialPollingPeriod = time.Minute

	// PollingTimeout is the timeout for retrieving a license from upstream.
	PollingTimeout = time.Minute * 5
)

type KonnectLicenseClient interface {
	List(ctx context.Context, pageNumber int) (*license.ListLicenseResponse, error)
}

type AgentOpt func(*Agent)

// WithInitialPollingPeriod sets the initial polling period for the license agent.
func WithInitialPollingPeriod(initialPollingPeriod time.Duration) AgentOpt {
	return func(a *Agent) {
		a.initialPollingPeriod = initialPollingPeriod
	}
}

// WithPollingPeriod sets the regular polling period for the license agent.
func WithPollingPeriod(regularPollingPeriod time.Duration) AgentOpt {
	return func(a *Agent) {
		a.regularPollingPeriod = regularPollingPeriod
	}
}

// NewAgent creates a new license agent that retrieves a license from the given url once every given period.
func NewAgent(
	konnectLicenseClient KonnectLicenseClient,
	logger logr.Logger,
	opts ...AgentOpt,
) *Agent {
	a := &Agent{
		logger:               logger,
		konnectLicenseClient: konnectLicenseClient,
		initialPollingPeriod: DefaultInitialPollingPeriod,
		regularPollingPeriod: DefaultPollingPeriod,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Agent handles retrieving a Kong license and providing it to other KIC subsystems.
type Agent struct {
	logger               logr.Logger
	konnectLicenseClient KonnectLicenseClient
	initialPollingPeriod time.Duration
	regularPollingPeriod time.Duration

	// cachedLicense is the current license retrieved from upstream.
	cachedLicense mo.Option[license.Item]
	mutex         sync.RWMutex
}

// NeedLeaderElection indicates if the Agent requires leadership to runPollingLoop. It always returns true.
func (a *Agent) NeedLeaderElection() bool {
	return true
}

// Start starts the Agent. It attempts to pull an initial license from upstream, and then polls for updates on a
// regular period defined by DefaultRegularPollingInternval.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.V(util.DebugLevel).Info("starting license agent")

	err := a.reconcileLicenseWithKonnect(ctx)
	if err != nil {
		// If that happens, GetLicense() will return no license until we retrieve a valid one in polling.
		a.logger.Error(err, "could not retrieve license from upstream")
	}

	return a.runPollingLoop(ctx)
}

// GetLicense returns the agent's current license as a go-kong License struct. It omits the origin timestamps,
// as Kong will auto-populate these when adding the license to its config database.
// If the agent has not yet retrieved a license, it returns an empty struct and false.
func (a *Agent) GetLicense() (kong.License, bool) {
	a.logger.V(util.DebugLevel).Info("retrieving license from cache")
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if licenseItem, ok := a.cachedLicense.Get(); ok {
		return kong.License{
			ID:      kong.String(licenseItem.ID),
			Payload: kong.String(licenseItem.License),
		}, true
	}

	return kong.License{}, false
}

// runPollingLoop updates the license on a regular period until the context is cancelled.
// It will run at a faster period initially, and then switch to the regular period once a license is retrieved.
func (a *Agent) runPollingLoop(ctx context.Context) error {
	ticker := time.NewTicker(a.resolvePollingPeriod())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.logger.V(util.DebugLevel).Info("retrieving license from external service")
			if err := a.reconcileLicenseWithKonnect(ctx); err != nil {
				a.logger.Error(err, "could not reconcile license with Konnect")
			}
			// Reset the ticker to run with the expected period which may change after we receive the license.
			ticker.Reset(a.resolvePollingPeriod())
		case <-ctx.Done():
			a.logger.Info("context done, shutting down license agent")
			return ctx.Err()
		}
	}
}

func (a *Agent) resolvePollingPeriod() time.Duration {
	// If we already have a license, start with the regular polling period (happy path) ...
	if a.cachedLicense.IsPresent() {
		return a.regularPollingPeriod
	}
	// ... otherwise, start with the initial polling period which is shorter by default (to get a license faster
	// when it appears, e.g. when a user upgrades from Free to Enterprise tier).
	return a.initialPollingPeriod
}

// reconcileLicenseWithKonnect retrieves a license from upstream and caches it if it is newer than the cached license.
// When it's the first time retrieving a license, it will always cache it.
func (a *Agent) reconcileLicenseWithKonnect(ctx context.Context) error {
	updatedAtAsString := func(updatedAt uint64) string {
		return time.Unix(int64(updatedAt), 0).String()
	}

	retrievedLicense, err := a.retrieveLicenseFromUpstream(ctx)
	if err != nil {
		// If the license is not found, we do not return an error as it's expected when a customer is on the Free tier.
		if errors.Is(err, license.ErrNotFound) {
			a.logger.V(util.DebugLevel).Info("no license found in Konnect")
			return nil
		}
		return err
	}

	if a.cachedLicense.IsAbsent() {
		a.logger.V(util.InfoLevel).Info("caching license retrieved from the upstream",
			"updated_at", updatedAtAsString(retrievedLicense.UpdatedAt),
		)
		a.updateCache(retrievedLicense)
	} else if cachedLicense, ok := a.cachedLicense.Get(); ok && retrievedLicense.UpdatedAt > cachedLicense.UpdatedAt {
		a.logger.V(util.InfoLevel).Info("caching license retrieved from the upstream as it is newer than the cached one",
			"cached_updated_at", updatedAtAsString(cachedLicense.UpdatedAt),
			"retrieved_updated_at", updatedAtAsString(retrievedLicense.UpdatedAt),
		)
		a.updateCache(retrievedLicense)
	} else {
		a.logger.V(util.DebugLevel).Info("license cache is up to date")
	}

	return nil
}

func (a *Agent) retrieveLicenseFromUpstream(ctx context.Context) (*license.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, PollingTimeout)
	defer cancel()

	// This is an array because it's a Kong entity collection, even though we only expect to have exactly one license.
	licenses, err := a.konnectLicenseClient.List(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve license: %w", err)
	}
	if len(licenses.Items) == 0 {
		return nil, fmt.Errorf("received empty license response")
	}
	return licenses.Items[0], nil
}

func (a *Agent) updateCache(license *license.Item) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.cachedLicense = mo.Some(*license)
}
