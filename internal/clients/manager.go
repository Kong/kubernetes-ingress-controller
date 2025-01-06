package clients

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
)

const (
	// DefaultReadinessReconciliationInterval is the interval at which the manager will run readiness reconciliation loop.
	// It's the same as the default interval of a Kubernetes container's readiness probe.
	DefaultReadinessReconciliationInterval = 10 * time.Second
	// MinReadinessReconciliationInterval is the minimum interval of readiness reconciliation loop.
	MinReadinessReconciliationInterval = 3 * time.Second
)

// ClientFactory is responsible for creating Admin API clients.
type ClientFactory interface {
	CreateAdminAPIClient(ctx context.Context, address adminapi.DiscoveredAdminAPI) (*adminapi.Client, error)
}

// AdminAPIClientsProvider allows fetching the most recent list of Admin API clients of Gateways that
// we should configure.
type AdminAPIClientsProvider interface {
	KonnectClient() *adminapi.KonnectClient
	GatewayClients() []*adminapi.Client
	GatewayClientsToConfigure() []*adminapi.Client
}

// Ticker is an interface that allows to control a ticker.
type Ticker interface {
	Stop()
	Channel() <-chan time.Time
	Reset(d time.Duration)
}

// AdminAPIClientsManager keeps track of current Admin API clients of Gateways that we should configure.
// In particular, it can be notified about the discovered clients' list with use of Notify method, and queried
// for the latest slice of ready to be configured clients with use of GatewayClients method. It also runs periodic
// readiness reconciliation loop which is responsible for checking readiness of the clients.
type AdminAPIClientsManager struct {
	// discoveredAdminAPIsNotifyChan is used for notifications that contain Admin API
	// endpoints list that should be used for configuring the dataplane.
	discoveredAdminAPIsNotifyChan    chan []adminapi.DiscoveredAdminAPI
	gatewayClientsChangesSubscribers []chan struct{}

	dbMode dpconf.DBMode

	ctx                   context.Context
	onceNotifyLoopRunning sync.Once
	runningChan           chan struct{}
	isRunning             bool

	// readyGatewayClients represent all Kong Gateway data-planes that are ready to be configured.
	readyGatewayClients map[string]*adminapi.Client

	// pendingGatewayClients represent all Kong Gateway data-planes that were discovered but are not ready to be
	// configured.
	pendingGatewayClients map[string]adminapi.DiscoveredAdminAPI

	readinessReconciliationInterval time.Duration

	// readinessChecker is used to check readiness of the clients.
	readinessChecker ReadinessChecker

	// readinessReconciliationTicker is used to run readiness reconciliation loop.
	readinessReconciliationTicker Ticker

	// konnectClient represents a special-case of the data-plane which is Konnect cloud.
	// This client is used to synchronise configuration with Konnect's Control Plane Admin API.
	konnectClient *adminapi.KonnectClient

	// lock prevents concurrent access to the manager's fields.
	lock sync.RWMutex

	logger logr.Logger
}

type AdminAPIClientsManagerOption func(*AdminAPIClientsManager)

// WithReadinessReconciliationTicker allows to set a custom ticker for readiness reconciliation loop.
func WithReadinessReconciliationTicker(ticker Ticker) AdminAPIClientsManagerOption {
	return func(m *AdminAPIClientsManager) {
		m.readinessReconciliationTicker = ticker
	}
}

// WithDBMode allows to set the DBMode of the Kong gateway instances behind the admin API service.
func (c *AdminAPIClientsManager) WithDBMode(dbMode dpconf.DBMode) *AdminAPIClientsManager {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dbMode = dbMode
	return c
}

// WithReconciliationInterval allows to set the Reconciliation interval to check readiness of clients.
func (c *AdminAPIClientsManager) WithReconciliationInterval(d time.Duration) *AdminAPIClientsManager {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.readinessReconciliationInterval = d
	return c
}

func NewAdminAPIClientsManager(
	ctx context.Context,
	logger logr.Logger,
	initialClients []*adminapi.Client,
	readinessChecker ReadinessChecker,
	opts ...AdminAPIClientsManagerOption,
) (*AdminAPIClientsManager, error) {
	if len(initialClients) == 0 {
		return nil, errors.New("at least one initial client must be provided")
	}

	readyClients := lo.SliceToMap(initialClients, func(c *adminapi.Client) (string, *adminapi.Client) {
		return c.BaseRootURL(), c
	})
	c := &AdminAPIClientsManager{
		readyGatewayClients:             readyClients,
		pendingGatewayClients:           make(map[string]adminapi.DiscoveredAdminAPI),
		readinessReconciliationInterval: DefaultReadinessReconciliationInterval,
		readinessChecker:                readinessChecker,
		readinessReconciliationTicker:   clock.NewTicker(),
		discoveredAdminAPIsNotifyChan:   make(chan []adminapi.DiscoveredAdminAPI),
		ctx:                             ctx,
		runningChan:                     make(chan struct{}),
		logger:                          logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Running returns a channel that is closed when the manager's background tasks are already running.
func (c *AdminAPIClientsManager) Running() chan struct{} {
	return c.runningChan
}

// Run runs a goroutine that will dynamically ingest new addresses of Kong Admin API endpoints.
// It should only be called when Gateway Discovery is enabled.
func (c *AdminAPIClientsManager) Run() {
	c.onceNotifyLoopRunning.Do(func() {
		go c.gatewayClientsReconciliationLoop()

		c.lock.Lock()
		defer c.lock.Unlock()
		c.isRunning = true
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

// SetKonnectClient sets a client that will be used to communicate with Konnect Control Plane Admin API.
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

// GatewayClientsToConfigure returns the gateway clients which need to be configured with the new configuration.
// In DBLess mode, it returns ALL gateway clients
// because we need to update configurations of each gateway instance.
// In DB-backed mode, it returns ONE random gateway client
// because we only need to send configurations to one gateway instance
// while others will be synced using the DB.
func (c *AdminAPIClientsManager) GatewayClientsToConfigure() []*adminapi.Client {
	c.lock.RLock()
	defer c.lock.RUnlock()
	readyGatewayClients := lo.Values(c.readyGatewayClients)
	// With DB-less mode, we should send the configuration to ALL gateway instances.
	if c.dbMode.IsDBLessMode() {
		return readyGatewayClients
	}
	// When a gateway is DB-backed, we return a random client
	// since KIC only needs to send requests to one instance.
	// If there are no ready gateway clients, we return an empty list.
	if len(readyGatewayClients) == 0 {
		return []*adminapi.Client{}
	}
	return readyGatewayClients[:1]
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
	if !c.isRunning {
		return nil, false
	}

	ch := make(chan struct{}, 1)
	c.gatewayClientsChangesSubscribers = append(c.gatewayClientsChangesSubscribers, ch)
	return ch, true
}

// gatewayClientsReconciliationLoop is an inner loop listening on:
// - discoveredAdminAPIsNotifyChan - triggered on every Notify() call.
// - readinessReconciliationTicker - triggered on every readinessReconciliationTicker tick.
func (c *AdminAPIClientsManager) gatewayClientsReconciliationLoop() {
	c.readinessReconciliationTicker.Reset(c.readinessReconciliationInterval)
	defer c.readinessReconciliationTicker.Stop()

	close(c.runningChan)
	for {
		select {
		case <-c.ctx.Done():
			c.logger.V(logging.InfoLevel).Info("Closing AdminAPIClientsManager", "reason", c.ctx.Err())
			c.closeGatewayClientsSubscribers()
			return
		case discoveredAdminAPIs := <-c.discoveredAdminAPIsNotifyChan:
			c.onDiscoveredAdminAPIsNotification(discoveredAdminAPIs)
		case <-c.readinessReconciliationTicker.Channel():
			c.onReadinessReconciliationTick()
		}
	}
}

// onDiscoveredAdminAPIsNotification is called when a new notification about Admin API addresses change is received.
// It will adjust lists of gateway clients and notify subscribers about the change if readyGatewayClients list has
// changed.
func (c *AdminAPIClientsManager) onDiscoveredAdminAPIsNotification(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) {
	c.logger.V(logging.DebugLevel).Info("Received notification about Admin API addresses change")

	clientsChanged := c.adjustGatewayClients(discoveredAdminAPIs)
	readinessChanged := c.reconcileGatewayClientsReadiness()
	if clientsChanged || readinessChanged {
		c.notifyGatewayClientsSubscribers()
	}
}

// onReadinessReconciliationTick is called on every readinessReconciliationTicker tick. It will reconcile readiness
// of all gateway clients and notify subscribers about the change if readyGatewayClients list has changed.
func (c *AdminAPIClientsManager) onReadinessReconciliationTick() {
	c.logger.V(logging.DebugLevel).Info("Reconciling readiness of gateway clients")

	if changed := c.reconcileGatewayClientsReadiness(); changed {
		c.notifyGatewayClientsSubscribers()
	}
}

// adjustGatewayClients adjusts internally stored clients slice based on the provided
// discovered Admin APIs slice. It consults BaseRootURLs of already stored clients with each
// of the discovered Admin APIs and creates only those clients that we don't have.
// It returns true if the gatewayClients slice has been changed, false otherwise.
func (c *AdminAPIClientsManager) adjustGatewayClients(discoveredAdminAPIs []adminapi.DiscoveredAdminAPI) (changed bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Short circuit.
	if len(discoveredAdminAPIs) == 0 {
		// If we have no clients and the provided list is empty, it means we're in sync. No change was made.
		if len(c.readyGatewayClients) == 0 && len(c.pendingGatewayClients) == 0 {
			return false
		}
		// Otherwise, we have to clear the clients and return true to indicate that the change was made.
		clear(c.readyGatewayClients)
		clear(c.pendingGatewayClients)
		return true
	}

	// Make sure all discovered clients that are not in the ready list are in the pending list.
	for _, d := range discoveredAdminAPIs {
		if _, ok := c.readyGatewayClients[d.Address]; !ok {
			c.pendingGatewayClients[d.Address] = d
		}
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

	return changed
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

	c.lock.Lock()
	defer c.lock.Unlock()

	// Short circuit.
	if len(c.readyGatewayClients) == 0 && len(c.pendingGatewayClients) == 0 {
		return false
	}

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
	c.logger.V(logging.DebugLevel).Info("Notifying subscribers about gateway clients change")
	for _, sub := range c.gatewayClientsChangesSubscribers {
		select {
		case <-c.ctx.Done():
			c.logger.V(logging.InfoLevel).Info("Not sending notification to subscribers as the context is done")
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
