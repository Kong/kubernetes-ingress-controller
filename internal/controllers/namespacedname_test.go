package controllers

import (
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestOptionalNamespacedName(t *testing.T) {
	t.Run("Get empty", func(t *testing.T) {
		emptyOptional := OptionalNamespacedName{}
		nn, ok := emptyOptional.Get()
		require.False(t, ok)
		assert.Equal(t, k8stypes.NamespacedName{}, nn)
	})

	t.Run("Get nonempty", func(t *testing.T) {
		namespacedName := k8stypes.NamespacedName{
			Name:      "example",
			Namespace: "default",
		}
		presentOptional := NewOptionalNamespacedName(mo.Some(namespacedName))
		nn, ok := presentOptional.Get()
		require.True(t, ok)
		assert.Equal(t, namespacedName, nn)
	})

	t.Run("Matches", func(t *testing.T) {
		namespacedName := k8stypes.NamespacedName{
			Name:      "example",
			Namespace: "default",
		}
		presentOptional := NewOptionalNamespacedName(mo.Some(namespacedName))
		assert.True(t, presentOptional.Matches(&gatewayapi.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
			},
		}))
		assert.False(t, presentOptional.Matches(&gatewayapi.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy",
				Namespace: "default",
			},
		}))
		assert.False(t, presentOptional.Matches(&gatewayapi.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "dummy",
			},
		}))
	})
}
