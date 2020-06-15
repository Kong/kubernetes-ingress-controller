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
	"context"
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
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/eapache/channels"
	"github.com/fatih/color"
	"github.com/golang/glog"
	"github.com/hashicorp/go-uuid"
	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/admission"
	configclientv1 "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/clientset/versioned"
	configinformer "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/informers/externalversions"
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
	knativeclient "knative.dev/serving/pkg/client/clientset/versioned"
	knativeinformer "knative.dev/serving/pkg/client/informers/externalversions"
)

func controllerConfigFromCLIConfig(cliConfig cliConfig) controller.Configuration {
	return controller.Configuration{
		Kong: controller.Kong{
			URL:         cliConfig.KongAdminURL,
			FilterTags:  cliConfig.KongAdminFilterTags,
			Concurrency: cliConfig.KongAdminConcurrency,
		},
		KongCustomEntitiesSecret: cliConfig.KongCustomEntitiesSecret,

		ResyncPeriod:      cliConfig.SyncPeriod,
		SyncRateLimit:     cliConfig.SyncRateLimit,
		EnableReverseSync: cliConfig.EnableReverseSync,

		Namespace: cliConfig.WatchNamespace,

		IngressClass: cliConfig.IngressClass,

		PublishService:       cliConfig.PublishService,
		PublishStatusAddress: cliConfig.PublishStatusAddress,

		UpdateStatus:           cliConfig.UpdateStatus,
		UpdateStatusOnShutdown: cliConfig.UpdateStatusOnShutdown,
		ElectionID:             cliConfig.ElectionID,
	}
}

func main() {
	color.Output = ioutil.Discard
	rand.Seed(time.Now().UnixNano())

	fmt.Println(version())

	cliConfig, err := parseFlags()
	if err != nil {
		glog.Fatal(err)
	}

	if cliConfig.ShowVersion {
		os.Exit(0)
	}

	if cliConfig.PublishService == "" && cliConfig.PublishStatusAddress == "" {
		glog.Fatal("either --publish-service or --publish-status-address",
			"must be specified")
	}

	if cliConfig.SyncPeriod.Seconds() < 10 {
		glog.Fatalf("resync period (%vs) is too low", cliConfig.SyncPeriod.Seconds())
	}

	if cliConfig.KongAdminConcurrency < 1 {
		glog.Fatalf("kong-admin-concurrency (%v) cannot be less than 1",
			cliConfig.KongAdminConcurrency)
	}

	kubeCfg, kubeClient, err := createApiserverClient(cliConfig.APIServerHost,
		cliConfig.KubeConfigFilePath)
	if err != nil {
		handleFatalInitError(err)
	}

	if cliConfig.PublishService != "" {
		svc := cliConfig.PublishService
		ns, name, err := utils.ParseNameNS(svc)
		if err != nil {
			glog.Fatal(err)
		}
		_, err = kubeClient.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
		if err != nil {
			glog.Fatalf("unexpected error getting information about service %v: %v", svc, err)
		}
	}

	if cliConfig.WatchNamespace != "" {
		_, err = kubeClient.CoreV1().Namespaces().Get(cliConfig.WatchNamespace,
			metav1.GetOptions{})
		if err != nil {
			glog.Fatalf("no namespace with name %v found: %v",
				cliConfig.WatchNamespace, err)
		}
	}

	controllerConfig := controllerConfigFromCLIConfig(cliConfig)

	controllerConfig.KubeClient = kubeClient

	defaultTransport := http.DefaultTransport.(*http.Transport)

	var tlsConfig tls.Config

	if cliConfig.KongAdminTLSSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	if cliConfig.KongAdminTLSServerName != "" {
		tlsConfig.ServerName = cliConfig.KongAdminTLSServerName
	}

	if cliConfig.KongAdminCACertPath != "" && cliConfig.KongAdminCACert != "" {
		glog.Fatalf("both --kong-admin-ca-cert-path and --kong-admin-ca-cert" +
			"are set. Please remove one or the other")
	}
	if cliConfig.KongAdminCACert != "" {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(cliConfig.KongAdminCACert))
		if !ok {
			glog.Fatalf("failed to load CACert")
		}
		tlsConfig.RootCAs = certPool
	}
	if cliConfig.KongAdminCACertPath != "" {
		certPath := cliConfig.KongAdminCACertPath
		certPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile(certPath)
		if err != nil {
			glog.Fatalf("failed to read CACert: %s", certPath)
		}
		ok := certPool.AppendCertsFromPEM(cert)
		if !ok {
			glog.Fatalf("failed to load CACert: %s", certPath)
		}
		tlsConfig.RootCAs = certPool
	}
	defaultTransport.TLSClientConfig = &tlsConfig
	c := http.DefaultClient
	c.Transport = &HeaderRoundTripper{
		headers: cliConfig.KongAdminHeaders,
		rt:      defaultTransport,
	}

	kongClient, err := kong.NewClient(kong.String(cliConfig.KongAdminURL), c)
	if err != nil {
		glog.Fatalf("Error creating Kong Rest client: %v", err)
	}

	root, err := kongClient.Root(context.Background())
	if err != nil {
		glog.Fatalf("%v", err)
	}
	v, err := getSemVerVer(root["version"].(string))
	if err != nil {
		glog.Fatalf("Error determining Kong version: %v", err)
	}

	glog.Infof("kong version: %s", v)
	kongConfiguration := root["configuration"].(map[string]interface{})
	controllerConfig.Kong.Version = v

	if strings.Contains(root["version"].(string), "enterprise") {
		controllerConfig.Kong.Enterprise = true
	}

	kongDB := kongConfiguration["database"].(string)
	glog.Infof("Kong datastore: %s", kongDB)

	if kongDB == "off" {
		controllerConfig.Kong.InMemory = true
	}
	if kongDB == "cassandra" {
		glog.Error("running controller with kong backed by cassandra is " +
			"deprecated; please consider using postgres or in-memory mode")
	}

	req, _ := http.NewRequest("GET",
		cliConfig.KongAdminURL+"/tags", nil)
	res, err := kongClient.Do(context.Background(), req, nil)
	if err == nil && res.StatusCode == 200 {
		controllerConfig.Kong.HasTagSupport = true
	}

	// setup workspace in Kong Enterprise
	if cliConfig.KongWorkspace != "" {
		// ensure the workspace exists or try creating it
		err := ensureWorkspace(kongClient, cliConfig.KongWorkspace)
		if err != nil {
			glog.Fatalf("Error ensuring workspace: %v", err)
		}
		kongClient, err = kong.NewClient(kong.String(cliConfig.KongAdminURL+"/"+cliConfig.KongWorkspace), c)
		if err != nil {
			glog.Fatalf("Error creating Kong Rest client: %v", err)
		}
	}
	controllerConfig.Kong.Client = kongClient

	err = discovery.ServerSupportsVersion(kubeClient.Discovery(), schema.GroupVersion{
		Group:   "networking.k8s.io",
		Version: "v1beta1",
	})
	if err == nil {
		controllerConfig.UseNetworkingV1beta1 = true
	}
	coreInformerFactory := informers.NewSharedInformerFactoryWithOptions(
		kubeClient,
		cliConfig.SyncPeriod,
		informers.WithNamespace(cliConfig.WatchNamespace),
	)
	confClient, _ := configclientv1.NewForConfig(kubeCfg)
	controllerConfig.KongConfigClient = confClient

	kongInformerFactory := configinformer.NewSharedInformerFactoryWithOptions(
		confClient,
		cliConfig.SyncPeriod,
		configinformer.WithNamespace(cliConfig.WatchNamespace),
	)

	knativeClient, _ := knativeclient.NewForConfig(kubeCfg)

	var knativeInformerFactory knativeinformer.SharedInformerFactory
	err = discovery.ServerSupportsVersion(knativeClient.Discovery(), schema.GroupVersion{
		Group:   "networking.internal.knative.dev",
		Version: "v1alpha1",
	})
	if err == nil {
		controllerConfig.EnableKnativeIngressSupport = true
		controllerConfig.KnativeClient = knativeClient
		knativeInformerFactory = knativeinformer.NewSharedInformerFactoryWithOptions(
			knativeClient,
			cliConfig.SyncPeriod,
			knativeinformer.WithNamespace(cliConfig.WatchNamespace),
		)
	}

	var synced []cache.InformerSynced
	updateChannel := channels.NewRingChannel(1024)
	reh := controller.ResourceEventHandler{
		UpdateCh:           updateChannel,
		IsValidIngresClass: annotations.IngressClassValidatorFunc(cliConfig.IngressClass),
	}
	var informers []cache.SharedIndexInformer
	var cacheStores store.CacheStores

	var ingInformer cache.SharedIndexInformer
	if controllerConfig.UseNetworkingV1beta1 {
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

	tcpIngressInformer := kongInformerFactory.Configuration().V1beta1().TCPIngresses().Informer()
	tcpIngressInformer.AddEventHandler(reh)
	cacheStores.TCPIngress = tcpIngressInformer.GetStore()
	informers = append(informers, tcpIngressInformer)

	kongIngressInformer := kongInformerFactory.Configuration().V1().KongIngresses().Informer()
	kongIngressInformer.AddEventHandler(reh)
	cacheStores.Configuration = kongIngressInformer.GetStore()
	informers = append(informers, kongIngressInformer)

	kongPluginInformer := kongInformerFactory.Configuration().V1().KongPlugins().Informer()
	kongPluginInformer.AddEventHandler(reh)
	cacheStores.Plugin = kongPluginInformer.GetStore()
	informers = append(informers, kongPluginInformer)

	kongClusterPluginInformer := kongInformerFactory.Configuration().V1().KongClusterPlugins().Informer()
	kongClusterPluginInformer.AddEventHandler(reh)
	cacheStores.ClusterPlugin = kongClusterPluginInformer.GetStore()
	informers = append(informers, kongClusterPluginInformer)

	kongConsumerInformer := kongInformerFactory.Configuration().V1().KongConsumers().Informer()
	kongConsumerInformer.AddEventHandler(reh)
	cacheStores.Consumer = kongConsumerInformer.GetStore()
	informers = append(informers, kongConsumerInformer)

	kongCredentialInformer := kongInformerFactory.Configuration().V1().KongCredentials().Informer()
	kongCredentialInformer.AddEventHandler(reh)
	cacheStores.Credential = kongCredentialInformer.GetStore()
	informers = append(informers, kongCredentialInformer)

	if controllerConfig.EnableKnativeIngressSupport {
		knativeIngressInformer := knativeInformerFactory.Networking().V1alpha1().Ingresses().Informer()
		knativeIngressInformer.AddEventHandler(reh)
		cacheStores.KnativeIngress = knativeIngressInformer.GetStore()
		informers = append(informers, knativeIngressInformer)
	}

	stopCh := make(chan struct{})
	for _, informer := range informers {
		go informer.Run(stopCh)
		synced = append(synced, informer.HasSynced)
	}
	if !cache.WaitForCacheSync(stopCh, synced...) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	}

	store := store.New(cacheStores, cliConfig.IngressClass)
	kong, err := controller.NewKongController(&controllerConfig, updateChannel,
		store)
	if err != nil {
		glog.Fatal(err)
	}

	exitCh := make(chan int, 1)
	var wg sync.WaitGroup
	mux := http.NewServeMux()
	wg.Add(1)
	go func() {
		defer wg.Done()
		serveHTTP(cliConfig.EnableProfiling, 10254, mux, stopCh)
	}()
	go handleSigterm(kong, stopCh, exitCh)

	if cliConfig.AnonymousReports {
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
			KubernetesVersion: k8sVersion.String(),
			Hostname:          hostname,
			ID:                uuid,
			KongDB:            kongDB,
		}
		reporter := utils.NewReporter(info)
		go reporter.Run(stopCh)
	}
	if cliConfig.AdmissionWebhookListen != "off" {
		admissionServer := admission.Server{
			Validator: admission.KongHTTPValidator{
				Client: kongClient,
			},
		}
		var cert tls.Certificate
		if cliConfig.AdmissionWebhookCertPath != defaultAdmissionWebhookCertPath && cliConfig.AdmissionWebhookCert != "" {
			glog.Fatalf("both --admission-webhook-cert-file and --admission-webhook-cert" +
				"are set. Please remove one or the other")
		}
		if cliConfig.AdmissionWebhookKeyPath != defaultAdmissionWebhookKeyPath && cliConfig.AdmissionWebhookKey != "" {
			glog.Fatalf("both --admission-webhook-cert-key and --admission-webhook-key" +
				"are set. Please remove one or the other")
		}
		if cliConfig.AdmissionWebhookCert != "" {
			var err error
			cert, err = tls.X509KeyPair([]byte(cliConfig.AdmissionWebhookCert), []byte(cliConfig.AdmissionWebhookKey))
			if err != nil {
				glog.Fatalf("failed to load admission webhook cert: %s", err)
			}
		}
		// although this is partially checked earlier, that check does not fail if it sees the default path
		// we don't want to overwrite any certs set by admission-webhook-cert, but also don't want to run this
		// first, as it can potentially result in a fatal error
		if cliConfig.AdmissionWebhookCertPath != "" && cliConfig.AdmissionWebhookCert == "" {
			var err error
			cert, err = tls.LoadX509KeyPair(cliConfig.AdmissionWebhookCertPath, cliConfig.AdmissionWebhookKeyPath)
			if err != nil {
				glog.Fatalf("failed to load admission webhook cert: %s", err)
			}
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		server := http.Server{
			Addr:      cliConfig.AdmissionWebhookListen,
			TLSConfig: tlsConfig,
			Handler:   admissionServer,
		}
		go func() {
			glog.Error("error running the admission controller server:",
				server.ListenAndServeTLS("", ""))
		}()
	}
	kong.Start()
	wg.Wait()
	os.Exit(<-exitCh)
}

func handleSigterm(kong *controller.KongController, stopCh chan<- struct{}, exitCh chan<- int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	glog.Infof("Received SIGTERM, shutting down")

	exitCode := 0
	close(stopCh)
	if err := kong.Stop(); err != nil {
		glog.Warningf("Stopping Kong controller failed: %v", err)
		exitCode = 1
	}
	exitCh <- exitCode
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

func serveHTTP(enableProfiling bool, port int, mux *http.ServeMux, stop <-chan struct{}) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(version())
		if _, err := w.Write(b); err != nil {
			glog.Errorf("profiling server /build: %v", err)
		}
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
	serveDone := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-stop:
			// Allow the server to drain for as long as it takes.
			if err := server.Shutdown(context.Background()); err != nil {
				// We know the error wasn't due to a timeout.
				glog.Warningf("Shutting down HTTP server failed: %v", err)
			}
		case <-serveDone:
		}
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		glog.Fatalf("HTTP server failed: %v", err)
		close(serveDone)
	}
	wg.Wait()
}
