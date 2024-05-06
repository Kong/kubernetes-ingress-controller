package license

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/tidwall/gjson"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
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

// KonnectLicense is a license retrieved from Konnect.
type KonnectLicense struct {
	ID        string
	Payload   string
	UpdatedAt time.Time
}

type KonnectLicenseClient interface {
	Get(ctx context.Context) (mo.Option[KonnectLicense], error)
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

type Ticker interface {
	Stop()
	Channel() <-chan time.Time
	Reset(d time.Duration)
}

// WithTicker sets the ticker in Agent. This is useful for testing.
// Ticker doesn't define the period, it defines the implementation of ticking.
func WithTicker(t Ticker) AgentOpt {
	return func(a *Agent) {
		a.ticker = t
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
		// Note: the ticker defines the implementation of ticking, not the period.
		ticker:    clock.NewTicker(),
		startedCh: make(chan struct{}),
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
	ticker               Ticker
	startedCh            chan struct{}

	// cachedLicense is the current license retrieved from upstream. It's optional because we may not have retrieved a
	// license yet.
	cachedLicense mo.Option[KonnectLicense]
	mutex         sync.RWMutex
}

// NeedLeaderElection indicates if the Agent requires leadership to run. It always returns true.
func (a *Agent) NeedLeaderElection() bool {
	return true
}

// Start starts the Agent. It attempts to pull an initial license from upstream, and then polls for updates on a
// regular period, either the agent's initialPollingPeriod if it has not yet obtained a license or regularPollingPeriod if it has.
func (a *Agent) Start(ctx context.Context) error {
	a.logger.V(util.DebugLevel).Info("Starting license agent")

	err := a.reconcileLicenseWithKonnect(ctx)
	if err != nil {
		// If that happens, GetLicense() will return no license until we retrieve a valid one in polling.
		a.logger.Error(err, "Could not retrieve license from upstream")
	}

	return a.runPollingLoop(ctx)
}

// GetLicense returns the agent's current license as a go-kong License struct. It omits the origin timestamps,
// as Kong will auto-populate these when adding the license to its config database.
// It's optional because we may not have retrieved a license yet.
func (a *Agent) GetLicense() mo.Option[kong.License] {
	a.logger.V(util.DebugLevel).Info("Retrieving license from cache")
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if cachedLicense, ok := a.cachedLicense.Get(); ok {
		return mo.Some(kong.License{
			ID:      lo.ToPtr(cachedLicense.ID),
			Payload: lo.ToPtr(cachedLicense.Payload),
		})
	}

	return mo.None[kong.License]()
}

// Started returns a channel which will be closed when the Agent has started.
func (a *Agent) Started() <-chan struct{} {
	return a.startedCh
}

// runPollingLoop updates the license on a regular period until the context is cancelled.
// It will run at a faster period initially, and then switch to the regular period once a license is retrieved.
func (a *Agent) runPollingLoop(ctx context.Context) error {
	a.ticker.Reset(a.initialPollingPeriod)
	defer a.ticker.Stop()

	ch := a.ticker.Channel()
	close(a.startedCh)
	for {
		select {
		case <-ch:
			a.logger.V(util.DebugLevel).Info("Retrieving license from external service")
			if err := a.reconcileLicenseWithKonnect(ctx); err != nil {
				a.logger.Error(err, "Could not reconcile license with Konnect")
			}
			// Reset the ticker to run with the expected period which may change after we receive the license.
			a.ticker.Reset(a.resolvePollingPeriod())
		case <-ctx.Done():
			a.logger.Info("Context done, shutting down license agent")
			return ctx.Err()
		}
	}
}

func (a *Agent) resolvePollingPeriod() time.Duration {
	cached, ok := a.cachedLicense.Get()
	// With no license available, update frequently to retrieve it when it appears, e.g. when a user upgrades from Free
	// to Enterprise tier).
	if !ok {
		return a.initialPollingPeriod
	}
	// The current license is expired, update more often to try and fix it.
	if IsExpiredLicense(cached.Payload) {
		return a.initialPollingPeriod
	}
	// We already have a license, update at the slower interval.
	return a.regularPollingPeriod
}

// reconcileLicenseWithKonnect retrieves a license from upstream and caches it if it is newer than the cached license or there is no cached license.
func (a *Agent) reconcileLicenseWithKonnect(ctx context.Context) error {
	retrievedLicenseOpt, err := a.retrieveLicenseFromUpstream(ctx)
	if err != nil {
		return err
	}

	retrievedLicense, retrievedLicenseOk := retrievedLicenseOpt.Get()
	if !retrievedLicenseOk {
		// If we get no license from Konnect, we cannot do anything.
		a.logger.V(util.DebugLevel).Info("No license found in Konnect")
		return nil
	}

	if a.cachedLicense.IsAbsent() {
		a.logger.V(util.InfoLevel).Info("Caching initial license retrieved from the upstream",
			"updated_at", retrievedLicense.UpdatedAt.String(),
		)
		a.updateCache(retrievedLicense)
	} else if cachedLicense, ok := a.cachedLicense.Get(); ok {
		if IsExpiredLicense(cachedLicense.Payload) {
			a.logger.V(util.InfoLevel).Info("Caching license retrieved from the upstream as current license is expired",
				"cached_updated_at", cachedLicense.UpdatedAt.String(),
				"retrieved_updated_at", retrievedLicense.UpdatedAt.String(),
			)
			a.updateCache(retrievedLicense)
		} else if retrievedLicense.UpdatedAt.After(cachedLicense.UpdatedAt) {
			a.logger.V(util.InfoLevel).Info("Caching license retrieved from the upstream as it is newer than the cached one",
				"cached_updated_at", cachedLicense.UpdatedAt.String(),
				"retrieved_updated_at", retrievedLicense.UpdatedAt.String(),
			)
			a.updateCache(retrievedLicense)
		}
	} else {
		a.logger.V(util.DebugLevel).Info("License cache is up to date")
	}

	return nil
}

func (a *Agent) retrieveLicenseFromUpstream(ctx context.Context) (mo.Option[KonnectLicense], error) {
	ctx, cancel := context.WithTimeout(ctx, PollingTimeout)
	defer cancel()
	return a.konnectLicenseClient.Get(ctx)
}

func (a *Agent) updateCache(license KonnectLicense) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.cachedLicense = mo.Some(license)
}

// IsExpiredLicense returns true if given a well-formatted license string with an expiration date before now. It
// returns false if the license expiration is after now, or if the license status is unknown due to a parse error.
func IsExpiredLicense(license string) bool {
	// Licenses from Kong APIs do not have formally-defined schemas in those APIs. License objects consist only of an ID,
	// created/updated times, and an opaque payload string that we expect to be valid JSON with an expiration field.
	// As we don't care about the rest of the license contents, this performs a basic path extraction and tries to parse
	// a date out of it.
	if !gjson.Valid(license) {
		return false
	}
	expiry := gjson.Get(license, "license.payload.license_expiration_date")
	if !expiry.Exists() {
		return false
	}
	date, err := time.Parse(time.DateOnly, expiry.String())
	if err != nil {
		return false
	}
	if date.Before(time.Now()) {
		return true
	}
	return false
}
