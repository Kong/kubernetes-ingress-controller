package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// WithTypeMeta adds type meta to the given object based on its Go type.
func WithTypeMeta[T runtime.Object](t *testing.T, obj T) T {
	s, err := scheme.Get()
	require.NoError(t, err)
	err = util.PopulateTypeMeta(obj, s)
	require.NoError(t, err)
	return obj
}
