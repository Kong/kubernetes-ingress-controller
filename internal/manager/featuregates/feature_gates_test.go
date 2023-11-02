package featuregates

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFeatureGates(t *testing.T) {
	t.Log("Setting up configurations and logging for feature gates testing")
	setupLog := zapr.NewLogger(zap.NewNop())

	t.Log("Verifying feature gates setup defaults when no feature gates are configured")
	fgs, err := New(setupLog, nil)
	assert.NoError(t, err)
	assert.Len(t, fgs, len(GetFeatureGatesDefaults()))

	t.Log("Verifying feature gates setup results when valid feature gates options are present")
	featureGates := map[string]bool{GatewayAlphaFeature: true}
	fgs, err = New(setupLog, featureGates)
	assert.NoError(t, err)
	assert.True(t, fgs[GatewayAlphaFeature])

	t.Log("Configuring several invalid feature gates options")
	featureGates = map[string]bool{"invalidGateway": true}

	t.Log("Verifying feature gates setup results when invalid feature gates options are present")
	_, err = New(setupLog, featureGates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalidGateway is not a valid feature")
}
