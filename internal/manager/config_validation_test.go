package manager_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
)

func TestConfigValidatedVars(t *testing.T) {
	type testCase struct {
		Input                      string
		ExpectedValue              any
		ExtractValueFn             func(c manager.Config) any
		ExpectedErrorContains      string
		ExpectedUsageAdditionalMsg string
	}

	testCasesGroupedByFlag := map[string][]testCase{
		"--gateway-api-controller-name": {
			{
				Input: "example.com/controller-name",
				ExtractValueFn: func(c manager.Config) any {
					return c.GatewayAPIControllerName
				},
				ExpectedValue: "example.com/controller-name",
			},
			{
				Input: "",
				ExtractValueFn: func(c manager.Config) any {
					return c.GatewayAPIControllerName
				},
				ExpectedValue: string(gateway.GetControllerName()),
			},
			{
				Input:                 "%invalid_controller_name$",
				ExpectedErrorContains: "the expected format is example.com/controller-name",
			},
		},
		"--publish-service": {
			{
				Input: "namespace/servicename",
				ExtractValueFn: func(c manager.Config) any {
					return c.PublishService
				},
				ExpectedValue: mo.Some(k8stypes.NamespacedName{Namespace: "namespace", Name: "servicename"}),
			},
			{
				Input:                 "servicename",
				ExpectedErrorContains: "the expected format is namespace/name",
			},
		},
		"--kong-admin-svc": {
			{
				Input: "namespace/servicename",
				ExtractValueFn: func(c manager.Config) any {
					return c.KongAdminSvc
				},
				ExpectedValue: mo.Some(k8stypes.NamespacedName{Namespace: "namespace", Name: "servicename"}),
			},
			{
				Input:                 "namespace/",
				ExpectedErrorContains: "name cannot be empty",
			},
			{
				Input:                 "/name",
				ExpectedErrorContains: "namespace cannot be empty",
			},
		},
		"--konnect-runtime-group-id": {
			{
				Input: "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
				ExtractValueFn: func(c manager.Config) any {
					return c.Konnect.ControlPlaneID
				},
				ExpectedValue:              "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
				ExpectedUsageAdditionalMsg: "Flag --konnect-runtime-group-id has been deprecated, Use --konnect-control-plane-id instead.\n",
			},
		},
		"--konnect-control-plane-id": {
			{
				Input: "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
				ExtractValueFn: func(c manager.Config) any {
					return c.Konnect.ControlPlaneID
				},
				ExpectedValue: "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
			},
		},
		"--gateway-to-reconcile": {
			{
				Input: "namespace/gatewayname",
				ExtractValueFn: func(c manager.Config) any {
					return c.GatewayToReconcile
				},
				ExpectedValue: mo.Some(k8stypes.NamespacedName{Namespace: "namespace", Name: "gatewayname"}),
			},
			{
				Input:                 "namespace/",
				ExpectedErrorContains: "name cannot be empty",
			},
			{
				Input:                 "/name",
				ExpectedErrorContains: "namespace cannot be empty",
			},
		},
	}

	for flag, flagTestCases := range testCasesGroupedByFlag {
		for _, tc := range flagTestCases {
			t.Run(fmt.Sprintf("%s=%s", flag, tc.Input), func(t *testing.T) {
				var c manager.Config
				var input []string
				if tc.Input != "" {
					input = []string{flag, tc.Input}
				}

				flagSet := c.FlagSet()
				var usageAdditionalMsg bytes.Buffer
				flagSet.SetOutput(&usageAdditionalMsg)

				err := flagSet.Parse(input)
				if tc.ExpectedErrorContains != "" {
					require.ErrorContains(t, err, tc.ExpectedErrorContains)
					return
				}

				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, tc.ExtractValueFn(c))
				require.Equal(t, tc.ExpectedUsageAdditionalMsg, usageAdditionalMsg.String())
			})
		}
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("konnect", func(t *testing.T) {
		validEnabled := func() *manager.Config {
			return &manager.Config{
				KongAdminSvc: mo.Some(k8stypes.NamespacedName{Name: "admin-svc", Namespace: "ns"}),
				Konnect: adminapi.KonnectConfig{
					ConfigSynchronizationEnabled: true,
					ControlPlaneID:               "fbd3036f-0f1c-4e98-b71c-d4cd61213f90",
					Address:                      "https://us.kic.api.konghq.tech",
					TLSClient: adminapi.TLSClientConfig{
						// We do not set valid cert or key, and it's still considered valid as at this level we only care
						// about them being not empty. Their validity is to be verified later on by the Admin API client
						// constructor.
						Cert: "not-empty-cert",
						Key:  "not-empty-key",
					},
					UploadConfigPeriod: konnect.DefaultConfigUploadPeriod,
				},
				GatewayDiscoveryReadinessCheckInterval: clients.DefaultReadinessReconciliationInterval,
			}
		}

		t.Run("disabled should not require other vars to be set", func(t *testing.T) {
			c := &manager.Config{Konnect: adminapi.KonnectConfig{ConfigSynchronizationEnabled: false}}
			require.NoError(t, c.Validate())
		})

		t.Run("enabled should require tls client config", func(t *testing.T) {
			require.NoError(t, validEnabled().Validate())
		})

		t.Run("enabled with no tls config is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.TLSClient = adminapi.TLSClientConfig{}
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
			c.KongAdminSvc = manager.OptionalNamespacedName{}
			require.ErrorContains(t, c.Validate(), "--kong-admin-svc has to be set when using --konnect-sync-enabled")
		})

		t.Run("enabled with too small upload config period is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.UploadConfigPeriod = time.Second
			require.ErrorContains(t, c.Validate(), "cannot set upload config period to be smaller than 10s")
		})
	})

	t.Run("Admin API", func(t *testing.T) {
		validWithClientTLS := func() manager.Config {
			return manager.Config{
				KongAdminAPIConfig: adminapi.HTTPClientOpts{
					TLSClient: adminapi.TLSClientConfig{
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
			c := manager.Config{
				KongAdminAPIConfig: adminapi.HTTPClientOpts{
					TLSClient: adminapi.TLSClientConfig{},
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
		validWithToken := func() manager.Config {
			return manager.Config{
				KongAdminToken: "non-empty-token",
			}
		}

		t.Run("admin token accepted", func(t *testing.T) {
			c := validWithToken()
			require.NoError(t, c.Validate())
		})
	})

	t.Run("Admin Token Path", func(t *testing.T) {
		validWithTokenPath := func() manager.Config {
			return manager.Config{
				KongAdminTokenPath: "non-empty-token-path",
			}
		}

		t.Run("admin token and token path rejected", func(t *testing.T) {
			c := validWithTokenPath()
			c.KongAdminToken = "non-empty-token"
			require.ErrorContains(t, c.Validate(), "both admin token and admin token file specified, only one allowed")
		})
	})

	t.Run("--use-last-valid-config-for-fallback", func(t *testing.T) {
		t.Run("enabled without feature gate is rejected", func(t *testing.T) {
			c := manager.Config{
				UseLastValidConfigForFallback: true,
			}
			require.ErrorContains(t, c.Validate(), "--use-last-valid-config-for-fallback or CONTROLLER_USE_LAST_VALID_CONFIG_FOR_FALLBACK can only be used with FallbackConfiguration feature gate enabled")
		})
		t.Run("enabled with feature gate is accepted", func(t *testing.T) {
			c := manager.Config{
				UseLastValidConfigForFallback: true,
				FeatureGates: map[string]bool{
					featuregates.FallbackConfiguration: true,
				},
			}
			require.NoError(t, c.Validate())
		})
	})

	t.Run("gateway discovery", func(t *testing.T) {
		validEnabled := func() *manager.Config {
			return &manager.Config{
				KongAdminSvc:                           mo.Some(k8stypes.NamespacedName{Name: "admin-svc", Namespace: "ns"}),
				GatewayDiscoveryReadinessCheckInterval: clients.DefaultReadinessReconciliationInterval,
				GatewayDiscoveryReadinessCheckTimeout:  clients.DefaultReadinessCheckTimeout,
			}
		}

		t.Run("disabled should not check other fields to set", func(t *testing.T) {
			c := &manager.Config{}
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
			c.GatewayDiscoveryReadinessCheckTimeout = clients.DefaultReadinessReconciliationInterval
			require.ErrorContains(t, c.Validate(), "Readiness check timeout must be less than readiness check recociliation interval")
		})
	})
}
