package manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

func TestConfigValidatedVars(t *testing.T) {
	type testCase struct {
		Input                 string
		ExpectedValue         any
		ExtractValueFn        func(c Config) any
		ExpectedErrorContains string
	}

	testCases := map[string][]testCase{
		"--gateway-api-controller-name": {
			{
				Input: "example.com/controller-name",
				ExtractValueFn: func(c Config) any {
					return c.GatewayAPIControllerName
				},
				ExpectedValue: "example.com/controller-name",
			},
			{
				Input: "",
				ExtractValueFn: func(c Config) any {
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
				ExtractValueFn: func(c Config) any {
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

	for flag, flagTestCases := range testCases {
		for _, tc := range flagTestCases {
			t.Run(fmt.Sprintf("%s=%s", flag, tc.Input), func(t *testing.T) {
				var c Config
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
