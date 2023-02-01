package adminapi

import (
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// Client is a wrapper around *kong.Client. It's needed to be able to distinguish between clients
// that are to be used with a regular Kong Gateway Admin API, and the ones that are to be used with
// Konnect Runtime Group Admin API.
// The distinction is needed to be able to tell what protocol (deck or dbless) should be used when
// updating configuration using the client.
// Each client holds its own PluginSchemaStore to cache plugins' schemas as they may theoretically differ between
// instances.
type Client struct {
	adminAPIClient      *kong.Client
	pluginSchemaStore   *util.PluginSchemaStore
	isKonnect           bool
	konnectRuntimeGroup string

	lastConfigSHA []byte
}

// NewClient creates an Admin API client that is to be used with a regular Admin API exposed by Kong Gateways.
func NewClient(c *kong.Client) Client {
	return Client{
		adminAPIClient:    c,
		pluginSchemaStore: util.NewPluginSchemaStore(c),
	}
}

// NewKonnectClient creates an Admin API client that is to be used with a Konnect Runtime Group Admin API.
func NewKonnectClient(c *kong.Client, runtimeGroup string) Client {
	return Client{
		adminAPIClient:      c,
		isKonnect:           true,
		konnectRuntimeGroup: runtimeGroup,
		pluginSchemaStore:   util.NewPluginSchemaStore(c),
	}
}

// AdminAPIClient returns an underlying go-kong's Admin API client.
func (c *Client) AdminAPIClient() *kong.Client {
	return c.adminAPIClient
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
