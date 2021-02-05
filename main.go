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
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	konghqcomv1 "github.com/kong/railgun/api/v1"
	"github.com/kong/railgun/controllers"
	"github.com/kong/railgun/controllers/configuration"
	"github.com/kong/railgun/controllers/inputs"
	//+kubebuilder:scaffold:imports
)

//go:generate go run github.com/kong/railgun/cmd/generators/controllers/inputs

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(konghqcomv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	var kongURL string
	var filterTag string
	var concurrency int
	var secretName string
	var secretNamespace string
	flag.StringVar(&kongURL, "kong-url", "http://localhost:8001", "TODO")
	flag.StringVar(&filterTag, "kong-filter-tag", "managed-by-railgun", "TODO")
	flag.IntVar(&concurrency, "kong-concurrency", 10, "TODO")
	flag.StringVar(&secretName, "secret-name", "kong-config", "TODO")
	flag.StringVar(&secretNamespace, "secret-namespace", "kube-system", "TODO")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if v := os.Getenv(controllers.CtrlNamespaceEnv); v == "" {
		setupLog.Error(fmt.Errorf("kong can not be configured because the required %s env var is not present", controllers.CtrlNamespaceEnv), "could not start controller manager")
		os.Exit(25)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "5b374a9e.my.domain",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	/* TODO: re-enable once fixed
	if err = (&inputs.KongIngressReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongIngress"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongIngress")
		os.Exit(1)
	}
	if err = (&inputs.KongClusterPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongClusterPlugin")
		os.Exit(1)
	}
	if err = (&inputs.KongPluginReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongPlugin"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongPlugin")
		os.Exit(1)
	}
	if err = (&inputs.KongConsumerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KongConsumer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KongConsumer")
		os.Exit(1)
	}
	*/

	kongClient, err := kong.NewClient(&kongURL, http.DefaultClient)
	if err != nil {
		setupLog.Error(err, "unable to create kongClient")
		os.Exit(1)
	}

	if err = (&configuration.SecretReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Secret"),
		Scheme: mgr.GetScheme(),
		Params: configuration.SecretReconcilerParams{
			WatchName:      secretName,
			WatchNamespace: secretNamespace,
			KongConfig: sendconfig.Kong{
				URL:         kongURL,
				FilterTags:  []string{filterTag},
				Concurrency: concurrency,
				Client:      kongClient,
			},
		},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Secret")
		os.Exit(1)
	}

	if err := inputs.SetupIngressControllers(mgr); err != nil {
		setupLog.Error(err, "unable to create controllers", "controllers", "Ingress")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
