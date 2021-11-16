package manager

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission"
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

func setupControllerOptions(logger logr.Logger, c *Config, scheme *runtime.Scheme) (ctrl.Options, error) {
	// some controllers may require additional namespaces to be cached and this
	// is currently done using the global manager client cache.
	//
	// See: https://github.com/Kong/kubernetes-ingress-controller/issues/2004
	requiredCacheNamespaces := make([]string, 0)

	// if publish service has been provided the namespace for it should be
	// watched so that controllers can see updates to the service.
	if c.PublishService != "" {
		publishServiceSplit := strings.SplitN(c.PublishService, "/", 3)
		if len(publishServiceSplit) != 2 {
			return ctrl.Options{}, fmt.Errorf("--publish-service was expected to be in format <namespace>/<name> but got %s", c.PublishService)
		}
		requiredCacheNamespaces = append(requiredCacheNamespaces, publishServiceSplit[0])
	}

	// configure the general controller options
	controllerOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         c.EnableLeaderElection,
		LeaderElectionID:       c.LeaderElectionID,
		SyncPeriod:             &c.SyncPeriod,
	}

	// configure the controller caching options
	if len(c.WatchNamespaces) == 0 {
		// if there are no configured watch namespaces, then we're watching ALL namespaces
		// and we don't have to bother individually caching any particular namespaces
		controllerOpts.Namespace = corev1.NamespaceAll
	} else {
		// in all other cases we are a multi-namespace setup and must watch all the
		// c.WatchNamespaces and additionalNamespacesToCache defined namespaces.
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("manager set up with multiple namespaces", "namespaces", c.WatchNamespaces)
		controllerOpts.NewCache = cache.MultiNamespacedCacheBuilder(append(c.WatchNamespaces, requiredCacheNamespaces...))
	}

	return controllerOpts, nil
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

func setupAdmissionServer(ctx context.Context, managerConfig *Config, managerClient client.Client) error {
	log, err := util.MakeLogger(managerConfig.LogLevel, managerConfig.LogFormat)
	if err != nil {
		return err
	}

	if managerConfig.AdmissionServer.ListenAddr == "off" {
		log.Info("admission webhook server disabled")
		return nil
	}

	logger := log.WithField("component", "admission-server")

	kongclient, err := managerConfig.GetKongClient(ctx)
	if err != nil {
		return err
	}
	srv, err := admission.MakeTLSServer(&managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			kongclient.Consumers,
			kongclient.Plugins,
			log,
			managerClient,
			managerConfig.IngressClassName,
		),
		Logger: logger,
	})
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServeTLS("", "")
		log.WithError(err).Error("admission webhook server stopped")
	}()
	return nil
}
