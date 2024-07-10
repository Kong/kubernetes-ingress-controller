package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestTypeMetaFromK8sObject(t *testing.T) {
	testCases := []struct {
		name     string
		obj      client.Object
		typeMeta metav1.TypeMeta
	}{
		{
			name: "empty group",
			obj: &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "svc",
					Namespace: "default",
				},
			},
			typeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
		},
		{
			name: "non-empty group",
			obj: &netv1.Ingress{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "networking.k8s.io/v1",
					Kind:       "Ingress",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ing",
					Namespace: "default",
				},
			},
			typeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			typeMeta := TypeMetaFromK8sObject(tc.obj)
			require.Equal(t, tc.typeMeta, typeMeta)
		})
	}
}
