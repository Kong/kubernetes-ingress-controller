//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

var (
	tcpMutex sync.Mutex
	tlsMutex sync.Mutex
)

func TestTCPIngressEssentials(t *testing.T) {
	RunWhenKongExpressionRouter(t)
	ctx := context.Background()

	t.Parallel()
	// Ensure no other TCP tests run concurrently to avoid fights over the port
	t.Log("locking TCP port")
	tcpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	})

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	deployment := generators.NewDeploymentForContainer(generators.NewContainer(testName, test.HTTPBinImage, test.HTTPBinPort))
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("routing to service %s via TCPIngress", service.Name)
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: ktfkong.DefaultTCPServicePort,
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
	cleaner.Add(tcp)

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
	tcpProxyURL := fmt.Sprintf("http://%s:", proxyTCPURL)
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(tcpProxyURL)
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
		resp, err := helpers.DefaultHTTPClient().Get(tcpProxyURL)
		if err != nil {
			return true
		}
		defer resp.Body.Close()
		return false
	}, ingressWait, waitTick)
}

func TestTCPIngressTLS(t *testing.T) {
	RunWhenKongExpressionRouter(t)
	t.Parallel()

	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress-%s"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	testServiceSuffixes := []string{"alpha", "bravo", "charlie"}
	testServices := make(map[string]*corev1.Service)

	for _, i := range testServiceSuffixes {
		localTestName := fmt.Sprintf(testName, i)
		t.Log("deploying a minimal TCP container deployment to test Ingress routes")
		container := generators.NewContainer(localTestName, test.EchoImage, test.EchoTCPPort)
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
		cleaner.Add(deployment)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)
		testServices[i] = service
		cleaner.Add(service)
	}

	t.Log("adding TCPIngresses")
	tcpX := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(testName, "x"),
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[0] + ".example",
					Port: ktfkong.DefaultTLSServicePort,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[0]].Name,
						ServicePort: test.EchoTCPPort,
					},
				},
				{
					Host: testServiceSuffixes[1] + ".example",
					Port: ktfkong.DefaultTLSServicePort,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[1]].Name,
						ServicePort: test.EchoTCPPort,
					},
				},
			},
		},
	}
	tcpX, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcpX, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tcpX)

	tcpY := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(testName, "y"),
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[2] + ".example",
					Port: ktfkong.DefaultTLSServicePort,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[2]].Name,
						ServicePort: test.EchoTCPPort,
					},
				},
			},
		},
	}
	tcpY, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcpY, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tcpY)
	for _, i := range testServiceSuffixes {
		t.Logf("verifying TCP Ingress for %s.example operational", i)
		require.Eventually(t, func() bool {
			conn, err := tls.Dial("tcp", proxyTLSURL, &tls.Config{
				InsecureSkipVerify: true,
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
		conn, err := tls.Dial("tcp", proxyTLSURL, &tls.Config{
			InsecureSkipVerify: true,
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
	t.Parallel()

	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("setting up the TCPIngress TLS passthrough tests")

	const tlsExampleHostname = "tlsroute.kong.example"

	t.Log("configuring secrets")
	exampleTLSCert, exampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(
		certificate.WithCommonName(tlsExampleHostname), certificate.WithDNSNames(tlsExampleHostname),
	)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tlsSecretName,
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"tls.crt": exampleTLSCert,
			"tls.key": exampleTLSKey,
		},
	}

	k8sClient := env.Cluster().Client()
	t.Log("deploying secrets")
	_, err := k8sClient.CoreV1().Secrets(ns.Name).Create(ctx, secret, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Log("creating a tcpecho deployment to test TLSRoute traffic routing")
	testUUID := uuid.NewString() // go-echo sends a "Running on Pod <UUID>." immediately on connecting
	deployment := generators.NewDeploymentForContainer(createTLSEchoContainer(tlsEchoPort, testUUID))
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: tlsSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: tlsSecretName,
			},
		},
	})
	deployment, err = k8sClient.AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = k8sClient.CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Log("adding TCPIngress")
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey:                             consts.IngressClass,
				annotations.AnnotationPrefix + annotations.ProtocolsKey: "tls_passthrough",
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Host: tlsExampleHostname,
					Port: ktfkong.DefaultTLSServicePort,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: tlsEchoPort,
					},
				},
			},
		},
	}
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	tcp, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcp, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tcp)

	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		err := tlsEchoResponds(proxyTLSURL, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		assert.NoError(c, err)
	}, ingressWait, waitTick)
}
