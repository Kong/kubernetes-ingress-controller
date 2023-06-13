package clients

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// DefaultReadinessReconciliationInterval is the interval at which the manager will run readiness reconciliation loop.
// It's the same as the default interval of the readiness probe.
const DefaultReadinessReconciliationInterval = 10 * time.Second

type ClientFactory interface {
	CreateAdminAPIClient(ctx context.Context, discoveredAdminAPI adminapi.DiscoveredAdminAPI) (*adminapi.Client, error)
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

	// activeGatewayClients represent all Kong Gateway data-planes that are ready to be configured.
	activeGatewayClients []*adminapi.Client

	// pendingGatewayClients represent all Kong Gateway data-planes that were discovered but are not ready to be
	// configured.
	pendingGatewayClients []adminapi.DiscoveredAdminAPI

	// readinessReconciliationInterval is the interval at which the manager will run clients' readiness reconciliation.
	readinessReconciliationInterval time.Duration

	// konnectClient represents a special-case of the data-plane which is Konnect cloud.
	// This client is used to synchronise configuration with Konnect's Runtime Group Admin API.
	konnectClient *adminapi.KonnectClient

	// lock prevents concurrent access to the manager's fields.
	lock sync.RWMutex

	logger logrus.FieldLogger
}

func NewAdminAPIClientsManager(
	ctx context.Context,
	logger logrus.FieldLogger,
	initialClients []*adminapi.Client,
	kongClientFactory ClientFactory,
	readinessReconciliationInterval time.Duration,
) (*AdminAPIClientsManager, error) {
	if len(initialClients) == 0 {
		return nil, errors.New("at least one initial client must be provided")
	}

	return &AdminAPIClientsManager{
		activeGatewayClients:            initialClients,
		adminAPIClientFactory:           kongClientFactory,
		discoveredAdminAPIsNotifyChan:   make(chan []adminapi.DiscoveredAdminAPI),
		ctx:                             ctx,
		notifyLoopRunningCh:             make(chan struct{}),
		logger:                          logger,
		readinessReconciliationInterval: readinessReconciliationInterval,
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
func (c *AdminAPIClientsManager) SetKonnectClient(client *adminapi.KonnectClient) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.konnectClient = client
}

func (c *AdminAPIClientsManager) KonnectClient() *adminapi.KonnectClient {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.konnectClient
}

// GatewayClients returns a copy of current client's slice. Konnect client won't be included.
// This method can be used when some actions need to be performed only against Kong Gateway clients.
func (c *AdminAPIClientsManager) GatewayClients() []*adminapi.Client {
	c.lock.RLock()
	defer c.lock.RUnlock()

	copied := make([]*adminapi.Client, len(c.activeGatewayClients))
	copy(copied, c.activeGatewayClients)
	return copied
}

func (c *AdminAPIClientsManager) GatewayClientsCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.activeGatewayClients)
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

	readinessReconciliationTicker := time.NewTicker(c.readinessReconciliationInterval)
	defer readinessReconciliationTicker.Stop()

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

			// Log the current state of clients.
			c.logClientsLists("notify")

		case <-readinessReconciliationTicker.C:
			if clientsChanged := c.reconcileReadiness(); clientsChanged {
				// Notify subscribers that the clients list has been updated.
				c.notifyGatewayClientsSubscribers()

				// Log the current state of clients.
				c.logClientsLists("readiness reconciliation")
			}
		}
	}
}

// adjustGatewayClients adjusts internally stored clients slice based on the provided
// discovered Admin APIs slice. It consults BaseRootURLs of already stored clients with each
// of the discovered Admin APIs and creates only those clients that we don't have.
func (c *AdminAPIClientsManager) adjustGatewayClients(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Cleanup the pending clients list.
	c.pendingGatewayClients = c.pendingGatewayClients[:0]

	// Short circuit
	if len(discoveredAdminAPIs) == 0 {
		c.activeGatewayClients = c.activeGatewayClients[:0]
		return
	}

	toAdd := lo.Filter(discoveredAdminAPIs, func(api adminapi.DiscoveredAdminAPI, _ int) bool {
		// If we already have a client with a provided address then great, no need
		// to do anything.

		// If we don't have a client with new address then filter it and add
		// a client for this address.
		return !lo.ContainsBy(c.activeGatewayClients, func(cl *adminapi.Client) bool {
			return api.Address == cl.BaseRootURL()
		})
	})

	var idxToRemove []int
	for i, cl := range c.activeGatewayClients {
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
		c.activeGatewayClients = append(c.activeGatewayClients[:idx], c.activeGatewayClients[idx+1:]...)
	}

	for _, adminAPI := range toAdd {
		client, err := c.adminAPIClientFactory.CreateAdminAPIClient(c.ctx, adminAPI)
		if err != nil {
			if errors.As(err, &adminapi.KongClientNotReadyError{}) {
				c.logger.WithError(err).Debugf("client for %q is not ready yet", adminAPI.Address)
			} else {
				c.logger.WithError(err).Errorf("failed to create a client for %q", adminAPI.Address)
			}

			// Despite the error because we still want to keep the client in the pending list to retry later.
			c.pendingGatewayClients = append(c.pendingGatewayClients, adminAPI)
			continue
		}

		// Client is ready, add it to the active list.
		c.activeGatewayClients = append(c.activeGatewayClients, client)
	}
}

// notifyGatewayClientsSubscribers sends notifications to all subscribers that have called SubscribeToGatewayClientsChanges.
func (c *AdminAPIClientsManager) notifyGatewayClientsSubscribers() {
	c.logger.Info("notifying subscribers about the changes in the gateway clients list")
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

// reconcileReadiness checks the readiness of all active and pending clients and moves them to the appropriate list.
func (c *AdminAPIClientsManager) reconcileReadiness() (clientsChanged bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkActiveGatewayClients() || c.checkPendingGatewayClients()
}

func (c *AdminAPIClientsManager) checkPendingGatewayClients() (clientsChanged bool) {
	var idxToRemove []int
	for i, adminAPI := range c.pendingGatewayClients {
		client, err := c.adminAPIClientFactory.CreateAdminAPIClient(c.ctx, adminAPI)
		if err != nil {
			// Despite the error reason we still want to keep the client in the pending list to retry later.
			c.logger.WithError(err).Debugf("pending client for %q is not ready yet", adminAPI.Address)
			continue
		}

		// Client is ready, move it to the active list...
		c.activeGatewayClients = append(c.activeGatewayClients, client)
		// ... and keep the index to remove it from the pending list.
		idxToRemove = append(idxToRemove, i)
	}

	// Remove all clients that are now active.
	for i := len(idxToRemove) - 1; i >= 0; i-- {
		idx := idxToRemove[i]
		c.pendingGatewayClients = append(c.pendingGatewayClients[:idx], c.pendingGatewayClients[idx+1:]...)
	}

	return len(idxToRemove) > 0
}

func (c *AdminAPIClientsManager) checkActiveGatewayClients() (clientsChanged bool) {
	var idxToMoveToPending []int
	for i, client := range c.activeGatewayClients {
		_, err := client.AdminAPIClient().Status(c.ctx)
		if err != nil {
			// Despite the error reason we still want to keep the client in the pending list to retry later.
			c.logger.WithError(err).Debugf("active client for %q is not ready, moving to pending", client.BaseRootURL())
			idxToMoveToPending = append(idxToMoveToPending, i)
		}
	}

	for i := len(idxToMoveToPending) - 1; i >= 0; i-- {
		idx := idxToMoveToPending[i]
		client := c.activeGatewayClients[idx]

		podRef, ok := client.PodReference()
		if !ok {
			// This should never happen, but if it does, we want to log it.
			c.logger.Errorf("failed to get PodReference for client %q", client.BaseRootURL())
		}

		// Add the client to the pending list.
		c.pendingGatewayClients = append(c.pendingGatewayClients, adminapi.DiscoveredAdminAPI{
			Address: client.BaseRootURL(),
			PodRef:  podRef,
		})

		// Remove the client from the active list.
		c.activeGatewayClients = append(c.activeGatewayClients[:idx], c.activeGatewayClients[idx+1:]...)
	}

	return len(idxToMoveToPending) > 0
}

// logClientsLists logs the current state of the active and pending clients, for debugging purposes.
func (c *AdminAPIClientsManager) logClientsLists(afterFunctionName string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	active := lo.Map(c.activeGatewayClients, func(c *adminapi.Client, _ int) string { return c.BaseRootURL() })
	pending := lo.Map(c.pendingGatewayClients, func(c adminapi.DiscoveredAdminAPI, _ int) string { return c.Address })
	c.logger.Debugf("after %s, gateway clients: active: %v, pending: %v", afterFunctionName, active, pending)
}

// AdminAPIClientsProvider allows fetching the most recent list of Admin API clients of Gateways that
// we should configure.
type AdminAPIClientsProvider interface {
	KonnectClient() *adminapi.KonnectClient
	GatewayClients() []*adminapi.Client
}
