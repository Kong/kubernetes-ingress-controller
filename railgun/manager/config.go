package manager

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/pkg/adminapi"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/spf13/pflag"
	apiv1 "k8s.io/api/core/v1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Config
// -----------------------------------------------------------------------------

// Config collects all configuration that the controller manager takes from the environment.
// BUG: the above is not 100% accurate today - controllers read some settings from environment variables directly
type Config struct {
	// See flag definitions in RegisterFlags(...) for documentation of the fields defined here.

	// Logging configurations
	LogLevel  string
	LogFormat string

	// Kong high-level controller manager configurations
	KongAdminAPIConfig adminapi.HTTPClientOpts
	KongAdminToken     string
	KongStateEnabled   util.EnablementStatus
	KongWorkspace      string
	AnonymousReports   bool

	// Kong Proxy configurations
	APIServerHost string
	MetricsAddr   string
	ProbeAddr     string
	KongAdminURL  string

	// Kubernetes configurations
	KubeconfigPath       string
	IngressClassName     string
	EnableLeaderElection bool
	LeaderElectionID     string
	Concurrency          int
	FilterTag            string
	WatchNamespace       string

	// Kubernetes API toggling
	IngressExtV1beta1Enabled util.EnablementStatus
	IngressNetV1beta1Enabled util.EnablementStatus
	IngressNetV1Enabled      util.EnablementStatus
	UDPIngressEnabled        util.EnablementStatus
	TCPIngressEnabled        util.EnablementStatus
	KongIngressEnabled       util.EnablementStatus
	KongClusterPluginEnabled util.EnablementStatus
	KongPluginEnabled        util.EnablementStatus
	KongConsumerEnabled      util.EnablementStatus
	ServiceEnabled           util.EnablementStatus

	// "Classless" API support
	ProcessClasslessIngressV1      bool
	ProcessClasslessIngressV1Beta1 bool
	ProcessClasslessKongConsumer   bool
}

// -----------------------------------------------------------------------------
// Controller Manager - Config - Constants & Vars
// -----------------------------------------------------------------------------

// onOffUsage is used to indicate what textual options are used to enable or disable a feature.
const onOffUsage = "Can be one of [enabled, disabled]."

// -----------------------------------------------------------------------------
// Controller Manager - Config - Methods
// -----------------------------------------------------------------------------

// FlagSet binds the provided Config to commandline flags.
func (c *Config) FlagSet() *pflag.FlagSet {

	flagSet := flagSet{*pflag.NewFlagSet("", pflag.ExitOnError)}

	// Logging configurations
	flagSet.StringVar(&c.LogLevel, "log-level", "info", `Level of logging for the controller. Allowed values are trace, debug, info, warn, error, fatal and panic.`)
	flagSet.StringVar(&c.LogFormat, "log-format", "text", `Format of logs of the controller. Allowed values are text and json.`)

	// Kong high-level controller manager configurations
	flagSet.BoolVar(&c.KongAdminAPIConfig.TLSSkipVerify, "kong-admin-tls-skip-verify", false, "Disable verification of TLS certificate of Kong's Admin endpoint.")
	flagSet.StringVar(&c.KongAdminAPIConfig.TLSServerName, "kong-admin-tls-server-name", "", "SNI name to use to verify the certificate presented by Kong in TLS.")
	flagSet.StringVar(&c.KongAdminAPIConfig.CACertPath, "kong-admin-ca-cert-file", "", `Path to PEM-encoded CA certificate file to verify Kong's Admin SSL certificate.`)
	flagSet.StringVar(&c.KongAdminAPIConfig.CACert, "kong-admin-ca-cert", "", `PEM-encoded CA certificate to verify Kong's Admin SSL certificate.`)
	flagSet.StringSliceVar(&c.KongAdminAPIConfig.Headers, "kong-admin-header", nil, `add a header (key:value) to every Admin API call, this flag can be used multiple times to specify multiple headers`)
	flagSet.StringVar(&c.KongAdminToken, "kong-admin-token", "", `The Kong Enterprise RBAC token used by the controller.`)
	flagSet.enablementStatusVar(&c.KongStateEnabled, "controller-kongstate", util.EnablementStatusEnabled, "Enable or disable the KongState controller. "+onOffUsage)
	flagSet.StringVar(&c.KongWorkspace, "kong-workspace", "", "Kong Enterprise workspace to configure. Leave this empty if not using Kong workspaces.")
	flagSet.BoolVar(&c.AnonymousReports, "anonymous-reports", true, `Send anonymized usage data to help improve Kong`)

	// Kong Proxy configurations
	flagSet.StringVar(&c.APIServerHost, "apiserver-host", "", `The Kubernetes API server URL. If not set, the controller will use cluster config discovery.`)
	flagSet.StringVar(&c.MetricsAddr, "metrics-bind-address", fmt.Sprintf(":%v", MetricsPort), "The address the metric endpoint binds to.")
	flagSet.StringVar(&c.ProbeAddr, "health-probe-bind-address", fmt.Sprintf(":%v", HealthzPort), "The address the probe endpoint binds to.")
	flagSet.StringVar(&c.KongAdminURL, "kong-admin-url", "http://localhost:8001", `The Kong Admin URL to connect to in the format "protocol://address:port".`)

	// Kubernetes configurations
	flagSet.StringVar(&c.KubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file.")
	flagSet.StringVar(&c.IngressClassName, "ingress-class", annotations.DefaultIngressClass, `Name of the ingress class to route through this controller.`)
	flagSet.BoolVar(&c.EnableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flagSet.StringVar(&c.LeaderElectionID, "election-id", "5b374a9e.konghq.com", `Election id to use for status update.`)
	flagSet.StringVar(&c.FilterTag, "kong-filter-tag", "managed-by-railgun", "TODO")
	flagSet.IntVar(&c.Concurrency, "kong-concurrency", 10, "TODO")
	flagSet.StringVar(&c.WatchNamespace, "watch-namespace", apiv1.NamespaceAll, "Namespace to watch for Kubernetes resources. Defaults to all namespaces.")

	// Kubernetes API toggling
	flagSet.enablementStatusVar(&c.IngressNetV1Enabled, "controller-ingress-networkingv1", util.EnablementStatusEnabled, "Enable or disable the Ingress controller (using API version networking.k8s.io/v1)."+onOffUsage)
	flagSet.enablementStatusVar(&c.IngressNetV1beta1Enabled, "controller-ingress-networkingv1beta1", util.EnablementStatusDisabled, "Enable or disable the Ingress controller (using API version networking.k8s.io/v1beta1)."+onOffUsage)
	flagSet.enablementStatusVar(&c.IngressExtV1beta1Enabled, "controller-ingress-extensionsv1beta1", util.EnablementStatusDisabled, "Enable or disable the Ingress controller (using API version extensions/v1beta1)."+onOffUsage)
	flagSet.enablementStatusVar(&c.UDPIngressEnabled, "controller-udpingress", util.EnablementStatusDisabled, "Enable or disable the UDPIngress controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.TCPIngressEnabled, "controller-tcpingress", util.EnablementStatusDisabled, "Enable or disable the TCPIngress controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.KongIngressEnabled, "controller-kongingress", util.EnablementStatusEnabled, "Enable or disable the KongIngress controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.KongClusterPluginEnabled, "controller-kongclusterplugin", util.EnablementStatusDisabled, "Enable or disable the KongClusterPlugin controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.KongPluginEnabled, "controller-kongplugin", util.EnablementStatusDisabled, "Enable or disable the KongPlugin controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.KongConsumerEnabled, "controller-kongconsumer", util.EnablementStatusDisabled, "Enable or disable the KongConsumer controller. "+onOffUsage)
	flagSet.enablementStatusVar(&c.ServiceEnabled, "controller-service", util.EnablementStatusEnabled, "Enable or disable the Service controller. "+onOffUsage)

	// "Classless" API support
	flagSet.BoolVar(&c.ProcessClasslessIngressV1Beta1, "process-classless-ingress-v1beta1", false, `Process v1beta1 Ingress resources with no class annotation.`)
	flagSet.BoolVar(&c.ProcessClasslessIngressV1, "process-classless-ingress-v1", false, `Process v1 Ingress resources with no class annotation.`)
	flagSet.BoolVar(&c.ProcessClasslessKongConsumer, "process-classless-kong-consumer", false, `Process KongConsumer resources with no class annotation.`)

	return &flagSet.FlagSet
}
