package manager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	managerinternal "github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// Manager is an object representing an instance of the Kong Ingress Controller.
type Manager struct {
	id      ID
	logger  logr.Logger
	manager *managerinternal.Manager
}

// NewManager creates a new instance of the Kong Ingress Controller. It does not start the controller.
func NewManager(ctx context.Context, id ID, logger logr.Logger, cfg managercfg.Config) (*Manager, error) {
	m, err := managerinternal.New(ctx, id, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	return &Manager{
		id: id,

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

// GetKubeconfig returns the Kubernetes REST config object associated with the instance.
func (m *Manager) GetKubeconfig() *rest.Config {
	return m.manager.GetKubeconfig()
}

// GetClientsManager returns the clients manager associated with the instance.
func (m *Manager) GetClientsManager() *clients.AdminAPIClientsManager {
	return m.manager.GetClientsManager()
}

// DiagnosticsHandler returns the diagnostics handler of the manager if available. Otherwise, it returns nil.
func (m *Manager) DiagnosticsHandler() http.Handler {
	return m.manager.DiagnosticsHandler()
}
