//go:build integration_tests

package isolated

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	kongaddon "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

var (
	// TODO: support specifying images by env.
	oldKICImageRepo = "kong/kubernetes-ingress-controller"
	oldKICImageTag  = "3.4.4"
	newKICImageRepo = "kong/kubernetes-ingress-controller"
	newKICImageTag  = ""
)

func TestUpgradeKICWithExistingPlugins(t *testing.T) {
	const serviceName = "http-echo"
	const pluginName = "response-transformer-add-header"
	const echoPath = "/echo"
	testUUID := uuid.New()

	f := features.New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.KindIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongPlugin).
		WithSetup("Install Kong and KIC by helm", setUpKongAndKIC).
		WithSetup("Install an echo service and an ingress with a response-transformer plugin",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				kongClient, err := clientset.NewForConfig(cfg.Client().RESTConfig())
				assert.NoError(t, err)
				ctx = SetInCtxForT(ctx, t, kongClient)

				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				cluster := GetClusterFromCtx(ctx)

				t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
				container := generators.NewContainer(serviceName, test.EchoImage, test.EchoHTTPPort)
				// App go-echo sends a "Running on Pod <UUID>." immediately on connecting.
				container.Env = []corev1.EnvVar{
					{
						Name:  "POD_NAME",
						Value: testUUID.String(),
					},
				}
				deployment := generators.NewDeploymentForContainer(container)
				deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(deployment)

				t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
				service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
				service.Name = serviceName
				// Use the same port as the default TCP port from the Kong Gateway deployment
				// to the tcpecho port, as this is what will be used to route the traffic at the Gateway.
				service.Spec.Ports = []corev1.ServicePort{{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       test.EchoHTTPPort,
					TargetPort: intstr.FromInt(test.EchoHTTPPort),
				}}
				service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(service)

				plugin := &configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      pluginName,
					},
					PluginName: "response-transformer",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"add":{"headers":["Kic-Added:Test"]}}`),
					},
				}

				plugin, err = kongClient.ConfigurationV1().KongPlugins(namespace).Create(ctx, plugin, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(plugin)

				ingress := generators.NewIngressForService(
					echoPath, map[string]string{
						"konghq.com/plugins": pluginName,
					}, service)
				ingress.Spec.IngressClassName = lo.ToPtr(GetIngressClassFromCtx(ctx))
				ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(ingress)

				return ctx
			},
		).
		Assess("Verify that the ingress can be accessed and the response-transformer plugin works",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				proxyURL := GetHTTPURLFromCtx(ctx)
				assert.NotNil(t, proxyURL)

				t.Logf("sendind HTTP GET request to %s%s to verify that ingress and plugin are configured",
					proxyURL.Host, echoPath)
				helpers.EventuallyGETPath(
					t, proxyURL,
					proxyURL.Host,
					echoPath,
					nil,
					http.StatusOK,
					testUUID.String(),
					nil,
					consts.IngressWait,
					consts.WaitTick,
					func(resp *http.Response, _ string) (key string, ok bool) {
						return "header 'Kic-Added' should be added and has expected value", resp.Header.Get("Kic-Added") == "Test"
					},
				)
				return ctx
			}).
		Assess("Upgrade KIC and verify it restarts successfully", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)

			if testenv.KongEnterpriseEnabled() {
				t.Logf("cleaning up existing secrets for enterprise")
				secretClient := cluster.Client().CoreV1().Secrets(namespace)
				_, getErr := secretClient.Get(ctx, kongaddon.DefaultEnterpriseLicenseSecretName, metav1.GetOptions{})
				if getErr == nil {
					delErr := secretClient.Delete(ctx, kongaddon.DefaultEnterpriseLicenseSecretName, metav1.DeleteOptions{})
					assert.NoError(t, delErr, "failed to delete existing secret for Kong enterprise license")
				}
				_, getErr = secretClient.Get(ctx, kongaddon.DefaultEnterpriseAdminPasswordSecretName, metav1.GetOptions{})
				if getErr == nil {
					delErr := secretClient.Delete(ctx, kongaddon.DefaultEnterpriseAdminPasswordSecretName, metav1.DeleteOptions{})
					assert.NoError(t, delErr, "failed to delete existing secret for Kong enterprise admin password")
				}
				_, getErr = secretClient.Get(ctx, kongaddon.DefaultAdminGUISessionConfSecretName, metav1.GetOptions{})
				if getErr == nil {
					delErr := secretClient.Delete(ctx, kongaddon.DefaultAdminGUISessionConfSecretName, metav1.DeleteOptions{})
					assert.NoError(t, delErr, "failed to delete existing secret for Kong enterprise admin GUI session conifg")
				}
			}

			runControllerManager := true
			if newKICImageRepo != "" && newKICImageTag != "" {
				runControllerManager = false
				t.Logf("Upgrading KIC to %s:%s", newKICImageRepo, newKICImageTag)
			} else {
				t.Log("Run controller manager from code base, disabling KIC in Kong addon")
			}

			ctx = deployKongAddon(ctx, t, deployKongAddonCfg{
				// Do NOT deploy KIC when we run controller manager locally.
				deployControllerInKongAddon: !runControllerManager,
				controllerImageRepository:   newKICImageRepo,
				controllerImageTag:          newKICImageTag,
				kongProxyEnvVars:            map[string]string{},
			})

			kongAddon := GetFromCtxForT[*kongaddon.Addon](ctx, t)
			if runControllerManager {
				startControllerManager(ctx, t, featureSetupCfg{
					ingressClassName: "kong",
				}, kongAddon, nil)
			}

			return ctx
		}).
		Assess("Verify that the old configuration remains unchanged", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			assert.NotNil(t, proxyURL)

			t.Logf("sendind HTTP GET request to %s%s to verify that ingress and plugin are configured",
				proxyURL.Host, echoPath)
			helpers.EventuallyGETPath(
				t, proxyURL,
				proxyURL.Host,
				echoPath,
				nil,
				http.StatusOK,
				testUUID.String(),
				nil,
				consts.IngressWait,
				consts.WaitTick,
				func(resp *http.Response, _ string) (key string, ok bool) {
					return "header 'Kic-Added' should be added and has expected value", resp.Header.Get("Kic-Added") == "Test"
				},
			)
			return ctx
		}).
		Assess("Update the plugin and verify that the new configuration works", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			kongClient := GetFromCtxForT[*clientset.Clientset](ctx, t)
			assert.NotNil(t, kongClient)
			cluster := GetClusterFromCtx(ctx)

			namespace := GetNamespaceForT(ctx, t)

			plugin, err := kongClient.ConfigurationV1().KongPlugins(namespace).Get(ctx, pluginName, metav1.GetOptions{})
			assert.NoErrorf(t, err, "failed to get plugin %s/%s", namespace, pluginName)

			plugin.Config = apiextensionsv1.JSON{
				Raw: []byte(`{"add":{"headers":["Kic-Added:Another-Test"]}}`),
			}
			_, err = kongClient.ConfigurationV1().KongPlugins(namespace).Update(ctx, plugin, metav1.UpdateOptions{})
			assert.NoErrorf(t, err, "failed to update plugin %s/%s", namespace, pluginName)

			service, err := cluster.Client().CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
			assert.NoErrorf(t, err, "failed to get service %s/%S", namespace, serviceName)

			newEchoPath := "/echo-new"
			newIngress := generators.NewIngressForService(
				newEchoPath,
				map[string]string{
					"konghq.com/plugins": pluginName,
				}, service,
			)
			newIngress.Name = serviceName + "-new"
			newIngress.Spec.IngressClassName = lo.ToPtr("kong")
			_, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, newIngress, metav1.CreateOptions{})
			assert.NoError(t, err)

			proxyURL := GetHTTPURLFromCtx(ctx)
			assert.NotNil(t, proxyURL)

			t.Logf("sendind HTTP GET request to %s%s to verify that ingress and plugin are configured",
				proxyURL.Host, newEchoPath)
			helpers.EventuallyGETPath(
				t, proxyURL,
				proxyURL.Host,
				newEchoPath,
				nil,
				http.StatusOK,
				testUUID.String(),
				nil,
				consts.IngressWait,
				consts.WaitTick,
				func(resp *http.Response, _ string) (key string, ok bool) {
					t.Logf("header Kic-Added: %s", resp.Header.Get("Kic-Added"))
					return "header 'Kic-Added' should be added and has expected value", resp.Header.Get("Kic-Added") == "Another-Test"
				},
			)

			return ctx
		},
		).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func setUpKongAndKIC(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	cluster := GetClusterFromCtx(ctx)
	ctx = setUpNamespaceAndCleaner(ctx, t, cfg)
	// TODO: configure ingress class in the helm installation.
	ingressClass := "kong"
	t.Logf("deploying the controller's IngressClass %q", ingressClass)
	if !assert.NoError(t, helpers.CreateIngressClass(ctx, ingressClass, cluster.Client()), "failed creating IngressClass") {
		return ctx
	}
	ctx = setInCtx(ctx, _ingressClass{}, ingressClass)

	// TODO: extract setting up Kong addon to a standalone setup function.
	ctx = deployKongAddon(ctx, t, deployKongAddonCfg{
		deployControllerInKongAddon: true,
		controllerImageRepository:   oldKICImageRepo,
		controllerImageTag:          oldKICImageTag,
		kongProxyEnvVars:            map[string]string{},
	})

	// Set a dummy cancel function in ctx for cleanup.
	ctx = SetInCtxForT(ctx, t, func() {})
	return ctx
}
