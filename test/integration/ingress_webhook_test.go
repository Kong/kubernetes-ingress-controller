//go:build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
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
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr(consts.IngressClass),
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/foo"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "an invalid ingress passes validation when Ingress class is not set to KIC's (it's not ours)",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr("third-party-ingress-class"),
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/foo"),
										constructIngressPathImplSpecific("/~/foo[[["),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "valid Ingress with multiple hosts, paths (with valid regex expressions) passes validation",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "foo.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/foo"),
										constructIngressPathImplSpecific("/bar[1-9]"),
									},
								},
							},
						},
						{
							Host: "bar.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/baz"),
									},
								},
							},
						},
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/test"),
										constructIngressPathImplSpecific("/~/foo[1-9]"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "fail when path in Ingress does not start with '/' (K8s builtin Ingress validation)",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr(consts.IngressClass),
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("~/foo[1-9]"),
										constructIngressPathImplSpecific("/bar"),
									},
								},
							},
						},
					},
				},
			},
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
		Ingress: &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.NewString(),
			},
			Spec: netv1.IngressSpec{
				IngressClassName: lo.ToPtr(consts.IngressClass),
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									constructIngressPathImplSpecific("/bar"),
									constructIngressPathImplSpecific("/~/baz[1-9]"),
								},
							},
						},
					},
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									constructIngressPathImplSpecific(`/~/foo[[[`),
								},
							},
						},
					},
				},
			},
		},
		WantCreateErrSubstring: wantCreateErrSubstring,
	}
}

func TestIngressValidationWebhookTraditionalRouter(t *testing.T) {
	skipTestForNonKindCluster(t)
	skipTestForRouterFlavors(t, expressions)

	ctx := context.Background()
	namespace := setUpEnvForTestingIngressValidationWebhook(ctx, t)
	testCases := append(
		commonIngressValidationTestCases(),
		invalidRegexInIngressPathTestCase(`invalid regex: '/foo[[['`),
		testCaseIngressValidation{
			Name: "path should start with '/' or '~/' (regex path) (Kong Gateway requirement for non-expressions router)",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/bar"),
										constructIngressPathImplSpecific("/~foo[1-9]"),
									},
								},
							},
						},
					},
				},
			},
			WantCreateErrSubstring: `should start with: / (fixed path) or ~/ (regex path)`,
		},
	)
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
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr(consts.IngressClass),
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/bar"),
										constructIngressPathImplSpecific("/~baz[1-9]"),
									},
								},
							},
						},
					},
				},
			},
		},
		testCaseIngressValidation{
			Name: "invalid regex path fails validation",
			Ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr(consts.IngressClass),
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/bar"),
										constructIngressPathImplSpecific("/~baz[1-9]"),
									},
								},
							},
						},
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										constructIngressPathImplSpecific("/~foo[[["),
									},
								},
							},
						},
					},
				},
			},
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

func constructIngressPathImplSpecific(path string) netv1.HTTPIngressPath {
	return netv1.HTTPIngressPath{
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
	}
}
