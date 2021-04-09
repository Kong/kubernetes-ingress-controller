package manager

import (
	"net/http"
	"os"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	kongctrl "github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type Config struct {
	MetricsAddr          string
	EnableLeaderElection bool
	ProbeAddr            string
	KongURL              string
	FilterTag            string
	Concurrency          int
	SecretName           string
	SecretNamespace      string

	ZapOptions zap.Options
}

func Run(c *Config) error {
	setupLog := ctrl.Log.WithName("setup")

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1alpha1.AddToScheme(scheme))
	utilruntime.Must(configurationv1beta1.AddToScheme(scheme))

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&c.ZapOptions)))

	// TODO: we might want to change how this works in the future, rather than just assuming the default ns
	if v := os.Getenv(controllers.CtrlNamespaceEnv); v == "" {
		os.Setenv(controllers.CtrlNamespaceEnv, controllers.DefaultNamespace)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         c.EnableLeaderElection,
		LeaderElectionID:       "5b374a9e.konghq.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	/* TODO: re-enable once fixed
	if err = (&kongctrl.KongIngressReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongIngress")
		return err
	}
	if err = (&kongctrl.KongClusterPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongClusterPlugin")
		return err
	}
	if err = (&kongctrl.KongPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongPlugin")
		return err
	}
	if err = (&kongctrl.KongConsumerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongConsumer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongConsumer")
		return err
	}
	*/

	kongClient, err := kong.NewClient(&c.KongURL, http.DefaultClient)
	if err != nil {
		setupLog.Error(err, "unable to create kongClient")
		return err
	}

	if err = (&kongctrl.SecretReconciler{
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
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Secret")
		return err
	}

	// TODO - we've got a couple places in here and below where we "short circuit" controllers if the relevant API isn't available.
	// This is convenient for testing, but maintainers should reconsider this before we release KIC 2.0.
	// SEE: https://github.com/Kong/kubernetes-ingress-controller/issues/1101
	if err := kongctrl.SetupIngressControllers(mgr); err != nil {
		setupLog.Error(err, "unable to create controllers", "controllers", "Ingress")
		return err
	}

	// TODO - similar to above, we're short circuiting here. It's convenient, but let's discuss if this is what we want ultimately.
	// SEE: https://github.com/Kong/kubernetes-ingress-controller/issues/1101
	udpIngressAvailable, err := kongctrl.IsAPIAvailable(mgr, &v1alpha1.UDPIngress{})
	if !udpIngressAvailable {
		setupLog.Error(err, "API configuration.konghq.com/v1alpha1/UDPIngress is not available, skipping controller")
	} else {
		if err = (&kongctrl.KongV1Alpha1UDPIngressReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("UDPIngress"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "UDPIngress")
			return err
		}
	}

	tcpIngressAvailable, err := kongctrl.IsAPIAvailable(mgr, &configurationv1beta1.TCPIngress{})
	if !tcpIngressAvailable {
		setupLog.Error(err, "API configuration.konghq.com/v1alpha1/TCPIngress is not available, skipping controller")
	} else {
		if err = (&kongctrl.KongV1Beta1TCPIngressReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("TCPIngress"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "TCPIngress")
			return err
		}
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	setupLog.Info("starting manager")
	return mgr.Start(ctrl.SetupSignalHandler())
}
