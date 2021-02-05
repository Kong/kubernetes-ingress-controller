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
	"github.com/hashicorp/go-uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller"
	configuration "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configclientv1 "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/clientset/versioned"
	configinformer "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/informers/externalversions"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog"
	knativeclient "knative.dev/networking/pkg/client/clientset/versioned"
	knativeinformer "knative.dev/networking/pkg/client/informers/externalversions"
)

var (
	logrusLevel = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}
	logrusFormat = map[string]logrus.Formatter{
		"text": &logrus.TextFormatter{},
		"json": &logrus.JSONFormatter{},
	}
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

		DumpConfig: cliConfig.DumpConfig,
	}
}

func init() {
	// initialize for dependencies
	klog.InitFlags(nil)
}

func main() {
	ctx := context.Background()

	color.Output = ioutil.Discard
	rand.Seed(time.Now().UnixNano())

	fmt.Println(version())

	cliConfig, err := parseFlags()
	if err != nil {
		logrus.Fatalf("failed to parse configuration: %v", err)
	}
	log := logrus.New()
	level, ok := logrusLevel[cliConfig.LogLevel]
	if !ok {
		logrus.Fatalf("invalid log-level: %v", cliConfig.LogLevel)
	}
	log.Level = level

	format, ok := logrusFormat[cliConfig.LogFormat]
	if !ok {
		logrus.Fatalf("invalid log-format: %v", cliConfig.LogFormat)
	}
	log.Formatter = format

	for key, value := range viper.AllSettings() {
		log.WithField(key, fmt.Sprintf("%v", value)).Debug("input flag")
	}

	if cliConfig.ShowVersion {
		os.Exit(0)
	}

	invalidConfErrPrefix := "invalid configuration: "
	if cliConfig.PublishService == "" && cliConfig.PublishStatusAddress == "" {
		log.Fatal(invalidConfErrPrefix + "either --publish-service or --publish-status-address must be specified")
	}

	if cliConfig.SyncPeriod.Seconds() < 10 {
		log.Fatalf(invalidConfErrPrefix+"resync period (%vs) is too low", cliConfig.SyncPeriod.Seconds())
	}

	if cliConfig.KongAdminConcurrency < 1 {
		log.Fatalf(invalidConfErrPrefix+"kong-admin-concurrency (%v) cannot be less than 1", cliConfig.KongAdminConcurrency)
	}

	kubeCfg, kubeClient, err := createApiserverClient(cliConfig.APIServerHost,
		cliConfig.KubeConfigFilePath, log)
	if err != nil {
		log.Fatalf("failed to connect to Kubernetes api-server,"+
			"this most likely means that the cluster is misconfigured (e.g., it has "+
			"invalid apiserver certificates or service accounts configuration); error: %v", err)
	}

	if cliConfig.PublishService != "" {
		svc := cliConfig.PublishService
		ns, name, err := util.ParseNameNS(svc)
		if err != nil {
			log.Fatalf(invalidConfErrPrefix+"publish-service: %v", err)
		}
		_, err = kubeClient.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			log.WithFields(logrus.Fields{
				"service_name":      name,
				"service_namespace": ns,
			}).Fatalf("failed to fetch publish-service: %v", err)
		}
	}

	if cliConfig.WatchNamespace != "" {
		_, err = kubeClient.CoreV1().Namespaces().Get(ctx, cliConfig.WatchNamespace,
			metav1.GetOptions{})
		if err != nil {
			log.Fatalf("failed to fetch watch-namespace '%s': %v", cliConfig.WatchNamespace, err)
		}
	}

	controllerConfig := controllerConfigFromCLIConfig(cliConfig)
	controllerConfig.Logger = log.WithField("component", "controller")

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
		log.Fatalf(invalidConfErrPrefix + "both --kong-admin-ca-cert-path and --kong-admin-ca-cert" +
			"are set; please remove one or the other")
	}
	if cliConfig.KongAdminCACert != "" {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(cliConfig.KongAdminCACert))
		if !ok {
			// TODO give user an error to make this actionable
			log.Fatalf("failed to load kong-admin-ca-cert")
		}
		tlsConfig.RootCAs = certPool
	}
	if cliConfig.KongAdminCACertPath != "" {
		certPath := cliConfig.KongAdminCACertPath
		certPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile(certPath)
		if err != nil {
			log.Fatalf("failed to read kong-admin-ca-cert from path '%s': %v", certPath, err)
		}
		ok := certPool.AppendCertsFromPEM(cert)
		if !ok {
			// TODO give user an error to make this actionable
			log.Fatalf("failed to load kong-admin-ca-cert from path '%s'", certPath)
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
		log.Fatalf("failed to create kong client: %v", err)
	}

	var root map[string]interface{}
	backoff := flowcontrol.NewBackOff(1*time.Second, 15*time.Second)
	const backoffID = "kong-admin-api"
	retryCount := 0
	for {
		root, err = rootWithTimeout(ctx, kongClient)
		if err == nil {
			break
		}
		if retryCount > 5 {
			log.Fatalf("failed to fetch metadata from kong: %v", err)
		}
		backoff.Next(backoffID, backoff.Clock.Now())
		delay := backoff.Get(backoffID)
		time.Sleep(delay)
		retryCount++
		log.Infof("retry %d to fetch metadata from kong: %v", retryCount, err)
		continue
	}
	v, err := getSemVerVer(root["version"].(string))
	if err != nil {
		log.Fatalf("failed to determine version of kong: %v", err)
	}
	log.WithField("kong_version", v).Infof("kong version: %s", v)
	controllerConfig.Kong.Version = v

	if strings.Contains(root["version"].(string), "enterprise") {
		log.Debug("enterprise version of kong detected")
		controllerConfig.Kong.Enterprise = true
	}

	kongConfiguration := root["configuration"].(map[string]interface{})
	kongDB := kongConfiguration["database"].(string)
	log.Infof("datastore strategy for kong: %s", kongDB)

	if kongDB == "off" {
		controllerConfig.Kong.InMemory = true
	}
	if kongDB == "cassandra" {
		log.Fatalf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported;" +
			"you must migrate to a Postgres-backed or DB-less deployment")
	}

	req, _ := http.NewRequest("GET",
		cliConfig.KongAdminURL+"/tags", nil)
	res, err := kongClient.Do(ctx, req, nil)
	if err == nil && res.StatusCode == 200 {
		controllerConfig.Kong.HasTagSupport = true
	}

	// setup workspace in Kong Enterprise
	if cliConfig.KongWorkspace != "" {
		// ensure the workspace exists or try creating it
		err := ensureWorkspace(ctx, kongClient, cliConfig.KongWorkspace)
		if err != nil {
			log.Fatalf("failed to ensure workspace in kong: %v", err)
		}
		kongClient, err = kong.NewClient(kong.String(cliConfig.KongAdminURL+"/"+cliConfig.KongWorkspace), c)
		if err != nil {
			log.Fatalf("failed to create kong client: %v", err)
		}
	}
	controllerConfig.Kong.Client = kongClient

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

	if cliConfig.DumpConfig != util.ConfigDumpModeOff {
		controllerConfig.DumpDir, err = ioutil.TempDir("", "controller")
		if err != nil {
			log.Fatalf("failed to create a dump directory: %v", err)
		}
	}

	var synced []cache.InformerSynced
	updateChannel := channels.NewRingChannel(1024)
	reh := controller.ResourceEventHandler{
		UpdateCh: updateChannel,
	}

	var preferredIngressAPIs []util.IngressAPI
	if !cliConfig.DisableIngressNetworkingV1 {
		preferredIngressAPIs = append(preferredIngressAPIs, util.NetworkingV1)
	}
	if !cliConfig.DisableIngressNetworkingV1beta1 {
		preferredIngressAPIs = append(preferredIngressAPIs, util.NetworkingV1beta1)
	}
	if !cliConfig.DisableIngressExtensionsV1beta1 {
		preferredIngressAPIs = append(preferredIngressAPIs, util.ExtensionsV1beta1)
	}
	controllerConfig.IngressAPI, err = util.NegotiateResourceAPI(kubeClient, "Ingress", preferredIngressAPIs)
	if err != nil {
		log.Fatalf("NegotiateIngressAPI failed: %v, tried: %+v", err, preferredIngressAPIs)
	}
	log.Infof("chosen Ingress API version: %v", controllerConfig.IngressAPI)

	var informers []cache.SharedIndexInformer
	var cacheStores store.CacheStores

	var ingInformer cache.SharedIndexInformer
	switch controllerConfig.IngressAPI {
	case util.NetworkingV1:
		ingInformer = coreInformerFactory.Networking().V1().Ingresses().Informer()
		cacheStores.IngressV1 = ingInformer.GetStore()
		cacheStores.IngressV1beta1 = newEmptyStore()
	case util.NetworkingV1beta1:
		ingInformer = coreInformerFactory.Networking().V1beta1().Ingresses().Informer()
		cacheStores.IngressV1 = newEmptyStore()
		cacheStores.IngressV1beta1 = ingInformer.GetStore()
	case util.ExtensionsV1beta1:
		ingInformer = coreInformerFactory.Extensions().V1beta1().Ingresses().Informer()
		cacheStores.IngressV1 = newEmptyStore()
		cacheStores.IngressV1beta1 = ingInformer.GetStore()
	}
	ingInformer.AddEventHandler(reh)
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

	hasKongClusterPlugin, err := util.ServerHasGVK(kubeClient.Discovery(),
		configuration.SchemeGroupVersion.String(), "KongClusterPlugin")

	if hasKongClusterPlugin {
		kongClusterPluginInformer := kongInformerFactory.Configuration().V1().KongClusterPlugins().Informer()
		kongClusterPluginInformer.AddEventHandler(reh)
		cacheStores.ClusterPlugin = kongClusterPluginInformer.GetStore()
		informers = append(informers, kongClusterPluginInformer)
	} else {
		if err != nil {
			log.Fatalf("failed to retrieve KongClusterPlugin availability: %s", err)
		}
		log.Warn("KongClusterPlugin CRD not detected. Disabling KongClusterPlugin functionality.")
		cacheStores.ClusterPlugin = newEmptyStore()
	}

	kongConsumerInformer := kongInformerFactory.Configuration().V1().KongConsumers().Informer()
	kongConsumerInformer.AddEventHandler(reh)
	cacheStores.Consumer = kongConsumerInformer.GetStore()
	informers = append(informers, kongConsumerInformer)

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

	store := store.New(cacheStores, cliConfig.IngressClass, cliConfig.ProcessClasslessIngressV1Beta1,
		cliConfig.ProcessClasslessIngressV1, cliConfig.ProcessClasslessKongConsumer, log.WithField("component", "store"))

	kong, err := controller.NewKongController(ctx, &controllerConfig, updateChannel,
		store)
	if err != nil {
		log.Fatalf("failed to create a controller: %v", err)
	}

	exitCh := make(chan int, 1)
	var wg sync.WaitGroup
	mux := http.NewServeMux()
	wg.Add(1)
	go func() {
		defer wg.Done()
		serveHTTP(cliConfig.EnableProfiling,
			10254, mux, stopCh,
			log.WithField("component", "metadata-server"))
	}()
	go handleSigterm(kong, stopCh, exitCh, log.WithField("component", "signal-handler"))

	if cliConfig.AnonymousReports {
		logger := log.WithField("component", "reporter")
		hostname, err := os.Hostname()
		if err != nil {
			logger.Warnf("failed to fetch hostname: %v", err)
		}
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			logger.Warnf("failed to generate a random uuid: %v", err)
		}
		k8sVersion, err := kubeClient.Discovery().ServerVersion()
		if err != nil {
			logger.Warnf("failed to fetch k8s api-server version: %v", err)
		}
		info := util.Info{
			KongVersion:       root["version"].(string),
			KICVersion:        RELEASE,
			KubernetesVersion: k8sVersion.String(),
			Hostname:          hostname,
			ID:                uuid,
			KongDB:            kongDB,
		}
		reporter := util.Reporter{
			Info:   info,
			Logger: logger,
		}
		reporter.Logger = logger
		go reporter.Run(stopCh)
	}
	if cliConfig.AdmissionWebhookListen != "off" {
		logger := log.WithField("component", "admission-server")
		admissionServer := admission.Server{
			Validator: admission.KongHTTPValidator{
				Client: kongClient,
				Logger: logger,
			},
			Logger: logger,
		}
		var cert tls.Certificate
		if cliConfig.AdmissionWebhookCertPath != defaultAdmissionWebhookCertPath && cliConfig.AdmissionWebhookCert != "" {
			logger.Fatalf(invalidConfErrPrefix + "both --admission-webhook-cert-file and --admission-webhook-cert" +
				"are set; please remove one or the other")
		}
		if cliConfig.AdmissionWebhookKeyPath != defaultAdmissionWebhookKeyPath && cliConfig.AdmissionWebhookKey != "" {
			logger.Fatalf(invalidConfErrPrefix + "both --admission-webhook-cert-key and --admission-webhook-key" +
				"are set; please remove one or the other")
		}
		if cliConfig.AdmissionWebhookCert != "" {
			var err error
			cert, err = tls.X509KeyPair([]byte(cliConfig.AdmissionWebhookCert), []byte(cliConfig.AdmissionWebhookKey))
			if err != nil {
				logger.Fatalf("failed to load admission webhook certificate: %s", err)
			}
		}
		// although this is partially checked earlier, that check does not fail if it sees the default path
		// we don't want to overwrite any certs set by admission-webhook-cert, but also don't want to run this
		// first, as it can potentially result in a fatal error
		if cliConfig.AdmissionWebhookCertPath != "" && cliConfig.AdmissionWebhookCert == "" {
			var err error
			cert, err = tls.LoadX509KeyPair(cliConfig.AdmissionWebhookCertPath, cliConfig.AdmissionWebhookKeyPath)
			if err != nil {
				logger.Fatalf("failed to load admission webhook certificate: %s", err)
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
			err := server.ListenAndServeTLS("", "")
			logger.Errorf("server stopped with err: %v", err)
		}()
	}
	kong.Start()
	wg.Wait()
	os.Exit(<-exitCh)
}

func handleSigterm(kong *controller.KongController,
	stopCh chan<- struct{},
	exitCh chan<- int,
	logger logrus.FieldLogger) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	logger.Infof("Received SIGTERM, shutting down")

	exitCode := 0
	close(stopCh)
	if err := kong.Stop(); err != nil {
		logger.Errorf("failed to stop controller: %v", err)
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
func createApiserverClient(apiserverHost string, kubeConfig string,
	logger logrus.FieldLogger) (*rest.Config, *kubernetes.Clientset, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(apiserverHost, kubeConfig)
	if err != nil {
		return nil, nil, err
	}

	cfg.QPS = defaultQPS
	cfg.Burst = defaultBurst

	// cfg.ContentType = "application/vnd.kubernetes.protobuf"

	logger = logger.WithField("api-server-host", cfg.Host)
	logger.Debugf("creating api client")

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
	retryCount := 0
	logger.Debug("attempting to discover version of kubernetes")
	err = wait.ExponentialBackoff(defaultRetry, func() (bool, error) {
		v, err = client.Discovery().ServerVersion()

		if err == nil {
			return true, nil
		}

		lastErr = err
		logger.WithField("retry_count", retryCount).Warnf("failed to fetch version of kubernetes api-server: %v", err)
		retryCount++
		return false, nil
	})

	// err is not null only if there was a timeout in the exponential backoff (ErrWaitTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch version of kubernetes api-server: %w", lastErr)
	}

	logger.WithFields(logrus.Fields{
		"major":          v.Major,
		"minor":          v.Minor,
		"git_version":    v.GitVersion,
		"git_tree_state": v.GitTreeState,
		"git_commit":     v.GitCommit,
		"platform":       v.Platform,
	}).Infof("version of kubernetes api-server: %v.%v", v.Major, v.Minor)

	return cfg, client, nil
}

func rootWithTimeout(ctx context.Context, kc *kong.Client) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return kc.Root(ctx)
}

const (
	// High enough QPS to fit all expected use cases. QPS=0 is not set here, because
	// client code is overriding it.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases. Burst=0 is not set here, because
	// client code is overriding it.
	defaultBurst = 1e6
)

func serveHTTP(enableProfiling bool,
	port int,
	mux *http.ServeMux,
	stop <-chan struct{},
	logger logrus.FieldLogger) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(version())
		if _, err := w.Write(b); err != nil {
			logger.WithField("endpoint", "/build").Errorf("failed to write response: %v", err)
		}
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			logger.WithField("endpoint", "/stop").Errorf("failed to send SIGTERM to self: %v", err)
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
				logger.Errorf("failed to shut down server: %v", err)
			}
		case <-serveDone:
		}
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("server stopped with err: %v", err)
		close(serveDone)
	}
	wg.Wait()
}
