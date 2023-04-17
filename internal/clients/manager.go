package clients

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

	// discoveredAdminAPIsNotifyChan is used for notifications that contain Admin API
	// endpoints list that should be used for configuring the dataplane.
	discoveredAdminAPIsNotifyChan    chan []adminapi.DiscoveredAdminAPI
	gatewayClientsChangesSubscribers []chan struct{}

	ctx                   context.Context
	onceNotifyLoopRunning sync.Once
	notifyLoopRunningCh   chan struct{}
	isNotifyLoopRunning   bool

	// gatewayClients represent all Kong Gateway data-planes that are configured by this KIC instance with use of
	// their Admin API.
	gatewayClients []*adminapi.Client

	// konnectClient represents a special-case of the data-plane which is Konnect cloud.
	// This client is used to synchronise configuration with Konnect's Runtime Group Admin API.
	konnectClient *adminapi.Client

	// lock prevents concurrent access to the manager's fields.
	lock sync.RWMutex

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
		gatewayClients:                initialClients,
		adminAPIClientFactory:         kongClientFactory,
		discoveredAdminAPIsNotifyChan: make(chan []adminapi.DiscoveredAdminAPI),
		ctx:                           ctx,
		notifyLoopRunningCh:           make(chan struct{}),
		logger:                        logger,
	}, nil
}

// Running returns a channel that is closed when the manager's background tasks are already running.
func (c *AdminAPIClientsManager) Running() chan struct{} {
	return c.notifyLoopRunningCh
}

// RunNotifyLoop runs a goroutine that will dynamically ingest new addresses of Kong Admin API endpoints.
func (c *AdminAPIClientsManager) RunNotifyLoop() {
	c.onceNotifyLoopRunning.Do(func() {
		go c.adminAPIAddressNotifyLoop()

		c.lock.Lock()
		defer c.lock.Unlock()
		c.isNotifyLoopRunning = true
	})
}

// Notify receives a list of addresses that KongClient should use from now on as
// a list of Kong Admin API endpoints.
func (c *AdminAPIClientsManager) Notify(discoveredAPIs []adminapi.DiscoveredAdminAPI) {
	// Ensure here that we're not done.
	select {
	case <-c.ctx.Done():
		return
	default:
	}

	// And here also listen on c.ctx.Done() to allow the notification to be interrupted.
	select {
	case <-c.ctx.Done():
	case c.discoveredAdminAPIsNotifyChan <- discoveredAPIs:
	}
}

// SetKonnectClient sets a client that will be used to communicate with Konnect Runtime Group Admin API.
// If called multiple times, it will override the client.
func (c *AdminAPIClientsManager) SetKonnectClient(client *adminapi.Client) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.konnectClient = client
}

// AllClients returns a copy of current client's slice. It will also include Konnect client if set.
func (c *AdminAPIClientsManager) AllClients() []*adminapi.Client {
	c.lock.RLock()
	defer c.lock.RUnlock()

	copied := make([]*adminapi.Client, len(c.gatewayClients))
	copy(copied, c.gatewayClients)

	if c.konnectClient != nil {
		copied = append(copied, c.konnectClient)
	}

	return copied
}

// GatewayClients returns a copy of current client's slice. Konnect client won't be included.
// This method can be used when some actions need to be performed only against Kong Gateway clients.
func (c *AdminAPIClientsManager) GatewayClients() []*adminapi.Client {
	c.lock.RLock()
	defer c.lock.RUnlock()

	copied := make([]*adminapi.Client, len(c.gatewayClients))
	copy(copied, c.gatewayClients)
	return copied
}

func (c *AdminAPIClientsManager) GatewayClientsCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.gatewayClients)
}

// SubscribeToGatewayClientsChanges returns a channel that will receive a notification on every Gateway clients update.
// Can be used to receive a signal when immediate reaction to the changes is needed. After receiving the notification,
// GatewayClients call will return an already updated slice of clients.
// It will return `false` as a second result in case the notifications loop is not running (e.g. static clients setup
// is used and no updates are going to happen).
func (c *AdminAPIClientsManager) SubscribeToGatewayClientsChanges() (<-chan struct{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Context is already done, no subscriptions should be created.
	if c.ctx.Err() != nil {
		return nil, false
	}

	// No notify loop running, there will be no updates, let's propagate that to the caller.
	if !c.isNotifyLoopRunning {
		return nil, false
	}

	ch := make(chan struct{}, 1)
	c.gatewayClientsChangesSubscribers = append(c.gatewayClientsChangesSubscribers, ch)
	return ch, true
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
	close(c.notifyLoopRunningCh)
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("closing AdminAPIClientsManager: %s", c.ctx.Err())
			c.discoveredAdminAPIsNotifyChan = nil
			c.closeGatewayClientsSubscribers()
			return

		case discoveredAdminAPIs := <-c.discoveredAdminAPIsNotifyChan:
			// This call will only log errors e.g. during creation of new clients.
			// If need be we might consider propagating those errors up the stack.
			c.adjustGatewayClients(discoveredAdminAPIs)

			// Notify subscribers that the clients list has been updated.
			c.notifyGatewayClientsSubscribers()
		}
	}
}

// adjustGatewayClients adjusts internally stored clients slice based on the provided
// discovered Admin APIs slice. It consults BaseRootURLs of already stored clients with each
// of the discovered Admin APIs and creates only those clients that we don't have.
func (c *AdminAPIClientsManager) adjustGatewayClients(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Short circuit
	if len(discoveredAdminAPIs) == 0 {
		c.gatewayClients = c.gatewayClients[:0]
		return
	}

	toAdd := lo.Filter(discoveredAdminAPIs, func(api adminapi.DiscoveredAdminAPI, _ int) bool {
		// If we already have a client with a provided address then great, no need
		// to do anything.

		// If we don't have a client with new address then filter it and add
		// a client for this address.
		return !lo.ContainsBy(c.gatewayClients, func(cl *adminapi.Client) bool {
			return api.Address == cl.BaseRootURL()
		})
	})

	var idxToRemove []int
	for i, cl := range c.gatewayClients {
		// If the new address set contains a client that we already have then
		// good, no need to do anything for it.
		if lo.ContainsBy(discoveredAdminAPIs, func(api adminapi.DiscoveredAdminAPI) bool {
			return api.Address == cl.BaseRootURL()
		}) {
			continue
		}
		// If the new address set does not contain an address that we already
		// have then remove it.
		idxToRemove = append(idxToRemove, i)
	}

	for i := len(idxToRemove) - 1; i >= 0; i-- {
		idx := idxToRemove[i]
		c.gatewayClients = append(c.gatewayClients[:idx], c.gatewayClients[idx+1:]...)
	}

	for _, adminAPI := range toAdd {
		client, err := c.adminAPIClientFactory.CreateAdminAPIClient(c.ctx, adminAPI.Address)
		if err != nil {
			c.logger.WithError(err).Errorf("failed to create a client for %s", adminAPI)
			continue
		}
		client.AttachPodReference(adminAPI.PodRef)

		c.gatewayClients = append(c.gatewayClients, client)
	}
}

// notifyGatewayClientsSubscribers sends notifications to all subscribers that have called SubscribeToGatewayClientsChanges.
func (c *AdminAPIClientsManager) notifyGatewayClientsSubscribers() {
	for _, sub := range c.gatewayClientsChangesSubscribers {
		select {
		case <-c.ctx.Done():
			c.logger.Info("not sending notification to subscribers as the context is done")
			return
		case sub <- struct{}{}:
		}
	}
}

func (c *AdminAPIClientsManager) closeGatewayClientsSubscribers() {
	for _, sub := range c.gatewayClientsChangesSubscribers {
		close(sub)
	}
}

// AdminAPIClientsProvider allows fetching the most recent list of Admin API clients of Gateways that
// we should configure.
type AdminAPIClientsProvider interface {
	AllClients() []*adminapi.Client
	GatewayClients() []*adminapi.Client
}
