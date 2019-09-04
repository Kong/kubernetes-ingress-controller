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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blang/semver"
	"github.com/eapache/channels"
	"github.com/golang/glog"
	"github.com/hashicorp/go-uuid"
	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/admission"
	configurationclientv1 "github.com/kong/kubernetes-ingress-controller/internal/client/configuration/clientset/versioned"
	configurationinformer "github.com/kong/kubernetes-ingress-controller/internal/client/configuration/informers/externalversions"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sVersion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println(version())

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

	if conf.PublishService == "" {
		glog.Fatal("flag --publish-address is mandatory")
	}

	ns, name, err := utils.ParseNameNS(conf.PublishService)
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

	defaultTransport := http.DefaultTransport.(*http.Transport)

	var tlsConfig tls.Config

	if conf.Kong.TLSSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	if conf.Kong.TLSServerName != "" {
		tlsConfig.ServerName = conf.Kong.TLSServerName
	}

	if conf.Kong.CACert != "" {
		certPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile(conf.Kong.CACert)
		if err != nil {
			glog.Fatalf("failed to read CACert: %s", conf.Kong.CACert)
		}
		ok := certPool.AppendCertsFromPEM([]byte(cert))
		if !ok {
			glog.Fatalf("failed to load CACert: %s", conf.Kong.CACert)
		}
		tlsConfig.RootCAs = certPool
	}
	defaultTransport.TLSClientConfig = &tlsConfig
	c := http.DefaultClient
	c.Transport = &HeaderRoundTripper{
		headers: conf.Kong.Headers,
		rt:      defaultTransport,
	}

	kongClient, err := kong.NewClient(kong.String(conf.Kong.URL), c)
	if err != nil {
		glog.Fatalf("Error creating Kong Rest client: %v", err)
	}

	root, err := kongClient.Root(nil)
	if err != nil {
		glog.Fatalf("%v", err)
	}
	v, err := getSemVerVer(root["version"].(string))

	if !(v.GTE(semver.MustParse("0.13.0")) || v.GTE(semver.MustParse("0.32.0"))) {
		glog.Fatalf("The version %s is not compatible with the Kong Ingress Controller. It requires Kong CE 0.13.0 or higher, or Kong EE 0.32 or higher.", v)
	}

	glog.Infof("kong version: %s", v)
	kongConfiguration := root["configuration"].(map[string]interface{})
	conf.Kong.Version = v
	kongDB := kongConfiguration["database"].(string)
	glog.Infof("Kong datastore: %s", kongDB)

	if kongDB == "off" {
		conf.Kong.InMemory = true
	}
	req, _ := http.NewRequest("GET",
		conf.Kong.URL+"/tags", nil)
	res, err := kongClient.Do(nil, req, nil)
	if err == nil && res.StatusCode == 200 {
		conf.Kong.HasTagSupport = true
	}

	// setup workspace in Kong Enterprise
	if conf.Kong.Workspace != "" {
		kongClient, err = kong.NewClient(kong.String(conf.Kong.URL+"/"+conf.Kong.Workspace), c)
		if err != nil {
			glog.Fatalf("Error creating Kong Rest client: %v", err)
		}
	}
	conf.Kong.Client = kongClient

	err = discovery.ServerSupportsVersion(kubeClient.Discovery(), schema.GroupVersion{
		Group:   "networking.k8s.io",
		Version: "v1beta1",
	})
	if err == nil {
		conf.UseNetworkingV1beta1 = true
	}
	coreInformerFactory := informers.NewSharedInformerFactoryWithOptions(
		kubeClient,
		conf.ResyncPeriod,
		informers.WithNamespace(conf.Namespace),
	)
	confClient, _ := configurationclientv1.NewForConfig(conf.KubeConf)
	kongInformerFactory := configurationinformer.NewSharedInformerFactoryWithOptions(
		confClient,
		conf.ResyncPeriod,
		configurationinformer.WithNamespace(conf.Namespace),
	)

	var synced []cache.InformerSynced
	updateChannel := channels.NewRingChannel(1024)
	reh := controller.ResourceEventHandler{
		UpdateCh:           updateChannel,
		IsValidIngresClass: annotations.IngressClassValidatorFunc(conf.IngressClass),
	}
	var informers []cache.SharedIndexInformer
	var cacheStores store.CacheStores

	var ingInformer cache.SharedIndexInformer
	if conf.UseNetworkingV1beta1 {
		ingInformer = coreInformerFactory.Networking().V1beta1().Ingresses().Informer()
	} else {
		ingInformer = coreInformerFactory.Extensions().V1beta1().Ingresses().Informer()
	}

	ingInformer.AddEventHandler(reh)
	cacheStores.Ingress = ingInformer.GetStore()
	informers = append(informers, ingInformer)

	endpointsInformer := coreInformerFactory.Core().V1().Endpoints().Informer()
	endpointsInformer.AddEventHandler(controller.EndpointsEventHandler{
		UpdateCh: updateChannel,
	})
	cacheStores.Endpoint = endpointsInformer.GetStore()
	informers = append(informers, endpointsInformer)

	secretsInformer := coreInformerFactory.Core().V1().Secrets().Informer()
	secretsInformer.AddEventHandler(reh)
	cacheStores.Secret = secretsInformer.GetStore()
	informers = append(informers, secretsInformer)

	servicesInformer := coreInformerFactory.Core().V1().Services().Informer()
	servicesInformer.AddEventHandler(reh)
	cacheStores.Service = servicesInformer.GetStore()
	informers = append(informers, servicesInformer)

	kongIngressInformer := kongInformerFactory.Configuration().V1().KongIngresses().Informer()
	kongIngressInformer.AddEventHandler(reh)
	cacheStores.Configuration = kongIngressInformer.GetStore()
	informers = append(informers, kongIngressInformer)

	kongPluginInformer := kongInformerFactory.Configuration().V1().KongPlugins().Informer()
	kongPluginInformer.AddEventHandler(reh)
	cacheStores.Plugin = kongPluginInformer.GetStore()
	informers = append(informers, kongPluginInformer)

	kongConsumerInformer := kongInformerFactory.Configuration().V1().KongConsumers().Informer()
	kongConsumerInformer.AddEventHandler(reh)
	cacheStores.Consumer = kongConsumerInformer.GetStore()
	informers = append(informers, kongConsumerInformer)

	kongCredentialInformer := kongInformerFactory.Configuration().V1().KongCredentials().Informer()
	kongCredentialInformer.AddEventHandler(reh)
	cacheStores.Credential = kongCredentialInformer.GetStore()
	informers = append(informers, kongCredentialInformer)

	stopCh := make(chan struct{})
	for _, informer := range informers {
		go informer.Run(stopCh)
		synced = append(synced, informer.HasSynced)
	}
	if !cache.WaitForCacheSync(stopCh, synced...) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	}

	store := store.New(
		cacheStores,
		annotations.IngressClassValidatorFuncFromObjectMeta(conf.IngressClass),
	)
	kong, err := controller.NewKongController(conf, updateChannel, store)
	if err != nil {
		glog.Fatal(err)
	}

	go handleSigterm(kong, stopCh, func(code int) {
		os.Exit(code)
	})

	mux := http.NewServeMux()
	go registerHandlers(conf.EnableProfiling, 10254, kong, mux)

	if "off" != os.Getenv("KONG_ANONYMOUS_REPORTS") {
		hostname, err := os.Hostname()
		if err != nil {
			glog.Error(err)
		}
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			glog.Error(err)
		}
		k8sVersion, err := kubeClient.Discovery().ServerVersion()
		if err != nil {
			glog.Error(err)
		}
		info := utils.Info{
			KongVersion:       root["version"].(string),
			KICVersion:        RELEASE,
			KubernetesVersion: fmt.Sprintf("%s", k8sVersion),
			Hostname:          hostname,
			ID:                uuid,
			KongDB:            kongDB,
		}
		reporter := utils.NewReporter(info)
		go reporter.Run(stopCh)
	}
	if admissionWebhookListen != "off" {
		admissionServer := admission.Server{
			Validator: admission.KongHTTPValidator{
				Client: kongClient,
			},
		}
		go func() {
			glog.Error("error running the admission controller server:",
				http.ListenAndServeTLS(
					admissionWebhookListen,
					admissionWebhookCertPath,
					admissionWebhookKeyPath,
					admissionServer,
				))
		}()
	}
	kong.Start()
}

type exiter func(code int)

func handleSigterm(kong *controller.KongController, stopCh chan struct{},
	exit exiter) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	glog.Infof("Received SIGTERM, shutting down")

	exitCode := 0
	close(stopCh)
	if err := kong.Stop(); err != nil {
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

	var v *k8sVersion.Info

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

func registerHandlers(enableProfiling bool, port int, ic *controller.KongController, mux *http.ServeMux) {

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(version())
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
