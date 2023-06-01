package parser_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

var (
	// updateGolden tells whether to update golden files using the current output of the parser.
	updateGolden = os.Getenv("UPDATE_GOLDEN") == "true"

	// defaultKongVersion is the default Kong version to use in tests. Can be overridden in a test case.
	defaultKongVersion = semver.MustParse("3.3.0")

	// defaultFeatureFlags is the default set of feature flags to use in tests. Can be overridden in a test case.
	defaultFeatureFlags = func() parser.FeatureFlags {
		return parser.FeatureFlags{
			CombinedServiceRoutes: true,
			RegexPathPrefix:       true,
		}
	}
)

// TestParser_GoldenTests runs the golden tests for the parser.
//
// Command to update the golden files:
// $ make test.golden.update
//
// Data for the test cases is stored in the ./testdata/golden directory. It's recommended to group test cases into
// subdirectories based on the Kubernetes input that they run against so that each of the subdirectories has:
//   - an input file that represents the input with Kubernetes objects to be loaded into the store (in.yaml),
//   - a set of expected golden *.yaml files (in Deck format) where each file represents an expected output for a given
//     parser configuration (feature flags, Kong version, etc.).
//
// The test case is executed by loading the in.yaml file into the store, then running the parser on the store,
// and finally comparing the output of the parser with the expected golden file.
//
// When adding a new test case, you can follow these steps:
// - Add a test case to the testCases slice below along with.
// - Add a new directory ./testdata/golden/<your-dir> with the in.yaml that you want to test against.
// - Specify the path for the golden file (./testdata/golden/<your-dir>/<your-parser-config>.yaml).
// - Run `make test.golden.update` to generate the golden file.
// - Inspect the generated golden file and make sure it's correct.
//
// If you introduce a change that may affect many test cases, and you're sure about it correctness, you can run the
// command as well to update all golden files at once.
func TestParser_GoldenTests(t *testing.T) {
	const (
		ingressV1SingleServiceInMultipleIngressesK8sConfigFile = "testdata/golden/ingress-v1-single-service-in-multiple-ingresses/in.yaml"
		ingressV1WithDefaultBackendK8sConfigFile               = "testdata/golden/ingress-v1-with-default-backend/in.yaml"
		ingressV1RegexPrefixExactRuleK8sConfig                 = "testdata/golden/ingress-v1-regex-prefix-exact-rule/in.yaml"
		ingressV1RuleWithTLSK8sConfig                          = "testdata/golden/ingress-v1-rule-with-tls/in.yaml"
		ingressV1WithAcmeLikePathK8sConfig                     = "testdata/golden/ingress-v1-with-acme-like-path/in.yaml"
		ingressV1EmptyPathK8sConfig                            = "testdata/golden/ingress-v1-empty-path/in.yaml"
		ingressV1MultiplePortsForOneServiceK8sConfig           = "testdata/golden/ingress-v1-multiple-ports-for-one-service/in.yaml"
		ingressV1RegexPrefixedPathK8sConfig                    = "testdata/golden/ingress-v1-regex-prefixed-path/in.yaml"
		ingressV1PortsDefinedByNameK8sConfig                   = "testdata/golden/ingress-v1-ports-defined-by-name/in.yaml"
	)

	testCases := []parserGoldenTestCase{
		// Test cases for testdata/golden/ingress-v1-single-service-in-multiple-ingresses/in.yaml.
		{
			k8sConfigFile: ingressV1SingleServiceInMultipleIngressesK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-single-service-in-multiple-ingresses/default.yaml",
		},
		{
			k8sConfigFile: ingressV1SingleServiceInMultipleIngressesK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-single-service-in-multiple-ingresses/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1SingleServiceInMultipleIngressesK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-single-service-in-multiple-ingresses/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
		// Test cases for testdata/golden/ingress-v1-with-default-backend/in.yaml.
		{
			k8sConfigFile: ingressV1WithDefaultBackendK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-with-default-backend/default.yaml",
		},
		{
			k8sConfigFile: ingressV1WithDefaultBackendK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-with-default-backend/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1WithDefaultBackendK8sConfigFile,
			goldenFile:    "testdata/golden/ingress-v1-with-default-backend/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
		// Test cases for testdata/golden/ingress-v1-regex-prefix-exact-rule/in.yaml.
		{
			k8sConfigFile: ingressV1RegexPrefixExactRuleK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-regex-prefix-exact-rule/default.yaml",
		},
		{
			k8sConfigFile: ingressV1RegexPrefixExactRuleK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-regex-prefix-exact-rule/regex-path-prefix-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.RegexPathPrefix = false
			},
		},
		// Test cases for testdata/golden/ingress-v1-rule-with-tls/in.yaml.
		{
			k8sConfigFile: ingressV1RuleWithTLSK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-rule-with-tls/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1RuleWithTLSK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-rule-with-tls/combined-routes-on.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = true
			},
		},
		{
			k8sConfigFile: ingressV1RuleWithTLSK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-rule-with-tls/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
		// Test cases for testdata/golden/ingress-v1-with-acme-like-path/in.yaml.
		{
			k8sConfigFile: ingressV1WithAcmeLikePathK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-with-acme-like-path/default.yaml",
		},
		// Test cases for testdata/golden/ingress-v1-empty-path/in.yaml.
		{
			k8sConfigFile: ingressV1EmptyPathK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-empty-path/default.yaml",
		},
		{
			k8sConfigFile: ingressV1EmptyPathK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-empty-path/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1EmptyPathK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-empty-path/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
		// Test cases for testdata/golden/ingress-v1-multiple-ports-for-one-service/in.yaml.
		{
			k8sConfigFile: ingressV1MultiplePortsForOneServiceK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-multiple-ports-for-one-service/default.yaml",
		},
		{
			k8sConfigFile: ingressV1MultiplePortsForOneServiceK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-multiple-ports-for-one-service/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1MultiplePortsForOneServiceK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-multiple-ports-for-one-service/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
		// Test cases for testdata/golden/ingress-v1-regex-prefixed-path/in.yaml.
		{
			k8sConfigFile: ingressV1RegexPrefixedPathK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-regex-prefixed-path/default.yaml",
		},
		// Test cases for testdata/golden/ingress-v1-ports-defined-by-name/in.yaml.
		{
			k8sConfigFile: ingressV1PortsDefinedByNameK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-ports-defined-by-name/default.yaml",
		},
		{
			k8sConfigFile: ingressV1PortsDefinedByNameK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-ports-defined-by-name/combined-routes-off.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServiceRoutes = false
			},
		},
		{
			k8sConfigFile: ingressV1PortsDefinedByNameK8sConfig,
			goldenFile:    "testdata/golden/ingress-v1-ports-defined-by-name/combined-services.yaml",
			featureFlagsModifier: func(flags *parser.FeatureFlags) {
				flags.CombinedServices = true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("in=%s,out=%s", tc.k8sConfigFile, tc.goldenFile), func(t *testing.T) {
			runParserGoldenTest(t, tc)
		})
	}
}

type parserGoldenTestCase struct {
	k8sConfigFile        string
	goldenFile           string
	featureFlagsModifier func(flags *parser.FeatureFlags)
	kongVersion          *semver.Version
}

func runParserGoldenTest(t *testing.T, tc parserGoldenTestCase) {
	logger := logrus.New()

	// Load the K8s objects from the YAML file.
	objects := extractObjectsFromYAML(t, tc.k8sConfigFile)
	t.Logf("Found %d K8s objects to be loaded into the store", len(objects))

	// Load the K8s objects into the store.
	cacheStores, err := store.NewCacheStoresFromObjYAML(objects...)
	require.NoError(t, err, "Failed creating cache stores")

	// Determine the feature flags to use.
	featureFlags := defaultFeatureFlags()

	// Apply test case's feature flags modifier if defined.
	if tc.featureFlagsModifier != nil {
		tc.featureFlagsModifier(&featureFlags)
	}

	// Determine the Kong version to use.
	kongVersion := defaultKongVersion
	if tc.kongVersion != nil {
		kongVersion = *tc.kongVersion
	}

	// Create the parser.
	s := store.New(cacheStores, "kong", logger)
	p, err := parser.NewParser(logger, s, featureFlags, kongVersion)
	require.NoError(t, err, "Failed creating parser")

	// Build the Kong configuration.
	result := p.BuildKongConfig()
	targetConfig := deckgen.ToDeckContent(context.Background(),
		logger,
		result.KongState,
		deckgen.GenerateDeckContentParams{
			FormatVersion:    "3.0",
			ExpressionRoutes: featureFlags.ExpressionRoutes,
			PluginSchemas:    pluginsSchemaStoreStub{},
		},
	)

	// Marshal the result into YAML bytes for comparison.
	resultB, err := yaml.Marshal(targetConfig)
	require.NoError(t, err, "Failed marshalling result")

	// If the update flag is set, update the golden file with the result...
	if updateGolden {
		err = os.WriteFile(tc.goldenFile, resultB, 0o600)
		require.NoError(t, err, "Failed writing to golden file")
		t.Logf("Updated golden file %s", tc.goldenFile)
	} else {
		// ...otherwise, compare the result to the golden file.

		commandToRegenerateGoldenFile := fmt.Sprintf("UPDATE_GOLDEN=true go test -run %s ./internal/dataplane/parser", t.Name())

		goldenB, err := os.ReadFile(tc.goldenFile)
		require.NoError(t, err, "Failed reading golden file.\n"+
			"If it's missing, you can generate it by running:\n"+
			"$ %s\n"+
			"Make sure to carefully inspect the generated golden file output\n"+
			"to ensure it matches the expectations.", commandToRegenerateGoldenFile)

		require.Equalf(t, string(goldenB), string(resultB),
			"Golden file %s does not match the result. \n"+
				"If you are sure the result is correct, run the test "+
				"with the -update-golden flag to update the golden file: \n"+
				"$ %s", tc.goldenFile, commandToRegenerateGoldenFile)
		t.Logf("Successfully compared result to golden file %s", tc.goldenFile)
	}
}

func extractObjectsFromYAML(t *testing.T, filePath string) [][]byte {
	y, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed reading input file")

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
