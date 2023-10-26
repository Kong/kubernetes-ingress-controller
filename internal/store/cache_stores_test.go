package store_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestCacheStores(t *testing.T) {
	t.Run("KongUpstreamPolicy", func(t *testing.T) {
		s := store.NewCacheStores()
		upstreamPolicy := &kongv1beta1.KongUpstreamPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upstream-policy",
				Namespace: "default",
			},
			Spec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("least-connections"),
			},
		}

		err := s.Add(upstreamPolicy)
		require.NoError(t, err)

		objKey := &kongv1beta1.KongUpstreamPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upstream-policy",
				Namespace: "default",
			},
		}
		storedObj, ok, err := s.Get(objKey)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, upstreamPolicy, storedObj)

		err = s.Delete(objKey)
		require.NoError(t, err)

		_, ok, err = s.Get(objKey)
		require.NoError(t, err, err)
		require.False(t, ok)
	})
}
