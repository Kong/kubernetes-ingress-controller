package konnect

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
)

const (
	MinRefreshNodePeriod     = 30 * time.Second
	DefaultRefreshNodePeriod = 60 * time.Second
)

// GatewayInstance is a controlled kong gateway instance.
// its hostname and version will be used to update status of nodes corresponding to the instance in konnect.
type GatewayInstance struct {
	Hostname string
	Version  string
	NodeID   string
}

// GatewayInstanceGetter is the interface to get currently running gateway instances in the kubernetes cluster.
type GatewayInstanceGetter interface {
	GetGatewayInstances(ctx context.Context) ([]GatewayInstance, error)
}

type GatewayClientsChangesNotifier interface {
	SubscribeToGatewayClientsChanges() (<-chan struct{}, bool)
}

type ManagerInstanceIDProvider interface {
	GetID() uuid.UUID
}

// NodeClient is the interface to Konnect Control Plane Node API.
type NodeClient interface {
	CreateNode(ctx context.Context, req *nodes.CreateNodeRequest) (*nodes.CreateNodeResponse, error)
	UpdateNode(ctx context.Context, nodeID string, req *nodes.UpdateNodeRequest) (*nodes.UpdateNodeResponse, error)
	DeleteNode(ctx context.Context, nodeID string) error
	ListAllNodes(ctx context.Context) ([]*nodes.NodeItem, error)
}

type Ticker interface {
	Stop()
	Channel() <-chan time.Time
	Reset(time.Duration)
}

// NodeAgent gets the running status of KIC node and controlled kong gateway nodes,
// and update their statuses to konnect.
type NodeAgent struct {
	hostname string
	version  string

	logger logr.Logger

	nodeClient    NodeClient
	refreshPeriod time.Duration
	refreshTicker Ticker

	gatewayConfigStatus    clients.GatewayConfigApplyStatus
	konnectConfigStatus    clients.KonnectConfigUploadStatus
	configStatus           atomic.Value
	configStatusSubscriber clients.ConfigStatusSubscriber

	gatewayInstanceGetter         GatewayInstanceGetter
	gatewayClientsChangesNotifier GatewayClientsChangesNotifier
	managerInstanceIDProvider     ManagerInstanceIDProvider
}

type NodeAgentOpt func(*NodeAgent)

// WithRefreshTicker sets the refresh ticker of node agent.
func WithRefreshTicker(ticker Ticker) NodeAgentOpt {
	return func(a *NodeAgent) {
		a.refreshTicker = ticker
	}
}

// NewNodeAgent creates a new node agent.
// hostname and version are hostname and version of KIC.
func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	logger logr.Logger,
	client NodeClient,
	configStatusSubscriber clients.ConfigStatusSubscriber,
	gatewayGetter GatewayInstanceGetter,
	gatewayClientsChangesNotifier GatewayClientsChangesNotifier,
	managerInstanceIDProvider ManagerInstanceIDProvider,
	opts ...NodeAgentOpt,
) *NodeAgent {
	if refreshPeriod < MinRefreshNodePeriod {
		refreshPeriod = MinRefreshNodePeriod
	}
	a := &NodeAgent{
		hostname:                      hostname,
		version:                       version,
		logger:                        logger.WithName("konnect-node-agent"),
		nodeClient:                    client,
		refreshPeriod:                 refreshPeriod,
		refreshTicker:                 clock.NewTicker(),
		configStatusSubscriber:        configStatusSubscriber,
		gatewayInstanceGetter:         gatewayGetter,
		gatewayClientsChangesNotifier: gatewayClientsChangesNotifier,
		managerInstanceIDProvider:     managerInstanceIDProvider,
	}
	a.configStatus.Store(clients.ConfigStatusOK)

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Start runs the process of maintaining and uploading of KIC and kong gateway nodes.
func (a *NodeAgent) Start(ctx context.Context) error {
	a.logger.Info("Starting Konnect NodeAgent")

	if err := a.updateNodes(ctx); err != nil {
		a.logger.Error(err, "Failed to run initial update of Konnect nodes, no further updates will be performed")
		// Do not return here as we don't want NodeAgent to affect the manager health.
		// If we returned, it would cause the manager to exit.
	} else {
		// Run the goroutines only in case we succeeded to run initial update of nodes.
		go a.updateNodeLoop(ctx)
		go a.subscribeConfigStatus(ctx)
		go a.subscribeToGatewayClientsChanges(ctx)
	}

	// We're waiting here as that's the manager.Runnable interface requirement to block until the context is done.
	<-ctx.Done()
	return nil
}

// NeedLeaderElection implements LeaderElectionRunnable interface to ensure that the node agent is run only when
// the KIC instance is elected a leader.
func (a *NodeAgent) NeedLeaderElection() bool {
	return true
}

// updateNodeLoop runs the loop to update status of KIC and kong gateway nods periodically.
func (a *NodeAgent) updateNodeLoop(ctx context.Context) {
	a.refreshTicker.Reset(a.refreshPeriod)
	defer a.refreshTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			a.logger.Info("Update node loop stopped", "message", err.Error())
			return
		case <-a.refreshTicker.Channel():
			a.logger.V(logging.DebugLevel).Info("Updating nodes on tick")
			err := a.updateNodes(ctx)
			if err != nil {
				a.logger.Error(err, "Failed to update nodes")
			}
		}
	}
}

// sortNodesByLastPing sort nodes by descending order of last ping time
// so that nodes are sorted by the newest order.
func sortNodesByLastPing(nodes []*nodes.NodeItem) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].LastPing > nodes[j].LastPing
	})
}

// subscribeConfigStatus subscribes and updates KIC status on translating and applying configurations to kong gateway.
func (a *NodeAgent) subscribeConfigStatus(ctx context.Context) {
	gatewayStatusCh := a.configStatusSubscriber.SubscribeGatewayConfigStatus()
	konnectStatusCh := a.configStatusSubscriber.SubscribeKonnectConfigStatus()
	chDone := ctx.Done()

	for {
		select {
		case <-chDone:
			a.logger.Info("Subscribe loop stopped", "message", ctx.Err().Error())
			return
		case gatewayConfigStatus := <-gatewayStatusCh:
			if a.gatewayConfigStatus != gatewayConfigStatus {
				a.logger.V(logging.DebugLevel).Info("Gateway config status changed")
				a.gatewayConfigStatus = gatewayConfigStatus
				a.maybeUpdateConfigStatus(ctx)
			}
		case konnectConfigStatus := <-konnectStatusCh:
			if a.konnectConfigStatus != konnectConfigStatus {
				a.logger.V(logging.DebugLevel).Info("Konnect config status changed")
				a.konnectConfigStatus = konnectConfigStatus
				a.maybeUpdateConfigStatus(ctx)
			}
		}
	}
}

func (a *NodeAgent) maybeUpdateConfigStatus(ctx context.Context) {
	configStatus := clients.CalculateConfigStatus(a.gatewayConfigStatus, a.konnectConfigStatus)
	if configStatus == a.configStatus.Load() {
		a.logger.V(logging.DebugLevel).Info("Config status not changed, skipping update")
		return
	}
	a.logger.V(logging.DebugLevel).Info("Config status changed, updating nodes")
	a.configStatus.Store(configStatus)
	if err := a.updateNodes(ctx); err != nil {
		a.logger.Error(err, "Failed to update nodes after config status changed")
	}
}

func (a *NodeAgent) subscribeToGatewayClientsChanges(ctx context.Context) {
	gatewayClientsChangedCh, changesAreExpected := a.gatewayClientsChangesNotifier.SubscribeToGatewayClientsChanges()
	if !changesAreExpected {
		// There are no changes of gateway clients going to happen, we don't have to watch them.
		return
	}

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Subscribe gateway clients changes loop stopped", "message", ctx.Err().Error())
			return
		case <-gatewayClientsChangedCh:
			a.logger.V(logging.DebugLevel).Info("Gateway clients changed, updating nodes")
			if err := a.updateNodes(ctx); err != nil {
				a.logger.Error(err, "Failed to update nodes after gateway clients changed")
			}
		}
	}
}

// updateKICNode updates status of KIC node in konnect.
func (a *NodeAgent) updateKICNode(ctx context.Context, existingNodes []*nodes.NodeItem) error {
	nodesWithSameName := []*nodes.NodeItem{}
	for _, node := range existingNodes {
		if node.Type != nodes.NodeTypeIngressController {
			continue
		}

		if node.Hostname == a.hostname {
			// save all nodes with same name as current KIC node, update the latest one and delete others.
			nodesWithSameName = append(nodesWithSameName, node)
		} else {
			// delete the nodes with different name of the current node, since only on KIC node is allowed in the control plane.
			a.logger.V(logging.DebugLevel).Info("Remove outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.nodeClient.DeleteNode(ctx, node.ID)
			if err != nil {
				a.logger.Error(err, "Failed to delete KIC node", "node_id", node.ID, "hostname", node.Hostname)
				continue
			}
		}
	}
	// sort nodes by last ping and reserve the latest node.
	sortNodesByLastPing(nodesWithSameName)

	var ingressControllerStatus nodes.IngressControllerState
	configStatus := a.configStatus.Load().(clients.ConfigStatus)
	switch configStatus {
	case clients.ConfigStatusOK:
		ingressControllerStatus = nodes.IngressControllerStateOperational
	case clients.ConfigStatusTranslationErrorHappened:
		ingressControllerStatus = nodes.IngressControllerStatePartialConfigFail
	case clients.ConfigStatusApplyFailed:
		ingressControllerStatus = nodes.IngressControllerStateInoperable
	case clients.ConfigStatusOKKonnectApplyFailed:
		ingressControllerStatus = nodes.IngressControllerStateOperationalKonnectOutOfSync
	case clients.ConfigStatusTranslationErrorHappenedKonnectApplyFailed:
		ingressControllerStatus = nodes.IngressControllerStatePartialConfigFailKonnectOutOfSync
	case clients.ConfigStatusApplyFailedKonnectApplyFailed:
		ingressControllerStatus = nodes.IngressControllerStateInoperableKonnectOutOfSync
	case clients.ConfigStatusUnknown:
	default:
		ingressControllerStatus = nodes.IngressControllerStateUnknown
	}

	// create a new node if there is no existing node with same name as the current KIC node.
	if len(nodesWithSameName) == 0 {
		a.logger.V(logging.DebugLevel).Info("No nodes found for KIC pod, should create one", "hostname", a.hostname)
		createNodeReq := &nodes.CreateNodeRequest{
			ID:       a.managerInstanceIDProvider.GetID().String(),
			Hostname: a.hostname,
			Version:  a.version,
			Type:     nodes.NodeTypeIngressController,
			LastPing: time.Now().Unix(),
			Status:   string(ingressControllerStatus),
		}
		resp, err := a.nodeClient.CreateNode(ctx, createNodeReq)
		if err != nil {
			return fmt.Errorf("failed to create KIC node, hostname %s: %w", a.hostname, err)
		}
		a.logger.Info("Created KIC node", "node_id", resp.Item.ID, "hostname", a.hostname)
		return nil
	}

	// update the node with latest last ping time.
	latestNode := nodesWithSameName[0]
	updateNodeReq := &nodes.UpdateNodeRequest{
		Hostname: a.hostname,
		Type:     nodes.NodeTypeIngressController,
		Version:  a.version,
		LastPing: time.Now().Unix(),
		Status:   string(ingressControllerStatus),
	}
	_, err := a.nodeClient.UpdateNode(ctx, latestNode.ID, updateNodeReq)
	if err != nil {
		a.logger.Error(err, "Failed to update node for KIC")
		return err
	}
	a.logger.V(logging.DebugLevel).Info("Updated last ping time of node for KIC", "node_id", latestNode.ID, "hostname", a.hostname)

	// treat more nodes with the same name as outdated, and remove them.
	for i := 1; i < len(nodesWithSameName); i++ {
		node := nodesWithSameName[i]
		err := a.nodeClient.DeleteNode(ctx, node.ID)
		if err != nil {
			a.logger.Error(err, "Failed to delete outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			continue
		}
		a.logger.V(logging.DebugLevel).Info("Removed outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
	}
	return nil
}

// updateGatewayNodes updates status of controlled kong gateway nodes to konnect.
func (a *NodeAgent) updateGatewayNodes(ctx context.Context, existingNodes []*nodes.NodeItem) error {
	gatewayInstances, err := a.gatewayInstanceGetter.GetGatewayInstances(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get controlled kong gateway pods: %w", err)
	}
	gatewayInstanceMap := make(map[string]struct{})

	nodeType := nodes.NodeTypeKongProxy

	existingNodeMap := make(map[string][]*nodes.NodeItem)
	for _, node := range existingNodes {
		if node.Type == nodeType {
			existingNodeMap[node.Hostname] = append(existingNodeMap[node.Hostname], node)
		}
	}

	for _, gateway := range gatewayInstances {
		gatewayInstanceMap[gateway.Hostname] = struct{}{}
		ns, ok := existingNodeMap[gateway.Hostname]

		// hostname in existing nodes, should create a new node.
		if !ok || len(ns) == 0 {
			createNodeReq := &nodes.CreateNodeRequest{
				ID:       gateway.NodeID,
				Hostname: gateway.Hostname,
				Version:  gateway.Version,
				Type:     nodeType,
				LastPing: time.Now().Unix(),
			}
			newNode, err := a.nodeClient.CreateNode(ctx, createNodeReq)
			if err != nil {
				a.logger.Error(err, "Failed to create kong gateway node", "hostname", gateway.Hostname)
			} else {
				a.logger.Info("Created kong gateway node", "hostname", gateway.Hostname, "node_id", newNode.Item.ID)
			}
			continue
		}

		// sort the nodes by last ping, and only reserve the latest node.
		sortNodesByLastPing(ns)
		updateNodeReq := &nodes.UpdateNodeRequest{
			Hostname: gateway.Hostname,
			Version:  gateway.Version,
			Type:     nodeType,
			LastPing: time.Now().Unix(),
		}
		// update the latest node.
		latestNode := ns[0]
		_, err := a.nodeClient.UpdateNode(ctx, latestNode.ID, updateNodeReq)
		if err != nil {
			a.logger.Error(err, "Failed to update kong gateway node", "hostname", gateway.Hostname, "node_id", latestNode.ID)
			continue
		}
		a.logger.V(logging.DebugLevel).Info("Updated kong gateway node", "hostname", gateway.Hostname, "node_id", latestNode.ID)
		// succeeded to update node, remove the other outdated nodes.
		for i := 1; i < len(ns); i++ {
			node := ns[i]
			err := a.nodeClient.DeleteNode(ctx, node.ID)
			if err != nil {
				a.logger.Error(err, "Failed to delete outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
				continue
			}
			a.logger.V(logging.DebugLevel).Info("Removed outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
		}

	}

	// delete nodes with no corresponding gateway pod.
	for hostname, ns := range existingNodeMap {
		if _, ok := gatewayInstanceMap[hostname]; !ok {
			for _, node := range ns {
				err := a.nodeClient.DeleteNode(ctx, node.ID)
				if err != nil {
					a.logger.Error(err, "Failed to delete outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
					continue
				}
				a.logger.V(logging.DebugLevel).Info("Removed outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
			}
		}
	}

	return nil
}

// updateNodes updates current status of KIC and controlled kong gateway nodes.
func (a *NodeAgent) updateNodes(ctx context.Context) error {
	// Reset the ticker after updating nodes to make sure the next ticker-triggered update happens after refreshPeriod.
	defer func() {
		a.refreshTicker.Reset(a.refreshPeriod)
	}()

	existingNodes, err := a.nodeClient.ListAllNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list existing nodes: %w", err)
	}

	err = a.updateKICNode(ctx, existingNodes)
	if err != nil {
		return fmt.Errorf("failed to update KIC node: %w", err)
	}

	err = a.updateGatewayNodes(ctx, existingNodes)
	if err != nil {
		return fmt.Errorf("failed to update controlled kong gateway nodes: %w", err)
	}

	return nil
}

// GatewayClientGetter gets gateway instances from admin API clients.
type GatewayClientGetter struct {
	logger          logr.Logger
	clientsProvider clients.AdminAPIClientsProvider
}

var _ GatewayInstanceGetter = &GatewayClientGetter{}

// NewGatewayClientGetter creates a GatewayClientGetter to get gateway instances from client provider.
func NewGatewayClientGetter(logger logr.Logger, clientsProvider clients.AdminAPIClientsProvider) *GatewayClientGetter {
	return &GatewayClientGetter{
		logger:          logger.WithName("gateway-admin-api-getter"),
		clientsProvider: clientsProvider,
	}
}

// GetGatewayInstances gets gateway instances from currently available gateway API clients.
func (p *GatewayClientGetter) GetGatewayInstances(ctx context.Context) ([]GatewayInstance, error) {
	gatewayClients := p.clientsProvider.GatewayClients()
	// TODO: get version of each kong gateway instance behind clients:
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3590
	kongVersion := ""
	if len(gatewayClients) != 0 {
		v, err := gatewayClients[0].GetKongVersion(ctx)
		if err != nil {
			p.logger.Error(err, "Failed to get kong version")
		} else {
			kongVersion = v
		}
	}

	gatewayInstances := make([]GatewayInstance, 0, len(gatewayClients))
	for _, client := range gatewayClients {
		var hostname string
		podNN, ok := client.PodReference()
		if ok {
			hostname = podNN.String()
		} else {
			rootURL := client.BaseRootURL()
			u, err := url.Parse(rootURL)
			if err != nil {
				p.logger.Error(err, "Failed to parse URL of gateway admin API from raw URL, skipping", "url", rootURL)
				continue
			}
			// use "gateway_address" as hostname of konnect node.
			hostname = "gateway" + "_" + u.Host
		}

		nodeID, err := client.NodeID(ctx)
		if err != nil {
			p.logger.Error(err, "Failed to get node ID from gateway admin API, skipping", "url", client.BaseRootURL())
			continue
		}

		gatewayInstances = append(gatewayInstances, GatewayInstance{
			Hostname: hostname,
			Version:  kongVersion,
			NodeID:   nodeID,
		})
	}

	return gatewayInstances, nil
}
