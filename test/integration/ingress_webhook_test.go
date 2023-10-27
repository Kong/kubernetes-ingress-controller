//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

type testCaseIngressValidation struct {
	Name                   string
	Ingress                *netv1.Ingress
	WantCreateErrSubstring string
}

// commonIngressValidationTestCases returns a list of test cases for validating Ingress that are common
// to both traditional and expressions routers (in case of an expected error the same message is returned).
func commonIngressValidationTestCases() []testCaseIngressValidation {
	return []testCaseIngressValidation{
		{
			Name: "a valid ingress passes validation",
			Ingress: builder.NewIngress(uuid.NewString(), consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/foo"),
			).Build(),
		},
		{
			Name: "an invalid ingress passes validation when Ingress class is not set to KIC's (it's not ours)",
			Ingress: builder.NewIngress(uuid.NewString(), "third-party-ingress-class").WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/foo", "/~/foo[[["),
			).Build(),
		},
		{
			Name: "an invalid ingress passes validation when Ingress class is not set to KIC's (it's not ours), usage of legacy annotation",
			Ingress: builder.NewIngress(uuid.NewString(), "").WithLegacyClassAnnotation("third-party-ingress-class").WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/foo", "/~/foo[[["),
			).Build(),
		},
		{
			Name: "valid Ingress with multiple hosts, paths (with valid regex expressions) passes validation",
			Ingress: builder.NewIngress(uuid.NewString(), "").WithLegacyClassAnnotation("third-party-ingress-class").WithRules(
				constructIngressRuleWithPathsImplSpecific("foo.com", "/foo", "/bar[1-9]"),
				constructIngressRuleWithPathsImplSpecific("bar.com", "/baz"),
				constructIngressRuleWithPathsImplSpecific("", "/test", "/~/foo[1-9]"),
			).Build(),
		},
		{
			Name: "fail when path in Ingress does not start with '/' (K8s builtin Ingress validation)",
			Ingress: builder.NewIngress(uuid.NewString(), consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "~/foo[1-9]", "/bar"),
			).Build(),
			WantCreateErrSubstring: "Invalid value: \"~/foo[1-9]\": must be an absolute path",
		},
	}
}

// invalidRegexInIngressPathTestCase returns a test case for a Ingress with an invalid regex in the path,
// in the format that is common for both traditional and expressions routers. Error message is different
// for router flavors, thus it has passed by caller.
func invalidRegexInIngressPathTestCase(wantCreateErrSubstring string) testCaseIngressValidation {
	return testCaseIngressValidation{
		Name: "valid path format with invalid regex expression fails validation",
		Ingress: builder.NewIngress(uuid.NewString(), consts.IngressClass).WithRules(
			constructIngressRuleWithPathsImplSpecific("", "/bar", "/~/baz[1-9]"),
			constructIngressRuleWithPathsImplSpecific("", "/~/foo[[["),
		).Build(),
		WantCreateErrSubstring: wantCreateErrSubstring,
	}
}

func TestIngressValidationWebhookTraditionalRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, expressions)
	RunWhenKongVersion(t, fmt.Sprintf(">=%s", versions.ExplicitRegexPathVersionCutoff))

	ctx := context.Background()
	namespace := setUpEnvForTestingIngressValidationWebhook(ctx, t)
	testCases := append(
		commonIngressValidationTestCases(),
		invalidRegexInIngressPathTestCase(`invalid regex: '/foo[[['`),
		testCaseIngressValidation{
			Name: "path should start with '/' or '~/' (regex path) (Kong Gateway requirement for non-expressions router)",
			Ingress: builder.NewIngress(uuid.NewString(), "").WithLegacyClassAnnotation(consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/bar", "/~foo[1-9]"),
			).Build(),
			WantCreateErrSubstring: `should start with: / (fixed path) or ~/ (regex path)`,
		},
	)
	testIngressValidationWebhook(ctx, t, namespace, testCases)
}

func TestIngressValidationWebhookTraditionalRouterBeforeRequiringExplicitRegexPath(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, expressions)
	RunWhenKongVersion(t, fmt.Sprintf("<%s", versions.ExplicitRegexPathVersionCutoff))

	ctx := context.Background()
	namespace := setUpEnvForTestingIngressValidationWebhook(ctx, t)
	testCases := append(
		commonIngressValidationTestCases(),
		invalidRegexInIngressPathTestCase("should start with: /"),
		testCaseIngressValidation{
			Name: "path should start with / without any explicit regex prefixes",
			Ingress: builder.NewIngress(uuid.NewString(), "").WithLegacyClassAnnotation(consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/bar", "/~foo[1-9]"),
			).Build(),
			WantCreateErrSubstring: "should start with: /",
		},
		testCaseIngressValidation{
			Name: "path with invalid regex should fail validation",
			Ingress: builder.NewIngress(uuid.NewString(), "").WithLegacyClassAnnotation(consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/foo[[["),
			).Build(),
			WantCreateErrSubstring: `invalid regex: '/foo[[['`,
		},
	)
	testIngressValidationWebhook(ctx, t, namespace, testCases)
}

func TestIngressValidationWebhookExpressionsRouter_Reproduce(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, traditional, traditionalCompatible)
	ctx := context.Background()
	namespace := setUpEnvForTestingIngressValidationWebhook(ctx, t)

	testCases := []testCaseIngressValidation{
		{
			Name: "valid regex path passes validation",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ingress-prefix",
					Annotations: map[string]string{
						"konghq.com/strip-path": "true",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     "/",
											PathType: lo.ToPtr(netv1.PathTypePrefix),
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "httpbin-deployment",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	testIngressValidationWebhook(ctx, t, namespace, testCases)
}

func TestIngressValidationWebhookExpressionsRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, traditional, traditionalCompatible)

	ctx := context.Background()
	namespace := setUpEnvForTestingIngressValidationWebhook(ctx, t)
	testCases := append(
		commonIngressValidationTestCases(),
		invalidRegexInIngressPathTestCase("regex parse error:\n    ^/foo[[[\n           ^\nerror: unclosed character class"),
		testCaseIngressValidation{
			Name: "valid regex path passes validation",
			Ingress: builder.NewIngress(uuid.NewString(), consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/bar", "/~baz[1-9]"),
			).Build(),
		},
		testCaseIngressValidation{
			Name: "invalid regex path fails validation",
			Ingress: builder.NewIngress(uuid.NewString(), consts.IngressClass).WithRules(
				constructIngressRuleWithPathsImplSpecific("", "/bar", "/~baz[1-9]"),
				constructIngressRuleWithPathsImplSpecific("", "/~foo[[["),
			).Build(),
			WantCreateErrSubstring: "regex parse error:\n    ^foo[[[\n          ^\nerror: unclosed character class",
		},
	)
	testIngressValidationWebhook(ctx, t, namespace, testCases)
}

// setUpEnvForTestingIngressValidationWebhook sets up the environment for testing Ingress validation webhook,
// it sets it only for objects applied to namespace specified as argument.
func setUpEnvForTestingIngressValidationWebhook(ctx context.Context, t *testing.T) (namespace string) {
	ns, _ := helpers.Setup(ctx, t, env)
	namespace = ns.Name
	const webhookName = "kong-validations-ingress"
	ensureAdmissionRegistration(
		ctx,
		t,
		namespace,
		webhookName,
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"networking.k8s.io"},
					APIVersions: []string{"v1"},
					Resources:   []string{"ingresses"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	ensureWebhookServiceIsConnective(ctx, t, webhookName)
	return namespace
}

// testIngressValidationWebhook tries to create the given Ingress (passed in testCaseIngressValidation) and asserts expected results.
func testIngressValidationWebhook(
	ctx context.Context, t *testing.T, namespace string, testCases []testCaseIngressValidation,
) {
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			_, err := env.Cluster().Client().NetworkingV1().Ingresses(namespace).Create(ctx, tC.Ingress, metav1.CreateOptions{})
			if tC.WantCreateErrSubstring != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tC.WantCreateErrSubstring)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func constructIngressRuleWithPathsImplSpecific(host string, paths ...string) netv1.IngressRule {
	var pathsToSet []netv1.HTTPIngressPath
	for _, path := range paths {
		pathsToSet = append(
			pathsToSet,
			netv1.HTTPIngressPath{
				Path:     path,
				PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
				Backend: netv1.IngressBackend{
					Service: &netv1.IngressServiceBackend{
						Name: "foo",
						Port: netv1.ServiceBackendPort{
							Number: 80,
						},
					},
				},
			},
		)
	}
	return netv1.IngressRule{
		Host: host,
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: pathsToSet,
			},
		},
	}
}
