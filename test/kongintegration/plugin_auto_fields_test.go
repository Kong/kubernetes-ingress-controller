package kongintegration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/kubectl/pkg/cmd/util"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

func TestKongClientGoldenTestsOutputs_AutoFields_OSS(t *testing.T) {
	t.Parallel()

	runPluginAutoFieldCases(t, []pluginAutoFieldCase{
		{
			name:          "oauth2 provision_key",
			fixture:       "oauth2.yaml",
			forbiddenNull: []string{"provision_key"},
		},
		{
			name:          "jwt credential key secret",
			forbiddenNull: []string{"key", "secret"},
			buildContent: func(_ context.Context, _ *testing.T, _ logr.Logger, _ *adminapi.Client) *file.Content {
				return &file.Content{
					FormatVersion: "3.0",
					Plugins: []file.FPlugin{
						{
							Plugin: kong.Plugin{
								Name: kong.String("jwt"),
							},
						},
					},
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("consumer-jwt"),
							},
							JWTAuths: []*kong.JWTAuth{
								{
									Algorithm: kong.String("HS256"),
								},
							},
						},
					},
				}
			},
		},
		{
			name:          "key-auth credential key",
			forbiddenNull: []string{"key"},
			buildContent: func(_ context.Context, _ *testing.T, _ logr.Logger, _ *adminapi.Client) *file.Content {
				return &file.Content{
					FormatVersion: "3.0",
					Plugins: []file.FPlugin{
						{
							Plugin: kong.Plugin{
								Name: kong.String("key-auth"),
							},
						},
					},
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("consumer-key-auth"),
							},
							KeyAuths: []*kong.KeyAuth{{}},
						},
					},
				}
			},
		},
		{
			name:          "hmac-auth credential secret",
			fixture:       "hmac-auth.yaml",
			forbiddenNull: []string{"secret"},
		},
	})
}

func TestKongClientGoldenTestsOutputs_AutoFields_EE(t *testing.T) {
	if !testenv.KongEnterpriseEnabled() {
		t.Skip("enterprise auto-field coverage requires TEST_KONG_ENTERPRISE=true")
	}
	if testenv.KongLicenseData() == "" {
		t.Skip("enterprise auto-field coverage requires KONG_LICENSE_DATA")
	}

	t.Parallel()

	runPluginAutoFieldCases(t, []pluginAutoFieldCase{
		{
			name:          "rate-limiting-advanced namespace",
			fixture:       "rate-limiting-advanced-ee.yaml",
			forbiddenNull: []string{"namespace"},
		},
		{
			name:          "graphql-rate-limiting-advanced namespace",
			fixture:       "graphql-rate-limiting-advanced-ee.yaml",
			forbiddenNull: []string{"namespace"},
		},
		{
			name:          "kafka-upstream cluster_name",
			fixture:       "kafka-upstream-ee.yaml",
			forbiddenNull: []string{"cluster_name"},
		},
		{
			name:          "kafka-log cluster_name",
			fixture:       "kafka-log-ee.yaml",
			forbiddenNull: []string{"cluster_name"},
		},
		{
			name:          "openid-connect cache_tokens_salt",
			fixture:       "openid-connect-ee.yaml",
			forbiddenNull: []string{"cache_tokens_salt"},
		},
	})
}

type pluginAutoFieldCase struct {
	name          string
	fixture       string
	forbiddenNull []string
	buildContent  func(context.Context, *testing.T, logr.Logger, *adminapi.Client) *file.Content
}

type unavailableSchemaServiceProvider struct{}

func (unavailableSchemaServiceProvider) GetSchemaService() kong.AbstractSchemaService {
	return translator.UnavailableSchemaService{}
}

func runPluginAutoFieldCases(t *testing.T, testCases []pluginAutoFieldCase) {
	t.Helper()

	const (
		timeout = 5 * time.Second
		period  = 200 * time.Millisecond
	)

	ctx := t.Context()
	kongC := containers.NewKong(ctx, t)
	kongAdminClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), managercfg.AdminAPIClientConfig{}, "")
	require.NoError(t, err)
	adminClient := adminapi.NewClient(kongAdminClient)

	logbase, err := zap.NewDevelopment()
	require.NoError(t, err)
	logger := zapr.NewLogger(logbase)
	updateStrategy := sendconfig.NewUpdateStrategyInMemory(
		kongAdminClient,
		sendconfig.DefaultContentToDBLessConfigConverter{},
		logger,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := buildContentForCase(ctx, t, logger, adminClient, tc)
			assertNoAutoFields(t, content, tc.forbiddenNull)

			require.EventuallyWithT(t, func(t *assert.CollectT) {
				configSize, err := updateStrategy.Update(ctx, sendconfig.ContentWithHash{Content: content})
				if !assert.NoError(t, err) {
					return
				}
				assert.NotZero(t, configSize)
			}, timeout, period)
		})
	}
}

func buildContentForCase(
	ctx context.Context,
	t *testing.T,
	logger logr.Logger,
	adminClient *adminapi.Client,
	tc pluginAutoFieldCase,
) *file.Content {
	t.Helper()

	if tc.buildContent != nil {
		return tc.buildContent(ctx, t, logger, adminClient)
	}

	return buildContentFromFixture(ctx, t, logger, adminClient, filepath.Join("testdata", "plugin-auto-fields", tc.fixture))
}

func buildContentFromFixture(
	ctx context.Context,
	t *testing.T,
	logger logr.Logger,
	adminClient *adminapi.Client,
	fixturePath string,
) *file.Content {
	t.Helper()

	objects := extractObjectsFromYAML(t, fixturePath)
	cacheStores, err := store.NewCacheStoresFromObjYAML(objects...)
	require.NoError(t, err)

	translatorInstance, err := translator.NewTranslator(
		logger,
		store.New(cacheStores, "kong", logger),
		"",
		translator.FeatureFlags{},
		unavailableSchemaServiceProvider{},
		translator.Config{},
	)
	require.NoError(t, err)

	result := translatorInstance.BuildKongConfig() //nolint:contextcheck
	require.Empty(t, result.TranslationFailures)

	return deckgen.ToDeckContent(
		ctx,
		logger,
		result.KongState,
		deckgen.GenerateDeckContentParams{
			PluginSchemas: adminClient.PluginSchemaStore(),
		},
	)
}

func assertNoAutoFields(t *testing.T, content any, fieldNames []string) {
	t.Helper()

	contentJSON, err := json.Marshal(content)
	require.NoError(t, err)

	for _, fieldName := range fieldNames {
		require.NotContains(t, string(contentJSON), fmt.Sprintf("\"%s\"", fieldName))
	}
}

func extractObjectsFromYAML(t *testing.T, filePath string) [][]byte {
	t.Helper()

	yamlBytes, err := os.ReadFile(filePath)
	require.NoErrorf(t, err, "failed reading input file: %s", filePath)

	stripped := util.ManualStrip(yamlBytes)
	objects := bytes.Split(stripped, []byte("---"))

	filtered := make([][]byte, 0, len(objects))
	for _, object := range objects {
		if len(strings.TrimSpace(string(object))) == 0 {
			continue
		}
		filtered = append(filtered, object)
	}

	return filtered
}
