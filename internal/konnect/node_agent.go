package konnect

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const defaultRefreshNodeInterval = 30 * time.Second

type NodeAgent struct {
	NodeID   string
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient   *Client
	refreshInterval time.Duration
}

func NewNodeAgent(hostname string, version string, logger logr.Logger, client *Client) *NodeAgent {
	return &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", client.RuntimeGroupID),
		konnectClient: client,
		// TODO: set refresh interval by some flag
		// https://github.com/Kong/kubernetes-ingress-controller/issues/3515
		refreshInterval: defaultRefreshNodeInterval,
	}
}

func (a *NodeAgent) createNode() error {
	createNodeReq := &CreateNodeRequest{
		ID:       a.NodeID,
		Hostname: a.Hostname,
		Version:  a.Version,
		Type:     NodeTypeIngressController,
		LastPing: time.Now().Unix(),
	}
	resp, err := a.konnectClient.CreateNode(createNodeReq)
	if err != nil {
		return fmt.Errorf("failed to update node, hostname %s: %w", a.Hostname, err)
	}

	a.NodeID = resp.Item.ID
	a.Logger.V(util.DebugLevel).Info("created node for KIC", "node_id", a.NodeID, "hostname", a.Hostname)
	return nil
}

func (a *NodeAgent) clearOutdatedNodes() error {
	nodes, err := a.konnectClient.ListNodes()
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	for _, node := range nodes.Items {
		if node.Type == NodeTypeIngressController && node.Hostname != a.Hostname {
			a.Logger.V(util.DebugLevel).Info("remove outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.konnectClient.DeleteNode(node.ID)
			if err != nil {
				return fmt.Errorf("failed to delete node %s: %w", node.ID, err)
			}
		}
	}
	return nil
}

func (a *NodeAgent) updateNode() error {
	err := a.clearOutdatedNodes()
	if err != nil {
		a.Logger.Error(err, "failed to clear outdated nodes")
		return err
	}

	// TODO: retrieve the real state of KIC
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3515
	ingressControllerStatus := IngressControllerStateOperational

	updateNodeReq := &UpdateNodeRequest{
		Hostname: a.Hostname,
		Type:     NodeTypeIngressController,
		Version:  a.Version,
		LastPing: time.Now().Unix(),

		Status: string(ingressControllerStatus),
	}
	_, err = a.konnectClient.UpdateNode(a.NodeID, updateNodeReq)
	if err != nil {
		a.Logger.Error(err, "failed to update node for KIC")
		return err
	}
	a.Logger.V(util.DebugLevel).Info("updated last ping time of node for KIC", "node_id", a.NodeID, "hostname", a.Hostname)
	return nil
}

func (a *NodeAgent) updateNodeLoop() {
	ticker := time.NewTicker(a.refreshInterval)
	defer ticker.Stop()
	// TODO: add some mechanism to break the loop
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3515
	for range ticker.C {
		err := a.updateNode()
		if err != nil {
			a.Logger.Error(err, "failed to update node", "node_id", a.NodeID)
		}
	}
}

func (a *NodeAgent) Run() {
	err := a.createNode()
	if err != nil {
		a.Logger.Error(err, "failed to create node, agent abort")
		return
	}
	go a.updateNodeLoop()
}
