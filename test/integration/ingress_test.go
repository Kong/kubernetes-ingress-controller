//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
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
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

// extraIngressNamespace is the name of an alternative namespace used for ingress tests
const extraIngressNamespace = "elsewhere"

func TestIngressEssentials(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))

	defer func() {
		t.Log("cleaning up Ingress resource")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for updated ingress status to include IP")
	require.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ingress)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, ingressWait, waitTick)

	t.Log("waiting for routes from Ingress to be operational")
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

	t.Logf("removing the ingress.class annotation %q from ingress", ingressClass)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
		_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	t.Logf("verifying that removing the ingress.class annotation %q from ingress causes routes to disconnect", ingressClass)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the ingress.class annotation %q back on ingress", ingressClass)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
		_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	t.Log("waiting for routes from Ingress to be operational after reintroducing ingress class annotation")
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

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
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

func TestGRPCIngressEssentials(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("grpcbin", "moul/grpcbin", 9001)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	// as of KTF 0.9,0 NewServiceForDeployment doesn't initialize annotations itself, need to do it outside
	service.ObjectMeta.Annotations = map[string]string{annotations.AnnotationPrefix + annotations.ProtocolKey: "grpc"}
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/", map[string]string{
		annotations.IngressClassKey:                             ingressClass,
		annotations.AnnotationPrefix + annotations.ProtocolsKey: "grpc,grpcs",
		annotations.AnnotationPrefix + annotations.StripPathKey: "false",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))

	defer func() {
		t.Log("cleaning up Ingress resource")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for updated ingress status to include IP")
	require.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ingress)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, ingressWait, waitTick)

	// So far this only tests that the ingress is created and receives status information, to confirm the fix for
	// https://github.com/Kong/kubernetes-ingress-controller/issues/1991
	// It does not test routing, though the status implementation implies it (we only add status after we confirm
	// configuration is present in the proxy). This test could be expanded to better confirm routing with a suitable
	// gRPC test client.

	t.Log("deleting Ingress")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
}

func TestIngressClassNameSpec(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	if clusterVersion.Major < uint64(2) && clusterVersion.Minor < uint64(19) {
		t.Skip("ingress spec tests can not be properly validated against old clusters")
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes using the IngressClassName spec")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/httpbin", map[string]string{"konghq.com/strip-path": "true"}, service)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		obj.Spec.IngressClassName = kong.String(ingressClass)
	case *netv1beta1.Ingress:
		obj.Spec.IngressClassName = kong.String(ingressClass)
	}
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))

	defer func() {
		t.Log("ensuring that Ingress is cleaned up")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for routes from Ingress to be operational")
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

	t.Logf("removing the IngressClassName %q from ingress", ingressClass)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.Spec.IngressClassName = nil
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.Spec.IngressClassName = nil
		_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	t.Logf("verifying that removing the IngressClassName %q from ingress causes routes to disconnect", ingressClass)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the IngressClassName %q back on ingress", ingressClass)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.Spec.IngressClassName = kong.String(ingressClass)
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ingress.Spec.IngressClassName = kong.String(ingressClass)
		_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	t.Log("waiting for routes from Ingress to be operational after reintroducing ingress class annotation")
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

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
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
	t.Parallel()

	t.Log("creating extra testing namespaces")
	elsewhereNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: extraIngressNamespace}}
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
	elsewhereDeployment, err := env.Cluster().Client().AppsV1().Deployments(extraIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	nowhereDeployment, err := env.Cluster().Client().AppsV1().Deployments(nowhere).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", elsewhereDeployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(extraIngressNamespace).Delete(ctx, elsewhereDeployment.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(nowhere).Delete(ctx, nowhereDeployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(extraIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = env.Cluster().Client().CoreV1().Services(nowhere).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(extraIngressNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(nowhere).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	elsewhereIngress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/elsewhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	nowhereIngress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/nowhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), extraIngressNamespace, elsewhereIngress))
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), nowhere, nowhereIngress))

	defer func() {
		t.Log("ensuring that Ingress resources are cleaned up")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), extraIngressNamespace, elsewhereIngress); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		if err := clusters.DeleteIngress(ctx, env.Cluster(), nowhere, nowhereIngress); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for routes from Ingress to be operational")
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

func TestIngressStatusUpdatesExtended(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	if clusterVersion.Major == uint64(1) && clusterVersion.Minor < uint64(19) {
		t.Skip("status test disabled for old cluster versions")
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Log("creating a variety of Ingress resources for the service to verify status updates")
	pathType := netv1.PathTypePrefix
	ingNameWithPeriods := "status.check.with.periods"
	ingNameExtraForService := "statuscheck1"
	for _, ing := range []*netv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: ingNameWithPeriods,
				Annotations: map[string]string{
					annotations.IngressClassKey: ingressClass,
					"konghq.com/strip-path":     "true",
				},
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     "/statuscheck2",
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: service.Name,
												Port: netv1.ServiceBackendPort{
													Number: service.Spec.Ports[0].Port,
												}}}}}}}}}},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: ingNameExtraForService,
				Annotations: map[string]string{
					annotations.IngressClassKey: ingressClass,
					"konghq.com/strip-path":     "true",
				},
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     "/statuscheck1",
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: service.Name,
												Port: netv1.ServiceBackendPort{
													Number: service.Spec.Ports[0].Port,
												}}}}}}}}}},
		},
	} {
		t.Logf("creating ingress %s and verifying status updates", ing.Name)
		createdIngress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Create(ctx, ing, metav1.CreateOptions{})
		require.NoError(t, err)

		defer func() {
			t.Logf("cleaning up ingress %s", createdIngress.Name)
			if err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Delete(ctx, createdIngress.Name, metav1.DeleteOptions{}); err != nil {
				if !errors.IsNotFound(err) {
					assert.NoError(t, err)
				}
			}
		}()
	}

	t.Log("verifying that an ingress with periods in the name has its status populated")
	require.Eventually(t, func() bool {
		ing, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingNameWithPeriods, metav1.GetOptions{})
		if err != nil {
			return false
		}
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ing)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, time.Minute, time.Second)

	t.Log("verifying that when a service has more than one ingress, the status updates for those beyond the first")
	require.Eventually(t, func() bool {
		ing, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingNameExtraForService, metav1.GetOptions{})
		if err != nil {
			return false
		}
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ing)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, time.Minute, time.Second)
}
