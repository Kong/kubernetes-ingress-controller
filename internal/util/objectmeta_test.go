package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFromK8sObject(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   client.Object
		want K8sObjectInfo
	}{
		{
			name: "empty annotations",
			in: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
			},
			want: K8sObjectInfo{
				Name:      "name",
				Namespace: "namespace",
			},
		},
		{
			name: "has annotations",
			in: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "name",
					Namespace:   "namespace",
					Annotations: map[string]string{"a": "1", "b": "2"},
				},
			},
			want: K8sObjectInfo{
				Name:        "name",
				Namespace:   "namespace",
				Annotations: map[string]string{"a": "1", "b": "2"},
			},
		},
		{
			name: "with group version kind",
			in: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "name",
					Namespace:   "namespace",
					Annotations: map[string]string{"a": "1", "b": "2"},
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
				},
			},
			want: K8sObjectInfo{
				Name:        "name",
				Namespace:   "namespace",
				Annotations: map[string]string{"a": "1", "b": "2"},
				GroupVersionKind: schema.GroupVersionKind{
					Group:   "networking.k8s.io",
					Version: "v1",
					Kind:    "Ingress",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := FromK8sObject(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkFromK8sObject(b *testing.B) {
	in := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "name",
			Namespace:   "namespace",
			Annotations: map[string]string{"a": "1", "b": "2"},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := FromK8sObject(in)
		_ = out
	}
}
