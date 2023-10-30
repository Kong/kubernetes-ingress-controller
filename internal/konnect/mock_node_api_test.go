package konnect_test

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
)

// mockNodeClient is a mock implementation of the NodeClient interface.
type mockNodeClient struct {
	nodes                    map[string]*nodes.NodeItem
	returnErrorFromListNodes bool
	wasListAllNodesCalled    bool
	nodeUpdatesCount         atomic.Int32
	lock                     sync.RWMutex
}

func newMockNodeClient(initialNodes []*nodes.NodeItem) *mockNodeClient {
	nodesMap := lo.SliceToMap(initialNodes, func(i *nodes.NodeItem) (string, *nodes.NodeItem) {
		return i.ID, i
	})
	return &mockNodeClient{nodes: nodesMap}
}

func (m *mockNodeClient) CreateNode(_ context.Context, req *nodes.CreateNodeRequest) (*nodes.CreateNodeResponse, error) {
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

func (m *mockNodeClient) UpdateNode(_ context.Context, nodeID string, req *nodes.UpdateNodeRequest) (*nodes.UpdateNodeResponse, error) {
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

func (m *mockNodeClient) DeleteNode(_ context.Context, nodeID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.nodes, nodeID)
	return nil
}

func (m *mockNodeClient) ListAllNodes(_ context.Context) ([]*nodes.NodeItem, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.wasListAllNodesCalled = true
	return lo.MapToSlice(m.nodes, func(_ string, i *nodes.NodeItem) *nodes.NodeItem {
		return i
	}), nil
}

func (m *mockNodeClient) upsertNode(node *nodes.NodeItem) *nodes.NodeItem {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.nodeUpdatesCount.Add(1)

	if node.ID == "" {
		node.ID = uuid.New().String()
	}
	m.nodes[node.ID] = node
	return node
}

func (m *mockNodeClient) MustAllNodes() []*nodes.NodeItem {
	ns, err := m.ListAllNodes(context.Background())
	if err != nil {
		panic(err)
	}
	return ns
}

func (m *mockNodeClient) WasListAllNodesCalled() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.wasListAllNodesCalled
}

func (m *mockNodeClient) ReturnErrorFromListAllNodes(v bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.returnErrorFromListNodes = v
}

func (m *mockNodeClient) NodesUpdatesCount() int {
	return int(m.nodeUpdatesCount.Load())
}
