//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"crypto/tls"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
)

// -----------------------------------------------------------------------------
// E2E feature tests
//
// These tests test features that are not easily testable using integration
// tests due to environment requirements (e.g. needing to mount volumes) or
// conflicts with the integration configuration.
// -----------------------------------------------------------------------------

// TLSPair is a PEM certificate+key pair
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

var (
	// openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout cert.key -out cert.crt -days 3650 -subj '/CN=first.example/'
	// openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout cert.key -out cert.crt -days 3650 -subj '/CN=first.example/'
	tlsPairs = []TLSPair{
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
)

// TestWebhookUpdate checks that the webhook updates the certificate indicated by --admission-webhook-cert-file when
// the mounted Secret updates. This requires E2E because we can't mount Secrets with the locally-run integration
// test controller instance.
func TestWebhookUpdate(t *testing.T) {
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
		assert.NoError(t, env.Cleanup(ctx))
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
							SecretName: "admission-cert"},
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
			&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}) // nolint:gosec
		if err != nil {
			return false
		}
		return conn.ConnectionState().PeerCertificates[0].Subject.CommonName == "first.example"
	}, time.Minute*2, time.Second)

	t.Log("changing certificate")
	_, err = env.Cluster().Client().CoreV1().Secrets(kongNamespace).Update(ctx, secondCertificate, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("checking second certificate")
	require.Eventually(t, func() bool {
		conn, err := tls.Dial("tcp", admissionAddress+":443",
			&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}) // nolint:gosec
		if err != nil {
			return false
		}
		return conn.ConnectionState().PeerCertificates[0].Subject.CommonName == "second.example"
	}, time.Minute*10, time.Second)
}

// TestDeployAllInOneDBLESSGateway tests the Gateway feature flag and the admission controller with no user-provided
// certificate (all other tests with the controller provide certificates, so that behavior isn't tested otherwise)
func TestDeployAllInOneDBLESSGateway(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test for Gateway")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up environment for cluster %s", env.Cluster().Name())
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Logf("deploying Gateway APIs CRDs from %s", gatewayCRDsURL)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), gatewayCRDsURL))

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("running the admission webhook setup script")
	cmd := exec.Command("bash", admissionScriptPath)
	require.NoError(t, cmd.Run())

	deployment, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	require.NoError(t, err)
	t.Log("updating kong deployment to enable Gateway feature gate and admission controller")
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "ingress-controller" {
			deployment.Spec.Template.Spec.Containers[i].Env = append(deployment.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: "CONTROLLER_FEATURE_GATES", Value: "Gateway=true"})
		}
	}

	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Update(ctx,
		deployment, metav1.UpdateOptions{})

	// vov it's easier than tracking the deployment state
	t.Log("creating a consumer to ensure the admission webhook is online")
	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nihoniy",
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Username: "nihoniy",
	}

	kongClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.Eventually(t, func() bool {
		_, err = kongClient.ConfigurationV1().KongConsumers(namespace).Create(ctx, consumer, metav1.CreateOptions{})
		return err == nil
	}, time.Minute*2, time.Second*1)

	t.Log("verifying controller updates associated Gateway resoures")
	gw := deployGateway(ctx, t, env)
	verifyGateway(ctx, t, env, gw)
	deployHTTPRoute(ctx, t, env, gw)
	verifyHTTPRoute(ctx, t, env)

	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)
	gw.Spec.Listeners = append(gw.Spec.Listeners,
		gatewayv1alpha2.Listener{
			Name:     "badhttp",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(9999),
		},
		gatewayv1alpha2.Listener{
			Name:     "badudp",
			Protocol: gatewayv1alpha2.UDPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
		},
	)

	t.Log("verifying that unsupported listeners indicate correct status")
	gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
		var http, udp bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badhttp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionDetached) {
						if condition.Reason == string(gatewayv1alpha2.ListenerReasonPortUnavailable) {
							http = true
						}
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionDetached) {
						if condition.Reason == string(gatewayv1alpha2.ListenerReasonUnsupportedProtocol) {
							return false
						}
					}
				}
			}
			if lstatus.Name == "badudp" {
				for _, condition := range lstatus.Conditions {
					// no check against the other reason here: this gets both the port and protocol condition
					if condition.Type == string(gatewayv1alpha2.ListenerConditionDetached) {
						if condition.Reason == string(gatewayv1alpha2.ListenerReasonUnsupportedProtocol) {
							udp = true
						}
					}
				}
			}
		}
		return http == udp == true
	}, time.Minute*2, time.Second*5)

	gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)
}

// Unsatisfied LoadBalancers have special handling, see
// https://github.com/Kong/kubernetes-ingress-controller/issues/2001
func TestDeployAllInOneDBLESSNoLoadBalancer(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	_ = deployKong(ctx, t, env, manifest)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)
}
