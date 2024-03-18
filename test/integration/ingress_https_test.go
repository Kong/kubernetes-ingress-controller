//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestHTTPSRedirect(t *testing.T) {
	RunWhenKongExpressionRouter(t)
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("creating an HTTP container via deployment to test redirect functionality")
	container := generators.NewContainer("alsohttpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	opts := metav1.CreateOptions{}
	_, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, opts)
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via Service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, opts)
	assert.NoError(t, err)
	cleaner.Add(service)

	t.Logf("exposing Service %s via Ingress", service.Name)
	ingress := generators.NewIngressForService("/test_https_redirect", map[string]string{
		"konghq.com/protocols":                  "https",
		"konghq.com/https-redirect-status-code": "301",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	assert.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for Ingress to be operational and properly redirect")
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}
	assert.Eventually(t, func() bool {
		resp, err := client.Get(fmt.Sprintf("%s/test_https_redirect", proxyURL))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusMovedPermanently
	}, ingressWait, waitTick)
}

func TestHTTPSIngress(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	testTransport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "foo.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			if addr == "bar.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			if addr == "baz.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			return dialer.DialContext(ctx, network, addr)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpcStatic := http.Client{
		Timeout:   httpcTimeout,
		Transport: &testTransport,
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, consts.IngressClass)
	ingress1 := generators.NewIngressForService("/foo", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress1.Spec.IngressClassName = kong.String(consts.IngressClass)
	ingress2 := generators.NewIngressForService("/bar", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress2.Spec.IngressClassName = kong.String(consts.IngressClass)

	t.Log("configuring ingress tls spec")
	ingress1.Spec.TLS = []netv1.IngressTLS{{SecretName: "secret1", Hosts: []string{"foo.example"}}}
	ingress1.ObjectMeta.Name = "ingress1"
	ingress2.Spec.TLS = []netv1.IngressTLS{{SecretName: "secret2", Hosts: []string{"bar.example"}}}
	ingress2.ObjectMeta.Name = "ingress2"

	t.Log("configuring secrets")
	fooExampleTLSCert, fooExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(
		certificate.WithCommonName("secure-foo-bar"), certificate.WithDNSNames("secure-foo-bar", "foo.example"),
	)
	barExampleTLSCert, barExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(
		certificate.WithCommonName("foo.com"), certificate.WithDNSNames("foo.com", "bar.example"),
	)

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Name:      "secret1",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": fooExampleTLSCert,
				"tls.key": fooExampleTLSKey,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e5"),
				Name:      "secret2",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": barExampleTLSCert,
				"tls.key": barExampleTLSKey,
			},
		},
	}

	// Since we updated the logic of secret controller to only process secrets that are referred by
	// other controlled objects (service, ingress, gateway, ...), we should make sure that ingresses
	// created before and after referred secret created both works.
	// so here we interleave the creating process of deploying 2 ingresses and secrets.
	t.Log("deploying secrets and ingresses")
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress1))
	cleaner.Add(ingress1)

	secret1, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secrets[0], metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(secret1)

	secret2, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secrets[1], metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(secret2)

	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress2))
	cleaner.Add(ingress2)

	t.Log("checking first ingress status readiness")
	require.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ingress1)
		if err != nil {
			return false
		}
		for _, ingress := range lbstatus.Ingress {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("networkingv1 ingress1 hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Log("checking second ingress status readiness")
	assert.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ingress2)
		if err != nil {
			return false
		}
		for _, ingress := range lbstatus.Ingress {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("networkingv1 ingress2 hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Log("waiting for routes from Ingress to be operational with expected certificate")
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://foo.example:443/foo")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://foo.example:443/foo: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>") && resp.TLS.PeerCertificates[0].Subject.CommonName == "secure-foo-bar"
		}
		return false
	}, ingressWait, waitTick, true)

	t.Log("waiting for routes from Ingress to be operational with expected certificate")
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://bar.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/bar: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>") && resp.TLS.PeerCertificates[0].Subject.CommonName == "foo.com"
		}
		return false
	}, ingressWait, waitTick, true)

	// This should work currently. generators.NewIngressForService() only creates path rules by default, so while we don't
	// do anything for baz.example other than add fake DNS for it, the /bar still routes it through ingress2's route.
	// We're going to break it later, but need to confirm it does work first.
	t.Log("confirm Ingress path routes available on other hostnames")
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://baz.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/baz: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	ingress2, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress2.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	ingress2.ObjectMeta.Annotations["konghq.com/snis"] = "bar.example"
	_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress2, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Log("confirm Ingress no longer routes without matching SNI")
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://baz.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://baz.example:443/bar: %v", err)
			return false
		}

		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)

	t.Log("confirm Ingress still routes with matching SNI")
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://bar.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/bar: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)
}
