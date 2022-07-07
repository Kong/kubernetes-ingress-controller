package manager

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/go-logr/logr"
	"github.com/kong/deck/cprint"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup Utility Functions
// -----------------------------------------------------------------------------

func setupLoggers(c *Config) (logrus.FieldLogger, logr.Logger, error) {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat)
	if err != nil {
		return nil, logr.Logger{}, fmt.Errorf("failed to make logger: %w", err)
	}

	if c.LogReduceRedundancy {
		deprecatedLogger.Info("WARNING: log stifling has been enabled (experimental)")
		deprecatedLogger = util.MakeDebugLoggerWithReducedRedudancy(os.Stdout, &logrus.TextFormatter{}, 3, time.Second*30)
	}

	logger := logrusr.New(deprecatedLogger)
	ctrl.SetLogger(logger)

	if c.LogLevel != "trace" && c.LogLevel != "debug" {
		// disable deck's per-change diff output
		cprint.DisableOutput = true
	}

	return deprecatedLogger, logger, nil
}

func setupControllerOptions(logger logr.Logger, c *Config, scheme *runtime.Scheme,
	dbmode string,
) (ctrl.Options, error) {
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

	var leaderElection bool
	if dbmode == "off" {
		logger.Info("DB-less mode detected, disabling leader election")
		leaderElection = false
	} else {
		logger.Info("Database mode detected, enabling leader election")
		leaderElection = true
	}

	// configure the general controller options
	controllerOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         leaderElection,
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

	if len(c.LeaderElectionNamespace) > 0 {
		controllerOpts.LeaderElectionNamespace = c.LeaderElectionNamespace
	}

	return controllerOpts, nil
}

func setupKongConfig(ctx context.Context, kongClient *kong.Client, logger logr.Logger, c *Config) sendconfig.Kong {
	var filterTags []string
	if ok, err := kongClient.Tags.Exists(ctx); err != nil {
		logger.Error(err, "tag filtering disabled because Kong Admin API does not support tags")
	} else if ok {
		logger.Info("tag filtering enabled", "tags", c.FilterTags)
		filterTags = c.FilterTags
	}

	return sendconfig.Kong{
		URL:               c.KongAdminURL,
		FilterTags:        filterTags,
		Concurrency:       c.Concurrency,
		Client:            kongClient,
		PluginSchemaStore: util.NewPluginSchemaStore(kongClient),
	}
}

func setupDataplaneSynchronizer(
	logger logr.Logger,
	fieldLogger logrus.FieldLogger,
	mgr manager.Manager,
	dataplaneClient dataplane.Client,
	c *Config,
) (*dataplane.Synchronizer, error) {
	if c.ProxySyncSeconds < dataplane.DefaultSyncSeconds {
		logger.Info(fmt.Sprintf("WARNING: --proxy-sync-seconds is configured for %fs, in DBLESS mode this may result in"+
			" problems of inconsistency in the proxy state. For DBLESS mode %fs+ is recommended (3s is the default).",
			c.ProxySyncSeconds, dataplane.DefaultSyncSeconds,
		))
	}

	syncTickDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxySyncSeconds))
	if err != nil {
		logger.Error(err, "%s is not a valid number of seconds to stagger the proxy server synchronization")
		return nil, err
	}

	dataplaneSynchronizer, err := dataplane.NewSynchronizerWithStagger(
		fieldLogger.WithField("subsystem", "dataplane-synchronizer"),
		dataplaneClient,
		syncTickDuration,
	)
	if err != nil {
		return nil, err
	}

	err = mgr.Add(dataplaneSynchronizer)
	if err != nil {
		return nil, err
	}

	return dataplaneSynchronizer, nil
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
	srv, err := admission.MakeTLSServer(ctx, &managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			kongclient.Consumers,
			kongclient.Plugins,
			log,
			managerClient,
			managerConfig.IngressClassName,
		),
		Logger: logger,
	}, log)
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServeTLS("", "")
		log.WithError(err).Error("admission webhook server stopped")
	}()
	return nil
}

func setupDataplaneAddressFinder(ctx context.Context, mgrc client.Client, c *Config) (*dataplane.AddressFinder, error) {
	dataplaneAddressFinder := dataplane.NewAddressFinder()
	if c.UpdateStatus {
		if overrideAddrs := c.PublishStatusAddress; len(overrideAddrs) > 0 {
			dataplaneAddressFinder.SetOverrides(overrideAddrs)
		} else if c.PublishService != "" {
			parts := strings.Split(c.PublishService, "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("publish service %s is invalid, expecting <namespace>/<name>", c.PublishService)
			}
			nsn := types.NamespacedName{
				Namespace: parts[0],
				Name:      parts[1],
			}
			dataplaneAddressFinder.SetGetter(func() ([]string, error) {
				svc := new(corev1.Service)
				if err := mgrc.Get(ctx, nsn, svc); err != nil {
					return nil, err
				}

				var addrs []string
				switch svc.Spec.Type { //nolint:exhaustive
				case corev1.ServiceTypeLoadBalancer:
					for _, lbaddr := range svc.Status.LoadBalancer.Ingress {
						if lbaddr.IP != "" {
							addrs = append(addrs, lbaddr.IP)
						}
						if lbaddr.Hostname != "" {
							addrs = append(addrs, lbaddr.Hostname)
						}
					}
				default:
					addrs = append(addrs, svc.Spec.ClusterIPs...)
				}

				if len(addrs) == 0 {
					return nil, fmt.Errorf("waiting for addresses to be provisioned for publish service %s/%s", nsn.Namespace, nsn.Name)
				}

				return addrs, nil
			})
		} else {
			return nil, fmt.Errorf("status updates enabled but no method to determine data-plane addresses, need either --publish-service or --publish-status-address")
		}
	}

	return dataplaneAddressFinder, nil
}
