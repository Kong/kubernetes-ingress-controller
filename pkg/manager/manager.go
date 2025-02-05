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
	id      ID
	config  managercfg.Config
	logger  logr.Logger
	manager *managerinternal.Manager
}

// NewManager creates a new instance of the Kong Ingress Controller. It does not start the controller.
func NewManager(ctx context.Context, id ID, logger logr.Logger, cfg managercfg.Config) (*Manager, error) {
	m, err := managerinternal.New(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	return &Manager{
		id:      id,
		config:  cfg,
		logger:  logger.WithValues("managerID", id.String()),
		manager: m,
	}, nil
}

// Run starts the Kong Ingress Controller. It blocks until the context is cancelled.
// It should be called only once per Manager instance.
func (m *Manager) Run(ctx context.Context) error {
	return m.manager.Run(ctx)
}

// IsReady checks if the controller manager is ready to manage resources.
// It's only valid to call this method after the controller manager has been started
// with method Run(ctx).
func (m *Manager) IsReady() error {
	return m.manager.IsReady()
}

// ID returns the unique identifier of the manager.
func (m *Manager) ID() ID {
	return m.id
}

// Config returns the configuration of the manager.
func (m *Manager) Config() managercfg.Config {
	return m.config
}
