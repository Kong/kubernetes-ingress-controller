package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCaching(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("configuring and starting a new proxy server")
	proxyInterface, err := NewCacheBasedProxy(ctx, logger, fakeK8sClient, fakeKongConfig, "kongtests", false)
	assert.NoError(t, err)

	t.Log("ensuring the integrity of the proxy server")
	proxy, ok := proxyInterface.(*clientgoCachedProxyResolver)
	assert.True(t, ok)
	assert.NotNil(t, proxy.cache)

	t.Log("intentionally freezing async updates to inspect cache state during tests")
	proxy.syncTicker.Reset(time.Minute * 1)

	t.Log("generating 10 new objects to the proxy cache server")
	testObjects := make([]client.Object, 10)
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		deployment := k8s.NewDeploymentForContainer(k8s.NewContainer(name, name, 8080))
		service := k8s.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
		ingress := k8s.NewIngressForService("/testing", nil, service)
		testObjects[i] = ingress
	}

	t.Logf("adding %d new objects to the proxy cache server", len(testObjects))
	assert.Len(t, proxy.cache.IngressV1.List(), 0)
	for _, testObject := range testObjects {
		proxy.UpdateObject(testObject)
	}

	t.Log("ensuring the consistency of the underlying object cache (that objects were added properly)")
	assert.Eventually(t, func() bool {
		return len(proxy.cache.IngressV1.List()) == len(testObjects)
	}, time.Second*10, time.Millisecond*200)

	t.Log("flushing the cache state to kong admin api")
	previousUpdateCount := fakeKongAdminUpdateCount()
	proxy.syncTicker.Reset(time.Millisecond * 200)
	assert.Eventually(t, func() bool { return fakeKongAdminUpdateCount() == previousUpdateCount+1 }, time.Second*10, time.Millisecond*200)
}
