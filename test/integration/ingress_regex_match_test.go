//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

func TestIngressRegexMatchPath(t *testing.T) {
	if !versions.GetKongVersion().MajorOnly().GTE(versions.ExplicitRegexPathVersionCutoff) {
		t.Skip("regex prefixes are only relevant for Kong 3.0+")
	}
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	pathRegexPrefix := "/~"
	pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
	testCases := []struct {
		pathRegex     string
		matchPaths    []string
		notMatchPaths []string
		description   string
	}{
		{
			pathRegex:     "/test-ingress-regex-match-path/\\*\\**",
			matchPaths:    []string{"/test-ingress-regex-match-path/*", "/test-ingress-regex-match-path/***"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/", "/test-ingress-regex-match-path/a*"},
			description:   "match arbitrary number of *s (at lease 1) after /",
		},
		{
			pathRegex:     "/test-ingress-regex-match-path/[\\\\|\\.]a$",
			matchPaths:    []string{"/test-ingress-regex-match-path/\\a"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/", "/test-ingress-regex-match-path/a*"},
			description:   "match path '/\\a' or '/.a'",
		},
		{
			pathRegex:     "/test-ingress-regex-match-path/\\d{3}/?$",
			matchPaths:    []string{"/test-ingress-regex-match-path/123", "/test-ingress-regex-match-path/456/"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/1234", "/test-ingress-regex-match-path/56/"},
			description:   "match 3 digits, and maybe a trailing /",
		},
		{
			pathRegex:     "/test-ingress-regex-match-path/\\d$",
			matchPaths:    []string{"/test-ingress-regex-match-path/1", "/test-ingress-regex-match-path/9"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/12", "/test-ingress-regex-match-path/3/"},
			description:   "match 1 digit, and no trailing /",
		},
		{
			pathRegex:     "/test-ingress-regex-match-path/\\\\d$",
			matchPaths:    []string{"/test-ingress-regex-match-path/\\d"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/1", "/test-ingress-regex-match-path/d"},
			description:   "match literal '/\\d'",
		},
		{
			pathRegex:     "/test-ingress-regex-match-path/[ab][cd]$",
			matchPaths:    []string{"/test-ingress-regex-match-path/ac", "/test-ingress-regex-match-path/bd"},
			notMatchPaths: []string{"/test-ingress-regex-match-path/ab", "/test-ingress-regex-match-path/ac/"},
			description:   "match charset: first from {a,b}, second from {c,d}",
		},
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	for i, tc := range testCases {
		index := i
		tc := tc
		t.Run(fmt.Sprintf("case-%d: %s", index, tc.pathRegex), func(t *testing.T) {
			t.Log("create an ingress")
			ingress := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "ingress-regex-path-" + strconv.Itoa(index),
					Annotations: map[string]string{
						"konghq.com/strip-path": "true",
					},
				},
				Spec: netv1.IngressSpec{
					IngressClassName: &ingressClass,
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     pathRegexPrefix + tc.pathRegex,
											PathType: &pathTypeImplementationSpecific,
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: service.Name,
													Port: netv1.ServiceBackendPort{Number: int32(80)},
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
			_, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Create(ctx, ingress, metav1.CreateOptions{})
			require.NoError(t, err)
			defer func() {
				err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Delete(ctx, ingress.Name, metav1.DeleteOptions{})
				require.NoError(t, err)
			}()

			t.Log("testing paths expected to match")
			notMatchedPaths := []string{}
			require.Eventually(t, func() bool {
				notMatchedPaths = []string{}
				for _, path := range tc.matchPaths {
					resp, err := httpc.Get(fmt.Sprintf("%s%s", proxyURL, path))
					if err != nil {
						t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
						notMatchedPaths = append(notMatchedPaths, path)
						return false
					}
					defer resp.Body.Close()
					// returns false if one path is not matched.
					if resp.StatusCode == http.StatusOK {
						b := new(bytes.Buffer)
						n, err := b.ReadFrom(resp.Body)
						require.NoError(t, err)
						require.True(t, n > 0)
						if !strings.Contains(b.String(), "<title>httpbin.org</title>") {
							notMatchedPaths = append(notMatchedPaths, path)
							return false
						}
					} else {
						notMatchedPaths = append(notMatchedPaths, path)
						return false
					}
				}
				// returns true if all testing paths matched.
				return true
			}, ingressWait, waitTick,
				fmt.Sprintf("paths %v not matched %s", notMatchedPaths, tc.description))

			t.Log("testing paths expected not to match")
			for _, path := range tc.notMatchPaths {
				resp, err := httpc.Get(fmt.Sprintf("%s%s", proxyURL, path))
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equalf(t, http.StatusNotFound, resp.StatusCode, "should not match path %s: %s", path, tc.description)
			}
		})
	}
}

func TestIngressRegexMatchHeader(t *testing.T) {
	if !versions.GetKongVersion().MajorOnly().GTE(versions.ExplicitRegexPathVersionCutoff) {
		t.Skip("regex prefixes are only relevant for Kong 3.0+")
	}
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	headerRegexPrefix := "~*"
	matchHeaderKey := "X-Kic-Test-Match"
	pathTypePrefix := netv1.PathTypePrefix
	testCases := []struct {
		headerRegex     string
		matchHeaders    []string
		notMatchHeaders []string
	}{
		{
			headerRegex:     "^[abc]+$",
			matchHeaders:    []string{"a", "aaa", "abc"},
			notMatchHeaders: []string{"", "abcd"},
		},
		{
			headerRegex:     "^kong\\.",
			matchHeaders:    []string{"kong.", "kong.abc", "kong.foo.bar"},
			notMatchHeaders: []string{"kong", "akong."},
		},
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	for i, tc := range testCases {
		index := i
		tc := tc
		t.Run(fmt.Sprintf("case-%d: %s", index, tc.headerRegex), func(t *testing.T) {
			t.Log("create an ingress")
			ingress := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "ingress-regex-header-" + strconv.Itoa(index),
					Annotations: map[string]string{
						"konghq.com/strip-path":                                 "true",
						"konghq.com/headers." + strings.ToLower(matchHeaderKey): headerRegexPrefix + tc.headerRegex,
					},
				},
				Spec: netv1.IngressSpec{
					IngressClassName: &ingressClass,
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     "/",
											PathType: &pathTypePrefix,
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: service.Name,
													Port: netv1.ServiceBackendPort{Number: int32(80)},
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
			_, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Create(ctx, ingress, metav1.CreateOptions{})
			require.NoError(t, err)
			defer func() {
				err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Delete(ctx, ingress.Name, metav1.DeleteOptions{})
				require.NoError(t, err)
			}()

			t.Log("testing headers expected to match")
			require.Eventually(t, func() bool {
				for _, header := range tc.matchHeaders {
					req, err := http.NewRequest("GET", proxyURL.String(), nil)
					req.Header.Add(matchHeaderKey, header)
					require.NoError(t, err)
					resp, err := httpc.Do(req)
					if err != nil {
						t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
						return false
					}
					defer resp.Body.Close()
					// returns false if one of test requests is not matched.
					if resp.StatusCode == http.StatusOK {
						b := new(bytes.Buffer)
						n, err := b.ReadFrom(resp.Body)
						require.NoError(t, err)
						require.True(t, n > 0)
						if !strings.Contains(b.String(), "<title>httpbin.org</title>") {
							return false
						}
					} else {
						return false
					}
				}
				// returns true if all testing requests matched.
				return true
			}, ingressWait, waitTick)

			t.Log("testing headers expected not to match")
			for _, header := range tc.notMatchHeaders {
				req, err := http.NewRequest("GET", proxyURL.String(), nil)
				req.Header.Add(matchHeaderKey, header)
				require.NoError(t, err)
				resp, err := httpc.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equalf(t, http.StatusNotFound, resp.StatusCode, "should not match host %s", header)
			}
		})
	}
}
