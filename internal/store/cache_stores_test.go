package store_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

func TestCacheStores(t *testing.T) {
	testCases := []struct {
		name          string
		objectToStore client.Object
	}{
		{
			name: "KongUpstreamPolicy",
			objectToStore: &kongv1beta1.KongUpstreamPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "upstream-policy",
					Namespace: "default",
				},
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("least-connections"),
				},
			},
		},
		{
			name: "KongServiceFacade",
			objectToStore: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-facade",
					Namespace: "default",
				},
				Spec: incubatorv1alpha1.KongServiceFacadeSpec{
					Backend: incubatorv1alpha1.KongServiceFacadeBackend{
						Name: "backend",
						Port: 80,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := store.NewCacheStores()
			err := s.Add(tc.objectToStore)
			require.NoError(t, err)

			storedObj, ok, err := s.Get(tc.objectToStore)
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, tc.objectToStore, storedObj)

			err = s.Delete(tc.objectToStore)
			require.NoError(t, err)

			_, ok, err = s.Get(tc.objectToStore)
			require.NoError(t, err, err)
			require.False(t, ok)
		})
	}
}
