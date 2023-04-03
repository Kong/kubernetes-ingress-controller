package konnect

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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

// NodeAgent gets the running status of KIC node and controlled kong gateway nodes,
// and update their statuses to konnect.
type NodeAgent struct {
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient *NodeAPIClient
	refreshPeriod time.Duration

	configStatus           atomic.Uint32
	configStatusSubscriber dataplane.ConfigStatusSubscriber

	gatewayInstanceGetter         GatewayInstanceGetter
	gatewayClientsChangesNotifier GatewayClientsChangesNotifier
}

// NewNodeAgent creates a new node agent.
// hostname and version are hostname and version of KIC.
func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	logger logr.Logger,
	client *NodeAPIClient,
	configStatusSubscriber dataplane.ConfigStatusSubscriber,
	gatewayGetter GatewayInstanceGetter,
	gatewayClientsChangesNotifier GatewayClientsChangesNotifier,
) *NodeAgent {
	if refreshPeriod < MinRefreshNodePeriod {
		refreshPeriod = MinRefreshNodePeriod
	}
	a := &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", client.RuntimeGroupID),
		konnectClient:                 client,
		refreshPeriod:                 refreshPeriod,
		configStatusSubscriber:        configStatusSubscriber,
		gatewayInstanceGetter:         gatewayGetter,
		gatewayClientsChangesNotifier: gatewayClientsChangesNotifier,
	}
	a.configStatus.Store(uint32(dataplane.ConfigStatusOK))
	return a
}

// Start runs the process of maintaining and uploading of KIC and kong gateway nodes.
func (a *NodeAgent) Start(ctx context.Context) error {
	a.Logger.Info("Starting Konnect NodeAgent")

	if err := a.updateNodes(ctx); err != nil {
		a.Logger.Error(err, "Failed to run initial update of Konnect nodes, no further updates will be performed")
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

// sortNodesByLastPing sort nodes by descending order of last ping time
// so that nodes are sorted by the newest order.
func sortNodesByLastPing(nodes []*NodeItem) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].LastPing > nodes[j].LastPing
	})
}

// subscribeConfigStatus subscribes and updates KIC status on translating and applying configurations to kong gateway.
func (a *NodeAgent) subscribeConfigStatus(ctx context.Context) {
	ch := a.configStatusSubscriber.SubscribeConfigStatus()
	chDone := ctx.Done()

	for {
		select {
		case <-chDone:
			a.Logger.Info("subscribe loop stopped", "message", ctx.Err().Error())
			return
		case configStatus := <-ch:
			a.configStatus.Store(uint32(configStatus))
		}
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
			a.Logger.Info("subscribe gateway clients changes loop stopped", "message", ctx.Err().Error())
			return
		case <-gatewayClientsChangedCh:
			if err := a.updateNodes(ctx); err != nil {
				a.Logger.Error(err, "failed to update nodes after gateway clients changed")
			}
		}
	}
}

// updateKICNode updates status of KIC node in konnect.
func (a *NodeAgent) updateKICNode(ctx context.Context, existingNodes []*NodeItem) error {
	nodesWithSameName := []*NodeItem{}
	for _, node := range existingNodes {
		if node.Type != NodeTypeIngressController {
			continue
		}

		if node.Hostname == a.Hostname {
			// save all nodes with same name as current KIC node, update the latest one and delete others.
			nodesWithSameName = append(nodesWithSameName, node)
		} else {
			// delete the nodes with different name of the current node, since only on KIC node is allowed in the runtime group.
			a.Logger.V(util.DebugLevel).Info("remove outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.konnectClient.DeleteNode(ctx, node.ID)
			if err != nil {
				a.Logger.Error(err, "failed to delete KIC node", "node_id", node.ID, "hostname", node.Hostname)
				continue
			}
		}
	}
	// sort nodes by last ping and reserve the latest node.
	sortNodesByLastPing(nodesWithSameName)

	var ingressControllerStatus IngressControllerState
	configStatus := int(a.configStatus.Load())
	switch dataplane.ConfigStatus(configStatus) {
	case dataplane.ConfigStatusOK:
		ingressControllerStatus = IngressControllerStateOperational
	case dataplane.ConfigStatusTranslationErrorHappened:
		ingressControllerStatus = IngressControllerStatePartialConfigFail
	case dataplane.ConfigStatusApplyFailed:
		ingressControllerStatus = IngressControllerStateInoperable
	default:
		ingressControllerStatus = IngressControllerStateUnknown
	}

	// create a new node if there is no existing node with same name as the current KIC node.
	if len(nodesWithSameName) == 0 {
		a.Logger.V(util.DebugLevel).Info("no nodes found for KIC pod, should create one", "hostname", a.Hostname)
		createNodeReq := &CreateNodeRequest{
			Hostname: a.Hostname,
			Version:  a.Version,
			Type:     NodeTypeIngressController,
			LastPing: time.Now().Unix(),
			Status:   string(ingressControllerStatus),
		}
		resp, err := a.konnectClient.CreateNode(ctx, createNodeReq)
		if err != nil {
			return fmt.Errorf("failed to create KIC node, hostname %s: %w", a.Hostname, err)
		}
		a.Logger.Info("created KIC node", "node_id", resp.Item.ID, "hostname", a.Hostname)
		return nil
	}

	// update the node with latest last ping time.
	latestNode := nodesWithSameName[0]
	updateNodeReq := &UpdateNodeRequest{
		Hostname: a.Hostname,
		Type:     NodeTypeIngressController,
		Version:  a.Version,
		LastPing: time.Now().Unix(),
		Status:   string(ingressControllerStatus),
	}
	_, err := a.konnectClient.UpdateNode(ctx, latestNode.ID, updateNodeReq)
	if err != nil {
		a.Logger.Error(err, "failed to update node for KIC")
		return err
	}
	a.Logger.V(util.DebugLevel).Info("updated last ping time of node for KIC", "node_id", latestNode.ID, "hostname", a.Hostname)

	// treat more nodes with the same name as outdated, and remove them.
	for i := 1; i < len(nodesWithSameName); i++ {
		node := nodesWithSameName[i]
		err := a.konnectClient.DeleteNode(ctx, node.ID)
		if err != nil {
			a.Logger.Error(err, "failed to delete outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			continue
		}
		a.Logger.V(util.DebugLevel).Info("removed outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
	}
	return nil
}

// updateGatewayNodes updates status of controlled kong gateway nodes to konnect.
func (a *NodeAgent) updateGatewayNodes(ctx context.Context, existingNodes []*NodeItem) error {
	gatewayInstances, err := a.gatewayInstanceGetter.GetGatewayInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to get controlled kong gateway pods: %w", err)
	}
	gatewayInstanceMap := make(map[string]struct{})

	nodeType := NodeTypeKongProxy

	existingNodeMap := make(map[string][]*NodeItem)
	for _, node := range existingNodes {
		if node.Type == nodeType {
			existingNodeMap[node.Hostname] = append(existingNodeMap[node.Hostname], node)
		}
	}

	for _, gateway := range gatewayInstances {
		gatewayInstanceMap[gateway.Hostname] = struct{}{}
		nodes, ok := existingNodeMap[gateway.Hostname]

		// hostname in existing nodes, should create a new node.
		if !ok || len(nodes) == 0 {
			createNodeReq := &CreateNodeRequest{
				ID:       gateway.NodeID,
				Hostname: gateway.Hostname,
				Version:  gateway.Version,
				Type:     nodeType,
				LastPing: time.Now().Unix(),
			}
			newNode, err := a.konnectClient.CreateNode(ctx, createNodeReq)
			if err != nil {
				a.Logger.Error(err, "failed to create kong gateway node", "hostname", gateway.Hostname)
			} else {
				a.Logger.Info("created kong gateway node", "hostname", gateway.Hostname, "node_id", newNode.Item.ID)
			}
			continue
		}

		// sort the nodes by last ping, and only reserve the latest node.
		sortNodesByLastPing(nodes)
		updateNodeReq := &UpdateNodeRequest{
			Hostname: gateway.Hostname,
			Version:  gateway.Version,
			Type:     nodeType,
			LastPing: time.Now().Unix(),
		}
		// update the latest node.
		latestNode := nodes[0]
		_, err := a.konnectClient.UpdateNode(ctx, latestNode.ID, updateNodeReq)
		if err != nil {
			a.Logger.Error(err, "failed to update kong gateway node", "hostname", gateway.Hostname, "node_id", latestNode.ID)
			continue
		}
		a.Logger.V(util.DebugLevel).Info("updated kong gateway node", "hostname", gateway.Hostname, "node_id", latestNode.ID)
		// succeeded to update node, remove the other outdated nodes.
		for i := 1; i < len(nodes); i++ {
			node := nodes[i]
			err := a.konnectClient.DeleteNode(ctx, node.ID)
			if err != nil {
				a.Logger.Error(err, "failed to delete outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
				continue
			}
			a.Logger.V(util.DebugLevel).Info("removed outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
		}

	}

	// delete nodes with no corresponding gateway pod.
	for hostname, nodes := range existingNodeMap {
		if _, ok := gatewayInstanceMap[hostname]; !ok {
			for _, node := range nodes {
				err := a.konnectClient.DeleteNode(ctx, node.ID)
				if err != nil {
					a.Logger.Error(err, "failed to delete outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
					continue
				}
				a.Logger.V(util.DebugLevel).Info("removed outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
			}
		}
	}

	return nil
}

// updateNodes updates current status of KIC and controlled kong gateway nodes.
func (a *NodeAgent) updateNodes(ctx context.Context) error {
	existingNodes, err := a.konnectClient.ListAllNodes(ctx)
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

// updateNodeLoop runs the loop to update status of KIC and kong gateway nods periodically.
func (a *NodeAgent) updateNodeLoop(ctx context.Context) {
	ticker := time.NewTicker(a.refreshPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			a.Logger.Info("update node loop stopped", "message", err.Error())
			return
		case <-ticker.C:
			err := a.updateNodes(ctx)
			if err != nil {
				a.Logger.Error(err, "failed to update nodes")
			}
		}
	}
}

// GatewayClientGetter gets gateway instances from admin API clients.
type GatewayClientGetter struct {
	logger          logr.Logger
	clientsProvider dataplane.AdminAPIClientsProvider
}

var _ GatewayInstanceGetter = &GatewayClientGetter{}

// NewGatewayClientGetter creates a GatewayClientGetter to get gateway instances from client provider.
func NewGatewayClientGetter(logger logr.Logger, clientsProvider dataplane.AdminAPIClientsProvider) *GatewayClientGetter {
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
			p.logger.Error(err, "failed to get kong version")
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
				p.logger.Error(err, "failed to parse URL of gateway admin API from raw URL, skipping", "url", rootURL)
				continue
			}
			// use "gateway_address" as hostname of konnect node.
			hostname = "gateway" + "_" + u.Host
		}

		nodeID, err := client.NodeID(ctx)
		if err != nil {
			p.logger.Error(err, "failed to get node ID from gateway admin API, skipping", "url", client.BaseRootURL())
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
