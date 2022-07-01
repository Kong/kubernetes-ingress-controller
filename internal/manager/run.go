// Package manager implements the controller manager for all controllers
package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/kong/go-kong/kong"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/metadata"
	mgrutils "github.com/kong/kubernetes-ingress-controller/v2/internal/manager/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *Config, diagnostic util.ConfigDumpDiagnostic) error {
	deprecatedLogger, _, err := setupLoggers(c)
	if err != nil {
		return err
	}
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("starting controller manager", "release", metadata.Release, "repo", metadata.Repo, "commit", metadata.Commit)
	setupLog.V(util.DebugLevel).Info("the ingress class name has been set", "value", c.IngressClassName)
	setupLog.V(util.DebugLevel).Info("building the manager runtime scheme and loading apis into the scheme")
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1beta1.AddToScheme(scheme))
	utilruntime.Must(knativev1alpha1.AddToScheme(scheme))
	utilruntime.Must(gatewayv1alpha2.AddToScheme(scheme))

	if c.EnableLeaderElection {
		setupLog.V(0).Info("the --leader-elect flag is deprecated and no longer has any effect: leader election is set based on the Kong database setting")
	}

	setupLog.Info("getting enabled options and features")
	featureGates, err := setupFeatureGates(setupLog, c)
	if err != nil {
		return fmt.Errorf("failed to configure feature gates: %w", err)
	}

	setupLog.Info("getting the kubernetes client configuration")
	kubeconfig, err := c.GetKubeconfig()
	if err != nil {
		return fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
	}

	setupLog.Info("getting the kong admin api client configuration")
	adminClient, err := c.GetKongClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to build kong api client: %w", err)
	}

	var kongRoot map[string]interface{}
	err = retry.Do(
		func() error {
			kongRoot, err = adminClient.Root(ctx)
			return err
		},
		retry.Attempts(c.KongAdminInitializationRetries),
		retry.Delay(c.KongAdminInitializationRetryDelay),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			setupLog.Info("Retrying kong admin api client call after error",
				"retries", fmt.Sprintf("%d/%d", n, c.KongAdminInitializationRetries),
				"error", err.Error(),
			)
		}),
	)

	if err != nil {
		return fmt.Errorf("could not retrieve Kong admin root: %w", err)
	}

	kongConfig := setupKongConfig(ctx, adminClient, setupLog, c)
	kongVersion, err := kong.ParseSemanticVersion(kong.VersionFromInfo(kongRoot))
	if err != nil {
		setupLog.V(util.WarnLevel).Info("could not parse Kong version, version-specific behavior disabled", "error", err)
	} else {
		util.SetKongVersion(kongVersion)
	}
	kongRootConfig, ok := kongRoot["configuration"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid root configuration, expected a map[string]interface{} got %T",
			kongRoot["configuration"])
	}
	dbmode, ok := kongRootConfig["database"].(string)
	if !ok {
		return fmt.Errorf("invalid database configuration, expected a string got %T", kongRootConfig["database"])
	}
	if dbmode == "off" && c.SkipCACertificates {
		return fmt.Errorf("--skip-ca-certificates is not available for use with DB-less Kong instances")
	}

	setupLog.Info("configuring and building the controller manager")
	controllerOpts, err := setupControllerOptions(setupLog, c, scheme, dbmode)
	if err != nil {
		return fmt.Errorf("unable to setup controller options: %w", err)
	}
	mgr, err := ctrl.NewManager(kubeconfig, controllerOpts)
	if err != nil {
		return fmt.Errorf("unable to start controller manager: %w", err)
	}

	setupLog.Info("Starting Admission Server")
	if err := setupAdmissionServer(ctx, c, mgr.GetClient()); err != nil {
		return err
	}

	setupLog.Info("Initializing Dataplane Client")
	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxyTimeoutSeconds))
	if err != nil {
		return fmt.Errorf("%f is not a valid number of seconds to the timeout config for the kong client: %w", c.ProxyTimeoutSeconds, err)
	}
	dataplaneClient, err := dataplane.NewKongClient(deprecatedLogger, timeoutDuration, c.IngressClassName, c.EnableReverseSync, c.SkipCACertificates, diagnostic, kongConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize kong data-plane client: %w", err)
	}

	setupLog.Info("Initializing Dataplane Synchronizer")
	synchronizer, err := setupDataplaneSynchronizer(setupLog, deprecatedLogger, mgr, dataplaneClient, c)
	if err != nil {
		return fmt.Errorf("unable to initialize dataplane synchronizer: %w", err)
	}

	if enabled, ok := featureGates[combinedRoutesFeature]; ok && enabled {
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
	dataplaneAddressFinder, err := setupDataplaneAddressFinder(ctx, mgr.GetClient(), c)
	if err != nil {
		return err
	}

	setupLog.Info("Starting Enabled Controllers")
	controllers, err := setupControllers(mgr, dataplaneClient, dataplaneAddressFinder, kubernetesStatusQueue, c, featureGates)
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
	//+kubebuilder:scaffold:builder

	setupLog.Info("Starting health check servers")
	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return fmt.Errorf("unable to setup healthz: %w", err)
	}
	if err := mgr.AddReadyzCheck("check", func(_ *http.Request) error {
		if !synchronizer.IsReady() {
			return errors.New("synchronizer not yet configured")
		}
		return nil
	}); err != nil {
		return fmt.Errorf("unable to setup readyz: %w", err)
	}

	if c.AnonymousReports {
		setupLog.Info("Starting anonymous reports")
		if err := mgrutils.RunReport(ctx, kubeconfig, kongConfig, c.PublishService, metadata.Release, featureGates); err != nil {
			setupLog.Error(err, "anonymous reporting failed")
		}
	} else {
		setupLog.Info("anonymous reports disabled, skipping")
	}

	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
}
