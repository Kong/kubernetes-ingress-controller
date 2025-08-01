package config

import (
	"fmt"
	"time"

	"github.com/samber/mo"
	"github.com/spf13/pflag"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/flags"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

type OptionalNamespacedName = mo.Option[k8stypes.NamespacedName]

// Type override to be used with OptionalNamespacedName variables to override their type name printed in the help text.
var nnTypeNameOverride = flags.WithTypeNameOverride[OptionalNamespacedName]("namespaced-name")

// CLIConfig collects all configuration that the controller manager takes from the environment.
type CLIConfig struct {
	// Embed the public managercfg.Config to expose the configuration fields.
	*managercfg.Config

	// See flag definitions in FlagSet(...) for documentation of the fields defined here.
	flagSet *pflag.FlagSet
}

func NewCLIConfig() *CLIConfig {
	cfg := &CLIConfig{
		Config: &managercfg.Config{},
	}
	cfg.bindFlagSet()
	return cfg
}

func (c *CLIConfig) FlagSet() *pflag.FlagSet {
	return c.flagSet
}

// bindFlagSet binds the provided CLIConfig to command-line flags.
func (c *CLIConfig) bindFlagSet() {
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)

	// Logging configurations.
	flagSet.StringVar(&c.LogLevel, "log-level", "info", `Level of logging for the controller. Allowed values are trace, debug, info, and error.`)
	flagSet.StringVar(&c.LogFormat, "log-format", "text", `Format of logs of the controller. Allowed values are text and json.`)

	// Kong high-level controller manager configurations.
	flagSet.BoolVar(&c.KongAdminAPIConfig.TLSSkipVerify, "kong-admin-tls-skip-verify", false, "Disable verification of TLS certificate of Kong's Admin endpoint.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSServerName, "kong-admin-tls-server-name", "", "SNI name to use to verify the certificate presented by Kong in TLS.")
	flagSet.StringVar(&c.KongAdminAPIConfig.CACertPath, "kong-admin-ca-cert-file", "", `Path to PEM-encoded CA certificate file to verify Kong's Admin TLS certificate. Mutually exclusive with --kong-admin-ca-cert.`)
	flagSet.StringVar(&c.KongAdminAPIConfig.CACert, "kong-admin-ca-cert", "", `PEM-encoded CA certificate to verify Kong's Admin TLS certificate. Mutually exclusive with --kong-admin-ca-cert-file.`)

	flagSet.StringSliceVar(&c.KongAdminAPIConfig.Headers, "kong-admin-header", nil, `Header(s) (key:value) in comma-separated format (or specify this flag multiple times) to add to every Admin API call.`)
	flagSet.UintVar(&c.KongAdminInitializationRetries, "kong-admin-init-retries", 60, "Number of attempts that will be made initially on controller startup to connect to the Kong Admin API.")
	flagSet.DurationVar(&c.KongAdminInitializationRetryDelay, "kong-admin-init-retry-delay", time.Second, "The time delay between every attempt (on controller startup) to connect to the Kong Admin API.")
	flagSet.StringVar(&c.KongAdminToken, "kong-admin-token", "", `The Kong Enterprise RBAC token used by the controller. Mutually exclusive with --kong-admin-token-file.`)
	flagSet.StringVar(&c.KongAdminTokenPath, "kong-admin-token-file", "", `Path to the Kong Enterprise RBAC token file used by the controller. Mutually exclusive with --kong-admin-token.`)
	flagSet.StringVar(&c.KongWorkspace, "kong-workspace", "", "Kong Enterprise workspace to configure. Leave this empty if not using Kong workspaces.")
	flagSet.BoolVar(&c.AnonymousReports, "anonymous-reports", true, `Send anonymized usage data to help improve Kong.`)
	flagSet.BoolVar(&c.EnableReverseSync, "enable-reverse-sync", false, `Send configuration to Kong even if the configuration checksum has not changed since previous update.`)
	// TODO: When FallbackConfiguration graduates we should remove the feature gate mention from the help text.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/6170
	flagSet.BoolVar(&c.UseLastValidConfigForFallback, "use-last-valid-config-for-fallback", false, fmt.Sprintf(`When recovering from config push failures, use the last valid configuration cache to backfill broken objects. It can only be used with the %s feature gate enabled.`, managercfg.FallbackConfigurationFeature))
	// Default has to be explicitly passed to generate the proper docs. See https://github.com/kubernetes-sigs/controller-runtime/blob/f1c5dd3851ce3df8b4b7830d9b6eae6271f6932d/pkg/cache/cache.go#L146-L151.
	flagSet.DurationVar(&c.SyncPeriod, "sync-period", 10*time.Hour, `Determine the minimum frequency at which watched resources are reconciled. Set to 0 to use default from controller-runtime.`)
	flagSet.BoolVar(&c.SkipCACertificates, "skip-ca-certificates", false, `Disable syncing CA certificate syncing (for use with multi-workspace environments).`)
	// Default has to be explicitly passed to generate the proper docs. See https://github.com/kubernetes-sigs/controller-runtime/blob/f1c5dd3851ce3df8b4b7830d9b6eae6271f6932d/pkg/config/controller.go#L38-L39.
	flagSet.DurationVar(&c.CacheSyncTimeout, "cache-sync-timeout", 2*time.Minute, `The time limit set to wait for syncing controllers' caches. Set to 0 to use default from controller-runtime.`)

	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClient.CertFile, "kong-admin-tls-client-cert-file", "", "Mutual TLS (mTLS) client certificate file for authentication. Mutually exclusive with --kong-admin-tls-client-cert.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClient.KeyFile, "kong-admin-tls-client-key-file", "", "Mutual TLS (mTLS) client key file for authentication. Mutually exclusive with --kong-admin-tls-client-key.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClient.Cert, "kong-admin-tls-client-cert", "", "Mutual TLS (mTLS) client certificate for authentication. Mutually exclusive with --kong-admin-tls-client-cert-file.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClient.Key, "kong-admin-tls-client-key", "", "Mutual TLS (mTLS) client key for authentication. Mutually exclusive with --kong-admin-tls-client-key-file.")

	// Kong Admin API configuration.
	flagSet.StringSliceVar(&c.KongAdminURLs, "kong-admin-url", []string{"http://localhost:8001"},
		`Kong Admin URL(s) in comma-separated format (or specify this flag multiple times) to connect to in the format "protocol://address:port".`)
	flagSet.Var(flags.NewValidatedValue(&c.KongAdminSvc, namespacedNameFromFlagValue, nnTypeNameOverride), "kong-admin-svc",
		`Kong Admin API Service namespaced name in "namespace/name" format, to use for Kong Gateway service discovery.`)
	flagSet.StringSliceVar(&c.KongAdminSvcPortNames, "kong-admin-svc-port-names", []string{"admin-tls", "kong-admin-tls"},
		"Name(s) of ports on Kong Admin API service in comma-separated format (or specify this flag multiple times) to take into account when doing gateway discovery.")
	flagSet.DurationVar(&c.GatewayDiscoveryReadinessCheckInterval, "gateway-discovery-readiness-check-interval", managercfg.DefaultDataPlanesReadinessReconciliationInterval,
		"Interval of readiness checks on gateway admin API clients for discovery.")
	flagSet.DurationVar(&c.GatewayDiscoveryReadinessCheckTimeout, "gateway-discovery-readiness-check-timeout", managercfg.DefaultDataPlanesReadinessCheckTimeout,
		"Timeout of readiness checks on gateway admin clients.")

	// Kong Proxy and Proxy Cache configurations
	flagSet.StringVar(&c.APIServerHost, "apiserver-host", "", `The Kubernetes API server URL. If not set, the controller will use cluster config discovery.`)
	flagSet.IntVar(&c.APIServerQPS, "apiserver-qps", 100, "The Kubernetes API RateLimiter maximum queries per second.")
	flagSet.IntVar(&c.APIServerBurst, "apiserver-burst", 300, "The Kubernetes API RateLimiter maximum burst queries per second.")
	flagSet.StringVar(&c.MetricsAddr, "metrics-bind-address", fmt.Sprintf(":%v", consts.MetricsPort), "The address the metric endpoint binds to.")
	flagSet.Var(flags.NewValidatedValue(&c.MetricsAccessFilter, metricsAccessFilterFromFlagValue, flags.WithDefault(managercfg.MetricsAccessFilterOff)), "metrics-access-filter", "Specifies the filter access function to be used for accessing the metrics endpoint (possible values: off, rbac).")
	flagSet.StringVar(&c.ProbeAddr, "health-probe-bind-address", fmt.Sprintf(":%v", consts.HealthzPort), "The address the probe endpoint binds to.")
	flagSet.Float32Var(&c.ProxySyncSeconds, "proxy-sync-seconds", dataplane.DefaultSyncSeconds,
		"Define the rate (in seconds) in which configuration updates will be applied to the Kong Admin API.")
	flagSet.DurationVar(&c.InitCacheSyncDuration, "init-cache-sync-duration", dataplane.DefaultCacheSyncWaitDuration, `The initial delay to wait for Kubernetes object caches to be synced before the initial configuration.`)
	flagSet.Float32Var(&c.ProxyTimeoutSeconds, "proxy-timeout-seconds", dataplane.DefaultTimeoutSeconds,
		"Sets the timeout (in seconds) for all requests to Kong's Admin API.")

	// Kubernetes configurations
	flagSet.Var(flags.NewValidatedValue(&c.GatewayAPIControllerName, gatewayAPIControllerNameFromFlagValue, flags.WithDefault(string(gateway.GetControllerName()))), "gateway-api-controller-name", "The controller name to match on Gateway API resources.")
	flagSet.StringVar(&c.KubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file.")
	flagSet.StringVar(&c.IngressClassName, "ingress-class", annotations.DefaultIngressClass, `Name of the ingress class to route through this controller.`)
	flagSet.StringVar(&c.LeaderElectionID, "election-id", "5b374a9e.konghq.com", `Election id to use for status update.`)
	flagSet.StringVar(&c.LeaderElectionNamespace, "election-namespace", "", `Leader election namespace to use when running outside a cluster.`)
	flagSet.StringVar(&c.LeaderElectionForce, "force-leader-election", "", `Set to "enabled" or "disabled" to force a leader election behavior. Behavior is normally determined automatically from other settings.`)
	_ = flagSet.MarkHidden("force-leader-election")
	flagSet.StringSliceVar(&c.FilterTags, "kong-admin-filter-tag", []string{"managed-by-ingress-controller"},
		"Tag(s) in comma-separated format (or specify this flag multiple times). They are used to manage and filter entities in Kong. "+
			"This setting will be silently ignored if the Kong instance has no tags support.")
	flagSet.IntVar(&c.Concurrency, "kong-admin-concurrency", 10, "Max number of concurrent requests sent to Kong's Admin API.")
	flagSet.StringSliceVar(&c.WatchNamespaces, "watch-namespace", nil,
		`Namespace(s) in comma-separated format (or specify this flag multiple times) to watch for Kubernetes resources. Defaults to all namespaces.`)
	flagSet.BoolVar(&c.EmitKubernetesEvents, "emit-kubernetes-events", true, `Emit Kubernetes events for successful configuration applies, translation failures and configuration apply failures on managed objects.`)
	flagSet.StringVar(&c.ClusterDomain, "cluster-domain", consts.DefaultClusterDomain, `The cluster domain. This is used e.g. in generating addresses for upstream services.`)

	// Ingress status
	flagSet.Var(flags.NewValidatedValue(&c.PublishService, namespacedNameFromFlagValue, nnTypeNameOverride), "publish-service",
		`Service fronting Ingress resources in "namespace/name" format. The controller will update Ingress status information with this Service's endpoints.`)
	flagSet.StringSliceVar(&c.PublishStatusAddress, "publish-status-address", []string{},
		`Addresses in comma-separated format (or specify this flag multiple times), for use in lieu of "publish-service" `+
			`when that Service lacks useful address information (for example, in bare-metal environments).`)
	flagSet.Var(flags.NewValidatedValue(&c.PublishServiceUDP, namespacedNameFromFlagValue, nnTypeNameOverride), "publish-service-udp", `Service fronting UDP routing resources in `+
		`"namespace/name" format. The controller will update UDP route status information with this Service's `+
		`endpoints. If omitted, the same Service will be used for both TCP and UDP routes.`)
	flagSet.StringSliceVar(&c.PublishStatusAddressUDP, "publish-status-address-udp", []string{},
		`Addresses in comma-separated format (or specify this flag multiple times), for use in lieu of "publish-service-udp" `+
			`when that Service lacks useful address information (for example, in bare-metal environments).`)

	flagSet.BoolVar(&c.UpdateStatus, "update-status", true,
		`Indicates if the ingress controller should update the status of resources (e.g. IP/Hostname for v1.Ingress, etc.).`)
	flagSet.IntVar(&c.UpdateStatusQueueBufferSize, "update-status-queue-buffer-size", status.DefaultBufferSize, "Buffer size of the underlying channels used to update the status of resources.")

	// Kubernetes API toggling.
	flagSet.BoolVar(&c.IngressNetV1Enabled, "enable-controller-ingress-networkingv1", true, "Enable the networking.k8s.io/v1 Ingress controller.")
	flagSet.BoolVar(&c.IngressClassNetV1Enabled, "enable-controller-ingress-class-networkingv1", true, "Enable the networking.k8s.io/v1 IngressClass controller.")
	flagSet.BoolVar(&c.IngressClassParametersEnabled, "enable-controller-ingress-class-parameters", true, "Enable the IngressClassParameters controller.")
	flagSet.BoolVar(&c.UDPIngressEnabled, "enable-controller-udpingress", true, "Enable the UDPIngress controller.")
	flagSet.BoolVar(&c.TCPIngressEnabled, "enable-controller-tcpingress", true, "Enable the TCPIngress controller.")
	flagSet.BoolVar(&c.KongIngressEnabled, "enable-controller-kongingress", true, "Enable the KongIngress controller.")
	flagSet.BoolVar(&c.KongClusterPluginEnabled, "enable-controller-kongclusterplugin", true, "Enable the KongClusterPlugin controller.")
	flagSet.BoolVar(&c.KongPluginEnabled, "enable-controller-kongplugin", true, "Enable the KongPlugin controller.")
	flagSet.BoolVar(&c.KongConsumerEnabled, "enable-controller-kongconsumer", true, "Enable the KongConsumer controller.")
	flagSet.BoolVar(&c.ServiceEnabled, "enable-controller-service", true, "Enable the Service controller.")
	flagSet.BoolVar(&c.KongUpstreamPolicyEnabled, "enable-controller-kong-upstream-policy", true, "Enable the KongUpstreamPolicy controller.")
	flagSet.BoolVar(&c.GatewayAPIGatewayController, "enable-controller-gwapi-gateway", true, "Enable the Gateway API Gateway controller.")
	flagSet.BoolVar(&c.GatewayAPIHTTPRouteController, "enable-controller-gwapi-httproute", true, "Enable the Gateway API HTTPRoute controller.")
	flagSet.BoolVar(&c.GatewayAPIReferenceGrantController, "enable-controller-gwapi-reference-grant", true, "Enable the Gateway API ReferenceGrant controller.")
	flagSet.BoolVar(&c.GatewayAPIGRPCRouteController, "enable-controller-gwapi-grpcroute", true, "Enable the Gateway API GRPCRoute controller.")
	flagSet.Var(flags.NewValidatedValue(&c.GatewayToReconcile, namespacedNameFromFlagValue, nnTypeNameOverride), "gateway-to-reconcile",
		`Gateway namespaced name in "namespace/name" format. Makes KIC reconcile only the specified Gateway.`)
	flagSet.StringVar(&c.SecretLabelSelector, "secret-label-selector", "",
		`Limits the secrets ingested to those having this label set to "true". If not specified, all secrets are ingested.`)
	flagSet.StringVar(&c.ConfigMapLabelSelector, "configmap-label-selector", consts.DefaultConfigMapSelector,
		`Limits the configmaps ingested to those having this label set to "true".`)
	flagSet.BoolVar(&c.KongServiceFacadeEnabled, "enable-controller-kong-service-facade", true, "Enable the KongServiceFacade controller.")
	flagSet.BoolVar(&c.KongVaultEnabled, "enable-controller-kong-vault", true, "Enable the KongVault controller.")
	flagSet.BoolVar(&c.KongLicenseEnabled, "enable-controller-kong-license", true, "Enable the KongLicense controller.")
	flagSet.BoolVar(&c.KongCustomEntityEnabled, "enable-controller-kong-custom-entity", true, "Enable the KongCustomEntity controller.")

	// Admission Webhook server config
	flagSet.StringVar(&c.AdmissionServer.ListenAddr, "admission-webhook-listen", "off",
		`The address to start admission controller on (ip:port). Setting it to 'off' disables the admission controller.`)
	flagSet.StringVar(&c.AdmissionServer.CertPath, "admission-webhook-cert-file", "",
		`Admission server PEM certificate file path. `+
			fmt.Sprintf(`If both this and the cert value is unset, defaults to %s. `, admission.DefaultAdmissionWebhookCertPath)+`Mutually exclusive with --admission-webhook-cert.`)
	flagSet.StringVar(&c.AdmissionServer.KeyPath, "admission-webhook-key-file", "",
		`Admission server PEM private key file path. `+
			fmt.Sprintf(`If both this and the key value is unset, defaults to %s. `, admission.DefaultAdmissionWebhookKeyPath)+`Mutually exclusive with --admission-webhook-key.`)
	flagSet.StringVar(&c.AdmissionServer.Cert, "admission-webhook-cert", "",
		`Admission server PEM certificate value. Mutually exclusive with --admission-webhook-cert-file.`)
	flagSet.StringVar(&c.AdmissionServer.Key, "admission-webhook-key", "",
		`Admission server PEM private key value. Mutually exclusive with --admission-webhook-key-file.`)

	// Diagnostics
	flagSet.BoolVar(&c.EnableProfiling, "profiling", false, fmt.Sprintf("Enable profiling via web interface host:%v/debug/pprof/.", consts.DiagnosticsPort))
	flagSet.BoolVar(&c.EnableConfigDumps, "dump-config", false, fmt.Sprintf("Enable config dumps via web interface host:%v/debug/config.", consts.DiagnosticsPort))
	flagSet.BoolVar(&c.DumpSensitiveConfig, "dump-sensitive-config", false, "Include credentials and TLS secrets in configs exposed with --dump-config flag.")
	flagSet.IntVar(&c.DiagnosticServerPort, "diagnostic-server-port", consts.DiagnosticsPort, "The port to listen on for the profiling and config dump server.")
	_ = flagSet.MarkHidden("diagnostic-server-port")

	// Drain support
	flagSet.BoolVar(&c.EnableDrainSupport, "enable-drain-support", consts.DefaultEnableDrainSupport, "Include terminating endpoints in Kong upstreams with weight=0 for graceful connection draining.")

	// Combined services from different HTTPRoutes
	flagSet.BoolVar(&c.CombinedServicesFromDifferentHTTPRoutes, "combined-services-from-different-httproutes", false, "Combine rules from different HTTPRoutes that are sharing the same combination of backends to one Kong service to reduce total number of Kong services.")

	// Feature Gates (see FEATURE_GATES.md).
	flagSet.Var(flags.NewMapStringBoolForFeatureGatesWithDefaults(&c.FeatureGates), "feature-gates", "A set of comma separated key=value pairs that describe feature gates for alpha/beta/experimental features. "+
		fmt.Sprintf("See the Feature Gates documentation for information and available options: %s.", managercfg.DocsURL))

	// SIGTERM or SIGINT signal delay.
	flagSet.DurationVar(&c.TermDelay, "term-delay", 0, "The time delay to sleep before SIGTERM or SIGINT will shut down the ingress controller.")

	// Konnect
	flagSet.BoolVar(&c.Konnect.ConfigSynchronizationEnabled, "konnect-sync-enabled", false, "Enable synchronization of data plane configuration with a Konnect control plane.")
	flagSet.BoolVar(&c.Konnect.LicenseSynchronizationEnabled, "konnect-licensing-enabled", false, "Retrieve licenses from Konnect if available. Overrides licenses provided via the environment.")
	flagSet.BoolVar(&c.Konnect.LicenseStorageEnabled, "konnect-license-storage-enabled", true, "Store licenses fetched from Konnect to Secrets locally to use them later when connection to Konnect is broken. Only effective when --konnect-licensing-enabled is true.")
	flagSet.DurationVar(&c.Konnect.InitialLicensePollingPeriod, "konnect-initial-license-polling-period", license.DefaultInitialPollingPeriod, "Polling period to be used before the first license is retrieved.")
	flagSet.DurationVar(&c.Konnect.LicensePollingPeriod, "konnect-license-polling-period", license.DefaultPollingPeriod, "Polling period to be used after the first license is retrieved.")
	flagSet.StringVar(&c.Konnect.ControlPlaneID, "konnect-control-plane-id", "", "An ID of a control plane that is to be synchronized with data plane configuration.")
	flagSet.StringVar(&c.Konnect.Address, "konnect-address", "https://us.kic.api.konghq.com", "Base address of Konnect API.")
	flagSet.StringVar(&c.Konnect.TLSClient.Cert, "konnect-tls-client-cert", "", "Konnect TLS client certificate.")
	flagSet.StringVar(&c.Konnect.TLSClient.CertFile, "konnect-tls-client-cert-file", "", "Konnect TLS client certificate file path.")
	flagSet.StringVar(&c.Konnect.TLSClient.Key, "konnect-tls-client-key", "", "Konnect TLS client key.")
	flagSet.StringVar(&c.Konnect.TLSClient.KeyFile, "konnect-tls-client-key-file", "", "Konnect TLS client key file path.")
	flagSet.DurationVar(&c.Konnect.UploadConfigPeriod, "konnect-upload-config-period", managercfg.DefaultKonnectConfigUploadPeriod, "Period of uploading Kong configuration.")
	flagSet.DurationVar(&c.Konnect.RefreshNodePeriod, "konnect-refresh-node-period", konnect.DefaultRefreshNodePeriod, "Period of uploading status of KIC and controlled Kong instances.")
	flagSet.BoolVar(&c.Konnect.ConsumersSyncDisabled, "konnect-disable-consumers-sync", false, "Disable synchronization of consumers with Konnect.")

	// Deprecated flags.
	flagSet.StringVar(&c.Konnect.ControlPlaneID, "konnect-runtime-group-id", "", "Use --konnect-control-plane-id instead.")
	_ = flagSet.MarkDeprecated("konnect-runtime-group-id", "Use --konnect-control-plane-id instead.")

	_ = flagSet.String("gateway-discovery-dns-strategy", "", "DNS strategy to use when creating Gateway's Admin API addresses. One of: ip, service, pod.")
	_ = flagSet.MarkDeprecated("gateway-discovery-dns-strategy", "this setting is deprecated and has no effect, now it always works out of the box (without adjustments).")

	c.flagSet = flagSet
}
