package dataplane

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Public Types
// -----------------------------------------------------------------------------

// KongClient is a threadsafe high level API client for the Kong data-plane
// which parses Kubernetes object caches into Kong Admin configurations and
// sends them as updates to the data-plane (Kong Admin API).
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
	kongConfig sendconfig.Kong

	// dbmode indicates the current database mode of the backend Kong Admin API
	dbmode string

	// lastConfigSHA is a checksum of the last successful update to the data-plane
	lastConfigSHA []byte

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
	kubernetesObjectReportsFilter k8sobj.Set
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
	kongConfig sendconfig.Kong,
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
	}

	// download the kong root configuration (and validate connectivity to the proxy API)
	root, err := c.RootWithTimeout()
	if err != nil {
		return nil, err
	}

	// pull the proxy configuration out of the root config and validate it
	proxyConfig, ok := root["configuration"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid root configuration, expected a map[string]interface{} got %T", proxyConfig["configuration"])
	}

	// validate the database configuration for the proxy and check for supported database configurations
	dbmode, ok := proxyConfig["database"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid database configuration, expected a string got %t", proxyConfig["database"])
	}
	switch dbmode {
	case "off", "":
		c.kongConfig.InMemory = true
	case "postgres":
		c.kongConfig.InMemory = false
	case "cassandra":
		return nil, fmt.Errorf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported; you must migrate to a Postgres-backed or DB-less deployment")
	default:
		return nil, fmt.Errorf("%s is not a supported database backend", dbmode)
	}

	// validate the proxy version
	proxySemver, err := kong.ParseSemanticVersion(kong.VersionFromInfo(root))
	if err != nil {
		return nil, err
	}

	// store the gathered configuration options
	c.kongConfig.Version = proxySemver
	c.dbmode = dbmode

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

// Listeners retrieves the currently configured listeners from the
// underlying proxy so that callers can gather this metadata to
// know which ports and protocols are in use by the proxy.
func (c *KongClient) Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error) {
	return c.kongConfig.Client.Listeners(ctx)
}

// RootWithTimeout provides the root configuration from Kong, but uses a configurable timeout to avoid long waits if the Admin API
// is not yet ready to respond. If a timeout error occurs, the caller is responsible for providing a retry mechanism.
func (c *KongClient) RootWithTimeout() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	return c.kongConfig.Client.Root(ctx)
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
	return c.kubernetesObjectReportsFilter.Has(obj)
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

// DBMode indicates which database the Kong Gateway is using
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
	storer := store.New(*c.cache, c.ingressClass, false, false, false, c.logger)

	// initialize a parser
	c.logger.Debug("parsing kubernetes objects into data-plane configuration")
	p := parser.NewParser(c.logger, storer)
	if c.AreKubernetesObjectReportsEnabled() {
		p.EnableKubernetesObjectReports()
	}
	if c.AreCombinedServiceRoutesEnabled() {
		p.EnableCombinedServiceRoutes()
	}

	// parse the Kubernetes objects from the storer into Kong configuration
	kongstate, err := p.Build()
	if err != nil {
		c.prometheusMetrics.TranslationCount.With(prometheus.Labels{
			metrics.SuccessKey: metrics.SuccessFalse,
		}).Inc()
		return err
	}
	c.prometheusMetrics.TranslationCount.With(prometheus.Labels{
		metrics.SuccessKey: metrics.SuccessTrue,
	}).Inc()
	c.logger.Debug("successfully built data-plane configuration")

	// generate the deck configuration to be applied to the admin API
	c.logger.Debug("converting configuration to deck config")
	targetConfig := deckgen.ToDeckContent(ctx,
		c.logger, kongstate,
		c.kongConfig.PluginSchemaStore,
		c.kongConfig.FilterTags,
	)

	// generate diagnostic configuration if enabled
	// "diagnostic" will be empty if --dump-config is not set
	var diagnosticConfig *file.Content
	if c.diagnostic != (util.ConfigDumpDiagnostic{}) {
		if !c.diagnostic.DumpsIncludeSensitive {
			redactedConfig := deckgen.ToDeckContent(ctx,
				c.logger,
				kongstate.SanitizedCopy(),
				c.kongConfig.PluginSchemaStore,
				c.kongConfig.FilterTags,
			)
			diagnosticConfig = redactedConfig
		} else {
			diagnosticConfig = targetConfig
		}
	}

	// apply the configuration update in Kong
	c.logger.Debug("sending configuration to Kong Admin API")
	timedCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	newConfigSHA, err := sendconfig.PerformUpdate(timedCtx,
		c.logger,
		&c.kongConfig,
		c.kongConfig.InMemory,
		c.enableReverseSync,
		c.skipCACertificates,
		targetConfig,
		c.kongConfig.FilterTags,
		nil,
		c.lastConfigSHA,
		c.prometheusMetrics,
	)
	if err != nil {
		if expired, ok := timedCtx.Deadline(); ok && time.Now().After(expired) {
			c.logger.Warn("exceeded Kong API timeout, consider increasing --proxy-timeout-seconds")
		}
		// ship diagnostics if enabled
		if c.diagnostic != (util.ConfigDumpDiagnostic{}) {
			select {
			case c.diagnostic.Configs <- util.ConfigDump{Failed: true, Config: *diagnosticConfig}:
				c.logger.Debug("shipping config to diagnostic server")
			default:
				c.logger.Error("config diagnostic buffer full, dropping diagnostic config")
			}
		}
		return err
	}

	// ship diagnostics if enabled
	if c.diagnostic != (util.ConfigDumpDiagnostic{}) {
		select {
		case c.diagnostic.Configs <- util.ConfigDump{Failed: false, Config: *diagnosticConfig}:
			c.logger.Debug("shipping config to diagnostic server")
		default:
			c.logger.Error("config diagnostic buffer full, dropping diagnostic config")
		}
	}

	// report on configured Kubernetes objects if enabled
	if c.AreKubernetesObjectReportsEnabled() {
		if string(c.lastConfigSHA) != string(newConfigSHA) {
			report := p.GenerateKubernetesObjectReport()
			c.logger.Debugf("triggering report for %d configured Kubernetes objects", len(report))
			c.triggerKubernetesObjectReport(report...)
		} else {
			c.logger.Debug("no configuration change, skipping kubernetes object report")
		}
	}

	// update the lastConfigSHA with the new updated checksum
	c.lastConfigSHA = newConfigSHA
	return nil
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Private
// -----------------------------------------------------------------------------

// triggerKubernetesObjectReport will update the KongClient with a set which
// enables filtering for which objects are currently applied to the data-plane,
// as well as updating the c.kubernetesObjectStatusQueue to queue those objects
// for reconciliation so their statuses can be properly updated.
func (c *KongClient) triggerKubernetesObjectReport(objs ...client.Object) {
	// first a new set of the included objects for the most recent configuration
	// needs to be generated.
	set := k8sobj.Set{}
	for _, obj := range objs {
		set.Insert(obj)
	}

	c.updateKubernetesObjectReportFilter(set)

	// after the filter has been updated we signal the status queue so that the
	// control-plane can update the Kubernetes object statuses for affected objs.
	// this has to be done in a separate loop so that the filter is in place
	// before the objects are enqueued, as the filter is used by the control-plane
	for _, obj := range objs {
		c.kubernetesObjectStatusQueue.Publish(obj)
	}
}

// updateKubernetesObjectReportFilter overrides the internal object set with
// a new provided set.
func (c *KongClient) updateKubernetesObjectReportFilter(set k8sobj.Set) {
	c.kubernetesObjectReportLock.Lock()
	defer c.kubernetesObjectReportLock.Unlock()
	c.kubernetesObjectReportsFilter = set
}
