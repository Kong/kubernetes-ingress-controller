package license

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// NewLicenseAgent creates a new license agent that retrieves a license from the given url once every given period.
func NewLicenseAgent(ctx context.Context, period time.Duration, url string) *Agent {
	client := http.Client{} // TODO pass in a konnect client instead
	return &Agent{
		logger:      logrus.New(), // TODO figure out how we're supposed to actually create new loggers
		client:      client,
		upstreamURL: url,
		ticker:      time.NewTicker(period),
	}
}

// License represents the response body of the upstream license API
type License struct {
	License string    `json:"license,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

// Agent handles retrieving a Kong license and providing it to other KIC subsystems.
type Agent struct {
	license     License // TODO maybe separate types
	logger      logrus.FieldLogger
	client      http.Client // TODO this needs to be a Konnect client eventually
	upstreamURL string
	ticker      *time.Ticker
}

// NeedLeaderElection indicates if the Agent requires leadership to run. It always returns true.
func (a *Agent) NeedLeaderElection() bool {
	return true
}

// Start starts the Agent. It attempts to pull an initial license from upstream, and failing that, pulls it from local
// cache. If both fail, startup fails.
func (a *Agent) Start(ctx context.Context) error {
	updateTimeout, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()
	err := a.UpdateLicense(updateTimeout)
	if err != nil {
		a.logger.WithError(err).Errorf("could not retrieve license from upstream")
		err := a.UpdateLicenseFromCache(ctx)
		if err != nil {
			return err
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
			a.logger.Infof("context done, shutting down license agent")
			a.ticker.Stop()
			return
		case <-a.ticker.C:
			updateTimeout, cancel := context.WithTimeout(ctx, time.Minute*5)
			defer cancel()
			if err := a.UpdateLicense(updateTimeout); err != nil {
				a.logger.WithError(err).Errorf("could not update license")
			}
		}
	}
}

// Update retrievs a license from an outside system. If it successfully retrieves a license, it TODO.
func (a *Agent) UpdateLicense(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, "GET", a.upstreamURL, nil)
	if err != nil {
		return err
	}
	response, err := a.client.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var license License
	err = json.Unmarshal(body, &license)
	if err != nil {
		return err
	}
	if license.Updated.After(a.license.Updated) {
		a.license = license

		err = persistLicense(license.License)
		if err != nil {
			a.logger.WithError(err).Errorf("could not store license in Secret")
		}
	}
	return nil
}

// UpdateLicenseFromCache retrieves a license from a local cache.
func (a *Agent) UpdateLicenseFromCache(ctx context.Context) error {
	// TODO make this not a stub
	return nil
}

// GetLicense returns the agent's current license.
func (a *Agent) GetLicense() string {
	return a.license.License
}

// PersistLicense saves the current license to a Secret.
func persistLicense(license string) error {
	// TODO make this not a stub
	return nil
}
