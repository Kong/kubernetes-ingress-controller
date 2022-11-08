//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blang/semver/v4"
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
	"k8s.io/client-go/util/retry"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

// extraIngressNamespace is the name of an alternative namespace used for ingress tests.
const extraIngressNamespace = "elsewhere"

var (
	statusWait        = time.Minute * 3
	ingressClassMutex = sync.Mutex{}
)

func TestIngressEssentials(t *testing.T) {
	t.Parallel()
	t.Log("locking IngressClass management")
	ingressClassMutex.Lock()
	defer func() {
		t.Log("unlocking IngressClass management")
		ingressClassMutex.Unlock()
	}()
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/test_ingress_essentials", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	addIngressToCleaner(cleaner, ingress)

	t.Log("waiting for updated ingress status to include IP")
	require.Eventually(t, func() bool {
		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, env.Cluster(), ns.Name, ingress)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, statusWait, waitTick)

	t.Log("waiting for routes from Ingress to be operational")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_essentials", proxyURL))
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
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
			_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err
		})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			delete(ingress.ObjectMeta.Annotations, annotations.IngressClassKey)
			_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err
		})
		require.NoError(t, err)
	}

	t.Logf("verifying that removing the ingress.class annotation %q from ingress causes routes to disconnect", ingressClass)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_essentials", proxyURL))
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
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
			_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err
		})
		require.NoError(t, err)
	case *netv1beta1.Ingress:
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			ingress.ObjectMeta.Annotations[annotations.IngressClassKey] = ingressClass
			_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err
		})
		require.NoError(t, err)
	}

	t.Log("waiting for routes from Ingress to be operational after reintroducing ingress class annotation")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_essentials", proxyURL))
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
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_essentials", proxyURL))
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
	}, statusWait, waitTick)

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
	t.Log("locking IngressClass management")
	ingressClassMutex.Lock()
	defer func() {
		t.Log("unlocking IngressClass management")
		ingressClassMutex.Unlock()
	}()
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	if clusterVersion.Major < uint64(2) && clusterVersion.Minor < uint64(19) {
		t.Skip("ingress spec tests can not be properly validated against old clusters")
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes using the IngressClassName spec")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/test_ingressclassname_spec/", map[string]string{"konghq.com/strip-path": "true"}, service)
	switch obj := ingress.(type) {
	case *netv1.Ingress:
		obj.Spec.IngressClassName = kong.String(ingressClass)
	case *netv1beta1.Ingress:
		obj.Spec.IngressClassName = kong.String(ingressClass)
	}
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	addIngressToCleaner(cleaner, ingress)

	t.Log("waiting for routes from Ingress to be operational")
	defer func() {
		if t.Failed() {
			resp, err := httpc.Get(fmt.Sprintf("%s/test_ingressclassname_spec", proxyURL))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			}
			t.Logf("TestIngressClassNameSpec failed, current GET %s/test_ingressclassname_spec status code is %d",
				proxyURL, resp.StatusCode)
			resp.Body.Close()
		}
	}()

	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingressclassname_spec", proxyURL))
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
	err = setIngressClassNameWithRetry(ctx, ns.Name, ingress, nil)
	require.NoError(t, err)

	t.Logf("verifying that removing the IngressClassName %q from ingress causes routes to disconnect", ingressClass)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingressclassname_spec", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)

	t.Logf("putting the IngressClassName %q back on ingress", ingressClass)
	err = setIngressClassNameWithRetry(ctx, ns.Name, ingress, kong.String(ingressClass))
	require.NoError(t, err)

	t.Log("waiting for routes from Ingress to be operational after reintroducing ingress class annotation")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingressclassname_spec", proxyURL))
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
	}, ingressWait, waitTick, // TODO: dump status of kong gateway here.
	)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingressclassname_spec", proxyURL))
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
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, elsewhereNamespace, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("cleaning up namespace %s", elsewhereNamespace.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, elsewhereNamespace.Name, metav1.DeleteOptions{}))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	elsewhereDeployment, err := env.Cluster().Client().AppsV1().Deployments(extraIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", elsewhereDeployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(extraIngressNamespace).Delete(ctx, elsewhereDeployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(extraIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(extraIngressNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	elsewhereIngress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/elsewhere", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), extraIngressNamespace, elsewhereIngress))

	defer func() {
		t.Log("ensuring that Ingress resources are cleaned up")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), extraIngressNamespace, elsewhereIngress); err != nil {
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
}

func TestIngressStatusUpdatesExtended(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	if clusterVersion.Major == uint64(1) && clusterVersion.Minor < uint64(19) {
		t.Skip("status test disabled for old cluster versions")
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
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
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
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
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
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
	}, statusWait, waitTick)

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
	}, statusWait, waitTick)
}

// TestIngressClassRegexToggle tests if the controller adds the 3.x "~" regular expression path prefix to Ingress
// paths that match the 2.x heuristic when their IngressClass has the EnableLegacyRegexDetection flag set. It IS NOT
// parallel: parts of the test may add this route _without_ the prefix, and the 3.x router really hates this and will
// stop working altogether.
func TestIngressClassRegexToggle(t *testing.T) {
	// the manager runs in a goroutine and may not have pulled the version before this test starts
	require.Eventually(t, func() bool {
		return !versions.GetKongVersion().Full().EQ(semver.MustParse("0.0.0"))
	}, time.Minute, time.Second)
	if !versions.GetKongVersion().MajorOnly().GTE(versions.ExplicitRegexPathVersionCutoff) {
		t.Skip("legacy regex detection is only relevant for Kong 3.0+")
	}

	// skip the test if the cluster does not support namespaced ingress class parameter (<=1.21).
	// since 1.21 is End of Life now.
	namespacedIngressClassParameterMinKubernetesVersion := semver.MustParse("1.22.0")
	if clusterVersion.LT(namespacedIngressClassParameterMinKubernetesVersion) {
		t.Skipf("kubernetes cluster version %s does not support namespaced ingress class parameters", clusterVersion.String())
	}

	t.Log("locking IngressClass management")
	ingressClassMutex.Lock()
	defer func() {
		t.Log("unlocking IngressClass management")
		ingressClassMutex.Unlock()
	}()
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an IngressClassParameters with legacy regex detection enabled")
	params := &v1alpha1.IngressClassParameters{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressClass,
		},
		Spec: v1alpha1.IngressClassParametersSpec{
			EnableLegacyRegexDetection: true,
		},
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	params, err = c.ConfigurationV1alpha1().IngressClassParameterses(ns.Name).Create(ctx, params, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(params)

	class, err := env.Cluster().Client().NetworkingV1().IngressClasses().Get(ctx, ingressClass, metav1.GetOptions{})
	require.NoError(t, err)
	t.Logf("adding legacy regex IngressClassParameters to the %s IngressClass", class.Name)
	class.Spec.Parameters = &netv1.IngressClassParametersReference{
		APIGroup:  &v1alpha1.GroupVersion.Group,
		Kind:      v1alpha1.IngressClassParametersKind,
		Name:      params.Name,
		Scope:     kong.String(netv1.IngressClassParametersReferenceScopeNamespace),
		Namespace: &params.Namespace,
	}
	_, err = env.Cluster().Client().NetworkingV1().IngressClasses().Update(ctx, class, metav1.UpdateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("removing parameters from IngressClass %s", class.ObjectMeta.Name)
		class, err := env.Cluster().Client().NetworkingV1().IngressClasses().Get(ctx, ingressClass, metav1.GetOptions{})
		require.NoError(t, err)
		class.Spec.Parameters = nil
		_, err = env.Cluster().Client().NetworkingV1().IngressClasses().Update(ctx, class, metav1.UpdateOptions{})
		require.NoError(t, err)
	}()

	t.Logf("creating an ingress for service %s", service.Name)
	require.NoError(t, err)
	pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-toggle",
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     `/test_ingress_class_regex_toggle/\d+`,
									PathType: &pathTypeImplementationSpecific,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	// we only test the positive case here, and assume the negative case (without the toggle, this will not route)
	// based on prior experience. unfortunately the effect of the negative case is that it breaks router rebuilds
	// entirely, which would be bad for other tests.
	t.Log("waiting for ingress path to become available")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_class_regex_toggle/999", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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

func TestIngressRegexPrefix(t *testing.T) {
	if !versions.GetKongVersion().MajorOnly().GTE(versions.ExplicitRegexPathVersionCutoff) {
		t.Skip("regex prefixes are only relevant for Kong 3.0+")
	}
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s", service.Name)
	require.NoError(t, err)
	pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-prefix",
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     `/~/test_ingress_regex_prefix/\d+`,
									PathType: &pathTypeImplementationSpecific,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Logf("creating an ingress with a non-default prefix for service %s", service.Name)
	require.NoError(t, err)
	ingressMod := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-prefix-ns",
			Annotations: map[string]string{
				"konghq.com/strip-path":   "true",
				"konghq.com/regex-prefix": "/@",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     `/@/test_ingress_regex_prefix_nonstandard/\d+`,
									PathType: &pathTypeImplementationSpecific,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
								{
									Path:     `/~/test_ingress_regex_prefix_nonstandard_default`,
									PathType: &pathTypeImplementationSpecific,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingressMod))
	cleaner.Add(ingressMod)

	t.Log("waiting for ingress path to become available")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_regex_prefix/999", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/test_ingress_regex_prefix_nonstandard/999", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/~/test_ingress_regex_prefix_nonstandard_default", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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

func TestIngressRecoverFromInvalidPath(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, test ns %s, dumped diagnostics to %s", t.Name(), ns.Name, output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	// TODO: run this separately, make it not to affect other tests for sharing Kong.
	if !runInvalidConfigTests {
		t.Skipf("the case %s should be run separately; please set TEST_RUN_INVALID_CONFIG_CASES to true to run this case", t.Name())
	}

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Log("create an ingress")
	pathTypePrefix := netv1.PathTypePrefix
	pathTypeImplementationSpecific := netv1.PathTypeImplementationSpecific
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-prefix-ns",
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									PathType: &pathTypePrefix,
									Path:     "/foo",
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Create(ctx, ingress, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(ingress)

	t.Log("waiting for ingress path to become available")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/foo/", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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

	t.Log("add an invalid path to ingress")
	ingressInvalid := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-prefix-ns",
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									PathType: &pathTypePrefix,
									Path:     "/bar",
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
								{
									PathType: &pathTypeImplementationSpecific,
									Path:     `/~^^/*$`, // invalid regex
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingressInvalid, metav1.UpdateOptions{})
	require.NoError(t, err)
	t.Log("verifying new configuration is not applied to kong proxy")
	require.Never(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/bar/", proxyURL))
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, time.Minute, waitTick)

	t.Log("verifying routes configured before invalid config is still available")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/foo/", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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

	t.Log("reconfigure ingress with valid paths")
	ingressRecover := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "regex-prefix-ns",
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(ingressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									PathType: &pathTypePrefix,
									Path:     "/bar",
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: service.Spec.Ports[0].Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingressRecover, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("waiting for ingress path to recover and new path available")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/bar/", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
