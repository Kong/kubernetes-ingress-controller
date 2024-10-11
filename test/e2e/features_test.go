//go:build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// -----------------------------------------------------------------------------
// E2E feature tests
//
// These tests test features that are not easily testable using integration
// tests due to environment requirements (e.g. needing to mount volumes) or
// conflicts with the integration configuration.
// -----------------------------------------------------------------------------

const (
	// webhookKINDConfig is a KIND configuration used for TestWebhookUpdate. KIND, when running in GitHub Actions, is
	// a bit wonky with handling Secret updates, and they do not propagate to container filesystems in a reasonable
	// amount of time (>10m) when running this in the complete test suite, even though the actual sync frequency/update
	// propagation should be 1m by default. These changes force Secret updates to go directly to the API server and
	// update containers much more often. The latter causes significant performance degradation elsewhere, and Pods take
	// much longer to start, but once they do Secret updates show up more quickly, enough for the test to complete in time.
	webhookKINDConfig = `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: KubeletConfiguration
    configMapAndSecretChangeDetectionStrategy: Get
    syncFrequency: 3s
`
	validationWebhookName = "kong-validation-webhook"
	kongNamespace         = "kong"
)

// TestWebhookUpdate checks that the webhook updates the certificate indicated by --admission-webhook-cert-file when
// the mounted Secret updates. This requires E2E because we can't mount Secrets with the locally-run integration
// test controller instance.
func TestWebhookUpdate(t *testing.T) {
	// On KIND, this test requires webhookKINDConfig. the generic getEnvironmentBuilder we use for most tests doesn't
	// support this: the configuration is specific to KIND but should not be used by default, and the scaffolding isn't
	// flexible enough to support tests building their own clusters or passing additional builder functions. this still
	// uses the setup style from before getEnvironmentBuilder/GKE support as such, and just skips if it's attempting
	// to run on GKE.
	runOnlyOnKindClusters(t)

	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	clusterBuilder := kind.NewBuilder()
	clusterBuilder.WithConfigReader(strings.NewReader(webhookKINDConfig))
	if testenv.ClusterVersion() != "" {
		clusterVersion, err := semver.ParseTolerant(testenv.ClusterVersion())
		require.NoError(t, err)
		t.Logf("k8s cluster version is set to %v", clusterVersion)
		clusterBuilder.WithClusterVersion(clusterVersion)
	}
	cluster, err := clusterBuilder.Build(ctx)
	require.NoError(t, err)
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())
	if testenv.ClusterLoadImages() == "true" {
		if b, err := loadimage.NewBuilder().WithImage(testenv.ControllerImageTag()); err == nil {
			addons = append(addons, b.Build())
		} else {
			require.NoError(t, err)
		}
	}
	builder := environments.NewBuilder().WithExistingCluster(cluster).WithAddons(addons...)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	cluster = env.Cluster()
	logClusterInfo(t, cluster)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	ManifestDeploy{Path: dblessPath}.Run(ctx, t, env)

	certPool := x509.NewCertPool()
	const firstCertificateHostName = "first.example"
	firstCertificateCrt, firstCertificateKey := certificate.MustGenerateSelfSignedCertPEMFormat(
		certificate.WithCommonName(firstCertificateHostName),
		certificate.WithDNSNames(firstCertificateHostName),
	)
	require.True(t, certPool.AppendCertsFromPEM(firstCertificateCrt))
	firstCertificate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admission-cert",
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": firstCertificateCrt,
			"tls.key": firstCertificateKey,
		},
	}

	const secondCertificateHostName = "second.example"
	secondCertificateCrt, secondCertificateKey := certificate.MustGenerateSelfSignedCertPEMFormat(
		certificate.WithCommonName(secondCertificateHostName),
		certificate.WithDNSNames(secondCertificateHostName),
	)
	require.True(t, certPool.AppendCertsFromPEM(secondCertificateCrt))
	secondCertificate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admission-cert",
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": secondCertificateCrt,
			"tls.key": secondCertificateKey,
		},
	}

	_, err = env.Cluster().Client().CoreV1().Secrets(kongNamespace).Create(ctx, firstCertificate, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("exposing admission service to the test environment")
	admission, err := env.Cluster().Client().CoreV1().Services(kongNamespace).Get(ctx, validationWebhookName,
		metav1.GetOptions{})
	require.NoError(t, err)
	admission.Spec.Type = corev1.ServiceTypeLoadBalancer
	_, err = env.Cluster().Client().CoreV1().Services(kongNamespace).Update(ctx, admission, metav1.UpdateOptions{})
	require.NoError(t, err)
	var admissionAddress string
	require.Eventually(t, func() bool {
		admission, err = env.Cluster().Client().CoreV1().Services(kongNamespace).Get(ctx, validationWebhookName,
			metav1.GetOptions{})
		if err != nil {
			return false
		}
		if len(admission.Status.LoadBalancer.Ingress) > 0 {
			admissionAddress = admission.Status.LoadBalancer.Ingress[0].IP
			return true
		}
		return false
	}, time.Minute, time.Second)

	t.Log("updating kong deployment to use admission certificate")
	deployment := getManifestDeployments(dblessPath).GetController(ctx, t, env)
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == controllerContainerName {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "CONTROLLER_ADMISSION_WEBHOOK_CERT_FILE", Value: "/admission-webhook/tls.crt"},
				corev1.EnvVar{Name: "CONTROLLER_ADMISSION_WEBHOOK_KEY_FILE", Value: "/admission-webhook/tls.key"},
				corev1.EnvVar{Name: "CONTROLLER_ADMISSION_WEBHOOK_LISTEN", Value: ":8080"})

			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes,
				corev1.Volume{
					Name: "admission-cert",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "admission-cert",
						},
					},
				})

			deployment.Spec.Template.Spec.Containers[i].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[i].VolumeMounts,
				corev1.VolumeMount{Name: "admission-cert", MountPath: "/admission-webhook"})
		}
	}

	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Update(ctx,
		deployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	checkCertificate := func(hostname string) {
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			_, err := tls.Dial("tcp", admissionAddress+":443", &tls.Config{
				MinVersion: tls.VersionTLS12,
				RootCAs:    certPool,
				ServerName: hostname,
			})
			assert.NoError(c, err)
		}, 1*time.Minute, time.Second)
	}

	t.Log("checking initial certificate")
	checkCertificate(firstCertificateHostName)

	t.Log("changing certificate")
	_, err = env.Cluster().Client().CoreV1().Secrets(kongNamespace).Update(ctx, secondCertificate, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("checking second certificate")
	checkCertificate(secondCertificateHostName)
}

// TestDeployAllInOneDBLESSGateway tests the Gateway feature flag and the admission controller with no user-provided
// certificate (all other tests with the controller provide certificates, so that behavior isn't tested otherwise).
func TestDeployAllInOneDBLESSGateway(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test for Gateway")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{Path: dblessPath}.Run(ctx, t, env)
	controllerDeploymentNN := deployments.ControllerNN
	controllerDeploymentListOptions := metav1.ListOptions{
		LabelSelector: "app=" + controllerDeploymentNN.Name,
	}

	t.Log("updating controller deployment to enable alpha Gateway feature gate")
	controllerDeployment := deployments.GetController(ctx, t, env)
	for i, container := range controllerDeployment.Spec.Template.Spec.Containers {
		if container.Name == controllerContainerName {
			controllerDeployment.Spec.Template.Spec.Containers[i].Env = append(
				controllerDeployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{
					Name:  "CONTROLLER_FEATURE_GATES",
					Value: fmt.Sprintf("%s=true", featuregates.GatewayAlphaFeature),
				},
			)
		}
	}

	_, err := env.Cluster().Client().AppsV1().Deployments(namespace).Update(ctx, controllerDeployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that KIC waits for Gateway API CRDs and prints proper log")
	require.Eventually(t, func() bool {
		pods, err := env.Cluster().Client().CoreV1().Pods(controllerDeploymentNN.Namespace).List(ctx, controllerDeploymentListOptions)
		require.NoError(t, err)

		expectedMsg := "Required CustomResourceDefinitions are not installed, setting up a watch for them in case they are installed afterward"
		t.Logf("checking logs of #%d pods", len(pods.Items))
		for _, pod := range pods.Items {
			logs, err := getPodLogs(ctx, t, env, pod.Namespace, pod.Name)
			if err != nil {
				t.Logf("Failed to get logs of pods %s/%s, error %v", pod.Namespace, pod.Name, err)
				return false
			}
			if !strings.Contains(logs, expectedMsg) {
				return false
			}
		}
		return true
	}, time.Minute, 3*time.Second)

	t.Logf("deploying Gateway APIs CRDs in standard channel from %s", consts.GatewayStandardCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayStandardCRDsKustomizeURL))

	t.Log("verifying controller updates associated Gateway resoures")
	gw := deployGateway(ctx, t, env)
	verifyGateway(ctx, t, env, gw)
	deployHTTPRoute(ctx, t, env, gw)
	verifyHTTPRoute(ctx, t, env)

	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gw, err = gc.GatewayV1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)
	gw.Spec.Listeners = append(gw.Spec.Listeners,
		gatewayapi.Listener{
			Name:     "badhttp",
			Protocol: gatewayapi.HTTPProtocolType,
			Port:     gatewayapi.PortNumber(9999),
		},
		gatewayapi.Listener{
			Name:     "badudp",
			Protocol: gatewayapi.UDPProtocolType,
			Port:     gatewayapi.PortNumber(80),
		},
	)

	t.Log("verifying that unsupported listeners indicate correct status")
	gw, err = gc.GatewayV1().Gateways(corev1.NamespaceDefault).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		gw, err = gc.GatewayV1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
		var http, udp bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badhttp" {
				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayapi.ListenerConditionAccepted),
					util.ConditionReason(gatewayapi.ListenerReasonPortUnavailable),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					http = true
				}

				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayapi.ListenerConditionAccepted),
					util.ConditionReason(gatewayapi.ListenerReasonUnsupportedProtocol),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					return false
				}
			}
			if lstatus.Name == "badudp" {
				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayapi.ListenerConditionAccepted),
					util.ConditionReason(gatewayapi.ListenerReasonUnsupportedProtocol),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					udp = true
				}
			}
		}
		return http == udp == true
	}, time.Minute*2, time.Second*5)

	gw, err = gc.GatewayV1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Logf("deploying Gateway APIs CRDs in experimental channel from %s", consts.GatewayExperimentalCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayExperimentalCRDsKustomizeURL))

	t.Log("updating proxy deployment to enable TCP listener")
	proxyDeployment := deployments.GetProxy(ctx, t, env)
	for i, container := range proxyDeployment.Spec.Template.Spec.Containers {
		if container.Name == proxyContainerName {
			proxyDeployment.Spec.Template.Spec.Containers[i].Env = append(proxyDeployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "KONG_STREAM_LISTEN", Value: fmt.Sprintf("0.0.0.0:%d", tcpListenerPort)})
		}
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(namespace).Update(ctx, proxyDeployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("updating kong proxy service to enable TCP listener")
	proxyService, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	proxyService.Spec.Ports = append(proxyService.Spec.Ports, corev1.ServicePort{
		Name:       "stream-tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       tcpListenerPort,
		TargetPort: intstr.FromInt(tcpListenerPort),
	})
	_, err = env.Cluster().Client().CoreV1().Services(namespace).Update(ctx, proxyService, metav1.UpdateOptions{})
	require.NoError(t, err)

	gw = deployGatewayWithTCPListener(ctx, t, env)
	verifyGateway(ctx, t, env, gw)

	deployTCPRoute(ctx, t, env, gw)
	verifyTCPRoute(ctx, t, env)
}

// Unsatisfied LoadBalancers have special handling, see
// https://github.com/Kong/kubernetes-ingress-controller/issues/2001
func TestDeployAllInOneDBLESSNoLoadBalancer(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	ManifestDeploy{Path: dblessPath}.Run(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	// ensure that Gateways with no addresses come online and start ingesting routes
	t.Logf("deploying Gateway APIs CRDs from %s", consts.GatewayExperimentalCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayExperimentalCRDsKustomizeURL))

	deployment := getManifestDeployments(dblessPath).GetController(ctx, t, env)
	t.Log("updating controller deployment to enable Gateway feature gate")
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == controllerContainerName {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "CONTROLLER_FEATURE_GATES", Value: consts.DefaultFeatureGates})
		}
	}
	_, err := env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Update(ctx,
		deployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("updating service type to NodePort")
	svc, err := env.Cluster().Client().CoreV1().Services(deployment.Namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	svc.Spec.Type = corev1.ServiceTypeNodePort
	_, err = env.Cluster().Client().CoreV1().Services(deployment.Namespace).Update(ctx, svc, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying controller updates associated Gateway resoures")
	gw := deployGateway(ctx, t, env)
	verifyGateway(ctx, t, env, gw)
	deployHTTPRoute(ctx, t, env, gw)
	verifyHTTPRoute(ctx, t, env)
}

// TestDefaultIngressClass tests functionality related to the default Ingress class, which loads resources that have
// no class information. This is in E2E because loading classless resources interferes with integration tests that
// expect the opposite, as the integration test controller cannot use a different IngressClass than it selected at
// startup.
func TestDefaultIngressClass(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{Path: dblessPath}.Run(ctx, t, env)
	kongDeployment := deployments.ControllerNN

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(kongDeployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(kongDeployment.Namespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating a classless ingress for service %s", service.Name)
	ingress := generators.NewIngressForService("/abbosiysaltanati", map[string]string{
		annotations.AnnotationPrefix + annotations.StripPathKey: "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), kongDeployment.Namespace, ingress))

	proxyURL := "http://" + getKongProxyIP(ctx, t, env)
	t.Log("ensuring Ingress does not become live")
	require.Never(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/abbosiysaltanati", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode < http.StatusBadRequest {
			t.Logf("unexpected status when checking Ingress status: %v", resp.StatusCode)
			return true
		}
		return false
	}, time.Minute, time.Second)

	t.Logf("making our class a default IngressClass")
	class, err := env.Cluster().Client().NetworkingV1().IngressClasses().Get(ctx, "kong", metav1.GetOptions{})
	require.NoError(t, err)
	class.ObjectMeta.Annotations["ingressclass.kubernetes.io/is-default-class"] = "true"
	_, err = env.Cluster().Client().NetworkingV1().IngressClasses().Update(ctx, class, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("waiting for updated ingress status to include IP")
	require.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), kongDeployment.Namespace, ingress)
		if err != nil {
			return false
		}
		if len(lbstatus.Ingress) == 0 || lbstatus.Ingress[0].IP == "" {
			return false
		}
		return true
	}, ingressWait, time.Second)

	t.Log("getting kong proxy IP after LB provisioning")
	proxyURLForDefaultIngress := "http://" + getKongProxyIP(ctx, t, env)

	t.Log("creating classless global KongClusterPlugin")
	kongclusterplugin := &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"global": "true",
			},
		},
		PluginName: "cors",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"origins": ["example.com"]}`),
		},
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	_, err = c.ConfigurationV1().KongClusterPlugins().Create(ctx, kongclusterplugin, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("waiting for routes from Ingress to be operational")
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/abbosiysaltanati", proxyURLForDefaultIngress))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURLForDefaultIngress, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			if value, ok := resp.Header["Access-Control-Allow-Origin"]; ok {
				return strings.Contains(b.String(), "<title>httpbin.org</title>") && value[0] == "example.com"
			}
		}
		return false
	}, ingressWait, time.Second)
}
