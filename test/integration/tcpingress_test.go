//go:build integration_tests

package integration

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

func TestTCPIngressTLS(t *testing.T) {
	ctx := context.Background()
	RunWhenKongExpressionRouter(ctx, t)
	t.Parallel()

	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress-%s"
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	testServiceSuffixes := []string{"alpha", "bravo", "charlie"}
	const domain = ".example"
	var certOpts []certificate.SelfSignedCertificateOption
	for _, tss := range testServiceSuffixes {
		certOpts = append(certOpts, certificate.WithDNSNames(tss+domain), certificate.WithCommonName(tss+domain))
	}
	exampleTLSCert, exampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certOpts...)
	tlsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: tlsSecretName,
		},
		Data: map[string][]byte{
			"tls.crt": exampleTLSCert,
			"tls.key": exampleTLSKey,
		},
	}
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, tlsSecret, metav1.CreateOptions{})
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	require.True(t, certPool.AppendCertsFromPEM(exampleTLSCert))

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
			TLS: []kongv1beta1.IngressTLS{
				{
					Hosts:      []string{testServiceSuffixes[0] + domain, testServiceSuffixes[1] + domain},
					SecretName: tlsSecretName,
				},
			},
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[0] + domain,
					Port: ktfkong.DefaultTLSServicePort,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: testServices[testServiceSuffixes[0]].Name,
						ServicePort: test.EchoTCPPort,
					},
				},
				{
					Host: testServiceSuffixes[1] + domain,
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
			TLS: []kongv1beta1.IngressTLS{
				{
					Hosts:      []string{testServiceSuffixes[2] + domain},
					SecretName: tlsSecretName,
				},
			},
			Rules: []kongv1beta1.IngressRule{
				{
					Host: testServiceSuffixes[2] + domain,
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
				MinVersion: tls.VersionTLS12,
				ServerName: fmt.Sprintf("%s.example", i),
				RootCAs:    certPool,
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
			MinVersion: tls.VersionTLS12,
			ServerName: fmt.Sprintf("%s.example", testServiceSuffixes[0]),
			RootCAs:    certPool,
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
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(exampleTLSCert)
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
		err := tlsEchoResponds(proxyTLSURL, testUUID, tlsRouteHostname, certPool, true)
		assert.NoError(c, err)
	}, ingressWait, waitTick)
}
