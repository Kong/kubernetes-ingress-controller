package translator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
)

func TestObjectsCollector(t *testing.T) {
	t.Run("adds and pops objects", func(t *testing.T) {
		objects := []client.Object{
			&corev1.Pod{},
			&corev1.Service{},
		}

		collector := translator.NewObjectsCollector()
		for _, obj := range objects {
			collector.Add(obj)
		}

		got := collector.Pop()
		require.ElementsMatch(t, objects, got)

		secondGot := collector.Pop()
		require.Nil(t, secondGot)
	})

	t.Run("gracefully handles nil receiver", func(t *testing.T) {
		var collector *translator.ObjectsCollector

		require.NotPanics(t, func() {
			collector.Add(&corev1.Pod{})
			got := collector.Pop()
			require.Nil(t, got)
		})
	})
}
