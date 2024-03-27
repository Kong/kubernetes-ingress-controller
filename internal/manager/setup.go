package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/blang/semver/v4"
	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/cprint"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	dataplaneutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/dataplane"
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

	if c.LogLevel != "trace" && c.LogLevel != "debug" {
		// disable deck's per-change diff output
		cprint.DisableOutput = true
	}

	return deprecatedLogger, logger, nil
}

func setupControllerOptions(ctx context.Context, logger logr.Logger, c *Config, dbmode string, featureGates map[string]bool) (ctrl.Options, error) {
	logger.Info("building the manager runtime scheme and loading apis into the scheme")
	scheme, err := scheme.Get(featureGates)
	if err != nil {
		return ctrl.Options{}, err
	}

	// configure the general controller options
	controllerOpts := ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: c.MetricsAddr,
		},
		WebhookServer:    webhook.NewServer(webhook.Options{Port: 9443}),
		LeaderElection:   leaderElectionEnabled(logger, c, dbmode),
		LeaderElectionID: c.LeaderElectionID,
		Cache: cache.Options{
			SyncPeriod: &c.SyncPeriod,
		},
		Logger: ctrl.LoggerFrom(ctx),
	}

	// If there are no configured watch namespaces, then we're watching ALL namespaces,
	// and we don't have to bother individually caching any particular namespaces.
	// This is the default behavior of the controller-runtime manager.
	// If there are configured watch namespaces, then we're watching only those namespaces.
	if len(c.WatchNamespaces) > 0 {
		watchNamespaces := c.WatchNamespaces

		// in all other cases we are a multi-namespace setup and must watch all the
		// c.WatchNamespaces.
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("manager set up with multiple namespaces", "namespaces", watchNamespaces)

		// if publish service has been provided the namespace for it should be
		// watched so that controllers can see updates to the service.
		if s, ok := c.PublishService.Get(); ok {
			watchNamespaces = append(c.WatchNamespaces, s.Namespace)
		}
		watched := make(map[string]cache.Config)
		for _, n := range sets.NewString(watchNamespaces...).List() {
			watched[n] = cache.Config{}
		}
		controllerOpts.Cache.DefaultNamespaces = watched
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

	if dataplaneutil.IsDBLessMode(dbmode) {
		if c.KongAdminSvc.IsPresent() {
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
	initCacheSyncWait time.Duration,
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
		dataplane.WithInitCacheSyncDuration(initCacheSyncWait),
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
	clientsManager *clients.AdminAPIClientsManager,
	managerClient client.Client,
	deprecatedLogger logrus.FieldLogger,
	parserFeatures parser.FeatureFlags,
	kongVersion semver.Version,
) error {
	logger := deprecatedLogger.WithField("component", "admission-server")

	if managerConfig.AdmissionServer.ListenAddr == "off" {
		logger.Info("admission webhook server disabled")
		return nil
	}

	adminAPIServicesProvider := admission.NewDefaultAdminAPIServicesProvider(clientsManager)
	srv, err := admission.MakeTLSServer(ctx, &managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			logger,
			managerClient,
			managerConfig.IngressClassName,
			adminAPIServicesProvider,
			parserFeatures,
			kongVersion,
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

func buildDataplaneAddressFinder(mgrc client.Client, publishStatusAddress []string, publishServiceNN OptionalNamespacedName) (*dataplane.AddressFinder, error) {
	addressFinder := dataplane.NewAddressFinder()

	if len(publishStatusAddress) > 0 {
		addressFinder.SetOverrides(publishStatusAddress)
		return addressFinder, nil
	}
	if serviceNN, ok := publishServiceNN.Get(); ok {
		addressFinder.SetGetter(generateAddressFinderGetter(mgrc, serviceNN))
		return addressFinder, nil
	}

	return nil, errors.New("no publish status address or publish service were provided")
}

func generateAddressFinderGetter(mgrc client.Client, publishServiceNn k8stypes.NamespacedName) func(context.Context) ([]string, error) {
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

// adminAPIClients returns the kong clients given the config.
// When a list of URLs is provided via --kong-admin-url then those are used
// to create the list of clients.
// When a headless service name is provided via --kong-admin-svc then that is used
// to obtain a list of endpoints via EndpointSlice lookup in kubernetes API.
func (c *Config) adminAPIClients(
	ctx context.Context,
	logger logr.Logger,
	discoverer *adminapi.Discoverer,
	factory adminapi.ClientFactory,
) ([]*adminapi.Client, error) {
	httpclient, err := adminapi.MakeHTTPClient(&c.KongAdminAPIConfig, c.KongAdminToken)
	if err != nil {
		return nil, err
	}

	// If kong-admin-svc flag has been specified then use it to get the list
	// of Kong Admin API endpoints.
	if kongAdminSvc, ok := c.KongAdminSvc.Get(); ok {
		kubeClient, err := c.GetKubeClient()
		if err != nil {
			return nil, fmt.Errorf("failed to get kubernetes client: %w", err)
		}
		return AdminAPIClientFromServiceDiscovery(ctx, logger, kongAdminSvc, kubeClient, discoverer, factory)
	}

	// Otherwise fallback to the list of kong admin URLs.
	addresses := c.KongAdminURLs
	clients := make([]*adminapi.Client, 0, len(addresses))
	for _, address := range addresses {
		cl, err := adminapi.NewKongClientForWorkspace(ctx, address, c.KongWorkspace, httpclient)
		if err != nil {
			return nil, err
		}

		clients = append(clients, cl)
	}

	return clients, nil
}

type NoAvailableEndpointsError struct {
	serviceNN k8stypes.NamespacedName
}

func (e NoAvailableEndpointsError) Error() string {
	return fmt.Sprintf("no endpoints for service: %q", e.serviceNN)
}

type AdminAPIsDiscoverer interface {
	GetAdminAPIsForService(context.Context, client.Client, k8stypes.NamespacedName) (sets.Set[adminapi.DiscoveredAdminAPI], error)
}

type AdminAPIClientFactory interface {
	CreateAdminAPIClient(context.Context, adminapi.DiscoveredAdminAPI) (*adminapi.Client, error)
}

func AdminAPIClientFromServiceDiscovery(
	ctx context.Context,
	logger logr.Logger,
	kongAdminSvcNN k8stypes.NamespacedName,
	kubeClient client.Client,
	discoverer AdminAPIsDiscoverer,
	factory AdminAPIClientFactory,
	retryOpts ...retry.Option,
) ([]*adminapi.Client, error) {
	// Retry this as we may encounter an error of getting 0 addresses,
	// which can mean that Kong instances meant to be configured by this controller
	// are not yet ready.
	// If we end up in a situation where none of them are ready then bail
	// because we have more code that relies on the configuration of Kong
	// instance and without an address and there's no way to initialize the
	// configuration validation and sending code.
	retryOpts = append([]retry.Option{
		retry.Context(ctx),
		retry.Attempts(0),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Second),
		retry.OnRetry(func(_ uint, err error) {
			// log the error if the error is NOT caused by 0 available gateway endpoints.
			if !errors.As(err, &NoAvailableEndpointsError{}) {
				logger.Error(err, "failed to create kong client(s)")
			}
		}),
	}, retryOpts...)

	var adminAPIs []adminapi.DiscoveredAdminAPI
	err := retry.Do(func() error {
		s, err := discoverer.GetAdminAPIsForService(ctx, kubeClient, kongAdminSvcNN)
		if err != nil {
			return retry.Unrecoverable(err)
		}
		if s.Len() == 0 {
			return NoAvailableEndpointsError{serviceNN: kongAdminSvcNN}
		}
		adminAPIs = s.UnsortedList()
		return nil
	},
		retryOpts...,
	)
	if err != nil {
		return nil, err
	}

	clients := make([]*adminapi.Client, 0, len(adminAPIs))
	for _, adminAPI := range adminAPIs {
		cl, err := factory.CreateAdminAPIClient(ctx, adminAPI)
		if err != nil {
			return nil, err
		}
		clients = append(clients, cl)
	}

	return clients, nil
}
