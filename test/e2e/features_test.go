//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
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
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

// -----------------------------------------------------------------------------
// E2E feature tests
//
// These tests test features that are not easily testable using integration
// tests due to environment requirements (e.g. needing to mount volumes) or
// conflicts with the integration configuration.
// -----------------------------------------------------------------------------

// TLSPair is a PEM certificate+key pair.
type TLSPair struct {
	Key, Cert string
}

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
	admissionScriptPath   = "../../hack/deploy-admission-controller.sh"
)

// openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout cert.key -out cert.crt -days 3650 -subj '/CN=first.example/'
// openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout cert.key -out cert.crt -days 3650 -subj '/CN=first.example/'.
var tlsPairs = []TLSPair{
	{
		Cert: `-----BEGIN CERTIFICATE-----
MIICTDCCAdKgAwIBAgIUOe9HN8v1eedsZXur5uXAwJkOSG4wCgYIKoZIzj0EAwIw
XTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEWMBQGA1UEAwwNZmlyc3QuZXhhbXBsZTAe
Fw0yMjA2MTAxOTIzNDhaFw0zMjAyMDgxOTIzNDhaMF0xCzAJBgNVBAYTAkFVMRMw
EQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0
eSBMdGQxFjAUBgNVBAMMDWZpcnN0LmV4YW1wbGUwdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAAR2pbLcSQhX4gD6IyPJiRN7lxZ8aPbi6qyPyjvoTJc6DPjMuJuJgkdSC8wy
e1XFsI295WGl5gbqJsXQyJOqU6pHg6mjTEeyRxN9HbfEpH+Zp7GZ2KuTTGzi3wnh
CPqzic6jUzBRMB0GA1UdDgQWBBTPOtLEjQvk5/iy4/dhxIWWEoSJbTAfBgNVHSME
GDAWgBTPOtLEjQvk5/iy4/dhxIWWEoSJbTAPBgNVHRMBAf8EBTADAQH/MAoGCCqG
SM49BAMCA2gAMGUCMQC7rKXFcTAfoTSw5m2/ALseXru/xZC5t3Y7yQ+zSaneFMvQ
KvXcO0/RGYeqLmS58C4CMGoJva3Ad5LaZ7qgMkahhLdopePb0U/GAQqIsWhHfjOT
Il2dwxMvntBECtd0uXeKHQ==
-----END CERTIFICATE-----`,
		Key: `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDAA9OHUgH4O/xF0/qyQ
t3ZSX0/6IDilnyM1ayoUSUOfNcELUd2UZVAuZgP10f6cMUWhZANiAAR2pbLcSQhX
4gD6IyPJiRN7lxZ8aPbi6qyPyjvoTJc6DPjMuJuJgkdSC8wye1XFsI295WGl5gbq
JsXQyJOqU6pHg6mjTEeyRxN9HbfEpH+Zp7GZ2KuTTGzi3wnhCPqzic4=
-----END PRIVATE KEY-----`,
	},
	{
		Cert: `-----BEGIN CERTIFICATE-----
MIICTzCCAdSgAwIBAgIUOOTCdVckt76c9OSeGHyf+OrLU+YwCgYIKoZIzj0EAwIw
XjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEXMBUGA1UEAwwOc2Vjb25kLmV4YW1wbGUw
HhcNMjIwMjEwMTkyNTMwWhcNMzIwMjA4MTkyNTMwWjBeMQswCQYDVQQGEwJBVTET
MBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQ
dHkgTHRkMRcwFQYDVQQDDA5zZWNvbmQuZXhhbXBsZTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABHCTYbqp3P2v5aDuhkO+1rVNAidb0UcnCdtyoZx0+Oqz35Auq/GNaLvZ
RYsyW6SHVGaRWhPh3jQ8zFnc28TCGwmAMnzYPs5RHYbvBm2BSP9YWPXhc6h+lkma
HNNCu1tu56NTMFEwHQYDVR0OBBYEFEG94gMq4SvGtTs48Nw5BzVnPK69MB8GA1Ud
IwQYMBaAFEG94gMq4SvGtTs48Nw5BzVnPK69MA8GA1UdEwEB/wQFMAMBAf8wCgYI
KoZIzj0EAwIDaQAwZgIxAPRJkWfSdIQMr2R77RgCicR+adD/mMxZra2SoL7qSMyq
3iXLIXauNP9ar3tt1uZE8wIxAM4C6G4uoQ0dydhcgQVhlgB6GaqO18AEDYPzQjir
dV2Bs8EBkYBx87PmZ+e/S7g9Ug==
-----END CERTIFICATE-----`,
		Key: `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDBVtvjDBFke/k2Skezl
h63g1q5IHCQM7wr1T43m5ACKZQt0ZDE1jfm1BYKk1omNpeChZANiAARwk2G6qdz9
r+Wg7oZDvta1TQInW9FHJwnbcqGcdPjqs9+QLqvxjWi72UWLMlukh1RmkVoT4d40
PMxZ3NvEwhsJgDJ82D7OUR2G7wZtgUj/WFj14XOofpZJmhzTQrtbbuc=
-----END PRIVATE KEY-----`,
	},
}

// TestWebhookUpdate checks that the webhook updates the certificate indicated by --admission-webhook-cert-file when
// the mounted Secret updates. This requires E2E because we can't mount Secrets with the locally-run integration
// test controller instance.
func TestWebhookUpdate(t *testing.T) {
	// on KIND, this test requires webhookKINDConfig. the generic getEnvironmentBuilder we use for most tests doesn't
	// support this: the configuration is specific to KIND but should not be used by default, and the scaffolding isn't
	// flexible enough to support tests building their own clusters or passing additional builder functions. this still
	// uses the setup style from before getEnvironmentBuilder/GKE support as such, and just skips if it's attempting
	// to run on GKE
	runOnlyOnKindClusters(t)

	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	configFile, err := os.CreateTemp(os.TempDir(), "webhook-kind-config-")
	require.NoError(t, err)
	defer os.Remove(configFile.Name())
	defer configFile.Close()
	written, err := configFile.Write([]byte(webhookKINDConfig))
	require.NoError(t, err)
	require.Equal(t, len(webhookKINDConfig), written)

	clusterBuilder := kind.NewBuilder()
	clusterBuilder.WithConfig(configFile.Name())
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		t.Logf("k8s cluster version is set to %v", clusterVersion)
		clusterBuilder.WithClusterVersion(clusterVersion)
	}
	cluster, err := clusterBuilder.Build(ctx)
	require.NoError(t, err)
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithExistingCluster(cluster).WithAddons(addons...)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	firstCertificate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admission-cert",
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": []byte(tlsPairs[0].Cert),
			"tls.key": []byte(tlsPairs[0].Key),
		},
	}

	secondCertificate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admission-cert",
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": []byte(tlsPairs[1].Cert),
			"tls.key": []byte(tlsPairs[1].Key),
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
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "ingress-controller" {
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

	t.Log("checking initial certificate")
	require.Eventually(t, func() bool {
		conn, err := tls.Dial("tcp", admissionAddress+":443",
			&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}) //nolint:gosec
		if err != nil {
			t.Logf("failed to dial %s:443, error %v", admissionAddress, err)
			return false
		}
		certCommonName := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
		t.Logf("subject common name of certificate: %s", certCommonName)
		return certCommonName == "first.example"
	}, time.Minute*2, time.Second)

	t.Log("changing certificate")
	_, err = env.Cluster().Client().CoreV1().Secrets(kongNamespace).Update(ctx, secondCertificate, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("checking second certificate")
	require.Eventually(t, func() bool {
		conn, err := tls.Dial("tcp", admissionAddress+":443",
			&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}) //nolint:gosec
		if err != nil {
			t.Logf("failed to dial %s:443, error %v", admissionAddress, err)
			return false
		}
		certCommonName := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
		t.Logf("subject common name of certificate: %s", certCommonName)
		return certCommonName == "second.example"
	}, time.Minute*10, time.Second)
}

// TestDeployAllInOneDBLESSGateway tests the Gateway feature flag and the admission controller with no user-provided
// certificate (all other tests with the controller provide certificates, so that behavior isn't tested otherwise).
func TestDeployAllInOneDBLESSGateway(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test for Gateway")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)
	deploymentListOptions := metav1.ListOptions{
		LabelSelector: "app=" + deployment.Name,
	}

	t.Log("verifying that KIC disabled controllers for Gateway API and printed proper log")
	require.Eventually(t, func() bool {
		pods, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, deploymentListOptions)
		require.NoError(t, err)
		for _, pod := range pods.Items {
			logs, err := getPodLogs(ctx, t, env, pod.Namespace, pod.Name)
			if err != nil {
				t.Logf("failed to get logs of pods %s/%s, error %v", pod.Namespace, pod.Name, err)
				return false
			}
			if !strings.Contains(logs, "disabling the 'gateways' controller due to missing CRD installation") {
				return false
			}
		}
		return true
	}, time.Minute, 5*time.Second)

	t.Logf("deploying Gateway APIs CRDs in standard channel from %s", consts.GatewayStandardCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayStandardCRDsKustomizeURL))

	t.Logf("deleting KIC pods to restart them after Gateway APIs CRDs installed")
	pods, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, deploymentListOptions)
	require.NoError(t, err)
	for _, pod := range pods.Items {
		err = env.Cluster().Client().CoreV1().Pods(deployment.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
		require.NoError(t, err)
	}

	t.Log("verifying controller updates associated Gateway resoures")
	gw := deployGateway(ctx, t, env)
	verifyGateway(ctx, t, env, gw)
	deployHTTPRoute(ctx, t, env, gw)
	verifyHTTPRoute(ctx, t, env)

	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)
	gw.Spec.Listeners = append(gw.Spec.Listeners,
		gatewayv1beta1.Listener{
			Name:     "badhttp",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(9999),
		},
		gatewayv1beta1.Listener{
			Name:     "badudp",
			Protocol: gatewayv1beta1.UDPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
		},
	)

	t.Log("verifying that unsupported listeners indicate correct status")
	gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
		var http, udp bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badhttp" {
				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayv1beta1.ListenerConditionAccepted),
					util.ConditionReason(gatewayv1beta1.ListenerReasonPortUnavailable),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					http = true
				}

				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayv1beta1.ListenerConditionAccepted),
					util.ConditionReason(gatewayv1beta1.ListenerReasonUnsupportedProtocol),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					return false
				}
			}
			if lstatus.Name == "badudp" {
				if util.CheckCondition(
					lstatus.Conditions,
					util.ConditionType(gatewayv1beta1.ListenerConditionAccepted),
					util.ConditionReason(gatewayv1beta1.ListenerReasonUnsupportedProtocol),
					metav1.ConditionTrue,
					gw.Generation,
				) {
					udp = true
				}
			}
		}
		return http == udp == true
	}, time.Minute*2, time.Second*5)

	gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Logf("deploying Gateway APIs CRDs in experimental channel from %s", consts.GatewayExperimentalCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayExperimentalCRDsKustomizeURL))

	deployment, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	require.NoError(t, err)
	t.Log("updating kong deployment to enable Gateway feature gate and admission controller")
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "ingress-controller" {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "CONTROLLER_FEATURE_GATES", Value: consts.DefaultFeatureGates})
		}
		if container.Name == "proxy" {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "KONG_STREAM_LISTEN", Value: fmt.Sprintf("0.0.0.0:%d", tcpListnerPort)})
		}
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("updating kong proxy service to enable TCP listener")
	proxyService, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	proxyService.Spec.Ports = append(proxyService.Spec.Ports, corev1.ServicePort{
		Name:       "stream-tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       tcpListnerPort,
		TargetPort: intstr.FromInt(tcpListnerPort),
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	// ensure that Gateways with no addresses come online and start ingesting routes
	t.Logf("deploying Gateway APIs CRDs from %s", consts.GatewayExperimentalCRDsKustomizeURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), consts.GatewayExperimentalCRDsKustomizeURL))

	deployment, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	require.NoError(t, err)
	t.Log("updating kong deployment to enable Gateway feature gate")
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "ingress-controller" {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "CONTROLLER_FEATURE_GATES", Value: consts.DefaultFeatureGates})
		}
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Update(ctx,
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	kongDeployment := deployKong(ctx, t, env, manifest)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(kongDeployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(kongDeployment.Namespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating a classless ingress for service %s", service.Name)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/abbosiysaltanati", map[string]string{
		"konghq.com/strip-path": "true",
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
	proxyURL = "http://" + getKongProxyIP(ctx, t, env)

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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/abbosiysaltanati", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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

// TestMissingCRDsDontCrashTheController ensures that in case of missing CRDs installation in the cluster, specific
// controllers are disabled, this fact is properly logged, and the controller doesn't crash.
func TestMissingCRDsDontCrashTheController(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)

	manifest = stripCRDs(t, manifest)

	// reducing controllers' cache synchronisation timeout in order to trigger the possible process crash quicker
	cacheSyncTimeout := time.Second
	manifest = addControllerEnv(t, manifest, "CONTROLLER_CACHE_SYNC_TIMEOUT", cacheSyncTimeout.String())

	deployment := deployKong(ctx, t, env, manifest)

	t.Log("ensuring pod's ready and controller didn't crash")
	require.Never(t, func() bool {
		pods, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
		})
		if err != nil || len(pods.Items) == 0 {
			return true
		}

		pod := pods.Items[0]
		if !containerDidntCrash(pod, "ingress-controller") {
			t.Log("controller crashed")
			return true
		}

		if !isPodReady(pod) {
			t.Log("pod is not ready")
			return true
		}

		return false
	}, cacheSyncTimeout+time.Second*5, time.Second)

	t.Log("waiting for pod to output required logs")
	require.Eventually(t, func() bool {
		pods, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
		})
		if err != nil || len(pods.Items) == 0 {
			return false
		}

		podName := pods.Items[0].Name
		logs, err := getPodLogs(ctx, t, env, deployment.Namespace, podName)
		if err != nil {
			return false
		}

		resources := []string{
			"udpingresses",
			"tcpingresses",
			"kongingresses",
			"ingressclassparameterses",
			"kongplugins",
			"kongconsumers",
			"kongclusterplugins",
			"ingresses",
			"gateways",
			"httproutes",
		}
		for _, resource := range resources {
			if !strings.Contains(logs, fmt.Sprintf("disabling the '%s' controller due to missing CRD installation", resource)) {
				return false
			}
		}

		return true
	}, time.Minute, time.Second)
}
