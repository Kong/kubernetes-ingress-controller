package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	ctrllicense "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/license"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/k8s"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	konnectLicense "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/scheme"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup Utility Functions
// -----------------------------------------------------------------------------

func setupManagerOptions(ctx context.Context, logger logr.Logger, c *managercfg.Config, dbmode dpconf.DBMode) ctrl.Options {
	logger.Info("Building the manager runtime scheme and loading apis into the scheme")

	// configure the general manager options
	managerOpts := ctrl.Options{
		Controller: config.Controller{
			// This is needed because controller-runtime keeps a global list of controller
			// names and panics if there are duplicates.
			// This is a workaround for that in tests.
			// Ref: https://github.com/kubernetes-sigs/controller-runtime/pull/2902#issuecomment-2284194683
			SkipNameValidation: lo.ToPtr(true),
		},
		GracefulShutdownTimeout: c.GracefulShutdownTimeout,
		Scheme:                  scheme.Get(),
		Metrics: metricsserver.Options{
			BindAddress: c.MetricsAddr,
			FilterProvider: func() func(c *rest.Config, httpClient *http.Client) (metricsserver.Filter, error) {
				switch c.MetricsAccessFilter {
				case managercfg.MetricsAccessFilterOff:
					return nil
				case managercfg.MetricsAccessFilterRBAC:
					return filters.WithAuthenticationAndAuthorization
				default:
					// This is checked in flags validation so this should never happen.
					panic("unsupported metrics filter")
				}
			}(),
		},
		WebhookServer:    webhook.NewServer(webhook.Options{Port: 9443}),
		LeaderElection:   leaderElectionEnabled(logger, *c, dbmode),
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
		watchNamespaces := sets.NewString(c.WatchNamespaces...)

		// In all other cases we are a multi-namespace setup and must watch all the
		// c.WatchNamespaces.
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		logger.Info("Manager set up with multiple namespaces", "namespaces", watchNamespaces)

		// If ingress service has been provided the namespace for it should be
		// watched so that controllers can see updates to the service.
		if s, ok := c.PublishService.Get(); ok {
			watchNamespaces.Insert(s.Namespace)
		}

		watched := make(map[string]cache.Config)
		for _, n := range watchNamespaces.List() {
			watched[n] = cache.Config{}
		}
		managerOpts.Cache.DefaultNamespaces = watched
	}

	if len(c.LeaderElectionNamespace) > 0 {
		managerOpts.LeaderElectionNamespace = c.LeaderElectionNamespace
	}

	return managerOpts
}

func leaderElectionEnabled(logger logr.Logger, c managercfg.Config, dbmode dpconf.DBMode) bool {
	if c.LeaderElectionForce == managercfg.LeaderElectionEnabled {
		logger.Info("leader election forcibly enabled")
		return true
	}
	if c.LeaderElectionForce == managercfg.LeaderElectionDisabled {
		logger.Info("leader election forcibly disabled")
		return false
	}
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

func (m *Manager) setupAdmissionServer(
	ctx context.Context,
	referenceIndexers ctrlref.CacheIndexers,
	translatorFeatures translator.FeatureFlags,
	storer store.Storer,
) error {
	admissionLogger := ctrl.LoggerFrom(ctx).WithName("admission-server")

	if m.cfg.AdmissionServer.ListenAddr == "off" {
		admissionLogger.Info("Admission webhook server disabled")
		return nil
	}

	adminAPIServicesProvider := admission.NewDefaultAdminAPIServicesProvider(m.clientsManager)
	srv, err := admission.MakeTLSServer(m.cfg.AdmissionServer, &admission.RequestHandler{
		Validator: admission.NewKongHTTPValidator(
			admissionLogger,
			m.m.GetClient(),
			m.cfg.IngressClassName,
			adminAPIServicesProvider,
			translatorFeatures,
			storer,
		),
		ReferenceIndexers: referenceIndexers,
		Logger:            admissionLogger,
	})
	if err != nil {
		return err
	}

	m.admissionServer = mo.Some(srv)
	return nil
}

// setupDataplaneAddressFinder returns a default and UDP address finder. These finders return the override addresses if
// set or the publish service addresses if no overrides are set. If no UDP overrides or UDP publish service are set,
// the UDP finder will also return the default addresses. If no override or publish service is set, this function
// returns nil finders and an error.
func setupDataplaneAddressFinder(mgrc client.Client, c managercfg.Config, log logr.Logger) (*dataplane.AddressFinder, *dataplane.AddressFinder, error) {
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

func buildDataplaneAddressFinder(mgrc client.Client, publishStatusAddress []string, publishServiceNN mo.Option[k8stypes.NamespacedName]) (*dataplane.AddressFinder, error) {
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
func adminAPIClients(
	ctx context.Context,
	c managercfg.Config,
	logger logr.Logger,
	discoverer *adminapi.Discoverer,
	factory adminapi.ClientFactory,
) ([]*adminapi.Client, error) {
	// If kong-admin-svc flag has been specified then use it to get the list
	// of Kong Admin API endpoints.
	if kongAdminSvc, ok := c.KongAdminSvc.Get(); ok {
		kubeClient, err := k8s.GetKubeClient(c)
		if err != nil {
			return nil, fmt.Errorf("failed to get kubernetes client: %w", err)
		}
		return AdminAPIClientFromServiceDiscovery(ctx, logger, kongAdminSvc, kubeClient, discoverer, factory,
			retry.Attempts(c.KongAdminInitializationRetries), retry.Delay(c.KongAdminInitializationRetryDelay))
	}

	// Otherwise fallback to the list of kong admin URLs.
	addresses := c.KongAdminURLs
	clients := make([]*adminapi.Client, 0, len(addresses))
	for _, address := range addresses {
		err := retry.Do(
			func() error {
				cl, err := adminapi.NewKongClientForWorkspace(ctx, address, c.KongWorkspace, c.KongAdminAPIConfig, c.KongAdminToken)
				if err != nil {
					return err
				}
				clients = append(clients, cl)
				return nil
			},
			retry.Attempts(c.KongAdminInitializationRetries),
			retry.Delay(c.KongAdminInitializationRetryDelay),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create admin API client to %s: %w", address, err)
		}
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
		delay                        = time.Second
		createAdminAPIClientAttempts = 60
	)

	// Retry this as we may encounter an error of getting 0 addresses,
	// which can mean that Kong instances meant to be configured by this controller
	// are not yet ready.
	// If we end up in a situation where none of them are ready then bail
	// because we have more code that relies on the configuration of Kong
	// instance and without an address and there's no way to initialize the
	// configuration validation and sending code.
	fetchEndpointsRetryOptions := append([]retry.Option{
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
		fetchEndpointsRetryOptions...,
	)
	if err != nil {
		return nil, err
	}

	createAdminAPIClientRetryOptions := append([]retry.Option{
		retry.Context(ctx),
		retry.Attempts(createAdminAPIClientAttempts),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(delay),
	}, retryOpts...)
	clients := make([]*adminapi.Client, 0, len(adminAPIs))
	for _, adminAPI := range adminAPIs {
		var client *adminapi.Client
		err := retry.Do(
			func() error {
				cl, err := factory.CreateAdminAPIClient(ctx, adminAPI)
				if err != nil {
					return err
				}
				client = cl
				return nil
			},
			createAdminAPIClientRetryOptions...,
		)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// setupLicenseGetter sets up a license getter to get Kong license from Konnect or `KongLicense` CRD.
// If synchoroniztion license from Konnect is enabled, it sets up and returns a Konnect license agent.
// If controller of `KongLicense` CRD is enabled and sync license with Konnect is disabled,
// it starts and returns a KongLicense controller.
func setupLicenseGetter(
	ctx context.Context,
	c managercfg.Config,
	setupLog logr.Logger,
	mgr manager.Manager,
	statusQueue *status.Queue,
) (license.Getter, error) {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/3922
	// This requires the Konnect client, which currently requires c.Konnect.ConfigSynchronizationEnabled also.
	// We need to figure out exactly how that config surface works. Initial direction says add a separate toggle, but
	// we probably want to avoid that long term. If we do have separate toggles, we need an AND condition that sets up
	// the client and makes it available to all Konnect-related subsystems.
	if c.Konnect.LicenseSynchronizationEnabled {
		setupLog.Info("Creating konnect license client")
		konnectLicenseAPIClient, err := konnectLicense.NewClient(
			c.Konnect,
			ctrl.LoggerFrom(ctx).WithName("konnect-license-client"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed creating konnect client: %w", err)
		}

		if c.Konnect.LicenseStorageEnabled {
			setupLog.Info("Creating a storage to store fetched Konnect license")
			nn, err := util.GetPodNN()
			if err != nil {
				return nil, err
			}
			licenseStore := konnectLicense.NewSecretLicenseStore(
				mgr.GetClient(), nn.Namespace, c.Konnect.ControlPlaneID,
			)
			konnectLicenseAPIClient.WithLicenseStore(licenseStore)
			konnectLicenseAPIClient = konnectLicenseAPIClient.WithLicenseStore(licenseStore)
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

// setupKonnectConfigSynchronizerWithMgr sets up Konnect config synchronizer and adds it to the manager runnables.
func setupKonnectConfigSynchronizerWithMgr(
	ctx context.Context,
	mgr manager.Manager,
	cfg managercfg.Config,
	kongConfig sendconfig.Config,
	updateStrategyResolver sendconfig.UpdateStrategyResolver,
	configStatusNotifier clients.ConfigStatusNotifier,
	metricsRecorder metrics.Recorder,
) (*konnect.ConfigSynchronizer, error) {
	s := konnect.NewConfigSynchronizer(
		konnect.ConfigSynchronizerParams{
			Logger:                 ctrl.LoggerFrom(ctx).WithName("konnect-config-synchronizer"),
			KongConfig:             kongConfig,
			ConfigUploadTicker:     clock.NewTickerWithDuration(cfg.Konnect.UploadConfigPeriod),
			KonnectClientFactory:   adminapi.NewKonnectClientFactory(cfg.Konnect, ctrl.LoggerFrom(ctx).WithName("konnect-client-factory")),
			UpdateStrategyResolver: updateStrategyResolver,
			ConfigChangeDetector:   sendconfig.NewKonnectConfigurationChangeDetector(),
			ConfigStatusNotifier:   configStatusNotifier,
			MetricsRecorder:        metricsRecorder,
		},
	)
	err := mgr.Add(s)
	if err != nil {
		return nil, fmt.Errorf("could not add Konnect config synchronizer to manager: %w", err)
	}
	return s, nil
}
