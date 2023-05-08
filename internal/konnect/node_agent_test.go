package konnect_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/nodes"
)

const (
	testKicVersion  = "2.9.0"
	testKongVersion = "3.2.0.0"
	testHostname    = "ingress-0"
)

type mockGatewayInstanceGetter struct {
	gatewayInstances []konnect.GatewayInstance
}

func newMockGatewayInstanceGetter(instances []konnect.GatewayInstance) *mockGatewayInstanceGetter {
	return &mockGatewayInstanceGetter{gatewayInstances: instances}
}

func (m *mockGatewayInstanceGetter) GetGatewayInstances(context.Context) ([]konnect.GatewayInstance, error) {
	return m.gatewayInstances, nil
}

type mockGatewayClientsNotifier struct {
	ch chan struct{}
}

func newMockGatewayClientsNotifier() *mockGatewayClientsNotifier {
	return &mockGatewayClientsNotifier{
		ch: make(chan struct{}),
	}
}

func (m *mockGatewayClientsNotifier) SubscribeToGatewayClientsChanges() (<-chan struct{}, bool) {
	return m.ch, true
}

func (m *mockGatewayClientsNotifier) Notify() {
	m.ch <- struct{}{}
}

type mockManagerInstanceIDProvider struct {
	instanceID uuid.UUID
}

func newMockManagerInstanceIDProvider(instanceID uuid.UUID) *mockManagerInstanceIDProvider {
	return &mockManagerInstanceIDProvider{instanceID: instanceID}
}

func (m *mockManagerInstanceIDProvider) GetID() uuid.UUID {
	return m.instanceID
}

type mockConfigStatusQueue struct {
	ch chan clients.ConfigStatus
}

func newMockConfigStatusNotifier() *mockConfigStatusQueue {
	return &mockConfigStatusQueue{
		ch: make(chan clients.ConfigStatus),
	}
}

func (m mockConfigStatusQueue) SubscribeConfigStatus() chan clients.ConfigStatus {
	return m.ch
}

func (m mockConfigStatusQueue) NotifyConfigStatus(_ context.Context, status clients.ConfigStatus) {
	m.ch <- status
}

func TestNodeAgentUpdateNodes(t *testing.T) {
	const (
		timeout = 10 * time.Second
		tick    = 10 * time.Millisecond
	)

	testManagerID := uuid.New()
	testNodeIDs := lo.Map(lo.Range(3), func(_, _ int) string { return uuid.NewString() })

	testCases := []struct {
		name                  string
		initialNodesInNodeAPI []*nodes.NodeItem
		// When configStatus is non-nil, notify the status to node agent in the test case.
		configStatus     *clients.ConfigStatus
		gatewayInstances []konnect.GatewayInstance

		containNodes    []*nodes.NodeItem
		notContainNodes []*nodes.NodeItem
		numNodes        int
	}{
		{
			name: "create kic node",
			// no existing nodes
			initialNodesInNodeAPI: nil,
			configStatus:          lo.ToPtr(clients.ConfigStatusOK),
			containNodes: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testManagerID.String(),
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
			},
			numNodes: 1,
		},
		{
			name: "update status existing kic node",
			initialNodesInNodeAPI: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testNodeIDs[0],
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
			},
			configStatus: lo.ToPtr(clients.ConfigStatusTranslationErrorHappened),
			containNodes: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testNodeIDs[0],
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStatePartialConfigFail),
					Version:  testKicVersion,
				},
			},
			numNodes: 1,
		},
		{
			name: "remove outdated KIC nodes",
			initialNodesInNodeAPI: []*nodes.NodeItem{
				// older node with same hostname, should delete this.
				{
					Hostname: testHostname,
					ID:       testNodeIDs[0],
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStatePartialConfigFail),
					Version:  testKicVersion,
					LastPing: time.Now().Unix() - 10,
				},
				// newer node, should reserve this.
				{
					Hostname: testHostname,
					ID:       testManagerID.String(),
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
					LastPing: time.Now().Unix() - 3,
				},
				// KIC node with other name, should delete this.
				{
					Hostname: "ingress-1",
					ID:       testNodeIDs[2],
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
					LastPing: time.Now().Unix() - 3,
				},
			},
			containNodes: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testManagerID.String(),
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
			},
			notContainNodes: []*nodes.NodeItem{
				{
					Hostname: "ingress-1",
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
			},
			numNodes: 1,
		},
		{
			name: "update gateway nodes and remove outdated nodes",
			initialNodesInNodeAPI: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testManagerID.String(),
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
				{
					Hostname: testHostname,
					ID:       testNodeIDs[0],
					Type:     nodes.NodeTypeKongProxy,
					Version:  testKongVersion,
				},
				// 2 gateway nodes with same name, should reserve newer one.
				{
					Hostname: "proxy-1",
					ID:       testNodeIDs[1],
					Type:     nodes.NodeTypeKongProxy,
					Version:  testKongVersion,
					LastPing: time.Now().Unix() - 10,
				},
				{
					Hostname: "proxy-1",
					ID:       testNodeIDs[2],
					Type:     nodes.NodeTypeKongProxy,
					Version:  testKongVersion,
					LastPing: time.Now().Unix() - 5,
				},
			},
			gatewayInstances: []konnect.GatewayInstance{
				{Hostname: "proxy-1", Version: testKongVersion, NodeID: testNodeIDs[2]},
			},
			containNodes: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					ID:       testManagerID.String(),
					Type:     nodes.NodeTypeIngressController,
					Status:   string(nodes.IngressControllerStateOperational),
					Version:  testKicVersion,
				},
				{
					Hostname: "proxy-1",
					ID:       testNodeIDs[2],
					Type:     nodes.NodeTypeKongProxy,
					Version:  testKongVersion,
				},
			},
			notContainNodes: []*nodes.NodeItem{
				{
					Hostname: testHostname,
					Type:     nodes.NodeTypeKongProxy,
					Version:  testKongVersion,
				},
			},
			numNodes: 2,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			nodeClient := newMockNodeAPIClient(tc.initialNodesInNodeAPI)
			configStatusQueue := newMockConfigStatusNotifier()
			gatewayClientsChangesNotifier := newMockGatewayClientsNotifier()

			nodeAgent := konnect.NewNodeAgent(
				testHostname,
				testKicVersion,
				konnect.DefaultRefreshNodePeriod,
				logr.Discard(),
				nodeClient,
				configStatusQueue,
				newMockGatewayInstanceGetter(tc.gatewayInstances),
				gatewayClientsChangesNotifier,
				newMockManagerInstanceIDProvider(testManagerID),
			)

			ctx, cancel := context.WithCancel(context.Background())
			agentReturned := make(chan struct{})
			go func() {
				require.NoError(t, nodeAgent.Start(ctx))
				close(agentReturned)
			}()

			if tc.configStatus != nil {
				configStatusQueue.NotifyConfigStatus(ctx, *tc.configStatus)
			}

			require.Eventually(t, func() bool {
				// Notify gateway clients changes.
				gatewayClientsChangesNotifier.Notify()

				// Check number of nodes in RG.
				ns := nodeClient.MustAllNodes()
				if len(ns) != tc.numNodes {
					t.Logf("expected %d nodes, got %d", tc.numNodes, len(ns))
					return false
				}

				// Check for nodes that must be included in RG by hostname, type, version, status, and ID.
				for _, expectedNode := range tc.containNodes {
					if !lo.ContainsBy(
						ns,
						func(n *nodes.NodeItem) bool {
							return n.Hostname == expectedNode.Hostname &&
								n.Type == expectedNode.Type &&
								n.Version == expectedNode.Version &&
								n.Status == expectedNode.Status &&
								n.ID == expectedNode.ID
						}) {
						t.Logf("expected node %+v not found", expectedNode)
						return false
					}
				}
				// Check for nodes that must not be included by hostname and type.
				for _, node := range tc.notContainNodes {
					if lo.ContainsBy(
						ns,
						func(n *nodes.NodeItem) bool {
							return n.Hostname == node.Hostname && n.Type == node.Type
						}) {
						t.Logf("unexpected node %+v found", node)
						return false
					}
				}

				return true
			}, timeout, tick)

			// Cancel the context and wait for the nodeAgent.Start() to return.
			cancel()
			select {
			case <-time.After(timeout):
				t.Fatal("expected the agent to return after the context was cancelled")
			case <-agentReturned:
			}
		})
	}
}

func TestNodeAgent_StartDoesntReturnUntilContextGetsCancelled(t *testing.T) {
	t.Parallel()

	nodeClient := newMockNodeAPIClient(nil)
	// Always return errors from ListNodes to ensure that the agent doesn't propagate it to the Start() caller.
	// ListNodes is the first call made by the agent in Start(), so we care only about this one.
	nodeClient.ReturnErrorFromListAllNodes(true)

	nodeAgent := konnect.NewNodeAgent(
		testHostname,
		testKicVersion,
		konnect.DefaultRefreshNodePeriod,
		logr.Discard(),
		nodeClient,
		newMockConfigStatusNotifier(),
		newMockGatewayInstanceGetter(nil),
		newMockGatewayClientsNotifier(),
		newMockManagerInstanceIDProvider(uuid.New()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	agentReturned := make(chan struct{})
	go func() {
		err := nodeAgent.Start(ctx)
		assert.NoError(t, err, "expected no error even when the context is cancelled")
		close(agentReturned)
	}()

	require.Eventually(t, func() bool {
		return nodeClient.WasListAllNodesCalled()
	}, time.Second, time.Millisecond, "expected list nodes to be called when starting the agent")

	// Ensure that after list nodes returned an error, the agent didn't return.
	select {
	case <-agentReturned:
		t.Fatal("expected the agent to not return yet")
	default:
	}

	// Cancel the context and wait for the nodeAgent.Start() to return.
	cancel()
	select {
	case <-time.After(time.Second):
		t.Fatal("expected the agent to return after the context was cancelled")
	case <-agentReturned:
	}
}
