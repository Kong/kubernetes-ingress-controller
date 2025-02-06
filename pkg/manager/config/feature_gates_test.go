package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureGates(t *testing.T) {
	t.Log("Verifying feature gates setup defaults when no feature gates are configured")
	fgs, err := NewFeatureGates(nil)
	assert.NoError(t, err)
	assert.Len(t, fgs, len(GetFeatureGatesDefaults()))

	t.Log("Verifying feature gates setup results when valid feature gates options are present")
	featureGates := map[string]bool{GatewayAlphaFeature: true}
	fgs, err = NewFeatureGates(featureGates)
	assert.NoError(t, err)
	assert.True(t, fgs[GatewayAlphaFeature])

	t.Log("Verifying feature gates setup will return error when settings has conflicts")
	featureGates = map[string]bool{KongCustomEntityFeature: true, FillIDsFeature: false}
	_, err = NewFeatureGates(featureGates)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("%s is required if %s is enabled", FillIDsFeature, KongCustomEntityFeature))

	t.Log("Configuring several invalid feature gates options")
	featureGates = map[string]bool{"invalidGateway": true}

	t.Log("Verifying feature gates setup results when invalid feature gates options are present")
	_, err = NewFeatureGates(featureGates)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalidGateway is not a valid feature")
}
