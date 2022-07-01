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

var tcpMutex sync.Mutex
var tlsMutex sync.Mutex

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
			"tls.crt": []byte(tlsPairs[0].Cert),
			"tls.key": []byte(tlsPairs[0].Key),
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
			return false
		}
		defer conn.Close()
		err = conn.Handshake()
		if err != nil {
			return false
		}
		cert := conn.ConnectionState().PeerCertificates[0]
		return cert.Subject.CommonName == "secure-foo-bar"
	}, ingressWait, waitTick)
}
