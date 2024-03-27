//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const invalidRegexPath = "/foo[[[["

type testCaseHTTPRouteValidation struct {
	Name                   string
	Route                  *gatewayv1.HTTPRoute
	WantCreateErrSubstring string
}

// commonHTTPRouteValidationTestCases returns a list of test cases for validating HTTPRoutes
// that are common to both traditional and expressions routers (the same error message is returned).
func commonHTTPRouteValidationTestCases(
	managedGateway *gatewayv1.Gateway, unmanagedGateway *gatewayv1.Gateway,
) []testCaseHTTPRouteValidation {
	return []testCaseHTTPRouteValidation{
		{
			Name: "a valid httproute linked to a managed gateway passes validation",
			Route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Namespace: (*gatewayv1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1.ObjectName(managedGateway.Name),
						}},
					},
				},
			},
		},
		{
			Name: "a httproute linked to a non-existent gateway fails validation",
			Route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Namespace: (*gatewayv1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1.ObjectName("fake-gateway"),
						}},
					},
				},
			},
			WantCreateErrSubstring: `Gateway.gateway.networking.k8s.io \"fake-gateway\" not found`,
		},
		{
			Name: "an invalid httproute will pass validation if it's not linked to a managed controller (it's not ours)",
			Route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex(invalidRegexPath).Build(),
						},
					}},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Namespace: (*gatewayv1.Namespace)(&unmanagedGateway.Namespace),
							Name:      gatewayv1.ObjectName(unmanagedGateway.Name),
						}},
					},
				},
			},
		},
		{
			Name: "a httproute with valid regex expressions for a path and a header pass validation",
			Route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{"foo.com"},
					Rules: []gatewayv1.HTTPRouteRule{
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathRegex("/path[1-8]").Build(),
								builder.NewHTTPRouteMatch().WithHeaderRegex("foo", "bar[1-8]").Build(),
							},
						},
					},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Namespace: (*gatewayv1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1.ObjectName(managedGateway.Name),
						}},
					},
				},
			},
		},
	}
}

// invalidRegexInPathTestCase returns a test case for a HTTPRoute with an invalid regex in the path.
// The expected error substring is different for traditional and expressions routers, thus it has
// passed by caller.
func invalidRegexInPathTestCase(
	managedGateway *gatewayv1.Gateway, wantCreateErrSubstring string,
) testCaseHTTPRouteValidation {
	return testCaseHTTPRouteValidation{
		Name: "a httproute with invalid regex for path does not pass validation",
		Route: &gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.NewString(),
			},
			Spec: gatewayv1.HTTPRouteSpec{
				Hostnames: []gatewayv1.Hostname{"foo.com"},
				Rules: []gatewayv1.HTTPRouteRule{
					{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
							builder.NewHTTPRouteMatch().WithPathRegex(invalidRegexPath).Build(),
							builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
						},
					},
				},
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{{
						Namespace: (*gatewayv1.Namespace)(&managedGateway.Namespace),
						Name:      gatewayv1.ObjectName(managedGateway.Name),
					}},
				},
			},
		},
		WantCreateErrSubstring: wantCreateErrSubstring,
	}
}

func TestHTTPRouteValidationWebhookTraditionalRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, expressions)

	ctx := context.Background()
	namespace, gatewayClient, managedGateway, unmanagedGateway := setUpEnvForTestingHTTPRouteValidationWebhook(ctx, t)
	testCases := append(
		commonHTTPRouteValidationTestCases(managedGateway, unmanagedGateway),
		invalidRegexInPathTestCase(managedGateway, `invalid regex: '/foo[[[['`),
		// No test case for invalid regex in header, because Kong Gateway doesn't return any error in such case (it works only for expressions router).
	)
	testHTTPRouteValidationWebhook(ctx, t, namespace, gatewayClient, testCases)
}

func TestHTTPRouteValidationWebhookExpressionsRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, traditional, traditionalCompatible)

	ctx := context.Background()
	namespace, gatewayClient, managedGateway, unmanagedGateway := setUpEnvForTestingHTTPRouteValidationWebhook(ctx, t)
	testCases := append(
		commonHTTPRouteValidationTestCases(managedGateway, unmanagedGateway),
		invalidRegexInPathTestCase(managedGateway, "regex parse error:\n    ^/foo[[[[\n            ^\nerror: unclosed character class)"),
		testCaseHTTPRouteValidation{
			Name: "a httproute with invalid regex for header does not pass validation",
			Route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{"foo.com"},
					Rules: []gatewayv1.HTTPRouteRule{
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
								builder.NewHTTPRouteMatch().WithHeaderRegex("foo", "bar[[").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
							},
						},
					},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Namespace: (*gatewayv1.Namespace)(&managedGateway.Namespace),
							Name:      gatewayv1.ObjectName(managedGateway.Name),
						}},
					},
				},
			},
			WantCreateErrSubstring: "regex parse error:\n    bar[[\n        ^\nerror: unclosed character class)",
		},
	)
	testHTTPRouteValidationWebhook(ctx, t, namespace, gatewayClient, testCases)
}

// setUpEnvForTestingHTTPRouteValidationWebhook sets up the environment for testing HTTPRoute validation webhook,
// it sets it only for objects applied to namespace specified as argument.
func setUpEnvForTestingHTTPRouteValidationWebhook(ctx context.Context, t *testing.T) (
	namespace string,
	gatewayClient *gatewayclient.Clientset,
	managedGateway *gatewayv1.Gateway,
	unmanagedGateway *gatewayv1.Gateway,
) {
	ns, cleaner := helpers.Setup(ctx, t, env)
	namespace = ns.Name
	const webhookName = "kong-validations-gateway"
	ensureAdmissionRegistration(
		ctx,
		t,
		namespace,
		webhookName,
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"gateway.networking.k8s.io"},
					APIVersions: []string{"v1beta1"},
					Resources:   []string{"httproutes"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)

	t.Log("creating a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("creating a managed gateway")
	managedGateway, err = DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName, func(g *gatewayv1.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(managedGateway)
	t.Logf("created managed gateway: %q", managedGateway.Name)

	t.Logf("creating an unmanaged gatewayclass")
	unmanagedGatewayClass, err := DeployGatewayClass(ctx, gatewayClient, uuid.NewString(), func(gc *gatewayv1.GatewayClass) {
		gc.Spec.ControllerName = unsupportedControllerName
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGatewayClass)
	t.Logf("created unmanaged gatewayclass: %q", unmanagedGatewayClass.Name)

	t.Log("creating an unmanaged gateway")
	unmanagedGateway, err = DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClass.Name, func(g *gatewayv1.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGateway)
	t.Logf("created unmanaged gateway: %q", unmanagedGateway.Name)

	t.Log("waiting for webhook service to be connective")
	ensureWebhookServiceIsConnective(ctx, t, webhookName)

	return namespace, gatewayClient, managedGateway, unmanagedGateway
}

// testHTTPRouteValidationWebhook tries to create the given HTTPRoutes (passed in testCaseHTTPRouteValidation) and asserts expected results.
func testHTTPRouteValidationWebhook(
	ctx context.Context, t *testing.T, namespace string, gatewayClient *gatewayclient.Clientset, testCases []testCaseHTTPRouteValidation,
) {
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			_, err := gatewayClient.GatewayV1beta1().HTTPRoutes(namespace).Create(ctx, tC.Route, metav1.CreateOptions{})
			if tC.WantCreateErrSubstring != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tC.WantCreateErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
