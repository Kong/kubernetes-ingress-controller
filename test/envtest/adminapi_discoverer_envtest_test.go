//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config/types"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestDiscoverer_GetAdminAPIsForServiceReturnsAllAddressesCorrectlyPagingThroughResults(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t)
	cfg := Setup(t, scheme)
	client := NewControllerClient(t, scheme, cfg)

	// In tests below we use a deferred cancel to stop the manager and not wait
	// for its timeout.

	testcases := []struct {
		subnetC int
		subnetD int
	}{
		{subnetC: 1, subnetD: 100},
		{subnetC: 1, subnetD: 101},
		{subnetC: 1, subnetD: 250},
		{subnetC: 2, subnetD: 250},
		{subnetC: 5, subnetD: 250},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%dx%d", tc.subnetC, tc.subnetD), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var (
				ns          = CreateNamespace(ctx, t, client)
				serviceName = uuid.NewString()
				service     = k8stypes.NamespacedName{
					Namespace: ns.Name,
					Name:      serviceName,
				}
			)

			for i := 0; i < tc.subnetC; i++ {
				for j := 0; j < tc.subnetD; j++ {
					es := discoveryv1.EndpointSlice{
						ObjectMeta: metav1.ObjectMeta{
							Name:      uuid.NewString(),
							Namespace: ns.Name,
							Labels: map[string]string{
								"kubernetes.io/service-name": serviceName,
							},
						},
						AddressType: discoveryv1.AddressTypeIPv4,
						Endpoints: []discoveryv1.Endpoint{
							{
								Addresses: []string{fmt.Sprintf("10.0.%d.%d", i, j)},
								Conditions: discoveryv1.EndpointConditions{
									Ready:       lo.ToPtr(true),
									Terminating: lo.ToPtr(false),
								},
								TargetRef: testPodReference("pod-1", ns.Name),
							},
						},
						Ports: builder.NewEndpointPort(8444).WithName("admin").IntoSlice(),
					}
					require.NoError(t, client.Create(ctx, &es))
				}
			}

			discoverer, err := adminapi.NewDiscoverer(sets.New("admin"), cfgtypes.IPDNSStrategy)
			require.NoError(t, err)

			got, err := discoverer.GetAdminAPIsForService(ctx, client, service)
			require.NoError(t, err)
			require.Len(t, got, tc.subnetD*tc.subnetC, "GetAdminAPIsForService should return all valid addresses")
		})
	}
}

func testPodReference(name, ns string) *corev1.ObjectReference {
	return &corev1.ObjectReference{
		Kind:      "Pod",
		Namespace: ns,
		Name:      name,
	}
}
