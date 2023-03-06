// Package manager implements the controller manager for all controllers
package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/telemetry"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/utils/kongconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *Config, diagnostic util.ConfigDumpDiagnostic, deprecatedLogger logrus.FieldLogger) error {
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("starting controller manager", "release", metadata.Release, "repo", metadata.Repo, "commit", metadata.Commit)
	setupLog.Info("the ingress class name has been set", "value", c.IngressClassName)

	gateway.SetControllerName(gatewayv1beta1.GatewayController(c.GatewayAPIControllerName))

	setupLog.Info("getting enabled options and features")
	featureGates, err := featuregates.Setup(setupLog, c.FeatureGates)
	if err != nil {
		return fmt.Errorf("failed to configure feature gates: %w", err)
	}

	setupLog.Info("getting the kubernetes client configuration")
	kubeconfig, err := c.GetKubeconfig()
	if err != nil {
		return fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
	}

	setupLog.Info("starting standalone health check server")
	healthServer := &healthCheckServer{}
	healthServer.setHealthzCheck(healthz.Ping)
	healthServer.Start(ctx, c.ProbeAddr, setupLog.WithName("health-check"))

	setupLog.Info("getting the kong admin api client configuration")
	initialKongClients, err := c.adminAPIClients(ctx, setupLog.WithName("initialize-kong-clients"))
	if err != nil {
		return fmt.Errorf("unable to build kong api client(s): %w", err)
	}

	// Get Kong configuration root(s) to validate them and extract Kong's version.
	kongRoots, err := kongconfig.GetRoots(ctx, setupLog, c.KongAdminInitializationRetries, c.KongAdminInitializationRetryDelay, initialKongClients)
	if err != nil {
		return fmt.Errorf("could not retrieve Kong admin root(s): %w", err)
	}

	dbMode, v, err := kongconfig.ValidateRoots(kongRoots, c.SkipCACertificates)
	if err != nil {
		return fmt.Errorf("could not validate Kong admin root(s) configuration: %w", err)
	}

	err = c.ValidateGatewayDiscovery(dbMode)
	if err != nil {
		return err
	}

	semV := semver.Version{Major: v.Major(), Minor: v.Minor(), Patch: v.Patch()}
	versions.SetKongVersion(semV)

	kongConfig := sendconfig.Config{
		Version:            semV,
		InMemory:           (dbMode == "off") || (dbMode == ""),
		Concurrency:        c.Concurrency,
		FilterTags:         c.FilterTags,
		SkipCACertificates: c.SkipCACertificates,
	}
	kongConfig.Init(ctx, setupLog, initialKongClients)

	setupLog.Info("configuring and building the controller manager")
	controllerOpts, err := setupControllerOptions(setupLog, c, dbMode, featureGates)
	if err != nil {
		return fmt.Errorf("unable to setup controller options: %w", err)
	}

	mgr, err := ctrl.NewManager(kubeconfig, controllerOpts)
	if err != nil {
		return fmt.Errorf("unable to start controller manager: %w", err)
	}

	setupLog.Info("Initializing Dataplane Client")
	eventRecorder := mgr.GetEventRecorderFor(KongClientEventRecorderComponentName)

	clientsManager, err := dataplane.NewAdminAPIClientsManager(
		ctx,
		deprecatedLogger,
		initialKongClients,
		adminapi.NewClientFactoryForWorkspace(c.KongWorkspace, c.KongAdminAPIConfig, c.KongAdminToken),
	)
	if err != nil {
		return fmt.Errorf("failed to create AdminAPIClientsManager: %w", err)
	}
	if c.KongAdminSvc.IsPresent() {
		setupLog.Info("Running AdminAPIClientsManager notify loop")
		clientsManager.RunNotifyLoop()
	}

	setupLog.Info("Starting Admission Server")
	if err := setupAdmissionServer(ctx, c, clientsManager, mgr.GetClient(), deprecatedLogger); err != nil {
		return err
	}

	updateStrategyResolver := sendconfig.NewDefaultUpdateStrategyResolver(kongConfig, deprecatedLogger)
	configurationChangeDetector := sendconfig.NewDefaultClientConfigurationChangeDetector(deprecatedLogger)
	dataplaneClient, err := dataplane.NewKongClient(
		deprecatedLogger,
		time.Duration(c.ProxyTimeoutSeconds*float32(time.Second)),
		c.IngressClassName,
		diagnostic,
		kongConfig,
		eventRecorder,
		dbMode,
		clientsManager,
		updateStrategyResolver,
		configurationChangeDetector,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize kong data-plane client: %w", err)
	}

	setupLog.Info("Initializing Dataplane Synchronizer")
	synchronizer, err := setupDataplaneSynchronizer(setupLog, deprecatedLogger, mgr, dataplaneClient, c.ProxySyncSeconds)
	if err != nil {
		return fmt.Errorf("unable to initialize dataplane synchronizer: %w", err)
	}

	if enabled, ok := featureGates[featuregates.CombinedRoutesFeature]; ok && enabled {
		dataplaneClient.EnableCombinedServiceRoutes()
		setupLog.Info("combined routes mode has been enabled")
	}

	var kubernetesStatusQueue *status.Queue
	if c.UpdateStatus {
		setupLog.Info("Starting Status Updater")
		kubernetesStatusQueue = status.NewQueue()
		dataplaneClient.EnableKubernetesObjectReports(kubernetesStatusQueue)
	} else {
		setupLog.Info("status updates disabled, skipping status updater")
	}

	setupLog.Info("Initializing Dataplane Address Discovery")
	dataplaneAddressFinder, udpDataplaneAddressFinder, err := setupDataplaneAddressFinder(mgr.GetClient(), c, setupLog)
	if err != nil {
		return err
	}

	setupLog.Info("Starting Enabled Controllers")
	controllers, err := setupControllers(mgr, dataplaneClient,
		dataplaneAddressFinder, udpDataplaneAddressFinder, kubernetesStatusQueue, c, featureGates, clientsManager)
	if err != nil {
		return fmt.Errorf("unable to setup controller as expected %w", err)
	}
	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}

	// BUG: kubebuilder (at the time of writing - 3.0.0-rc.1) does not allow this tag anywhere else than main.go
	// See https://github.com/kubernetes-sigs/kubebuilder/issues/932
	// +kubebuilder:scaffold:builder

	// use standalone health check server instead of servers inside manager
	// because manager depends on initial kong clients:
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3592
	// After we implement the feature that manager does not depend on initial kong clients,
	// we should move back to the health check server inside the manager.
	/*
		setupLog.Info("Starting health check servers")
		if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
			return fmt.Errorf("unable to setup healthz: %w", err)
		}
		if err := mgr.AddReadyzCheck("check", readyzHandler(mgr, synchronizer)); err != nil {
			return fmt.Errorf("unable to setup readyz: %w", err)
		}
	*/

	setupLog.Info("Add readiness probe to health server")
	healthServer.setReadyzCheck(readyzHandler(mgr, synchronizer))

	if c.Konnect.ConfigSynchronizationEnabled {
		// In case of failures when building Konnect related objects, we're not returning errors as Konnect is not
		// considered critical feature, and it should not break the basic functionality of the controller.

		// Run the Konnect Admin API client initialization in a separate goroutine to not block while ensuring
		// connection.
		go setupKonnectAdminAPIClientWithClientsMgr(ctx, c.Konnect, clientsManager, setupLog)

		// Setup Konnect NodeAgent with manager.
		if err := setupKonnectNodeAgentWithMgr(
			c,
			mgr,
			dataplaneClient,
			clientsManager,
			setupLog,
		); err != nil {
			setupLog.Error(err, "Failed to setup Konnect NodeAgent with manager, skipping")
		}
	}

	if c.AnonymousReports {
		stopAnonymousReports, err := telemetry.SetupAnonymousReports(
			ctx,
			kubeconfig,
			clientsManager,
			telemetry.ReportValues{
				PublishServiceNN:               c.PublishService.OrEmpty(),
				FeatureGates:                   featureGates,
				MeshDetection:                  len(c.WatchNamespaces) == 0,
				KonnectSyncEnabled:             c.Konnect.ConfigSynchronizationEnabled,
				GatewayServiceDiscoveryEnabled: c.KongAdminSvc.IsPresent(),
			},
		)
		if err != nil {
			setupLog.Error(err, "failed setting up anonymous reports")
		} else {
			defer stopAnonymousReports()
		}
		setupLog.Info("anonymous reports enabled")
	} else {
		setupLog.Info("anonymous reports disabled, skipping")
	}

	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
}

// setupKonnectNodeAgentWithMgr creates and adds Konnect NodeAgent as the manager's Runnable.
// Returns error if failed to create Konnect NodeAgent.
func setupKonnectNodeAgentWithMgr(
	c *Config,
	mgr manager.Manager,
	dataplaneClient *dataplane.KongClient,
	clientsManager *dataplane.AdminAPIClientsManager,
	logger logr.Logger,
) error {
	konnectNodeAPIClient, err := konnect.NewNodeAPIClient(c.Konnect)
	if err != nil {
		return fmt.Errorf("failed creating konnect client: %w", err)
	}
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
	configStatusNotifier := dataplane.NewChannelConfigNotifier(logger)
	dataplaneClient.SetConfigStatusNotifier(configStatusNotifier)

	agent := konnect.NewNodeAgent(
		hostname,
		version,
		c.Konnect.RefreshNodePeriod,
		logger,
		konnectNodeAPIClient,
		configStatusNotifier,
		konnect.NewGatewayClientGetter(logger, clientsManager),
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
	clientsManager *dataplane.AdminAPIClientsManager,
	logger logr.Logger,
) {
	konnectAdminAPIClient, err := adminapi.NewKongClientForKonnectRuntimeGroup(config)
	if err != nil {
		logger.Error(err, "Failed creating Konnect Runtime Group Admin API client, skipping synchronisation")
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
