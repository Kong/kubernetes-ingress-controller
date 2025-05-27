//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
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

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	testconsts "github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	testhelpers "github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/test/util"
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
		Assess("Upgrade KIC and verify it restarts successfully", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			runControllerManager := true
			kongBuilder, err := helpers.GenerateKongBuilderWithController()
			if !assert.NoError(t, err) {
				return ctx
			}
			if newKICImageRepo != "" && newKICImageTag != "" {
				runControllerManager = false
				t.Logf("Upgrading KIC to %s:%s", newKICImageRepo, newKICImageTag)
				kongBuilder.WithControllerImage(newKICImageRepo, newKICImageTag)
			}

			if testenv.KongImage() != "" && testenv.KongTag() != "" {
				fmt.Printf("INFO: custom kong image specified via env: %s:%s\n", testenv.KongImage(), testenv.KongTag())
			}
			namespace := GetNamespaceForT(ctx, t)

			kongBuilder.WithHelmChartVersion(testenv.KongHelmChartVersion())
			kongBuilder.WithNamespace(namespace)
			kongBuilder.WithName(NameFromT(t))
			kongBuilder.WithAdditionalValue("readinessProbe.initialDelaySeconds", "1")

			if runControllerManager {
				t.Log("Run controller manager from code base, disabling KIC in Kong addon")
				kongBuilder.WithControllerDisabled()
			}

			kongAddon := kongBuilder.Build()
			t.Logf("deploying kong addon to cluster %s in namespace %s", cluster.Name(), namespace)
			if !assert.NoError(t, cluster.DeployAddon(ctx, kongAddon)) {
				return ctx
			}

			t.Log("Waiting for Kong addon to be ready")
			if !assert.Eventually(t, func() bool {
				_, ok, err := kongAddon.Ready(ctx, cluster)
				if err != nil {
					t.Logf("error checking if kong addon is ready: %v", err)
					return false
				}

				return ok
			}, time.Minute*3, 100*time.Millisecond, "failed waiting for kong addon to become ready") {
				return ctx
			}

			if runControllerManager {
				logger, logOutput, err := testutils.SetupLoggers("trace", "text")
				if !assert.NoError(t, err, "failed to setup loggers") {
					return ctx
				}
				if logOutput != "" {
					t.Logf("writing manager logs to %s", logOutput)
				}

				featureGates := testconsts.DefaultFeatureGates

				t.Logf("feature gates enabled: %s", featureGates)

				t.Logf("starting the controller manager")
				cert, key := certificate.GetKongSystemSelfSignedCerts()
				metricsPort := testhelpers.GetFreePort(t)
				healthProbePort := testhelpers.GetFreePort(t)
				ingressClass := "kong"
				extraControllerArgs := []string{}
				if testenv.DBMode() != testenv.DBModeOff {
					extraControllerArgs = append(extraControllerArgs,
						fmt.Sprintf("--kong-admin-token=%s", testconsts.KongTestPassword),
						fmt.Sprintf("--kong-workspace=%s", testconsts.KongTestWorkspace),
					)
				}

				standardControllerArgs := []string{
					fmt.Sprintf("--health-probe-bind-address=localhost:%d", healthProbePort),
					fmt.Sprintf("--metrics-bind-address=localhost:%d", metricsPort),
					fmt.Sprintf("--ingress-class=%s", ingressClass),
					fmt.Sprintf("--admission-webhook-cert=%s", cert),
					fmt.Sprintf("--admission-webhook-key=%s", key),
					fmt.Sprintf("--admission-webhook-listen=0.0.0.0:%d", testutils.AdmissionWebhookListenPort),
					"--anonymous-reports=false",
					"--log-level=trace",
					"--dump-config=true",
					"--dump-sensitive-config=true",
					fmt.Sprintf("--feature-gates=%s", featureGates),
					// Use fixed election namespace `kong` because RBAC roles for leader election are in the namespace,
					// so we create resources for leader election in the namespace to make sure that KIC can operate these resources.
					fmt.Sprintf("--election-namespace=%s", testconsts.ControllerNamespace),
					fmt.Sprintf("--watch-namespace=%s", kongAddon.Namespace()),
				}
				allControllerArgs := slices.Concat(standardControllerArgs, extraControllerArgs)

				gracefulShutdownWithoutTimeoutOpt := func(c *managercfg.Config) {
					// Set the GracefulShutdownTimeout to -1 to keep graceful shutdown enabled but disable the timeout.
					// This prevents the errors:
					// failed waiting for all runnables to end within grace period of 30s: context deadline exceeded
					c.GracefulShutdownTimeout = lo.ToPtr(time.Duration(-1))
				}

				cancel, err := testutils.DeployControllerManagerForCluster(ctx, logger, cluster, kongAddon, allControllerArgs, gracefulShutdownWithoutTimeoutOpt)
				t.Cleanup(func() { cancel() })
				if !assert.NoError(t, err, "failed deploying controller manager") {
					return ctx
				}

				// TODO refactor. Perhaps there's a better way than just storing the cancel func in context.
				ctx = SetInCtxForT(ctx, t, cancel)
			}

			return ctx

		}).
		Assess("Verify that the old configuration remains unchanged", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
		Assess("Update the plugin and verify that the new configuration works", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			kongClient := GetFromCtxForT[*clientset.Clientset](ctx, t)
			assert.NotNil(t, kongClient)
			cluster := GetClusterFromCtx(ctx)

			namespace := GetNamespaceForT(ctx, t)

			plugin, err := kongClient.ConfigurationV1().KongPlugins(namespace).Get(ctx, pluginName, metav1.GetOptions{})
			assert.NoErrorf(t, err, "failed to get plugin %s/%s", namespace, pluginName)

			plugin.Config = apiextensionsv1.JSON{
				Raw: []byte(`{"add":{"headers":["Kic-Added:Another-Test"]}}`),
			}
			plugin, err = kongClient.ConfigurationV1().KongPlugins(namespace).Update(ctx, plugin, metav1.UpdateOptions{})
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
			cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, newIngress, metav1.CreateOptions{})

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
	runID := GetRunIDFromCtx(ctx)

	// TODO: configure ingress class in the helm installation.
	ingressClass := "kong"
	t.Logf("deploying the controller's IngressClass %q", ingressClass)
	if !assert.NoError(t, helpers.CreateIngressClass(ctx, ingressClass, cluster.Client()), "failed creating IngressClass") {
		return ctx
	}
	ctx = setInCtx(ctx, _ingressClass{}, ingressClass)

	ctx, err := CreateNSForTest(ctx, cfg, t, runID)
	if !assert.NoError(t, err) {
		return ctx
	}
	// TODO: extract setting up Kong addon to a standalone setup function.
	t.Logf("setting up test environment")
	kongBuilder, err := helpers.GenerateKongBuilderWithController()
	if !assert.NoError(t, err) {
		return ctx
	}
	kongBuilder.WithControllerImage(oldKICImageRepo, oldKICImageTag)
	if testenv.KongImage() != "" && testenv.KongTag() != "" {
		fmt.Printf("INFO: custom kong image specified via env: %s:%s\n", testenv.KongImage(), testenv.KongTag())
	}

	namespace := GetNamespaceForT(ctx, t)

	kongBuilder.WithHelmChartVersion(testenv.KongHelmChartVersion())
	kongBuilder.WithNamespace(namespace)
	kongBuilder.WithName(NameFromT(t))
	kongBuilder.WithAdditionalValue("readinessProbe.initialDelaySeconds", "1")

	kongAddon := kongBuilder.Build()
	t.Logf("deploying kong addon to cluster %s in namespace %s", cluster.Name(), namespace)
	if !assert.NoError(t, cluster.DeployAddon(ctx, kongAddon)) {
		return ctx
	}
	ctx = SetInCtxForT(ctx, t, kongAddon)

	cleaner := clusters.NewCleaner(cluster)
	t.Cleanup(func() {
		helpers.DumpDiagnosticsIfFailed(ctx, t, cluster)
		t.Logf("Start cleanup for test %s", t.Name())
		if err := cleaner.Cleanup(context.Background()); err != nil { //nolint:contextcheck
			fmt.Printf("ERROR: failed cleaning up the cluster: %v\n", err)
		}
	})
	ctx = SetInCtxForT(ctx, t, cleaner)

	t.Log("Waiting for Kong addon to be ready")
	if !assert.Eventually(t, func() bool {
		_, ok, err := kongAddon.Ready(ctx, cluster)
		if err != nil {
			t.Logf("error checking if kong addon is ready: %v", err)
			return false
		}

		return ok
	}, time.Minute*3, 100*time.Millisecond, "failed waiting for kong addon to become ready") {
		return ctx
	}

	t.Logf("collecting urls from the kong proxy deployment in namespace: %s", namespace)
	proxyAdminURL, err := kongAddon.ProxyAdminURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	ctx = SetAdminURLInCtx(ctx, proxyAdminURL)

	proxyUDPURL, err := kongAddon.ProxyUDPURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	ctx = SetUDPURLInCtx(ctx, proxyUDPURL)

	proxyTCPURL, err := kongAddon.ProxyTCPURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	ctx = SetTCPURLInCtx(ctx, proxyTCPURL)

	proxyTLSURL, err := kongAddon.ProxyTLSURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	ctx = SetTLSURLInCtx(ctx, proxyTLSURL)

	proxyHTTPURL, err := kongAddon.ProxyHTTPURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	t.Log("proxy HTTP URL:", proxyHTTPURL.String())
	ctx = SetHTTPURLInCtx(ctx, proxyHTTPURL)

	proxyHTTPSURL, err := kongAddon.ProxyHTTPSURL(ctx, cluster)
	if !assert.NoError(t, err) {
		return ctx
	}
	ctx = SetHTTPSURLInCtx(ctx, proxyHTTPSURL)

	if !assert.NoError(t, retry.Do(
		func() error {
			version, err := helpers.GetKongVersion(ctx, proxyAdminURL, testconsts.KongTestPassword)
			if err != nil {
				return err
			}
			t.Logf("using Kong instance (version: %s) reachable at %s\n", version, proxyAdminURL)
			return nil
		},
		retry.OnRetry(
			func(n uint, err error) {
				t.Logf("WARNING: try to get Kong Gateway version attempt %d/10 - error: %s\n", n+1, err)
			},
		),
		retry.LastErrorOnly(true),
		retry.Attempts(10),
	), "failed getting Kong's version") {
		return ctx
	}

	ctx = SetInCtxForT(ctx, t, func() {})
	return ctx
}
