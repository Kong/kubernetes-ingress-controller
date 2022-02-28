package dataplane

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blang/semver/v4"
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
	// Logger is the log writer that will be used to log information about
	// data-plane configuration runtime events.
	Logger logrus.FieldLogger

	// IngressClass indicates which Kubernetes IngressClass is supported for
	// objects which define IngressClass to indicate inclusion.
	IngressClass string

	// Cache is the Kubernetes object cache which is used to list Kubernetes
	// objects for parsing into Kong objects.
	Cache *store.CacheStores

	// KongConfig is the client configuration for the Kong Admin API
	KongConfig sendconfig.Kong

	// EnableReverseSync indicates that reverse sync should be enabled for
	// updates to the data-plane.
	EnableReverseSync bool

	// RequestTimeout is the maximum amount of time that should be waited for
	// requests to the data-plane to receive a response.
	RequestTimeout time.Duration

	// Diagnostic is the client and configuration for reporting diagnostic
	// information during data-plane update runtime.
	Diagnostic util.ConfigDumpDiagnostic

	// PrometheusMetrics is the client for shipping metrics information
	// updates to the prometheus exporter.
	PrometheusMetrics *metrics.CtrlFuncMetrics

	version       semver.Version
	dbmode        string
	lastConfigSHA []byte
	lock          sync.RWMutex
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Public Methods
// -----------------------------------------------------------------------------

// Initialize connects to the Kong Admin API to determine metadata and set
// configuration options based on the backend gateway's configuration.
func (c *KongClient) Initialize() error {
	// download the kong root configuration (and validate connectivity to the proxy API)
	root, err := c.RootWithTimeout()
	if err != nil {
		return err
	}

	// pull the proxy configuration out of the root config and validate it
	proxyConfig, ok := root["configuration"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid root configuration, expected a map[string]interface{} got %T", proxyConfig["configuration"])
	}

	// validate the database configuration for the proxy and check for supported database configurations
	dbmode, ok := proxyConfig["database"].(string)
	if !ok {
		return fmt.Errorf("invalid database configuration, expected a string got %t", proxyConfig["database"])
	}
	switch dbmode {
	case "off", "":
		c.KongConfig.InMemory = true
	case "postgres":
		c.KongConfig.InMemory = false
	case "cassandra":
		return fmt.Errorf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported; you must migrate to a Postgres-backed or DB-less deployment")
	default:
		return fmt.Errorf("%s is not a supported database backend", dbmode)
	}

	// validate the proxy version
	proxySemver, err := kong.ParseSemanticVersion(kong.VersionFromInfo(root))
	if err != nil {
		return err
	}

	// store the gathered configuration options
	c.KongConfig.Version = proxySemver
	c.dbmode = dbmode

	return nil
}

// UpdateObject accepts a Kubernetes controller-runtime client.Object and adds/updates that to the configuration cache.
// It will be asynchronously converted into the upstream Kong DSL and applied to the Kong Admin API.
// A status will later be added to the object whether the configuration update succeeds or fails.
func (c *KongClient) UpdateObject(obj client.Object) error {
	return c.Cache.Add(obj)
}

// DeleteObject accepts a Kubernetes controller-runtime client.Object and removes it from the configuration cache.
// The delete action will asynchronously be converted to Kong DSL and applied to the Kong Admin API.
// A status will later be added to the object whether the configuration update succeeds or fails.
func (c *KongClient) DeleteObject(obj client.Object) error {
	return c.Cache.Delete(obj)
}

// ObjectExists indicates whether or not any version of the provided object is already present in the proxy.
func (c *KongClient) ObjectExists(obj client.Object) (bool, error) {
	_, exists, err := c.Cache.Get(obj)
	return exists, err
}

// Listeners retrieves the currently configured listeners from the
// underlying proxy so that callers can gather this metadata to
// know which ports and protocols are in use by the proxy.
func (c *KongClient) Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error) {
	return c.KongConfig.Client.Listeners(ctx)
}

// RootWithTimeout provides the root configuration from Kong, but uses a configurable timeout to avoid long waits if the Admin API
// is not yet ready to respond. If a timeout error occurs, the caller is responsible for providing a retry mechanism.
func (c *KongClient) RootWithTimeout() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.RequestTimeout)
	defer cancel()
	return c.KongConfig.Client.Root(ctx)
}

// DBMode indicates which database the Kong Gateway is using
func (c *KongClient) DBMode() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.dbmode
}

// -----------------------------------------------------------------------------
// Dataplane Client - Kong - Interface Implementation
// -----------------------------------------------------------------------------

// Update parses the Cache present in the client and converts current
// Kubernetes state into Kong objects and state, and then ships the
// resulting configuration to the data-plane (Kong Admin API).
func (c *KongClient) Update(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*c.Cache, c.IngressClass, false, false, false, c.Logger)
	kongstate, err := parser.Build(c.Logger, storer)
	if err != nil {
		c.PrometheusMetrics.TranslationCount.With(prometheus.Labels{
			metrics.SuccessKey: metrics.SuccessFalse,
		}).Inc()
		return err
	}
	c.PrometheusMetrics.TranslationCount.With(prometheus.Labels{
		metrics.SuccessKey: metrics.SuccessTrue,
	}).Inc()

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(ctx,
		c.Logger, kongstate,
		c.KongConfig.PluginSchemaStore,
		c.KongConfig.FilterTags,
	)

	// generate diagnostic configuration if enabled
	// "diagnostic" will be empty if --dump-config is not set
	var diagnosticConfig *file.Content
	if c.Diagnostic != (util.ConfigDumpDiagnostic{}) {
		if !c.Diagnostic.DumpsIncludeSensitive {
			redactedConfig := deckgen.ToDeckContent(ctx,
				c.Logger,
				kongstate.SanitizedCopy(),
				c.KongConfig.PluginSchemaStore,
				c.KongConfig.FilterTags,
			)
			diagnosticConfig = redactedConfig
		} else {
			diagnosticConfig = targetConfig
		}
	}

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, c.RequestTimeout)
	defer cancel()
	newConfigSHA, err := sendconfig.PerformUpdate(timedCtx,
		c.Logger,
		&c.KongConfig,
		c.KongConfig.InMemory,
		c.EnableReverseSync,
		targetConfig,
		c.KongConfig.FilterTags,
		nil,
		c.lastConfigSHA,
		false,
		c.PrometheusMetrics,
	)
	if err != nil {
		// ship diagnostics if enabled
		if c.Diagnostic != (util.ConfigDumpDiagnostic{}) {
			select {
			case c.Diagnostic.Configs <- util.ConfigDump{Failed: true, Config: *diagnosticConfig}:
				c.Logger.Debug("shipping config to diagnostic server")
			default:
				c.Logger.Error("config diagnostic buffer full, dropping diagnostic config")
			}
		}
		return err
	}

	// ship diagnostics if enabled
	if c.Diagnostic != (util.ConfigDumpDiagnostic{}) {
		select {
		case c.Diagnostic.Configs <- util.ConfigDump{Failed: false, Config: *diagnosticConfig}:
			c.Logger.Debug("shipping config to diagnostic server")
		default:
			c.Logger.Error("config diagnostic buffer full, dropping diagnostic config")
		}
	}

	// update the lastConfigSHA with the new updated checksum
	c.lastConfigSHA = newConfigSHA
	return nil
}
