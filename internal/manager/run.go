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
	mgrutils "github.com/kong/kubernetes-ingress-controller/v2/internal/manager/utils"
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
	setupLog.Info("getting the kong admin api client configuration")
	initialKongClients, err := c.getKongClients(ctx)
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

	setupLog.Info("Starting Admission Server")
	if err := setupAdmissionServer(ctx, c, mgr.GetClient(), deprecatedLogger); err != nil {
		return err
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
	if c.KongAdminSvc.Name != "" {
		setupLog.Info("Running AdminAPIClientsManager notify loop")
		clientsManager.RunNotifyLoop()
	}

	dataplaneClient, err := dataplane.NewKongClient(
		deprecatedLogger,
		time.Duration(c.ProxyTimeoutSeconds*float32(time.Second)),
		c.IngressClassName,
		c.EnableReverseSync,
		c.SkipCACertificates,
		diagnostic,
		kongConfig,
		eventRecorder,
		dbMode,
		clientsManager,
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

	setupLog.Info("Starting health check servers")
	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return fmt.Errorf("unable to setup healthz: %w", err)
	}
	if err := mgr.AddReadyzCheck("check", readyzHandler(mgr, synchronizer)); err != nil {
		return fmt.Errorf("unable to setup readyz: %w", err)
	}

	if c.Konnect.ConfigSynchronizationEnabled {
		// In case of failures when building Konnect related objects, we're not returning errors as Konnect is not
		// considered critical feature, and it should not break the basic functionality of the controller.

		setupLog.Info("Start Konnect client to register runtime instances to Konnect")
		konnectNodeAPIClient, err := konnect.NewNodeAPIClient(c.Konnect)
		if err != nil {
			setupLog.Error(err, "failed creating konnect client, skipping running NodeAgent")
		} else {
			hostname, _ := os.Hostname()
			version := metadata.Release
			agent := konnect.NewNodeAgent(hostname, version, setupLog, konnectNodeAPIClient)
			agent.Run()
		}

		konnectAdminAPIClient, err := adminapi.NewKongClientForKonnectRuntimeGroup(ctx, c.Konnect)
		if err != nil {
			setupLog.Error(err, "failed creating Konnect Runtime Group Admin API client, skipping synchronisation")
		} else {
			setupLog.Info("Initialized Konnect Admin API client")
			clientsManager.SetKonnectClient(konnectAdminAPIClient)
		}
	}

	if c.AnonymousReports {
		setupLog.Info("Starting anonymous reports")
		// the argument checking the watch namespaces length enables or disables mesh detection. the mesh detect client
		// attempts to use all namespaces and can't utilize a manager multi-namespaced cache, so if we need to limit
		// namespace access we just disable mesh detection altogether.
		if err := mgrutils.RunReport(
			ctx,
			kubeconfig,
			c.PublishService.String(),
			metadata.Release,
			len(c.WatchNamespaces) == 0,
			featureGates,
			clientsManager,
		); err != nil {
			setupLog.Error(err, "anonymous reporting failed")
		}
	} else {
		setupLog.Info("anonymous reports disabled, skipping")
	}

	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
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
