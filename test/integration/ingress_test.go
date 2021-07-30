//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
)

var testIngressEssentialsNamespace = "ingress-essentials"
var testIngressClassNameSpecNamespace = "ingress-class-name-spec"

func TestIngressEssentials(t *testing.T) {
	ctx := context.Background()

	t.Logf("creating namespace %s for testing", testIngressEssentialsNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testIngressEssentialsNamespace}}
	ns, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testIngressEssentialsNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(testIngressEssentialsNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(testIngressEssentialsNamespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(testIngressEssentialsNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(testIngressEssentialsNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Create(ctx, ingress, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for updated ingress status to include IP")
	require.Eventually(t, func() bool {
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return len(ingress.Status.LoadBalancer.Ingress) > 0
	}, ingressWait, waitTick)

	t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("removing the ingress.class annotation %q from ingress %s", ingressClass, ingress.Name)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Get(ctx, ingress.Name, metav1.GetOptions{})
	require.NoError(t, err)
	delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Update(ctx, ingress, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Logf("verifying that removing the ingress.class annotation %q from ingress %s causes routes to disconnect", ingressClass, ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the ingress.class annotation %q back on ingress %s", ingressClass, ingress.Name)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Get(ctx, ingress.Name, metav1.GetOptions{})
	require.NoError(t, err)
	ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Update(ctx, ingress, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Logf("waiting for routes from Ingress %s to be operational after reintroducing ingress class annotation", ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	require.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(testIngressEssentialsNamespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)
}

func TestIngressClassNameSpec(t *testing.T) {
	ctx := context.Background()

	t.Logf("creating namespace %s for testing", testIngressClassNameSpecNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testIngressClassNameSpecNamespace}}
	ns, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testIngressClassNameSpecNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes using the IngressClassName spec")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(testIngressClassNameSpecNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(testIngressClassNameSpecNamespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(testIngressClassNameSpecNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(testIngressClassNameSpecNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{"konghq.com/strip-path": "true"}, service)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Create(ctx, ingress, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("removing the IngressClassName %q from ingress %s", ingressClass, ingress.Name)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Get(ctx, ingress.Name, metav1.GetOptions{})
	require.NoError(t, err)
	ingress.Spec.IngressClassName = nil
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Update(ctx, ingress, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Logf("verifying that removing the IngressClassName %q from ingress %s causes routes to disconnect", ingressClass, ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the IngressClassName %q back on ingress %s", ingressClass, ingress.Name)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Get(ctx, ingress.Name, metav1.GetOptions{})
	require.NoError(t, err)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Update(ctx, ingress, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Logf("waiting for routes from Ingress %s to be operational after reintroducing ingress class annotation", ingress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	require.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(testIngressClassNameSpecNamespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

}

func TestIngressNamespaces(t *testing.T) {
	ctx := context.Background()

	// ensure the alternative namespace is created
	elsewhereNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: elsewhere}}
	nowhere := "nowhere"
	nowhereNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nowhere}}
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, nowhereNamespace, metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = env.Cluster().Client().CoreV1().Namespaces().Create(ctx, elsewhereNamespace, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("cleaning up namespace %s", elsewhereNamespace.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, elsewhereNamespace.Name, metav1.DeleteOptions{}))
		t.Logf("cleaning up namespace %s", nowhereNamespace.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, nowhereNamespace.Name, metav1.DeleteOptions{}))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	elsewhereDeployment, err := env.Cluster().Client().AppsV1().Deployments(elsewhere).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	nowhereDeployment, err := env.Cluster().Client().AppsV1().Deployments(nowhere).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", elsewhereDeployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(elsewhere).Delete(ctx, elsewhereDeployment.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(nowhere).Delete(ctx, nowhereDeployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(elsewhere).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = env.Cluster().Client().CoreV1().Services(nowhere).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(elsewhere).Delete(ctx, service.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(nowhere).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	elsewhereIngress := generators.NewIngressForService("/elsewhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	nowhereIngress := generators.NewIngressForService("/nowhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	elsewhereIngress, err = env.Cluster().Client().NetworkingV1().Ingresses(elsewhere).Create(ctx, elsewhereIngress, metav1.CreateOptions{})
	require.NoError(t, err)
	nowhereIngress, err = env.Cluster().Client().NetworkingV1().Ingresses(nowhere).Create(ctx, nowhereIngress, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", elsewhereIngress.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(elsewhere).Delete(ctx, elsewhereIngress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		if err := env.Cluster().Client().NetworkingV1().Ingresses(nowhere).Delete(ctx, nowhereIngress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", elsewhereIngress.Name)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/elsewhere", proxyURL))
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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/nowhere", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)
}
