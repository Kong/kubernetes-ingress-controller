package kongintegration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/kongintegration/containers"
)

// TestGoldenTestsOutputs ensures that the Parser's golden tests outputs are accepted by Kong.
func TestParsersGoldenTestsOutputs(t *testing.T) {
	t.Parallel()

	const (
		timeout = 5 * time.Second
		tick    = time.Millisecond * 100
	)

	ctx := context.Background()

	const goldenTestsOutputsGlob = "../../internal/dataplane/parser/testdata/golden/*/*_golden.yaml"
	goldenTestsOutputsPaths, err := filepath.Glob(goldenTestsOutputsGlob)
	require.NoError(t, err)
	require.NotEmpty(t, goldenTestsOutputsPaths, "no golden tests outputs found")

	// TODO: Test EE features as well (requires kong/kong-gateway + license).
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4815
	goldenTestsOutputsPaths = lo.Filter(goldenTestsOutputsPaths, func(path string, _ int) bool {
		return !strings.Contains(path, "-ee/") // Skip Enterprise tests.
	})

	expressionRoutesOutputsPaths := lo.Filter(goldenTestsOutputsPaths, func(path string, _ int) bool {
		return strings.Contains(path, "expression-routes-on_")
	})
	defaultOutputsPaths := lo.Filter(goldenTestsOutputsPaths, func(path string, _ int) bool {
		return strings.Contains(path, "default_")
	})

	t.Logf("will test %d expression routes outputs and %d default ones", len(goldenTestsOutputsPaths), len(defaultOutputsPaths))

	runTest := func(t *testing.T, goldenTestOutputPath string, sut sendconfig.UpdateStrategyInMemory) {
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

	t.Run("expressions router", func(t *testing.T) {
		t.Parallel()

		kongC := containers.NewKong(ctx, t, containers.KongWithRouterFlavor("expressions"))
		kongClient, err := kong.NewClient(kong.String(kongC.AdminURL(ctx, t)), helpers.DefaultHTTPClient())
		require.NoError(t, err)

		sut := sendconfig.NewUpdateStrategyInMemory(
			kongClient,
			sendconfig.DefaultContentToDBLessConfigConverter{},
			logr.Discard(),
		)

		for _, goldenTestOutputPath := range expressionRoutesOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				runTest(t, goldenTestOutputPath, sut)
			})
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		kongC := containers.NewKong(ctx, t, containers.KongWithRouterFlavor("traditional"))
		kongClient, err := kong.NewClient(kong.String(kongC.AdminURL(ctx, t)), helpers.DefaultHTTPClient())
		require.NoError(t, err)

		sut := sendconfig.NewUpdateStrategyInMemory(
			kongClient,
			sendconfig.DefaultContentToDBLessConfigConverter{},
			logr.Discard(),
		)

		for _, goldenTestOutputPath := range defaultOutputsPaths {
			t.Run(goldenTestOutputPath, func(t *testing.T) {
				runTest(t, goldenTestOutputPath, sut)
			})
		}
	})
}
