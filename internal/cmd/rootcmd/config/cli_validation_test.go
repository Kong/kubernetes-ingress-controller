package config_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/cmd/rootcmd/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	mgrconfig "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

func TestConfigValidatedVars(t *testing.T) {
	type testCase struct {
		Input                      string
		ExpectedValue              any
		ExtractValueFn             func(c config.CLIConfig) any
		ExpectedErrorContains      string
		ExpectedUsageAdditionalMsg string
	}

	testCasesGroupedByFlag := map[string][]testCase{
		"--gateway-api-controller-name": {
			{
				Input: "example.com/controller-name",
				ExtractValueFn: func(c config.CLIConfig) any {
					return c.GatewayAPIControllerName
				},
				ExpectedValue: "example.com/controller-name",
			},
			{
				Input: "",
				ExtractValueFn: func(c config.CLIConfig) any {
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
				ExtractValueFn: func(c config.CLIConfig) any {
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
				ExtractValueFn: func(c config.CLIConfig) any {
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
				ExtractValueFn: func(c config.CLIConfig) any {
					return c.Konnect.ControlPlaneID
				},
				ExpectedValue:              "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
				ExpectedUsageAdditionalMsg: "Flag --konnect-runtime-group-id has been deprecated, Use --konnect-control-plane-id instead.\n",
			},
		},
		"--konnect-control-plane-id": {
			{
				Input: "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
				ExtractValueFn: func(c config.CLIConfig) any {
					return c.Konnect.ControlPlaneID
				},
				ExpectedValue: "5ef731c0-6081-49d6-b3ec-d4f85e58b956",
			},
		},
		"--gateway-to-reconcile": {
			{
				Input: "namespace/gatewayname",
				ExtractValueFn: func(c config.CLIConfig) any {
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
		"--secret-label-selector": {
			{
				Input: "konghq.com/label-for-caching",
				ExtractValueFn: func(c config.CLIConfig) any {
					return c.SecretLabelSelector
				},
				ExpectedValue: "konghq.com/label-for-caching",
			},
		},
		"--configmap-label-selector": {
			{
				Input: "konghq.com/label-for-caching",
				ExtractValueFn: func(c config.CLIConfig) any {
					return c.ConfigMapLabelSelector
				},
				ExpectedValue: "konghq.com/label-for-caching",
			},
		},
	}

	for flag, flagTestCases := range testCasesGroupedByFlag {
		for _, tc := range flagTestCases {
			t.Run(fmt.Sprintf("%s=%s", flag, tc.Input), func(t *testing.T) {
				c := config.NewCLIConfig()
				flagSet := c.FlagSet()

				var input []string
				if tc.Input != "" {
					input = []string{flag, tc.Input}
				}

				var usageAdditionalMsg bytes.Buffer
				flagSet.SetOutput(&usageAdditionalMsg)

				err := flagSet.Parse(input)
				if tc.ExpectedErrorContains != "" {
					require.ErrorContains(t, err, tc.ExpectedErrorContains)
					return
				}

				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, tc.ExtractValueFn(*c))
				require.Equal(t, tc.ExpectedUsageAdditionalMsg, usageAdditionalMsg.String())
			})
		}
	}
}

func TesCLIArgumentParser(t *testing.T) {
	c := config.NewCLIConfig()
	flagSet := c.FlagSet()
	flagSet.SetOutput(io.Discard)

	err := flagSet.Parse([]string{
		"--feature-gates=GatewayAlpha=true,RewriteURIs=true",
		"--health-probe-bind-address=:4321",
	})
	require.NoError(t, err)
	require.Equal(t, true, c.FeatureGates.Enabled(mgrconfig.GatewayAlphaFeature))
	require.True(t, c.FeatureGates.Enabled(mgrconfig.RewriteURIsFeature))
	for key, value := range mgrconfig.GetFeatureGatesDefaults() {
		if key == mgrconfig.GatewayAlphaFeature || key == mgrconfig.RewriteURIsFeature {
			continue
		}
		require.Equal(t, value, c.FeatureGates.Enabled(key))
	}
	require.Equal(t, c.ProbeAddr, ":4321")
}
