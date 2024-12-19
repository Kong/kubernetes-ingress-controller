package mocks

import (
	"context"
	"sync"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
)

// ConfigStatusNotifier is a mock implementation of clients.ConfigStatusNotifier.
type ConfigStatusNotifier struct {
	lastGatewayConfigStatus     clients.GatewayConfigApplyStatus
	receivedKonnectConfigStatus []clients.KonnectConfigUploadStatus
	lock                        sync.RWMutex
}

func (c *ConfigStatusNotifier) NotifyGatewayConfigStatus(_ context.Context, status clients.GatewayConfigApplyStatus) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lastGatewayConfigStatus = status
}

func (c *ConfigStatusNotifier) NotifyKonnectConfigStatus(_ context.Context, status clients.KonnectConfigUploadStatus) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.receivedKonnectConfigStatus = append(c.receivedKonnectConfigStatus, status)
}

func (c *ConfigStatusNotifier) LastGatewayConfigStatus() clients.GatewayConfigApplyStatus {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lastGatewayConfigStatus
}

func (c *ConfigStatusNotifier) FirstKonnectConfigStatus() (clients.KonnectConfigUploadStatus, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if len(c.receivedKonnectConfigStatus) == 0 {
		return clients.KonnectConfigUploadStatus{}, false
	}
	return c.receivedKonnectConfigStatus[0], true
}
