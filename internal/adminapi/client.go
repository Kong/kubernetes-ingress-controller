package adminapi

import (
	"context"
	"net/http"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
	konnectRuntimeGroup string
	lastConfigSHA       []byte
}

// NewClient creates an Admin API client that is to be used with a regular Admin API exposed by Kong Gateways.
func NewClient(c *kong.Client) *Client {
	return &Client{
		adminAPIClient:    c,
		pluginSchemaStore: util.NewPluginSchemaStore(c),
	}
}

// NewKonnectClient creates an Admin API client that is to be used with a Konnect Runtime Group Admin API.
func NewKonnectClient(c *kong.Client, runtimeGroup string) *Client {
	return &Client{
		adminAPIClient:      c,
		isKonnect:           true,
		konnectRuntimeGroup: runtimeGroup,
		pluginSchemaStore:   util.NewPluginSchemaStore(c),
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

// AdminAPIClient returns an underlying go-kong's Admin API client.
func (c *Client) AdminAPIClient() *kong.Client {
	return c.adminAPIClient
}

// BaseRootURL returns a base address used for communicating with the Admin API.
func (c *Client) BaseRootURL() string {
	return c.adminAPIClient.BaseRootURL()
}

// PluginSchemaStore returns client's PluginSchemaStore.
func (c *Client) PluginSchemaStore() *util.PluginSchemaStore {
	return c.pluginSchemaStore
}

// IsKonnect tells if a client is used for communication with Konnect Runtime Group Admin API.
func (c *Client) IsKonnect() bool {
	return c.isKonnect
}

// KonnectRuntimeGroup gets a unique identifier of a Konnect's Runtime Group that config should
// be synchronised with. Empty in case of non-Konnect clients.
func (c *Client) KonnectRuntimeGroup() string {
	if !c.isKonnect {
		return ""
	}

	return c.konnectRuntimeGroup
}

// SetLastConfigSHA overrides last config SHA.
func (c *Client) SetLastConfigSHA(s []byte) {
	c.lastConfigSHA = s
}

// LastConfigSHA returns a checksum of the last successful configuration push.
func (c *Client) LastConfigSHA() []byte {
	return c.lastConfigSHA
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

func (cf ClientFactory) CreateAdminAPIClient(ctx context.Context, address string) (*Client, error) {
	httpclient, err := MakeHTTPClient(&cf.httpClientOpts, cf.adminToken)
	if err != nil {
		return nil, err
	}
	return NewKongClientForWorkspace(ctx, address, cf.workspace, httpclient)
}
