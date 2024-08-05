package dataplane

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

var (
	// updateGolden tells whether to update golden files using the current config received by the Admin API.
	updateGolden = flag.Bool("update", false, "update golden files")

	// defaultFeatureFlags is the default set of Translator feature flags to use in tests. Can be overridden in a test case.
	defaultFeatureFlags = func() translator.FeatureFlags {
		defaults := featuregates.GetFeatureGatesDefaults()
		return translator.FeatureFlags{
			// We do not verify configuration reports in golden tests.
			ReportConfiguredKubernetesObjects: false,

			// Feature flags that are directly propagated from the feature gates get their defaults.
			FillIDs:           defaults.Enabled(featuregates.FillIDsFeature),
			KongServiceFacade: defaults.Enabled(featuregates.KongServiceFacade),
			KongCustomEntity:  defaults.Enabled(featuregates.KongCustomEntity),
		}
	}
)

const (
	goldenDir          = "testdata/golden"
	inFileName         = "in.yaml"
	goldenFileSuffix   = "_golden.yaml"
	settingsFileSuffix = "_settings.yaml"
)

// TestKongClient_GoldenTests runs the golden tests for the KongClient.
//
// Command to update the golden files:
// $ make test.golden.update
//
// Data for the test cases is stored in the "./testdata/golden" directory. Test cases are grouped into subdirectories
// based on the Kubernetes input that they run against so that each of the subdirectories has:
//   - an input file that represents the input with Kubernetes objects to be loaded into the store: "in.yaml",
//   - a set of "<settings-name>_settings.yaml" files that define settings for a given test case (i.e. translator feature flags),
//   - a set of expected golden "<settings-name>_golden.yaml" files (in declarative config format) where each file represents an
//     expected output for a given translator configuration defined in "<settings-name>_settings.yaml".
//
// The test case is executed by loading the in.yaml file into the store, then running KongClient.Update method with
// the store injected. KongClient pushes configuration to a mock Admin API HTTP server. We fetch the last received
// configuration from the server and compare the output with the expected golden file.
//
// When adding a new test case, you can follow these steps:
//  1. Add a new directory ./testdata/golden/<your-dir> with the "in.yaml" that you want to test against.
//  2. (Optional) Define a set of "<settings-name>_settings.yaml" files that define the translator configuration you want to
//     test. If you don't define any settings files, the test will run with default settings.
//  3. Run `make test.golden.update` to generate the golden files.
//  4. Inspect the generated golden files and make sure they're correct. If they are, commit them.
//
// If you introduce a change that may affect many test cases, and you're sure about it correctness, you can run the
// update command as well to update all golden files at once.
//
// If you want to make the mocked Admin API server return errors for specific objects, you can add an annotation
// "test.konghq.com/broken: true" to the object in the in.yaml file. If there's at least one object with this annotation,
// the test will expect an error from the KongClient.Update method and will turn the FallbackConfiguration feature on.
func TestKongClient_GoldenTests(t *testing.T) {
	// First, let's prepare the test cases basing on the testdata/golden directory contents.
	var testCases []kongClientGoldenTestCase
	testCasesDirectories, err := os.ReadDir(goldenDir)
	require.NoError(t, err, "failed to iterate over files in testdata/golden")

	for _, testCaseDir := range testCasesDirectories {
		testCaseDirPath := filepath.Join(goldenDir, testCaseDir.Name())
		require.True(t, testCaseDir.IsDir(),
			"%s is not a directory, while we expect testdata/golden/* to include only directories", testCaseDirPath)

		if *updateGolden {
			// If we're updating the golden files, let's first prune test case directories.
			pruneTestCaseDirectory(t, testCaseDirPath)
		}

		// Then, let's iterate over all settings files in the directory and add a test case for each of them.
		// If there are no settings files, we'll add just a single test case with default settings.
		for _, settings := range resolveSetsOfSettingsForTestCaseDir(t, testCaseDirPath) {
			testCases = append(testCases, kongClientGoldenTestCase{
				k8sConfigFile: filepath.Join(testCaseDirPath, inFileName),
				goldenFile:    filepath.Join(testCaseDirPath, fmt.Sprintf("%s%s", settings.name, goldenFileSuffix)),
				featureFlags:  settings.featureFlags,
			})
		}
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("in=%s,out=%s", tc.k8sConfigFile, tc.goldenFile), func(t *testing.T) {
			runKongClientGoldenTest(t, tc)
		})
	}
}

// pruneTestCaseDirectory removes all files from the test case directory except for in.yaml and *_settings.yaml files.
func pruneTestCaseDirectory(t *testing.T, path string) {
	filesInDirectory, err := os.ReadDir(path)
	require.NoError(t, err, "failed to iterate over files in test case directory %s", path)

	for _, fileInDirectory := range filesInDirectory {
		// First, let's skip the files we want to keep.
		if fileInDirectory.Name() == inFileName || strings.HasSuffix(fileInDirectory.Name(), settingsFileSuffix) {
			continue
		}

		// Then, let's remove any other files.
		err = os.Remove(filepath.Join(path, fileInDirectory.Name()))
		require.NoError(t, err, "failed to remove file %s", filepath.Join(path, fileInDirectory.Name()))
	}
}

// resolveSetsOfSettingsForTestCaseDir returns a slice of testCaseSettings, each of which represents a combination of
// feature flags and Kong version.
// The function iterates over a test case directory containing zero or more files named "<name>_settings.yaml".
// If it doesn't find any settings files, it returns a single testCaseSettings with default feature flags and Kong version.
func resolveSetsOfSettingsForTestCaseDir(t *testing.T, path string) []testCaseSettings {
	t.Helper()

	// Iterate over all files in the directory and look for settings files.
	files, err := os.ReadDir(path)
	require.NoErrorf(t, err, "failed to iterate over files in test case directory %s", path)

	setsOfTranslatorSettings := []testCaseSettings{
		// Always include a testCaseSettings with default feature flags and Kong version.
		{
			name:         "default",
			featureFlags: defaultFeatureFlags(),
		},
	}

	// Iterate over all settings files and create a testCaseSettings for each.
	for _, file := range files {
		require.False(t, file.IsDir(), "unexpected directory %s in test case directory %s", file.Name(), path)

		// Skip files that are not settings files.
		if !strings.HasSuffix(file.Name(), settingsFileSuffix) {
			continue
		}

		require.NotEqual(t, "default_settings.yaml", file.Name(),
			"settings file name must not be default_settings.yaml - it's reserved for the default settings")

		// Load the settings file and use it.
		settingsFile := filepath.Join(path, file.Name())
		setsOfTranslatorSettings = append(setsOfTranslatorSettings, unmarshalSettingsFile(t, settingsFile))
	}

	return setsOfTranslatorSettings
}

type testCaseSettings struct {
	name         string
	featureFlags translator.FeatureFlags
}

// unmarshalSettingsFile unmarshals a settings file and returns a testCaseSettings struct.
// All feature flags and Kong version specified in the settings file will be used to override the defaults.
func unmarshalSettingsFile(t *testing.T, path string) testCaseSettings {
	// It specifies only the json tags, because we're using "sigs.k8s.io/yaml" to unmarshal the file and that
	// package respects only json tags: "Unmarshal converts YAML to JSON then uses JSON to unmarshal into an object".
	type settingsFile struct {
		FeatureFlags map[string]bool `json:"feature_flags"`
	}

	b, err := os.ReadFile(path)
	require.NoErrorf(t, err, "Failed to read settings file %s", path)

	var settings settingsFile
	err = yaml.Unmarshal(b, &settings)
	require.NoErrorf(t, err, "Failed to unmarshal settings file %s", path)

	// Construct translator settings name from the file name without the extension.
	settingsName := strings.TrimSuffix(filepath.Base(path), settingsFileSuffix)

	featureFlags := defaultFeatureFlags()
	// Override the feature flags if specified in the settings file.
	for featureFlagName, featureFlagValue := range settings.FeatureFlags {
		field := reflect.ValueOf(&featureFlags).Elem().FieldByName(featureFlagName)
		require.Truef(t, field.IsValid(),
			"invalid feature flag %s from %s, its name has to match one of translator.FeatureFlag's fields", featureFlagName, path)

		t.Logf("%s: Setting feature flag %s to %v", path, featureFlagName, featureFlagValue)
		field.SetBool(featureFlagValue)
	}

	return testCaseSettings{
		name:         settingsName,
		featureFlags: featureFlags,
	}
}

// kongClientGoldenTestCase represents a single test case for the KongClient with an input file and an expected output golden
// file for a specific combination of feature flags.
type kongClientGoldenTestCase struct {
	// k8sConfigFile is the path to the input file with K8s objects to be loaded into the store.
	k8sConfigFile string
	// goldenFile is the path to the expected output golden file.
	goldenFile string
	// featureFlags is the set of Translator feature flags to use in the test case.
	featureFlags translator.FeatureFlags
}

// runKongClientGoldenTest runs a single golden test case for the KongClient.
func runKongClientGoldenTest(t *testing.T, tc kongClientGoldenTestCase) {
	t.Helper()

	t.Logf("Running test case with input file %s and golden file %s", tc.k8sConfigFile, tc.goldenFile)
	t.Logf("Feature flags: %+v", tc.featureFlags)

	// Load the K8s objects from the YAML file.
	objects := extractObjectsFromYAML(t, tc.k8sConfigFile)
	t.Logf("Found %d K8s objects to be loaded into the store", len(objects))

	// Load the K8s objects into the store.
	cacheStores, err := store.NewCacheStoresFromObjYAML(objects...)
	require.NoError(t, err, "failed creating cache stores")

	var objectsToBeConsideredBroken []client.Object
	for _, s := range cacheStores.ListAllStores() {
		for _, o := range s.List() {
			o := o.(client.Object)
			if o.GetAnnotations()["test.konghq.com/broken"] == "true" {
				objectsToBeConsideredBroken = append(objectsToBeConsideredBroken, o)
			}
		}
	}

	// Create the translator.
	logger := zapr.NewLogger(zap.NewNop())
	s := store.New(cacheStores, "kong", logger)
	p, err := translator.NewTranslator(logger, s, "", tc.featureFlags, fakeSchemaServiceProvier{})
	require.NoError(t, err, "failed creating translator")

	// Start a mock Admin API server and create an Admin API client for inspecting the configuration.
	t.Log("Starting mock Admin API server")
	var adminAPIOpts []mocks.AdminAPIHandlerOpt
	if len(objectsToBeConsideredBroken) > 0 {
		t.Logf("Configuring the mock Admin API server to return errors for broken objects: %v", objectsToBeConsideredBroken)
		adminAPIOpts = append(adminAPIOpts,
			mocks.WithConfigPostError(buildPostConfigErrorResponseWithBrokenObjects(objectsToBeConsideredBroken)),
			mocks.WithConfigPostErrorOnlyOnFirstRequest(),
		)
	}

	adminAPIHandler := mocks.NewAdminAPIHandler(t, adminAPIOpts...)
	adminAPIServer := httptest.NewServer(adminAPIHandler)
	defer adminAPIServer.Close()

	t.Log("Creating Admin API client")
	adminAPIClient, err := adminapi.NewTestClient(adminAPIServer.URL)
	require.NoError(t, err)

	// Create the KongClient using _mostly_ real dependencies' implementations (except for the clients provider
	// as we want to avoid spinning up a real Kong Gateway to keep the tests fast).
	t.Log("Building KongClient")
	const timeout = time.Second
	cfg := sendconfig.Config{
		InMemory:              true, // We're running in DB-less mode only for now. In the future, we may want to test DB mode as well.
		ExpressionRoutes:      tc.featureFlags.ExpressionRoutes,
		FallbackConfiguration: len(objectsToBeConsideredBroken) > 0,
	}
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{adminAPIClient},
	}
	updateStrategyResolver := sendconfig.NewDefaultUpdateStrategyResolver(cfg, logger)
	lastValidConfigFetcher := configfetcher.NewDefaultKongLastGoodConfigFetcher(tc.featureFlags.FillIDs, "default")
	fallbackConfigGenerator := fallback.NewGenerator(fallback.NewDefaultCacheGraphProvider(), logger)
	kongClient, err := NewKongClient(
		logger,
		timeout,
		diagnostics.ClientDiagnostic{},
		cfg,
		mocks.NewEventRecorder(),
		dpconf.DBModeOff, // Test will run in DB-less mode only for now. In the future, we may want to test DB mode as well.
		clientsProvider,
		updateStrategyResolver,
		sendconfig.NewDefaultConfigurationChangeDetector(logger),
		lastValidConfigFetcher,
		p,
		&cacheStores,
		fallbackConfigGenerator,
	)
	require.NoError(t, err)

	t.Log("Triggering KongClient.Update")
	ctx := context.Background()
	err = kongClient.Update(ctx)
	if len(objectsToBeConsideredBroken) > 0 {
		require.Error(t, err, "expected an error when fallback configuration is enabled")
	} else {
		require.NoError(t, err, "failed updating Kong configuration")
	}

	t.Log("Fetching the last received configuration from the Admin API")
	resultB, err := adminAPIClient.AdminAPIClient().Config(ctx)
	require.NoError(t, err)

	// If the update flag is set, update the golden file with the result...
	if *updateGolden {
		err = os.WriteFile(tc.goldenFile, resultB, 0o600)
		require.NoError(t, err, "failed writing to golden file")
		t.Logf("Updated golden file %s", tc.goldenFile)
	} else {
		// ...otherwise, compare the result to the golden file.
		const commandToRegenerateGoldenFile = "make test.golden.update"

		goldenB, err := os.ReadFile(tc.goldenFile)
		require.NoError(t, err, "Failed reading golden file.\n"+
			"If it's missing, you can generate it by running:\n"+
			"$ %s\n"+
			"Make sure to carefully inspect the generated golden file output\n"+
			"to ensure it matches the expectations.", commandToRegenerateGoldenFile)

		require.Equalf(t, string(goldenB), string(resultB),
			"Golden file %s does not match the result. \n"+
				"If you are sure the result is correct, update the golden file: \n"+
				"$ %s", tc.goldenFile, commandToRegenerateGoldenFile)
		t.Logf("Successfully compared result to golden file %s", tc.goldenFile)
	}
}

func extractObjectsFromYAML(t *testing.T, filePath string) [][]byte {
	y, err := os.ReadFile(filePath)
	require.NoErrorf(t, err, "failed reading input file: %s", filePath)

	// Strip out the YAML comments.
	f := util.ManualStrip(y)

	// Split the YAML by the document separator.
	objects := bytes.Split(f, []byte("---"))

	// Filter out empty YAML documents.
	return lo.Filter(objects, func(o []byte, _ int) bool {
		return len(bytes.TrimSpace(o)) > 0
	})
}

// fakeSchemaServiceProvier is a stub implementation of the SchemaServiceProvider interface that returns an
// UnavailableSchemaService. It's used to avoid hitting the Kong Admin API during tests.
type fakeSchemaServiceProvier struct{}

func (p fakeSchemaServiceProvier) GetSchemaService() kong.AbstractSchemaService {
	return fakeSchemaService{}
}

// fakeSchemaService is a stub implementation of the kong.AbstractSchemaService interface returning hardcoded schemas
// for testing purposes.
type fakeSchemaService struct{}

func (f fakeSchemaService) Get(_ context.Context, entityType string) (kong.Schema, error) {
	switch entityType {
	case "degraphql_routes":
		return kong.Schema{
			"fields": []interface{}{
				map[string]interface{}{
					"service": map[string]interface{}{
						"type":      "foreign",
						"reference": "services",
					},
				},
			},
		}, nil
	default:
		return kong.Schema{}, nil
	}
}

func (f fakeSchemaService) Validate(context.Context, kong.EntityType, any) (bool, string, error) {
	return true, "", nil
}

func buildPostConfigErrorResponseWithBrokenObjects(brokenObjects []client.Object) []byte {
	var flattenedErrors []string
	for _, o := range brokenObjects {
		gvk := o.GetObjectKind().GroupVersionKind()
		flattenedError := fmt.Sprintf(`{"errors": [{"messages": ["broken object"]}], "entity_tags": ["k8s-name:%s","k8s-namespace:%s","k8s-kind:%s","k8s-group:%s", "k8s-version:%s", "k8s-uid:%s"]}`,
			o.GetName(), o.GetNamespace(), gvk.Kind, gvk.Group, gvk.Version, o.GetUID(),
		)
		flattenedErrors = append(flattenedErrors, flattenedError)
	}

	return []byte(fmt.Sprintf(`{"flattened_errors": [%s]}`, strings.Join(flattenedErrors, ",")))
}
