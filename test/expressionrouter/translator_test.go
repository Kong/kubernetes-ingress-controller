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
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestExpressionRouterTranslateIngress(t *testing.T) {
	httpClient := helpers.DefaultHTTPClient()

	ip, port := exposeKongAdminService(ctx, t, env, consts.ControllerNamespace, "ingress-controller-kong-admin")

	t.Logf("creating a kong client to %s:%d", ip, port)
	kongAdminURL := fmt.Sprintf("http://%s:%d", ip, port)
	kongClient, err := kong.NewClient(&kongAdminURL, httpClient)
	require.NoError(t, err)

	ns, cleaner := helpers.Setup(ctx, t, env)

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

	ingressbackend := netv1.IngressBackend{
		Service: &netv1.IngressServiceBackend{
			Name: service.Name,
			Port: netv1.ServiceBackendPort{
				Number: 80,
			},
		},
	}

	pathTypeExact := netv1.PathType(netv1.PathTypeExact)

	testCases := []struct {
		name            string
		ingress         *netv1.Ingress
		matchRequests   []*http.Request
		unmatchRequests []*http.Request
	}{
		{
			name: "simple ingress",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "simple-ingress",
					Namespace: ns.Name,
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
											PathType: &pathTypeExact,
											Path:     "/foo",
											Backend:  ingressbackend,
										},
									},
								},
							},
						},
					},
				},
			},
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo/", nil),
			},
		},
		{
			name: "ingress with host match and host alias",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ingress-host-match",
					Namespace: ns.Name,
					Annotations: map[string]string{
						"konghq.com/strip-path":   "true",
						"konghq.com/host-aliases": "a.bar.com",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "*.foo.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											PathType: &pathTypeExact,
											Path:     "/foo",
											Backend:  ingressbackend,
										},
									},
								},
							},
						},
					},
				},
			},
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com"), "foo", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.foo.com"), "foo", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.bar.com"), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.bar.com"), "foo", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.bla.com"), "foo", nil),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			kongServices := translators.TranslateIngressATC(tc.ingress)
			req, err := kongClient.NewRequest("POST", "/config", nil, marshalKongStateServices(t, kongServices, serviceIP))
			require.NoError(t, err)

			resp, err := kongClient.DoRAW(context.Background(), req)
			require.NoError(t, err)
			resp.Body.Close()

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
