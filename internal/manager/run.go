// Package manager implements the controller manager for all controllers
package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/telemetry"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/utils/kongconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(
	ctx context.Context,
	c *Config,
	diagnostic diagnostics.ClientDiagnostic,
	logger logr.Logger,
) error {
	setupLog := ctrl.LoggerFrom(ctx).WithName("setup")
	setupLog.Info("Starting controller manager", "release", metadata.Release, "repo", metadata.Repo, "commit", metadata.Commit)
	setupLog.Info("The ingress class name has been set", "value", c.IngressClassName)

	gateway.SetControllerName(gatewayapi.GatewayController(c.GatewayAPIControllerName))

	setupLog.Info("Getting enabled options and features")
	featureGates, err := featuregates.New(setupLog, c.FeatureGates)
	if err != nil {
		return fmt.Errorf("failed to configure feature gates: %w", err)
	}
	setupLog.Info("Getting the kubernetes client configuration")
	kubeconfig, err := c.GetKubeconfig()
	if err != nil {
		return fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
	}

	setupLog.Info("Starting standalone health check server")
	healthServer := &healthCheckServer{}
	healthServer.setHealthzCheck(healthz.Ping)
	healthServer.Start(ctx, c.ProbeAddr, setupLog.WithName("health-check"))

	adminAPIsDiscoverer, err := adminapi.NewDiscoverer(sets.New(c.KongAdminSvcPortNames...), c.GatewayDiscoveryDNSStrategy)
	if err != nil {
		return fmt.Errorf("failed to create admin apis discoverer: %w", err)
	}

	err = c.Resolve()
	if err != nil {
		return fmt.Errorf("failed to resolve configuration: %w", err)
	}

	adminAPIClientsFactory := adminapi.NewClientFactoryForWorkspace(c.KongWorkspace, c.KongAdminAPIConfig, c.KongAdminToken)

	setupLog.Info("Getting the kong admin api client configuration")
	initialKongClients, err := c.adminAPIClients(
		ctx,
		setupLog.WithName("initialize-kong-clients"),
		adminAPIsDiscoverer,
		adminAPIClientsFactory,
	)
	if err != nil {
		return fmt.Errorf("unable to build kong api client(s): %w", err)
	}

	// Get Kong configuration root(s) to validate them and extract Kong's version.
	kongRoots, err := kongconfig.GetRoots(ctx, setupLog, c.KongAdminInitializationRetries, c.KongAdminInitializationRetryDelay, initialKongClients)
	if err != nil {
		return fmt.Errorf("could not retrieve Kong admin root(s): %w", err)
	}

	kongStartUpConfig, err := kongconfig.ValidateRoots(kongRoots, c.SkipCACertificates)
	if err != nil {
		return fmt.Errorf("could not validate Kong admin root(s) configuration: %w", err)
	}
	dbMode := kongStartUpConfig.DBMode
	routerFlavor := kongStartUpConfig.RouterFlavor
	v := kongStartUpConfig.Version

	kongSemVersion := semver.Version{Major: v.Major(), Minor: v.Minor(), Patch: v.Patch()}

	kongConfig := sendconfig.Config{
		Version:                       kongSemVersion,
		InMemory:                      dbMode.IsDBLessMode(),
		Concurrency:                   c.Concurrency,
		FilterTags:                    c.FilterTags,
		SkipCACertificates:            c.SkipCACertificates,
		EnableReverseSync:             c.EnableReverseSync,
		ExpressionRoutes:              dpconf.ShouldEnableExpressionRoutes(routerFlavor),
		SanitizeKonnectConfigDumps:    featureGates.Enabled(featuregates.SanitizeKonnectConfigDumps),
		FallbackConfiguration:         featureGates.Enabled(featuregates.FallbackConfiguration),
		UseLastValidConfigForFallback: c.UseLastValidConfigForFallback,
	}

	setupLog.Info("Configuring and building the controller manager")
	managerOpts, err := setupManagerOptions(ctx, setupLog, c, dbMode)
	if err != nil {
		return fmt.Errorf("unable to setup manager options: %w", err)
	}

	mgr, err := ctrl.NewManager(kubeconfig, managerOpts)
	if err != nil {
		return fmt.Errorf("unable to create controller manager: %w", err)
	}

	if err := waitForKubernetesAPIReadiness(ctx, setupLog, mgr); err != nil {
		return fmt.Errorf("unable to connect to Kubernetes API: %w", err)
	}

	setupLog.Info("Initializing Dataplane Client")
	var eventRecorder record.EventRecorder
	if c.EmitKubernetesEvents {
		setupLog.Info("Emitting Kubernetes events enabled, creating an event recorder for " + KongClientEventRecorderComponentName)
		eventRecorder = mgr.GetEventRecorderFor(KongClientEventRecorderComponentName)
	} else {
		setupLog.Info("Emitting Kubernetes events disabled, discarding all events")
		// Create an empty record.FakeRecorder with no Events channel to discard all events.
		eventRecorder = &record.FakeRecorder{}
	}

	readinessChecker := clients.NewDefaultReadinessChecker(adminAPIClientsFactory, setupLog.WithName("readiness-checker"))
	clientsManager, err := clients.NewAdminAPIClientsManager(
		ctx,
		logger,
		initialKongClients,
		readinessChecker,
	)
	if err != nil {
		return fmt.Errorf("failed to create AdminAPIClientsManager: %w", err)
	}
	clientsManager = clientsManager.WithDBMode(dbMode)

	if c.KongAdminSvc.IsPresent() {
		setupLog.Info("Running AdminAPIClientsManager loop")
		clientsManager.Run()
	}

	translatorFeatureFlags := translator.NewFeatureFlags(
		featureGates,
		routerFlavor,
		c.UpdateStatus,
		kongStartUpConfig.Version.IsKongGatewayEnterprise(),
	)

	referenceIndexers := ctrlref.NewCacheIndexers(setupLog.WithName("reference-indexers"))
	cache := store.NewCacheStores()
	storer := store.New(cache, c.IngressClassName, logger)
	configTranslator, err := translator.NewTranslator(logger, storer, c.KongWorkspace, translatorFeatureFlags, NewSchemaServiceGetter(clientsManager))
	if err != nil {
		return fmt.Errorf("failed to create translator: %w", err)
	}

	setupLog.Info("Starting Admission Server")
	if err := setupAdmissionServer(ctx, c, clientsManager, referenceIndexers, mgr.GetClient(), logger, translatorFeatureFlags, storer); err != nil {
		return err
	}

	updateStrategyResolver := sendconfig.NewDefaultUpdateStrategyResolver(kongConfig, logger)
	configurationChangeDetector := sendconfig.NewDefaultConfigurationChangeDetector(logger)
	kongConfigFetcher := configfetcher.NewDefaultKongLastGoodConfigFetcher(translatorFeatureFlags.FillIDs, c.KongWorkspace)
	fallbackConfigGenerator := fallback.NewGenerator(fallback.NewDefaultCacheGraphProvider(), logger)
	dataplaneClient, err := dataplane.NewKongClient(
		logger,
		time.Duration(c.ProxyTimeoutSeconds*float32(time.Second)),
		diagnostic,
		kongConfig,
		eventRecorder,
		dbMode,
		clientsManager,
		updateStrategyResolver,
		configurationChangeDetector,
		kongConfigFetcher,
		configTranslator,
		&cache,
		fallbackConfigGenerator,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize kong data-plane client: %w", err)
	}

	setupLog.Info("Initializing Dataplane Synchronizer")
	synchronizer, err := setupDataplaneSynchronizer(logger, mgr, dataplaneClient, c.ProxySyncSeconds, c.InitCacheSyncDuration)
	if err != nil {
		return fmt.Errorf("unable to initialize dataplane synchronizer: %w", err)
	}

	var kubernetesStatusQueue *status.Queue
	if c.UpdateStatus {
		setupLog.Info("Starting Status Updater")
		kubernetesStatusQueue = status.NewQueue(status.WithBufferSize(c.UpdateStatusQueueBufferSize))
		dataplaneClient.EnableKubernetesObjectReports(kubernetesStatusQueue)
	} else {
		setupLog.Info("Status updates disabled, skipping status updater")
	}

	setupLog.Info("Initializing Dataplane address Discovery")
	dataplaneAddressFinder, udpDataplaneAddressFinder, err := setupDataplaneAddressFinder(mgr.GetClient(), c, setupLog)
	if err != nil {
		return err
	}

	setupLog.Info("Starting Enabled Controllers")
	controllers := setupControllers(
		ctx,
		mgr,
		dataplaneClient,
		referenceIndexers,
		dataplaneAddressFinder,
		udpDataplaneAddressFinder,
		kubernetesStatusQueue,
		c,
		featureGates,
		clientsManager,
		adminAPIsDiscoverer,
	)
	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}

	// BUG: kubebuilder (at the time of writing - 3.0.0-rc.1) does not allow this tag anywhere else than main.go
	// See https://github.com/kubernetes-sigs/kubebuilder/issues/932
	// +kubebuilder:scaffold:builder

	setupLog.Info("Add readiness probe to health server")
	healthServer.setReadyzCheck(readyzHandler(mgr, synchronizer))
	instanceIDProvider := NewInstanceIDProvider()

	if c.Konnect.ConfigSynchronizationEnabled {
		konnectNodesAPIClient, err := nodes.NewClient(c.Konnect)
		if err != nil {
			return fmt.Errorf("failed creating konnect client: %w", err)
		}
		// In case of failures when building Konnect related objects, we're not returning errors as Konnect is not
		// considered critical feature, and it should not break the basic functionality of the controller.

		// Run the Konnect Admin API client initialization in a separate goroutine to not block while ensuring
		// connection.
		go setupKonnectAdminAPIClientWithClientsMgr(ctx, c.Konnect, clientsManager, setupLog)

		// Setup Konnect NodeAgent with manager.
		if err := setupKonnectNodeAgentWithMgr(
			c,
			mgr,
			konnectNodesAPIClient,
			dataplaneClient,
			clientsManager,
			setupLog,
			instanceIDProvider,
		); err != nil {
			setupLog.Error(err, "Failed to setup Konnect NodeAgent with manager, skipping")
		}
	}

	// Setup and inject license getter.
	licenseGetter, err := setupLicenseGetter(
		ctx,
		c,
		setupLog,
		mgr,
		kubernetesStatusQueue,
	)
	if err != nil {
		setupLog.Error(err, "Failed to create a license getter from configuration")
		return err
	}
	if licenseGetter != nil {
		setupLog.Info("Inject license getter to config translator",
			"license_getter_type", fmt.Sprintf("%T", licenseGetter))
		configTranslator.InjectLicenseGetter(licenseGetter)
		kongConfigFetcher.InjectLicenseGetter(licenseGetter)
	}

	if c.AnonymousReports {
		stopAnonymousReports, err := telemetry.SetupAnonymousReports(
			ctx,
			logger.WithName("telemetry"),
			kubeconfig,
			clientsManager,
			telemetry.ReportConfig{
				SplunkEndpoint:                   c.SplunkEndpoint,
				SplunkEndpointInsecureSkipVerify: c.SplunkEndpointInsecureSkipVerify,
				TelemetryPeriod:                  c.TelemetryPeriod,
				ReportValues: telemetry.ReportValues{
					PublishServiceNN:               c.PublishService.OrEmpty(),
					FeatureGates:                   featureGates,
					MeshDetection:                  len(c.WatchNamespaces) == 0,
					KonnectSyncEnabled:             c.Konnect.ConfigSynchronizationEnabled,
					GatewayServiceDiscoveryEnabled: c.KongAdminSvc.IsPresent(),
				},
			},
			instanceIDProvider,
		)
		if err != nil {
			setupLog.Error(err, "Failed setting up anonymous reports")
		} else {
			defer stopAnonymousReports()
		}
		setupLog.Info("Anonymous reports enabled")
	} else {
		setupLog.Info("Anonymous reports disabled, skipping")
	}

	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
}

// waitForKubernetesAPIReadiness waits for the Kubernetes API to be ready. It's used as a prerequisite to run any
// controller components (i.e. Manager along with its Runnables).
// It retries with a timeout of 1m and a fixed delay of 1s.
func waitForKubernetesAPIReadiness(ctx context.Context, logger logr.Logger, mgr manager.Manager) error {
	const (
		timeout = time.Minute
		delay   = time.Second
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	readinessEndpointURL, err := url.JoinPath(mgr.GetConfig().Host, "readyz")
	if err != nil {
		return fmt.Errorf("failed to build readiness check URL: %w", err)
	}

	return retry.Do(func() error {
		// Call the readiness check of the Kubernetes API server: https://kubernetes.io/docs/reference/using-api/health-checks/.
		resp, err := mgr.GetHTTPClient().Get(readinessEndpointURL)
		if err != nil {
			return fmt.Errorf("failed to connect to %q: %w", readinessEndpointURL, err)
		}
		defer resp.Body.Close()
		// We're waiting for the readiness check to return status 200.
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("readiness check %q returned status %d", readinessEndpointURL, resp.StatusCode)
		}
		return nil
	},
		retry.Context(ctx),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(0), // We're using a context with timeout, so we don't need to limit the number of attempts.
		retry.LastErrorOnly(true),
		retry.OnRetry(func(_ uint, err error) {
			logger.Info("Retrying Kubernetes API readiness check after error", "error", err.Error())
		}),
	)
}

// setupKonnectNodeAgentWithMgr creates and adds Konnect NodeAgent as the manager's Runnable.
// Returns error if failed to create Konnect NodeAgent.
func setupKonnectNodeAgentWithMgr(
	c *Config,
	mgr manager.Manager,
	konnectNodeAPIClient *nodes.Client,
	dataplaneClient *dataplane.KongClient,
	clientsManager *clients.AdminAPIClientsManager,
	logger logr.Logger,
	instanceIDProvider *InstanceIDProvider,
) error {
	var hostname string
	nn, err := util.GetPodNN()
	if err != nil {
		logger.Error(err, "Failed getting pod name and/or namespace, fallback to use hostname as node name in Konnect")
		hostname, _ = os.Hostname()
	} else {
		hostname = nn.String()
		logger.Info(fmt.Sprintf("Using %s as controller's node name in Konnect", hostname))
	}
	version := metadata.Release

	// Set channel to send config status.
	configStatusNotifier := clients.NewChannelConfigNotifier(logger)
	dataplaneClient.SetConfigStatusNotifier(configStatusNotifier)

	agent := konnect.NewNodeAgent(
		hostname,
		version,
		c.Konnect.RefreshNodePeriod,
		logger,
		konnectNodeAPIClient,
		configStatusNotifier,
		konnect.NewGatewayClientGetter(logger, clientsManager),
		clientsManager,
		instanceIDProvider,
	)
	if err := mgr.Add(agent); err != nil {
		return fmt.Errorf("failed adding konnect.NodeAgent runnable to the manager: %w", err)
	}
	return nil
}

// setupKonnectAdminAPIClientWithClientsMgr initializes Konnect Admin API client and sets it to clientsManager.
// If it fails to initialize the client, it logs the error and returns.
func setupKonnectAdminAPIClientWithClientsMgr(
	ctx context.Context,
	config adminapi.KonnectConfig,
	clientsManager *clients.AdminAPIClientsManager,
	logger logr.Logger,
) {
	konnectAdminAPIClient, err := adminapi.NewKongClientForKonnectControlPlane(config)
	if err != nil {
		logger.Error(err, "Failed creating Konnect Control Plane Admin API client, skipping synchronisation")
		return
	}
	if err := adminapi.EnsureKonnectConnection(ctx, konnectAdminAPIClient.AdminAPIClient(), logger); err != nil {
		logger.Error(err, "Failed to ensure connection to Konnect Admin API, skipping synchronisation")
		return
	}

	clientsManager.SetKonnectClient(konnectAdminAPIClient)
	logger.Info("Initialized Konnect Admin API client")
}

type IsReady interface {
	IsReady() bool
}

func readyzHandler(mgr manager.Manager, dataplaneSynchronizer IsReady) func(*http.Request) error {
	return func(_ *http.Request) error {
		select {
		// If we're elected as leader then report readiness based on the readiness
		// of dataplane synchronizer.
		case <-mgr.Elected():
			if !dataplaneSynchronizer.IsReady() {
				return errors.New("synchronizer not yet configured")
			}
		// If we're not the leader then just report as ready.
		default:
		}
		return nil
	}
}
