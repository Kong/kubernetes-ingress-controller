/*
Copyright 2015 The Kubernetes Authors.

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
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	discovery "k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kong "github.com/kong/kubernetes-ingress-controller/internal/apis/admin"
	consumerintscheme "github.com/kong/kubernetes-ingress-controller/internal/client/plugin/clientset/versioned/scheme"
	pluginintscheme "github.com/kong/kubernetes-ingress-controller/internal/client/plugin/clientset/versioned/scheme"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller"
	"github.com/kong/kubernetes-ingress-controller/internal/k8s"
	"github.com/kong/kubernetes-ingress-controller/version"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println(version.String())

	showVersion, conf, err := parseFlags()
	if showVersion {
		os.Exit(0)
	}

	if err != nil {
		glog.Fatal(err)
	}

	kubeCfg, kubeClient, err := createApiserverClient(conf.APIServerHost, conf.KubeConfigFile)
	if err != nil {
		handleFatalInitError(err)
	}

	// Add types to the default Kubernetes Scheme
	pluginintscheme.AddToScheme(scheme.Scheme)
	consumerintscheme.AddToScheme(scheme.Scheme)

	ns, name, err := k8s.ParseNameNS(conf.DefaultService)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = kubeClient.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "cannot get services in the namespace") {
			glog.Fatalf("✖ It seems the cluster it is running with Authorization enabled (like RBAC) and there is no permissions for the ingress controller. Please check the configuration")
		}
		glog.Fatalf("no service with name %v found: %v", conf.DefaultService, err)
	}
	glog.Infof("validated %v as the default backend", conf.DefaultService)

	if conf.PublishService == "" {
		glog.Fatal("flag --publish-address is mandatory")
	}

	ns, name, err = k8s.ParseNameNS(conf.PublishService)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = kubeClient.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		glog.Fatalf("unexpected error getting information about service %v: %v", conf.PublishService, err)
	}

	if conf.Namespace != "" {
		_, err = kubeClient.CoreV1().Namespaces().Get(conf.Namespace, metav1.GetOptions{})
		if err != nil {
			glog.Fatalf("no namespace with name %v found: %v", conf.Namespace, err)
		}
	}

	if conf.ResyncPeriod.Seconds() < 10 {
		glog.Fatalf("resync period (%vs) is too low", conf.ResyncPeriod.Seconds())
	}

	conf.KubeClient = kubeClient
	conf.KubeConf = kubeCfg

	kongClient, err := kong.NewRESTClient(&rest.Config{
		Host:    conf.Kong.URL,
		Timeout: 0,
	})
	if err != nil {
		glog.Fatalf("Error creating Kong Rest client: %v", err)
	}

	v, err := kongClient.GetVersion()
	if err != nil {
		glog.Fatalf("%v", err)
	}

	if !(v.GTE(semver.MustParse("0.13.0")) || v.GTE(semver.MustParse("0.32.0"))) {
		glog.Fatalf("The version %s is not compatible with the Kong Ingress Controller. It requires Kong CE 0.13.0 or higher, or Kong EE 0.32 or higher.", v)
	}

	glog.Infof("kong version: %s", v)
	conf.Kong.Client = kongClient

	ngx := controller.NewNGINXController(conf)

	go handleSigterm(ngx, func(code int) {
		os.Exit(code)
	})

	mux := http.NewServeMux()
	go registerHandlers(conf.EnableProfiling, 10254, ngx, mux)

	ngx.Start()
}

type exiter func(code int)

func handleSigterm(ngx *controller.NGINXController, exit exiter) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	glog.Infof("Received SIGTERM, shutting down")

	exitCode := 0
	if err := ngx.Stop(); err != nil {
		glog.Infof("Error during shutdown %v", err)
		exitCode = 1
	}

	glog.Infof("Handled quit, awaiting pod deletion")
	time.Sleep(10 * time.Second)

	glog.Infof("Exiting with %v", exitCode)
	exit(exitCode)
}

// createApiserverClient creates new Kubernetes Apiserver client. When kubeconfig or apiserverHost param is empty
// the function assumes that it is running inside a Kubernetes cluster and attempts to
// discover the Apiserver. Otherwise, it connects to the Apiserver specified.
//
// apiserverHost param is in the format of protocol://address:port/pathPrefix, e.g.http://localhost:8001.
// kubeConfig location of kubeconfig file
func createApiserverClient(apiserverHost string, kubeConfig string) (*rest.Config, *kubernetes.Clientset, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(apiserverHost, kubeConfig)
	if err != nil {
		return nil, nil, err
	}

	cfg.QPS = defaultQPS
	cfg.Burst = defaultBurst

	// cfg.ContentType = "application/vnd.kubernetes.protobuf"

	glog.Infof("Creating API client for %s", cfg.Host)

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	var v *discovery.Info

	// In some environments is possible the client cannot connect the API server in the first request
	// https://github.com/kubernetes/ingress-nginx/issues/1968
	defaultRetry := wait.Backoff{
		Steps:    10,
		Duration: 1 * time.Second,
		Factor:   1.5,
		Jitter:   0.1,
	}

	var lastErr error
	retries := 0
	glog.V(2).Info("trying to discover Kubernetes version")
	err = wait.ExponentialBackoff(defaultRetry, func() (bool, error) {
		v, err = client.Discovery().ServerVersion()

		if err == nil {
			return true, nil
		}

		lastErr = err
		glog.V(2).Infof("unexpected error discovering Kubernetes version (attempt %v): %v", err, retries)
		retries++
		return false, nil
	})

	// err is not null only if there was a timeout in the exponential backoff (ErrWaitTimeout)
	if err != nil {
		return nil, nil, lastErr
	}

	// this should not happen, warn the user
	if retries > 0 {
		glog.Warningf("it was required to retry %v times before reaching the API server", retries)
	}

	glog.Infof("Running in Kubernetes Cluster version v%v.%v (%v) - git (%v) commit %v - platform %v",
		v.Major, v.Minor, v.GitVersion, v.GitTreeState, v.GitCommit, v.Platform)

	return cfg, client, nil
}

const (
	// High enough QPS to fit all expected use cases. QPS=0 is not set here, because
	// client code is overriding it.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases. Burst=0 is not set here, because
	// client code is overriding it.
	defaultBurst = 1e6

	fakeCertificate = "default-fake-certificate"
)

/**
 * Handles fatal init error that prevents server from doing any work. Prints verbose error
 * message and quits the server.
 */
func handleFatalInitError(err error) {
	glog.Fatalf("Error while initializing connection to Kubernetes apiserver. "+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service accounts configuration). Reason: %s\n"+
		"Refer to the troubleshooting guide for more information: "+
		"https://github.com/kubernetes/ingress-nginx/blob/master/docs/troubleshooting.md", err)
}

func registerHandlers(enableProfiling bool, port int, ic *controller.NGINXController, mux *http.ServeMux) {
	// expose health check endpoint (/healthz)
	healthz.InstallHandler(mux,
		healthz.PingHealthz,
	)

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(version.String())
		w.Write(b)
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			glog.Errorf("unexpected error: %v", err)
		}
	})

	if enableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/heap", pprof.Index)
		mux.HandleFunc("/debug/pprof/mutex", pprof.Index)
		mux.HandleFunc("/debug/pprof/goroutine", pprof.Index)
		mux.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
		mux.HandleFunc("/debug/pprof/block", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	glog.Fatal(server.ListenAndServe())
}
