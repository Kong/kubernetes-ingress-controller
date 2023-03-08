//go:build conformance_tests
// +build conformance_tests

package conformance

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

const (
	showDebug                  = true
	shouldCleanup              = true
	enableAllSupportedFeatures = true
)

var conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", consts.GatewayRawRepoURL)

func TestGatewayConformance(t *testing.T) {
	t.Log("configuring environment for gateway conformance tests")
	client, err := client.New(env.Cluster().Config(), client.Options{})
	require.NoError(t, err)
	require.NoError(t, gatewayv1alpha2.AddToScheme(client.Scheme()))
	require.NoError(t, gatewayv1beta1.AddToScheme(client.Scheme()))

	t.Log("starting the controller manager")
	args := []string{
		fmt.Sprintf("--ingress-class=%s", ingressClass),
		fmt.Sprintf("--admission-webhook-cert=%s", testutils.KongSystemServiceCert),
		fmt.Sprintf("--admission-webhook-key=%s", testutils.KongSystemServiceKey),
		fmt.Sprintf("--admission-webhook-listen=%s:%d", testutils.AdmissionWebhookListenHost, testutils.AdmissionWebhookListenPort),
		"--profiling",
		"--dump-config",
		"--log-level=trace",
		"--debug-log-reduce-redundancy",
		fmt.Sprintf("--feature-gates=%s", consts.DefaultFeatureGates),
		"--anonymous-reports=false",
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

	t.Log("starting the gateway conformance test suite")
	cSuite := suite.New(suite.Options{
		Client:                     client,
		GatewayClassName:           gatewayClass.Name,
		Debug:                      showDebug,
		CleanupBaseResources:       shouldCleanup,
		EnableAllSupportedFeatures: enableAllSupportedFeatures,
		BaseManifests:              conformanceTestsBaseManifests,
		SkipTests: []string{
			// this test is currently fixed but cannot be re-enabled yet due to an upstream issue
			// https://github.com/kubernetes-sigs/gateway-api/pull/1745
			tests.GatewaySecretReferenceGrantSpecific.ShortName,

			// standard conformance
			tests.HTTPRouteHeaderMatching.ShortName,
			tests.HTTPRouteRedirectHostAndStatus.ShortName,

			// extended conformance
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
			// https://github.com/Kong/kubernetes-ingress-controller/issues/3683
			tests.HTTPRouteResponseHeaderModifier.ShortName,

			// experimental conformance
			// https://github.com/Kong/kubernetes-ingress-controller/issues/3684
			tests.HTTPRouteRedirectPath.ShortName,
			// https://github.com/Kong/kubernetes-ingress-controller/issues/3685
			tests.HTTPRouteRewriteHost.ShortName,
			// https://github.com/Kong/kubernetes-ingress-controller/issues/3686
			tests.HTTPRouteRewritePath.ShortName,
		},
	})
	cSuite.Setup(t)
	cSuite.Run(t, tests.ConformanceTests)
}
