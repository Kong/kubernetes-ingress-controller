//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	gatewayapi "github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestHTTPRouteConsumerGroups(t *testing.T) {
	t.Parallel()

	RunWhenKongEnterprise(t)

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	// path is the basic path used for most of the test
	path := "/test-gw-consumer-group/basic"
	// multiPath is the path used to test consumer group + route plugins
	multiPath := "/test-gw-consumer-group/multi"

	t.Log("getting a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	// Start GWAPI boilerplate base environment stolen from HTTPRouteEssentials

	t.Log("deploying a new gatewayClass")
	gatewayClassName := uuid.NewString()
	gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a new gateway")
	gatewayName := uuid.NewString()
	gateway, err := helpers.DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	deployment, service, ingress, keyauthPlugin := deployMinimalSvcWithKeyAuth(ctx, t, ns.Name, path)
	cleaner.Add(deployment)
	cleaner.Add(service)
	cleaner.Add(keyauthPlugin)

	// This borrows a lot from TestConsumerGroup, since it wants to test the same things for the GWAPI case, and reusing
	// test helpers simplifies that. deployMinimalSvcWithKeyAuth always creates an Ingress though, which we don't want.
	// Deleting it immediately is a bit silly, but simpler than refactoring the whole helper to be more modular. We should
	// consider such a refactor if we start using it more widely.
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))

	t.Log("creating plugins to attach to the route and groups")
	addedHeader := header{
		K: "X-Test-Header",
		V: "added",
	}
	// Use the same header key as plugin name.
	pluginRespTrans := configurePlugin(
		ctx, t, ns.Name, "response-transformer-advanced-1", "response-transformer-advanced", fmt.Sprintf(
			`{
				"add": {
					"headers": [
						"%s: %s"
					]
				}
			 }`,
			addedHeader.K, addedHeader.V,
		),
	)
	cleaner.Add(pluginRespTrans)

	const rateLimitValue = 100
	pluginRateLimit := configurePlugin(
		ctx, t, ns.Name, "rate-limiting-advanced-1", "rate-limiting-advanced", fmt.Sprintf(
			`{
				"limit": [%d],
				"window_size": [100],
				"namespace": "test",
				"sync_rate": -1,
				"window_type": "fixed"
			 }`,
			rateLimitValue,
		),
	)
	cleaner.Add(pluginRateLimit)

	t.Logf("creating an httproute to access deployment %s via kong", deployment.Name)
	httpPort := gatewayapi.PortNumber(service.Spec.Ports[0].Port)
	pathMatchPrefix := gatewayapi.PathMatchPathPrefix
	httpRoute := &gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				annotations.AnnotationPrefix + annotations.PluginsKey:   strings.Join([]string{keyauthPlugin.Name}, ","),
			},
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name: gatewayapi.ObjectName(gateway.Name),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				Matches: []gatewayapi.HTTPRouteMatch{
					{
						Path: &gatewayapi.HTTPPathMatch{
							Type:  &pathMatchPrefix,
							Value: kong.String(path),
						},
					},
				},
				BackendRefs: []gatewayapi.HTTPBackendRef{{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service.Name),
							Port: &httpPort,
							Kind: util.StringToGatewayAPIKindPtr("Service"),
						},
					},
				}},
			}},
		},
	}
	httpRoute, err = gatewayClient.GatewayV1().HTTPRoutes(ns.Name).Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	// Currently Gateway supports only one plugin per consumer group, read more
	// https://konghq.atlassian.net/browse/FTI-5282, but it silently accepts many
	// but it runs only one of them. So we have to be careful.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4472 tracks extending
	// this test once the limitation on Kong side is fixed.
	addHeaderGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-1", pluginRespTrans.Name,
	)
	cleaner.Add(addHeaderGroup)
	rateLimitGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-2", pluginRateLimit.Name,
	)
	cleaner.Add(rateLimitGroup)
	// 3 has consumers but no plugins
	nothingGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-3",
	)
	cleaner.Add(nothingGroup)
	addHeaderRouteGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-4", pluginRespTrans.Name,
	)
	cleaner.Add(addHeaderRouteGroup)

	rateLimitHeader := header{
		K: "RateLimit-Limit",
		V: fmt.Sprintf("%d", rateLimitValue),
	}
	consumers := [...]struct {
		Name                string
		ConsumerGroups      []string
		ExpectedHeaders     []header
		ForbiddenHeaderKeys []string
	}{
		{
			Name:                "test-consumer-1",
			ConsumerGroups:      []string{addHeaderGroup.Name},
			ExpectedHeaders:     []header{addedHeader},
			ForbiddenHeaderKeys: []string{rateLimitHeader.K},
		},
		{
			Name:                "test-consumer-2",
			ConsumerGroups:      []string{rateLimitGroup.Name},
			ExpectedHeaders:     []header{rateLimitHeader},
			ForbiddenHeaderKeys: []string{addedHeader.K},
		},
		{
			Name:                "test-consumer-3",
			ConsumerGroups:      nil,
			ExpectedHeaders:     nil,
			ForbiddenHeaderKeys: []string{addedHeader.K, rateLimitHeader.K},
		},
	}
	t.Log("creating consumers to be created")
	for _, consumer := range consumers {
		c, s := configureConsumerWithAPIKey(
			ctx, t, ns.Name, consumer.Name, consumer.ConsumerGroups...,
		)
		cleaner.Add(c)
		cleaner.Add(s)
	}
	t.Log("checking if consumer has plugin configured correctly based on consumer group membership")
	for _, consumer := range consumers {
		require.Eventually(t, func() bool {
			req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, path, map[string]string{
				"apikey": consumer.Name,
			})
			resp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(req)
			if err != nil {
				t.Logf("WARNING: consumer %q failed to make a request: %v", consumer.Name, err)
				return false
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Logf("WARNING: consumer %q expected status code %d, got %d", consumer.Name, http.StatusOK, resp.StatusCode)
				return false
			}
			for _, hk := range consumer.ForbiddenHeaderKeys {
				if hv := resp.Header.Get(hk); hv != "" {
					t.Logf("WARNING: consumer %q expected header %q to be empty, got %q", consumer.Name, hk, hv)
					return false
				}
			}
			for _, eh := range consumer.ExpectedHeaders {
				if hv := resp.Header.Get(eh.K); hv != eh.V {
					t.Logf("WARNING: consumer %q expected header %q to be %q, got %q", consumer.Name, eh.K, eh.V, hv)
					return false
				}
			}
			return true
		}, ingressWait, waitTick)
	}

	t.Log("checking plugins attached to a consumer group and route only apply when request matches both")
	four, fourSecret := configureConsumerWithAPIKey(ctx, t, ns.Name, "test-consumer-4", "test-consumer-group-4")
	cleaner.Add(four)
	cleaner.Add(fourSecret)

	multihttpRoute := &gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				annotations.AnnotationPrefix + annotations.PluginsKey:   strings.Join([]string{keyauthPlugin.Name, pluginRespTrans.Name}, ","),
			},
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name: gatewayapi.ObjectName(gateway.Name),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				Matches: []gatewayapi.HTTPRouteMatch{
					{
						Path: &gatewayapi.HTTPPathMatch{
							Type:  &pathMatchPrefix,
							Value: kong.String(multiPath),
						},
					},
				},
				BackendRefs: []gatewayapi.HTTPBackendRef{{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service.Name),
							Port: &httpPort,
							Kind: util.StringToGatewayAPIKindPtr("Service"),
						},
					},
				}},
			}},
		},
	}
	multihttpRoute, err = gatewayClient.GatewayV1().HTTPRoutes(ns.Name).Create(ctx, multihttpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(multihttpRoute)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		// this should see the header, it uses a consumer in the group on the associated route
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, multiPath, map[string]string{
			"apikey": four.Name,
		})
		resp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(req)
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()
		if !assert.Equal(c, resp.StatusCode, http.StatusOK) {
			return
		}
		hv := resp.Header.Get(addedHeader.K)
		if !assert.Equal(c, addedHeader.V, hv) {
			return
		}

		// this should not see the header, it uses a consumer in the group on another route
		clear := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, path, map[string]string{
			"apikey": four.Name,
		})
		clearResp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(clear)
		if !assert.NoError(c, err) {
			return
		}
		defer clearResp.Body.Close()
		if !assert.Equal(c, clearResp.StatusCode, http.StatusOK) {
			return
		}
		hv = clearResp.Header.Get(addedHeader.K)
		if !assert.NotEqual(c, addedHeader.V, hv) {
			return
		}

		// this should not see the header, it uses a consumer outside the group on the associated route
		empty := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, multiPath, map[string]string{
			"apikey": "test-consumer-3",
		})
		emptyResp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(empty)
		if !assert.NoError(c, err) {
			return
		}
		defer emptyResp.Body.Close()
		if !assert.Equal(c, emptyResp.StatusCode, http.StatusOK) {
			return
		}
		hv = emptyResp.Header.Get(addedHeader.K)
		if !assert.NotEqual(c, addedHeader.V, hv) {
			return
		}
	}, ingressWait, waitTick)
}

func TestConsumerGroup(t *testing.T) {
	t.Parallel()

	RunWhenKongEnterprise(t)

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	// path is the basic path used for most of the test
	path := "/test-consumer-group/basic"
	// multiPath is the path used to test consumer group + route plugins
	multiPath := "/test-consumer-group/multi"

	deployment, service, ingress, keyauthPlugin := deployMinimalSvcWithKeyAuth(ctx, t, ns.Name, path)
	cleaner.Add(deployment)
	cleaner.Add(service)
	cleaner.Add(ingress)
	cleaner.Add(keyauthPlugin)

	addedHeader := header{
		K: "X-Test-Header",
		V: "added",
	}
	// Use the same header key as plugin name.
	pluginRespTrans := configurePlugin(
		ctx, t, ns.Name, "response-transformer-advanced-1", "response-transformer-advanced", fmt.Sprintf(
			`{
				"add": {
					"headers": [
						"%s: %s"
					]
				}
			 }`,
			addedHeader.K, addedHeader.V,
		),
	)
	cleaner.Add(pluginRespTrans)

	const rateLimitValue = 100
	pluginRateLimit := configurePlugin(
		ctx, t, ns.Name, "rate-limiting-advanced-1", "rate-limiting-advanced", fmt.Sprintf(
			`{
				"limit": [%d],
				"window_size": [100],
				"namespace": "test",
				"sync_rate": -1,
				"window_type": "fixed"
			 }`,
			rateLimitValue,
		),
	)
	cleaner.Add(pluginRateLimit)

	// Currently Gateway supports only one plugin per consumer group, read more
	// https://konghq.atlassian.net/browse/FTI-5282, but it silently accepts many
	// but it runs only one of them. So we have to be careful.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4472 tracks extending
	// this test once the limitation on Kong side is fixed.
	addHeaderGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-1", pluginRespTrans.Name,
	)
	cleaner.Add(addHeaderGroup)
	rateLimitGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-2", pluginRateLimit.Name,
	)
	cleaner.Add(rateLimitGroup)
	// 3 has consumers but no plugins
	nothingGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-3",
	)
	cleaner.Add(nothingGroup)
	addHeaderRouteGroup := configureConsumerGroupWithPlugins(
		ctx, t, ns.Name, "test-consumer-group-4", pluginRespTrans.Name,
	)
	cleaner.Add(addHeaderRouteGroup)

	rateLimitHeader := header{
		K: "RateLimit-Limit",
		V: fmt.Sprintf("%d", rateLimitValue),
	}
	consumers := [...]struct {
		Name                string
		ConsumerGroups      []string
		ExpectedHeaders     []header
		ForbiddenHeaderKeys []string
	}{
		{
			Name:                "test-consumer-1",
			ConsumerGroups:      []string{addHeaderGroup.Name},
			ExpectedHeaders:     []header{addedHeader},
			ForbiddenHeaderKeys: []string{rateLimitHeader.K},
		},
		{
			Name:                "test-consumer-2",
			ConsumerGroups:      []string{rateLimitGroup.Name},
			ExpectedHeaders:     []header{rateLimitHeader},
			ForbiddenHeaderKeys: []string{addedHeader.K},
		},
		{
			Name:                "test-consumer-3",
			ConsumerGroups:      nil,
			ExpectedHeaders:     nil,
			ForbiddenHeaderKeys: []string{addedHeader.K, rateLimitHeader.K},
		},
	}
	t.Log("creating consumers to be created")
	for _, consumer := range consumers {
		c, s := configureConsumerWithAPIKey(
			ctx, t, ns.Name, consumer.Name, consumer.ConsumerGroups...,
		)
		cleaner.Add(c)
		cleaner.Add(s)
	}
	t.Log("checking if consumer has plugin configured correctly based on consumer group membership")
	for _, consumer := range consumers {
		require.Eventually(t, func() bool {
			req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, path, map[string]string{
				"apikey": consumer.Name,
			})
			resp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(req)
			if err != nil {
				t.Logf("WARNING: consumer %q failed to make a request: %v", consumer.Name, err)
				return false
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Logf("WARNING: consumer %q expected status code %d, got %d", consumer.Name, http.StatusOK, resp.StatusCode)
				return false
			}
			for _, hk := range consumer.ForbiddenHeaderKeys {
				if hv := resp.Header.Get(hk); hv != "" {
					t.Logf("WARNING: consumer %q expected header %q to be empty, got %q", consumer.Name, hk, hv)
					return false
				}
			}
			for _, eh := range consumer.ExpectedHeaders {
				if hv := resp.Header.Get(eh.K); hv != eh.V {
					t.Logf("WARNING: consumer %q expected header %q to be %q, got %q", consumer.Name, eh.K, eh.V, hv)
					return false
				}
			}
			return true
		}, ingressWait, waitTick)
	}

	t.Log("checking plugins attached to a consumer group and route only apply when request matches both")
	four, fourSecret := configureConsumerWithAPIKey(ctx, t, ns.Name, "test-consumer-4", "test-consumer-group-4")
	cleaner.Add(four)
	cleaner.Add(fourSecret)

	multiIngress := generators.NewIngressForService(multiPath, map[string]string{
		annotations.AnnotationPrefix + annotations.StripPathKey: "true",
		annotations.AnnotationPrefix + annotations.PluginsKey:   strings.Join([]string{keyauthPlugin.Name, pluginRespTrans.Name}, ","),
	}, service)
	multiIngress.Spec.IngressClassName = kong.String(consts.IngressClass)
	multiIngress.Name = "multi"
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, multiIngress))
	cleaner.Add(multiIngress)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		// this should see the header, it uses a consumer in the group on the associated route
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, multiPath, map[string]string{
			"apikey": four.Name,
		})
		resp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(req)
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()
		if !assert.Equal(c, resp.StatusCode, http.StatusOK) {
			return
		}
		hv := resp.Header.Get(addedHeader.K)
		if !assert.Equal(c, addedHeader.V, hv) {
			return
		}

		// this should not see the header, it uses a consumer in the group on another route
		clear := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, path, map[string]string{
			"apikey": four.Name,
		})
		clearResp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(clear)
		if !assert.NoError(c, err) {
			return
		}
		defer clearResp.Body.Close()
		if !assert.Equal(c, clearResp.StatusCode, http.StatusOK) {
			return
		}
		hv = clearResp.Header.Get(addedHeader.K)
		if !assert.NotEqual(c, addedHeader.V, hv) {
			return
		}

		// this should not see the header, it uses a consumer outside the group on the associated route
		empty := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, multiPath, map[string]string{
			"apikey": "test-consumer-3",
		})
		emptyResp, err := helpers.DefaultHTTPClientWithProxy(proxyHTTPURL).Do(empty)
		if !assert.NoError(c, err) {
			return
		}
		defer emptyResp.Body.Close()
		if !assert.Equal(c, emptyResp.StatusCode, http.StatusOK) {
			return
		}
		hv = emptyResp.Header.Get(addedHeader.K)
		if !assert.NotEqual(c, addedHeader.V, hv) {
			return
		}
	}, ingressWait, waitTick)
}

func deployMinimalSvcWithKeyAuth(
	ctx context.Context, t *testing.T, namespace, path string,
) (*appsv1.Deployment, *corev1.Service, *netv1.Ingress, *kongv1.KongPlugin) {
	const pluginKeyAuthName = "key-auth"
	t.Logf("configuring plugin %q (to give consumers an identity)", pluginKeyAuthName)
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	pluginKeyAuth, err := c.ConfigurationV1().KongPlugins(namespace).Create(
		ctx,
		&kongv1.KongPlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name: pluginKeyAuthName,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
			PluginName: "key-auth",
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	t.Log("deploying a minimal HTTP container")
	deployment := generators.NewDeploymentForContainer(
		generators.NewContainer("echo", test.EchoImage, test.EchoHTTPPort),
	)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %q via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	_, err = env.Cluster().Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an ingress for service %q with plugin %q attached", service.Name, pluginKeyAuthName)
	ingress := generators.NewIngressForService(path, map[string]string{
		annotations.AnnotationPrefix + annotations.StripPathKey: "true",
		annotations.AnnotationPrefix + annotations.PluginsKey:   pluginKeyAuthName,
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), namespace, ingress))
	return deployment, service, ingress, pluginKeyAuth
}

func configurePlugin(
	ctx context.Context, t *testing.T, namespace string, name string, pluginName string, cfg string,
) *kongv1.KongPlugin {
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	t.Logf("configuring plugin %q (%q)", name, pluginName)
	pluginRespTrans, err := c.ConfigurationV1().KongPlugins(namespace).Create(
		ctx,
		&kongv1.KongPlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
			PluginName: pluginName,
			Config: apiextensionsv1.JSON{
				Raw: []byte(cfg),
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	return pluginRespTrans
}

func configureConsumerGroupWithPlugins(
	ctx context.Context, t *testing.T, namespace string, name string, pluginName ...string,
) *kongv1beta1.KongConsumerGroup {
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	a := map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
	}
	if plugins := strings.Join(pluginName, ","); plugins != "" {
		a[annotations.AnnotationPrefix+annotations.PluginsKey] = plugins
		t.Logf("configuring consumer group %q with attached plugins: %s", name, plugins)
	} else {
		t.Logf("configuring consumer group %q with no plugins attached", name)
	}
	cg, err := c.ConfigurationV1beta1().KongConsumerGroups(namespace).Create(
		ctx,
		&kongv1beta1.KongConsumerGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Annotations: a,
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	return cg
}

// configureConsumerWithAPIKey creates a consumer with a key-auth credential set to the consumer's name.
// Assign consumer to specified consumer groups.
func configureConsumerWithAPIKey(
	ctx context.Context, t *testing.T, namespace string, name string, consumerGroups ...string,
) (*kongv1.KongConsumer, *corev1.Secret) {
	t.Logf(
		"creating a consumer: %q with api-key and consumer groups: %s configured",
		name, strings.Join(consumerGroups, ","),
	)
	secret, err := env.Cluster().Client().CoreV1().Secrets(namespace).Create(
		ctx,
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			StringData: map[string]string{
				"key": name,
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	consumer, err := c.ConfigurationV1().KongConsumers(namespace).Create(
		ctx,
		&kongv1.KongConsumer{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
			Username:       name,
			ConsumerGroups: consumerGroups,
			Credentials:    []string{name},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	return consumer, secret
}

type header struct {
	K string
	V string
}
