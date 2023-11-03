//go:build conformance_tests

package conformance

import (
	"testing"

	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

var commonSkippedTests = []string{
	// extended conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4563
	tests.GatewayWithAttachedRoutesWithPort8080.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4913
	tests.GatewayStaticAddresses.ShortName,

	tests.HTTPRouteRedirectPortAndScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMirror.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMultipleMirrors.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4914
	tests.HTTPRouteTimeoutBackendRequest.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4915
	tests.HTTPRouteTimeoutRequest.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4983
	tests.HTTPRouteBackendProtocolH2C.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4984
	tests.HTTPRouteBackendProtocolWebSocket.ShortName,

	// experimental conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3684
	tests.HTTPRouteRedirectPath.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3685
	tests.HTTPRouteRewriteHost.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3686
	tests.HTTPRouteRewritePath.ShortName,

	// TLS
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4562
	tests.TLSRouteInvalidReferenceGrant.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3678
	tests.TLSRouteSimpleSameNamespace.ShortName,
}

var (
	skippedTestsForExpressionRoutes  = commonSkippedTests
	skippedTestsForTraditionalRoutes = append(
		commonSkippedTests,
		// core conformance
		tests.HTTPRouteHeaderMatching.ShortName,
		// extended conformance
		// only 10 and 11 broken because traditional/traditional_compatible router
		// cannot support the path > method > header precedence,
		// but no way to omit individual cases.
		tests.HTTPRouteMethodMatching.ShortName,
		// only expression router supports query param matches
		tests.HTTPRouteQueryParamMatching.ShortName,
	)
)

func TestGatewayConformance(t *testing.T) {
	if shouldRunExperimentalConformance() {
		t.Skip("skipping standard conformance tests")
	}

	k8sClient, gatewayClassName := prepareEnvForGatewayConformanceTests(t)
	// Conformance tests are run for both configs with and without
	// KONG_TEST_EXPRESSION_ROUTES='true'.
	skipTests := skippedTestsForTraditionalRoutes
	if testenv.ExpressionRoutesEnabled() {
		skipTests = skippedTestsForExpressionRoutes
	}

	cSuite := suite.New(suite.Options{
		Client:                     k8sClient,
		GatewayClassName:           gatewayClassName,
		Debug:                      true,
		CleanupBaseResources:       !testenv.IsCI(),
		EnableAllSupportedFeatures: true,
		ExemptFeatures:             suite.MeshCoreFeatures,
		BaseManifests:              conformanceTestsBaseManifests,
		SkipTests:                  skipTests,
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
