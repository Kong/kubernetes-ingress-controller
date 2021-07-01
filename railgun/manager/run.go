// Package manager implements the controller manager for all controllers in Railgun.
package manager

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"

	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	kongctrl "github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Setup & Run
// -----------------------------------------------------------------------------

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *config.Config) error {
	var deprecatedLogger logrus.FieldLogger
	var err error

	if v := os.Getenv("KONG_TEST_ENVIRONMENT"); v != "" {
		deprecatedLogger = util.MakeDebugLoggerWithReducedRedudancy(os.Stdout, &logrus.TextFormatter{}, 3, time.Second*30)
		deprecatedLogger.Info("detected that the controller is running in an automated testing environment: " +
			"log stifling has been enabled")
	} else {
		deprecatedLogger, err = util.MakeLogger(c.LogLevel, c.LogFormat)
		if err != nil {
			return fmt.Errorf("failed to make logger: %w", err)
		}
	}
	var logger logr.Logger = logrusr.NewLogger(deprecatedLogger)

	ctrl.SetLogger(logger)
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("starting controller manager", "release", Release, "repo", Repo, "commit", Commit)

	kubeconfig, err := c.GetKubeconfig()
	if err != nil {
		return fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
	}

	// set "kubernetes.io/ingress.class" to be used by controllers (defaults to "kong")
	setupLog.Info(`the ingress class name has been set`, "value", c.IngressClassName)

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1beta1.AddToScheme(scheme))
	utilruntime.Must(knativev1alpha1.AddToScheme(scheme))
	utilruntime.Must(apiextv1beta1.AddToScheme(scheme))

	controllerOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         c.EnableLeaderElection,
		LeaderElectionID:       c.LeaderElectionID,
	}

	// determine how to configure namespace watchers
	if strings.Contains(c.WatchNamespace, ",") {
		setupLog.Info("manager set up with multiple namespaces", "namespaces", c.WatchNamespace)
		// this mode does not set the Namespace option, so the manager will default to watching all namespaces
		// MultiNamespacedCacheBuilder imposes a filter on top of that watch to retrieve scoped resources
		// from the watched namespaces only.
		controllerOpts.NewCache = cache.MultiNamespacedCacheBuilder(strings.Split(c.WatchNamespace, ","))
	} else {
		controllerOpts.Namespace = c.WatchNamespace
	}

	// build the controller manager
	mgr, err := ctrl.NewManager(kubeconfig, controllerOpts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	kongClient, err := c.GetKongClient(ctx)
	if err != nil {
		setupLog.Error(err, "cannot create a Kong Admin API client")
		return err
	}

	var filterTags []string
	if ok, err := kongClient.Tags.Exists(ctx); err != nil {
		setupLog.Error(err, "tag filtering disabled because Kong Admin API does not support tags")
	} else if ok {
		setupLog.Info("tag filtering enabled", "tags", c.FilterTags)
		filterTags = c.FilterTags
	}

	// configure the kong client
	kongConfig := sendconfig.Kong{
		URL:               c.KongAdminURL,
		FilterTags:        filterTags,
		Concurrency:       c.Concurrency,
		Client:            kongClient,
		PluginSchemaStore: util.NewPluginSchemaStore(kongClient),
	}

	// determine the proxy synchronization strategy
	if c.ProxySyncSeconds < proxy.DefaultSyncSeconds {
		setupLog.Info(fmt.Sprintf("WARNING: --proxy-sync-seconds is configured for %fs, in DBLESS mode this may result in"+
			" problems of inconsistency in the proxy state. For DBLESS mode %fs+ is recommended (3s is the default).",
			c.ProxySyncSeconds, proxy.DefaultSyncSeconds,
		))
	}
	syncTickDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxySyncSeconds))
	if err != nil {
		setupLog.Error(err, "%s is not a valid number of seconds to stagger the proxy server synchronization")
		return err
	}

	// determine the proxy timeout
	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%gs", c.ProxyTimeoutSeconds))
	if err != nil {
		setupLog.Error(err, "%s is not a valid number of seconds to the timeout config for the kong client")
		return err
	}

	// start the proxy cache server
	prx, err := proxy.NewCacheBasedProxyWithStagger(ctx,
		// NOTE: logr-based loggers use the "logger" field instead of "subsystem". When replacing logrus with logr, replace
		// WithField("subsystem", ...) with WithName(...).
		deprecatedLogger.WithField("subsystem", "proxy-cache-resolver"),
		mgr.GetClient(),
		kongConfig,
		c.IngressClassName,
		c.EnableReverseSync,
		syncTickDuration,
		timeoutDuration,
		sendconfig.UpdateKongAdminSimple)
	if err != nil {
		setupLog.Error(err, "unable to start proxy cache server")
		return err
	}

	alwaysEnabled := util.EnablementStatusEnabled
	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------

		{
			IsEnabled: &c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme: mgr.GetScheme(),
				Proxy:  prx,
			},
		},
		{
			IsEnabled: &c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme: mgr.GetScheme(),
				Proxy:  prx,
			},
		},
		{
			IsEnabled: &alwaysEnabled,
			Controller: &configuration.CoreV1SecretReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Secrets"),
				Scheme: mgr.GetScheme(),
				Proxy:  prx,
			},
		},

		// ---------------------------------------------------------------------------
		// Kong API Controllers
		// ---------------------------------------------------------------------------
		{
			IsEnabled: &c.UDPIngressEnabled,
			Controller: &kongctrl.KongV1Beta1UDPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: &c.TCPIngressEnabled,
			Controller: &kongctrl.KongV1Beta1TCPIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
		{
			IsEnabled: &c.KongIngressEnabled,
			Controller: &kongctrl.KongV1KongIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme: mgr.GetScheme(),
				Proxy:  prx,
			},
		},
		{
			IsEnabled: &c.KongPluginEnabled,
			Controller: &kongctrl.KongV1KongPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme: mgr.GetScheme(),
				Proxy:  prx,
			},
		},
		{
			IsEnabled: &c.KongConsumerEnabled,
			Controller: &kongctrl.KongV1KongConsumerReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
	}

	// Negotiate Ingress version
	ingressControllers := map[IngressAPI]ControllerDef{
		NetworkingV1: {
			IsEnabled: &c.IngressNetV1Enabled,
			Controller: &configuration.NetV1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
		NetworkingV1beta1: {
			IsEnabled: &c.IngressNetV1beta1Enabled,
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
		ExtensionsV1beta1: {
			IsEnabled: &c.IngressExtV1beta1Enabled,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		},
	}

	negotiatedIngressAPI, err := negotiateIngressAPI(c, mgr.GetClient())
	if err == nil {
		controllers = append(controllers, ingressControllers[negotiatedIngressAPI])
	} else {
		setupLog.Info(`no Ingress controllers enabled or no suitable Ingress version found.
		Disabling Ingress controller`)
	}

	kongClusterPluginGVR := schema.GroupVersionResource{
		Group:    konghqcomv1.SchemeGroupVersion.Group,
		Version:  konghqcomv1.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	if ctrlutils.CRDExists(mgr.GetClient(), kongClusterPluginGVR) == true {
		setupLog.Info("kongclusterplugins.configuration.konghq.com v1beta1 CRD available on cluster.")
		controller := ControllerDef{
			IsEnabled: &c.KongClusterPluginEnabled,
			Controller: &kongctrl.KongV1KongClusterPluginReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:           mgr.GetScheme(),
				Proxy:            prx,
				IngressClassName: c.IngressClassName,
			},
		}
		controllers = append(controllers, controller)
	} else {
		setupLog.Info(`kongclusterplugins.configuration.konghq.com v1beta1 CRD not available on cluster.
		Disabling KongClusterPlugin controller`)
	}

	// -------------------------------------------------------------------------
	// Knative Ingress Handling
	// -------------------------------------------------------------------------

	// if knative.Ingress has been specifically flagged as disabled, then we do nothing with it.
	// otherwise we attempt to load the relevant controller, and if the controller can not be loaded
	// during setup (e.g. the CRD is not present on the cluster) we will load our CRD controller which
	// will dynamically load the controller if/when the CRD is loaded and the API becomes available.
	if c.KnativeIngressEnabled != util.EnablementStatusDisabled {
		knativeGVR := schema.GroupVersionResource{
			Group:    knativev1alpha1.SchemeGroupVersion.Group,
			Version:  knativev1alpha1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		}
		knativeWasLoaded := false
		if ctrlutils.CRDExists(mgr.GetClient(), knativeGVR) {
			setupLog.Info("ingresses.networking.internal.knative.dev v1alpha1 CRD available on cluster.")
			controller := ControllerDef{
				IsEnabled: &c.KnativeIngressEnabled,
				Controller: &kongctrl.Knativev1alpha1IngressReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
					Scheme:           mgr.GetScheme(),
					Proxy:            prx,
					IngressClassName: c.IngressClassName,
				},
			}
			controllers = append(controllers, controller)
			knativeWasLoaded = true
		}

		if c.DynamicControllerLoaderEnabled == util.EnablementStatusDisabled {
			// if the end-user explicitly disabled dynamic controller loading, then let them be but warn them that
			// a restart will be required in order for the relevant controller to be started.
			setupLog.Info(`WARNING: the dynamic controller loader was explicitly disabled: knative support can not be
				 dynamically added to the cluster and if the knative.Ingress CRD later becomes available, the manager
				 will need to be restarted in order to pick it up and start the relevant controller(s).`)

		} else if !knativeWasLoaded {
			// if knative is not explicitly disabled but was not previously loaded we will run a controller that will
			// watch the cluster CRDs to determine if/when that API is loaded and when it becomes available we will
			// dynamically load the knative.Ingress controller into the controller manager at runtime.
			setupLog.Info(fmt.Sprintf("knative support set to %s, but the API is not yet available (CRD not found)",
				c.KnativeIngressEnabled.String()))

			setupLog.Info("starting the CRD controller to watch whether knative.Ingress becomes available later")
			controllers = append(controllers, ControllerDef{
				IsEnabled: &c.KnativeIngressEnabled,
				Controller: &configuration.CustomResourceDefinitionReconciler{
					Mgr:              mgr,
					Client:           mgr.GetClient(),
					Log:              ctrl.Log.WithName("controllers").WithName("CustomResourceDefinition").WithName("extv1beta1"),
					Scheme:           mgr.GetScheme(),
					Proxy:            prx,
					IngressClassName: c.IngressClassName,
				},
			})
		}
	}

	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}

	// BUG: kubebuilder (at the time of writing - 3.0.0-rc.1) does not allow this tag anywhere else than main.go
	// See https://github.com/kubernetes-sigs/kubebuilder/issues/932
	//+kubebuilder:scaffold:builder

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

	setupLog.Info("starting manager")
	return mgr.Start(ctx)
}
