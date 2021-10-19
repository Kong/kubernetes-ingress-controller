package manager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup Utility Functions
// -----------------------------------------------------------------------------

func setupLoggers(c *Config) (logrus.FieldLogger, logr.Logger, error) {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make logger: %w", err)
	}

	if c.LogReduceRedundancy {
		deprecatedLogger.Info("WARNING: log stifling has been enabled (experimental)")
		deprecatedLogger = util.MakeDebugLoggerWithReducedRedudancy(os.Stdout, &logrus.TextFormatter{}, 3, time.Second*30)
	}

	logger := logrusr.NewLogger(deprecatedLogger)
	ctrl.SetLogger(logger)

	return deprecatedLogger, logger, nil
}

func setupControllerOptions(logger logr.Logger, c *Config, scheme *runtime.Scheme) ctrl.Options {
	controllerOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         c.EnableLeaderElection,
		LeaderElectionID:       c.LeaderElectionID,
		SyncPeriod:             &c.SyncPeriod,
	}
	// determine how to configure namespace watchers
	switch len(c.WatchNamespaces) {
	case 0:
		// watch all namespaces
		controllerOpts.Namespace = corev1.NamespaceAll
	case 1:
		// watch one namespace
		controllerOpts.Namespace = c.WatchNamespaces[0]
	default:
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("manager set up with multiple namespaces", "namespaces", c.WatchNamespaces)
		controllerOpts.NewCache = cache.MultiNamespacedCacheBuilder(c.WatchNamespaces)
	}

	return controllerOpts
}

func setupKongConfig(ctx context.Context, logger logr.Logger, c *Config) (sendconfig.Kong, error) {
	kongClient, err := c.GetKongClient(ctx)
	if err != nil {
		return sendconfig.Kong{}, fmt.Errorf("unable to build kong api client: %w", err)
	}

	var filterTags []string
	if ok, err := kongClient.Tags.Exists(ctx); err != nil {
		logger.Error(err, "tag filtering disabled because Kong Admin API does not support tags")
	} else if ok {
		logger.Info("tag filtering enabled", "tags", c.FilterTags)
		filterTags = c.FilterTags
	}

	cfg := sendconfig.Kong{
		URL:               c.KongAdminURL,
		FilterTags:        filterTags,
		Concurrency:       c.Concurrency,
		Client:            kongClient,
		PluginSchemaStore: util.NewPluginSchemaStore(kongClient),
		ConfigDone:        make(chan file.Content),
	}

	return cfg, nil
}

func setupProxyServer(ctx context.Context,
	logger logr.Logger, fieldLogger logrus.FieldLogger,
	mgr manager.Manager, kongConfig sendconfig.Kong,
	diagnostic util.ConfigDumpDiagnostic, c *Config,
) (proxy.Proxy, error) {
	if c.ProxySyncSeconds < proxy.DefaultSyncSeconds {
		logger.Info(fmt.Sprintf("WARNING: --proxy-sync-seconds is configured for %fs, in DBLESS mode this may result in"+
			" problems of inconsistency in the proxy state. For DBLESS mode %fs+ is recommended (3s is the default).",
			c.ProxySyncSeconds, proxy.DefaultSyncSeconds,
		))
	}

	syncTickDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxySyncSeconds))
	if err != nil {
		logger.Error(err, "%s is not a valid number of seconds to stagger the proxy server synchronization")
		return nil, err
	}

	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxyTimeoutSeconds))
	if err != nil {
		logger.Error(err, "%s is not a valid number of seconds to the timeout config for the kong client")
		return nil, err
	}

	return proxy.NewCacheBasedProxyWithStagger(ctx,
		fieldLogger.WithField("subsystem", "proxy-cache-resolver"),
		mgr.GetClient(),
		kongConfig,
		c.IngressClassName,
		c.EnableReverseSync,
		syncTickDuration,
		timeoutDuration,
		diagnostic,
		sendconfig.UpdateKongAdminSimple)
}
