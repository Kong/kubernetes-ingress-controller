package mocks

import (
	"context"

	"github.com/kong/go-database-reconciler/pkg/file"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

// ConfigurationChangeDetector is a mock implementation of sendconfig.ConfigurationChangeDetector.
type ConfigurationChangeDetector struct {
	ConfigurationChanged bool
}

func (m ConfigurationChangeDetector) HasConfigurationChanged(
	context.Context, []byte, []byte, *file.Content, sendconfig.KonnectAwareClient, sendconfig.StatusClient,
) (bool, error) {
	return m.ConfigurationChanged, nil
}
