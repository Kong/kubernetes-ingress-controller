//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const invalidRegexPath = "/foo[[[["

type testCaseHTTPRouteValidation struct {
	Name                   string
	Route                  *gatewayapi.HTTPRoute
	WantCreateErrSubstring string
}

// commonHTTPRouteValidationTestCases returns a list of test cases for validating HTTPRoutes
// that are common to both traditional and expressions routers (the same error message is returned).
func commonHTTPRouteValidationTestCases(
	managedGateway *gatewayapi.Gateway, unmanagedGateway *gatewayapi.Gateway,
) []testCaseHTTPRouteValidation {
	return []testCaseHTTPRouteValidation{
		{
			Name: "a valid httproute linked to a managed gateway passes validation",
			Route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Namespace: (*gatewayapi.Namespace)(&managedGateway.Namespace),
							Name:      gatewayapi.ObjectName(managedGateway.Name),
						}},
					},
				},
			},
		},
		{
			Name: "a httproute linked to a non-existent gateway passes validation",
			Route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Namespace: (*gatewayapi.Namespace)(&managedGateway.Namespace),
							Name:      gatewayapi.ObjectName("fake-gateway"),
						}},
					},
				},
			},
		},
		{
			Name: "an invalid httproute will pass validation if it's not linked to a managed controller (it's not ours)",
			Route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex(invalidRegexPath).Build(),
						},
					}},
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Namespace: (*gatewayapi.Namespace)(&unmanagedGateway.Namespace),
							Name:      gatewayapi.ObjectName(unmanagedGateway.Name),
						}},
					},
				},
			},
		},
		{
			Name: "a httproute with valid regex expressions for a path and a header pass validation",
			Route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{"foo.com"},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathRegex("/path[1-8]").Build(),
								builder.NewHTTPRouteMatch().WithHeaderRegex("foo", "bar[1-8]").Build(),
							},
						},
					},
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Namespace: (*gatewayapi.Namespace)(&managedGateway.Namespace),
							Name:      gatewayapi.ObjectName(managedGateway.Name),
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
	managedGateway *gatewayapi.Gateway, wantCreateErrSubstring string,
) testCaseHTTPRouteValidation {
	return testCaseHTTPRouteValidation{
		Name: "a httproute with invalid regex for path does not pass validation",
		Route: &gatewayapi.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.NewString(),
			},
			Spec: gatewayapi.HTTPRouteSpec{
				Hostnames: []gatewayapi.Hostname{"foo.com"},
				Rules: []gatewayapi.HTTPRouteRule{
					{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
							builder.NewHTTPRouteMatch().WithPathRegex(invalidRegexPath).Build(),
							builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
						},
					},
				},
				CommonRouteSpec: gatewayapi.CommonRouteSpec{
					ParentRefs: []gatewayapi.ParentReference{{
						Namespace: (*gatewayapi.Namespace)(&managedGateway.Namespace),
						Name:      gatewayapi.ObjectName(managedGateway.Name),
					}},
				},
			},
		},
		WantCreateErrSubstring: wantCreateErrSubstring,
	}
}

func TestHTTPRouteValidationWebhookTraditionalRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(context.Background(), t, expressions)

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
	skipTestForRouterFlavors(context.Background(), t, traditional, traditionalCompatible)

	ctx := context.Background()
	namespace, gatewayClient, managedGateway, unmanagedGateway := setUpEnvForTestingHTTPRouteValidationWebhook(ctx, t)
	testCases := append(
		commonHTTPRouteValidationTestCases(managedGateway, unmanagedGateway),
		invalidRegexInPathTestCase(managedGateway, "regex parse error:\n    ^/foo[[[[\n            ^\nerror: unclosed character class)"),
		testCaseHTTPRouteValidation{
			Name: "a httproute with invalid regex for header does not pass validation",
			Route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{"foo.com"},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
								builder.NewHTTPRouteMatch().WithHeaderRegex("foo", "bar[[").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
							},
						},
					},
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Namespace: (*gatewayapi.Namespace)(&managedGateway.Namespace),
							Name:      gatewayapi.ObjectName(managedGateway.Name),
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
	managedGateway *gatewayapi.Gateway,
	unmanagedGateway *gatewayapi.Gateway,
) {
	ns, cleaner := helpers.Setup(ctx, t, env)
	namespace = ns.Name
	ensureAdmissionRegistration(ctx, t, env.Cluster().Client(), "kong-validations-gateway", ns.Name)

	t.Log("creating a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("creating a managed gateway")
	managedGateway, err = helpers.DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName, func(g *gatewayapi.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(managedGateway)
	t.Logf("created managed gateway: %q", managedGateway.Name)

	t.Logf("creating an unmanaged gatewayclass")
	unmanagedGatewayClass, err := helpers.DeployGatewayClass(ctx, gatewayClient, uuid.NewString(), func(gc *gatewayapi.GatewayClass) {
		gc.Spec.ControllerName = unsupportedControllerName
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGatewayClass)
	t.Logf("created unmanaged gatewayclass: %q", unmanagedGatewayClass.Name)

	t.Log("creating an unmanaged gateway")
	unmanagedGateway, err = helpers.DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClass.Name, func(g *gatewayapi.Gateway) {
		g.Name = uuid.NewString()
	})
	require.NoError(t, err)
	cleaner.Add(unmanagedGateway)
	t.Logf("created unmanaged gateway: %q", unmanagedGateway.Name)

	return namespace, gatewayClient, managedGateway, unmanagedGateway
}

// testHTTPRouteValidationWebhook tries to create the given HTTPRoutes (passed in testCaseHTTPRouteValidation) and asserts expected results.
func testHTTPRouteValidationWebhook(
	ctx context.Context, t *testing.T, namespace string, gatewayClient *gatewayclient.Clientset, testCases []testCaseHTTPRouteValidation,
) {
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			_, err := gatewayClient.GatewayV1().HTTPRoutes(namespace).Create(ctx, tC.Route, metav1.CreateOptions{})
			if tC.WantCreateErrSubstring != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tC.WantCreateErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
