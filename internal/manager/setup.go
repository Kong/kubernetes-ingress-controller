package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/avast/retry-go/v4"
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

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
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
		LeaderElection:         leaderElectionEnabled(logger, c, dbmode),
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
		if c.PublishService.String() != "" {
			watchNamespaces = append(c.WatchNamespaces, c.PublishService.Namespace)
		}
		controllerOpts.NewCache = cache.MultiNamespacedCacheBuilder(watchNamespaces)
	}

	if len(c.LeaderElectionNamespace) > 0 {
		controllerOpts.LeaderElectionNamespace = c.LeaderElectionNamespace
	}

	return controllerOpts, nil
}

func leaderElectionEnabled(logger logr.Logger, c *Config, dbmode string) bool {
	if c.Konnect.ConfigSynchronizationEnabled {
		logger.Info("Konnect config synchronisation enabled, enabling leader election")
		return true
	}

	if dbmode == "off" {
		if c.KongAdminSvc.Name != "" {
			logger.Info("DB-less mode detected with service detection, enabling leader election")
			return true
		}
		logger.Info("DB-less mode detected, disabling leader election")
		return false
	}

	logger.Info("Database mode detected, enabling leader election")
	return true
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

	kongclients, err := managerConfig.getKongClients(ctx)
	if err != nil {
		return err
	}
	// For now using first client is kind of OK. Using Consumer and Plugin
	// services from first kong client should theoretically return the same
	// results as for all other clients. There might be instances where
	// configurations in different Kong Gateways are ever so slightly
	// different but that shouldn't cause a fatal failure.
	//
	// TODO: We should take a look at this sooner rather than later.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3363
	designatedKongClient := kongclients[0].AdminAPIClient()
	srv, err := admission.MakeTLSServer(ctx, &managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			designatedKongClient.Consumers,
			designatedKongClient.Plugins,
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
func setupDataplaneAddressFinder(mgrc client.Client, c *Config, log logr.Logger) (*dataplane.AddressFinder, *dataplane.AddressFinder, error) {
	if !c.UpdateStatus {
		return nil, nil, nil
	}

	defaultAddressFinder, err := buildDataplaneAddressFinder(mgrc, c.PublishStatusAddress, c.PublishService)
	if err != nil {
		return nil, nil, fmt.Errorf("status updates enabled but no method to determine data-plane addresses: %w", err)
	}
	udpAddressFinder, err := buildDataplaneAddressFinder(mgrc, c.PublishStatusAddressUDP, c.PublishServiceUDP)
	if err != nil {
		log.Info("falling back to a default address finder for UDP", "reason", err.Error())
		udpAddressFinder = defaultAddressFinder
	}

	return defaultAddressFinder, udpAddressFinder, nil
}

func buildDataplaneAddressFinder(mgrc client.Client, publishStatusAddress []string, publishServiceNn types.NamespacedName) (*dataplane.AddressFinder, error) {
	addressFinder := dataplane.NewAddressFinder()

	if len(publishStatusAddress) > 0 {
		addressFinder.SetOverrides(publishStatusAddress)
		return addressFinder, nil
	}
	if publishServiceNn.String() != "" {
		addressFinder.SetGetter(generateAddressFinderGetter(mgrc, publishServiceNn))
		return addressFinder, nil
	}

	return nil, errors.New("no publish status address or publish service were provided")
}

func generateAddressFinderGetter(mgrc client.Client, publishServiceNn types.NamespacedName) func(context.Context) ([]string, error) {
	return func(ctx context.Context) ([]string, error) {
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
	}
}

// getKongClients returns the kong clients given the config.
// When a list of URLs is provided via --kong-admin-url then those are used
// to create the list of clients.
// When a headless service name is provided via --kong-admin-svc then that is used
// to obtain a list of endpoints via EndpointSlice lookup in kubernetes API.
func (c *Config) getKongClients(ctx context.Context) ([]*adminapi.Client, error) {
	httpclient, err := adminapi.MakeHTTPClient(&c.KongAdminAPIConfig, c.KongAdminToken)
	if err != nil {
		return nil, err
	}

	var addresses []string

	// If kong-admin-svc flag has been specified then use it to get the list
	// of Kong Admin API endpoints.
	if c.KongAdminSvc.Name != "" {
		kubeClient, err := c.GetKubeClient()
		if err != nil {
			return nil, err
		}

		// Retry this as we may encounter an error of getting 0 addresses,
		// which can mean that Kong instances meant to be configured by this controller
		// are not yet ready.
		// If we end up in a situation where none of them are ready then bail
		// because we have more code that relies on the configuration of Kong
		// instance and without an address and there's no way to initialize the
		// configuration validation and sending code.
		err = retry.Do(func() error {
			s, err := adminapi.GetURLsForService(ctx, kubeClient, c.KongAdminSvc)
			if err != nil {
				return err
			}
			if s.Len() == 0 {
				return fmt.Errorf("no endpoints for kong admin service: %q", c.KongAdminSvc)
			}
			addresses = s.UnsortedList()
			return nil
		},
			retry.Attempts(60),
			retry.DelayType(retry.FixedDelay),
			retry.Delay(time.Second),
			retry.OnRetry(func(_ uint, err error) {
				logrus.New().WithError(err).Error("failed to create kong client(s)")
			}),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Otherwise fallback to the list of kong admin URLs.
		addresses = c.KongAdminURLs
	}

	clients := make([]*adminapi.Client, 0, len(addresses))
	for _, address := range addresses {
		client, err := adminapi.NewKongClientForWorkspace(ctx, address, c.KongWorkspace, httpclient)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}
