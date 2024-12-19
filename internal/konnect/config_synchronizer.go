package konnect

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

const (
	// MinConfigUploadPeriod is the minimum period between operations to upload Kong configuration to Konnect.
	MinConfigUploadPeriod = 10 * time.Second
	// DefaultConfigUploadPeriod is the default period between operations to upload Kong configuration to Konnect.
	DefaultConfigUploadPeriod = 30 * time.Second
)

type ClientFactory interface {
	NewKonnectClient(ctx context.Context) (*adminapi.KonnectClient, error)
}

// ConfigSynchronizer runs a loop to upload the translated Kong configuration to Konnect periodically.
type ConfigSynchronizer struct {
	logger               logr.Logger
	kongConfig           sendconfig.Config
	konnectClientFactory ClientFactory

	metricsRecorder        metrics.Recorder
	updateStrategyResolver sendconfig.UpdateStrategyResolver
	configChangeDetector   sendconfig.ConfigurationChangeDetector
	configStatusNotifier   clients.ConfigStatusNotifier

	syncTicker *time.Ticker

	konnectAdminClient     *adminapi.KonnectClient
	konnectAdminClientLock sync.RWMutex

	// targetConfig is the latest configuration to be uploaded to Konnect.
	targetConfig targetConfig
	// configLock is used to prevent data
	configLock sync.RWMutex
}

type targetConfig struct {
	// Content the configuration to be uploaded to Konnect. It represents the latest state of the configuration
	// received from the KongClient.
	Content *file.Content

	// IsFallback indicates whether the configuration is a fallback configuration.
	IsFallback bool
}

type ConfigSynchronizerParams struct {
	Logger                 logr.Logger
	KongConfig             sendconfig.Config
	ConfigUploadPeriod     time.Duration
	KonnectClientFactory   ClientFactory
	UpdateStrategyResolver sendconfig.UpdateStrategyResolver
	ConfigChangeDetector   sendconfig.ConfigurationChangeDetector
	ConfigStatusNotifier   clients.ConfigStatusNotifier
	MetricsRecorder        metrics.Recorder
}

func NewConfigSynchronizer(p ConfigSynchronizerParams) *ConfigSynchronizer {
	return &ConfigSynchronizer{
		logger:                 p.Logger,
		kongConfig:             p.KongConfig,
		syncTicker:             time.NewTicker(p.ConfigUploadPeriod),
		konnectClientFactory:   p.KonnectClientFactory,
		updateStrategyResolver: p.UpdateStrategyResolver,
		configChangeDetector:   p.ConfigChangeDetector,
		configStatusNotifier:   p.ConfigStatusNotifier,
		metricsRecorder:        p.MetricsRecorder,
	}
}

var _ manager.LeaderElectionRunnable = &ConfigSynchronizer{}

// Start starts the loop to receive configuration and uplaod configuration to Konnect.
func (s *ConfigSynchronizer) Start(ctx context.Context) error {
	s.logger.Info("Starting Konnect configuration synchronizer")

	konnectAdminClient, err := s.konnectClientFactory.NewKonnectClient(ctx)
	if err != nil {
		s.logger.Error(err, "Failed to create Konnect client, skipping Konnect configuration synchronization")

		// We failed to set up Konnect client. We cannot proceed with running the synchronizer.
		// As it's a manager runnable, we'll wait for the context to be done and return only then to not break the
		// manager's start process.
		<-ctx.Done()
		return ctx.Err()
	}

	// Set the Konnect client to be used to upload configuration and start the synchronizer main loop.
	s.konnectAdminClientLock.Lock()
	s.konnectAdminClient = konnectAdminClient
	s.konnectAdminClientLock.Unlock()
	s.logger.Info("Konnect client initialized, starting Konnect configuration synchronization")
	s.run(ctx)

	return nil
}

// NeedLeaderElection returns true to indicate that this runnable requires leader election.
// This is required to ensure that only one instance of the synchronizer is running at a time.
func (s *ConfigSynchronizer) NeedLeaderElection() bool {
	return true
}

// UpdateKongState updates the Kong state to be uploaded to Konnect asynchronously. It may not update the state if
// the Konnect client is not initialized yet.
func (s *ConfigSynchronizer) UpdateKongState(
	ctx context.Context,
	ks *kongstate.KongState,
	isFallbackConfig bool,
) {
	// Running the update in a goroutine to not block the caller (i.e. KongClient) as we want to make Konnect updates
	// affect the critical path as little as possible.
	go func() {
		// Konnect client may not be initialized yet as that happens asynchronously after the synchronizer is started.
		// UpdateKongState may be called before the initialization completes.
		s.konnectAdminClientLock.RLock()
		defer s.konnectAdminClientLock.RUnlock()
		if s.konnectAdminClient == nil {
			s.logger.Info("Konnect client not initialized yet, skipping Kong state update")
			return
		}

		// Sanitize the configuration dumps if configured to do so.
		if s.kongConfig.SanitizeKonnectConfigDumps {
			ks = ks.SanitizedCopy(util.DefaultUUIDGenerator{})
		}

		// Generate the deck content to be uploaded to Konnect. It may issue some API calls to Konnect to get additional
		// information like plugin schemas.
		deckGenParams := deckgen.GenerateDeckContentParams{
			SelectorTags:     s.kongConfig.FilterTags,
			ExpressionRoutes: s.kongConfig.ExpressionRoutes,
			PluginSchemas:    s.konnectAdminClient.PluginSchemaStore(),
		}
		targetContent := deckgen.ToDeckContent(ctx, s.logger, ks, deckGenParams)

		// Update the target configuration to be picked up by the synchronizer loop.
		s.configLock.Lock()
		defer s.configLock.Unlock()
		s.targetConfig = targetConfig{
			Content:    targetContent,
			IsFallback: isFallbackConfig,
		}
	}()
}

// run starts the loop uploading the current configuration to Konnect.
func (s *ConfigSynchronizer) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context done: shutting down the Konnect configuration synchronizer")
			s.syncTicker.Stop()
		case <-s.syncTicker.C:
			s.handleConfigSynchronizationTick(ctx)
		}
	}
}

func (s *ConfigSynchronizer) handleConfigSynchronizationTick(ctx context.Context) {
	s.logger.V(logging.DebugLevel).Info("Start uploading configuration to Konnect")

	// Get the latest configuration copy to upload to Konnect. We don't want to hold the lock for a long time to prevent
	// blocking the update of the configuration.
	targetCfg := s.getTargetConfigCopy()
	if targetCfg.Content == nil {
		s.logger.Info("No configuration received yet, skipping Konnect configuration synchronization")
		return
	}

	// Upload the configuration to Konnect.
	err := s.uploadConfig(ctx, s.konnectAdminClient, targetCfg)
	if err != nil {
		s.logger.Error(err, "Failed to upload configuration to Konnect")
		logKonnectErrors(s.logger, err)
	}

	// Notify the status of the configuration upload to the system reporting it.
	s.configStatusNotifier.NotifyKonnectConfigStatus(ctx, clients.KonnectConfigUploadStatus{
		Failed: err != nil,
	})
}

// uploadConfig sends the given configuration to Konnect.
func (s *ConfigSynchronizer) uploadConfig(
	ctx context.Context,
	client *adminapi.KonnectClient,
	targetCfg targetConfig,
) error {
	// Remove consumers in target content if consumer sync is disabled.
	if client.ConsumersSyncDisabled() {
		targetCfg.Content.Consumers = []file.FConsumer{}
	}

	newSHA, err := sendconfig.PerformUpdate(
		ctx,
		s.logger,
		client,
		s.kongConfig,
		targetCfg.Content,
		// Konnect client does not upload custom entities.
		sendconfig.CustomEntitiesByType{},
		s.metricsRecorder,
		s.updateStrategyResolver,
		s.configChangeDetector,
		nil,
		targetCfg.IsFallback,
	)
	if err != nil {
		return err
	}
	client.SetLastConfigSHA(newSHA)
	return nil
}

// getTargetConfigCopy returns a copy of the latest configuration in `file.Content` format
// to prevent data race and long duration of occupying lock.
func (s *ConfigSynchronizer) getTargetConfigCopy() targetConfig {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return targetConfig{
		s.targetConfig.Content.DeepCopy(),
		s.targetConfig.IsFallback,
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
