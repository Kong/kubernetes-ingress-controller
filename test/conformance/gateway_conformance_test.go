//go:build conformance_tests

package conformance

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

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
}

var expressionRoutesSupportedFeatures = []suite.SupportedFeature{
	// core features
	suite.SupportGateway,
	suite.SupportHTTPRoute,
	// extended features
	suite.SupportHTTPRouteQueryParamMatching,
	suite.SupportHTTPRouteMethodMatching,
	suite.SupportHTTPRouteResponseHeaderModification,
}

func TestGatewayConformance(t *testing.T) {
	if shouldRunExperimentalConformance() {
		t.Skip("skipping standard conformance tests")
	}

	k8sClient, gatewayClassName := prepareEnvForGatewayConformanceTests(t)
	// Conformance tests are run for both configs with and without
	// KONG_TEST_EXPRESSION_ROUTES='true'.
	var skippedTests []string
	var supportedFeatures []suite.SupportedFeature
	if testenv.ExpressionRoutesEnabled() {
		supportedFeatures = expressionRoutesSupportedFeatures
	} else {
		skippedTests = skippedTestsForTraditionalRoutes
		supportedFeatures = traditionalRoutesSupportedFeatures
	}

	cSuite := suite.New(suite.Options{
		Client:               k8sClient,
		GatewayClassName:     gatewayClassName,
		Debug:                true,
		CleanupBaseResources: !testenv.IsCI(),
		SupportedFeatures:    sets.New(supportedFeatures...),
		ExemptFeatures:       suite.MeshCoreFeatures,
		BaseManifests:        conformanceTestsBaseManifests,
		SkipTests:            skippedTests,
	})
	t.Log("starting the gateway conformance test suite")
	cSuite.Setup(t)

	go patchGatewayClassToPassTestGatewayClassObservedGenerationBump(ctx, t, k8sClient)

	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	// cSuite.Run(t, []suite.ConformanceTest{tests.GatewayClassObservedGenerationBump})
	cSuite.Run(t, tests.ConformanceTests)
}
