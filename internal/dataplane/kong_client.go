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
}

// NewKongClient provides a new KongClient object after connecting to the
// data-plane API and verifying integrity.
func NewKongClient(
	logger logrus.FieldLogger,
	timeout time.Duration,
	ingressClass string,
	enableReverseSync bool,
	diagnostic util.ConfigDumpDiagnostic,
	kongConfig sendconfig.Kong,
) (*KongClient, error) {
	// build the client object
	cache := store.NewCacheStores()
	c := &KongClient{
		logger:            logger,
		ingressClass:      ingressClass,
		enableReverseSync: enableReverseSync,
		requestTimeout:    timeout,
		diagnostic:        diagnostic,
		prometheusMetrics: metrics.NewCtrlFuncMetrics(),
		cache:             &cache,
		kongConfig:        kongConfig,
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
	return c.cache.Add(obj)
}

// DeleteObject accepts a Kubernetes controller-runtime client.Object and removes it from the configuration cache.
// The delete action will asynchronously be converted to Kong DSL and applied to the Kong Admin API.
// A status will later be added to the object whether the configuration update succeeds or fails.
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

	// initialize a parser and convert the Kubernetes objects to Kong objects
	p := parser.NewParser(c.logger, storer)
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

	// generate the deck configuration to be applied to the admin API
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
	timedCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	newConfigSHA, err := sendconfig.PerformUpdate(timedCtx,
		c.logger,
		&c.kongConfig,
		c.kongConfig.InMemory,
		c.enableReverseSync,
		targetConfig,
		c.kongConfig.FilterTags,
		nil,
		c.lastConfigSHA,
		false,
		c.prometheusMetrics,
	)
	if err != nil {
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

	// update the lastConfigSHA with the new updated checksum
	c.lastConfigSHA = newConfigSHA
	return nil
}
