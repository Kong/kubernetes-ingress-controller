package clients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

const (
	readinessCheckTimeout = time.Second
)

// ReadinessCheckResult represents the result of a readiness check.
type ReadinessCheckResult struct {
	// ClientsTurnedReady are the clients that were pending and are now ready to be used.
	ClientsTurnedReady []*adminapi.Client
	// ClientsTurnedPending are the clients that were ready and are now pending to be created.
	ClientsTurnedPending []adminapi.DiscoveredAdminAPI
}

// HasChanges returns true if there are any changes in the readiness check result.
// When no changes are present, it means that the readiness check haven't successfully created any pending client
// nor detected any already created client that became not ready.
func (r ReadinessCheckResult) HasChanges() bool {
	return len(r.ClientsTurnedReady) > 0 || len(r.ClientsTurnedPending) > 0
}

// ReadinessChecker is responsible for checking the readiness of Admin API clients.
type ReadinessChecker interface {
	// CheckReadiness checks readiness of the provided clients:
	// - alreadyCreatedClients are the clients that have already been created. The readiness of these clients will be
	//   checked by their IsReady() method.
	// - pendingClients are the clients that have not been created yet and are pending to be created. The readiness of
	//   these clients will be checked by trying to create them.
	CheckReadiness(
		ctx context.Context,
		alreadyCreatedClients []AlreadyCreatedClient,
		pendingClients []adminapi.DiscoveredAdminAPI,
	) ReadinessCheckResult
}

// AlreadyCreatedClient represents an Admin API client that has already been created.
type AlreadyCreatedClient interface {
	IsReady(context.Context) error
	PodReference() (k8stypes.NamespacedName, bool)
	BaseRootURL() string
}

type DefaultReadinessChecker struct {
	factory ClientFactory
	logger  logr.Logger
}

func NewDefaultReadinessChecker(factory ClientFactory, logger logr.Logger) DefaultReadinessChecker {
	return DefaultReadinessChecker{
		factory: factory,
		logger:  logger,
	}
}

func (c DefaultReadinessChecker) CheckReadiness(
	ctx context.Context,
	readyClients []AlreadyCreatedClient,
	pendingClients []adminapi.DiscoveredAdminAPI,
) ReadinessCheckResult {
	return ReadinessCheckResult{
		ClientsTurnedReady:   c.checkPendingGatewayClients(ctx, pendingClients),
		ClientsTurnedPending: c.checkAlreadyExistingClients(ctx, readyClients),
	}
}

// checkPendingGatewayClients checks if the pending clients are ready to be used and returns the ones that are.
func (c DefaultReadinessChecker) checkPendingGatewayClients(ctx context.Context, lastPending []adminapi.DiscoveredAdminAPI) (turnedReady []*adminapi.Client) {
	for _, adminAPI := range lastPending {
		if client := c.checkPendingClient(ctx, adminAPI); client != nil {
			turnedReady = append(turnedReady, client)
		}
	}
	return turnedReady
}

// checkPendingClient indirectly check readiness of the client by trying to create it. If it succeeds then it
// means that the client is ready to be used. It returns a non-nil client if the client is ready to be used, otherwise
// nil is returned.
func (c DefaultReadinessChecker) checkPendingClient(
	ctx context.Context,
	pendingClient adminapi.DiscoveredAdminAPI,
) (client *adminapi.Client) {
	defer func() {
		c.logger.V(logging.DebugLevel).
			Info(fmt.Sprintf("Checking readiness of pending client for %q", pendingClient.Address),
				"ok", client != nil,
			)
	}()

	ctx, cancel := context.WithTimeout(ctx, readinessCheckTimeout)
	defer cancel()
	client, err := c.factory.CreateAdminAPIClient(ctx, pendingClient)
	if err != nil {
		// Despite the error reason we still want to keep the client in the pending list to retry later.
		c.logger.V(logging.DebugLevel).Info("Pending client is not ready yet",
			"reason", err.Error(),
			"address", pendingClient.Address,
		)
		return nil
	}

	return client
}

// checkAlreadyExistingClients checks if the already existing clients are still ready to be used and returns the ones
// that are not.
func (c DefaultReadinessChecker) checkAlreadyExistingClients(ctx context.Context, alreadyCreatedClients []AlreadyCreatedClient) (turnedPending []adminapi.DiscoveredAdminAPI) {
	for _, client := range alreadyCreatedClients {
		// For ready clients we check readiness by calling the Status endpoint.
		if ready := c.checkAlreadyCreatedClient(ctx, client); !ready {
			podRef, ok := client.PodReference()
			if !ok {
				// This should never happen, but if it does, we want to log it.
				c.logger.Error(
					errors.New("missing pod reference"),
					"Failed to get PodReference for client",
					"address", client.BaseRootURL(),
				)
				continue
			}
			turnedPending = append(turnedPending, adminapi.DiscoveredAdminAPI{
				Address: client.BaseRootURL(),
				PodRef:  podRef,
			})
		}
	}
	return turnedPending
}

func (c DefaultReadinessChecker) checkAlreadyCreatedClient(ctx context.Context, client AlreadyCreatedClient) (ready bool) {
	defer func() {
		c.logger.V(logging.DebugLevel).Info(
			fmt.Sprintf("Checking readiness of already created client for %q", client.BaseRootURL()),
			"ok", ready,
		)
	}()

	ctx, cancel := context.WithTimeout(ctx, readinessCheckTimeout)
	defer cancel()
	if err := client.IsReady(ctx); err != nil {
		// Despite the error reason we still want to keep the client in the pending list to retry later.
		c.logger.V(logging.DebugLevel).Info(
			"Already created client is not ready, moving to pending",
			"address", client.BaseRootURL(),
			"reason", err.Error(),
		)
		return false
	}

	return true
}
