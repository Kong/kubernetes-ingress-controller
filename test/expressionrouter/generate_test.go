//go:build expression_router_tests
// +build expression_router_tests

package expressionrouter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestExpressionRouterGenerateRoutes(t *testing.T) {
	httpClient := helpers.DefaultHTTPClient()

	ip, port := exposeKongAdminService(ctx, t, env, consts.ControllerNamespace, "ingress-controller-kong-admin")

	t.Logf("creating a kong client to %s:%d", ip, port)
	kongAdminURL := fmt.Sprintf("http://%s:%d", ip, port)
	kongClient, err := kong.NewClient(&kongAdminURL, httpClient)
	require.NoError(t, err)

	ns, cleaner := helpers.Setup(ctx, t, env)

	testCases := []struct {
		name            string
		matcher         atc.Matcher
		matchRequests   []*http.Request
		unmatchRequests []*http.Request
	}{
		{
			name:    "exact match on path",
			matcher: atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo/", nil),
			},
		},
		{
			name: "exact match on path and host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpEqual, "a.foo.com"),
			),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.foo.com"), "foo", nil),
			},
		},
		{
			name: "exact match on path and wildcard match on host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpSuffixMatch, ".foo.com"),
			),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.foo.com"), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.bar.com"), "foo", nil),
			},
		},
		{
			name: "match on header and path",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/foo"),
				atc.NewPredicateHTTPHeader("X-Kong-Test", atc.OpEqual, "bar"),
			),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", map[string]string{
					"X-Kong-Test": "bar",
				}),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", map[string]string{
					"X-Kong-Test": "baz",
				}),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
			},
		},
	}

	proxyIP := getKongProxyIP(ctx, t, env, consts.ControllerNamespace)
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyIP))
	proxyClient := helpers.DefaultHTTPClient()
	proxyClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	t.Log("deploying HTTP container deployment to test generating expression routes")
	// TODO: use another HTTP server image that can return 200 on any path
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)
	t.Logf("wait for deployment %s/%s to be ready", ns.Name, deployment.Name)
	// wait for deployment to be ready
	require.Eventually(t, func() bool {
		deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Get(ctx, deployment.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
	}, 2*time.Minute, 5*time.Second)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Get(ctx, service.Name, metav1.GetOptions{})
	require.NoError(t, err)
	serviceIP := service.Spec.ClusterIP

	s := &kong.Service{
		Host: kong.String(serviceIP),
		Path: kong.String("/"),
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := &kong.Route{
				StripPath: kong.Bool(true),
			}
			atc.ApplyExpression(r, tc.matcher, 1)
			req, err := kongClient.NewRequest("POST", "/config", nil, marshalSingleServiceRoute(t, *s, *r))
			require.NoError(t, err)

			resp, err := kongClient.DoRAW(context.Background(), req)
			require.NoError(t, err)
			resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)

			// matched requests should access upstream service
			require.Eventually(t, func() bool {
				for _, req := range tc.matchRequests {
					resp, err := proxyClient.Do(req)
					if err != nil {
						t.Logf("error happened on getting response from kong: %v", err)
						return false
					}
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						return false
					}
				}
				return true
			}, time.Minute, 5*time.Second)

			// unmatched requests should get a 404 from Kong
			for _, req := range tc.unmatchRequests {
				resp, err := proxyClient.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			}
		})

	}
}
