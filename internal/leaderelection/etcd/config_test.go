package etcd_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/leaderelection/etcd"
)

func TestNewConfigFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedErr    string
		expectedConfig *etcd.Config
	}{
		{
			name:        "missing ETCD_ENDPOINTS returns error",
			envVars:     map[string]string{},
			expectedErr: "environment variable ETCD_ENDPOINTS is required",
		},
		{
			name: "single endpoint",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "http://localhost:2379",
			},
			expectedConfig: &etcd.Config{
				Endpoints:      []string{"http://localhost:2379"},
				DialTimeout:    etcd.DefaultDialTimeout,
				SessionTTL:     etcd.DefaultSessionTTL,
				ElectionPrefix: etcd.DefaultElectionPrefix,
			},
		},
		{
			name: "multiple endpoints",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "http://etcd-0:2379,http://etcd-1:2379,http://etcd-2:2379",
			},
			expectedConfig: &etcd.Config{
				Endpoints:      []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
				DialTimeout:    etcd.DefaultDialTimeout,
				SessionTTL:     etcd.DefaultSessionTTL,
				ElectionPrefix: etcd.DefaultElectionPrefix,
			},
		},
		{
			name: "with TLS configuration",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "https://localhost:2379",
				"ETCD_CERT_FILE": "/path/to/client.crt",
				"ETCD_KEY_FILE":  "/path/to/client.key",
				"ETCD_CA_FILE":   "/path/to/ca.crt",
			},
			expectedConfig: &etcd.Config{
				Endpoints:      []string{"https://localhost:2379"},
				CertFile:       "/path/to/client.crt",
				KeyFile:        "/path/to/client.key",
				CAFile:         "/path/to/ca.crt",
				DialTimeout:    etcd.DefaultDialTimeout,
				SessionTTL:     etcd.DefaultSessionTTL,
				ElectionPrefix: etcd.DefaultElectionPrefix,
			},
		},
		{
			name: "with authentication",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "http://localhost:2379",
				"ETCD_USERNAME":  "testuser",
				"ETCD_PASSWORD":  "testpass",
			},
			expectedConfig: &etcd.Config{
				Endpoints:      []string{"http://localhost:2379"},
				Username:       "testuser",
				Password:       "testpass",
				DialTimeout:    etcd.DefaultDialTimeout,
				SessionTTL:     etcd.DefaultSessionTTL,
				ElectionPrefix: etcd.DefaultElectionPrefix,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all etcd-related env vars.
			clearEtcdEnvVars(t)

			// Set test env vars.
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := etcd.NewConfigFromEnv()

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedConfig.Endpoints, cfg.Endpoints)
			require.Equal(t, tt.expectedConfig.CertFile, cfg.CertFile)
			require.Equal(t, tt.expectedConfig.KeyFile, cfg.KeyFile)
			require.Equal(t, tt.expectedConfig.CAFile, cfg.CAFile)
			require.Equal(t, tt.expectedConfig.Username, cfg.Username)
			require.Equal(t, tt.expectedConfig.Password, cfg.Password)
			require.Equal(t, tt.expectedConfig.DialTimeout, cfg.DialTimeout)
			require.Equal(t, tt.expectedConfig.SessionTTL, cfg.SessionTTL)
			require.Equal(t, tt.expectedConfig.ElectionPrefix, cfg.ElectionPrefix)
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      etcd.Config
		expectedErr string
	}{
		{
			name: "valid config with single endpoint",
			config: etcd.Config{
				Endpoints: []string{"http://localhost:2379"},
			},
		},
		{
			name: "valid config with multiple endpoints",
			config: etcd.Config{
				Endpoints: []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			},
		},
		{
			name: "valid config with TLS",
			config: etcd.Config{
				Endpoints: []string{"https://localhost:2379"},
				CertFile:  "/path/to/client.crt",
				KeyFile:   "/path/to/client.key",
			},
		},
		{
			name: "empty endpoints returns error",
			config: etcd.Config{
				Endpoints: []string{},
			},
			expectedErr: "at least one etcd endpoint is required",
		},
		{
			name: "cert without key returns error",
			config: etcd.Config{
				Endpoints: []string{"https://localhost:2379"},
				CertFile:  "/path/to/client.crt",
			},
			expectedErr: "both ETCD_CERT_FILE and ETCD_KEY_FILE must be provided together",
		},
		{
			name: "key without cert returns error",
			config: etcd.Config{
				Endpoints: []string{"https://localhost:2379"},
				KeyFile:   "/path/to/client.key",
			},
			expectedErr: "both ETCD_CERT_FILE and ETCD_KEY_FILE must be provided together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestConfigToClientConfig(t *testing.T) {
	tests := []struct {
		name              string
		config            etcd.Config
		expectError       bool
		expectTLS         bool
		expectCredentials bool
	}{
		{
			name: "basic config without TLS",
			config: etcd.Config{
				Endpoints:   []string{"http://localhost:2379"},
				DialTimeout: etcd.DefaultDialTimeout,
			},
			expectTLS: false,
		},
		{
			name: "config with credentials",
			config: etcd.Config{
				Endpoints:   []string{"http://localhost:2379"},
				DialTimeout: etcd.DefaultDialTimeout,
				Username:    "testuser",
				Password:    "testpass",
			},
			expectCredentials: true,
		},
		{
			name: "config with non-existent TLS files returns error",
			config: etcd.Config{
				Endpoints:   []string{"https://localhost:2379"},
				DialTimeout: etcd.DefaultDialTimeout,
				CertFile:    "/nonexistent/client.crt",
				KeyFile:     "/nonexistent/client.key",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCfg, err := tt.config.ToClientConfig()

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.config.Endpoints, clientCfg.Endpoints)
			require.Equal(t, tt.config.DialTimeout, clientCfg.DialTimeout)

			if tt.expectTLS {
				require.NotNil(t, clientCfg.TLS)
			} else {
				require.Nil(t, clientCfg.TLS)
			}

			if tt.expectCredentials {
				require.Equal(t, tt.config.Username, clientCfg.Username)
				require.Equal(t, tt.config.Password, clientCfg.Password)
			}
		})
	}
}

func TestIsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name:     "not configured when ETCD_ENDPOINTS is not set",
			envVars:  map[string]string{},
			expected: false,
		},
		{
			name: "not configured when ETCD_ENDPOINTS is empty",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "",
			},
			expected: false,
		},
		{
			name: "configured when ETCD_ENDPOINTS is set",
			envVars: map[string]string{
				"ETCD_ENDPOINTS": "http://localhost:2379",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEtcdEnvVars(t)

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			result := etcd.IsConfigured()
			require.Equal(t, tt.expected, result)
		})
	}
}

// clearEtcdEnvVars clears all etcd-related environment variables for test isolation.
func clearEtcdEnvVars(t *testing.T) {
	t.Helper()
	envVars := []string{
		"ETCD_ENDPOINTS",
		"ETCD_CERT_FILE",
		"ETCD_KEY_FILE",
		"ETCD_CA_FILE",
		"ETCD_USERNAME",
		"ETCD_PASSWORD",
	}
	for _, v := range envVars {
		_ = os.Unsetenv(v)
	}
}
