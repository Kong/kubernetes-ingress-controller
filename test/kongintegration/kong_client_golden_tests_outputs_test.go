package kongintegration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

const (
	timeout = 5 * time.Second
	tick    = 100 * time.Millisecond
)

// TestKongClientGoldenTestsOutputs ensures that the KongClient's golden tests outputs are accepted by Kong.
func TestKongClientGoldenTestsOutputs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// By default, run only non-EE tests.
	goldenTestsOutputsPaths := lo.Filter(allGoldenTestsOutputsPaths(t), func(path string, _ int) bool {
		return !strings.Contains(path, "-ee/") // Skip Enterprise tests.
	})
	// If the Kong Enterprise is enabled, run all tests.
	if testenv.KongEnterpriseEnabled() {
		if testenv.KongLicenseData() == "" {
			t.Skip("Kong Enterprise enabled, but no license data provided")
		}
		goldenTestsOutputsPaths = allGoldenTestsOutputsPaths(t)
	}

	expressionRoutesOutputsPaths := lo.Filter(goldenTestsOutputsPaths, func(path string, _ int) bool {
		return strings.Contains(path, "expression-routes-on_")
	})
	defaultOutputsPaths := lo.Filter(goldenTestsOutputsPaths, func(path string, _ int) bool {
		return strings.Contains(path, "default_")
	})

	t.Logf("will test %d expression routes outputs and %d default ones", len(goldenTestsOutputsPaths), len(defaultOutputsPaths))

	t.Run("expressions router", func(t *testing.T) {
		t.Parallel()

		kongC := containers.NewKong(ctx, t, containers.KongWithRouterFlavor("expressions"))

		kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), helpers.DefaultHTTPClient())
		require.NoError(t, err)

		for _, goldenTestOutputPath := range expressionRoutesOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				ensureGoldenTestOutputIsAccepted(ctx, t, goldenTestOutputPath, kongClient)
			})
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		kongC := containers.NewKong(ctx, t, containers.KongWithRouterFlavor("traditional"))
		kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), helpers.DefaultHTTPClient())
		require.NoError(t, err)

		for _, goldenTestOutputPath := range defaultOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				ensureGoldenTestOutputIsAccepted(ctx, t, goldenTestOutputPath, kongClient)
			})
		}
	})
}

// TestKongClientGoldenTestsOutputs ensures that the KongClient's golden tests outputs are accepted by Konnect Control Plane
// Admin API.
func TestKongClientGoldenTestsOutputs_Konnect(t *testing.T) {
	konnect.SkipIfMissingRequiredKonnectEnvVariables(t)
	t.Parallel()

	ctx := context.Background()

	cpID := konnect.CreateTestControlPlane(ctx, t)
	cert, key := konnect.CreateClientCertificate(ctx, t, cpID)
	adminAPIClient := konnect.CreateKonnectAdminAPIClient(t, cpID, cert, key)
	updateStrategy := sendconfig.NewUpdateStrategyDBModeKonnect(adminAPIClient.AdminAPIClient(), dump.Config{
		SkipCACerts:         true,
		KonnectControlPlane: cpID,
	}, semver.MustParse("3.5.0"), 10, nil, logr.Discard())

	for _, goldenTestOutputPath := range allGoldenTestsOutputsPaths(t) {
		t.Run(goldenTestOutputPath, func(t *testing.T) {
			goldenTestOutput, err := os.ReadFile(goldenTestOutputPath)
			require.NoError(t, err)

			content := &file.Content{}
			err = yaml.Unmarshal(goldenTestOutput, content)
			require.NoError(t, err)

			require.EventuallyWithT(t, func(t *assert.CollectT) {
				err := updateStrategy.Update(ctx, sendconfig.ContentWithHash{Content: content})
				assert.NoError(t, err)
			}, timeout, tick)
		})
	}
}

func ensureGoldenTestOutputIsAccepted(ctx context.Context, t *testing.T, goldenTestOutputPath string, kongClient *kong.Client) {
	goldenTestOutput, err := os.ReadFile(goldenTestOutputPath)
	require.NoError(t, err)

	cfg := map[string]any{}
	err = yaml.Unmarshal(goldenTestOutput, &cfg)
	require.NoError(t, err)

	cfgAsJSON, err := json.Marshal(cfg)
	require.NoError(t, err)

	require.EventuallyWithT(t, func(t *assert.CollectT) {
		err := kongClient.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(cfgAsJSON), true, true)
		if !assert.NoErrorf(t, err, "failed to reload declarative config") {
			apiErr := &kong.APIError{}
			if errors.As(err, &apiErr) {
				t.Errorf("Kong Admin API response: %s", apiErr.Raw())
			}
		}
	}, timeout, tick)
}

func allGoldenTestsOutputsPaths(t *testing.T) []string {
	const goldenTestsOutputsGlob = "../../internal/dataplane/testdata/golden/*/*_golden.yaml"
	goldenTestsOutputsPaths, err := filepath.Glob(goldenTestsOutputsGlob)
	require.NoError(t, err)
	require.NotEmpty(t, goldenTestsOutputsPaths, "no golden tests outputs found")
	return goldenTestsOutputsPaths
}
