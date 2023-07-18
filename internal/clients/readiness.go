package clients

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/sirupsen/logrus"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

type ReadinessCheckResult struct {
	ClientsTurnedReady   []*adminapi.Client
	ClientsTurnedPending []adminapi.DiscoveredAdminAPI
}

func (r ReadinessCheckResult) HasChanges() bool {
	return len(r.ClientsTurnedReady) > 0 || len(r.ClientsTurnedPending) > 0
}

type ReadinessChecker interface {
	CheckReadiness(
		ctx context.Context,
		alreadyCreatedClients []AlreadyCreatedClient,
		pendingClients []adminapi.DiscoveredAdminAPI,
	) ReadinessCheckResult
}

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
		ClientsTurnedPending: c.checkActiveGatewayClients(ctx, readyClients),
	}
}

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

func (c DefaultReadinessChecker) checkActiveGatewayClients(ctx context.Context, lastActive []AlreadyCreatedClient) (turnedPending []adminapi.DiscoveredAdminAPI) {
	for _, client := range lastActive {
		// For active clients we check readiness by calling the Status endpoint.
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
