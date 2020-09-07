/*
Copyright 2017 The Kubernetes Authors.

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
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	apiv1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
)

const (
	defaultKongAdminURL             = "http://localhost:8001"
	defaultKongFilterTag            = "managed-by-ingress-controller"
	defaultAdmissionWebhookCertPath = "/admission-webhook/tls.crt"
	defaultAdmissionWebhookKeyPath  = "/admission-webhook/tls.key"
)

type cliConfig struct {
	// Admission controller server properties
	AdmissionWebhookListen   string
	AdmissionWebhookCertPath string
	AdmissionWebhookKeyPath  string
	AdmissionWebhookCert     string
	AdmissionWebhookKey      string

	// Kong connection details
	KongAdminURL             string
	KongWorkspace            string
	KongAdminConcurrency     int
	KongAdminFilterTags      []string
	KongAdminHeaders         []string
	KongAdminTLSSkipVerify   bool
	KongAdminTLSServerName   string
	KongAdminCACertPath      string
	KongAdminCACert          string
	KongCustomEntitiesSecret string

	// Resource filtering
	WatchNamespace                 string
	ProcessClasslessIngressV1beta1 bool
	ProcessClasslessKongConsumer   bool
	IngressClass                   string
	ElectionID                     string

	// Ingress Status publish resource
	PublishService         string
	PublishStatusAddress   string
	UpdateStatus           bool
	UpdateStatusOnShutdown bool

	// Runtime behavior
	SyncPeriod        time.Duration
	SyncRateLimit     float32
	EnableReverseSync bool

	// Logging
	LogLevel  string
	LogFormat string

	// k8s connection details
	APIServerHost      string
	KubeConfigFilePath string

	// Allowed Ingress resource versions
	AllowIngressExtensionsV1beta1 bool
	AllowIngressNetworkingV1beta1 bool
	AllowIngressNetworkingV1      bool

	// Performance
	EnableProfiling bool

	// Misc
	ShowVersion      bool
	AnonymousReports bool
}

func flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ExitOnError)

	// Admission controller server properties
	flags.String("admission-webhook-listen", "off",
		`The address to start admission controller on (ip:port).
Setting it to 'off' disables the admission controller.`)
	flags.String("admission-webhook-cert-file", defaultAdmissionWebhookCertPath,
		`Path to the PEM-encoded certificate file for
TLS handshake`)
	flags.String("admission-webhook-key-file", defaultAdmissionWebhookKeyPath,
		`Path to the PEM-encoded private key file for
TLS handshake`)
	flags.String("admission-webhook-cert", "",
		`PEM-encoded certificate for TLS handshake`)
	flags.String("admission-webhook-key", "",
		`PEM-encoded private key for TLS handshake`)

	// Kong connection details
	// deprecated
	flags.String("kong-url", "",
		`DEPRECATED, use --kong-admin-url
The address of the Kong Admin URL to connect to in the
format of protocol://address:port`)
	// new
	flags.String("kong-admin-url", defaultKongAdminURL,
		`The address of the Kong Admin URL to connect to in the
format of protocol://address:port`)

	flags.String("kong-workspace", "",
		"Workspace in Kong Enterprise to be configured")

	flag.Int("kong-admin-concurrency", 10,
		"Max number of concurrent requests sent to Kong's Admin API")

	flags.StringSlice("kong-admin-filter-tag", []string{defaultKongFilterTag},
		`The tag used to manage and filter entities in Kong
This flag can be specified multiple times to specify multiple tags.`)

	// deprecated
	flags.StringSlice("admin-header", nil,
		`DEPRECATED, use --kong-admin-header
add a header (key:value) to every Admin API call,
this flag can be used multiple times to specify multiple headers`)
	// new
	flags.StringSlice("kong-admin-header", nil,
		`add a header (key:value) to every Admin API call,
this flag can be used multiple times to specify multiple headers`)

	flags.String("kong-admin-token", "",
		`Sets the value of the 'kong-admin-token' header; useful for
authentication/authorization for Kong Enterprise environments`)

	// deprecated
	flags.Bool("admin-tls-skip-verify", false,
		`DEPRECATED, use --kong-admin-tls-skip-verify
Disable verification of TLS certificate of Kong's Admin endpoint.`)
	// new
	flags.Bool("kong-admin-tls-skip-verify", false,
		"Disable verification of TLS certificate of Kong's Admin endpoint.")

	// deprecated
	flags.String("admin-tls-server-name", "",
		`DEPRECATED, use --kong-admin-tls-server-name
SNI name to use to verify the certificate presented by Kong in TLS.`)
	// new
	flags.String("kong-admin-tls-server-name", "",
		"SNI name to use to verify the certificate presented by Kong in TLS.")

	// deprecated
	flags.String("admin-ca-cert-file", "",
		`DEPRECATED, use --kong-admin-ca-cert-file
Path to PEM-encoded CA certificate file to verify
Kong's Admin SSL certificate.`)
	// new
	flags.String("kong-admin-ca-cert-file", "",
		`Path to PEM-encoded CA certificate file to verify
Kong's Admin SSL certificate.`)

	flags.String("kong-admin-ca-cert", "",
		`PEM-encoded CA certificate to verify Kong's Admin SSL certificate.`)

	flags.String("kong-custom-entities-secret", "",
		`Secret containing custom entities that should be populated in DB-less
mode of Kong. Takes the form of namespace/name.`)

	// Resource filtering
	flags.String("watch-namespace", apiv1.NamespaceAll,
		`Namespace to watch for Ingress. Default is to watch all namespaces`)
	flags.Bool("skip-classless-ingress-v1beta1", false,
		`Skip non annotated Ingresses and Kong CRDs.`)
	flags.String("ingress-class", annotations.DefaultIngressClass,
		`Name of the ingress class to route through this controller.`)
	flags.String("election-id", "ingress-controller-leader",
		`Election id to use for status update.`)

	// Ingress Status publish resource
	flags.String("publish-service", "",
		`Service fronting the ingress controllers. Takes the form namespace/name.
The controller will set the endpoint records on the ingress objects
to reflect those on the service.`)
	flags.String("publish-status-address", "",
		`User customized address to be set in the status of ingress resources.
The controller will set the endpoint records on the
ingress using this address.`)
	flags.Bool("update-status", true, `Indicates if the ingress controller
should update the Ingress status IP/hostname.`)
	flags.Bool("update-status-on-shutdown", true,
		`Indicates if the ingress controller should update the Ingress status 
IP/hostname when the controller is being stopped.`)

	// Runtime behavior
	flags.Duration("sync-period", 600*time.Second,
		`Relist and confirm cloud resources this often.`)
	flags.Float32("sync-rate-limit", 0.3,
		`Define the sync frequency upper limit`)
	flag.Bool("enable-reverse-sync", false, `Enable reverse checks from Kong to Kubernetes`)

	// Logging
	flags.String("log-level", "info",
		`Level of logging for the controller. Allowed values are 
trace, debug, info, warn, error, fatal and panic.`)
	flags.String("log-format", "text",
		`Format of logs of the controller. Allowed values are 
text and json.`)

	// k8s connection details
	flags.String("apiserver-host", "",
		`The address of the Kubernetes Apiserver to connect to in the format of 
protocol://address:port, e.g., "http://localhost:8080.
If not specified, the assumption is that the binary runs inside a 
Kubernetes cluster and local discovery is attempted.`)
	flags.String("kubeconfig", "", "Path to kubeconfig file with "+
		"authorization and master location information.")

	// Allowed Ingress resource versions
	flags.Bool("allow-ingress-extensionsv1beta1", true,
		`If disabled, the ingress controller won't try extensions/v1beta1 when negotiating the newest supported
Ingress API with Kubernetes.`)
	flags.Bool("allow-ingress-networkingv1beta1", true,
		`If disabled, the ingress controller won't try networking.k8s.io/v1beta1 when negotiating the newest supported
Ingress API with Kubernetes.`)
	flags.Bool("allow-ingress-networkingv1", false,
		`If disabled, the ingress controller won't try networking/v1 when negotiating the newest supported
Ingress API with Kubernetes.`)

	// Misc
	flags.Bool("profiling", true, `Enable profiling via web interface host:port/debug/pprof/`)
	flags.Bool("version", false,
		`Shows release information about the Kong Ingress controller`)
	flags.Bool("anonymous-reports", true,
		`Send anonymized usage data to help improve Kong`)

	return flags
}

func parseFlags() (cliConfig, error) {

	flagSet := flagSet()

	flagSet.AddGoFlagSet(flag.CommandLine)
	if err := flagSet.Parse(os.Args); err != nil {
		return cliConfig{}, err
	}

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		return cliConfig{}, err
	}

	viper.SetEnvPrefix("CONTROLLER")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if err := viper.BindPFlags(flagSet); err != nil {
		return cliConfig{}, err
	}

	var config cliConfig
	// Admission controller server properties
	config.AdmissionWebhookListen = viper.GetString("admission-webhook-listen")
	config.AdmissionWebhookCertPath =
		viper.GetString("admission-webhook-cert-file")
	config.AdmissionWebhookKeyPath =
		viper.GetString("admission-webhook-key-file")
	config.AdmissionWebhookCert =
		viper.GetString("admission-webhook-cert")
	config.AdmissionWebhookKey =
		viper.GetString("admission-webhook-key")

	// Kong connection details
	kongAdminURL := defaultKongAdminURL
	oldURL := viper.GetString("kong-url")
	newURL := viper.GetString("kong-admin-url")
	if oldURL != "" {
		kongAdminURL = oldURL
	}
	if newURL != defaultKongAdminURL {
		kongAdminURL = newURL
	}
	config.KongAdminURL = kongAdminURL

	config.KongWorkspace = viper.GetString("kong-workspace")
	config.KongAdminConcurrency = viper.GetInt("kong-admin-concurrency")
	config.KongAdminFilterTags = viper.GetStringSlice("kong-admin-filter-tag")

	config.KongAdminHeaders = viper.GetStringSlice("admin-header")
	kongAdminHeaders := viper.GetStringSlice("kong-admin-header")
	if len(kongAdminHeaders) > 0 {
		config.KongAdminHeaders = kongAdminHeaders
	}

	kongAdminToken := viper.GetString("kong-admin-token")
	if kongAdminToken != "" {
		config.KongAdminHeaders = append(config.KongAdminHeaders,
			"kong-admin-token:"+kongAdminToken)
	}

	config.KongAdminTLSSkipVerify = viper.GetBool("admin-tls-skip-verify")
	kongAdminTLSSkipVerify := viper.GetBool("kong-admin-tls-skip-verify")
	if kongAdminTLSSkipVerify {
		config.KongAdminTLSSkipVerify = kongAdminTLSSkipVerify
	}

	config.KongAdminTLSServerName = viper.GetString("admin-tls-server-name")
	kongAdminTLSServerName := viper.GetString("kong-admin-tls-server-name")
	if kongAdminTLSServerName != "" {
		config.KongAdminTLSServerName = kongAdminTLSServerName
	}

	config.KongAdminCACertPath = viper.GetString("admin-ca-cert-file")
	kongAdminCACertPath := viper.GetString("kong-admin-ca-cert-file")
	if kongAdminCACertPath != "" {
		config.KongAdminCACertPath = kongAdminCACertPath
	}

	kongAdminCACert := viper.GetString("kong-admin-ca-cert")
	if kongAdminCACert != "" {
		config.KongAdminCACert = kongAdminCACert
	}

	config.KongCustomEntitiesSecret = viper.GetString(
		"kong-custom-entities-secret")

	// Resource filtering
	config.WatchNamespace = viper.GetString("watch-namespace")
	config.ProcessClasslessIngressV1beta1 = viper.GetBool("process-classless-ingress-v1beta1")
	config.ProcessClasslessKongConsumer = viper.GetBool("process-classless-kong-consumer")
	config.IngressClass = viper.GetString("ingress-class")
	config.ElectionID = viper.GetString("election-id")

	// Ingress Status publish resource
	config.PublishService = viper.GetString("publish-service")
	config.PublishStatusAddress = viper.GetString("publish-status-address")
	config.UpdateStatus = viper.GetBool("update-status")
	config.UpdateStatusOnShutdown = viper.GetBool("update-status-on-shutdown")

	// Rutnime behavior
	config.SyncPeriod = viper.GetDuration("sync-period")
	config.SyncRateLimit = (float32)(viper.GetFloat64("sync-rate-limit"))
	config.EnableReverseSync = viper.GetBool("enable-reverse-sync")

	// Logging
	config.LogLevel = viper.GetString("log-level")
	config.LogFormat = viper.GetString("log-format")

	// k8s connection details
	config.APIServerHost = viper.GetString("apiserver-host")
	config.KubeConfigFilePath = viper.GetString("kubeconfig")

	// Allowed Ingress resource versions
	config.AllowIngressExtensionsV1beta1 = viper.GetBool("allow-ingress-extensionsv1beta1")
	config.AllowIngressNetworkingV1beta1 = viper.GetBool("allow-ingress-networkingv1beta1")
	config.AllowIngressNetworkingV1 = viper.GetBool("allow-ingress-networkingv1")

	// Misc
	config.EnableProfiling = viper.GetBool("profiling")
	config.ShowVersion = viper.GetBool("version")
	config.AnonymousReports = viper.GetBool("anonymous-reports")
	return config, nil
}
