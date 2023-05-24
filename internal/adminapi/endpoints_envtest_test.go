//go:build envtest
// +build envtest

package adminapi_test

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
	"k8s.io/client-go/kubernetes/scheme"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	cfgtypes "github.com/kong/kubernetes-ingress-controller/v2/internal/manager/config/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/envtest"
)

func TestGetAdminAPIsForServiceReturnsAllAddressesCorrectlyPagingThroughResults(t *testing.T) {
	t.Parallel()

	var client ctrlclient.Client
	{
		cfg := envtest.Setup(t, scheme.Scheme)
		var err error
		client, err = ctrlclient.New(cfg, ctrlclient.Options{
			Scheme: scheme.Scheme,
		})
		require.NoError(t, err)
	}

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
		tc := tc
		t.Run(fmt.Sprintf("%dx%d", tc.subnetC, tc.subnetD), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var (
				ns          = envtest.CreateNamespace(ctx, t, client)
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

			got, err := adminapi.GetAdminAPIsForService(ctx, client, service, sets.New("admin"), cfgtypes.IPDNSStrategy)
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
