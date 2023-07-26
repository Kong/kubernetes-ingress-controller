//go:build conformance_tests

package conformance

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
)

const (
	showDebug                  = true
	shouldCleanup              = true
	enableAllSupportedFeatures = true
)

var conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", consts.GatewayRawRepoURL)

var skippedTestsForExpressionRoutes = []string{
	// extended conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4166
	// requires an 8080 listener, which our manually-built test gateway does not have
	tests.HTTPRouteRedirectPortAndScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3680
	tests.GatewayClassObservedGenerationBump.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3678
	tests.TLSRouteSimpleSameNamespace.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3679
	tests.HTTPRouteQueryParamMatching.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3681
	tests.HTTPRouteRedirectPort.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3682
	tests.HTTPRouteRedirectScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMirror.ShortName,

	// experimental conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3684
	tests.HTTPRouteRedirectPath.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3685
	tests.HTTPRouteRewriteHost.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3686
	tests.HTTPRouteRewritePath.ShortName,
}

var skippedTestsForTraditionalRoutes = []string{
	// core conformance
	tests.HTTPRouteHeaderMatching.ShortName,

	// extended conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4164
	// only 10 and 11 broken because traditional/traditional_compatible router
	// cannot support the path > method > header precedence,
	// but no way to omit individual cases.
	tests.HTTPRouteMethodMatching.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4166
	// requires an 8080 listener, which our manually-built test gateway does not have
	tests.HTTPRouteRedirectPortAndScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3680
	tests.GatewayClassObservedGenerationBump.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3678
	//tests.TLSRouteSimpleSameNamespace.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3679
	tests.HTTPRouteQueryParamMatching.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3681
	tests.HTTPRouteRedirectPort.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3682
	tests.HTTPRouteRedirectScheme.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4165
	tests.HTTPRouteRequestMirror.ShortName,

	// experimental conformance
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3684
	tests.HTTPRouteRedirectPath.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3685
	tests.HTTPRouteRewriteHost.ShortName,
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3686
	tests.HTTPRouteRewritePath.ShortName,
}

func TestGatewayConformance(t *testing.T) {
	t.Log("configuring environment for gateway conformance tests")
	client, err := client.New(env.Cluster().Config(), client.Options{})
	require.NoError(t, err)
	require.NoError(t, gatewayv1alpha2.AddToScheme(client.Scheme()))
	require.NoError(t, gatewayv1beta1.AddToScheme(client.Scheme()))

	featureGateFlag := fmt.Sprintf("--feature-gates=%s", consts.DefaultFeatureGates)
	if expressionRoutesEnabled() {
		t.Log("expression routes enabled")
		featureGateFlag = fmt.Sprintf("--feature-gates=%s", consts.ConformanceExpressionRoutesTestsFeatureGates)
	}

	t.Log("starting the controller manager")
	cert, key := certificate.GetKongSystemSelfSignedCerts()
	args := []string{
		fmt.Sprintf("--ingress-class=%s", ingressClass),
		fmt.Sprintf("--admission-webhook-cert=%s", cert),
		fmt.Sprintf("--admission-webhook-key=%s", key),
		fmt.Sprintf("--admission-webhook-listen=%s:%d", testutils.AdmissionWebhookListenHost, testutils.AdmissionWebhookListenPort),
		"--profiling",
		"--dump-config",
		"--log-level=trace",
		"--debug-log-reduce-redundancy",
		featureGateFlag,
		"--anonymous-reports=false",
		"--publish-service-tls=kong-system/ingress-controller-kong-tls-proxy",
	}

	require.NoError(t, testutils.DeployControllerManagerForCluster(ctx, globalDeprecatedLogger, globalLogger, env.Cluster(), args...))

	t.Log("creating GatewayClass for gateway conformance tests")
	gatewayClass := &gatewayv1beta1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayv1beta1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}
	require.NoError(t, client.Create(ctx, gatewayClass))
	t.Cleanup(func() { assert.NoError(t, client.Delete(ctx, gatewayClass)) })

	exemptFeatures := sets.New(suite.SupportMesh)

	t.Log("starting the gateway conformance test suite")
	skippedTests := skippedTestsForTraditionalRoutes
	if expressionRoutesEnabled() {
		skippedTests = skippedTestsForExpressionRoutes
	}
	cSuite := suite.New(suite.Options{
		Client:               client,
		GatewayClassName:     gatewayClass.Name,
		Debug:                showDebug,
		CleanupBaseResources: shouldCleanup,
		ExemptFeatures:       exemptFeatures,
		BaseManifests:        conformanceTestsBaseManifests,
		SkipTests:            skippedTests,
		SupportedFeatures:    suite.TLSCoreFeatures,
	})
	cSuite.Setup(t)
	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	//cSuite.Run(t, []suite.ConformanceTest{tests.HTTPRouteRedirectPortAndScheme})
	cSuite.Run(t, tests.ConformanceTests)
}
