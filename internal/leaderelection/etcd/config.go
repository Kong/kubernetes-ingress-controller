package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	// Environment variable names for etcd configuration.
	EnvEtcdEndpoints = "ETCD_ENDPOINTS"
	EnvEtcdCertFile  = "ETCD_CERT_FILE"
	EnvEtcdKeyFile   = "ETCD_KEY_FILE"
	EnvEtcdCAFile    = "ETCD_CA_FILE"
	EnvEtcdUsername  = "ETCD_USERNAME"
	EnvEtcdPassword  = "ETCD_PASSWORD"

	// Default values.
	DefaultDialTimeout      = 5 * time.Second
	DefaultSessionTTL       = 15 // seconds
	DefaultElectionPrefix   = "/kong-ingress-controller/leader-election"
	DefaultKeepAliveTime    = 10 * time.Second
	DefaultKeepAliveTimeout = 3 * time.Second
)

// Config holds the configuration for connecting to etcd.
type Config struct {
	// Endpoints is a list of etcd endpoints to connect to.
	Endpoints []string

	// TLS configuration.
	CertFile string
	KeyFile  string
	CAFile   string

	// Authentication.
	Username string
	Password string

	// Timeouts.
	DialTimeout time.Duration

	// Session TTL in seconds for the lease used in leader election.
	SessionTTL int

	// ElectionPrefix is the key prefix used for leader election in etcd.
	ElectionPrefix string
}

// NewConfigFromEnv creates a new Config from environment variables.
func NewConfigFromEnv() (*Config, error) {
	endpoints := os.Getenv(EnvEtcdEndpoints)
	if endpoints == "" {
		return nil, fmt.Errorf("environment variable %s is required", EnvEtcdEndpoints)
	}

	cfg := &Config{
		Endpoints:      strings.Split(endpoints, ","),
		CertFile:       os.Getenv(EnvEtcdCertFile),
		KeyFile:        os.Getenv(EnvEtcdKeyFile),
		CAFile:         os.Getenv(EnvEtcdCAFile),
		Username:       os.Getenv(EnvEtcdUsername),
		Password:       os.Getenv(EnvEtcdPassword),
		DialTimeout:    DefaultDialTimeout,
		SessionTTL:     DefaultSessionTTL,
		ElectionPrefix: DefaultElectionPrefix,
	}

	return cfg, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("at least one etcd endpoint is required")
	}

	// If TLS cert is provided, key must also be provided and vice versa.
	if (c.CertFile != "" && c.KeyFile == "") || (c.CertFile == "" && c.KeyFile != "") {
		return fmt.Errorf("both ETCD_CERT_FILE and ETCD_KEY_FILE must be provided together")
	}

	return nil
}

// ToClientConfig converts the Config to an etcd clientv3.Config.
func (c *Config) ToClientConfig() (clientv3.Config, error) {
	cfg := clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: c.DialTimeout,
		Username:    c.Username,
		Password:    c.Password,
	}

	// Configure TLS if certificates are provided.
	if c.CertFile != "" && c.KeyFile != "" {
		tlsConfig, err := c.buildTLSConfig()
		if err != nil {
			return clientv3.Config{}, fmt.Errorf("failed to build TLS config: %w", err)
		}
		cfg.TLS = tlsConfig
	}

	return cfg, nil
}

// buildTLSConfig creates a TLS configuration from the cert/key/CA files.
func (c *Config) buildTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Load CA certificate if provided.
	if c.CAFile != "" {
		caCert, err := os.ReadFile(c.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

// IsConfigured returns true if etcd configuration is available from environment.
func IsConfigured() bool {
	return os.Getenv(EnvEtcdEndpoints) != ""
}
