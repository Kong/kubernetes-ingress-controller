package store_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestCacheStores_TakeSnapshot(t *testing.T) {
	originalStores := getStoresForTests(t)
	t.Log("Taking a snapshot of the originalStores")
	snapshot, err := originalStores.TakeSnapshot()
	require.NoError(t, err)
	require.NotEqual(t, store.CacheStores{}, snapshot)

	testIfSnapshotIsIndependentFromSource(t, &originalStores, &snapshot)
}

func TestCacheStores_TakeSnapshotIfChanged(t *testing.T) {
	originalStores := getStoresForTests(t)
	t.Log("Taking a snapshot of the originalStores")
	originalStoresSnapshot, originalStoresHash, err := originalStores.TakeSnapshotIfChanged(store.SnapshotHashEmpty)
	require.NoError(t, err)
	require.Equal(t, store.SnapshotHash("4FU3CVGPMPGXTSC2I3L4UCKE6M46LWOJBSQMJ2N7HAYR46IVHNGQ===="), originalStoresHash)
	require.NotEqual(t, store.CacheStores{}, originalStoresSnapshot)

	t.Log("Taking again a snapshot of the originalStores")
	originalStoresSnapshot2, originalStoresHash2, err := originalStores.TakeSnapshotIfChanged(originalStoresHash)
	require.NoError(t, err)
	require.Empty(t, originalStoresHash2)
	require.Equal(t, store.CacheStores{}, originalStoresSnapshot2)

	testIfSnapshotIsIndependentFromSource(t, &originalStores, &originalStoresSnapshot) // This modifies the originalStores.

	t.Log("Taking a snapshot of the originalStores after modifying it")
	modStoresSnapshot, modStoresHash, err := originalStores.TakeSnapshotIfChanged(originalStoresHash)
	require.NoError(t, err)
	require.NotEmpty(t, modStoresHash)
	require.NotEqual(t, originalStoresHash, modStoresSnapshot)
	require.NotEqual(t, store.CacheStores{}, modStoresSnapshot)

	originalStores = getStoresForTests(t)
	t.Log("Taking a snapshot of the originalStores after resetting it")
}

func getStoresForTests(t *testing.T) store.CacheStores {
	t.Helper()
	originalStores, err := store.NewCacheStoresFromObjs([]runtime.Object{
		&netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					"foo": "bar",
				},
				ResourceVersion: "123",
				UID:             "19c1f570-b301-4b09-adcb-a1d29eb0b27e",
			},
		},
		&netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					"foo": "bar",
				},
				ResourceVersion: "456",
				UID:             "3a463a14-59ba-422d-9d7b-02f57ddb2800",
			},
		},
		&netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					"foo": "bar",
				},
				ResourceVersion: "789",
				UID:             "2f37d2d3-0e95-40e7-810a-31953df5ee69",
			},
		},
	}...)
	require.NoError(t, err)
	return originalStores
}

func testIfSnapshotIsIndependentFromSource(t *testing.T, originalStores, snapshot *store.CacheStores) {
	t.Helper()

	// We'll use ingressMeta in .Get() calls.
	ingressMeta := &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
	}

	t.Log("Modifying the original object")
	obj, ok, err := originalStores.IngressV1.Get(ingressMeta)
	require.NoError(t, err)
	require.True(t, ok)
	ingressFromOriginalStore, ok := obj.(*netv1.Ingress)
	require.True(t, ok)
	require.Equal(t, "bar", ingressFromOriginalStore.Annotations["foo"])
	ingressFromOriginalStore.Annotations["foo"] = "qux"
	ingressFromOriginalStore.ResourceVersion = "567"

	t.Log("Checking that the original store returns the modified object")
	obj, ok, err = originalStores.IngressV1.Get(ingressMeta)
	require.NoError(t, err)
	require.True(t, ok)
	ingressFromOriginalStore2, ok := obj.(*netv1.Ingress)
	require.True(t, ok)
	require.Same(t, ingressFromOriginalStore, ingressFromOriginalStore2, "Store should return pointers the same object")

	t.Log("Checking that the snapshot returns the unmodified object")
	obj, ok, err = snapshot.IngressV1.Get(obj)
	require.NoError(t, err)
	require.True(t, ok)
	ingressFromSnapshot, ok := obj.(*netv1.Ingress)
	require.True(t, ok)
	require.Equal(t, "bar", ingressFromSnapshot.Annotations["foo"], "Snapshot should not be affected by the change in the original store")
}

func BenchmarkCacheStores_TakeSnapshot(b *testing.B) {
	smallStores, err := store.NewCacheStoresFromObjs([]runtime.Object{
		&netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					"foo": "bar",
				},
				ResourceVersion: "123",
				UID:             k8stypes.UID(uuid.New().String()),
			},
		},
	}...)
	require.NoError(b, err)

	var k8sObjects []runtime.Object
	for i := 0; i < 1_000; i++ {
		route := &gatewayapi.HTTPRoute{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "gateway.networking.k8s.io/v1",
				Kind:       "HTTPRoute",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            fmt.Sprintf("route-%d", i),
				Namespace:       "default",
				ResourceVersion: fmt.Sprintf("12%d", i),
				UID:             k8stypes.UID(uuid.New().String()),
			},
			Spec: gatewayapi.HTTPRouteSpec{
				Rules: []gatewayapi.HTTPRouteRule{
					{
						Matches: []gatewayapi.HTTPRouteMatch{
							{
								Path: &gatewayapi.HTTPPathMatch{
									Type:  lo.ToPtr(gatewayapi.PathMatchExact),
									Value: lo.ToPtr("/test1"),
								},
								Method: lo.ToPtr(gatewayapi.HTTPMethodGet),
							},
							{
								Path: &gatewayapi.HTTPPathMatch{
									Type:  lo.ToPtr(gatewayapi.PathMatchExact),
									Value: lo.ToPtr("/test2"),
								},
								Method: lo.ToPtr(gatewayapi.HTTPMethodGet),
							},
						},
					},
				},
				CommonRouteSpec: gatewayapi.CommonRouteSpec{},
			},
		}
		k8sObjects = append(k8sObjects, route)

		service := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            fmt.Sprintf("service-%d", i),
				Namespace:       "default",
				ResourceVersion: fmt.Sprintf("12%d", i),
				UID:             k8stypes.UID(uuid.New().String()),
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: "http",
						Port: int32(i),
					},
				},
			},
		}
		k8sObjects = append(k8sObjects, service)

		ingress := &netv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            fmt.Sprintf("ingress-%d", i),
				Namespace:       "default",
				ResourceVersion: fmt.Sprintf("12%d", i),
				UID:             k8stypes.UID(uuid.New().String()),
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: fmt.Sprintf("host-%d", i),
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     fmt.Sprintf("/path-%d", i),
										PathType: lo.ToPtr(netv1.PathTypeExact),
									},
								},
							},
						},
					},
				},
			},
		}
		k8sObjects = append(k8sObjects, ingress)
	}
	bigStores, err := store.NewCacheStoresFromObjs(k8sObjects...)
	require.NoError(b, err)

	b.Run("Small_Without_Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s, err := smallStores.TakeSnapshot()
			if err != nil || (s == store.CacheStores{}) {
				b.Fatalf("unexpected error or empty snapshot err: %v, snapshot: %v", err, s)
			}
		}
	})
	b.Run("Big_Without_Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s, err := bigStores.TakeSnapshot()
			if err != nil || (s == store.CacheStores{}) {
				b.Fatalf("unexpected error or empty snapshot err: %v, snapshot: %v", err, s)
			}
		}
	})
	b.Run("Small_With_Cache", func(b *testing.B) {
		s, hash, err := smallStores.TakeSnapshotIfChanged(store.SnapshotHashEmpty)
		require.NoError(b, err)
		require.NotEmpty(b, hash)
		require.NotEqual(b, store.CacheStores{}, s)
		for i := 0; i < b.N; i++ {
			s, hashNew, err := smallStores.TakeSnapshotIfChanged(hash)
			if err != nil || hashNew != store.SnapshotHashEmpty || s != (store.CacheStores{}) {
				b.Fatalf("unexpected error or non-empty snapshot err: %v, hash: %s, snapshot: %v", err, hashNew, s)
			}
		}
	})

	b.Run("Big_With_Cache", func(b *testing.B) {
		s, hash, err := bigStores.TakeSnapshotIfChanged(store.SnapshotHashEmpty)
		require.NoError(b, err)
		require.NotEmpty(b, hash)
		require.NotEqual(b, store.CacheStores{}, s)
		for i := 0; i < b.N; i++ {
			s, hashNew, err := bigStores.TakeSnapshotIfChanged(hash)
			if err != nil || hashNew != store.SnapshotHashEmpty || s != (store.CacheStores{}) {
				b.Fatalf("unexpected error or non-empty snapshot err: %v, hash: %s, snapshot: %v", err, hashNew, s)
			}
		}
	})

	b.Run("Small_With_Missed_Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s, hashNew, err := smallStores.TakeSnapshotIfChanged(store.SnapshotHashEmpty)
			if err != nil || hashNew == store.SnapshotHashEmpty || s == (store.CacheStores{}) {
				b.Fatalf("unexpected error or empty snapshot err: %v, hash: %s, snapshot: %v", err, hashNew, s)
			}
		}
	})

	b.Run("Big_With_Missed_Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s, hashNew, err := bigStores.TakeSnapshotIfChanged(store.SnapshotHashEmpty)
			if err != nil || hashNew == store.SnapshotHashEmpty || s == (store.CacheStores{}) {
				b.Fatalf("unexpected error or empty snapshot err: %v, hash: %s, snapshot: %v", err, hashNew, s)
			}
		}
	})
}
