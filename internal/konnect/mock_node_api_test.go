package konnect_test

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/nodes"
)

// mockNodeAPIClient is a mock implementation of the NodeAPIClient interface.
type mockNodeAPIClient struct {
	nodes                    map[string]*nodes.NodeItem
	returnErrorFromListNodes bool
	wasListAllNodesCalled    bool
	lock                     sync.RWMutex
}

func newMockNodeAPIClient(initialNodes []*nodes.NodeItem) *mockNodeAPIClient {
	nodesMap := lo.SliceToMap(initialNodes, func(i *nodes.NodeItem) (string, *nodes.NodeItem) {
		return i.ID, i
	})
	return &mockNodeAPIClient{nodes: nodesMap}
}

func (m *mockNodeAPIClient) CreateNode(_ context.Context, req *nodes.CreateNodeRequest) (*nodes.CreateNodeResponse, error) {
	node := m.upsertNode(&nodes.NodeItem{
		ID:       req.ID,
		Version:  req.Version,
		Hostname: req.Hostname,
		LastPing: req.LastPing,
		Type:     req.Type,
		Status:   req.Status,
	})
	return &nodes.CreateNodeResponse{Item: node}, nil
}

func (m *mockNodeAPIClient) UpdateNode(_ context.Context, nodeID string, req *nodes.UpdateNodeRequest) (*nodes.UpdateNodeResponse, error) {
	node := m.upsertNode(&nodes.NodeItem{
		ID:       nodeID,
		Version:  req.Version,
		Hostname: req.Hostname,
		LastPing: req.LastPing,
		Type:     req.Type,
		Status:   req.Status,
	})
	return &nodes.UpdateNodeResponse{Item: node}, nil
}

func (m *mockNodeAPIClient) DeleteNode(_ context.Context, nodeID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.nodes, nodeID)
	return nil
}

func (m *mockNodeAPIClient) ListAllNodes(_ context.Context) ([]*nodes.NodeItem, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.wasListAllNodesCalled = true
	return lo.MapToSlice(m.nodes, func(_ string, i *nodes.NodeItem) *nodes.NodeItem {
		return i
	}), nil
}

func (m *mockNodeAPIClient) upsertNode(node *nodes.NodeItem) *nodes.NodeItem {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node.ID == "" {
		node.ID = uuid.New().String()
	}
	m.nodes[node.ID] = node
	return node
}

func (m *mockNodeAPIClient) MustAllNodes() []*nodes.NodeItem {
	ns, err := m.ListAllNodes(context.Background())
	if err != nil {
		panic(err)
	}
	return ns
}

func (m *mockNodeAPIClient) WasListAllNodesCalled() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.wasListAllNodesCalled
}

func (m *mockNodeAPIClient) ReturnErrorFromListAllNodes(v bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.returnErrorFromListNodes = v
}
