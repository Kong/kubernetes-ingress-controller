package kongintegration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

const (
	timeout = 5 * time.Second
	tick    = 100 * time.Millisecond
)

// TestGoldenTestsOutputs ensures that the translators' golden tests outputs are accepted by Kong.
func TestTranslatorsGoldenTestsOutputs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// TODO: Test EE features as well (requires kong/kong-gateway + license).
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4815
	goldenTestsOutputsPaths := lo.Filter(allGoldenTestsOutputsPaths(t), func(path string, _ int) bool {
		return !strings.Contains(path, "-ee/") // Skip Enterprise tests.
	})

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

		sut := sendconfig.NewUpdateStrategyInMemory(
			kongClient,
			sendconfig.DefaultContentToDBLessConfigConverter{},
			logr.Discard(),
		)

		for _, goldenTestOutputPath := range expressionRoutesOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				ensureGoldenTestOutputIsAccepted(ctx, t, goldenTestOutputPath, sut)
			})
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		kongC := containers.NewKong(ctx, t, containers.KongWithRouterFlavor("traditional"))
		kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), helpers.DefaultHTTPClient())
		require.NoError(t, err)

		sut := sendconfig.NewUpdateStrategyInMemory(
			kongClient,
			sendconfig.DefaultContentToDBLessConfigConverter{},
			logr.Discard(),
		)

		for _, goldenTestOutputPath := range defaultOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				ensureGoldenTestOutputIsAccepted(ctx, t, goldenTestOutputPath, sut)
			})
		}
	})
}

// TestGoldenTestsOutputs ensures that the translators' golden tests outputs are accepted by Konnect Control Plane
// Admin API.
func TestTranslatorsGoldenTestsOutputs_Konnect(t *testing.T) {
	konnect.SkipIfMissingRequiredKonnectEnvVariables(t)
	t.Parallel()

	ctx := context.Background()

	cpID := konnect.CreateTestControlPlane(ctx, t)
	cert, key := konnect.CreateClientCertificate(ctx, t, cpID)
	adminAPIClient := konnect.CreateKonnectAdminAPIClient(t, cpID, cert, key)
	updateStrategy := sendconfig.NewUpdateStrategyDBModeKonnect(adminAPIClient.AdminAPIClient(), dump.Config{
		SkipCACerts:         true,
		KonnectControlPlane: cpID,
	}, semver.MustParse("3.5.0"), 10)

	for _, goldenTestOutputPath := range allGoldenTestsOutputsPaths(t) {
		t.Run(goldenTestOutputPath, func(t *testing.T) {
			ensureGoldenTestOutputIsAccepted(ctx, t, goldenTestOutputPath, updateStrategy)
		})
	}
}

func ensureGoldenTestOutputIsAccepted(
	ctx context.Context,
	t *testing.T,
	goldenTestOutputPath string,
	sut sendconfig.UpdateStrategy,
) {
	goldenTestOutput, err := os.ReadFile(goldenTestOutputPath)
	require.NoError(t, err)

	content := &file.Content{}
	err = yaml.Unmarshal(goldenTestOutput, content)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		err, resourceErrors, parseErr := sut.Update(ctx, sendconfig.ContentWithHash{Content: content})
		if err != nil {
			t.Logf("error: %v", err)
			return false
		}
		if len(resourceErrors) > 0 {
			t.Logf("resource errors: %v", resourceErrors)
			return false
		}
		if parseErr != nil {
			t.Logf("parse error: %v", parseErr)
			return false
		}
		return true
	}, timeout, tick)
}

func allGoldenTestsOutputsPaths(t *testing.T) []string {
	const goldenTestsOutputsGlob = "../../internal/dataplane/translator/testdata/golden/*/*_golden.yaml"
	goldenTestsOutputsPaths, err := filepath.Glob(goldenTestsOutputsGlob)
	require.NoError(t, err)
	require.NotEmpty(t, goldenTestsOutputsPaths, "no golden tests outputs found")
	return goldenTestsOutputsPaths
}
