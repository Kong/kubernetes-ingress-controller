//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy/seeder"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestProxySeedRound(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("proxy seed functionality does not exist in pre-V2, skipping")
	}

	t.Log("configuration testing environment for proxy seed")
	_ = proxyReady()
	namespace := uuid.NewString()
	ingressClassName := uuid.NewString()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Logf("creating namespace %s to run the tests in", namespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	ns, err := cluster.Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testTCPIngressNamespace)
		assert.NoError(t, cluster.Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
	}()

	t.Log("creating a simple HTTP server container to test proxy seeding")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing service %s via Ingress using Ingress Class %s", service.Name, ingressClassName)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClassName,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses(namespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("seed a custom proxy server to verify object updates %s", ingress.Name)
	fakePrx := &fakeProxy{}
	seeder, err := seeder.NewBuilder(mgr.GetConfig(), fakePrx).
		WithFieldLogger(logrus.New().WithField("component", "integration_tests")).
		WithIngressClass(ingressClassName).
		Build()
	require.NoError(t, err)
	require.NoError(t, seeder.Seed(ctx))

	t.Log("verifying that the seeded ingress, its service and endpoints were all seeded properly into the proxy cache")
	require.True(t, len(fakePrx.objs) >= 3)
}

// -----------------------------------------------------------------------------
// Fake Proxy - Implementation For Tests
// -----------------------------------------------------------------------------

type fakeProxy struct {
	objs []client.Object
}

func (f *fakeProxy) UpdateObjects(objs ...client.Object) error {
	f.objs = append(f.objs, objs...)
	return nil
}

func (f *fakeProxy) DeleteObjects(deleteObjs ...client.Object) error {
	return fmt.Errorf("unimplemented")
}
