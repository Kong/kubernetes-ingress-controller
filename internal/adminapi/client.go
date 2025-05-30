package adminapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// Client is a wrapper around raw *kong.Client. It's advised to pass this wrapper across the codebase, and
// fallback to the underlying *kong.Client only when it's passed to external functions that require
// it. Also, where it's possible, use a specific Abstract*Service interfaces that *kong.Client includes.
// Each Client holds its own PluginSchemaStore to cache plugins' schemas as they may theoretically differ between
// instances.
type Client struct {
	adminAPIClient      *kong.Client
	pluginSchemaStore   *util.PluginSchemaStore
	isKonnect           bool
	konnectControlPlane string

	lastCacheStoresHash store.SnapshotHash

	lastConfigSHALock sync.RWMutex
	lastConfigSHA     []byte
	// podRef (optional) describes the Pod that the Client communicates with.
	podRef *k8stypes.NamespacedName
	// statusClient (optional) is used for status checks instead of the admin API client
	statusClient *StatusClient
}

// NewClient creates an Admin API client that is to be used with a regular Admin API exposed by Kong Gateways.
func NewClient(c *kong.Client) *Client {
	return &Client{
		adminAPIClient:    c,
		pluginSchemaStore: util.NewPluginSchemaStore(c),
	}
}

// NewTestClient creates a client for test purposes.
func NewTestClient(address string) (*Client, error) {
	kongClient, err := kong.NewTestClient(lo.ToPtr(address), &http.Client{})
	if err != nil {
		return nil, err
	}

	return NewClient(kongClient), nil
}

type KonnectClient struct {
	Client
	consumersSyncDisabled bool
	backoffStrategy       UpdateBackoffStrategy
}

// NewKonnectClient creates an Admin API client that is to be used with a Konnect Control Plane Admin API.
func NewKonnectClient(c *kong.Client, controlPlane string, consumersSyncDisabled bool) *KonnectClient {
	return &KonnectClient{
		Client: Client{
			adminAPIClient:      c,
			isKonnect:           true,
			konnectControlPlane: controlPlane,
			pluginSchemaStore:   util.NewPluginSchemaStore(c),
		},
		backoffStrategy:       NewKonnectBackoffStrategy(clock.System{}),
		consumersSyncDisabled: consumersSyncDisabled,
	}
}

func (c *KonnectClient) BackoffStrategy() UpdateBackoffStrategy {
	return c.backoffStrategy
}

func (c *KonnectClient) ConsumersSyncDisabled() bool {
	return c.consumersSyncDisabled
}

// AdminAPIClient returns an underlying go-kong's Admin API client.
func (c *Client) AdminAPIClient() *kong.Client {
	return c.adminAPIClient
}

// BaseRootURL returns a base address used for communicating with the Admin API.
func (c *Client) BaseRootURL() string {
	return c.adminAPIClient.BaseRootURL()
}

func (c *Client) NodeID(ctx context.Context) (string, error) {
	data, err := c.adminAPIClient.Root(ctx)
	if err != nil {
		return "", fmt.Errorf("failed fetching Kong client root: %w", err)
	}

	const nodeIDKey = "node_id"
	nodeID, err := extractStringFromRoot(data, nodeIDKey)
	if err != nil {
		return "", fmt.Errorf("malformed node ID found in Kong client root: %w", err)
	}

	return nodeID, nil
}

// IsReady returns nil if the Admin API is ready to serve requests.
// If a status client is attached, it will be used for the readiness check instead of the admin API.
func (c *Client) IsReady(ctx context.Context) error {
	if c.statusClient != nil {
		return c.statusClient.IsReady(ctx)
	}
	_, err := c.adminAPIClient.Status(ctx)
	return err
}

// GetKongVersion returns version of the kong gateway.
func (c *Client) GetKongVersion(ctx context.Context) (string, error) {
	if c.isKonnect {
		return "", errors.New("cannot get kong version from konnect")
	}
	rootConfig, err := c.adminAPIClient.Root(ctx)
	if err != nil {
		return "", fmt.Errorf("failed fetching Kong client root: %w", err)
	}

	versionKey := "version"
	version, err := extractStringFromRoot(rootConfig, versionKey)
	if err != nil {
		return "", fmt.Errorf("malformed Kong version found in Kong client root: %w", err)
	}

	return version, nil
}

func extractStringFromRoot(data map[string]interface{}, key string) (string, error) {
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("%q key not found", key)
	}

	valStr, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%q key is not a string, actual type: %T", key, val)
	}

	return valStr, nil
}

// PluginSchemaStore returns client's PluginSchemaStore.
func (c *Client) PluginSchemaStore() *util.PluginSchemaStore {
	return c.pluginSchemaStore
}

// IsKonnect tells if a client is used for communication with Konnect Control Plane Admin API.
func (c *Client) IsKonnect() bool {
	return c.isKonnect
}

// KonnectControlPlane gets a unique identifier of a Konnect's Control Plane that config should
// be synchronised with. Empty in case of non-Konnect clients.
func (c *Client) KonnectControlPlane() string {
	if !c.isKonnect {
		return ""
	}

	return c.konnectControlPlane
}

// SetLastCacheStoresHash overrides last cache stores hash.
func (c *Client) SetLastCacheStoresHash(s store.SnapshotHash) {
	c.lastCacheStoresHash = s
}

// LastCacheStoresHash returns a checksum of the last successful cache stores push.
func (c *Client) LastCacheStoresHash() store.SnapshotHash {
	return c.lastCacheStoresHash
}

// SetLastConfigSHA overrides last config SHA.
func (c *Client) SetLastConfigSHA(s []byte) {
	c.lastConfigSHALock.Lock()
	defer c.lastConfigSHALock.Unlock()
	c.lastConfigSHA = s
}

// LastConfigSHA returns a checksum of the last successful configuration push.
func (c *Client) LastConfigSHA() []byte {
	c.lastConfigSHALock.RLock()
	defer c.lastConfigSHALock.RUnlock()
	return c.lastConfigSHA
}

// AttachPodReference allows attaching a Pod reference to the client. Should be used in case we know what Pod the client
// will communicate with (e.g. when the gateway service discovery is used).
func (c *Client) AttachPodReference(podNN k8stypes.NamespacedName) {
	c.podRef = &podNN
}

// PodReference returns an optional reference to the Pod the client communicates with.
func (c *Client) PodReference() (k8stypes.NamespacedName, bool) {
	if c.podRef != nil {
		return *c.podRef, true
	}
	return k8stypes.NamespacedName{}, false
}

// AttachStatusClient allows attaching a status client to the admin API client for status checks.
func (c *Client) AttachStatusClient(statusClient *StatusClient) {
	c.statusClient = statusClient
}

type ClientFactory struct {
	logger               logr.Logger
	workspace            string
	opts                 managercfg.AdminAPIClientConfig
	adminToken           string
	statusAPIsDiscoverer *Discoverer
}

func NewClientFactoryForWorkspace(
	logger logr.Logger,
	workspace string,
	clientOpts managercfg.AdminAPIClientConfig,
	adminToken string,
) ClientFactory {
	return ClientFactory{
		logger:     logger,
		workspace:  workspace,
		opts:       clientOpts,
		adminToken: adminToken,
	}
}

func NewClientFactoryForWorkspaceWithStatusDiscoverer(
	logger logr.Logger,
	workspace string,
	clientOpts managercfg.AdminAPIClientConfig,
	adminToken string,
	statusAPIsDiscoverer *Discoverer,
) ClientFactory {
	return ClientFactory{
		logger:               logger,
		workspace:            workspace,
		opts:                 clientOpts,
		adminToken:           adminToken,
		statusAPIsDiscoverer: statusAPIsDiscoverer,
	}
}

func (cf ClientFactory) CreateAdminAPIClient(ctx context.Context, discoveredAdminAPI DiscoveredAdminAPI) (*Client, error) {
	cf.logger.V(logging.DebugLevel).Info(
		"Creating Kong Gateway Admin API client",
		"address", discoveredAdminAPI.Address, "tlsServerName", discoveredAdminAPI.TLSServerName,
	)
	opts := cf.opts
	opts.TLSServerName = discoveredAdminAPI.TLSServerName

	cl, err := NewKongClientForWorkspace(ctx, discoveredAdminAPI.Address, cf.workspace, opts, cf.adminToken)
	if err != nil {
		return nil, err
	}

	cl.AttachPodReference(discoveredAdminAPI.PodRef)

	// If we have a status APIs discoverer, try to find and attach a status client
	if cf.statusAPIsDiscoverer != nil {
		if statusClient := cf.tryCreateStatusClient(ctx, discoveredAdminAPI.PodRef); statusClient != nil {
			cl.AttachStatusClient(statusClient)
			cf.logger.V(logging.DebugLevel).Info(
				"Attached status client to admin API client",
				"adminAddress", discoveredAdminAPI.Address,
				"statusAddress", statusClient.BaseRootURL(),
			)
		}
	}

	return cl, nil
}

// tryCreateStatusClient attempts to create a status client for the same pod as the admin API client.
//
//nolint:unparam // This is a stub implementation that always returns nil for now
func (cf ClientFactory) tryCreateStatusClient(_ context.Context, podRef k8stypes.NamespacedName) *StatusClient {
	// Try to discover status APIs for the same service that the admin API belongs to
	// We'll use the pod reference to find the corresponding status endpoint

	// This is a simplified implementation that assumes the status API is on the same pod
	// but on a different port (8100 instead of 8444)

	// In a real implementation, you would use proper service discovery here
	// For now, we'll return nil to keep the existing behavior
	// The status client creation would need to be implemented based on your specific requirements

	cf.logger.V(logging.DebugLevel).Info(
		"Status client creation not yet implemented",
		"podRef", podRef,
	)

	return nil
}
