//go:build conformance_tests

package conformance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

var commonSkippedTests = []string{
	// HTTP
	// extended conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4563
	tests.GatewayWithAttachedRoutesWithPort8080.ShortName,
	tests.HTTPRouteRedirectPortAndScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3679
	tests.HTTPRouteQueryParamMatching.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3681
	tests.HTTPRouteRedirectPort.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3682
	tests.HTTPRouteRedirectScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMirror.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMultipleMirrors.ShortName,

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

	// We need to wait for the GatewayClass created by the test GatewayClassObservedGenerationBump
	// and patch it with the unmanaged annotation to make it reconciled by the GatewayClass
	// controller. The timeout and the tick are pretty loose because of the indeterministic
	// test order execution.
	go require.Eventually(t, func() bool {
		return ensureTestGatewayClassIsUnmanaged(ctx, k8sClient)
	}, 10*time.Minute, 10*time.Second)

	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	//cSuite.Run(t, []suite.ConformanceTest{tests.GatewayClassObservedGenerationBump})
	cSuite.Run(t, tests.ConformanceTests)
}

func ensureTestGatewayClassIsUnmanaged(ctx context.Context, k8sClient client.Client) bool {
	gwcNamespacedName := k8stypes.NamespacedName{Name: "gatewayclass-observed-generation-bump"}
	gwc := &gatewayv1.GatewayClass{}
	if err := k8sClient.Get(ctx, gwcNamespacedName, gwc); err != nil {
		return false
	}
	if gwc.Annotations == nil {
		gwc.Annotations = map[string]string{}
	}
	gwc.Annotations[annotations.GatewayClassUnmanagedAnnotation] = annotations.GatewayClassUnmanagedAnnotationValuePlaceholder
	if err := k8sClient.Update(ctx, gwc); err != nil {
		return false
	}
	return true
}
