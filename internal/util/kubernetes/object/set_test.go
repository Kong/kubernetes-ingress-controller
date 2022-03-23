package object

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestSet(t *testing.T) {
	t.Log("generating some objects to test the object set")
	ing1 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "test-ingress-1",
		},
	}
	ing1.SetGroupVersionKind(ingGVK)
	ing2 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "test-ingress-2",
		},
	}
	ing2.SetGroupVersionKind(ingGVK)
	ing3 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "other-namespace",
			Name:      "test-ingress-3",
		},
	}
	ing3.SetGroupVersionKind(ingGVK)
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "test-tcpingress",
		},
	}
	tcp.SetGroupVersionKind(tcpGVK)

	t.Log("verifying creation of an object set")
	set := &Set{}
	require.False(t, set.Has(ing1))
	require.False(t, set.Has(ing2))
	require.False(t, set.Has(ing3))
	require.False(t, set.Has(tcp))

	t.Log("verifying object set insertion")
	set.Insert(ing1)
	require.True(t, set.Has(ing1))
	require.False(t, set.Has(ing2))
	require.False(t, set.Has(ing3))
	require.False(t, set.Has(tcp))
	set.Insert(ing2)
	require.True(t, set.Has(ing1))
	require.True(t, set.Has(ing2))
	require.False(t, set.Has(ing3))
	require.False(t, set.Has(tcp))
	set.Insert(ing3)
	require.True(t, set.Has(ing1))
	require.True(t, set.Has(ing2))
	require.True(t, set.Has(ing3))
	require.False(t, set.Has(tcp))
	set.Insert(tcp)
	require.True(t, set.Has(ing1))
	require.True(t, set.Has(ing2))
	require.True(t, set.Has(ing3))
	require.True(t, set.Has(tcp))
}

// -----------------------------------------------------------------------------
// Testing Utilities
// -----------------------------------------------------------------------------

// initialized objects don't have GVK's, so we fake those for unit tests.
var (
	ingGVK = schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	}
	tcpGVK = schema.GroupVersionKind{
		Group:   "configuration.konghq.com",
		Version: "v1beta1",
		Kind:    "TCPIngress",
	}
)
