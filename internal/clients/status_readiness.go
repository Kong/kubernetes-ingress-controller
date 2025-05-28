package clients

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

// StatusReadinessCheckResult represents the result of a status readiness check.
type StatusReadinessCheckResult struct {
	// ClientsTurnedReady are the status clients that were pending and are now ready to be used.
	ClientsTurnedReady []*adminapi.StatusClient
	// ClientsTurnedPending are the status clients that were ready and are now pending to be created.
	ClientsTurnedPending []adminapi.DiscoveredAdminAPI
}

// HasChanges returns true if there are any changes in the status readiness check result.
func (r StatusReadinessCheckResult) HasChanges() bool {
	return len(r.ClientsTurnedReady) > 0 || len(r.ClientsTurnedPending) > 0
}

// StatusReadinessChecker is responsible for checking the readiness of Status API clients.
type StatusReadinessChecker interface {
	// CheckStatusReadiness checks readiness of the provided status clients.
	CheckStatusReadiness(
		ctx context.Context,
		alreadyCreatedClients []AlreadyCreatedStatusClient,
		pendingClients []adminapi.DiscoveredAdminAPI,
	) StatusReadinessCheckResult
}

// AlreadyCreatedStatusClient represents a Status API client that has already been created.
type AlreadyCreatedStatusClient interface {
	IsReady(context.Context) error
	PodReference() (k8stypes.NamespacedName, bool)
	BaseRootURL() string
}

// StatusClientFactory interface for creating status clients.
type StatusClientFactory interface {
	CreateStatusClient(ctx context.Context, discoveredStatusAPI adminapi.DiscoveredAdminAPI) (*adminapi.StatusClient, error)
}

// DefaultStatusReadinessChecker implements StatusReadinessChecker.
type DefaultStatusReadinessChecker struct {
	factory               StatusClientFactory
	readinessCheckTimeout time.Duration
	logger                logr.Logger
}

// NewDefaultStatusReadinessChecker creates a new DefaultStatusReadinessChecker.
func NewDefaultStatusReadinessChecker(factory StatusClientFactory, timeout time.Duration, logger logr.Logger) DefaultStatusReadinessChecker {
	return DefaultStatusReadinessChecker{
		factory:               factory,
		readinessCheckTimeout: timeout,
		logger:                logger,
	}
}

// CheckStatusReadiness checks readiness of status clients.
func (c DefaultStatusReadinessChecker) CheckStatusReadiness(
	ctx context.Context,
	readyClients []AlreadyCreatedStatusClient,
	pendingClients []adminapi.DiscoveredAdminAPI,
) StatusReadinessCheckResult {
	var (
		turnedReadyCh   = make(chan []*adminapi.StatusClient)
		turnedPendingCh = make(chan []adminapi.DiscoveredAdminAPI)
	)

	go func(ctx context.Context, pendingClients []adminapi.DiscoveredAdminAPI) {
		turnedReadyCh <- c.checkPendingStatusClients(ctx, pendingClients)
		close(turnedReadyCh)
	}(ctx, pendingClients)

	go func(ctx context.Context, readyClients []AlreadyCreatedStatusClient) {
		turnedPendingCh <- c.checkAlreadyExistingStatusClients(ctx, readyClients)
		close(turnedPendingCh)
	}(ctx, readyClients)

	return StatusReadinessCheckResult{
		ClientsTurnedReady:   <-turnedReadyCh,
		ClientsTurnedPending: <-turnedPendingCh,
	}
}

// checkPendingStatusClients checks if the pending status clients are ready to be used.
func (c DefaultStatusReadinessChecker) checkPendingStatusClients(ctx context.Context, lastPending []adminapi.DiscoveredAdminAPI) (turnedReady []*adminapi.StatusClient) {
	var (
		wg sync.WaitGroup
		ch = make(chan *adminapi.StatusClient)
	)
	for _, statusAPI := range lastPending {
		wg.Add(1)
		go func(statusAPI adminapi.DiscoveredAdminAPI) {
			defer wg.Done()
			if client := c.checkPendingStatusClient(ctx, statusAPI); client != nil {
				select {
				case ch <- client:
				case <-ctx.Done():
				}
			}
		}(statusAPI)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for client := range ch {
		turnedReady = append(turnedReady, client)
	}

	return turnedReady
}

// checkPendingStatusClient checks readiness of a pending status client by trying to create it.
func (c DefaultStatusReadinessChecker) checkPendingStatusClient(
	ctx context.Context,
	pendingClient adminapi.DiscoveredAdminAPI,
) (client *adminapi.StatusClient) {
	ctx, cancel := context.WithTimeout(ctx, c.readinessCheckTimeout)
	defer cancel()

	logger := c.logger.WithValues("address", pendingClient.Address)

	client, err := c.factory.CreateStatusClient(ctx, pendingClient)
	if err != nil {
		logger.V(logging.DebugLevel).Info(
			"Pending status client is not ready yet",
			"reason", err.Error(),
		)
		return nil
	}

	logger.V(logging.DebugLevel).Info(
		"Checked readiness of pending status client",
		"ok", client != nil,
	)

	return client
}

// checkAlreadyExistingStatusClients checks if already existing status clients are still ready.
func (c DefaultStatusReadinessChecker) checkAlreadyExistingStatusClients(ctx context.Context, alreadyCreatedClients []AlreadyCreatedStatusClient) (turnedPending []adminapi.DiscoveredAdminAPI) {
	var (
		wg          sync.WaitGroup
		pendingChan = make(chan adminapi.DiscoveredAdminAPI)
	)

	for _, client := range alreadyCreatedClients {
		wg.Add(1)
		go func(client AlreadyCreatedStatusClient) {
			defer wg.Done()

			if ready := c.checkAlreadyCreatedStatusClient(ctx, client); !ready {
				podRef, ok := client.PodReference()
				if !ok {
					c.logger.Error(
						nil,
						"Failed to get PodReference for status client",
						"address", client.BaseRootURL(),
					)
					return
				}
				select {
				case <-ctx.Done():
				case pendingChan <- adminapi.DiscoveredAdminAPI{
					Address: client.BaseRootURL(),
					PodRef:  podRef,
				}:
				}
			}
		}(client)
	}

	go func() {
		wg.Wait()
		close(pendingChan)
	}()

	for pendingClient := range pendingChan {
		turnedPending = append(turnedPending, pendingClient)
	}

	return turnedPending
}

// checkAlreadyCreatedStatusClient checks if an already created status client is still ready.
func (c DefaultStatusReadinessChecker) checkAlreadyCreatedStatusClient(ctx context.Context, client AlreadyCreatedStatusClient) (ready bool) {
	logger := c.logger.WithValues("address", client.BaseRootURL())

	ctx, cancel := context.WithTimeout(ctx, c.readinessCheckTimeout)
	defer cancel()
	if err := client.IsReady(ctx); err != nil {
		logger.V(logging.DebugLevel).Info(
			"Already created status client is not ready, moving to pending",
			"reason", err.Error(),
		)
		return false
	}

	logger.V(logging.DebugLevel).Info("Already created status client is ready")

	return true
}