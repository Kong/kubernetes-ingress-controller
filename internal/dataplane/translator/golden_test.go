package translator_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

var (
	// updateGolden tells whether to update golden files using the current output of the translator.
	updateGolden = flag.Bool("update", false, "update golden files")

	// defaultFeatureFlags is the default set of feature flags to use in tests. Can be overridden in a test case.
	defaultFeatureFlags = func() translator.FeatureFlags {
		defaults := featuregates.GetFeatureGatesDefaults()
		return translator.FeatureFlags{
			// We do not verify configuration reports in golden tests.
			ReportConfiguredKubernetesObjects: false,

			// Feature flags that are directly propagated from the feature gates get their defaults.
			FillIDs:           defaults.Enabled(featuregates.FillIDsFeature),
			KongServiceFacade: defaults.Enabled(featuregates.KongServiceFacade),
		}
	}
)

const (
	goldenDir          = "testdata/golden"
	inFileName         = "in.yaml"
	goldenFileSuffix   = "_golden.yaml"
	settingsFileSuffix = "_settings.yaml"
)

type fakeSchemaServiceProvier struct{}

func (p fakeSchemaServiceProvier) GetSchemaService() kong.AbstractSchemaService {
	return translator.UnavailableSchemaService{}
}

// TestTranslator_GoldenTests runs the golden tests for the translator.
//
// Command to update the golden files:
// $ make test.golden.update
//
// Data for the test cases is stored in the "./testdata/golden" directory. Test cases are grouped into subdirectories
// based on the Kubernetes input that they run against so that each of the subdirectories has:
//   - an input file that represents the input with Kubernetes objects to be loaded into the store: "in.yaml",
//   - a set of "<settings-name>_settings.yaml" files that define the translator configuration for a given test case,
//   - a set of expected golden "<settings-name>_golden.yaml" files (in Deck format) where each file represents an
//     expected output for a given translator configuration defined in "<settings-name>_settings.yaml".
//
// The test case is executed by loading the in.yaml file into the store, then running the translator on the store,
// and finally comparing the output of the translator with the expected golden file.
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
func TestTranslator_GoldenTests(t *testing.T) {
	// First, let's prepare the test cases basing on the testdata/golden directory contents.
	var testCases []translatorGoldenTestCase
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
		for _, translatorSettings := range resolveSetsOfTranslatorSettingsForTestCaseDir(t, testCaseDirPath) {
			testCases = append(testCases, translatorGoldenTestCase{
				k8sConfigFile: filepath.Join(testCaseDirPath, inFileName),
				goldenFile:    filepath.Join(testCaseDirPath, fmt.Sprintf("%s%s", translatorSettings.name, goldenFileSuffix)),
				featureFlags:  translatorSettings.featureFlags,
			})
		}
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("in=%s,out=%s", tc.k8sConfigFile, tc.goldenFile), func(t *testing.T) {
			runTranslatorGoldenTest(t, tc)
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

// resolveSetsOfTranslatorSettingsForTestCaseDir returns a slice of translatorSettings, each of which represents a combination of
// feature flags and Kong version.
// The function iterates over a test case directory containing zero or more files named "<name>_settings.yaml".
// If it doesn't find any settings files, it returns a single translatorSettings with default feature flags and Kong version.
func resolveSetsOfTranslatorSettingsForTestCaseDir(t *testing.T, path string) []translatorSettings {
	// Iterate over all files in the directory and look for settings files.
	files, err := os.ReadDir(path)
	require.NoErrorf(t, err, "failed to iterate over files in test case directory %s", path)

	setsOfTranslatorSettings := []translatorSettings{
		// Always include a translatorSettings with default feature flags and Kong version.
		{
			name:         "default",
			featureFlags: defaultFeatureFlags(),
		},
	}

	// Iterate over all settings files and create a translatorSettings for each.
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

type translatorSettings struct {
	name         string
	featureFlags translator.FeatureFlags
}

// unmarshalSettingsFile unmarshals a settings file and returns a translatorSettings struct.
// All feature flags and Kong version specified in the settings file will be used to override the defaults.
func unmarshalSettingsFile(t *testing.T, path string) translatorSettings {
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

	return translatorSettings{
		name:         settingsName,
		featureFlags: featureFlags,
	}
}

// translatorGoldenTestCase represents a single test case for the translator with an input file and an expected output golden
// file for a specific combination of feature flags and Kong version.
type translatorGoldenTestCase struct {
	k8sConfigFile string
	goldenFile    string
	featureFlags  translator.FeatureFlags
}

func runTranslatorGoldenTest(t *testing.T, tc translatorGoldenTestCase) {
	logger := zapr.NewLogger(zap.NewNop())

	// Load the K8s objects from the YAML file.
	objects := extractObjectsFromYAML(t, tc.k8sConfigFile)
	t.Logf("Found %d K8s objects to be loaded into the store", len(objects))

	// Load the K8s objects into the store.
	cacheStores, err := store.NewCacheStoresFromObjYAML(objects...)
	require.NoError(t, err, "failed creating cache stores")

	// Create the translator.
	s := store.New(cacheStores, "kong", logger)
	p, err := translator.NewTranslator(logger, s, "", tc.featureFlags, fakeSchemaServiceProvier{})
	require.NoError(t, err, "failed creating translator")

	// MustBuild the Kong configuration.
	result := p.BuildKongConfig()
	targetConfig := deckgen.ToDeckContent(context.Background(),
		logger,
		result.KongState,
		deckgen.GenerateDeckContentParams{
			ExpressionRoutes: tc.featureFlags.ExpressionRoutes,
			PluginSchemas:    pluginsSchemaStoreStub{},
		},
	)

	// Marshal the result into YAML bytes for comparison.
	resultB, err := yaml.Marshal(targetConfig)
	require.NoError(t, err, "failed marshalling result")

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

// pluginsSchemaStoreStub is a stub implementation of the plugins.SchemaStore interface that returns an empty schema
// for all plugins. It's used to avoid hitting the Kong Admin API during tests.
type pluginsSchemaStoreStub struct{}

func (p pluginsSchemaStoreStub) Schema(context.Context, string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
