package config_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/cmd/rootcmd/config"
	mgrconfig "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

func TestCLIArgumentAndFeatureGatesParser(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedErrMsg string
		expectedGates  map[string]bool
		expectedAddr   string
	}{
		{
			name:         "default values",
			args:         []string{},
			expectedAddr: ":10254",
		},
		{
			name: "enable GatewayAlpha and RewriteURIs",
			args: []string{
				"--feature-gates=GatewayAlpha=true,RewriteURIs=true",
				"--health-probe-bind-address=:4321",
			},
			expectedGates: map[string]bool{
				mgrconfig.GatewayAlphaFeature: true,
				mgrconfig.RewriteURIsFeature:  true,
			},
			expectedAddr: ":4321",
		},
		{
			name: "disable GatewayAlpha and enable RewriteURIs",
			args: []string{
				"--feature-gates=GatewayAlpha=false,RewriteURIs=true",
				"--health-probe-bind-address=:1234",
			},
			expectedGates: map[string]bool{
				mgrconfig.GatewayAlphaFeature: false,
				mgrconfig.RewriteURIsFeature:  true,
			},
			expectedAddr: ":1234",
		},
		{
			name: "for non-existing feature descriptive error message is returned",
			args: []string{
				"--feature-gates=GatewayAlpha=true,RewriteURIs=false,NonExistingGate=true",
				"--health-probe-bind-address=:5678",
			},
			expectedErrMsg: `invalid argument "GatewayAlpha=true,RewriteURIs=false,NonExistingGate=true" for "--feature-gates" flag: NonExistingGate is not a valid feature, please see the documentation: https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md`,
		},
		{
			name: "for non-existing option descriptive error message is returned",
			args: []string{
				"--non-existing-option=1234",
			},
			expectedErrMsg: `unknown flag: --non-existing-option`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := config.NewCLIConfig()
			flagSet := c.FlagSet()
			flagSet.SetOutput(io.Discard)

			err := flagSet.Parse(tt.args)
			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				return
			}
			require.NoError(t, err)

			for key, expectedValue := range tt.expectedGates {
				require.Equal(t, expectedValue, c.FeatureGates.Enabled(key), "feature gate not set according to configuration passed")
			}
			for key, value := range mgrconfig.GetFeatureGatesDefaults() {
				if ok := tt.expectedGates[key]; ok {
					continue
				}
				require.Equal(t, value, c.FeatureGates.Enabled(key), "not configured feature gate does not have default value")
			}

			require.Equal(t, tt.expectedAddr, c.ProbeAddr, "health probe address not set according to configuration passed")
		})
	}
}
