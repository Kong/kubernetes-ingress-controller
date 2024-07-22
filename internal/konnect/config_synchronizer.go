package konnect

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

var DefaultConfigUploadPeriod = 30 * time.Second

type ConfigSynchronizer struct {
	logger                 logr.Logger
	syncTicker             *time.Ticker
	kongConfig             sendconfig.Config
	clientsProvider        clients.AdminAPIClientsProvider
	prometheusMetrics      *metrics.CtrlFuncMetrics
	updateStrategyResolver sendconfig.UpdateStrategyResolver
	configChangeDetector   sendconfig.ConfigurationChangeDetector

	targetContent *file.Content

	lock sync.RWMutex
}

func NewConfigSynchronizer(
	logger logr.Logger,
	kongConfig sendconfig.Config,
	clientsProvider clients.AdminAPIClientsProvider,
	updateStrategyResolver sendconfig.UpdateStrategyResolver,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
) *ConfigSynchronizer {
	return &ConfigSynchronizer{
		logger:                 logger,
		syncTicker:             time.NewTicker(DefaultConfigUploadPeriod),
		kongConfig:             kongConfig,
		clientsProvider:        clientsProvider,
		prometheusMetrics:      metrics.NewCtrlFuncMetrics(),
		updateStrategyResolver: updateStrategyResolver,
		configChangeDetector:   configChangeDetector,
	}
}

func (s *ConfigSynchronizer) Start(ctx context.Context) error {
	s.logger.Info("Started Konnect configuration synchronizer")
	go s.runKonnectUpdateServer(ctx)
	return nil
}

func (s *ConfigSynchronizer) SetTargetContent(targetContent *file.Content) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.targetContent = targetContent
}

func (s *ConfigSynchronizer) CopyTargetContent() *file.Content {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.targetContent.DeepCopy()
}

func (s *ConfigSynchronizer) runKonnectUpdateServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context done: shutting down the Konnect update server")
			s.syncTicker.Stop()
		case <-s.syncTicker.C:
			s.logger.Info("Start uploading to Konnect")
			client := s.clientsProvider.KonnectClient()
			if client == nil {
				s.logger.Info("Konnect client not ready, skipping")
				continue
			}
			targetContent := s.CopyTargetContent()
			if targetContent == nil {
				s.logger.Info("No target content received, skipping")
				continue
			}
			err := s.uploadConfig(ctx, client, targetContent)
			if err != nil {
				s.logger.Error(err, "failed to upload configuration to Konnect")
				logKonnectErrors(s.logger, err)
			}
		}
	}
}

func (s *ConfigSynchronizer) uploadConfig(ctx context.Context, client *adminapi.KonnectClient, targetContent *file.Content) error {
	const isFallback = false
	// Remove consumers in target content if consumer sync is disabled.
	if client.ConsumersSyncDisabled() {
		targetContent.Consumers = []file.FConsumer{}
	}

	newSHA, err := sendconfig.PerformUpdate(
		ctx,
		s.logger,
		client,
		s.kongConfig,
		targetContent,
		// Konnect client does not upload custom entities.
		sendconfig.CustomEntitiesByType{},
		s.prometheusMetrics,
		s.updateStrategyResolver,
		s.configChangeDetector,
		isFallback,
	)
	if err != nil {
		return err
	}
	client.SetLastConfigSHA(newSHA)
	return nil
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
