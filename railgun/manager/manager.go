// Package manager implements the controller manager for all controllers in Railgun.
package manager

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/adminapi"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	"github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	kongctrl "github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Config collects all configuration that the controller manager takes from the environment.
// BUG: the above is not 100% accurate today - controllers read some settings from environment variables directly
type Config struct {
	// See flag definitions in RegisterFlags(...) for documentation of the fields defined here.

	MetricsAddr          string
	EnableLeaderElection bool
	LeaderElectionID     string
	ProbeAddr            string
	KongURL              string
	FilterTag            string
	Concurrency          int
	SecretName           string
	SecretNamespace      string
	KubeconfigPath       string

	KongAdminAPIConfig adminapi.HTTPClientOpts

	ZapOptions zap.Options

	KongStateEnabled         util.EnablementStatus
	IngressExtV1beta1Enabled util.EnablementStatus
	IngressNetV1beta1Enabled util.EnablementStatus
	IngressNetV1Enabled      util.EnablementStatus
	UDPIngressEnabled        util.EnablementStatus
	TCPIngressEnabled        util.EnablementStatus
	KongIngressEnabled       util.EnablementStatus
	KongClusterPluginEnabled util.EnablementStatus
	KongPluginEnabled        util.EnablementStatus
	KongConsumerEnabled      util.EnablementStatus
}

// MakeFlagSetFor binds the provided Config to commandline flags.
func MakeFlagSetFor(c *Config) *pflag.FlagSet {
	flagSet := flagSet{*pflag.NewFlagSet("", pflag.ExitOnError)}

	flagSet.StringVar(&c.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flagSet.StringVar(&c.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flagSet.BoolVar(&c.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flagSet.StringVar(&c.LeaderElectionID, "election-id", "5b374a9e.konghq.com", `Election id to use for status update.`)
	flagSet.StringVar(&c.KongURL, "kong-url", "http://localhost:8001", "TODO")
	flagSet.StringVar(&c.FilterTag, "kong-filter-tag", "managed-by-railgun", "TODO")
	flagSet.IntVar(&c.Concurrency, "kong-concurrency", 10, "TODO")
	flagSet.StringVar(&c.SecretName, "secret-name", "kong-config", "TODO")
	flagSet.StringVar(&c.SecretNamespace, "secret-namespace", controllers.DefaultNamespace, "TODO")
	flagSet.StringVar(&c.KubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file.")

	flagSet.BoolVar(&c.KongAdminAPIConfig.TLSSkipVerify, "kong-admin-tls-skip-verify", false,
		"Disable verification of TLS certificate of Kong's Admin endpoint.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSServerName, "kong-admin-tls-server-name", "",
		"SNI name to use to verify the certificate presented by Kong in TLS.")
	flagSet.StringVar(&c.KongAdminAPIConfig.CACertPath, "kong-admin-ca-cert-file", "",
		`Path to PEM-encoded CA certificate file to verify
Kong's Admin SSL certificate.`)
	flagSet.StringVar(&c.KongAdminAPIConfig.CACert, "kong-admin-ca-cert", "",
		`PEM-encoded CA certificate to verify Kong's Admin SSL certificate.`)
	flagSet.StringSliceVar(&c.KongAdminAPIConfig.Headers, "kong-admin-header", nil,
		`add a header (key:value) to every Admin API call, this flag can be used multiple times to specify multiple headers`)

	const onOffUsage = "Can be one of [enabled, disabled]."
	flagSet.EnablementStatusVar(&c.KongStateEnabled, "controller-kongstate", util.EnablementStatusEnabled,
		"Enable or disable the KongState controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.IngressNetV1Enabled, "controller-ingress-networkingv1", util.EnablementStatusEnabled,
		"Enable or disable the Ingress controller (using API version networking.k8s.io/v1)."+onOffUsage)
	flagSet.EnablementStatusVar(&c.IngressNetV1beta1Enabled, "controller-ingress-networkingv1beta1", util.EnablementStatusDisabled,
		"Enable or disable the Ingress controller (using API version networking.k8s.io/v1beta1)."+onOffUsage)
	flagSet.EnablementStatusVar(&c.IngressExtV1beta1Enabled, "controller-ingress-extensionsv1beta1", util.EnablementStatusDisabled,
		"Enable or disable the Ingress controller (using API version extensions/v1beta1)."+onOffUsage)
	flagSet.EnablementStatusVar(&c.UDPIngressEnabled, "controller-udpingress", util.EnablementStatusDisabled,
		"Enable or disable the UDPIngress controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.TCPIngressEnabled, "controller-tcpingress", util.EnablementStatusDisabled,
		"Enable or disable the TCPIngress controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.KongIngressEnabled, "controller-kongingress", util.EnablementStatusDisabled,
		"Enable or disable the KongIngress controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.KongClusterPluginEnabled, "controller-kongclusterplugin", util.EnablementStatusDisabled,
		"Enable or disable the KongClusterPlugin controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.KongPluginEnabled, "controller-kongplugin", util.EnablementStatusDisabled,
		"Enable or disable the KongPlugin controller. "+onOffUsage)
	flagSet.EnablementStatusVar(&c.KongConsumerEnabled, "controller-kongconsumer", util.EnablementStatusDisabled,
		"Enable or disable the KongConsumer controller. "+onOffUsage)

	zapFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	c.ZapOptions.BindFlags(zapFlagSet)
	flagSet.AddGoFlagSet(zapFlagSet)

	return &flagSet.FlagSet
}

// Controller is a Kubernetes controller that can be plugged into Manager.
type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

// AutoHandler decides whether the specific controller shall be enabled (true) or disabled (false).
type AutoHandler func(client.Reader) bool

// ControllerDef is a specification of a Controller that can be conditionally registered with Manager.
type ControllerDef struct {
	IsEnabled   *util.EnablementStatus
	AutoHandler AutoHandler
	Controller  Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// MaybeSetupWithManager runs SetupWithManager on the controller if its EnablementStatus is either "enabled", or "auto"
// and AutoHandler says that it should be enabled.
func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	switch *c.IsEnabled {
	case util.EnablementStatusDisabled:
		return nil

	case util.EnablementStatusAuto:
		if c.AutoHandler == nil {
			return fmt.Errorf("'auto' enablement not supported for controller %q", c.Name())
		}

		if enable := c.AutoHandler(mgr.GetAPIReader()); !enable {
			return nil
		}
		fallthrough

	default: // controller enabled
		return c.Controller.SetupWithManager(mgr)
	}
}

// Run starts the controller manager and blocks until it exits.
func Run(ctx context.Context, c *Config) error {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&c.ZapOptions)))
	setupLog := ctrl.Log.WithName("setup")

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1alpha1.AddToScheme(scheme))
	utilruntime.Must(configurationv1beta1.AddToScheme(scheme))

	// TODO: we might want to change how this works in the future, rather than just assuming the default ns
	if v := os.Getenv(controllers.CtrlNamespaceEnv); v == "" {
		os.Setenv(controllers.CtrlNamespaceEnv, controllers.DefaultNamespace)
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", c.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("get kubeconfig from file %q: %w", c.KubeconfigPath, err)
	}

	mgr, err := ctrl.NewManager(kubeconfig, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         c.EnableLeaderElection,
		LeaderElectionID:       c.LeaderElectionID,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	httpclient, err := adminapi.MakeHTTPClient(&c.KongAdminAPIConfig)
	if err != nil {
		setupLog.Error(err, "cannot create a Kong Admin API client")
	}

	kongClient, err := kong.NewClient(&c.KongURL, httpclient)
	if err != nil {
		setupLog.Error(err, "unable to create kongClient")
		return err
	}

	controllers := []ControllerDef{
		{
			IsEnabled: &c.KongStateEnabled,
			Controller: &kongctrl.SecretReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Secret"),
				Scheme: mgr.GetScheme(),
				Params: kongctrl.SecretReconcilerParams{
					WatchName:      c.SecretName,
					WatchNamespace: c.SecretNamespace,
					KongConfig: sendconfig.Kong{
						URL:         c.KongURL,
						FilterTags:  []string{c.FilterTag},
						Concurrency: c.Concurrency,
						Client:      kongClient,
					},
				},
			},
		},
		{
			IsEnabled: &c.IngressNetV1Enabled,
			Controller: &configuration.NetV1IngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Ingress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.IngressNetV1beta1Enabled,
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Ingress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.IngressExtV1beta1Enabled,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("Ingress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.UDPIngressEnabled,
			Controller: &kongctrl.KongV1Alpha1UDPIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.TCPIngressEnabled,
			Controller: &kongctrl.KongV1Beta1TCPIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.KongIngressEnabled,
			Controller: &kongctrl.KongV1KongIngressReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.KongClusterPluginEnabled,
			Controller: &kongctrl.KongV1KongClusterPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.KongPluginEnabled,
			Controller: &kongctrl.KongV1KongPluginReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme: mgr.GetScheme(),
			},
		},
		{
			IsEnabled: &c.KongConsumerEnabled,
			Controller: &kongctrl.KongV1KongConsumerReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme: mgr.GetScheme(),
			},
		},
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

	setupLog.Info("starting manager")
	return mgr.Start(ctx)
}
