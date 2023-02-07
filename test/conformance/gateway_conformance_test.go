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

var (
	showDebug     = true
	shouldCleanup = true

	conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", consts.GatewayRawRepoURL)
)

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
	gwc := &gatewayv1beta1.GatewayClass{
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
	require.NoError(t, client.Create(ctx, gwc))
	t.Cleanup(func() { assert.NoError(t, client.Delete(ctx, gwc)) })

	t.Log("starting the gateway conformance test suite")
	cSuite := suite.New(suite.Options{
		Client:               client,
		GatewayClassName:     gwc.Name,
		Debug:                showDebug,
		CleanupBaseResources: shouldCleanup,
		BaseManifests:        conformanceTestsBaseManifests,
		SupportedFeatures: map[suite.SupportedFeature]bool{
			suite.SupportReferenceGrant: true,

			// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/2778
			suite.SupportHTTPRouteQueryParamMatching: false,
		},
	})
	cSuite.Setup(t)

	t.Log("configuring gateway conformance tests")
	for i := range tests.ConformanceTests {
		for j, manifest := range tests.ConformanceTests[i].Manifests {
			tests.ConformanceTests[i].Manifests[j] = fmt.Sprintf("%s/conformance/%s", consts.GatewayRawRepoURL, manifest)
		}
	}

	// these tests are temporarily disabled to be able to bump the Gateway API to 0.6
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3305
	skipTests := map[string]struct{}{
		tests.GatewayInvalidRouteKind.ShortName:                   {},
		tests.GatewayInvalidTLSConfiguration.ShortName:            {},
		tests.GatewaySecretReferenceGrantAllInNamespace.ShortName: {},
		tests.GatewaySecretReferenceGrantSpecific.ShortName:       {},
		tests.HTTPRouteHeaderMatching.ShortName:                   {},
		tests.HTTPRouteObservedGenerationBump.ShortName:           {},
		tests.HTTPRouteRequestRedirect.ShortName:                  {},
	}

	t.Log("running gateway conformance tests")
	for _, tt := range tests.ConformanceTests {
		tt := tt
		t.Run(tt.Description, func(t *testing.T) {
			if _, ok := skipTests[tt.ShortName]; ok {
				t.Skipf("skipping %s", tt.ShortName)
			}
			tt.Run(t, cSuite)
		})
	}
}
