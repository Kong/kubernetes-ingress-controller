package manager_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestConfigValidatedVars(t *testing.T) {
	type testCase struct {
		Input                 string
		ExpectedValue         any
		ExtractValueFn        func(c manager.Config) any
		ExpectedErrorContains string
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
				ExpectedValue: string(gateway.ControllerName),
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
				ExpectedValue: types.NamespacedName{Namespace: "namespace", Name: "servicename"},
			},
			{
				Input:                 "servicename",
				ExpectedErrorContains: "the expected format is namespace/name",
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

				err := c.FlagSet().Parse(input)
				if tc.ExpectedErrorContains != "" {
					require.ErrorContains(t, err, tc.ExpectedErrorContains)
					return
				}

				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, tc.ExtractValueFn(c))
			})
		}
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("konnect", func(t *testing.T) {
		validEnabled := func() *manager.Config {
			return &manager.Config{
				Konnect: adminapi.KonnectConfig{
					ConfigSynchronizationEnabled: true,
					RuntimeGroup:                 "fbd3036f-0f1c-4e98-b71c-d4cd61213f90",
					Address:                      "https://us.kic.api.konghq.tech",
					ClientTLS: adminapi.TLSClientConfig{
						// We do not set valid cert or key, and it's still considered valid as at this level we only care
						// about them being not empty. Their validity is to be verified later on by the Admin API client
						// constructor.
						Cert: "not-empty-cert",
						Key:  "not-empty-key",
					},
				},
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
			c.Konnect.ClientTLS = adminapi.TLSClientConfig{}
			require.ErrorContains(t, c.Validate(), "TLS client config invalid")
		})

		t.Run("enabled with no tls cert is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.ClientTLS.Cert = ""
			require.ErrorContains(t, c.Validate(), "missing client cert")
		})

		t.Run("enabled with no tls key is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.ClientTLS.Key = ""
			require.ErrorContains(t, c.Validate(), "missing client key")
		})

		t.Run("enabled with tls cert file instead of cert is accepted", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.ClientTLS.Cert = ""
			c.Konnect.ClientTLS.CertFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})

		t.Run("enabled with tls key file instead of key is accepted", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.ClientTLS.Key = ""
			c.Konnect.ClientTLS.KeyFile = "non-empty-path"
			require.NoError(t, c.Validate())
		})

		t.Run("enabled with no runtime group is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.RuntimeGroup = ""
			require.ErrorContains(t, c.Validate(), "runtime group not specified")
		})

		t.Run("enabled with no address is rejected", func(t *testing.T) {
			c := validEnabled()
			c.Konnect.Address = ""
			require.ErrorContains(t, c.Validate(), "address not specified")
		})
	})
}
