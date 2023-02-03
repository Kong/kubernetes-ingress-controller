package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCtrlFuncMetricsDoesNotPanicWhenCalledTwice(t *testing.T) {
	require.NotPanics(t, func() {
		NewCtrlFuncMetrics()
	})
	require.NotPanics(t, func() {
		NewCtrlFuncMetrics()
	})
}
