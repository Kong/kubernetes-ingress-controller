package config

import (
	"fmt"
	"os"
	"time"

	"github.com/samber/mo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config/types"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

const (
	// LeaderElectionEnabled is a constant that represents a value that should be used to enable leader election.
	LeaderElectionEnabled = "enabled"
	// LeaderElectionDisabled is a constant that represents a value that should be used to disable leader election.
	LeaderElectionDisabled = "disabled"
)

// OptionalNamespacedName is a type that represents a NamespacedName that can be omitted in config.
type OptionalNamespacedName = mo.Option[k8stypes.NamespacedName]

// Opt is a function that modifies a Config.
type Opt func(*Config)

// Config is the configuration for the Kong Ingress Controller.
type Config struct {
	// Logging configurations
	LogLevel  string
	LogFormat string

	// Kong high-level controller manager configurations
	KongAdminAPIConfig                adminapi.ClientOpts
	KongAdminInitializationRetries    uint
	KongAdminInitializationRetryDelay time.Duration
	KongAdminToken                    string
	KongAdminTokenPath                string
	KongWorkspace                     string
	AnonymousReports                  bool
	EnableReverseSync                 bool
	UseLastValidConfigForFallback     bool
	SyncPeriod                        time.Duration
	SkipCACertificates                bool
	CacheSyncTimeout                  time.Duration
	GracefulShutdownTimeout           *time.Duration

	// Kong Proxy configurations
	APIServerHost                          string
	APIServerQPS                           int
	APIServerBurst                         int
	APIServerCAData                        []byte
	APIServerCertData                      []byte
	APIServerKeyData                       []byte
	MetricsAddr                            string
	MetricsAccessFilter                    cfgtypes.MetricsAccessFilter
	ProbeAddr                              string
	KongAdminURLs                          []string
	KongAdminSvc                           OptionalNamespacedName
	GatewayDiscoveryReadinessCheckInterval time.Duration
	GatewayDiscoveryReadinessCheckTimeout  time.Duration
	KongAdminSvcPortNames                  []string
	ProxySyncSeconds                       float32
	InitCacheSyncDuration                  time.Duration
	ProxyTimeoutSeconds                    float32

	// Kubernetes configurations
	KubeconfigPath           string
	IngressClassName         string
	LeaderElectionNamespace  string
	LeaderElectionID         string
	LeaderElectionForce      string
	Concurrency              int
	FilterTags               []string
	WatchNamespaces          []string
	GatewayAPIControllerName string
	Impersonate              string
	EmitKubernetesEvents     bool
	ClusterDomain            string

	// Ingress status
	PublishServiceUDP       OptionalNamespacedName
	PublishService          OptionalNamespacedName
	PublishStatusAddress    []string
	PublishStatusAddressUDP []string

	UpdateStatus                bool
	UpdateStatusQueueBufferSize int

	// Kubernetes API toggling
	IngressNetV1Enabled           bool
	IngressClassNetV1Enabled      bool
	IngressClassParametersEnabled bool
	UDPIngressEnabled             bool
	TCPIngressEnabled             bool
	KongIngressEnabled            bool
	KongClusterPluginEnabled      bool
	KongPluginEnabled             bool
	KongConsumerEnabled           bool
	ServiceEnabled                bool
	KongUpstreamPolicyEnabled     bool
	KongServiceFacadeEnabled      bool
	KongVaultEnabled              bool
	KongLicenseEnabled            bool
	KongCustomEntityEnabled       bool

	// Gateway API toggling.
	GatewayAPIGatewayController        bool
	GatewayAPIHTTPRouteController      bool
	GatewayAPIReferenceGrantController bool
	GatewayAPIGRPCRouteController      bool

	// GatewayToReconcile specifies the Gateway to be reconciled.
	GatewayToReconcile OptionalNamespacedName

	// SecretLabelSelector specifies the label which will be used to limit the ingestion of secrets. Only those that have this label set to "true" will be ingested.
	SecretLabelSelector string

	// ConfigMapLabelSelector specifies the label which will be used to limit the ingestion of configmaps. Only those that have this label set to "true" will be ingested.
	ConfigMapLabelSelector string

	// Admission Webhook server config
	AdmissionServer admission.ServerConfig

	// Diagnostics and performance
	EnableProfiling      bool
	EnableConfigDumps    bool
	DumpSensitiveConfig  bool
	DiagnosticServerPort int

	// Feature Gates
	FeatureGates map[string]bool

	// TermDelay is the time.Duration which the controller manager will wait
	// after receiving SIGTERM or SIGINT before shutting down. This can be
	// helpful for advanced cases with load-balancers so that the ingress
	// controller can be gracefully removed/drained from their rotation.
	TermDelay time.Duration

	Konnect adminapi.KonnectConfig

	// Override default telemetry settings (e.g. for testing). They aren't exposed in the CLI.
	SplunkEndpoint                   string
	SplunkEndpointInsecureSkipVerify bool
	TelemetryPeriod                  time.Duration
}

// Resolve modifies the Config object in place by resolving any values that are not set directly (e.g. reading a file
// for a token).
func (c *Config) Resolve() error {
	if c.KongAdminTokenPath != "" {
		token, err := os.ReadFile(c.KongAdminTokenPath)
		if err != nil {
			return fmt.Errorf("failed to read --kong-admin-token-file from path '%s': %w", c.KongAdminTokenPath, err)
		}
		c.KongAdminToken = string(token)
	}
	return nil
}

// GetKubeconfig returns a Kubernetes REST config object based on the configuration.
func GetKubeconfig(c Config) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags(c.APIServerHost, c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Configure k8s client rate-limiting
	config.QPS = float32(c.APIServerQPS)
	config.Burst = c.APIServerBurst

	if c.APIServerCertData != nil {
		config.CertData = c.APIServerCertData
	}
	if c.APIServerCAData != nil {
		config.CAData = c.APIServerCAData
	}
	if c.APIServerKeyData != nil {
		config.KeyData = c.APIServerKeyData
	}
	if c.Impersonate != "" {
		config.Impersonate.UserName = c.Impersonate
	}

	config.UserAgent = metadata.UserAgent()

	return config, err
}

// GetKubeClient returns a Kubernetes client based on the configuration.
func GetKubeClient(c Config) (client.Client, error) {
	conf, err := GetKubeconfig(c)
	if err != nil {
		return nil, err
	}
	return client.New(conf, client.Options{})
}
