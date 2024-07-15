package dataplane

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

// maybeSendOutToKonnectClient sends out the configuration to Konnect when KonnectClient is provided.
// It's a noop when Konnect integration is not enabled.
func (c *KongClient) maybeSendOutToKonnectClient(
	ctx context.Context,
	s *kongstate.KongState,
	config sendconfig.Config,
	isFallback bool,
) {
	konnectClient := c.clientsProvider.KonnectClient()
	// There's no KonnectClient configured, that's totally fine.
	if konnectClient == nil {
		return
	}

	// send Kong configuration to Konnect client in a new goroutine
	go c.sendOutToKonnectClient(
		ctx, konnectClient, s, config, isFallback,
	)
}

func (c *KongClient) sendOutToKonnectClient(
	ctx context.Context,
	konnectClient *adminapi.KonnectClient,
	s *kongstate.KongState,
	config sendconfig.Config,
	isFallback bool,
) {
	// In case users have many consumers, konnect sync can be very slow and cause dataplane sync issues.
	// For this reason, if the --disable-consumers-sync flag is set, we do not send consumers to Konnect.
	// TODO: modify the implementation here not to modify the shared `Consumers`.
	if konnectClient.ConsumersSyncDisabled() {
		s.Consumers = nil
	}

	if _, err := c.sendToClient(ctx, konnectClient, s, config, isFallback); err != nil {
		// In case of an error, we only log it since we don't want the Konnect to affect the basic functionality
		// of the controller.
		if errors.As(err, &sendconfig.UpdateSkippedDueToBackoffStrategyError{}) {
			c.logger.Info("Skipped pushing configuration to Konnect due to backoff strategy", "explanation", err.Error())
		}

		c.logger.Error(err, "Failed pushing configuration to Konnect")
		logKonnectErrors(c.logger, err)
		if isFallback {
			// If Konnect sync fails, we should log the error and carry on as it's not a critical error.
			c.logger.Error(err, "Failed to sync fallback configuration with Konnect")
		}
		// If we got any error, set Konnect failed to true.
		c.setSendToKonnectFailed(ctx, true)
		return
	}
	// If configuration is correctly sent to Konnect, set Konnect failed to false.
	c.setSendToKonnectFailed(ctx, false)
}

// setSendToKonnectFailed sets the current status of sending configuration to Konnect
// and update the general config status if the status is changed.
func (c *KongClient) setSendToKonnectFailed(ctx context.Context, failed bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.currentConfigStatusAttributes.KonnectFailed != failed {
		c.currentConfigStatusAttributes.KonnectFailed = failed
		c.updateConfigStatus(ctx, clients.CalculateConfigStatus(c.currentConfigStatusAttributes))
	}
}

// logKonnectErrors logs details of each error response returned from Konnect API.
func logKonnectErrors(logger logr.Logger, err error) {
	if crudActionErrors := deckerrors.ExtractCRUDActionErrors(err); len(crudActionErrors) > 0 {
		for _, actionErr := range crudActionErrors {
			apiErr := &kong.APIError{}
			if errors.As(actionErr.Err, &apiErr) {
				logger.Error(actionErr, "Failed to send request to Konnect",
					"operation_type", actionErr.OperationType.String(),
					"entity_kind", actionErr.Kind,
					"entity_name", actionErr.Name,
					"details", apiErr.Details())
			} else {
				logger.Error(actionErr, "Failed to send request to Konnect",
					"operation_type", actionErr.OperationType.String(),
					"entity_kind", actionErr.Kind,
					"entity_name", actionErr.Name,
				)
			}
		}
	}
}
