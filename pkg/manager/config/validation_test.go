package config_test

import (
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

func TestConfigValidate(t *testing.T) {
	t.Run("konnect", func(t *testing.T) {
		validEnabled := func() *managercfg.Config {
			return &managercfg.Config{
				KongAdminSvc: mo.Some(k8stypes.NamespacedName{Name: "admin-svc", Namespace: "ns"}),
				Konnect: managercfg.KonnectConfig{
					ConfigSynchronizationEnabled: true,
					ControlPlaneID:               "fbd3036f-0f1c-4e98-b71c-d4cd61213f90",
					Address:                      "https://us.kic.api.konghq.tech",
					ConfigSyncConcurrency:        managercfg.DefaultKonnectConfigSyncConcurrency,
					TLSClient: managercfg.TLSClientConfig{
						// We do not set valid cert or key, and it's still considered valid as at this level we only care
						// about them being not empty. Their validity is to be verified later on by the Admin API client
						// constructor.
						Cert: "not-empty-cert",
						Key:  "not-empty-key",
					},
					UploadConfigPeriod: managercfg.DefaultKonnectConfigUploadPeriod,
				},
				GatewayDiscoveryReadinessCheckInterval: managercfg.DefaultDataPlanesReadinessReconciliationInterval,
			}
		}

		t.Run("disabled should not require other vars to be set", func(t *testing.T) {
			c := &managercfg.Config{Konnect: managercfg.KonnectConfig{ConfigSynchronizationEnabled: false}}
			require.NoError(t, c.Validate())
		})

		t.Run("enabled should require tls client config", func(t *testing.T) {
			require.NoError(t, validEnabled().Validate())
		})

		t.Run("enabled with no tls config is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient = managercfg.TLSClientConfig{}
			require.ErrorContains(t, c.Validate(), "missing TLS client configuration")
		})

		t.Run("enabled with no tls key is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient.Key = ""
			require.ErrorContains(t, c.Validate(), "client certificate was provided, but the client key was not")
		})

		t.Run("enabled with no tls cert is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient.Cert = ""
			require.ErrorContains(t, c.Validate(), "client key was provided, but the client certificate was not")
		})

		t.Run("enabled with tls cert file instead of cert is accepted", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient.Cert = ""
			c.Konnect.TLSClient.CertFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})

		t.Run("enabled with tls key file instead of key is accepted", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient.Key = ""
			c.Konnect.TLSClient.KeyFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})

		t.Run("enabled with no control plane is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.ControlPlaneID = ""
			require.ErrorContains(t, c.Validate(), "control plane not specified")
		})

		t.Run("enabled with no address is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.Address = ""
			require.ErrorContains(t, c.Validate(), "address not specified")
		})

		t.Run("enabled with no gateway service discovery enabled", func(t *testing.T) {
			c := validEnabled()
			c.KongAdminSvc = managercfg.OptionalNamespacedName{}
			require.ErrorContains(t, c.Validate(), "--kong-admin-svc has to be set when using --konnect-sync-enabled")
		})

		t.Run("enabled with too small upload config period is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.UploadConfigPeriod = time.Second
			require.ErrorContains(t, c.Validate(), "cannot set upload config period to be smaller than 10s")
		})
	})

	t.Run("Admin API", func(t *testing.T) {
		validWithClientTLS := func() managercfg.Config {
			return managercfg.Config{
				KongAdminAPIConfig: managercfg.AdminAPIClientConfig{
					TLSClient: managercfg.TLSClientConfig{
						// We do not set valid cert or key, and it's still considered valid as at this level we only care
						// about them being not empty. Their validity is to be verified later on by the Admin API client
						// constructor.
						Cert: "not-empty-cert",
						Key:  "not-empty-key",
					},
				},
			}
		}

		t.Run("no TLS client is allowed", func(t *testing.T) {
			c := managercfg.Config{
				KongAdminAPIConfig: managercfg.AdminAPIClientConfig{
					TLSClient: managercfg.TLSClientConfig{},
				},
			}
			require.NoError(t, c.Validate())
		})

		t.Run("valid TLS client is allowed", func(t *testing.T) {
			c := validWithClientTLS()
			require.NoError(t, c.Validate())
		})

		t.Run("missing tls key is rejected", func(t *testing.T) {
			c := validWithClientTLS()
			c.KongAdminAPIConfig.TLSClient.Key = ""
			require.ErrorContains(t, c.Validate(), "client certificate was provided, but the client key was not")
		})

		t.Run("missing tls cert is rejected", func(t *testing.T) {
			c := validWithClientTLS()
			c.KongAdminAPIConfig.TLSClient.Cert = ""
			require.ErrorContains(t, c.Validate(), "client key was provided, but the client certificate was not")
		})

		t.Run("tls cert file instead of cert is accepted", func(t *testing.T) {
			c := validWithClientTLS()
			c.KongAdminAPIConfig.TLSClient.Cert = ""
			c.KongAdminAPIConfig.TLSClient.CertFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})

		t.Run("tls key file instead of key is accepted", func(t *testing.T) {
			c := validWithClientTLS()
			c.KongAdminAPIConfig.TLSClient.Key = ""
			c.KongAdminAPIConfig.TLSClient.KeyFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})
	})

	t.Run("Admin Token", func(t *testing.T) {
		validWithToken := func() managercfg.Config {
			return managercfg.Config{
				KongAdminToken: "non-empty-token",
			}
		}

		t.Run("admin token accepted", func(t *testing.T) {
			c := validWithToken()
			require.NoError(t, c.Validate())
		})
	})

	t.Run("Admin Token Path", func(t *testing.T) {
		validWithTokenPath := func() managercfg.Config {
			return managercfg.Config{
				KongAdminTokenPath: "non-empty-token-path",
			}
		}

		t.Run("admin token and token path rejected", func(t *testing.T) {
			c := validWithTokenPath()
			c.KongAdminToken = "non-empty-token"
			require.ErrorContains(t, c.Validate(), "both admin token and admin token file specified, only one allowed")
		})
	})

	t.Run("use last valid config for fallback", func(t *testing.T) {
		t.Run("enabled without feature gate is rejected", func(t *testing.T) {
			c := managercfg.Config{
				UseLastValidConfigForFallback: true,
			}
			require.ErrorContains(t, c.Validate(), "--use-last-valid-config-for-fallback or CONTROLLER_USE_LAST_VALID_CONFIG_FOR_FALLBACK can only be used with FallbackConfiguration feature gate enabled")
		})
		t.Run("enabled with feature gate is accepted", func(t *testing.T) {
			c := managercfg.Config{
				UseLastValidConfigForFallback: true,
				FeatureGates: map[string]bool{
					managercfg.FallbackConfigurationFeature: true,
				},
			}
			require.NoError(t, c.Validate())
		})
	})

	t.Run("gateway discovery", func(t *testing.T) {
		validEnabled := func() *managercfg.Config {
			return &managercfg.Config{
				KongAdminSvc:                           mo.Some(k8stypes.NamespacedName{Name: "admin-svc", Namespace: "ns"}),
				GatewayDiscoveryReadinessCheckInterval: managercfg.DefaultDataPlanesReadinessReconciliationInterval,
				GatewayDiscoveryReadinessCheckTimeout:  managercfg.DefaultDataPlanesReadinessCheckTimeout,
			}
		}

		t.Run("disabled should not check other fields to set", func(t *testing.T) {
			c := &managercfg.Config{}
			require.NoError(t, c.Validate())
		})

		t.Run("enabled with valid configuration should pass", func(t *testing.T) {
			c := validEnabled()
			require.NoError(t, c.Validate())
		})

		t.Run("too small reconciliation interval should not pass", func(t *testing.T) {
			c := validEnabled()
			c.GatewayDiscoveryReadinessCheckInterval = 2 * time.Second
			c.GatewayDiscoveryReadinessCheckTimeout = time.Second
			require.ErrorContains(t, c.Validate(), "Readiness check reconciliation interval cannot be less than 3s")
		})

		t.Run("readiness check timeout must be less than reconciliation interval", func(t *testing.T) {
			c := validEnabled()
			c.GatewayDiscoveryReadinessCheckTimeout = managercfg.DefaultDataPlanesReadinessReconciliationInterval
			require.ErrorContains(t, c.Validate(), "Readiness check timeout must be less than readiness check recociliation interval")
		})
	})
}
