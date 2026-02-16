package konnect

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/mo"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/tracing"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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

	syncTicker Ticker

	konnectAdminClient     *adminapi.KonnectClient
	konnectAdminClientLock sync.RWMutex

	targetKongState mo.Option[TargetKongState]
	configLock      sync.RWMutex

	// synchronizerID is the identifier to mark the ConfigSynchronizer instance.
	synchronizerID string
	// serialNumber is the serial number to mark the loop round of config synchronization.
	serialNumber atomic.Uint32
}

// TargetKongState wraps the Kong state to be uploaded to Konnect and indicates whether the configuration is a fallback
// configuration.
type TargetKongState struct {
	*kongstate.KongState

	// IsFallback indicates whether the configuration is a fallback configuration.
	IsFallback bool
}

// TargetContent wraps the deck content to be uploaded to Konnect and indicates whether the configuration is a fallback
// configuration.
type TargetContent struct {
	*file.Content

	// IsFallback indicates whether the configuration is a fallback configuration.
	IsFallback bool
}

type ConfigSynchronizerParams struct {
	Logger                 logr.Logger
	KongConfig             sendconfig.Config
	ConfigUploadTicker     Ticker
	KonnectClientFactory   ClientFactory
	UpdateStrategyResolver sendconfig.UpdateStrategyResolver
	ConfigChangeDetector   sendconfig.ConfigurationChangeDetector
	ConfigStatusNotifier   clients.ConfigStatusNotifier
	MetricsRecorder        metrics.Recorder
	SynchronizerID         string
}

func NewConfigSynchronizer(p ConfigSynchronizerParams) *ConfigSynchronizer {
	if p.SynchronizerID == "" {
		p.SynchronizerID = uuid.NewString()
	}
	return &ConfigSynchronizer{
		logger:                 p.Logger.WithValues("synchronizerID", p.SynchronizerID),
		kongConfig:             p.KongConfig,
		syncTicker:             p.ConfigUploadTicker,
		konnectClientFactory:   p.KonnectClientFactory,
		updateStrategyResolver: p.UpdateStrategyResolver,
		configChangeDetector:   p.ConfigChangeDetector,
		configStatusNotifier:   p.ConfigStatusNotifier,
		metricsRecorder:        p.MetricsRecorder,
		synchronizerID:         p.SynchronizerID,
	}
}

var _ manager.LeaderElectionRunnable = &ConfigSynchronizer{}

// Start starts the loop to receive configuration and upload configuration to Konnect.
func (s *ConfigSynchronizer) Start(ctx context.Context) error {
	s.logger.Info("Starting Konnect configuration synchronizer")
	ctx = context.WithValue(ctx, tracing.SynchronizerIDKey, s.synchronizerID)

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
func (s *ConfigSynchronizer) UpdateKongState(ks *kongstate.KongState, isFallbackConfig bool) {
	// Running the update in a goroutine to not block the caller (i.e. KongClient) as we want to make Konnect updates
	// affect the critical path as little as possible.
	go func() {
		// Update the target configuration to be picked up by the synchronizer loop.
		s.configLock.Lock()
		defer s.configLock.Unlock()
		s.targetKongState = mo.Some(TargetKongState{
			KongState:  ks,
			IsFallback: isFallbackConfig,
		})
	}()
}

// currentContent takes the current KongState (if available) and generates the deck content to be uploaded to Konnect.
// It returns the deck content and a boolean indicating whether the configuration is available or not.
func (s *ConfigSynchronizer) currentContent(ctx context.Context) (TargetContent, bool) {
	// Konnect client may not be initialized yet as that happens asynchronously after the synchronizer is started.
	// UpdateKongState may be called before the initialization completes.
	s.konnectAdminClientLock.RLock()
	defer s.konnectAdminClientLock.RUnlock()
	if s.konnectAdminClient == nil {
		// Konnect client not initialized yet. Cannot generate deck content yet.
		return TargetContent{}, false
	}

	s.configLock.RLock()
	defer s.configLock.RUnlock()
	targetKongState, ok := s.targetKongState.Get()
	if !ok {
		// No configuration received yet.
		return TargetContent{}, false
	}
	ks := targetKongState.KongState

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
	return TargetContent{
		Content:    deckgen.ToDeckContent(ctx, s.logger, ks, deckGenParams),
		IsFallback: targetKongState.IsFallback,
	}, true
}

// KonnectClientInitialized returns true if the Konnect client is initialized and ready to upload configuration.
func (s *ConfigSynchronizer) KonnectClientInitialized() bool {
	s.konnectAdminClientLock.RLock()
	defer s.konnectAdminClientLock.RUnlock()
	return s.konnectAdminClient != nil
}

// run starts the loop uploading the current configuration to Konnect.
func (s *ConfigSynchronizer) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context done: shutting down the Konnect configuration synchronizer")
			s.syncTicker.Stop()
			return
		case <-s.syncTicker.Channel():
			s.handleConfigSynchronizationTick(ctx)
		}
	}
}

func (s *ConfigSynchronizer) handleConfigSynchronizationTick(ctx context.Context) {
	// Add values about the sync round in the context.
	serialNumber := s.serialNumber.Add(1)
	startTimestamp := time.Now().Unix()
	syncRoundID := uuid.NewSHA1(uuid.Nil, fmt.Appendf([]byte{}, "%s:%d:%d", s.synchronizerID, serialNumber, startTimestamp))
	ctx = context.WithValue(ctx, tracing.SyncSerialNumberKey, serialNumber)
	ctx = context.WithValue(ctx, tracing.SyncStartTimestampKey, startTimestamp)
	ctx = context.WithValue(ctx, tracing.SyncRoundIDKey, syncRoundID)
	logger := s.logger.WithValues("syncRoundID", syncRoundID, "serialNumber", serialNumber)

	logger.V(logging.DebugLevel).Info("Start uploading configuration to Konnect")

	// Get the latest configuration copy to upload to Konnect. We don't want to hold the lock for a long time to prevent
	// blocking the update of the configuration.
	targetCfg, ok := s.currentContent(ctx)
	if !ok {
		logger.Info("No configuration received yet, skipping Konnect configuration synchronization")
		return
	}

	// Upload the configuration to Konnect.
	if err := s.uploadConfig(ctx, logger, s.konnectAdminClient, targetCfg); err != nil {
		logger.Error(err, "Failed to upload configuration to Konnect")
		logKonnectErrors(logger, err)
	}
}

// uploadConfig sends the given configuration to Konnect.
func (s *ConfigSynchronizer) uploadConfig(
	ctx context.Context,
	logger logr.Logger,
	client *adminapi.KonnectClient,
	targetContent TargetContent,
) error {
	// Remove consumers in target content if consumer sync is disabled.
	if client.ConsumersSyncDisabled() {
		targetContent.Consumers = []file.FConsumer{}
	}

	newSHA, err := sendconfig.PerformUpdate(
		ctx,
		logger,
		client,
		s.kongConfig,
		targetContent.Content,
		// Konnect client does not upload custom entities.
		sendconfig.CustomEntitiesByType{},
		s.metricsRecorder,
		s.updateStrategyResolver,
		s.configChangeDetector,
		nil,
		targetContent.IsFallback,
	)
	noConfigAcceptedYet := newSHA == nil
	s.configStatusNotifier.NotifyKonnectConfigStatus(ctx, clients.KonnectConfigUploadStatus{
		Failed: err != nil || noConfigAcceptedYet,
	})
	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
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
