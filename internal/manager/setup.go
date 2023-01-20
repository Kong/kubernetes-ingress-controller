package manager

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/go-logr/logr"
	"github.com/kong/deck/cprint"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup Utility Functions
// -----------------------------------------------------------------------------

// SetupLoggers sets up the loggers for the controller manager.
func SetupLoggers(c *Config, output io.Writer) (logrus.FieldLogger, logr.Logger, error) {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat, output)
	if err != nil {
		return nil, logr.Logger{}, fmt.Errorf("failed to make logger: %w", err)
	}

	if c.LogReduceRedundancy {
		deprecatedLogger.Info("WARNING: log stifling has been enabled (experimental)")
		deprecatedLogger = util.MakeDebugLoggerWithReducedRedudancy(output, &logrus.TextFormatter{}, 3, time.Second*30)
	}

	logger := logrusr.New(deprecatedLogger)
	ctrl.SetLogger(logger)

	if c.LogLevel != "trace" && c.LogLevel != "debug" {
		// disable deck's per-change diff output
		cprint.DisableOutput = true
	}

	return deprecatedLogger, logger, nil
}

func setupControllerOptions(logger logr.Logger, c *Config, dbmode string) (ctrl.Options, error) {
	var leaderElection bool
	if dbmode == "off" {
		logger.Info("DB-less mode detected, disabling leader election")
		leaderElection = false
	} else {
		logger.Info("Database mode detected, enabling leader election")
		leaderElection = true
	}

	logger.Info("building the manager runtime scheme and loading apis into the scheme")
	scheme, err := getScheme()
	if err != nil {
		return ctrl.Options{}, err
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
		watchNamespaces := c.WatchNamespaces

		// in all other cases we are a multi-namespace setup and must watch all the
		// c.WatchNamespaces.
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("manager set up with multiple namespaces", "namespaces", watchNamespaces)

		// if publish service has been provided the namespace for it should be
		// watched so that controllers can see updates to the service.
		if c.PublishService.NN.String() != "" {
			watchNamespaces = append(c.WatchNamespaces, c.PublishService.NN.Namespace)
		}
		controllerOpts.NewCache = cache.MultiNamespacedCacheBuilder(watchNamespaces)
	}

	if len(c.LeaderElectionNamespace) > 0 {
		controllerOpts.LeaderElectionNamespace = c.LeaderElectionNamespace
	}

	return controllerOpts, nil
}

func setupDataplaneSynchronizer(
	logger logr.Logger,
	fieldLogger logrus.FieldLogger,
	mgr manager.Manager,
	dataplaneClient dataplane.Client,
	proxySyncSeconds float32,
) (*dataplane.Synchronizer, error) {
	if proxySyncSeconds < dataplane.DefaultSyncSeconds {
		logger.Info(fmt.Sprintf(
			"WARNING: --proxy-sync-seconds is configured for %fs, in DBLESS mode this may result in"+
				" problems of inconsistency in the proxy state. For DBLESS mode %fs+ is recommended (3s is the default).",
			proxySyncSeconds, dataplane.DefaultSyncSeconds,
		))
	}

	dataplaneSynchronizer, err := dataplane.NewSynchronizer(
		fieldLogger.WithField("subsystem", "dataplane-synchronizer"),
		dataplaneClient,
		dataplane.WithStagger(time.Duration(proxySyncSeconds*float32(time.Second))),
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

func setupAdmissionServer(
	ctx context.Context,
	managerConfig *Config,
	managerClient client.Client,
	deprecatedLogger logrus.FieldLogger,
) error {
	logger := deprecatedLogger.WithField("component", "admission-server")

	if managerConfig.AdmissionServer.ListenAddr == "off" {
		logger.Info("admission webhook server disabled")
		return nil
	}

	kongclients, err := getKongClients(ctx,
		managerConfig.KongAdminURL,
		managerConfig.KongWorkspace,
		managerConfig.KongAdminAPIConfig,
	)
	if err != nil {
		return err
	}
	srv, err := admission.MakeTLSServer(ctx, &managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			// For now using first client is kind of OK. Using Consumer and Plugin
			// services from first kong client should theoretically return the same
			// results as for all other clients. There might be instances where
			// configurations in different Kong Gateways are ever so slightly
			// different but that shouldn't cause a fatal failure.
			//
			// TODO: We should take a look at this sooner rather than later.
			// https://github.com/Kong/kubernetes-ingress-controller/issues/3363
			kongclients[0].Consumers,
			kongclients[0].Plugins,
			logger,
			managerClient,
			managerConfig.IngressClassName,
		),
		Logger: logger,
	}, logger)
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServeTLS("", "")
		logger.WithError(err).Error("admission webhook server stopped")
	}()
	return nil
}

// setupDataplaneAddressFinder returns a default and UDP address finder. These finders return the override addresses if
// set or the publish service addresses if no overrides are set. If no UDP overrides or UDP publish service are set,
// the UDP finder will also return the default addresses. If no override or publish service is set, this function
// returns nil finders and an error.
func setupDataplaneAddressFinder(
	mgrc client.Client,
	c *Config,
) (*dataplane.AddressFinder, *dataplane.AddressFinder, error) {
	dataplaneAddressFinder := dataplane.NewAddressFinder()
	udpDataplaneAddressFinder := dataplane.NewAddressFinder()
	var getter func(ctx context.Context) ([]string, error)
	if c.UpdateStatus {
		// Default
		if overrideAddrs := c.PublishStatusAddress; len(overrideAddrs) > 0 {
			dataplaneAddressFinder.SetOverrides(overrideAddrs)
		} else if c.PublishService.String() != "" {
			publishServiceNn := c.PublishService.NN
			dataplaneAddressFinder.SetGetter(func(ctx context.Context) ([]string, error) {
				svc := new(corev1.Service)
				if err := mgrc.Get(ctx, publishServiceNn, svc); err != nil {
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
					return nil, fmt.Errorf("waiting for addresses to be provisioned for publish service %s", publishServiceNn)
				}

				return addrs, nil
			})
		} else {
			return nil, nil, fmt.Errorf("status updates enabled but no method to determine data-plane addresses, need either --publish-service or --publish-status-address")
		}

		// UDP. falls back to default if not configured
		if udpOverrideAddrs := c.PublishStatusAddressUDP; len(udpOverrideAddrs) > 0 {
			udpDataplaneAddressFinder.SetUDPOverrides(udpOverrideAddrs)
		} else if c.PublishServiceUDP.String() != "" {
			publishServiceNn := c.PublishServiceUDP.NN
			udpDataplaneAddressFinder.SetGetter(func(ctx context.Context) ([]string, error) {
				svc := new(corev1.Service)
				if err := mgrc.Get(ctx, publishServiceNn, svc); err != nil {
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
					return nil, fmt.Errorf("waiting for addresses to be provisioned for publish service %s", publishServiceNn)
				}

				return addrs, nil
			})
		} else {
			udpDataplaneAddressFinder.SetGetter(getter)
		}
	}

	return dataplaneAddressFinder, udpDataplaneAddressFinder, nil
}

func generateAddressFinderGetter(
	mgrc client.Client,
	nsn types.NamespacedName,
) func(ctx context.Context) ([]string, error) {
	return func(ctx context.Context) ([]string, error) {
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
	}
}
