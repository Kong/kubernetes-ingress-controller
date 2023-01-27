//go:build envtest
// +build envtest

package manager

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestEndpointsliceForService(t *testing.T) {
	ctx := context.Background()

	testEnv := &envtest.Environment{}

	cfg, err := testEnv.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, testEnv.Stop())
	})

	client, err := ctrlclient.New(cfg, ctrlclient.Options{})
	require.NoError(t, err)

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
	}
	require.NoError(t, client.Create(ctx, &ns, &ctrlclient.CreateOptions{}))

	service := types.NamespacedName{
		Name: "kong-admin",
	}

	endpoints := discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      uuid.NewString(),
			Namespace: ns.Name,
			Labels: map[string]string{
				"kubernetes.io/service-name": service.Name,
			},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"10.0.0.1"},
				Conditions: discoveryv1.EndpointConditions{
					Ready: lo.ToPtr(true),
				},
			},
		},
		Ports: []discoveryv1.EndpointPort{
			{
				Name: lo.ToPtr("admin"),
				Port: lo.ToPtr(int32(8080)),
			},
		},
	}
	require.NoError(t, client.Create(ctx, &endpoints, &ctrlclient.CreateOptions{}))

	require.NoError(t, client.Get(ctx, types.NamespacedName{Name: endpoints.Name, Namespace: endpoints.Namespace}, &endpoints))

	addresses, err := manager.GetEndpointslicesForService(ctx, client, service)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://10.0.0.1:8080"}, addresses)
}
