package adminapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
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
func (c *Client) IsReady(ctx context.Context) error {
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

type ClientFactory struct {
	workspace      string
	httpClientOpts HTTPClientOpts
	adminToken     string
}

func NewClientFactoryForWorkspace(workspace string, httpClientOpts HTTPClientOpts, adminToken string) ClientFactory {
	return ClientFactory{
		workspace:      workspace,
		httpClientOpts: httpClientOpts,
		adminToken:     adminToken,
	}
}

func (cf ClientFactory) CreateAdminAPIClient(ctx context.Context, discoveredAdminAPI DiscoveredAdminAPI) (*Client, error) {
	httpclient, err := MakeHTTPClient(&cf.httpClientOpts, cf.adminToken)
	if err != nil {
		return nil, err
	}
	cl, err := NewKongClientForWorkspace(ctx, discoveredAdminAPI.Address, cf.workspace, httpclient)
	if err != nil {
		return nil, err
	}
	cl.AttachPodReference(discoveredAdminAPI.PodRef)
	return cl, nil
}
