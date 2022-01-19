package manager

import (
	"bytes"
	"testing"

	"github.com/bombsimon/logrusr/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFeatureGates(t *testing.T) {
	t.Log("setting up configurations and logging for feature gates testing")
	out := new(bytes.Buffer)
	baseLogger := logrus.New()
	baseLogger.SetOutput(out)
	baseLogger.SetLevel(logrus.DebugLevel)
	setupLog := logrusr.New(baseLogger)
	config := new(Config)

	t.Log("verifying feature gates setup defaults when no feature gates are configured")
	fgs, err := setupFeatureGates(setupLog, config)
	assert.NoError(t, err)
	assert.Len(t, fgs, len(getFeatureGatesDefaults()))

	t.Log("verifying feature gates setup results when valid feature gates options are present")
	config.FeatureGates = map[string]bool{knativeFeature: true}
	fgs, err = setupFeatureGates(setupLog, config)
	assert.NoError(t, err)
	assert.True(t, fgs[knativeFeature])

	t.Log("configuring several invalid feature gates options")
	config.FeatureGates = map[string]bool{"invalidGateway": true}

	t.Log("verifying feature gates setup results when invalid feature gates options are present")
	_, err = setupFeatureGates(setupLog, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalidGateway is not a valid feature")
}
