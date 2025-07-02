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
	"github.com/google/uuid"
	"github.com/samber/mo"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
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
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/utils/kongconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/metadata"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

type InstanceID interface {
	String() string
}

type Manager struct {
	cfg                  managercfg.Config
	m                    manager.Manager
	synchronizer         *dataplane.Synchronizer
	diagnosticsServer    mo.Option[diagnostics.Server]
	diagnosticsCollector mo.Option[*diagnostics.Collector]
	diagnosticsHandler   mo.Option[*diagnostics.HTTPHandler]
	admissionServer      mo.Option[*admission.Server]
	kubeconfig           *rest.Config
	clientsManager       *clients.AdminAPIClientsManager
	licenseGetter        mo.Option[license.Getter]
}

// New configures the controller manager call Start.
func New(
	ctx context.Context,
	instanceID InstanceID,
	c managercfg.Config,
	logger logr.Logger,
) (*Manager, error) {
	// Inject logger into the context so it can be used by the controllers and other components without
	// passing it explicitly when they accept a context.
	ctx = ctrl.LoggerInto(ctx, logger)

	// Ensure that instance of config is valid, otherwise don't start the manager.
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config invalid: %w", err)
	}
	existingFeatureGates := managercfg.GetFeatureGatesDefaults()
	for feature, enabled := range c.FeatureGates {
		logger.Info("Found configuration option for gated feature", "feature", feature, "enabled", enabled)
		if _, ok := existingFeatureGates[feature]; !ok {
			return nil, fmt.Errorf("%s is not a valid feature, please see the documentation: %s", feature, managercfg.DocsURL)
		}
	}

	m := &Manager{
		cfg: c,
	}
	diagnosticsClient := m.setupDiagnostics(ctx, c)
	setupLog := logger.WithName("setup")
	setupLog.Info("Starting controller manager", "release", metadata.Release, "repo", metadata.Repo, "commit", metadata.Commit)
	setupLog.Info("The ingress class name has been set", "value", c.IngressClassName)

	gateway.SetControllerName(gatewayapi.GatewayController(c.GatewayAPIControllerName))

	setupLog.Info("Getting the Kubernetes client configuration")
	if c.KubeRestConfig != nil {
		setupLog.Info("Using KubeRestConfig from configuration")
		m.kubeconfig = c.KubeRestConfig
	} else {
		setupLog.Info("Using Kubeconfig based on fields from configuration")
		kubeconfigConstructed, err := utils.GetKubeconfig(c)
		if err != nil {
			return nil, fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
		}
		m.kubeconfig = kubeconfigConstructed
	}

	adminAPIsDiscoverer, err := adminapi.NewDiscoverer(sets.New(c.KongAdminSvcPortNames...))
	if err != nil {
		return nil, fmt.Errorf("failed to create admin apis discoverer: %w", err)
	}

	if err = c.Resolve(); err != nil {
		return nil, fmt.Errorf("failed to resolve configuration: %w", err)
	}

	adminAPIClientsFactory := adminapi.NewClientFactoryForWorkspace(logger, c.KongWorkspace, c.KongAdminAPIConfig, c.KongAdminToken)

	setupLog.Info("Getting the kong admin api client configuration")
	initialKongClients, err := adminAPIClients(
		ctx,
		c,
		setupLog.WithName("initialize-kong-clients"),
		adminAPIsDiscoverer,
		adminAPIClientsFactory,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to build kong api client(s): %w", err)
	}

	// Get Kong configuration root(s) to validate them and extract Kong's version.
	kongRoots, err := kongconfig.GetRoots(ctx, setupLog, c.KongAdminInitializationRetries, c.KongAdminInitializationRetryDelay, initialKongClients)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve Kong admin root(s): %w", err)
	}

	kongStartUpConfig, err := kongconfig.ValidateRoots(kongRoots, c.SkipCACertificates)
	if err != nil {
		return nil, fmt.Errorf("could not validate Kong admin root(s) configuration: %w", err)
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
		SanitizeKonnectConfigDumps:    c.FeatureGates.Enabled(managercfg.SanitizeKonnectConfigDumpsFeature),
		FallbackConfiguration:         c.FeatureGates.Enabled(managercfg.FallbackConfigurationFeature),
		UseLastValidConfigForFallback: c.UseLastValidConfigForFallback,
	}

	setupLog.Info("Configuring and building the controller manager")
	managerOpts := setupManagerOptions(ctx, setupLog, &c, dbMode)

	mgr, err := ctrl.NewManager(m.kubeconfig, managerOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to create controller manager: %w", err)
	}
	m.m = mgr

	if err := waitForKubernetesAPIReadiness(ctx, setupLog, mgr); err != nil {
		return nil, fmt.Errorf("unable to connect to Kubernetes API: %w", err)
	}

	setupLog.Info("Initializing Dataplane Client")
	var eventRecorder record.EventRecorder
	if c.EmitKubernetesEvents {
		setupLog.Info("Emitting Kubernetes events enabled, creating an event recorder for " + consts.KongClientEventRecorderComponentName)
		eventRecorder = newEventRecorderForInstance(instanceID, mgr.GetEventRecorderFor(consts.KongClientEventRecorderComponentName))
	} else {
		setupLog.Info("Emitting Kubernetes events disabled, discarding all events")
		// Create an empty record.FakeRecorder with no Events channel to discard all events.
		eventRecorder = &record.FakeRecorder{}
	}

	readinessChecker := clients.NewDefaultReadinessChecker(adminAPIClientsFactory, c.GatewayDiscoveryReadinessCheckTimeout, setupLog.WithName("readiness-checker"))
	clientsManager, err := clients.NewAdminAPIClientsManager(
		ctx,
		initialKongClients,
		readinessChecker,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AdminAPIClientsManager: %w", err)
	}
	clientsManager = clientsManager.WithDBMode(dbMode)
	clientsManager = clientsManager.WithReconciliationInterval(c.GatewayDiscoveryReadinessCheckInterval)
	m.clientsManager = clientsManager

	supportRedirectPlugin := kongSemVersion.GTE(versions.KongRedirectPluginCutoff)
	translatorFeatureFlags := translator.NewFeatureFlags(
		c.FeatureGates,
		routerFlavor,
		c.UpdateStatus,
		kongStartUpConfig.Version.IsKongGatewayEnterprise(),
		supportRedirectPlugin,
		c.CombinedServicesFromDifferentHTTPRoutes,
	)

	referenceIndexers := ctrlref.NewCacheIndexers(setupLog.WithName("reference-indexers"))
	cache := store.NewCacheStores()
	storer := store.New(cache, c.IngressClassName, logger)

	configTranslator, err := translator.NewTranslator(logger, storer, c.KongWorkspace, kongSemVersion, translatorFeatureFlags, NewSchemaServiceGetter(clientsManager),
		translator.Config{
			ClusterDomain:      c.ClusterDomain,
			EnableDrainSupport: c.EnableDrainSupport,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create translator: %w", err)
	}

	setupLog.Info("Setting up admission server")
	if err := m.setupAdmissionServer(ctx, referenceIndexers, translatorFeatureFlags, storer); err != nil {
		return nil, err
	}

	updateStrategyResolver := sendconfig.NewDefaultUpdateStrategyResolver(kongConfig, logger)
	configurationChangeDetector := sendconfig.NewKongGatewayConfigurationChangeDetector(logger)
	kongConfigFetcher := configfetcher.NewDefaultKongLastGoodConfigFetcher(translatorFeatureFlags.FillIDs, c.KongWorkspace)
	fallbackConfigGenerator := fallback.NewGenerator(fallback.NewDefaultCacheGraphProvider(), logger)
	metricsRecorder := metrics.NewGlobalCtrlRuntimeMetricsRecorder(instanceID)

	var dataplaneClientOpts []dataplane.KongClientOption
	if dc, ok := diagnosticsClient.Get(); ok {
		dataplaneClientOpts = append(dataplaneClientOpts, dataplane.WithDiagnosticsClient(dc))
	}
	dataplaneClient, err := dataplane.NewKongClient(
		logger,
		time.Duration(c.ProxyTimeoutSeconds*float32(time.Second)),
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
		metricsRecorder,
		dataplaneClientOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kong data-plane client: %w", err)
	}

	setupLog.Info("Initializing Dataplane Synchronizer")
	synchronizer, err := setupDataplaneSynchronizer(logger, mgr, dataplaneClient, c.ProxySyncSeconds, c.InitCacheSyncDuration)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize dataplane synchronizer: %w", err)
	}
	m.synchronizer = synchronizer

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
		return nil, err
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
		c.FeatureGates,
		clientsManager,
		adminAPIsDiscoverer,
	)
	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return nil, fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}

	// BUG: kubebuilder (at the time of writing - 3.0.0-rc.1) does not allow this tag anywhere else than main.go
	// See https://github.com/kubernetes-sigs/kubebuilder/issues/932
	// +kubebuilder:scaffold:builder

	if c.Konnect.ConfigSynchronizationEnabled {
		// In case of failures when building Konnect related objects, we're not returning errors as Konnect is not
		// considered critical feature, and it should not break the basic functionality of the controller.

		// Set up a config status notifier to be used by Konnect related components. Register it with the data plane
		// client so it can send status updates to Konnect.
		configStatusNotifier := clients.NewChannelConfigNotifier(logger)
		dataplaneClient.SetConfigStatusNotifier(configStatusNotifier)

		// Setup Konnect ConfigSynchronizer with manager.
		konnectConfigSynchronizer, err := setupKonnectConfigSynchronizerWithMgr(
			ctx,
			mgr,
			c,
			kongConfig,
			updateStrategyResolver,
			configStatusNotifier,
			metricsRecorder,
		)
		if err != nil {
			setupLog.Error(err, "Failed to setup Konnect configuration synchronizer with manager, skipping")
		} else {
			dataplaneClient.SetKonnectKongStateUpdater(konnectConfigSynchronizer)
		}

		// Setup Konnect NodeAgent with manager.
		if err := setupKonnectNodeAgentWithMgr(
			c,
			mgr,
			configStatusNotifier,
			clientsManager,
			setupLog,
			instanceID,
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
		return nil, err
	}
	if licenseGetter != nil {
		m.licenseGetter = mo.Some(licenseGetter)
		setupLog.Info("Inject license getter to config translator",
			"license_getter_type", fmt.Sprintf("%T", licenseGetter))
		configTranslator.InjectLicenseGetter(licenseGetter)
		kongConfigFetcher.InjectLicenseGetter(licenseGetter)
	}

	setupLog.Info("Finished setting up the controller manager")
	return m, nil
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
func setupKonnectNodeAgentWithMgr(
	c managercfg.Config,
	mgr manager.Manager,
	configStatusSubscriber clients.ConfigStatusSubscriber,
	clientsManager *clients.AdminAPIClientsManager,
	logger logr.Logger,
	instanceID InstanceID,
) error {
	konnectNodesAPIClient, err := nodes.NewClient(c.Konnect)
	if err != nil {
		return fmt.Errorf("failed creating konnect client: %w", err)
	}
	agent := konnect.NewNodeAgent(
		resolveControllerHostnameForKonnect(logger),
		metadata.UserAgent(),
		c.Konnect.RefreshNodePeriod,
		logger,
		konnectNodesAPIClient,
		configStatusSubscriber,
		konnect.NewGatewayClientGetter(logger, clientsManager),
		clientsManager,
		instanceID,
	)
	if err := mgr.Add(agent); err != nil {
		return fmt.Errorf("failed adding konnect.NodeAgent runnable to the manager: %w", err)
	}
	return nil
}

// resolveControllerHostnameForKonnect resolves the hostname to be used by Konnect NodeAgent. It tries to get the pod
// name and namespace, and if it fails, it falls back to using the hostname. If that fails too, it generates a random
// UUID.
func resolveControllerHostnameForKonnect(logger logr.Logger) string {
	nn, err := util.GetPodNN()
	if err != nil {
		logger.Error(err, "Failed getting pod name and/or namespace, falling back to use hostname as node name in Konnect")
		hostname, err := os.Hostname()
		if err != nil {
			logger.Error(err, "Failed getting hostname, falling back to random UUID as node name in Konnect")
			return uuid.NewString()
		}
		return hostname
	}
	logger.WithValues("hostname", nn.String()).Info("Resolved controller hostname for Konnect")
	return nn.String()
}

// setupDiagnostics creates diagnostics components (collector, server) if enabled in the configuration and prepares
// them to be run by the manager. It returns a non-empty diagnostics.Provider if config dumps are enabled.
func (m *Manager) setupDiagnostics(
	ctx context.Context,
	c managercfg.Config,
) mo.Option[diagnostics.Client] {
	logger := ctrl.LoggerFrom(ctx)

	// If neither profiling nor config dumps are enabled, we don't need to setup diagnostics at all.
	if !c.EnableProfiling && !c.EnableConfigDumps {
		logger.Info("Diagnostics disabled")
		return mo.None[diagnostics.Client]()
	}

	var serverOpts []diagnostics.ServerOption
	// If config dumps are enabled, we need to create a diagnostics collector, setup an HTTP handler exposing its
	// diagnostics, and pass it to the server options so it's plugged in.
	if c.EnableConfigDumps {
		diagnosticsCollector := diagnostics.NewCollector(logger, c)
		m.diagnosticsCollector = mo.Some(diagnosticsCollector)
		m.diagnosticsHandler = mo.Some(diagnostics.NewConfigDiagnosticsHTTPHandler(diagnosticsCollector, c.DumpSensitiveConfig))
		serverOpts = append(serverOpts, diagnostics.WithConfigDiagnostics(m.diagnosticsHandler.MustGet()))
	}

	if !c.DisableRunningDiagnosticsServer {
		m.diagnosticsServer = mo.Some(diagnostics.NewServer(logger, diagnostics.ServerConfig{
			ProfilingEnabled:    c.EnableProfiling,
			DumpSensitiveConfig: c.DumpSensitiveConfig,
			ListenerPort:        c.DiagnosticServerPort,
		}, serverOpts...))
	}

	// If diagnosticsCollector is set, it means that config dumps are enabled and we should return a diagnostics.Client.
	if dc, ok := m.diagnosticsCollector.Get(); ok {
		return mo.Some(dc.Client())
	}

	return mo.None[diagnostics.Client]()
}

// Run starts the Kong Ingress Controller. It blocks until the context is cancelled.
// It should be called only once per Manager instance.
func (m *Manager) Run(ctx context.Context) error {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Starting manager")

	if ds, ok := m.diagnosticsServer.Get(); ok {
		go func() {
			logger.Info("Starting diagnostics server")
			if err := ds.Listen(ctx); err != nil {
				logger.Error(err, "Diagnostics server exited")
			}
		}()
	}

	if dc, ok := m.diagnosticsCollector.Get(); ok {
		go func() {
			logger.Info("Starting diagnostics collector")
			if err := dc.Start(ctx); err != nil {
				logger.Error(err, "Diagnostics collector exited")
			}
		}()
	}

	if m.cfg.KongAdminSvc.IsPresent() {
		logger.Info("Starting AdminAPIClientsManager loop")
		m.clientsManager.Run(ctx)
	}

	if s, ok := m.admissionServer.Get(); ok {
		go func() {
			logger.Info("Starting admission server")
			if err := s.Start(ctx); err != nil {
				logger.Error(err, "Admission server exited")
			}
		}()
	}

	return m.m.Start(ctx)
}

// IsReady checks if the controller manager is ready to manage resources.
// It's only valid to call this method after the controller manager has been started
// with method Run(ctx).
func (m *Manager) IsReady() error {
	select {
	// If we're elected as leader then report readiness based on the readiness
	// of dataplane synchronizer.
	case <-m.m.Elected():
		if !m.synchronizer.IsReady() {
			return errors.New("synchronizer not yet configured")
		}
		// Do not mark the pod ready if KIC is configured to synchronize license from Konnect but no license can be found.
		if m.cfg.Konnect.LicenseSynchronizationEnabled {
			licenseGetter, present := m.licenseGetter.Get()
			if !present {
				return errors.New("Konnect license getter not present")
			}
			if !licenseGetter.GetLicense().IsPresent() {
				return errors.New("No Konnect license available")
			}
		}
	// If we're not the leader then just report as ready.
	default:
	}
	return nil
}

// DiagnosticsHandler returns the diagnostics HTTP handler if it's enabled in the configuration. Otherwise, it returns nil.
func (m *Manager) DiagnosticsHandler() http.Handler {
	if h, ok := m.diagnosticsHandler.Get(); ok {
		return h
	}
	return nil
}

// GetKubeconfig returns the Kubernetes client configuration used by the manager.
func (m *Manager) GetKubeconfig() *rest.Config {
	return m.kubeconfig
}

// GetClientsManager returns the AdminAPIClientsManager used by the manager.
// TODO: It is used by telemetry to calculate kongVersion, DB mode and router
// flavor. It shouldn't be exposed so wide, only relevant stuff.
func (m *Manager) GetClientsManager() *clients.AdminAPIClientsManager {
	return m.clientsManager
}
