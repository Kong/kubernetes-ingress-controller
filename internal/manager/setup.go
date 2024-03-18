package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/cprint"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	ctrllicense "github.com/kong/kubernetes-ingress-controller/v3/controllers/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	konnectLicense "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup Utility Functions
// -----------------------------------------------------------------------------

// SetupLoggers sets up the loggers for the controller manager.
func SetupLoggers(c *Config, output io.Writer) (logr.Logger, error) {
	zapBase, err := util.MakeLogger(c.LogLevel, c.LogFormat, output)
	if err != nil {
		return logr.Logger{}, fmt.Errorf("failed to make logger: %w", err)
	}
	logger := zapr.NewLoggerWithOptions(zapBase, zapr.LogInfoLevel("v"))

	if c.LogLevel != "trace" && c.LogLevel != "debug" {
		// disable deck's per-change diff output
		cprint.DisableOutput = true
	}

	// Prevents controller-runtime from logging
	// [controller-runtime] log.SetLogger(...) was never called; logs will not be displayed.
	ctrllog.SetLogger(logger)

	return logger, nil
}

func setupManagerOptions(ctx context.Context, logger logr.Logger, c *Config, dbmode dpconf.DBMode) (ctrl.Options, error) {
	logger.Info("Building the manager runtime scheme and loading apis into the scheme")
	scheme, err := scheme.Get()
	if err != nil {
		return ctrl.Options{}, err
	}

	// configure the general manager options
	managerOpts := ctrl.Options{
		GracefulShutdownTimeout: c.GracefulShutdownTimeout,
		Scheme:                  scheme,
		Metrics: metricsserver.Options{
			BindAddress: c.MetricsAddr,
		},
		WebhookServer:    webhook.NewServer(webhook.Options{Port: 9443}),
		LeaderElection:   leaderElectionEnabled(logger, c, dbmode),
		LeaderElectionID: c.LeaderElectionID,
		Cache: cache.Options{
			SyncPeriod: &c.SyncPeriod,
		},
		Logger:    ctrl.LoggerFrom(ctx),
		NewClient: newManagerClient,
	}

	// If there are no configured watch namespaces, then we're watching ALL namespaces,
	// and we don't have to bother individually caching any particular namespaces.
	// This is the default behavior of the controller-runtime manager.
	// If there are configured watch namespaces, then we're watching only those namespaces.
	if len(c.WatchNamespaces) > 0 {
		watchNamespaces := c.WatchNamespaces

		// In all other cases we are a multi-namespace setup and must watch all the
		// c.WatchNamespaces.
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("Manager set up with multiple namespaces", "namespaces", watchNamespaces)

		// If ingress service has been provided the namespace for it should be
		// watched so that controllers can see updates to the service.
		if s, ok := c.PublishService.Get(); ok {
			watchNamespaces = append(c.WatchNamespaces, s.Namespace)
		}
		watched := make(map[string]cache.Config)
		for _, n := range sets.NewString(watchNamespaces...).List() {
			watched[n] = cache.Config{}
		}
		managerOpts.Cache.DefaultNamespaces = watched
	}

	if len(c.LeaderElectionNamespace) > 0 {
		managerOpts.LeaderElectionNamespace = c.LeaderElectionNamespace
	}

	return managerOpts, nil
}

func leaderElectionEnabled(logger logr.Logger, c *Config, dbmode dpconf.DBMode) bool {
	if c.Konnect.ConfigSynchronizationEnabled {
		logger.Info("Konnect config synchronisation enabled, enabling leader election")
		return true
	}

	if dbmode.IsDBLessMode() {
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
		logger.WithName("dataplane-synchronizer"),
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
	referenceIndexers ctrlref.CacheIndexers,
	managerClient client.Client,
	logger logr.Logger,
	translatorFeatures translator.FeatureFlags,
	storer store.Storer,
) error {
	admissionLogger := logger.WithName("admission-server")

	if managerConfig.AdmissionServer.ListenAddr == "off" {
		logger.Info("Admission webhook server disabled")
		return nil
	}

	adminAPIServicesProvider := admission.NewDefaultAdminAPIServicesProvider(clientsManager)
	srv, err := admission.MakeTLSServer(ctx, &managerConfig.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			admissionLogger,
			managerClient,
			managerConfig.IngressClassName,
			adminAPIServicesProvider,
			translatorFeatures,
			storer,
		),
		ReferenceIndexers: referenceIndexers,
		Logger:            admissionLogger,
	}, admissionLogger)
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServeTLS("", "")
		logger.Error(err, "Admission webhook server stopped")
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
		log.Info("Falling back to a default address finder for UDP", "reason", err.Error())
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
		switch svc.Spec.Type {
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
	const (
		delay = time.Second
	)

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
		retry.Delay(delay),
		retry.OnRetry(func(_ uint, err error) {
			// log the error if the error is NOT caused by 0 available gateway endpoints.
			if !errors.As(err, &NoAvailableEndpointsError{}) {
				logger.Error(err, "Failed to create kong client(s)")
			}
			logger.Error(err, "Failed to create kong client(s), retrying...", "delay", delay)
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

// setupLicenseGetter sets up a license getter to get Kong license from Konnect or `KongLicense` CRD.
// If synchoroniztion license from Konnect is enabled, it sets up and returns a Konnect license agent.
// If controller of `KongLicense` CRD is enabled and sync license with Konnect is disabled,
// it starts and returns a KongLicense controller.
func setupLicenseGetter(
	ctx context.Context,
	c *Config,
	setupLog logr.Logger,
	mgr manager.Manager,
	statusQueue *status.Queue,
) (translator.LicenseGetter, error) {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/3922
	// This requires the Konnect client, which currently requires c.Konnect.ConfigSynchronizationEnabled also.
	// We need to figure out exactly how that config surface works. Initial direction says add a separate toggle, but
	// we probably want to avoid that long term. If we do have separate toggles, we need an AND condition that sets up
	// the client and makes it available to all Konnect-related subsystems.
	if c.Konnect.LicenseSynchronizationEnabled {
		konnectLicenseAPIClient, err := konnectLicense.NewClient(c.Konnect)
		if err != nil {
			return nil, fmt.Errorf("failed creating konnect client: %w", err)
		}
		setupLog.Info("Starting license agent")
		agent := license.NewAgent(
			konnectLicenseAPIClient,
			ctrl.LoggerFrom(ctx).WithName("license-agent"),
			license.WithInitialPollingPeriod(c.Konnect.InitialLicensePollingPeriod),
			license.WithPollingPeriod(c.Konnect.LicensePollingPeriod),
		)
		err = mgr.Add(agent)
		if err != nil {
			return nil, fmt.Errorf("could not add license agent to manager: %w", err)
		}
		return agent, nil
	}
	// Enable KongLicense controller if license synchornizition from Konnect is disabled.
	if c.KongLicenseEnabled && !c.Konnect.LicenseSynchronizationEnabled {
		setupLog.Info("Starting KongLicense controller")
		licenseController := ctrllicense.NewKongV1Alpha1KongLicenseReconciler(
			mgr.GetClient(),
			ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongLicense"),
			mgr.GetScheme(),
			ctrllicense.NewLicenseCache(),
			c.CacheSyncTimeout,
			statusQueue,
			ctrllicense.LicenseControllerTypeKIC,
			mo.Some(c.LeaderElectionID),
			mo.None[ctrllicense.ValidatorFunc](),
		)
		dynamicLicenseController := ctrllicense.WrapKongLicenseReconcilerToDynamicCRDController(
			ctx, mgr, licenseController,
		)
		err := dynamicLicenseController.SetupWithManager(mgr)
		if err != nil {
			return nil, fmt.Errorf("failed to start KongLicense controller: %w", err)
		}
		return licenseController, nil
	}

	return nil, nil
}
