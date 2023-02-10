package dataplane

import (
	"context"
	"errors"
	"sync"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

type ClientFactory interface {
	CreateAdminAPIClient(ctx context.Context, address string) (*adminapi.Client, error)
}

// AdminAPIClientsManager keeps track of current Admin API clients of Gateways that we should configure.
// In particular, it can be notified about the clients' list update with use of Notify method, and queried
// for the latest slice of those with use of Clients method.
type AdminAPIClientsManager struct {
	// adminAPIClientFactory is a factory used for creating Admin API clients.
	adminAPIClientFactory ClientFactory

	// adminAPIAddressNotifyChan is used for notifications that contain Admin API
	// endpoints list that should be used for configuring the dataplane.
	adminAPIAddressNotifyChan chan []string

	ctx         context.Context
	onceRunning sync.Once
	running     chan struct{}

	clients     []*adminapi.Client
	clientsLock sync.RWMutex

	logger logrus.FieldLogger
}

func NewAdminAPIClientsManager(
	ctx context.Context,
	logger logrus.FieldLogger,
	initialClients []*adminapi.Client,
	kongClientFactory ClientFactory,
) (*AdminAPIClientsManager, error) {
	if len(initialClients) == 0 {
		return nil, errors.New("at least one initial client must be provided")
	}

	return &AdminAPIClientsManager{
		clients:                   initialClients,
		adminAPIClientFactory:     kongClientFactory,
		adminAPIAddressNotifyChan: make(chan []string),
		ctx:                       ctx,
		running:                   make(chan struct{}),
		logger:                    logger,
	}, nil
}

// Running returns a channel that is closed when the manager's background tasks are already running.
func (c *AdminAPIClientsManager) Running() chan struct{} {
	return c.running
}

// RunNotifyLoop runs a goroutine that will dynamically ingest new addresses of Kong Admin API endpoints.
func (c *AdminAPIClientsManager) RunNotifyLoop() {
	c.onceRunning.Do(func() {
		go c.adminAPIAddressNotifyLoop()
	})
}

// Notify receives a list of addresses that KongClient should use from now on as
// a list of Kong Admin API endpoints.
func (c *AdminAPIClientsManager) Notify(addresses []string) {
	// Ensure here that we're not done.
	select {
	case <-c.ctx.Done():
		return
	default:
	}

	// And here also listen on c.ctx.Done() to allow the notification to be interrupted.
	select {
	case <-c.ctx.Done():
	case c.adminAPIAddressNotifyChan <- addresses:
	}
}

// Clients returns a copy of current clients slice.
func (c *AdminAPIClientsManager) Clients() []*adminapi.Client {
	c.clientsLock.RLock()
	defer c.clientsLock.RUnlock()

	copied := make([]*adminapi.Client, len(c.clients))
	copy(copied, c.clients)
	return copied
}

// adminAPIAddressNotifyLoop is an inner loop listening on notifyChan which are received via
// Notify() calls. Each time it receives on notifyChan tt will take the provided
// list of addresses and update the internally held list of clients such that:
//   - the internal list of kong clients contains only the provided addresses
//   - if a client for a provided address already exists it's not recreated again
//     (hence no external calls are made to check the provided endpoint if there
//     exists a client already using it)
//   - client that do not exist in the provided address list are removed if they
//     are present in the current state
//
// This function will acquire the internal lock to prevent the modification of
// internal clients list.
func (c *AdminAPIClientsManager) adminAPIAddressNotifyLoop() {
	close(c.running)
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("closing AdminAPIClientsManager: %s", c.ctx.Err())
			c.adminAPIAddressNotifyChan = nil
			return

		case addresses := <-c.adminAPIAddressNotifyChan:
			// This call will only log errors e.g. during creation of new clients.
			// If need be we might consider propagating those errors up the stack.
			c.adjustKongClients(addresses)
		}
	}
}

// adjustKongClients adjusts internally stored clients slice based on the provided
// addresses slice. It consults BaseRootURLs of already stored clients with each
// of the addreses and creates only those clients that we don't have.
func (c *AdminAPIClientsManager) adjustKongClients(addresses []string) {
	c.clientsLock.Lock()
	defer c.clientsLock.Unlock()

	toAdd := lo.Filter(addresses, func(addr string, _ int) bool {
		// If we already have a client with a provided address then great, no need
		// to do anything.

		// If we don't have a client with new address then filter it and add
		// a client for this address.
		return !lo.ContainsBy(c.clients, func(cl *adminapi.Client) bool {
			return addr == cl.BaseRootURL()
		})
	})

	var idxToRemove []int
	for i, cl := range c.clients {
		// If the new address set contains a client that we already have then
		// good, no need to do anything for it.
		if lo.Contains(addresses, cl.BaseRootURL()) {
			continue
		}
		// If the new address set does not contain an address that we already
		// have then remove it.
		idxToRemove = append(idxToRemove, i)
	}

	for i := len(idxToRemove) - 1; i >= 0; i-- {
		idx := idxToRemove[i]
		c.clients = append(c.clients[:idx], c.clients[idx+1:]...)
	}

	for _, addr := range toAdd {
		client, err := c.adminAPIClientFactory.CreateAdminAPIClient(c.ctx, addr)
		if err != nil {
			c.logger.WithError(err).Errorf("failed to create a client for %s", addr)
			continue
		}

		c.clients = append(c.clients, client)
	}
}
