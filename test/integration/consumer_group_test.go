//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestConsumerGroup(t *testing.T) {
	t.Parallel()

	RunWhenKongEnterprise(t)

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	d, s, i, p := deployMinimalSvcWithKeyAuth(ctx, t, ns.Name)
	cleaner.Add(d)
	cleaner.Add(s)
	cleaner.Add(i)
	cleaner.Add(p)

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
			req := helpers.MustHTTPRequest(t, http.MethodGet, proxyURL, "/", map[string]string{
				"apikey": consumer.Name,
			})
			resp, err := helpers.DefaultHTTPClient().Do(req)
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
}

func deployMinimalSvcWithKeyAuth(
	ctx context.Context, t *testing.T, namespace string,
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
	ingress := generators.NewIngressForService("/", map[string]string{
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
			},
			StringData: map[string]string{
				"key":          name,
				"kongCredType": "key-auth",
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
