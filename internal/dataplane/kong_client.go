package dataplane

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/iter"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

const (
	// KongConfigurationTranslationFailedEventReason defines an event reason used for creating all translation resource failure events.
	KongConfigurationTranslationFailedEventReason = "KongConfigurationTranslationFailed"
	// KongConfigurationApplyFailedEventReason defines an event reason used for creating all config apply resource failure events.
	KongConfigurationApplyFailedEventReason = "KongConfigurationApplyFailed"
)

// AdminAPIClientsProvider allows fetching the most recent list of Admin API clients of Gateways that
// we should configure.
type AdminAPIClientsProvider interface {
	Clients() []*adminapi.Client
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Public Types
// -----------------------------------------------------------------------------

// KongClient is a threadsafe high level API client for the Kong data-plane(s)
// which parses Kubernetes object caches into Kong Admin configurations and
// sends them as updates to the data-plane(s) (Kong Admin API).
type KongClient struct {
	logger logrus.FieldLogger

	// ingressClass indicates the Kubernetes ingress class that should be
	// used to qualify support for any given Kubernetes object to be parsed
	// into data-plane configuration.
	ingressClass string

	// enableReverseSync indicates that reverse sync should be enabled for
	// updates to the data-plane.
	enableReverseSync bool

	// enableCombinedServiceRoutes indicates that when translating Kubernetes
	// ingress objects into Kong Admin API configuration we should disable the
	// legacy logic which would create a single route per path and instead use
	// the newer logic which combines them.
	enableCombinedServiceRoutes bool

	// skipCACertificates disables CA certificates, to avoid fighting over configuration in multi-workspace
	// environments. See https://github.com/Kong/deck/pull/617
	skipCACertificates bool

	// requestTimeout is the maximum amount of time that should be waited for
	// requests to the data-plane to receive a response.
	requestTimeout time.Duration

	// cache is the Kubernetes object cache which is used to list Kubernetes
	// objects for parsing into Kong objects.
	cache *store.CacheStores

	// kongConfig is the client configuration for the Kong Admin API
	kongConfig sendconfig.Config

	// dbmode indicates the current database mode of the backend Kong Admin API
	dbmode string

	// lock is used to ensure threadsafety of the KongClient object
	lock sync.RWMutex

	// diagnostic is the client and configuration for reporting diagnostic
	// information during data-plane update runtime.
	diagnostic util.ConfigDumpDiagnostic

	// prometheusMetrics is the client for shipping metrics information
	// updates to the prometheus exporter.
	prometheusMetrics *metrics.CtrlFuncMetrics

	// kubernetesObjectReportLock is a mutex for thread-safety of
	// kubernetes object reporting functionality.
	kubernetesObjectReportLock sync.RWMutex

	// additionalFeaturesLock is a mutex to enable thread-safety of enabling or
	// disabling various features.
	additionalFeaturesLock sync.RWMutex

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
	clientsProvider AdminAPIClientsProvider
}

// NewKongClient provides a new KongClient object after connecting to the
// data-plane API and verifying integrity.
func NewKongClient(
	logger logrus.FieldLogger,
	timeout time.Duration,
	ingressClass string,
	enableReverseSync bool,
	skipCACertificates bool,
	diagnostic util.ConfigDumpDiagnostic,
	kongConfig sendconfig.Config,
	eventRecorder record.EventRecorder,
	dbMode string,
	clientsProvider AdminAPIClientsProvider,
) (*KongClient, error) {
	// build the client object
	cache := store.NewCacheStores()
	c := &KongClient{
		logger:             logger,
		ingressClass:       ingressClass,
		enableReverseSync:  enableReverseSync,
		skipCACertificates: skipCACertificates,
		requestTimeout:     timeout,
		diagnostic:         diagnostic,
		prometheusMetrics:  metrics.NewCtrlFuncMetrics(),
		cache:              &cache,
		kongConfig:         kongConfig,
		eventRecorder:      eventRecorder,
		dbmode:             dbMode,
		clientsProvider:    clientsProvider,
	}

	return c, nil
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
	for _, cl := range c.clientsProvider.Clients() {
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

func (c *KongClient) KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus {
	c.kubernetesObjectReportLock.RLock()
	defer c.kubernetesObjectReportLock.RUnlock()
	return c.kubernetesObjectReportsFilter.Get(obj)
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Optional Features
// -----------------------------------------------------------------------------

// EnableCombinedServiceRoutes turns on the combined service routes feature for
// the Kong Dataplane client.
func (c *KongClient) EnableCombinedServiceRoutes() {
	c.additionalFeaturesLock.Lock()
	defer c.additionalFeaturesLock.Unlock()
	c.enableCombinedServiceRoutes = true
}

// AreCombinedServiceRoutesEnabled determines whether the combined service
// routes translation mode has been enabled, or if the legacy logic is being
// used. When enabled this changes the logic to try and combine multiple paths
// into single routes, but it also changes the names of existing routes and so
// it should be considered disruptive as it will temporarily drop routes when
// it's first enabled.
func (c *KongClient) AreCombinedServiceRoutesEnabled() bool {
	c.additionalFeaturesLock.RLock()
	defer c.additionalFeaturesLock.RUnlock()
	return c.enableCombinedServiceRoutes
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Interface Implementation
// -----------------------------------------------------------------------------

// DBMode indicates which database the Kong Gateway is using.
func (c *KongClient) DBMode() string {
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

	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*c.cache, c.ingressClass, c.logger)

	// initialize a parser
	c.logger.Debug("parsing kubernetes objects into data-plane configuration")

	p, err := parser.NewParser(c.logger, storer)
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	if c.AreKubernetesObjectReportsEnabled() {
		p.EnableKubernetesObjectReports()
	}
	if c.AreCombinedServiceRoutesEnabled() {
		p.EnableCombinedServiceRoutes()
	}
	formatVersion := "1.1"
	if versions.GetKongVersion().MajorMinorOnly().GTE(versions.ExplicitRegexPathVersionCutoff) {
		p.EnableRegexPathPrefix()
		formatVersion = "3.0"
	}

	// parse the Kubernetes objects from the storer into Kong configuration
	kongstate, translationFailures := p.Build()
	if failuresCount := len(translationFailures); failuresCount > 0 {
		c.prometheusMetrics.RecordTranslationFailure()
		c.recordResourceFailureEvents(translationFailures, KongConfigurationTranslationFailedEventReason)
		c.logger.Debugf("%d translation failures have occurred when building data-plane configuration", failuresCount)
	} else {
		c.prometheusMetrics.RecordTranslationSuccess()
		c.logger.Debug("successfully built data-plane configuration")
	}

	shas, err := c.sendOutToClients(ctx, kongstate, formatVersion, c.kongConfig)
	if err != nil {
		return err
	}

	// report on configured Kubernetes objects if enabled
	if c.AreKubernetesObjectReportsEnabled() {
		// if the configuration SHAs that have just been pushed are different than
		// what's been previously pushed.
		if !slices.Equal(shas, c.SHAs) {
			report := p.GenerateKubernetesObjectReport()
			c.logger.Debugf("triggering report for %d configured Kubernetes objects", len(report))
			c.triggerKubernetesObjectReport(report, translationFailures)
		} else {
			c.logger.Debug("no configuration change, skipping kubernetes object report")
		}
	}
	return nil
}

// sendOutToClients will generate deck content (config) from the provided kong state
// and send it out to each of the configured clients.
func (c *KongClient) sendOutToClients(
	ctx context.Context, s *kongstate.KongState, formatVersion string, config sendconfig.Config,
) ([]string, error) {
	clients := c.clientsProvider.Clients()
	c.logger.Debugf("sending configuration to %d clients", len(clients))
	shas, err := iter.MapErr(clients, func(client **adminapi.Client) (string, error) {
		return c.sendToClient(ctx, *client, s, formatVersion, config)
	},
	)
	if err != nil {
		return nil, err
	}
	previousSHAs := c.SHAs

	sort.Strings(shas)
	c.SHAs = shas

	return previousSHAs, nil
}

func (c *KongClient) sendToClient(
	ctx context.Context,
	client *adminapi.Client,
	s *kongstate.KongState,
	formatVersion string,
	config sendconfig.Config,
) (string, error) {
	logger := c.logger.WithField("kong_url", client.AdminAPIClient().BaseRootURL())

	// generate the deck configuration to be applied to the admin API
	logger.Debug("converting configuration to deck config")
	targetConfig := deckgen.ToDeckContent(ctx,
		logger,
		s,
		client.PluginSchemaStore(),
		config.FilterTags,
		formatVersion,
	)

	sendDiagnostic := prepareSendDiagnosticFn(ctx, logger, c.diagnostic, s, targetConfig, client.PluginSchemaStore(), config.FilterTags, formatVersion)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	newConfigSHA, err := sendconfig.PerformUpdate(
		timedCtx,
		logger,
		client,
		config,
		targetConfig,
		c.prometheusMetrics,
	)
	sendDiagnostic(err != nil)
	if err != nil {
		if expired, ok := timedCtx.Deadline(); ok && time.Now().After(expired) {
			logger.Warn("exceeded Kong API timeout, consider increasing --proxy-timeout-seconds")
		}
		return "", fmt.Errorf("performing update for %s failed: %w", client.AdminAPIClient().BaseRootURL(), err)
	}

	// update the lastConfigSHA with the new updated checksum
	client.SetLastConfigSHA(newConfigSHA)

	return string(newConfigSHA), nil
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Private
// -----------------------------------------------------------------------------

type sendDiagnosticFn func(failed bool)

// prepareSendDiagnosticFn generates sendDiagnosticFn.
// Diagnostics are sent only when provided diagnostic config (--dump-config) is set.
func prepareSendDiagnosticFn(
	ctx context.Context,
	log logrus.FieldLogger,
	diagnosticConfig util.ConfigDumpDiagnostic,
	targetState *kongstate.KongState,
	targetContent *file.Content,
	pluginSchemaStore deckgen.PluginSchemaStore,
	filterTags []string,
	formatVersion string,
) sendDiagnosticFn {
	if diagnosticConfig == (util.ConfigDumpDiagnostic{}) {
		// noop, diagnostics won't be sent
		return func(bool) {}
	}

	var config *file.Content
	if diagnosticConfig.DumpsIncludeSensitive {
		redactedConfig := deckgen.ToDeckContent(ctx,
			log,
			targetState.SanitizedCopy(),
			pluginSchemaStore,
			filterTags,
			formatVersion,
		)
		config = redactedConfig
	} else {
		config = targetContent
	}

	return func(failed bool) {
		// Given that we can send multiple configs to this channel and
		// the fact that the API that exposes that can only expose 1 config
		// at a time it means that users utilizing the diagnostics API
		// might not see exactly what they intend to see i.e. come failures
		// or successfully send configs might be covered by those send
		// later on but we're OK with this limitation of said API.
		select {
		case diagnosticConfig.Configs <- util.ConfigDump{Failed: failed, Config: *config}:
			log.Debug("shipping config to diagnostic server")
		default:
			log.Error("config diagnostic buffer full, dropping diagnostic config")
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
	for _, obj := range uniqueObjects(reportedObjects, translationFailures) {
		c.kubernetesObjectStatusQueue.Publish(obj)
	}
}

func uniqueObjects(reportedObjects []client.Object, resourceFailures []failures.ResourceFailure) []client.Object {
	allCausingObjects := lo.FlatMap(resourceFailures, func(f failures.ResourceFailure, _ int) []client.Object {
		return f.CausingObjects()
	})
	allObjects := append(reportedObjects, allCausingObjects...)
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
			c.eventRecorder.Event(obj, corev1.EventTypeWarning, reason, failure.Message())
		}
	}
}
