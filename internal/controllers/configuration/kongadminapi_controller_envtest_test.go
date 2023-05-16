//go:build envtest
// +build envtest

package configuration_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v2/test/envtest"
)

type notifier struct {
	lock sync.RWMutex
	t    *testing.T
	last []adminapi.DiscoveredAdminAPI
}

func (n *notifier) Notify(adminAPIs []adminapi.DiscoveredAdminAPI) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.last = adminAPIs
}

func (n *notifier) LastNotified() []adminapi.DiscoveredAdminAPI {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.last
}

// startKongAdminAPIServiceReconciler starts KongAdminAPIServiceReconciler with
// the manager in a separate goroutine.
func startKongAdminAPIServiceReconciler(ctx context.Context, t *testing.T, client ctrlclient.Client, cfg *rest.Config) (
	adminService corev1.Service,
	adminPod corev1.Pod,
	n *notifier,
) {
	ns := envtest.CreateNamespace(ctx, t, client)
	adminPod = envtest.CreatePod(ctx, t, client, ns.Name)

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Logger:             logrusr.New(logrus.New()),
		Scheme:             client.Scheme(),
		SyncPeriod:         lo.ToPtr(2 * time.Second),
		MetricsBindAddress: "0",
	})
	require.NoError(t, err)

	adminService = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kong-admin",
			Namespace: ns.Name,
			UID:       k8stypes.UID(uuid.NewString()),
		},
	}

	n = &notifier{t: t}
	require.NoError(t,
		(&configuration.KongAdminAPIServiceReconciler{
			Client: mgr.GetClient(),
			ServiceNN: k8stypes.NamespacedName{
				Name:      adminService.Name,
				Namespace: adminService.Namespace,
			},
			PortNames:         sets.New("admin"),
			EndpointsNotifier: n,
			Log:               mgr.GetLogger(),
		}).SetupWithManager(mgr),
	)
	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NoError(t, mgr.Start(ctx))
	}()
	t.Cleanup(func() { wg.Wait() })

	return
}

func TestKongAdminAPIController(t *testing.T) {
	t.Parallel()

	// In tests below we use a deferred cancel to stop the manager and not wait
	// for its timeout.

	cfg := envtest.Setup(t, scheme.Scheme)
	client, err := ctrlclient.New(cfg, ctrlclient.Options{})
	require.NoError(t, err)

	t.Run("Endpoints are matched properly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		adminService, adminPod, n := startKongAdminAPIServiceReconciler(ctx, t, client, cfg)

		endpoints := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
			},
			Ports: []discoveryv1.EndpointPort{
				{
					Name: lo.ToPtr("admin"),
					Port: lo.ToPtr(int32(8080)),
				},
				{
					Name: lo.ToPtr("admin-tls"),
					Port: lo.ToPtr(int32(8444)),
				},
				{
					Name: lo.ToPtr("kong-admin-tls"),
					Port: lo.ToPtr(int32(8445)),
				},
			},
		}
		require.NoError(t, client.Create(ctx, &endpoints, &ctrlclient.CreateOptions{}))

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 2 }, 3*time.Second, time.Millisecond)
		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.2:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)
	})

	t.Run("not Ready Endpoints are not matched", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		adminService, adminPod, n := startKongAdminAPIServiceReconciler(ctx, t, client, cfg)

		endpoints := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(false),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
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

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 1 }, 3*time.Second, time.Millisecond)
		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.2:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)
	})

	t.Run("multiple EndpointSlices are matched properly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		adminService, adminPod, n := startKongAdminAPIServiceReconciler(ctx, t, client, cfg)

		endpoints := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
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

		endpoints2 := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.10"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.20"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
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
		require.NoError(t, client.Create(ctx, &endpoints2, &ctrlclient.CreateOptions{}))

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 4 }, 3*time.Second, time.Millisecond)
		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.2:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.10:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.20:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)
	})

	t.Run("with EndpointSlices changing over time the notifications are sent properly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		adminService, adminPod, n := startKongAdminAPIServiceReconciler(ctx, t, client, cfg)

		endpoints := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
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

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 2 }, 3*time.Second, time.Millisecond)
		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.2:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)

		// Update all endpoints so that they are not Ready.
		for i := range endpoints.Endpoints {
			endpoints.Endpoints[i].Conditions.Ready = lo.ToPtr(false)
		}
		require.NoError(t, client.Update(ctx, &endpoints, &ctrlclient.UpdateOptions{}))
		require.NoError(t, client.Get(ctx, k8stypes.NamespacedName{Name: endpoints.Name, Namespace: endpoints.Namespace}, &endpoints, &ctrlclient.GetOptions{}))
		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 0 }, 3*time.Second, time.Millisecond)

		// Update 1 endpoint so that that it's Ready.
		endpoints.Endpoints[0].Conditions.Ready = lo.ToPtr(true)

		require.NoError(t, client.Update(ctx, &endpoints, &ctrlclient.UpdateOptions{}))
		require.NoError(t, client.Get(ctx, k8stypes.NamespacedName{Name: endpoints.Name, Namespace: endpoints.Namespace}, &endpoints, &ctrlclient.GetOptions{}))
		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 1 }, 3*time.Second, time.Millisecond)

		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)
	})

	t.Run("when deleted EndpointsSlice is observed notifications are sent properly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		adminService, adminPod, n := startKongAdminAPIServiceReconciler(ctx, t, client, cfg)

		endpoints := discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Service",
						Name:       adminService.Name,
						APIVersion: "v1",
						UID:        adminService.UID,
					},
				},
				Name:      uuid.NewString(),
				Namespace: adminService.Namespace,
				Labels: map[string]string{
					"kubernetes.io/service-name": adminService.Name,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Name:      adminPod.Name,
						Namespace: adminPod.Namespace,
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

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 2 }, 3*time.Second, time.Millisecond)
		assert.ElementsMatch(t,
			[]adminapi.DiscoveredAdminAPI{
				{
					Address: "https://10.0.0.1:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
				{
					Address: "https://10.0.0.2:8080",
					PodRef: k8stypes.NamespacedName{
						Namespace: adminPod.Namespace,
						Name:      adminPod.Name,
					},
				},
			},
			n.LastNotified(),
		)

		// Mark EndpointSlice deleted
		require.NoError(t, client.Delete(ctx, &endpoints, &ctrlclient.DeleteOptions{}))

		assert.Eventually(t, func() bool { return len(n.LastNotified()) == 0 }, 3*time.Second, time.Millisecond)
		assert.Nil(t, n.LastNotified())
	})
}
