package adminapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	k8stypes "k8s.io/apimachinery/pkg/types"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// StatusClient is a client for checking Kong Gateway status via the dedicated status port.
type StatusClient struct {
	httpClient *http.Client
	baseURL    string
	podRef     *k8stypes.NamespacedName
}

// NewStatusClient creates a new status client for the given status API address.
func NewStatusClient(statusAPI DiscoveredAdminAPI, opts managercfg.AdminAPIClientConfig) (*StatusClient, error) {
	httpClient, err := makeHTTPClient(opts, "")
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client for status API: %w", err)
	}

	return &StatusClient{
		httpClient: httpClient,
		baseURL:    statusAPI.Address,
		podRef:     &statusAPI.PodRef,
	}, nil
}

// IsReady checks if the Kong Gateway is ready by calling the /status endpoint.
func (c *StatusClient) IsReady(ctx context.Context) error {
	statusURL := fmt.Sprintf("%s/status", c.baseURL)
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
	if err != nil {
		return fmt.Errorf("creating status request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status endpoint returned %d", resp.StatusCode)
	}

	return nil
}

// PodReference returns the Pod reference for this status client.
func (c *StatusClient) PodReference() (k8stypes.NamespacedName, bool) {
	if c.podRef != nil {
		return *c.podRef, true
	}
	return k8stypes.NamespacedName{}, false
}

// BaseRootURL returns the base URL for this status client.
func (c *StatusClient) BaseRootURL() string {
	return c.baseURL
}

// StatusClientFactory creates status clients for discovered status APIs.
type StatusClientFactory struct {
	logger logr.Logger
	opts   managercfg.AdminAPIClientConfig
}

// NewStatusClientFactory creates a new status client factory.
func NewStatusClientFactory(logger logr.Logger, opts managercfg.AdminAPIClientConfig) *StatusClientFactory {
	return &StatusClientFactory{
		logger: logger,
		opts:   opts,
	}
}

// CreateStatusClient creates a status client for the given discovered status API.
func (f *StatusClientFactory) CreateStatusClient(ctx context.Context, discoveredStatusAPI DiscoveredAdminAPI) (*StatusClient, error) {
	f.logger.V(1).Info(
		"Creating Kong Gateway Status API client",
		"address", discoveredStatusAPI.Address,
	)

	client, err := NewStatusClient(discoveredStatusAPI, f.opts)
	if err != nil {
		return nil, err
	}

	// Test the status client by calling IsReady
	if err := client.IsReady(ctx); err != nil {
		return nil, fmt.Errorf("status client not ready: %w", err)
	}

	return client, nil
}