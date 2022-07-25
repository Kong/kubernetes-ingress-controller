//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

var (
	tcpMutex sync.Mutex
	tlsMutex sync.Mutex
)

func TestTCPIngressEssentials(t *testing.T) {
	t.Parallel()
	// Ensure no other TCP tests run concurrently to avoid fights over the port
	// Free it when done
	t.Log("locking TCP port")
	tcpMutex.Lock()
	defer func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	}()

	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	deployment := generators.NewDeploymentForContainer(generators.NewContainer(testName, test.HTTPBinImage, 80))
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("routing to service %s via TCPIngress", service.Name)
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 8888,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: 80,
					},
				},
			},
		},
	}
	tcp, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcp, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("ensuring that TCPIngress %s is cleaned up", tcp.Name)
		if err := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("checking tcpingress %s status readiness.", tcp.Name)
	ingCli := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, tcp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("tcpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Logf("verifying TCP Ingress %s operational", tcp.Name)
	tcpProxyURL, err := url.Parse(fmt.Sprintf("http://%s:8888/", proxyURL.Hostname()))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(tcpProxyURL.String())
		if err != nil {
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("tearing down TCPIngress %s and ensuring that the relevant backend routes are removed", tcp.Name)
	require.NoError(t, gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(tcpProxyURL.String())
		if err != nil {
			return true
		}
		defer resp.Body.Close()
		return false
	}, ingressWait, waitTick)
}

func TestTCPIngressTLS(t *testing.T) {
	t.Parallel()
	t.Log("locking TLS port")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()

	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress-%s"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	testServiceSuffixes := []string{"alpha", "bravo", "charlie"}
	testServices := make(map[string]*corev1.Service)

	for _, i := range testServiceSuffixes {
		localTestName := fmt.Sprintf(testName, i)
		t.Log("deploying a minimal TCP container deployment to test Ingress routes")
		container := generators.NewContainer(localTestName, test.TCPEchoImage, tcpEchoPort)
		// go-echo sends a "Running on Pod POD_NAME." immediately on connecting
		container.Env = []corev1.EnvVar{
			{
				Name:  "POD_NAME",
				Value: i,
			},
		}
		deployment := generators.NewDeploymentForContainer(container)
		deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
		require.NoError(t, err)

		defer func() {
			t.Logf("cleaning up the deployment %s", deployment.Name)
			assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
		}()

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)
		testServices[i] = service

		defer func() {
			t.Logf("cleaning up the service %s", service.Name)
			assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
		}()
	}

	t.Log("adding TCPIngresses")
	tcpX := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(testName, "x"),
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[0] + ".example",
					Port: 8899,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[0]].Name,
						ServicePort: tcpEchoPort,
					},
				},
				{
					Host: testServiceSuffixes[1] + ".example",
					Port: 8899,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[1]].Name,
						ServicePort: tcpEchoPort,
					},
				},
			},
		},
	}
	tcpX, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcpX, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("ensuring that TCPIngress %s is cleaned up", tcpX.Name)
		if err := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcpX.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	tcpY := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(testName, "y"),
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[2] + ".example",
					Port: 8899,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[2]].Name,
						ServicePort: tcpEchoPort,
					},
				},
			},
		},
	}
	tcpY, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcpY, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("ensuring that TCPIngress %s is cleaned up", tcpY.Name)
		if err := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcpY.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	for _, i := range testServiceSuffixes {
		t.Logf("verifying TCP Ingress for %s.example operational", i)
		require.Eventually(t, func() bool {
			conn, err := tls.Dial("tcp", fmt.Sprintf("%s:8899", proxyURL.Hostname()), &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				ServerName:         fmt.Sprintf("%s.example", i),
			})
			if err != nil {
				return false
			}
			defer conn.Close()
			resp := make([]byte, 512)
			require.NoError(t, conn.SetDeadline(time.Now().Add(time.Second*5)))
			_, err = conn.Read(resp)
			if err != nil {
				return false
			}
			if strings.Contains(string(resp), i) {
				return true
			}
			return false
		}, ingressWait, waitTick)
	}

	// Update wipes out tcpY if actually assigned, breaking the deferred delete. we have no use for it, so discard it
	require.Eventually(t, func() bool {
		tcpY, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Get(ctx, tcpY.Name, metav1.GetOptions{})
		tcpY.Spec.Rules[0].Backend.ServiceName = testServiceSuffixes[0]
		_, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Update(ctx, tcpY, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Logf("verifying TCP Ingress routes to new upstream after update")
	require.Eventually(t, func() bool {
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:8899", proxyURL.Hostname()), &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			ServerName:         fmt.Sprintf("%s.example", testServiceSuffixes[0]),
		})
		if err != nil {
			return false
		}
		defer conn.Close()
		resp := make([]byte, 512)
		require.NoError(t, conn.SetDeadline(time.Now().Add(time.Second*5)))
		_, err = conn.Read(resp)
		if err != nil {
			return false
		}
		if strings.Contains(string(resp), testServiceSuffixes[0]) {
			return true
		}
		return false
	}, ingressWait, waitTick)
}

// once upon a time, our upstream Redis image vendor decided to add a passphrase setting that didn't check if it wasn't
// set and should therefore not be required, so now this test gets a special cert.
var (
	passphrasePair = TLSPair{
		Cert: `-----BEGIN CERTIFICATE-----
MIIFbTCCA1WgAwIBAgIUYbB5HeN2T1yaXlc/JtBRAcC9IuswDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMjA3MjUxNzQxMzhaGA8yMTIy
MDcwMTE3NDEzOFowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAK6Ka7/QCxbYWql+bQRSARWYcmy/0sZYzwhVdU+2
MMtHIgzeBwc3hAexPzhckskFz+Tmvue2pSNeIK14tYjp9SS6PjZxMP9/aZFS3wtY
3nd79VtuxDdA+jvjWeIAKqQIiYJtLW/UCByUc8b3dhyyE4X2/67tVOsxM9OCaBgn
XnyFRHSFR54ca/26iN3LB0fOtXDP+lxNnjs7A5v0Uuh2DvZ04hE+bjc5XRuB4TQj
Ii/3OU+ca3RsjZ9YRKc1DaokUiMpkL/ODYMwxgSmhcmEV8QnzbNQZKLPcABtVJw5
S4crjCXh0EXRtP6Di84xhbNQf/24djFjIZN0jCmH50isJc9SFVdOsjNZixb6MOhy
FF+57N9LYNNuShhUhC1f/OKpBg63mlGYFfE5qfLN6pVCgnOUgNmc01RGfw8/xbTe
BFZp8OA9PObY+HXxsMOBXDzJ0bea/IC8JmxkcOPZ3fqveBo7aTb2HVGts6THwhUI
d/Csl5q0+VgCfIsu2i5Mqjv/6mYL1FehfAsL5whPlICEaLgnVOPu2x4QS0dHn7dE
6BxeIhlq1CvZNFw0ONJayoqGfk0b3igc0PfmLXj9TFx0I5hn8F0/KQl0x3Al+gib
2ugH49efPJXL+IULwGySNB0CikrqQMIaqyf9rDkpW/L23XMEZ6nXE9asezK3N3Co
5FTnAgMBAAGjUzBRMB0GA1UdDgQWBBQQ8Kql9b8QuZVu/AoNEwXBScnCNTAfBgNV
HSMEGDAWgBQQ8Kql9b8QuZVu/AoNEwXBScnCNTAPBgNVHRMBAf8EBTADAQH/MA0G
CSqGSIb3DQEBCwUAA4ICAQAtxkL9ld5cqD3g/YJOE7g6s1lja6YDqR9nAxNurL5r
IaxIkLfXa6hx5MEOJ5qsBefpQiY1rhoXYawQ25HHSjPLqGecWDvTdZUFny9JPgs+
Q72Iphcclk10eOmeY5f33yY8/QEwl5P7BbICKIcp6EHk4yBbv1fCWPAeT40heqkY
A69gXlWmZIW61q1dR23ul07nv0uPcJEvN37WM0gjOPYez3bfMzV/jvldYqXcv0zb
GFuVvlaWJwWmn7FScCCq4J73F5bklx7fZJuknfmGegwgaYxtroe3JkaxfTnspomA
kcasWMaGzNtEO4ylB8JMoUUIXIZ53GXtjJDSA+AjiDduAxe9zvJsw2JE8bbyiSC0
ira9BKAn8MG4awwrDLVYj/oVdpZqwQlXOQ/8T72SLpNfgUCTQ24E/Hdw/Il7yC+h
2koyipD0Mkn7y46g6c/qcVjBqPTtpcep7iYmOohXhaADsdhPn7FNO1/CRabkLfGY
0FN9AWSMKRBpSWEXW1wAy8wrh8stA59zGoGEB7JuX6UrwZfZKGQdUH/ATA2G6TVZ
YLp4cgCwkL+OvWojOz1Mq9j0E8IYOaFQ2RIAgu/8BqbR2UVAY6nWqfg8eOIcY+OE
WMfexPAZa49uF0E7BA0h1jxzALtmz96CoKy+KI1EsF23ZUW67M8lq/rGYMlcBIiF
Tg==
-----END CERTIFICATE-----`,
		Key: `-----BEGIN ENCRYPTED PRIVATE KEY-----
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
8kBvpExREohI9skT5YquWQu3ClOPpcly
-----END ENCRYPTED PRIVATE KEY-----`,
	}
	passphrasePairPassphrase = "passphrase"
)

func TestTCPIngressTLSPassthrough(t *testing.T) {
	version, err := getKongVersion()
	if err != nil {
		t.Logf("attempting TLS passthrough test despite unknown kong version: %v", err)
	} else if version.LT(semver.MustParse("2.7.0")) {
		t.Skipf("kong version %s below minimum TLS passthrough version", version)
	}

	t.Parallel()
	t.Log("locking TLS port")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()

	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("setting up the TCPIngress TLS passthrough tests")
	testName := "tlspass"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("configuring secrets")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"tls.crt":        []byte(passphrasePair.Cert),
			"tls.key":        []byte(passphrasePair.Key),
			"tls.passphrase": []byte(passphrasePairPassphrase),
		},
	}

	t.Log("deploying secrets")
	secret, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secret, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Log("deploying Redis with certificate")
	container := generators.NewContainer(testName, redisImage, 6379)
	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "certificates",
			MountPath: "/opt/certs",
		},
	}
	container.Env = []corev1.EnvVar{
		{
			Name:  "REDIS_TLS_ENABLED",
			Value: "true",
		},
		{
			Name:  "REDIS_TLS_PORT",
			Value: "6379",
		},
		{
			Name:  "REDIS_TLS_CA_FILE",
			Value: "/opt/certs/tls.crt",
		},
		{
			Name:  "REDIS_TLS_CERT_FILE",
			Value: "/opt/certs/tls.crt",
		},
		{
			Name:  "REDIS_TLS_KEY_FILE",
			Value: "/opt/certs/tls.key",
		},
		{
			Name:  "REDIS_TLS_KEY_FILE_PASS",
			Value: "/opt/certs/tls.passphrase",
		},
		{
			Name:  "REDIS_PASSWORD",
			Value: "garbage",
		},
	}
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "certificates",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret.Name,
				},
			},
		},
	}
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Log("waiting for deployment to be ready")
	deploymentName := deployment.Name
	require.Eventually(t, func() bool {
		deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			t.Logf("failed to get deployment %s/%s , error %v", ns.Name, deploymentName, err)
			return false
		}
		if deployment.Status.Replicas == deployment.Status.AvailableReplicas {
			return true
		}
		t.Logf("deployment not ready: %d/%d pods available", deployment.Status.AvailableReplicas, deployment.Status.Replicas)
		return false
	}, ingressWait, waitTick, func() string {
		// dump status of all pods.
		podList, err := env.Cluster().Client().CoreV1().Pods(ns.Name).List(
			ctx, metav1.ListOptions{
				LabelSelector: "app=" + testName,
			})
		if err != nil {
			return err.Error()
		}
		podStatusString := []string{}
		for _, pod := range podList.Items {
			podStatusString = append(podStatusString, fmt.Sprintf("pod %s/%s: phase %s",
				pod.Namespace, pod.Name, pod.Status.Phase))
		}
		return strings.Join(podStatusString, "\n")
	}())

	t.Log("adding TCPIngress")
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey:                             ingressClass,
				annotations.AnnotationPrefix + annotations.ProtocolsKey: "tls_passthrough",
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: "redis.example",
					Port: 8899,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: 6379,
					},
				},
			},
		},
	}
	tcp, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcp, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Log("ensuring that TCPIngress is cleaned up", tcp.Name)
		if err := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("verifying TCP Ingress for redis.example operational")
	require.Eventually(t, func() bool {
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:8899", proxyURL.Hostname()), &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			ServerName:         "redis.example",
		})
		if err != nil {
			t.Logf("failed to connect to %s:8899, error %v, retrying...", proxyURL.Hostname(), err)
			return false
		}
		defer conn.Close()
		err = conn.Handshake()
		if err != nil {
			t.Logf("failed to do tls handshake to %s:8899, error %v, retrying...", proxyURL.Hostname(), err)
			return false
		}
		cert := conn.ConnectionState().PeerCertificates[0]
		return cert.Subject.CommonName == "secure-foo-bar"
	}, ingressWait, waitTick)
}
