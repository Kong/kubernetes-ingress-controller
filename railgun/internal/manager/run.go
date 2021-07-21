// Package manager implements the controller manager for all controllers in Railgun.
package manager

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *Config) error {
	deprecatedLogger, logger, err := setupLoggers(c)
	if err != nil {
		return err
	}
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("starting controller manager", "release", Release, "repo", Repo, "commit", Commit)
	setupLog.Info("the ingress class name has been set", "value", c.IngressClassName)

	setupLog.Info("building the manager runtime scheme and loading apis into the scheme")
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1beta1.AddToScheme(scheme))
	utilruntime.Must(knativev1alpha1.AddToScheme(scheme))

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
	controllerOpts := setupControllerOptions(setupLog, c, scheme)
	mgr, err := ctrl.NewManager(kubeconfig, controllerOpts)
	if err != nil {
		return fmt.Errorf("unable to start controller manager: %w", err)
	}

	setupLog.Info("configuring and building the proxy cache server")
	proxy, err := setupProxyServer(ctx, setupLog, deprecatedLogger, mgr, kongConfig, c)
	if err != nil {
		return fmt.Errorf("unable to start proxy cache server: %w", err)
	}

	setupLog.Info("deploying all enabled controllers")
	controllers := setupControllers(setupLog, mgr, proxy, c)
	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}

	// BUG: kubebuilder (at the time of writing - 3.0.0-rc.1) does not allow this tag anywhere else than main.go
	// See https://github.com/kubernetes-sigs/kubebuilder/issues/932
	//+kubebuilder:scaffold:builder

	setupLog.Info("enabling health checks")
	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return fmt.Errorf("unable to setup healthz: %w", err)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		return fmt.Errorf("unable to setup readyz: %w", err)
	}

	if c.AnonymousReports {
		setupLog.Info("running anonymous reports")
		if err := mgrutils.RunReport(ctx, kubeconfig, kongConfig, Release); err != nil {
			setupLog.Error(err, "anonymous reporting failed")
		}
	} else {
		setupLog.Info("anonymous reports disabled, skipping")
	}

	if c.UpdateStatus {
		setupLog.Info("status updates enabled, status update routine is being started in the background.")
		go ctrlutils.PullConfigUpdate(ctx, kongConfig, logger, kubeconfig, c.PublishService, c.PublishStatusAddress)
	} else {
		setupLog.Info("WARNING: status updates were disabled, resources like Ingress objects will not receive updates to their statuses.")
	}

	setupLog.Info("starting manager")
	return mgr.Start(ctx)
}
