package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
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
			in: &networkingv1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
			},
			want: K8sObjectInfo{
				Name:        "name",
				Namespace:   "namespace",
				Annotations: map[string]string{},
			},
		},
		{
			name: "has annotations",
			in: &networkingv1beta1.Ingress{
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
			assert.EqualValues(t, tt.want, got)
		})
	}
}
