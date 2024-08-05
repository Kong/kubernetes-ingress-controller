package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

type UpdateStrategyResolver interface {
	ResolveUpdateStrategy(client UpdateClient, diagnostic *diagnostics.ClientDiagnostic) UpdateStrategy
}

type AdminAPIClient interface {
	AdminAPIClient() *kong.Client
	LastConfigSHA() []byte
	SetLastConfigSHA([]byte)
	SetLastCacheStoresHash(store.SnapshotHash)
	BaseRootURL() string
	PluginSchemaStore() *util.PluginSchemaStore

	IsKonnect() bool
	KonnectControlPlane() string
}

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(
	ctx context.Context,
	logger logr.Logger,
	client AdminAPIClient,
	config Config,
	targetContent *file.Content,
	customEntities CustomEntitiesByType,
	promMetrics *metrics.CtrlFuncMetrics,
	updateStrategyResolver UpdateStrategyResolver,
	configChangeDetector ConfigurationChangeDetector,
	diagnostic *diagnostics.ClientDiagnostic,
	isFallback bool,
) ([]byte, error) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent, customEntities)
	if err != nil {
		return oldSHA, fmt.Errorf("failed to generate SHA for target content: %w", err)
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := configChangeDetector.HasConfigurationChanged(ctx, oldSHA, newSHA, targetContent, client, client.AdminAPIClient())
		if err != nil {
			return nil, fmt.Errorf("failed to detect configuration change: %w", err)
		}
		if !configurationChanged {
			if client.IsKonnect() {
				logger.V(logging.DebugLevel).Info("No configuration change, skipping sync to Konnect")
			} else {
				logger.V(logging.DebugLevel).Info("No configuration change, skipping sync to Kong")
			}
			return oldSHA, nil
		}
	}

	updateStrategy := updateStrategyResolver.ResolveUpdateStrategy(client, diagnostic)
	logger = logger.WithValues("update_strategy", updateStrategy.Type())
	timeStart := time.Now()
	err = updateStrategy.Update(ctx, ContentWithHash{
		Content:        targetContent,
		CustomEntities: customEntities,
		Hash:           newSHA,
	})
	duration := time.Since(timeStart)

	metricsProtocol := updateStrategy.MetricsProtocol()
	if err != nil {
		// For UpdateError, record the failure and return the error.
		var updateError UpdateError
		if errors.As(err, &updateError) {
			if isFallback {
				promMetrics.RecordFallbackPushFailure(metricsProtocol, duration, client.BaseRootURL(), len(updateError.ResourceFailures()), updateError.err)
			} else {
				promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), len(updateError.ResourceFailures()), updateError.err)
			}
			return nil, updateError
		}

		// Any other error, simply return it and skip metrics recording - we have no details to record.
		return nil, fmt.Errorf("config update failed: %w", err)
	}

	if isFallback {
		promMetrics.RecordFallbackPushSuccess(metricsProtocol, duration, client.BaseRootURL())
	} else {
		promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())
	}

	if client.IsKonnect() {
		logger.V(logging.InfoLevel).Info("Successfully synced configuration to Konnect", "duration", duration.Truncate(time.Millisecond).String())
	} else {
		logger.V(logging.InfoLevel).Info("Successfully synced configuration to Kong", "duration", duration.Truncate(time.Millisecond).String())
	}

	return newSHA, nil
}
