package proxy

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func Test_FetchCustomEntities(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid secret",
			args: args{
				secret: "default/validCustomEntities",
			},
			want:    []byte("carp"),
			wantErr: true,
		},
		{
			name: "incorrect name format",
			args: args{
				secret: "!",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-existent secret",
			args: args{
				secret: "default/nope",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "secret lacks config key",
			args: args{
				secret: "default/invalidCustomEntities",
			},
			want:    nil,
			wantErr: true,
		},
	}
	store, err := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "validCustomEntities",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"config": []byte("carp"),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalidCustomEntities",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"ohno": []byte("carp"),
				},
			},
		},
	})
	assert.Nil(err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchCustomEntities(tt.args.secret, store)
			if err != nil && !tt.wantErr {
				t.Errorf("kongPluginFromK8SClusterPlugin error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(tt.want, got)
		})
	}
}

func TestIsReady(t *testing.T) {
	fakePostgresKong := sendconfig.Kong{InMemory: false}
	fakeDblessKong := sendconfig.Kong{InMemory: true}
	postgresProxy := clientgoCachedProxyResolver{
		kongConfig: fakePostgresKong,
		dbmode:     "postgres",
	}
	dblessProxy := clientgoCachedProxyResolver{
		kongConfig: fakeDblessKong,
		dbmode:     "off",
	}

	t.Log("checking initial readiness state")
	assert.True(t, postgresProxy.IsReady())
	assert.False(t, dblessProxy.IsReady())

	t.Log("marking config applied and checking readiness after")
	postgresProxy.markConfigApplied()
	dblessProxy.markConfigApplied()
	assert.True(t, postgresProxy.IsReady())
	assert.True(t, dblessProxy.IsReady())
}

func TestCaching(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("configuring and starting a new proxy server")
	proxyInterface, err := NewCacheBasedProxy(ctx, logger, fakeK8sClient, fakeKongConfig,
		"kongtests", false, mockKongAdmin, util.ConfigDumpDiagnostic{}, time.Millisecond*300)
	assert.NoError(t, err)

	t.Log("ensuring the integrity of the proxy server")
	proxy, ok := proxyInterface.(*clientgoCachedProxyResolver)
	assert.True(t, ok)
	assert.NotNil(t, proxy.cache)

	t.Log("intentionally freezing async updates to inspect cache state during tests")
	proxy.syncTicker.Reset(time.Minute * 3)

	t.Log("generating 10 new objects to the proxy cache server")
	testObjects := make([]client.Object, 10)
	for i := 0; i < 10; i++ {
		name := uuid.NewString()
		deployment := generators.NewDeploymentForContainer(generators.NewContainer(name, name, 8080))
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
		ingress := generators.NewIngressForService("/testing", nil, service)
		testObjects[i] = ingress
	}

	t.Logf("adding %d new objects to the proxy cache server", len(testObjects))
	assert.Len(t, proxy.cache.IngressV1.List(), 0)
	for _, testObject := range testObjects {
		require.NoError(t, proxy.UpdateObject(testObject))
	}

	t.Log("ensuring the consistency of the underlying object cache (that objects were added properly)")
	assert.Eventually(t, func() bool {
		return len(proxy.cache.IngressV1.List()) == len(testObjects)
	}, time.Second*10, time.Millisecond*200)

	t.Log("verifying the integrity of the object cache items")
	matches := 0
	for _, testObj := range testObjects {
		for _, obj := range proxy.cache.IngressV1.List() {
			ing, ok := obj.(*netv1.Ingress)
			require.True(t, ok)
			if ing.Namespace == testObj.GetNamespace() && ing.Name == testObj.GetName() {
				matches++
			}
		}
	}
	require.Equal(t, len(testObjects), matches)
}

func TestProxyTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Log("configuring and starting a new proxy server")

	// mock the next Admin API response (which will be / to get the root config) to ensure
	// that it takes longer than the timeout we will set, in order to trigger the timeout.
	fakeKongAdminAPI.MockNextResponse(kong.AdminAPIResponse{
		Status:   http.StatusGatewayTimeout,
		Body:     []byte{},
		Callback: func() { time.Sleep(time.Millisecond * 30) },
	})

	// the timeout is shorter than the wait time for the http response, we should expect
	// to see the the context deadline for the http response triggered.
	timeout := time.Millisecond * 10

	_, err := NewCacheBasedProxy(ctx, logger, fakeK8sClient, fakeKongConfig, "kongtests", false, mockKongAdmin, util.ConfigDumpDiagnostic{}, timeout)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
