package dataplane

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sourcegraph/conc/iter"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

const (
	// KongConfigurationApplySucceededEventReason defines an event reason to tell the updating of Kong configuration succeeded.
	KongConfigurationApplySucceededEventReason = "KongConfigurationSucceeded"
	// KongConfigurationTranslationFailedEventReason defines an event reason used for creating all translation resource failure events.
	KongConfigurationTranslationFailedEventReason = "KongConfigurationTranslationFailed"
	// KongConfigurationApplyFailedEventReason defines an event reason used for creating all config apply resource failure events.
	KongConfigurationApplyFailedEventReason = "KongConfigurationApplyFailed"

	// FallbackKongConfigurationApplySucceededEventReason defines an event reason to tell the updating of fallback Kong configuration succeeded.
	FallbackKongConfigurationApplySucceededEventReason = "FallbackKongConfigurationSucceeded"
	// FallbackKongConfigurationTranslationFailedEventReason defines an event reason used for creating fallback translation resource failure events.
	FallbackKongConfigurationTranslationFailedEventReason = "FallbackKongConfigurationTranslationFailed"
	// FallbackKongConfigurationApplyFailedEventReason defines an event reason used for creating fallback config apply resource failure events.
	FallbackKongConfigurationApplyFailedEventReason = "FallbackKongConfigurationApplyFailed"
)

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Public Types
// -----------------------------------------------------------------------------

// KongConfigBuilder builds a Kong configuration from a Kubernetes object cache.
type KongConfigBuilder interface {
	BuildKongConfig() translator.KongConfigBuildingResult
	UpdateCache(store.CacheStores)
	CustomEntityTypes() []string
}

// FallbackConfigGenerator generates a fallback configuration based on a cache snapshot and a set of broken objects.
type FallbackConfigGenerator interface {
	GenerateExcludingBrokenObjects(
		store.CacheStores,
		[]fallback.ObjectHash,
	) (store.CacheStores, fallback.GeneratedCacheMetadata, error)
	GenerateBackfillingBrokenObjects(
		currentCache store.CacheStores,
		lastValidCache *store.CacheStores,
		brokenObjects []fallback.ObjectHash,
	) (store.CacheStores, fallback.GeneratedCacheMetadata, error)
}

// KongClient is a threadsafe high level API client for the Kong data-plane(s)
// which parses Kubernetes object caches into Kong Admin configurations and
// sends them as updates to the data-plane(s) (Kong Admin API).
type KongClient struct {
	logger logr.Logger

	// requestTimeout is the maximum amount of time that should be waited for
	// requests to the data-plane to receive a response.
	requestTimeout time.Duration

	// cache is the Kubernetes object cache which is used to list Kubernetes
	// objects for parsing into Kong objects.
	cache *store.CacheStores

	// kongConfig is the client configuration for the Kong Admin API
	kongConfig sendconfig.Config

	// dbmode indicates the current database mode of the backend Kong Admin API
	dbmode dpconf.DBMode

	// lock is used to ensure threadsafety of the KongClient object
	lock sync.RWMutex

	// diagnostic is the client and configuration for reporting diagnostic
	// information during data-plane update runtime.
	diagnostic diagnostics.ClientDiagnostic

	// prometheusMetrics is the client for shipping metrics information
	// updates to the prometheus exporter.
	prometheusMetrics *metrics.CtrlFuncMetrics

	// kubernetesObjectReportLock is a mutex for thread-safety of
	// kubernetes object reporting functionality.
	kubernetesObjectReportLock sync.RWMutex

	// kubernetesObjectStatusQueue is a queue that needs to be messaged whenever
	// a Kubernetes object has had configuration for itself successfully applied
	// to the data-plane: messages will trigger reconciliation in the control plane
	// so that status for the objects can be updated accordingly. This is only in
	// use when kubernetesObjectReportsEnabled is true.
	kubernetesObjectStatusQueue *status.Queue

	// kubernetesObjectReportsEnabled indicates whether the data-plane client will
	// file reports about Kubernetes objects which are successfully configured for
	// in the data-plane
	kubernetesObjectReportsEnabled bool

	// kubernetesObjectReportsFilter is a set of objects which were included
	// in the most recent Update(). This can be helpful for callers to determine
	// whether a Kubernetes object has corresponding data-plane configuration that
	// is actively configured (e.g. to know how to set the object status).
	kubernetesObjectReportsFilter k8sobj.ConfigurationStatusSet

	// eventRecorder is used to record warning events for resource failures.
	eventRecorder record.EventRecorder

	// SHAs is a slice is configuration hashes send in last batch send.
	SHAs []string

	// clientsProvider allows retrieving the most recent set of clients.
	clientsProvider clients.AdminAPIClientsProvider

	// configStatusNotifier notifies status of configuring kong gateway.
	configStatusNotifier clients.ConfigStatusNotifier

	// updateStrategyResolver resolves the update strategy for a given Kong Gateway.
	updateStrategyResolver sendconfig.UpdateStrategyResolver

	// configChangeDetector detects changes in the configuration.
	configChangeDetector sendconfig.ConfigurationChangeDetector

	// kongConfigBuilder is used to translate Kubernetes objects into Kong configuration.
	kongConfigBuilder KongConfigBuilder

	// kongConfigFetcher fetches the loaded configuration and status from a Kong node.
	kongConfigFetcher configfetcher.LastValidConfigFetcher

	// controllerPodReference is a reference to the controller pod this client is running in.
	// It may be empty if the client is not running in a pod (e.g. in a unit test).
	controllerPodReference mo.Option[k8stypes.NamespacedName]

	// fallbackConfigGenerator is used to generate a fallback configuration in case of sync failures.
	fallbackConfigGenerator FallbackConfigGenerator

	// lastProcessedSnapshotHash stores the hash of the last processed Kubernetes objects cache snapshot. It's used to determine configuration
	// changes. Please note it is always empty when the `FallbackConfiguration` feature gate is turned off.
	lastProcessedSnapshotHash store.SnapshotHash

	// lastValidCacheSnapshot stores the state of the cache that was last successfully synced with the gateways.
	// Please note it is only populated when the `FallbackConfiguration` feature gate is turned on and the
	// `--use-last-valid-config-for-fallback` flag is set.
	// lastValidCacheSnapshot and lastProcessedSnapshotHash do not always keep values related to the same cache snapshot.
	// While lastProcessedSnapshotHash keeps track of the last processed cache snapshot (the one kept in KongClient.cache),
	// lastValidCacheSnapshot can also represent the fallback cache snapshot that was successfully synced with gateways.
	lastValidCacheSnapshot *store.CacheStores
}

// NewKongClient provides a new KongClient object after connecting to the
// data-plane API and verifying integrity.
func NewKongClient(
	logger logr.Logger,
	timeout time.Duration,
	diagnostic diagnostics.ClientDiagnostic,
	kongConfig sendconfig.Config,
	eventRecorder record.EventRecorder,
	dbMode dpconf.DBMode,
	clientsProvider clients.AdminAPIClientsProvider,
	updateStrategyResolver sendconfig.UpdateStrategyResolver,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
	kongConfigFetcher configfetcher.LastValidConfigFetcher,
	kongConfigBuilder KongConfigBuilder,
	cacheStores *store.CacheStores,
	fallbackConfigGenerator FallbackConfigGenerator,
) (*KongClient, error) {
	c := &KongClient{
		logger:                  logger,
		requestTimeout:          timeout,
		diagnostic:              diagnostic,
		prometheusMetrics:       metrics.NewCtrlFuncMetrics(),
		cache:                   cacheStores,
		kongConfig:              kongConfig,
		eventRecorder:           eventRecorder,
		dbmode:                  dbMode,
		clientsProvider:         clientsProvider,
		configStatusNotifier:    clients.NoOpConfigStatusNotifier{},
		updateStrategyResolver:  updateStrategyResolver,
		configChangeDetector:    configChangeDetector,
		kongConfigBuilder:       kongConfigBuilder,
		kongConfigFetcher:       kongConfigFetcher,
		fallbackConfigGenerator: fallbackConfigGenerator,
	}
	c.initializeControllerPodReference()

	return c, nil
}

func (c *KongClient) initializeControllerPodReference() {
	podNN, err := util.GetPodNN()
	if err != nil {
		c.logger.Error(err, "Failed to resolve controller's pod to attach the apply configuration events to")
		return
	}
	c.controllerPodReference = mo.Some(podNN)
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Public Methods
// -----------------------------------------------------------------------------

// UpdateObject accepts a Kubernetes controller-runtime client.Object and adds/updates that to the configuration cache.
// It will be asynchronously converted into the upstream Kong DSL and applied to the Kong Admin API.
// A status will later be added to the object whether the configuration update succeeds or fails.
func (c *KongClient) UpdateObject(obj client.Object) error {
	// we do a deep copy of the object here so that the caller can continue to use
	// the original object in a threadsafe manner.
	return c.cache.Add(obj.DeepCopyObject())
}

// DeleteObject accepts a Kubernetes controller-runtime client.Object and removes it from the configuration cache.
// The delete action will asynchronously be converted to Kong DSL and applied to the Kong Admin API.
// A status will later be added to the object whether the configuration update succeeds or fails.
//
// under the hood the cache implementation will ignore deletions on objects
// that are not present in the cache, so in those cases this is a no-op.
func (c *KongClient) DeleteObject(obj client.Object) error {
	return c.cache.Delete(obj)
}

// ObjectExists indicates whether or not any version of the provided object is already present in the proxy.
func (c *KongClient) ObjectExists(obj client.Object) (bool, error) {
	_, exists, err := c.cache.Get(obj)
	return exists, err
}

// allEqual returns true if all provided objects are equal.
func allEqual[T any](objs ...T) bool {
	l := len(objs)
	if l == 0 || l == 1 {
		return true
	}

	obj := objs[0]
	for i := 1; i < l; i++ {
		if !reflect.DeepEqual(obj, objs[i]) {
			return false
		}
	}
	return true
}

// Listeners retrieves the currently configured listeners from the underlying
// proxy so that callers can gather this metadata to know which ports
// and protocols are in use by the proxy.
func (c *KongClient) Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error) {
	var (
		errg              errgroup.Group
		errgCollect       errgroup.Group
		listenersCh       = make(chan []kong.ProxyListener)
		listeners         = make([][]kong.ProxyListener, 0)
		streamListenersCh = make(chan []kong.StreamListener)
		streamListeners   = make([][]kong.StreamListener, 0)
	)

	errgCollect.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case l, ok := <-listenersCh:
				if !ok {
					return nil
				}
				listeners = append(listeners, l)
			}
		}
	})
	errgCollect.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sl, ok := <-streamListenersCh:
				if !ok {
					return nil
				}
				streamListeners = append(streamListeners, sl)
			}
		}
	})

	// This lock here (which is shared with .Update()) prevents a data race
	// between reading the client(s) and setting the last applied SHA via client's
	// SetLastConfigSHA() method. It's not ideal but it should do for now.
	c.lock.RLock()
	for _, cl := range c.clientsProvider.GatewayClients() {
		cl := cl
		errg.Go(func() error {
			listeners, streamListeners, err := cl.AdminAPIClient().Listeners(ctx)
			if err != nil {
				return fmt.Errorf("failed to get listeners from %s: %w", cl.BaseRootURL(), err)
			}
			listenersCh <- listeners
			streamListenersCh <- streamListeners

			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		c.lock.RUnlock()
		return nil, nil, err
	}
	c.lock.RUnlock()
	close(listenersCh)
	close(streamListenersCh)
	if err := errgCollect.Wait(); err != nil {
		return nil, nil, err
	}

	if !allEqual(listeners...) {
		return nil, nil, fmt.Errorf("not all listeners out of %d are the same", len(listeners))
	}

	if !allEqual(streamListeners...) {
		return nil, nil, fmt.Errorf("not all stream listeners out of %d are the same", len(streamListeners))
	}

	var (
		retListeners       []kong.ProxyListener
		retStreamListeners []kong.StreamListener
	)
	if len(listeners) > 0 {
		retListeners = listeners[0]
	}
	if len(streamListeners) > 0 {
		retStreamListeners = streamListeners[0]
	}

	return retListeners, retStreamListeners, nil
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Reporting
// -----------------------------------------------------------------------------

// EnableKubernetesObjectReports turns on reporting for Kubernetes objects which are
// configured as part of Update() operations. Enabling this makes it possible to use
// ObjectConfigured(obj) to determine whether an object has successfully been
// configured for on the data-plane.
func (c *KongClient) EnableKubernetesObjectReports(q *status.Queue) {
	c.kubernetesObjectReportLock.Lock()
	defer c.kubernetesObjectReportLock.Unlock()
	c.kubernetesObjectStatusQueue = q
	c.kubernetesObjectReportsEnabled = true
}

// AreKubernetesObjectReportsEnabled returns true or false whether this client has been
// configured to report on Kubernetes objects which have been successfully
// configured for in the data-plane.
func (c *KongClient) AreKubernetesObjectReportsEnabled() bool {
	c.kubernetesObjectReportLock.RLock()
	defer c.kubernetesObjectReportLock.RUnlock()
	return c.kubernetesObjectReportsEnabled
}

// KubernetesObjectIsConfigured reports whether the provided object has active
// configuration for itself successfully applied to the data-plane.
func (c *KongClient) KubernetesObjectIsConfigured(obj client.Object) bool {
	c.kubernetesObjectReportLock.RLock()
	defer c.kubernetesObjectReportLock.RUnlock()
	return c.kubernetesObjectReportsFilter.Get(obj) == k8sobj.ConfigurationStatusSucceeded
}

// KubernetesObjectConfigurationStatus reports the status of applying provided object's
// configuration to the data-plane.
func (c *KongClient) KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus {
	c.kubernetesObjectReportLock.RLock()
	defer c.kubernetesObjectReportLock.RUnlock()
	return c.kubernetesObjectReportsFilter.Get(obj)
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Interface Implementation
// -----------------------------------------------------------------------------

// DBMode indicates which database the Kong Gateway is using.
func (c *KongClient) DBMode() dpconf.DBMode {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.dbmode
}

// Update parses the Cache present in the client and converts current
// Kubernetes state into Kong objects and state, and then ships the
// resulting configuration to the data-plane (Kong Admin API).
func (c *KongClient) Update(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If Kong is running in dbless mode, we can fetch and store the last good configuration.
	if c.dbmode.IsDBLessMode() {
		// Fetch the last valid configuration from the proxy only in case there is no valid
		// configuration already stored in memory. This can happen when KIC restarts and there
		// already is a Kong Proxy with a valid configuration loaded.
		if _, found := c.kongConfigFetcher.LastValidConfig(); !found {
			if err := c.kongConfigFetcher.TryFetchingValidConfigFromGateways(ctx, c.logger, c.clientsProvider.GatewayClients(), c.kongConfigBuilder.CustomEntityTypes()); err != nil {
				// If the client fails to fetch the last good configuration, we log it
				// and carry on, as this is a condition that can be recovered with the following steps.
				c.logger.Error(err, "Failed to fetch last good configuration from gateways")
			}
		}
	}

	// If FallbackConfiguration is enabled, we take a snapshot of the cache so that we operate on a consistent
	// set of resources in case of failures being returned from Kong. As we're going to generate a fallback config
	// based on the cache contents, we need to ensure it is not modified during the process.
	var cacheSnapshot store.CacheStores
	if c.kongConfig.FallbackConfiguration {
		var newSnapshotHash store.SnapshotHash
		var err error
		// Empty snapshot hash means that the cache hasn't changed since the last snapshot was taken. That optimization can be used
		// in main code path to avoid unnecessary processing. TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6095
		cacheSnapshot, newSnapshotHash, err = c.cache.TakeSnapshotIfChanged(c.lastProcessedSnapshotHash)
		if err != nil {
			return fmt.Errorf("failed to take snapshot of cache: %w", err)
		}
		hasNewSnapshotToBeProcessed := newSnapshotHash != store.SnapshotHashEmpty
		if !hasNewSnapshotToBeProcessed {
			c.prometheusMetrics.RecordProcessedConfigSnapshotCacheHit()
		} else {
			c.prometheusMetrics.RecordProcessedConfigSnapshotCacheMiss()
		}
		if hasNewSnapshotToBeProcessed {
			c.logger.V(logging.DebugLevel).Info("New configuration snapshot detected", "hash", newSnapshotHash)
			c.lastProcessedSnapshotHash = newSnapshotHash
			c.kongConfigBuilder.UpdateCache(cacheSnapshot)
		}

		if allGatewaysAreInSync := lo.EveryBy(c.clientsProvider.GatewayClientsToConfigure(), func(cl *adminapi.Client) bool {
			return cl.LastCacheStoresHash() == c.lastProcessedSnapshotHash
		}); allGatewaysAreInSync {
			c.logger.V(logging.DebugLevel).Info("All gateways are in sync; pushing config is not necessary, skipping")
			return nil
		}
	}

	c.logger.V(logging.DebugLevel).Info("Parsing kubernetes objects into data-plane configuration")
	parsingResult := c.kongConfigBuilder.BuildKongConfig()
	if failuresCount := len(parsingResult.TranslationFailures); failuresCount > 0 {
		c.prometheusMetrics.RecordTranslationFailure()
		c.prometheusMetrics.RecordTranslationBrokenResources(failuresCount)
		c.recordResourceFailureEvents(parsingResult.TranslationFailures, KongConfigurationTranslationFailedEventReason)
		c.logger.V(logging.DebugLevel).Info("Translation failures occurred when building data-plane configuration", "count", failuresCount)
	} else {
		c.prometheusMetrics.RecordTranslationSuccess()
		c.prometheusMetrics.RecordTranslationBrokenResources(0)
		c.logger.V(logging.DebugLevel).Info("Successfully built data-plane configuration")
	}

	const isFallback = false
	shas, gatewaysSyncErr := c.sendOutToGatewayClients(ctx, parsingResult.KongState, c.kongConfig, isFallback)
	konnectSyncErr := c.maybeSendOutToKonnectClient(ctx, parsingResult.KongState, c.kongConfig, isFallback)

	// Taking into account the results of syncing configuration with Gateways and Konnect, and potential translation
	// failures, calculate the config status and update it.
	c.configStatusNotifier.NotifyGatewayConfigStatus(ctx, clients.GatewayConfigApplyStatus{
		TranslationFailuresOccurred: len(parsingResult.TranslationFailures) > 0,
		ApplyConfigFailed:           gatewaysSyncErr != nil,
	})
	c.configStatusNotifier.NotifyKonnectConfigStatus(ctx, clients.KonnectConfigUploadStatus{
		Failed: konnectSyncErr != nil,
	})

	// In case of a failure in syncing configuration with Gateways, propagate the error.
	if gatewaysSyncErr != nil {
		if recoveringErr := c.maybeTryRecoveringFromGatewaysSyncError(
			ctx,
			cacheSnapshot,
			gatewaysSyncErr,
		); recoveringErr != nil {
			return fmt.Errorf("failed to recover from gateways sync error: %w", recoveringErr)
		}
		// Update result is positive only if gateways were successfully synced with the current config, so we still
		// need to return the error here even if we succeeded recovering.
		return gatewaysSyncErr
	}

	// Gateways were successfully synced with the current configuration, so we can update the last valid cache snapshot.
	c.maybePreserveTheLastValidConfigCache(cacheSnapshot)

	// report on configured Kubernetes objects if enabled
	if c.AreKubernetesObjectReportsEnabled() {
		// if the configuration SHAs that have just been pushed are different than
		// what's been previously pushed.
		if !slices.Equal(shas, c.SHAs) {
			c.logger.V(logging.DebugLevel).Info("Triggering report for configured Kubernetes objects", "count",
				len(parsingResult.ConfiguredKubernetesObjects))
			c.triggerKubernetesObjectReport(parsingResult.ConfiguredKubernetesObjects, parsingResult.TranslationFailures)
		} else {
			c.logger.V(logging.DebugLevel).Info("No configuration change; resource status update not necessary, skipping")
		}
	}
	return nil
}

// maybePreserveTheLastValidConfigCache preserves the last valid configuration cache if the `FallbackConfiguration`
// feature gate is enabled and the `--enable-last-valid-config-fallback` flag is set.
func (c *KongClient) maybePreserveTheLastValidConfigCache(lastValidCache store.CacheStores) {
	if c.kongConfig.FallbackConfiguration && c.kongConfig.UseLastValidConfigForFallback {
		c.logger.V(logging.DebugLevel).Info("Preserving the last valid configuration cache")
		c.lastValidCacheSnapshot = &lastValidCache
	}
}

// maybeTryRecoveringFromGatewaysSyncError tries to recover from a configuration rejection if the error is of the expected
// UpdateError type. Otherwise, we assume the error is non-recoverable by means of fallback configuration generation
// or applying the last valid configuration, and we need to rely on the retry mechanism. It can happen in case of
// transient network issues, internal server errors, etc.
//
// The recovery can be handled in two ways:
// 1. Generating a fallback configuration and pushing it to the gateways if FallbackConfiguration feature is enabled and
// the UpdateError contains at least 1 broken object.
// 2. Applying the last valid configuration to the gateways if FallbackConfiguration is disabled, fallback
// configuration generation fails, or the UpdateError does not contain any broken objects (it can happen if a gateway
// returns 400 with no meaningful entities' errors).
func (c *KongClient) maybeTryRecoveringFromGatewaysSyncError(
	ctx context.Context,
	cacheSnapshot store.CacheStores,
	gatewaysSyncErr error,
) error {
	// If the error is not of the expected UpdateError type, we should log it and skip the recovery.
	updateErr := sendconfig.UpdateError{}
	if !errors.As(gatewaysSyncErr, &updateErr) {
		c.logger.V(logging.DebugLevel).Info("Skipping recovery from gateways sync error - not enough details to recover",
			"error", gatewaysSyncErr)
		return nil
	}

	// If configuration was rejected by the gateways, we identified at least one broken object and FallbackConfiguration
	// is enabled, we should generate a fallback configuration and push it to the gateways.
	brokenObjects := extractBrokenObjectsFromUpdateError(updateErr)
	if c.kongConfig.FallbackConfiguration && len(brokenObjects) > 0 {
		recoveringErr := c.tryRecoveringWithFallbackConfiguration(ctx, cacheSnapshot, brokenObjects)
		if recoveringErr == nil {
			c.logger.Info("Successfully recovered from configuration rejection with fallback configuration")
			return nil
		}
		// If we failed to recover using the fallback configuration, we should log the error and carry on with the last
		// valid configuration.
		c.logger.Error(recoveringErr, "Failed to recover from configuration rejection with fallback configuration")
	}

	// If FallbackConfiguration is disabled, we skipped or failed to recover using the fallback configuration, we should
	// apply the last valid configuration to the gateways.
	if state, found := c.kongConfigFetcher.LastValidConfig(); found {
		const isFallback = true
		if _, fallbackSyncErr := c.sendOutToGatewayClients(ctx, state, c.kongConfig, isFallback); fallbackSyncErr != nil {
			return errors.Join(gatewaysSyncErr, fallbackSyncErr)
		}
		c.logger.V(logging.DebugLevel).Info("Due to errors in the current config, the last valid config has been pushed to Gateways")
	}
	return nil
}

// tryRecoveringWithFallbackConfiguration tries to recover from a configuration rejection by generating a fallback
// configuration excluding affected objects from the cache.
func (c *KongClient) tryRecoveringWithFallbackConfiguration(
	ctx context.Context,
	currentCache store.CacheStores,
	brokenObjects []fallback.ObjectHash,
) error {
	// Generate a fallback cache snapshot.
	fallbackCache, generatedCacheMetadata, err := c.generateFallbackCache(currentCache, brokenObjects)
	if err != nil {
		return fmt.Errorf("failed to generate fallback configuration: %w", err)
	}
	c.logFallbackCacheMetadata(generatedCacheMetadata)
	if err := c.maybeSendFallbackConfigDiagnostics(ctx, generatedCacheMetadata); err != nil {
		return fmt.Errorf("failed to send fallback configuration diagnostics: %w", err)
	}

	// Update the KongConfigBuilder with the fallback configuration and build the KongConfig.
	c.kongConfigBuilder.UpdateCache(fallbackCache)
	fallbackParsingResult := c.kongConfigBuilder.BuildKongConfig()

	if failuresCount := len(fallbackParsingResult.TranslationFailures); failuresCount > 0 {
		c.recordResourceFailureEvents(fallbackParsingResult.TranslationFailures, FallbackKongConfigurationTranslationFailedEventReason)
		c.prometheusMetrics.RecordFallbackTranslationBrokenResources(failuresCount)
		c.prometheusMetrics.RecordFallbackTranslationFailure()
		c.logger.V(logging.DebugLevel).Info("Translation failures occurred when building fallback data-plane configuration", "count", failuresCount)
	} else {
		c.prometheusMetrics.RecordFallbackTranslationBrokenResources(0)
		c.prometheusMetrics.RecordFallbackTranslationSuccess()
	}

	const isFallback = true
	_, gatewaysSyncErr := c.sendOutToGatewayClients(ctx, fallbackParsingResult.KongState, c.kongConfig, isFallback)
	if gatewaysSyncErr != nil {
		return fmt.Errorf("failed to sync fallback configuration with gateways: %w", gatewaysSyncErr)
	}
	konnectSyncErr := c.maybeSendOutToKonnectClient(ctx, fallbackParsingResult.KongState, c.kongConfig, isFallback)
	if konnectSyncErr != nil {
		// If Konnect sync fails, we should log the error and carry on as it's not a critical error.
		c.logger.Error(konnectSyncErr, "Failed to sync fallback configuration with Konnect")
	}

	// Configuration was successfully recovered with the fallback configuration. Store the last valid configuration.
	c.maybePreserveTheLastValidConfigCache(fallbackCache)
	return nil
}

// generateFallbackCache generates a fallback configuration based on the current cache and a set of broken objects.
// It will either exclude the broken objects from the cache or backfill them from the last valid cache snapshot
// depending on the UseLastValidConfigForFallback flag.
func (c *KongClient) generateFallbackCache(
	currentCache store.CacheStores,
	brokenObjects []fallback.ObjectHash,
) (s store.CacheStores, metadata fallback.GeneratedCacheMetadata, err error) {
	start := time.Now()
	defer func() {
		c.prometheusMetrics.RecordFallbackCacheGenerationDuration(time.Since(start), err)
	}()
	if c.kongConfig.UseLastValidConfigForFallback {
		return c.fallbackConfigGenerator.GenerateBackfillingBrokenObjects(
			currentCache,
			c.lastValidCacheSnapshot,
			brokenObjects,
		)
	}
	return c.fallbackConfigGenerator.GenerateExcludingBrokenObjects(
		currentCache,
		brokenObjects,
	)
}

// extractBrokenObjectsFromUpdateError extracts broken objects from the UpdateError.
func extractBrokenObjectsFromUpdateError(err sendconfig.UpdateError) []fallback.ObjectHash {
	var brokenObjects []client.Object
	for _, resourceFailure := range err.ResourceFailures() {
		brokenObjects = append(brokenObjects, resourceFailure.CausingObjects()...)
	}
	return lo.Map(brokenObjects, func(obj client.Object, _ int) fallback.ObjectHash {
		return fallback.GetObjectHash(obj)
	})
}

// sendOutToGatewayClients will generate deck content (config) from the provided kong state
// and send it out to each of the configured gateway clients.
func (c *KongClient) sendOutToGatewayClients(
	ctx context.Context,
	s *kongstate.KongState,
	config sendconfig.Config,
	isFallback bool,
) ([]string, error) {
	gatewayClients := c.clientsProvider.GatewayClients()
	if len(gatewayClients) == 0 {
		c.logger.Error(
			errors.New("no ready gateway clients"),
			"Could not send configuration to gateways",
		)
		// Should not store the configuration in last valid config because the configuration is not validated on Kong gateway.
		return c.SHAs, nil
	}

	gatewayClientsToConfigure := c.clientsProvider.GatewayClientsToConfigure()
	configureGatewayClientURLs := lo.Map(gatewayClientsToConfigure, func(cl *adminapi.Client, _ int) string { return cl.BaseRootURL() })
	c.logger.V(logging.DebugLevel).Info("Sending configuration to gateway clients", "urls", configureGatewayClientURLs)

	shas, err := iter.MapErr(gatewayClientsToConfigure, func(client **adminapi.Client) (string, error) {
		return c.sendToClient(ctx, *client, s, config, isFallback)
	})
	if err != nil {
		return nil, err
	}

	// After a successful configuration update in DB mode,
	// since only ONE gateway client is chosen to send requests and store SHA of latest configurations,
	// we should propagate the SHA from the chosen client to other clients
	// as well as they will pick the configuration from the shared database.
	if c.dbmode.IsDBBacked() &&
		len(gatewayClients) > 1 {
		for _, client := range gatewayClients {
			client.SetLastConfigSHA([]byte(shas[0]))

			// If the last processed snapshot hash is not empty, we should set it to the clients as well.
			// It can be empty when FallbackConfiguration feature gate is off.
			if c.lastProcessedSnapshotHash != "" {
				client.SetLastCacheStoresHash(c.lastProcessedSnapshotHash)
			}
		}
	}

	previousSHAs := c.SHAs
	sort.Strings(shas)
	c.SHAs = shas

	c.kongConfigFetcher.StoreLastValidConfig(s)

	return previousSHAs, nil
}

// maybeSendOutToKonnectClient sends out the configuration to Konnect when KonnectClient is provided.
// It's a noop when Konnect integration is not enabled.
func (c *KongClient) maybeSendOutToKonnectClient(
	ctx context.Context,
	s *kongstate.KongState,
	config sendconfig.Config,
	isFallback bool,
) error {
	konnectClient := c.clientsProvider.KonnectClient()
	// There's no KonnectClient configured, that's totally fine.
	if konnectClient == nil {
		return nil
	}

	// In case users have many consumers, konnect sync can be very slow and cause dataplane sync issues.
	// For this reason, if the --disable-consumers-sync flag is set, we do not send consumers to Konnect.
	if konnectClient.ConsumersSyncDisabled() {
		s.Consumers = nil
	}

	if _, err := c.sendToClient(ctx, konnectClient, s, config, isFallback); err != nil {
		// In case of an error, we only log it since we don't want the Konnect to affect the basic functionality
		// of the controller.

		if errors.As(err, &sendconfig.UpdateSkippedDueToBackoffStrategyError{}) {
			c.logger.Info("Skipped pushing configuration to Konnect due to backoff strategy", "explanation", err.Error())
		} else {
			c.logger.Error(err, "Failed pushing configuration to Konnect")
			logKonnectErrors(c.logger, err)
		}
		return err
	}

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

func (c *KongClient) sendToClient(
	ctx context.Context,
	client sendconfig.AdminAPIClient,
	s *kongstate.KongState,
	config sendconfig.Config,
	isFallback bool,
) (string, error) {
	logger := c.logger.WithValues("url", client.AdminAPIClient().BaseRootURL())

	// If the client is Konnect and the feature flag is turned on,
	// we should sanitize the configuration before sending it out.
	if client.IsKonnect() && config.SanitizeKonnectConfigDumps {
		s = s.SanitizedCopy(util.DefaultUUIDGenerator{})
	}
	deckGenParams := deckgen.GenerateDeckContentParams{
		SelectorTags:                    config.FilterTags,
		ExpressionRoutes:                config.ExpressionRoutes,
		PluginSchemas:                   client.PluginSchemaStore(),
		AppendStubEntityWhenConfigEmpty: !client.IsKonnect() && config.InMemory,
	}
	targetContent := deckgen.ToDeckContent(ctx, logger, s, deckGenParams)
	customEntities := make(sendconfig.CustomEntitiesByType)
	for entityType, collection := range s.CustomEntities {
		for _, entity := range collection.Entities {
			customEntities[entityType] = append(customEntities[entityType], entity.Object)
		}
	}

	sendDiagnostic := prepareSendDiagnosticFn(ctx, logger, c.diagnostic, s, targetContent, deckGenParams, isFallback)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	newConfigSHA, err := sendconfig.PerformUpdate(
		timedCtx,
		logger,
		client,
		config,
		targetContent,
		customEntities,
		c.prometheusMetrics,
		c.updateStrategyResolver,
		c.configChangeDetector,
		&c.diagnostic,
		isFallback,
	)
	// Only record events on applying configuration to Kong gateway here.
	// Nil error is expected to be passed to indicate success.
	if !client.IsKonnect() {
		c.recordApplyConfigurationEvents(err, client.BaseRootURL(), isFallback)
	}
	if err != nil {
		var (
			rawResponseBody    []byte
			updateErr          sendconfig.UpdateError
			responseParsingErr sendconfig.ResponseParsingError
		)
		if errors.As(err, &updateErr) {
			reason := KongConfigurationApplyFailedEventReason
			if isFallback {
				reason = FallbackKongConfigurationApplyFailedEventReason
			}
			c.recordResourceFailureEvents(updateErr.ResourceFailures(), reason)
		}
		if errors.As(err, &responseParsingErr) {
			rawResponseBody = responseParsingErr.ResponseBody()
		}
		sendDiagnostic(diagnostics.DumpMeta{Failed: true, Hash: string(newConfigSHA)}, rawResponseBody)

		if err := ctx.Err(); err != nil {
			logger.Error(err, "Exceeded Kong API timeout, consider increasing --proxy-timeout-seconds")
		}
		return "", fmt.Errorf("performing update for %s failed: %w", client.BaseRootURL(), err)
	}
	sendDiagnostic(diagnostics.DumpMeta{Failed: false, Hash: string(newConfigSHA)}, nil) // No error occurred.
	// update the lastConfigSHA with the new updated checksum
	client.SetLastConfigSHA(newConfigSHA)
	client.SetLastCacheStoresHash(c.lastProcessedSnapshotHash)
	return string(newConfigSHA), nil
}

// SetConfigStatusNotifier sets a notifier which notifies subscribers about configuration sending results.
// Currently it is used for uploading the node status to konnect control plane.
func (c *KongClient) SetConfigStatusNotifier(n clients.ConfigStatusNotifier) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.configStatusNotifier = n
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Private
// -----------------------------------------------------------------------------

type sendDiagnosticFn func(meta diagnostics.DumpMeta, raw []byte)

// prepareSendDiagnosticFn generates sendDiagnosticFn.
// Diagnostics are sent only when provided diagnostic config (--dump-config) is set.
func prepareSendDiagnosticFn(
	ctx context.Context,
	logger logr.Logger,
	diagnosticConfig diagnostics.ClientDiagnostic,
	targetState *kongstate.KongState,
	targetContent *file.Content,
	deckGenParams deckgen.GenerateDeckContentParams,
	isFallback bool,
) sendDiagnosticFn {
	if diagnosticConfig == (diagnostics.ClientDiagnostic{}) {
		// noop, diagnostics won't be sent
		return func(diagnostics.DumpMeta, []byte) {}
	}

	var config *file.Content
	if diagnosticConfig.DumpsIncludeSensitive {
		config = targetContent
	} else {
		redactedConfig := deckgen.ToDeckContent(ctx,
			logger,
			targetState.SanitizedCopy(util.DefaultUUIDGenerator{}),
			deckGenParams,
		)
		config = redactedConfig
	}

	return func(meta diagnostics.DumpMeta, rawResponseBody []byte) {
		// Given that we can send multiple configs to this channel and
		// the fact that the API that exposes that can only expose 1 config
		// at a time it means that users utilizing the diagnostics API
		// might not see exactly what they intend to see i.e. come failures
		// or successfully send configs might be covered by those send
		// later on but we're OK with this limitation of said API.
		select {
		case diagnosticConfig.Configs <- diagnostics.ConfigDump{
			Meta: diagnostics.DumpMeta{
				Failed:   meta.Failed,
				Fallback: isFallback,
			},
			Config:          *config,
			RawResponseBody: rawResponseBody,
		}:
			logger.V(logging.DebugLevel).Info("Shipping config to diagnostic server")
		default:
			logger.Error(nil, "Config diagnostic buffer full, dropping diagnostic config")
		}
	}
}

// triggerKubernetesObjectReport will update the KongClient with a set which
// enables filtering for which objects are currently applied to the data-plane,
// as well as updating the c.kubernetesObjectStatusQueue to queue those objects
// for reconciliation so their statuses can be properly updated.
func (c *KongClient) triggerKubernetesObjectReport(reportedObjects []client.Object, translationFailures []failures.ResourceFailure) {
	// first a new set of the included objects for the most recent configuration
	// needs to be generated.
	set := k8sobj.ConfigurationStatusSet{}
	for _, obj := range reportedObjects {
		set.Insert(obj, true)
	}

	// in some situations, objects with translation failures are reported:
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3364
	// so we override the failed configuration status from translation failures.
	for _, translationFailure := range translationFailures {
		for _, obj := range translationFailure.CausingObjects() {
			set.Insert(obj, false)
		}
	}

	c.updateKubernetesObjectReportFilter(set)

	// after the filter has been updated we signal the status queue so that the
	// control-plane can update the Kubernetes object statuses for affected objs.
	// this has to be done in a separate loop so that the filter is in place
	// before the objects are enqueued, as the filter is used by the control-plane
	for _, obj := range UniqueObjects(reportedObjects, translationFailures) {
		c.kubernetesObjectStatusQueue.Publish(obj)
	}
}

func UniqueObjects(reportedObjects []client.Object, resourceFailures []failures.ResourceFailure) []client.Object {
	allCausingObjects := lo.FlatMap(resourceFailures, func(f failures.ResourceFailure, _ int) []client.Object {
		return f.CausingObjects()
	})
	allObjects := slices.Concat(reportedObjects, allCausingObjects)
	return lo.UniqBy(allObjects, func(obj client.Object) string {
		return obj.GetObjectKind().GroupVersionKind().String() + "/" +
			obj.GetNamespace() + "/" + obj.GetName()
	})
}

// updateKubernetesObjectReportFilter overrides the internal object set with
// a new provided set.
func (c *KongClient) updateKubernetesObjectReportFilter(set k8sobj.ConfigurationStatusSet) {
	c.kubernetesObjectReportLock.Lock()
	defer c.kubernetesObjectReportLock.Unlock()
	c.kubernetesObjectReportsFilter = set
}

// recordResourceFailureEvents records warning Events for each causing object in each input resource failure, with the
// provided reason.
func (c *KongClient) recordResourceFailureEvents(resourceFailures []failures.ResourceFailure, reason string) {
	for _, failure := range resourceFailures {
		for _, obj := range failure.CausingObjects() {
			gvk := obj.GetObjectKind().GroupVersionKind()
			c.logger.Error(
				errors.New("object failed to apply"),
				"recording a Warning event for object",
				"name", obj.GetName(),
				"namespace", obj.GetNamespace(),
				"kind", gvk.Kind,
				"apiVersion", gvk.Group+"/"+gvk.Version,
				"reason", reason,
				"message", failure.Message(),
			)
			c.eventRecorder.Event(obj, corev1.EventTypeWarning, reason, failure.Message())
		}
	}
}

// recordApplyConfigurationEvents records event attached to KIC pod after KIC applied Kong configuration.
func (c *KongClient) recordApplyConfigurationEvents(err error, rootURL string, isFallback bool) {
	podNN, ok := c.controllerPodReference.Get()
	if !ok {
		// Can't record an event without a controller pod reference to attach to.
		return
	}

	eventType := corev1.EventTypeNormal
	reason := KongConfigurationApplySucceededEventReason
	message := fmt.Sprintf("successfully applied Kong configuration to %s", rootURL)

	if isFallback {
		reason = FallbackKongConfigurationApplySucceededEventReason
		message = fmt.Sprintf("successfully applied fallback Kong configuration to %s", rootURL)
	}

	if err != nil {
		eventType = corev1.EventTypeWarning
		reason = KongConfigurationApplyFailedEventReason
		message = fmt.Sprintf("failed to apply Kong configuration to %s: %v", rootURL, err)

		if isFallback {
			reason = FallbackKongConfigurationApplyFailedEventReason
			message = fmt.Sprintf("failed to apply fallback Kong configuration to %s: %v", rootURL, err)
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podNN.Name,
			Namespace: podNN.Namespace,
		},
	}
	c.eventRecorder.Event(pod, eventType, reason, message)
}

func (c *KongClient) logFallbackCacheMetadata(metadata fallback.GeneratedCacheMetadata) {
	log := c.logger.WithName("fallback-cache-generator")

	// Log excluded objects.
	for _, excluded := range metadata.ExcludedObjects {
		gvk := excluded.Object.GetObjectKind().GroupVersionKind()
		obj := excluded.Object
		causingObjects := lo.Map(excluded.CausingObjects, func(causing fallback.ObjectHash, _ int) string {
			return causing.String()
		})
		log.V(logging.DebugLevel).Info("Excluded object from fallback cache",
			"kind", gvk.Kind,
			"group", gvk.Group,
			"namespace", obj.GetNamespace(),
			"name", obj.GetName(),
			"causing_objects", strings.Join(causingObjects, ","),
		)
	}

	// Log backfilled objects.
	for _, backfilled := range metadata.BackfilledObjects {
		gvk := backfilled.Object.GetObjectKind().GroupVersionKind()
		obj := backfilled.Object
		causingObjects := lo.Map(backfilled.CausingObjects, func(causing fallback.ObjectHash, _ int) string {
			return causing.String()
		})
		log.V(logging.DebugLevel).Info("Backfilled object in fallback cache",
			"kind", gvk.Kind,
			"group", gvk.Group,
			"namespace", obj.GetNamespace(),
			"name", obj.GetName(),
			"causing_objects", strings.Join(causingObjects, ","),
		)
	}
}

func (c *KongClient) maybeSendFallbackConfigDiagnostics(ctx context.Context, generatedCacheMetadata fallback.GeneratedCacheMetadata) error {
	if ch := c.diagnostic.FallbackCacheMetadata; ch != nil {
		select {
		case ch <- generatedCacheMetadata:
			c.logger.V(logging.DebugLevel).Info("Shipping fallback cache metadata to diagnostics server")
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.logger.Error(nil, "Fallback cache metadata buffer full, dropping diagnostics")
		}
	}
	return nil
}
