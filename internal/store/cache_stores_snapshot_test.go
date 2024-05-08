package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestCacheStores_TakeSnapshot(t *testing.T) {
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
			},
		},
	}...)
	require.NoError(t, err)

	t.Log("Taking a snapshot of the originalStores")
	snapshot, err := originalStores.TakeSnapshot()
	require.NoError(t, err)

	t.Log("Adding a new object to the original store")
	err = originalStores.Add(&netv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bar",
			Namespace: "default",
			Annotations: map[string]string{
				"foo": "baz",
			},
		},
	})
	require.NoError(t, err)

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
