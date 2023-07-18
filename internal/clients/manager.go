package clients

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/clock"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// DefaultReadinessReconciliationInterval is the interval at which the manager will run readiness reconciliation loop.
// It's the same as the default interval of the readiness probe.
const DefaultReadinessReconciliationInterval = 10 * time.Second

type ClientFactory interface {
	CreateAdminAPIClient(ctx context.Context, address adminapi.DiscoveredAdminAPI) (*adminapi.Client, error)
}

// AdminAPIClientsProvider allows fetching the most recent list of Admin API clients of Gateways that
// we should configure.
type AdminAPIClientsProvider interface {
	KonnectClient() *adminapi.KonnectClient
	GatewayClients() []*adminapi.Client
}

type Ticker interface {
	Stop()
	Channel() <-chan time.Time
	Reset(d time.Duration)
}

// AdminAPIClientsManager keeps track of current Admin API clients of Gateways that we should configure.
// In particular, it can be notified about the clients' list update with use of Notify method, and queried
// for the latest slice of those with use of Clients method.
type AdminAPIClientsManager struct {
	// discoveredAdminAPIsNotifyChan is used for notifications that contain Admin API
	// endpoints list that should be used for configuring the dataplane.
	discoveredAdminAPIsNotifyChan    chan []adminapi.DiscoveredAdminAPI
	gatewayClientsChangesSubscribers []chan struct{}

	ctx                   context.Context
	onceNotifyLoopRunning sync.Once
	running               chan struct{}
	isNotifyLoopRunning   bool

	// readyGatewayClients represent all Kong Gateway data-planes that are ready to be configured.
	readyGatewayClients map[string]*adminapi.Client

	// pendingGatewayClients represent all Kong Gateway data-planes that were discovered but are not ready to be
	// configured.
	pendingGatewayClients map[string]adminapi.DiscoveredAdminAPI

	// readinessChecker is used to check readiness of the clients.
	readinessChecker ReadinessChecker

	// readinessReconciliationInterval is the interval at which the manager will run readiness reconciliation loop.
	readinessReconciliationInterval time.Duration

	// readinessReconciliationTicker is used to run readiness reconciliation loop.
	readinessReconciliationTicker Ticker

	// konnectClient represents a special-case of the data-plane which is Konnect cloud.
	// This client is used to synchronise configuration with Konnect's Runtime Group Admin API.
	konnectClient *adminapi.KonnectClient

	// lock prevents concurrent access to the manager's fields.
	lock sync.RWMutex

	logger logrus.FieldLogger
}

type AdminAPIClientsManagerOption func(*AdminAPIClientsManager)

// WithReadinessReconciliationTicker allows to set a custom ticker for readiness reconciliation loop.
func WithReadinessReconciliationTicker(ticker Ticker) AdminAPIClientsManagerOption {
	return func(m *AdminAPIClientsManager) {
		m.readinessReconciliationTicker = ticker
	}
}

func NewAdminAPIClientsManager(
	ctx context.Context,
	logger logrus.FieldLogger,
	initialClients []*adminapi.Client,
	kongClientFactory ClientFactory,
	opts ...AdminAPIClientsManagerOption,
) (*AdminAPIClientsManager, error) {
	if len(initialClients) == 0 {
		return nil, errors.New("at least one initial client must be provided")
	}

	readyClients := lo.SliceToMap(initialClients, func(c *adminapi.Client) (string, *adminapi.Client) {
		return c.BaseRootURL(), c
	})
	c := &AdminAPIClientsManager{
		readyGatewayClients:   readyClients,
		pendingGatewayClients: make(map[string]adminapi.DiscoveredAdminAPI),
		// inject
		readinessChecker:                NewDefaultReadinessChecker(kongClientFactory, logger),
		readinessReconciliationTicker:   clock.NewTicker(),
		readinessReconciliationInterval: DefaultReadinessReconciliationInterval,
		discoveredAdminAPIsNotifyChan:   make(chan []adminapi.DiscoveredAdminAPI),
		ctx:                             ctx,
		running:                         make(chan struct{}),
		logger:                          logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Running returns a channel that is closed when the manager's background tasks are already running.
func (c *AdminAPIClientsManager) Running() chan struct{} {
	return c.running
}

// Run runs a goroutine that will dynamically ingest new addresses of Kong Admin API endpoints.
func (c *AdminAPIClientsManager) Run() {
	c.onceNotifyLoopRunning.Do(func() {
		go c.gatewayClientsReconciliationLoop()

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
	return lo.Values(c.readyGatewayClients)
}

func (c *AdminAPIClientsManager) GatewayClientsCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.readyGatewayClients)
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

// gatewayClientsReconciliationLoop is an inner loop listening on notifyChan which are received via
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
func (c *AdminAPIClientsManager) gatewayClientsReconciliationLoop() {
	close(c.running)
	c.readinessReconciliationTicker.Reset(c.readinessReconciliationInterval)
	defer c.readinessReconciliationTicker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("closing AdminAPIClientsManager: %s", c.ctx.Err())
			close(c.discoveredAdminAPIsNotifyChan)
			c.closeGatewayClientsSubscribers()
			return
		case discoveredAdminAPIs := <-c.discoveredAdminAPIsNotifyChan:
			c.onDiscoveredAdminAPIsNotification(discoveredAdminAPIs)
		case <-c.readinessReconciliationTicker.Channel():
			c.onReadinessReconciliationTick()
		}
	}
}

func (c *AdminAPIClientsManager) onDiscoveredAdminAPIsNotification(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.logger.Debug("received notification about Admin API addresses change")
	if clientsChanged := c.adjustGatewayClients(discoveredAdminAPIs); clientsChanged {
		// Notify subscribers that the clients list has been updated.
		c.logger.Debug("notifying subscribers about gateway clients change")
		c.notifyGatewayClientsSubscribers()
	}
}

func (c *AdminAPIClientsManager) onReadinessReconciliationTick() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.logger.Debug("reconciling readiness of gateway clients")
	if clientsChanged := c.reconcileGatewayClientsReadiness(); clientsChanged {
		// Notify subscribers that the clients list has been updated.
		c.logger.Debug("notifying subscribers about gateway clients change")
		c.notifyGatewayClientsSubscribers()
	}
}

// adjustGatewayClients adjusts internally stored clients slice based on the provided
// discovered Admin APIs slice. It consults BaseRootURLs of already stored clients with each
// of the discovered Admin APIs and creates only those clients that we don't have.
// It returns true if the gatewayClients slice has been changed, false otherwise.
func (c *AdminAPIClientsManager) adjustGatewayClients(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) (changed bool) {
	// Short circuit,
	if len(discoveredAdminAPIs) == 0 {
		// If we have no clients and the provided list is empty, it means we're in sync. No change was made.
		if len(c.readyGatewayClients) == 0 && len(c.pendingGatewayClients) == 0 {
			return false
		}
		// Otherwise, we have to clear the clients and return true to indicate that the change was made.
		maps.Clear(c.readyGatewayClients)
		maps.Clear(c.pendingGatewayClients)
		return true
	}

	// Make sure all discovered clients are in the pending list.
	for _, d := range discoveredAdminAPIs {
		c.pendingGatewayClients[d.Address] = d
	}

	// Remove ready clients that are not present in the discovered list.
	for _, cl := range c.readyGatewayClients {
		clientNotOnDiscoveredList := !lo.ContainsBy(discoveredAdminAPIs, func(d adminapi.DiscoveredAdminAPI) bool {
			return d.Address == cl.BaseRootURL()
		})
		if clientNotOnDiscoveredList {
			delete(c.readyGatewayClients, cl.BaseRootURL())
			changed = true
		}
	}

	// Remove pending clients that are not present in the discovered list.
	for _, cl := range c.pendingGatewayClients {
		clientNotOnDiscoveredList := !lo.ContainsBy(discoveredAdminAPIs, func(d adminapi.DiscoveredAdminAPI) bool {
			return d.Address == cl.Address
		})
		if clientNotOnDiscoveredList {
			delete(c.pendingGatewayClients, cl.Address)
			changed = true
		}
	}

	readinessChanged := c.reconcileGatewayClientsReadiness()
	return changed || readinessChanged
}

// reconcileGatewayClientsReadiness reconciles the readiness of the gateway clients. It ensures that the clients on the
// readyGatewayClients list are still ready and that the clients on the pendingGatewayClients list are still pending.
// If any of the clients is not ready anymore, it will be moved to the pendingGatewayClients list. If any of the clients
// is not pending anymore, it will be moved to the readyGatewayClients list. It returns true if any transition has been
// made, false otherwise.
func (c *AdminAPIClientsManager) reconcileGatewayClientsReadiness() bool {
	// Reset the ticker after each readiness reconciliation despite the trigger (whether it was a tick or a notification).
	// It's to ensure that the readiness is not reconciled too often when we receive a lot of notifications.
	defer c.readinessReconciliationTicker.Reset(c.readinessReconciliationInterval)

	readinessCheckResult := c.readinessChecker.CheckReadiness(
		c.ctx,
		lo.MapToSlice(c.readyGatewayClients, func(_ string, cl *adminapi.Client) AlreadyCreatedClient { return cl }),
		lo.Values(c.pendingGatewayClients),
	)

	for _, cl := range readinessCheckResult.ClientsTurnedReady {
		delete(c.pendingGatewayClients, cl.BaseRootURL())
		c.readyGatewayClients[cl.BaseRootURL()] = cl
	}
	for _, cl := range readinessCheckResult.ClientsTurnedPending {
		delete(c.readyGatewayClients, cl.Address)
		c.pendingGatewayClients[cl.Address] = cl
	}

	return readinessCheckResult.HasChanges()
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
