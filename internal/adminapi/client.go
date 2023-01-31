package adminapi

import (
	"github.com/kong/go-kong/kong"
)

// Client is a wrapper around *kong.Client. It's needed to be able to distinguish between clients
// that are to be used with a regular Kong Gateway Admin API, and the ones that are to be used with
// Konnect Runtime Group Admin API.
// The distinction is needed to be able to tell what protocol (deck or dbless) should be used when
// updating configuration using the client.
type Client struct {
	*kong.Client

	isKonnect           bool
	konnectRuntimeGroup string
}

// NewClient creates an Admin API client that is to be used with a regular Admin API exposed by Kong Gateways.
func NewClient(c *kong.Client) *Client {
	return &Client{Client: c}
}

// NewKonnectClient creates an Admin API client that is to be used with a Konnect Runtime Group Admin API.
func NewKonnectClient(c *kong.Client, runtimeGroup string) *Client {
	return &Client{
		Client:              c,
		isKonnect:           true,
		konnectRuntimeGroup: runtimeGroup,
	}
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
