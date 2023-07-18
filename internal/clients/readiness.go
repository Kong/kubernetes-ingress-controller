package clients

import (
	"context"

	"github.com/sirupsen/logrus"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
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
	logger  logrus.FieldLogger
}

func NewDefaultReadinessChecker(factory ClientFactory, logger logrus.FieldLogger) DefaultReadinessChecker {
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
		// We indirectly check readiness of the client by trying to create it. If it succeeds then it means that
		// the client is ready to be used.
		client, err := c.factory.CreateAdminAPIClient(ctx, adminAPI)
		if err != nil {
			// Despite the error reason we still want to keep the client in the pending list to retry later.
			c.logger.WithError(err).Debugf("pending client for %q is not ready yet", adminAPI.Address)
			continue
		}

		turnedReady = append(turnedReady, client)
	}
	return turnedReady
}

// checkAlreadyExistingClients checks if the already existing clients are still ready to be used and returns the ones
// that are not.
func (c DefaultReadinessChecker) checkAlreadyExistingClients(ctx context.Context, lastActive []AlreadyCreatedClient) (turnedPending []adminapi.DiscoveredAdminAPI) {
	for _, client := range lastActive {
		// For ready clients we check readiness by calling the Status endpoint.
		if err := client.IsReady(ctx); err != nil {
			// Despite the error reason we still want to keep the client in the pending list to retry later.
			c.logger.WithError(err).Debugf("active client for %q is not ready, moving to pending", client.BaseRootURL())

			podRef, ok := client.PodReference()
			if !ok {
				// This should never happen, but if it does, we want to log it.
				c.logger.Errorf("failed to get PodReference for client %q", client.BaseRootURL())
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
