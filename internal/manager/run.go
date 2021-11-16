// Package manager implements the controller manager for all controllers
package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metadata"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/mgrutils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *Config, diagnostic util.ConfigDumpDiagnostic) error {
	deprecatedLogger, logger, err := setupLoggers(c)
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
	kongConfig, err := setupKongConfig(ctx, setupLog, c)
	if err != nil {
		return fmt.Errorf("unable to build the kong admin api configuration: %w", err)
	}

	setupLog.Info("configuring and building the controller manager")
	controllerOpts, err := setupControllerOptions(setupLog, c, scheme)
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

	setupLog.Info("Starting Proxy Cache Server")
	proxy, err := setupProxyServer(ctx, setupLog, deprecatedLogger, mgr, kongConfig, diagnostic, c)
	if err != nil {
		return fmt.Errorf("unable to start proxy cache server: %w", err)
	}

	setupLog.Info("Starting Enabled Controllers")
	controllers, err := setupControllers(mgr, proxy, c, featureGates)
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
		if !proxy.IsReady() {
			return errors.New("proxy not yet configured")
		}
		return nil
	}); err != nil {
		return fmt.Errorf("unable to setup readyz: %w", err)
	}

	if c.AnonymousReports {
		setupLog.Info("Starting anonymous reports")
		if err := mgrutils.RunReport(ctx, kubeconfig, kongConfig, metadata.Release, featureGates); err != nil {
			setupLog.Error(err, "anonymous reporting failed")
		}
	} else {
		setupLog.Info("anonymous reports disabled, skipping")
	}

	if c.UpdateStatus {
		setupLog.Info("Starting resource status updater")
		go ctrlutils.PullConfigUpdate(ctx, kongConfig, logger, kubeconfig, c.PublishService, c.PublishStatusAddress)
	} else {
		setupLog.Info("WARNING: status updates were disabled, resources like Ingress objects will not receive updates to their statuses.")
	}

	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
}
