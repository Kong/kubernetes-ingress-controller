//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalIngress(t *testing.T) {
	ctx := context.Background()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err := cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
	p := proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("removing the ingress.class annotation %q from ingress %s", ingressClass, ingress.Name)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("verifying that removing the ingress.class annotation %q from ingress %s causes routes to disconnect", ingressClass, ingress.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the ingress.class annotation %q back on ingress %s", ingressClass, ingress.Name)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("waiting for routes from Ingress %s to be operational after reintroducing ingress class annotation", ingress.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	assert.NoError(t, cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)
}

func TestHTTPSRedirect(t *testing.T) {
	ctx := context.Background()
	opts := metav1.CreateOptions{}

	t.Log("creating an HTTP container via deployment to test redirect functionality")
	container := k8sgen.NewContainer("alsohttpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	_, err := cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via Service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	service, err = cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing Service %s via Ingress", service.Name)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey:             ingressClass,
		"konghq.com/protocols":                  "https",
		"konghq.com/https-redirect-status-code": "301",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Ingress %s", ingress.Name)
		assert.NoError(t, cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("waiting for Ingress %s to be operational and properly redirect", ingress.Name)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}
	assert.Eventually(t, func() bool {
		p := proxyReady()
		resp, err := client.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusMovedPermanently
	}, ingressWait, waitTick)
}

func TestIngressClassNameSpec(t *testing.T) {
	ctx := context.Background()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes using the IngressClassName spec")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err := cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{"konghq.com/strip-path": "true"}, service)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
	p := proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("removing the IngressClassName %q from ingress %s", ingressClass, ingress.Name)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	ingress.Spec.IngressClassName = nil
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("verifying that removing the IngressClassName %q from ingress %s causes routes to disconnect", ingressClass, ingress.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the IngressClassName %q back on ingress %s", ingressClass, ingress.Name)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("waiting for routes from Ingress %s to be operational after reintroducing ingress class annotation", ingress.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	assert.NoError(t, cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)

}

func TestIngressNamespaces(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("support for distinct namespace watches is not supported in legacy KIC")
	}
	ctx := context.Background()

	// ensure the alternative namespace is created
	elsewhereNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: elsewhere}}
	nowhere := "nowhere"
	nowhereNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nowhere}}
	_, err := cluster.Client().CoreV1().Namespaces().Create(ctx, nowhereNamespace, metav1.CreateOptions{})
	assert.NoError(t, err)
	_, err = cluster.Client().CoreV1().Namespaces().Create(ctx, elsewhereNamespace, metav1.CreateOptions{})
	assert.NoError(t, err)
	defer func() {
		t.Logf("cleaning up namespace %s", elsewhereNamespace.Name)
		assert.NoError(t, cluster.Client().CoreV1().Namespaces().Delete(ctx, elsewhereNamespace.Name, metav1.DeleteOptions{}))
		t.Logf("cleaning up namespace %s", nowhereNamespace.Name)
		assert.NoError(t, cluster.Client().CoreV1().Namespaces().Delete(ctx, nowhereNamespace.Name, metav1.DeleteOptions{}))
	}()
	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	elsewhereDeployment, err := cluster.Client().AppsV1().Deployments(elsewhere).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)
	nowhereDeployment, err := cluster.Client().AppsV1().Deployments(nowhere).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", elsewhereDeployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(elsewhere).Delete(ctx, elsewhereDeployment.Name, metav1.DeleteOptions{}))
		assert.NoError(t, cluster.Client().AppsV1().Deployments(nowhere).Delete(ctx, nowhereDeployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services(elsewhere).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	_, err = cluster.Client().CoreV1().Services(nowhere).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(elsewhere).Delete(ctx, service.Name, metav1.DeleteOptions{}))
		assert.NoError(t, cluster.Client().CoreV1().Services(nowhere).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	elsewhereIngress := k8sgen.NewIngressForService("/elsewhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	nowhereIngress := k8sgen.NewIngressForService("/nowhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	elsewhereIngress, err = cluster.Client().NetworkingV1().Ingresses(elsewhere).Create(ctx, elsewhereIngress, metav1.CreateOptions{})
	assert.NoError(t, err)
	nowhereIngress, err = cluster.Client().NetworkingV1().Ingresses(nowhere).Create(ctx, nowhereIngress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", elsewhereIngress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses(elsewhere).Delete(ctx, elsewhereIngress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		if err := cluster.Client().NetworkingV1().Ingresses(nowhere).Delete(ctx, nowhereIngress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", elsewhereIngress.Name)
	p := proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/elsewhere", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/nowhere", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)
}
