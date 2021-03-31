/*
Copyright 2021 Kong, Inc.Kong, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kong/go-kong/kong"
	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	kongctrl "github.com/kong/kubernetes-ingress-controller/railgun/controllers/configuration"
	//+kubebuilder:scaffold:imports
)

//go:generate go run github.com/kong/kubernetes-ingress-controller/railgun/cmd/generators/controllers/networking

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

type Config struct {
	MetricsAddr          string
	EnableLeaderElection bool
	ProbeAddr            string
	KongURL              string
	FilterTag            string
	Concurrency          int
	SecretNamespacedName types.NamespacedName

	ZapOptions zap.Options
}

var rootConfig Config

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	utilruntime.Must(configurationv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	initFlags()
}

func initFlags() {
	rootCmd.Flags().StringVar(&rootConfig.MetricsAddr,
		"metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	rootCmd.Flags().StringVar(&rootConfig.ProbeAddr,
		"health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	rootCmd.Flags().BoolVar(&rootConfig.EnableLeaderElection,
		"leader-elect", false, "Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	rootCmd.Flags().StringVar(&rootConfig.KongURL,
		"kong-url", "http://localhost:8001", "TODO")
	rootCmd.Flags().StringVar(&rootConfig.FilterTag,
		"kong-filter-tag", "managed-by-railgun", "TODO")
	rootCmd.Flags().IntVar(&rootConfig.Concurrency,
		"kong-concurrency", 10, "TODO")
	rootCmd.Flags().StringVar(&rootConfig.SecretNamespacedName.Name,
		"secret-name", "kong-config", "TODO")
	rootCmd.Flags().StringVar(&rootConfig.SecretNamespacedName.Namespace,
		"secret-namespace", "kong-system", "TODO")

	zapFlags := flag.NewFlagSet("", flag.ExitOnError)
	rootConfig.ZapOptions.BindFlags(zapFlags)
	rootCmd.Flags().AddGoFlagSet(zapFlags)
}

var rootCmd = &cobra.Command{
	Use:   "railgun",
	Short: "Kubernetes Ingress Controller (Railgun build)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runControllerManager(&rootConfig)
	},

	SilenceUsage:  true,
	SilenceErrors: true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		setupLog.Error(err, "rootCmd failed")
		os.Exit(1)
	}
}

func runControllerManager(config *Config) error {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&rootConfig.ZapOptions)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "5b374a9e.konghq.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	/* TODO: re-enable once fixed
	if err = (&kongctrl.KongIngressReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongIngress")
		os.Exit(1)
	}
	if err = (&kongctrl.KongClusterPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongClusterPlugin")
		os.Exit(1)
	}
	if err = (&kongctrl.KongPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongPlugin")
		os.Exit(1)
	}
	if err = (&kongctrl.KongConsumerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongConsumer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongConsumer")
		os.Exit(1)
	}
	*/

	kongClient, err := kong.NewClient(&config.KongURL, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("unable to create kongClient: %w", err)
	}

	if err = (&kongctrl.SecretReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Secret"),
		Scheme: mgr.GetScheme(),
		Params: kongctrl.SecretReconcilerParams{
			WatchNamespacedName: &config.SecretNamespacedName,
			KongConfig: sendconfig.Kong{
				URL:         config.KongURL,
				FilterTags:  []string{config.FilterTag},
				Concurrency: config.Concurrency,
				Client:      kongClient,
			},
		},
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create the secret controller: %w", err)
	}

	// TODO - we've got a couple places in here and below where we "short circuit" controllers if the relevant API isn't available.
	// This is convenient for testing, but maintainers should reconsider this before we release KIC 2.0.
	// SEE: https://github.com/Kong/kubernetes-ingress-controller/issues/1101
	if err := kongctrl.SetupIngressControllers(mgr); err != nil {
		return fmt.Errorf("unable to create the ingress controller: %w", err)
	}

	// TODO - similar to above, we're short circuiting here. It's convenient, but let's discuss if this is what we want ultimately.
	// SEE: https://github.com/Kong/kubernetes-ingress-controller/issues/1101
	udpIngressAvailable, err := kongctrl.IsAPIAvailable(mgr, &v1alpha1.UDPIngress{})
	if !udpIngressAvailable {
		setupLog.Error(err, "API configuration.konghq.com/v1alpha1/UDPIngress is not available, skipping controller")
	} else {
		if err = (&kongctrl.KongV1UDPIngressReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("UDPIngress"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create the udpingress controller: %w", err)
		}
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}
