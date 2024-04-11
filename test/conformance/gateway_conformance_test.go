//go:build conformance_tests

package conformance

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/sets"
	conformancev1alpha1 "sigs.k8s.io/gateway-api/conformance/apis/v1alpha1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/yaml"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

var skippedTestsForTraditionalRoutes = []string{
	// core conformance
	tests.HTTPRouteHeaderMatching.ShortName,
}

var traditionalRoutesSupportedFeatures = []suite.SupportedFeature{
	// core features
	suite.SupportGateway,
	suite.SupportHTTPRoute,
	// extended features
	suite.SupportHTTPRouteResponseHeaderModification,
	suite.SupportHTTPRouteBackendTimeout,
}

var expressionRoutesSupportedFeatures = []suite.SupportedFeature{
	// core features
	suite.SupportGateway,
	suite.SupportHTTPRoute,
	// extended features
	suite.SupportHTTPRouteQueryParamMatching,
	suite.SupportHTTPRouteMethodMatching,
	suite.SupportHTTPRouteResponseHeaderModification,
	suite.SupportHTTPRouteBackendTimeout,
}

func TestGatewayConformance(t *testing.T) {
	k8sClient, gatewayClassName := prepareEnvForGatewayConformanceTests(t)

	// Conformance tests are run for both available router flavours:
	// traditional_compatible and expressions.
	var (
		skippedTests      []string
		supportedFeatures []suite.SupportedFeature
	)
	switch rf := testenv.KongRouterFlavor(); rf {
	case dpconf.RouterFlavorTraditionalCompatible:
		skippedTests = skippedTestsForTraditionalRoutes
		supportedFeatures = traditionalRoutesSupportedFeatures
	case dpconf.RouterFlavorExpressions:
		supportedFeatures = expressionRoutesSupportedFeatures
	default:
		t.Fatalf("unsupported KongRouterFlavor: %s", rf)
	}

	cSuite, err := suite.NewExperimentalConformanceTestSuite(
		suite.ExperimentalConformanceOptions{
			Options: suite.Options{
				Client:               k8sClient,
				GatewayClassName:     gatewayClassName,
				Debug:                true,
				CleanupBaseResources: !testenv.IsCI(),
				BaseManifests:        conformanceTestsBaseManifests,
				SupportedFeatures:    sets.New(supportedFeatures...),
				SkipTests:            skippedTests,
			},
			ConformanceProfiles: sets.New(
				suite.HTTPConformanceProfileName,
			),
			Implementation: conformancev1alpha1.Implementation{
				Organization: metadata.Organization,
				Project:      metadata.ProjectName,
				URL:          metadata.ProjectURL,
				Version:      metadata.Release,
				Contact: []string{
					path.Join(metadata.ProjectURL, "/issues/new/choose"),
				},
			},
		},
	)
	require.NoError(t, err)

	t.Log("starting the gateway conformance test suite")
	cSuite.Setup(t)

	go patchGatewayClassToPassTestGatewayClassObservedGenerationBump(ctx, t, k8sClient)

	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	// cSuite.Run(t, []suite.ConformanceTest{tests.GatewayClassObservedGenerationBump})
	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	//cSuite.Run(t, []suite.ConformanceTest{tests.HTTPRouteRedirectPortAndScheme})
	require.NoError(t, cSuite.Run(t, tests.ConformanceTests))

	const reportFileName = "kong-kubernetes-ingress-controller.yaml"
	t.Log("saving the gateway conformance test report to file:", reportFileName)
	report, err := cSuite.Report()
	require.NoError(t, err)
	rawReport, err := yaml.Marshal(report)
	require.NoError(t, err)
	// Save report in root of the repository, file name is in .gitignore.
	require.NoError(t, os.WriteFile("../../"+reportFileName, rawReport, 0o600))
}
