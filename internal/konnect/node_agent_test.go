package konnect_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

const (
	testKicVersion = "2.9.0"
	testHostname   = "ingress-0"
)

// testKongVersion matches enterprise version format.
var testKongVersion = fmt.Sprintf("%s.0", versions.KICv3VersionCutoff)

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
	gatewayStatusCh chan clients.GatewayConfigApplyStatus
	konnetStatusCh  chan clients.KonnectConfigUploadStatus
	ch              chan clients.ConfigStatus
}

func newMockConfigStatusNotifier() *mockConfigStatusQueue {
	return &mockConfigStatusQueue{
		gatewayStatusCh: make(chan clients.GatewayConfigApplyStatus),
		konnetStatusCh:  make(chan clients.KonnectConfigUploadStatus),
		ch:              make(chan clients.ConfigStatus),
	}
}

func (m mockConfigStatusQueue) SubscribeConfigStatus() chan clients.ConfigStatus {
	return m.ch
}

func (m mockConfigStatusQueue) SubscribeGatewayConfigStatus() chan clients.GatewayConfigApplyStatus {
	return m.gatewayStatusCh
}

func (m mockConfigStatusQueue) SubscribeKonnectConfigStatus() chan clients.KonnectConfigUploadStatus {
	return m.konnetStatusCh
}

func (m mockConfigStatusQueue) NotifyGatewayConfigStatus(status clients.GatewayConfigApplyStatus) {
	m.gatewayStatusCh <- status
}

func (m mockConfigStatusQueue) NotifyKonnectConfigStatus(status clients.KonnectConfigUploadStatus) {
	m.konnetStatusCh <- status
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
		gatewayConfigStatus *clients.GatewayConfigApplyStatus
		gatewayInstances    []konnect.GatewayInstance

		containNodes    []*nodes.NodeItem
		notContainNodes []*nodes.NodeItem
		numNodes        int
	}{
		{
			name: "create kic node",
			// no existing nodes
			initialNodesInNodeAPI: nil,
			gatewayConfigStatus:   lo.ToPtr(clients.GatewayConfigApplyStatus{}),
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
			gatewayConfigStatus: lo.ToPtr(clients.GatewayConfigApplyStatus{TranslationFailuresOccurred: true}),
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
			nodeClient := newMockNodeClient(tc.initialNodesInNodeAPI)
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

			runAgent(t, nodeAgent)

			if tc.gatewayConfigStatus != nil {
				configStatusQueue.NotifyGatewayConfigStatus(*tc.gatewayConfigStatus)
			}

			require.Eventually(t, func() bool {
				// Notify gateway clients changes.
				gatewayClientsChangesNotifier.Notify()

				// Check number of nodes in CP.
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
		})
	}
}

func TestNodeAgent_StartDoesntReturnUntilContextGetsCancelled(t *testing.T) {
	nodeClient := newMockNodeClient(nil)
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

func TestNodeAgent_ControllerNodeStatusGetsUpdatedOnStatusNotification(t *testing.T) {
	nodeClient := newMockNodeClient(nil)
	configStatusQueue := newMockConfigStatusNotifier()
	gatewayClientsChangesNotifier := newMockGatewayClientsNotifier()

	nodeAgent := konnect.NewNodeAgent(
		testHostname,
		testKicVersion,
		konnect.DefaultRefreshNodePeriod,
		logr.Discard(),
		nodeClient,
		configStatusQueue,
		newMockGatewayInstanceGetter(nil),
		gatewayClientsChangesNotifier,
		newMockManagerInstanceIDProvider(uuid.New()),
	)

	runAgent(t, nodeAgent)

	testCases := []struct {
		expectedConfigStaus         clients.ConfigStatus
		notifiedGatewayConfigStatus clients.GatewayConfigApplyStatus
		notifiedKonnectConfigStatus clients.KonnectConfigUploadStatus
		expectedControllerState     nodes.IngressControllerState
	}{
		{
			expectedConfigStaus:         clients.ConfigStatusOK,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{},
			expectedControllerState:     nodes.IngressControllerStateOperational,
		},
		{
			expectedConfigStaus: clients.ConfigStatusTranslationErrorHappened,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{
				TranslationFailuresOccurred: true,
			},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{},
			expectedControllerState:     nodes.IngressControllerStatePartialConfigFail,
		},
		{
			expectedConfigStaus: clients.ConfigStatusApplyFailed,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{
				ApplyConfigFailed: true,
			},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{},
			expectedControllerState:     nodes.IngressControllerStateInoperable,
		},
		{
			expectedConfigStaus:         clients.ConfigStatusOKKonnectApplyFailed,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{
				Failed: true,
			},
			expectedControllerState: nodes.IngressControllerStateOperationalKonnectOutOfSync,
		},
		{
			expectedConfigStaus: clients.ConfigStatusTranslationErrorHappenedKonnectApplyFailed,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{
				TranslationFailuresOccurred: true,
			},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{
				Failed: true,
			},
			expectedControllerState: nodes.IngressControllerStatePartialConfigFailKonnectOutOfSync,
		},
		{
			expectedConfigStaus: clients.ConfigStatusApplyFailedKonnectApplyFailed,
			notifiedGatewayConfigStatus: clients.GatewayConfigApplyStatus{
				ApplyConfigFailed: true,
			},
			notifiedKonnectConfigStatus: clients.KonnectConfigUploadStatus{
				Failed: true,
			},
			expectedControllerState: nodes.IngressControllerStateInoperableKonnectOutOfSync,
		},
	}

	// expectedNodesUpdatesCount := 0
	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.expectedConfigStaus), func(t *testing.T) {
			configStatusQueue.NotifyGatewayConfigStatus(tc.notifiedGatewayConfigStatus)
			configStatusQueue.NotifyKonnectConfigStatus(tc.notifiedKonnectConfigStatus)

			require.Eventually(t, func() bool {
				controllerNode, ok := lo.Find(nodeClient.MustAllNodes(), func(n *nodes.NodeItem) bool {
					return n.Type == nodes.NodeTypeIngressController
				})
				if !ok {
					t.Log("controller node not found")
					return false
				}

				if controllerNode.Status != string(tc.expectedControllerState) {
					t.Logf("expected controller node status to be %q, got %q", tc.expectedControllerState, controllerNode.Status)
					return false
				}

				return true
			}, time.Second, time.Millisecond)

			// TODO: when we let node agent to subscribe gateway config status and konnect config status separately,
			// It will trigger two updates when they are both changed.
			// expectedNodesUpdatesCount++
			// require.Equal(t, expectedNodesUpdatesCount, nodeClient.NodesUpdatesCount(), "expected only one more node update")
		})
	}
}

func TestNodeAgent_ControllerNodeStatusGetsUpdatedOnlyWhenItChanges(t *testing.T) {
	nodeClient := newMockNodeClient(nil)
	configStatusQueue := newMockConfigStatusNotifier()
	gatewayClientsChangesNotifier := newMockGatewayClientsNotifier()

	nodeAgent := konnect.NewNodeAgent(
		testHostname,
		testKicVersion,
		konnect.DefaultRefreshNodePeriod,
		logr.Discard(),
		nodeClient,
		configStatusQueue,
		newMockGatewayInstanceGetter(nil),
		gatewayClientsChangesNotifier,
		newMockManagerInstanceIDProvider(uuid.New()),
	)

	runAgent(t, nodeAgent)

	// We'll use these two to toggle between when we want to trigger an update.
	statusOne := clients.GatewayConfigApplyStatus{}
	statusTwo := clients.GatewayConfigApplyStatus{TranslationFailuresOccurred: true}

	nodesUpdatesCountEventuallyEquals := func(count int) {
		require.Eventually(t, func() bool {
			return nodeClient.NodesUpdatesCount() == count
		}, time.Second, time.Millisecond)
	}

	// Notify the first status and wait for the node to be updated.
	configStatusQueue.NotifyGatewayConfigStatus(statusOne)
	nodesUpdatesCountEventuallyEquals(1)

	// Notify the same status again and ensure that the node wasn't updated.
	configStatusQueue.NotifyGatewayConfigStatus(statusOne)
	nodesUpdatesCountEventuallyEquals(1)

	// Notify the second status and ensure that the node was updated.
	configStatusQueue.NotifyGatewayConfigStatus(statusTwo)
	nodesUpdatesCountEventuallyEquals(2)

	// Notify the same status again and ensure that the node wasn't updated.
	configStatusQueue.NotifyGatewayConfigStatus(statusTwo)
	nodesUpdatesCountEventuallyEquals(2)
}

func TestNodeAgent_TickerResetsOnEveryNodesUpdate(t *testing.T) {
	const halfOfRefreshPeriod = konnect.DefaultRefreshNodePeriod / 2

	t.Run("config status notification", func(t *testing.T) {
		nodeClient := newMockNodeClient(nil)
		configStatusQueue := newMockConfigStatusNotifier()
		gatewayClientsChangesNotifier := newMockGatewayClientsNotifier()

		ticker := mocks.NewTicker()
		nodeAgent := konnect.NewNodeAgent(
			testHostname,
			testKicVersion,
			konnect.DefaultRefreshNodePeriod,
			logr.Discard(),
			nodeClient,
			configStatusQueue,
			newMockGatewayInstanceGetter(nil),
			gatewayClientsChangesNotifier,
			newMockManagerInstanceIDProvider(uuid.New()),
			konnect.WithRefreshTicker(ticker),
		)

		runAgent(t, nodeAgent)

		t.Log("wait for initial nodes update")
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 1 }, time.Second, time.Microsecond)

		t.Log("trigger update with config status notification")
		configStatusQueue.NotifyGatewayConfigStatus(clients.GatewayConfigApplyStatus{ApplyConfigFailed: true})
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 2 }, time.Second, time.Microsecond)

		t.Log("let another half of the period pass - no update should be triggered yet because of the notification")
		ticker.Add(halfOfRefreshPeriod)
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 2 }, time.Second, time.Microsecond)

		t.Log("trigger update with ticker")
		ticker.Add(konnect.DefaultRefreshNodePeriod)
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() > 2 }, time.Second, time.Microsecond)
	})

	t.Run("gateway clients changes notification", func(t *testing.T) {
		nodeClient := newMockNodeClient(nil)
		configStatusQueue := newMockConfigStatusNotifier()
		gatewayClientsChangesNotifier := newMockGatewayClientsNotifier()

		ticker := mocks.NewTicker()
		nodeAgent := konnect.NewNodeAgent(
			testHostname,
			testKicVersion,
			konnect.DefaultRefreshNodePeriod,
			logr.Discard(),
			nodeClient,
			configStatusQueue,
			newMockGatewayInstanceGetter(nil),
			gatewayClientsChangesNotifier,
			newMockManagerInstanceIDProvider(uuid.New()),
			konnect.WithRefreshTicker(ticker),
		)

		runAgent(t, nodeAgent)

		t.Log("wait for initial nodes update")
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 1 }, time.Second, time.Microsecond)
		t.Log("trigger update with gateway clients change notification")
		gatewayClientsChangesNotifier.Notify()
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 2 }, time.Second, time.Microsecond)

		t.Log("let another half of the period pass - no update should be triggered yet because of the notification")
		ticker.Add(halfOfRefreshPeriod)
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 2 }, time.Second, time.Microsecond)

		t.Log("trigger update with ticker")
		ticker.Add(konnect.DefaultRefreshNodePeriod)
		require.Eventually(t, func() bool { return nodeClient.NodesUpdatesCount() == 3 }, time.Second, time.Microsecond)
	})
}

// runAgent runs the agent in a goroutine and cancels the context after the test is done, ensuring that the agent
// doesn't return prematurely.
func runAgent(t *testing.T, nodeAgent *konnect.NodeAgent) {
	ctx, cancel := context.WithCancel(context.Background())

	// To be used as a barrier to ensure that the agent returned after the context was cancelled.
	agentReturned := make(chan struct{})
	go func() {
		err := nodeAgent.Start(ctx)
		require.NoError(t, err, "expected no error even when the context is cancelled")
		close(agentReturned)
	}()

	t.Cleanup(func() {
		cancel()
		select {
		case <-time.After(time.Second):
			t.Fatal("expected the agent to return after the context was cancelled")
		case <-agentReturned:
		}
	})
}
