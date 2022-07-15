package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
)

// -----------------------------------------------------------------------------
// Controller Manager - Config
// -----------------------------------------------------------------------------

// Config collects all configuration that the controller manager takes from the environment.
type Config struct {
	// See flag definitions in RegisterFlags(...) for documentation of the fields defined here.

	// Logging configurations
	LogLevel            string
	LogFormat           string
	LogReduceRedundancy bool

	// Kong high-level controller manager configurations
	KongAdminAPIConfig                adminapi.HTTPClientOpts
	KongAdminInitializationRetries    uint
	KongAdminInitializationRetryDelay time.Duration
	KongAdminToken                    string
	KongWorkspace                     string
	AnonymousReports                  bool
	EnableReverseSync                 bool
	SyncPeriod                        time.Duration
	SkipCACertificates                bool

	// Kong Proxy configurations
	APIServerHost            string
	APIServerQPS             int
	APIServerBurst           int
	MetricsAddr              string
	ProbeAddr                string
	KongAdminURL             string
	ProxySyncSeconds         float32
	ProxyTimeoutSeconds      float32
	KongCustomEntitiesSecret string

	// Kubernetes configurations
	KubeconfigPath          string
	IngressClassName        string
	EnableLeaderElection    bool
	LeaderElectionNamespace string
	LeaderElectionID        string
	Concurrency             int
	FilterTags              []string
	WatchNamespaces         []string

	// Ingress status
	PublishService       string
	PublishStatusAddress []string
	UpdateStatus         bool

	// Kubernetes API toggling
	IngressExtV1beta1Enabled bool
	IngressNetV1beta1Enabled bool
	IngressNetV1Enabled      bool
	IngressClassNetV1Enabled bool
	UDPIngressEnabled        bool
	TCPIngressEnabled        bool
	KongIngressEnabled       bool
	KnativeIngressEnabled    bool
	KongClusterPluginEnabled bool
	KongPluginEnabled        bool
	KongConsumerEnabled      bool
	ServiceEnabled           bool

	// Admission Webhook server config
	AdmissionServer admission.ServerConfig

	// Diagnostics and performance
	EnableProfiling     bool
	EnableConfigDumps   bool
	DumpSensitiveConfig bool

	// Feature Gates
	FeatureGates map[string]bool

	// TermDelay is the time.Duration which the controller manager will wait
	// after receiving SIGTERM or SIGINT before shutting down. This can be
	// helpful for advanced cases with load-balancers so that the ingress
	// controller can be gracefully removed/drained from their rotation.
	TermDelay time.Duration
}

// -----------------------------------------------------------------------------
// Controller Manager - Config - Methods
// -----------------------------------------------------------------------------

// FlagSet binds the provided Config to commandline flags.
func (c *Config) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("", pflag.ExitOnError)

	// Logging configurations
	flagSet.StringVar(&c.LogLevel, "log-level", "info", `Level of logging for the controller. Allowed values are trace, debug, info, warn, error, fatal and panic.`)
	flagSet.StringVar(&c.LogFormat, "log-format", "text", `Format of logs of the controller. Allowed values are text and json.`)
	flagSet.BoolVar(&c.LogReduceRedundancy, "debug-log-reduce-redundancy", false, `If enabled, repetitive log entries are suppressed. Built for testing environments - production use not recommended.`)
	flagSet.MarkHidden("debug-log-reduce-redundancy") //nolint:errcheck

	// Kong high-level controller manager configurations
	flagSet.BoolVar(&c.KongAdminAPIConfig.TLSSkipVerify, "kong-admin-tls-skip-verify", false, "Disable verification of TLS certificate of Kong's Admin endpoint.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSServerName, "kong-admin-tls-server-name", "", "SNI name to use to verify the certificate presented by Kong in TLS.")
	flagSet.StringVar(&c.KongAdminAPIConfig.CACertPath, "kong-admin-ca-cert-file", "", `Path to PEM-encoded CA certificate file to verify Kong's Admin SSL certificate.`)
	flagSet.StringVar(&c.KongAdminAPIConfig.CACert, "kong-admin-ca-cert", "", `PEM-encoded CA certificate to verify Kong's Admin SSL certificate.`)

	flagSet.StringSliceVar(&c.KongAdminAPIConfig.Headers, "kong-admin-header", nil, `add a header (key:value) to every Admin API call, this flag can be used multiple times to specify multiple headers`)
	flagSet.UintVar(&c.KongAdminInitializationRetries, "kong-admin-init-retries", 60, "Number of attempts that will be made initially on controller startup to connect to the Kong Admin API")
	flagSet.DurationVar(&c.KongAdminInitializationRetryDelay, "kong-admin-init-retry-delay", time.Second*1, "The time delay between every attempt (on controller startup) to connect to the Kong Admin API")
	flagSet.StringVar(&c.KongAdminToken, "kong-admin-token", "", `The Kong Enterprise RBAC token used by the controller.`)
	flagSet.StringVar(&c.KongWorkspace, "kong-workspace", "", "Kong Enterprise workspace to configure. Leave this empty if not using Kong workspaces.")
	flagSet.BoolVar(&c.AnonymousReports, "anonymous-reports", true, `Send anonymized usage data to help improve Kong`)
	flagSet.BoolVar(&c.EnableReverseSync, "enable-reverse-sync", false, `Send configuration to Kong even if the configuration checksum has not changed since previous update.`)
	flagSet.DurationVar(&c.SyncPeriod, "sync-period", time.Hour*48, `Relist and confirm cloud resources this often`) // 48 hours derived from controller-runtime defaults
	flagSet.BoolVar(&c.SkipCACertificates, "skip-ca-certificates", false, `disable syncing CA certificate syncing (for use with multi-workspace environments)`)

	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClientCertPath, "kong-admin-tls-client-cert-file", "", "mTLS client certificate file for authentication.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClientKeyPath, "kong-admin-tls-client-key-file", "", "mTLS client key file for authentication.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClientCert, "kong-admin-tls-client-cert", "", "mTLS client certificate for authentication.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSClientKey, "kong-admin-tls-client-key", "", "mTLS client key for authentication.")

	// Kong Proxy and Proxy Cache configurations
	flagSet.StringVar(&c.APIServerHost, "apiserver-host", "", `The Kubernetes API server URL. If not set, the controller will use cluster config discovery.`)
	flagSet.IntVar(&c.APIServerQPS, "apiserver-qps", 100, "The Kubernetes API RateLimiter maximum queries per second")
	flagSet.IntVar(&c.APIServerBurst, "apiserver-burst", 300, "The Kubernetes API RateLimiter maximum burst queries per second")
	flagSet.StringVar(&c.MetricsAddr, "metrics-bind-address", fmt.Sprintf(":%v", MetricsPort), "The address the metric endpoint binds to.")
	flagSet.StringVar(&c.ProbeAddr, "health-probe-bind-address", fmt.Sprintf(":%v", HealthzPort), "The address the probe endpoint binds to.")
	flagSet.StringVar(&c.KongAdminURL, "kong-admin-url", "http://localhost:8001", `The Kong Admin URL to connect to in the format "protocol://address:port".`)
	flagSet.Float32Var(&c.ProxySyncSeconds, "proxy-sync-seconds", dataplane.DefaultSyncSeconds,
		"Define the rate (in seconds) in which configuration updates will be applied to the Kong Admin API.",
	)
	flagSet.Float32Var(&c.ProxyTimeoutSeconds, "proxy-timeout-seconds", dataplane.DefaultTimeoutSeconds,
		"Sets the timeout (in seconds) for all requests to Kong's Admin API.",
	)
	flagSet.StringVar(&c.KongCustomEntitiesSecret, "kong-custom-entities-secret", "", `A Secret containing custom entities for DB-less mode, in "namespace/name" format`)

	// Kubernetes configurations
	flagSet.StringVar(&c.KubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file.")
	flagSet.StringVar(&c.IngressClassName, "ingress-class", annotations.DefaultIngressClass, `Name of the ingress class to route through this controller.`)
	flagSet.BoolVar(&c.EnableLeaderElection, "leader-elect", false, "DEPRECATED as of 2.1.0 leader election behavior is determined automatically and this flag has no effect")
	flagSet.StringVar(&c.LeaderElectionID, "election-id", "5b374a9e.konghq.com", `Election id to use for status update.`)
	flagSet.StringVar(&c.LeaderElectionNamespace, "election-namespace", "", `Leader election namespace to use when running outside a cluster`)
	flagSet.StringSliceVar(&c.FilterTags, "kong-admin-filter-tag", []string{"managed-by-ingress-controller"}, "The tag used to manage and filter entities in Kong. This flag can be specified multiple times to specify multiple tags. This setting will be silently ignored if the Kong instance has no tags support.")
	flagSet.IntVar(&c.Concurrency, "kong-admin-concurrency", 10, "Max number of concurrent requests sent to Kong's Admin API.")
	flagSet.StringSliceVar(&c.WatchNamespaces, "watch-namespace", nil,
		`Namespace(s) to watch for Kubernetes resources. Defaults to all namespaces. To watch multiple namespaces, use
		a comma-separated list of namespaces.`)

	// Ingress status
	flagSet.StringVar(&c.PublishService, "publish-service", "", `Service fronting Ingress resources in "namespace/name"
			format. The controller will update Ingress status information with this Service's endpoints.`)
	flagSet.StringSliceVar(&c.PublishStatusAddress, "publish-status-address", []string{}, `User-provided addresses in
			comma-separated string format, for use in lieu of "publish-service" when that Service lacks useful address
			information (for example, in bare-metal environments).`)
	flagSet.BoolVar(&c.UpdateStatus, "update-status", true,
		`Indicates if the ingress controller should update the status of resources (e.g. IP/Hostname for v1.Ingress, e.t.c.)`)

	// Kubernetes API toggling
	flagSet.BoolVar(&c.IngressNetV1Enabled, "enable-controller-ingress-networkingv1", true, "Enable the networking.k8s.io/v1 Ingress controller.")
	flagSet.BoolVar(&c.IngressClassNetV1Enabled, "enable-controller-ingress-class-networkingv1", true, "Enable the networking.k8s.io/v1 IngressClass controller.")
	flagSet.BoolVar(&c.IngressNetV1beta1Enabled, "enable-controller-ingress-networkingv1beta1", true, "Enable the networking.k8s.io/v1beta1 Ingress controller.")
	flagSet.BoolVar(&c.IngressExtV1beta1Enabled, "enable-controller-ingress-extensionsv1beta1", true, "Enable the extensions/v1beta1 Ingress controller.")
	flagSet.BoolVar(&c.UDPIngressEnabled, "enable-controller-udpingress", true, "Enable the UDPIngress controller.")
	flagSet.BoolVar(&c.TCPIngressEnabled, "enable-controller-tcpingress", true, "Enable the TCPIngress controller.")
	flagSet.BoolVar(&c.KnativeIngressEnabled, "enable-controller-knativeingress", true, "Enable the KnativeIngress controller.")
	flagSet.BoolVar(&c.KongIngressEnabled, "enable-controller-kongingress", true, "Enable the KongIngress controller.")
	flagSet.BoolVar(&c.KongClusterPluginEnabled, "enable-controller-kongclusterplugin", true, "Enable the KongClusterPlugin controller.")
	flagSet.BoolVar(&c.KongPluginEnabled, "enable-controller-kongplugin", true, "Enable the KongPlugin controller.")
	flagSet.BoolVar(&c.KongConsumerEnabled, "enable-controller-kongconsumer", true, "Enable the KongConsumer controller. ")
	flagSet.BoolVar(&c.ServiceEnabled, "enable-controller-service", true, "Enable the Service controller.")

	// Admission Webhook server config
	flagSet.StringVar(&c.AdmissionServer.ListenAddr, "admission-webhook-listen", "off",
		`The address to start admission controller on (ip:port).  Setting it to 'off' disables the admission controller.`)
	flagSet.StringVar(&c.AdmissionServer.CertPath, "admission-webhook-cert-file", "",
		`admission server PEM certificate file path; `+
			`if both this and the cert value is unset, defaults to `+admission.DefaultAdmissionWebhookCertPath)
	flagSet.StringVar(&c.AdmissionServer.KeyPath, "admission-webhook-key-file", "",
		`admission server PEM private key file path; `+
			`if both this and the key value is unset, defaults to `+admission.DefaultAdmissionWebhookKeyPath)
	flagSet.StringVar(&c.AdmissionServer.Cert, "admission-webhook-cert", "",
		`admission server PEM certificate value`)
	flagSet.StringVar(&c.AdmissionServer.Key, "admission-webhook-key", "",
		`admission server PEM private key value`)

	// Diagnostics
	flagSet.BoolVar(&c.EnableProfiling, "profiling", false, fmt.Sprintf("Enable profiling via web interface host:%v/debug/pprof/", DiagnosticsPort))
	flagSet.BoolVar(&c.EnableConfigDumps, "dump-config", false, fmt.Sprintf("Enable config dumps via web interface host:%v/debug/config", DiagnosticsPort))
	flagSet.BoolVar(&c.DumpSensitiveConfig, "dump-sensitive-config", false, "Include credentials and TLS secrets in configs exposed with --dump-config")

	// Feature Gates (see FEATURE_GATES.md)
	flagSet.Var(cliflag.NewMapStringBool(&c.FeatureGates), "feature-gates", "A set of key=value pairs that describe feature gates for alpha/beta/experimental features. "+
		fmt.Sprintf("See the Feature Gates documentation for information and available options: %s", featureGatesDocsURL))

	// SIGTERM or SIGINT signal delay
	flagSet.DurationVar(&c.TermDelay, "term-delay", time.Second*0, "The time delay to sleep before SIGTERM or SIGINT will shut down the Ingress Controller")

	// Deprecated (to be removed in future releases)
	flagSet.Float32Var(&c.ProxySyncSeconds, "sync-rate-limit", dataplane.DefaultSyncSeconds,
		"Define the rate (in seconds) in which configuration updates will be applied to the Kong Admin API (DEPRECATED, use --proxy-sync-seconds instead)",
	)
	flagSet.Int("stderrthreshold", 0, "DEPRECATED: has no effect and will be removed in future releases (see github issue #1297)")
	flagSet.Bool("update-status-on-shutdown", false, `DEPRECATED: no longer has any effect and will be removed in a later release (see github issue #1304)`)

	return flagSet
}

func (c *Config) GetKongClient(ctx context.Context) (*kong.Client, error) {
	if c.KongAdminToken != "" {
		c.KongAdminAPIConfig.Headers = append(c.KongAdminAPIConfig.Headers, "kong-admin-token:"+c.KongAdminToken)
	}
	httpclient, err := adminapi.MakeHTTPClient(&c.KongAdminAPIConfig)
	if err != nil {
		return nil, err
	}

	return adminapi.GetKongClientForWorkspace(ctx, c.KongAdminURL, c.KongWorkspace, httpclient)
}

func (c *Config) GetKubeconfig() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags(c.APIServerHost, c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Configure k8s client rate-limiting
	config.QPS = float32(c.APIServerQPS)
	config.Burst = c.APIServerBurst

	return config, err
}

func (c *Config) GetKubeClient() (client.Client, error) {
	conf, err := c.GetKubeconfig()
	if err != nil {
		return nil, err
	}
	return client.New(conf, client.Options{})
}
