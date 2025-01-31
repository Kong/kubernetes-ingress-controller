package manager

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	managerinternal "github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// Manager is an object representing an instance of the Kong Ingress Controller.
type Manager struct {
	id     ID
	config managercfg.Config
	logger logr.Logger
}

// NewManager creates a new instance of the Kong Ingress Controller. It does not start the controller.
func NewManager(id ID, logger logr.Logger, configOpts ...managercfg.Opt) (*Manager, error) {
	cfg, err := NewConfig(configOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create manager config: %w", err)
	}

	return &Manager{
		id:     id,
		config: cfg,
		logger: logger.WithValues("managerID", id.String()),
	}, nil
}

// Run starts the Kong Ingress Controller. It blocks until the context is cancelled.
// It should be called only once per Manager instance.
func (m *Manager) Run(ctx context.Context) error {
	return managerinternal.Run(ctx, m.config, m.logger)
}

// TODO(czeslavo): expose healthcheck/readiness check methods from the manager
