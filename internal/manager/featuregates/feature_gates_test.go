package featuregates

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFeatureGates(t *testing.T) {
	t.Log("setting up configurations and logging for feature gates testing")
	// TODO 1893 this was setting up a buffer, but then the tests don't actually check it, so a Nop is apparently fine
	setupLog := zapr.NewLogger(zap.NewNop())

	t.Log("verifying feature gates setup defaults when no feature gates are configured")
	fgs, err := New(setupLog, nil)
	assert.NoError(t, err)
	assert.Len(t, fgs, len(GetFeatureGatesDefaults()))

	t.Log("verifying feature gates setup results when valid feature gates options are present")
	featureGates := map[string]bool{GatewayFeature: true}
	fgs, err = New(setupLog, featureGates)
	assert.NoError(t, err)
	assert.True(t, fgs[GatewayFeature])

	t.Log("configuring several invalid feature gates options")
	featureGates = map[string]bool{"invalidGateway": true}

	t.Log("verifying feature gates setup results when invalid feature gates options are present")
	_, err = New(setupLog, featureGates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalidGateway is not a valid feature")
}
