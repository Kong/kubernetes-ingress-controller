package konnect

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var defaultRefreshNodeInterval = 15 * time.Second

type NodeAgent struct {
	NodeID   string
	Hostname string
	Version  string

	Logger logr.Logger

	adminClient     *AdminClient
	refreshInterval time.Duration
}

func NewNodeAgent(hostname string, version string, logger logr.Logger, adminClient *AdminClient) *NodeAgent {
	return &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", adminClient.RuntimeGroupID),
		adminClient: adminClient,
		// TODO: set refresh interval by flags/envvar
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
	resp, err := a.adminClient.CreateNode(createNodeReq)
	if err != nil {
		a.Logger.Error(err, "failed to create node")
		return err
	}

	a.NodeID = resp.Item.ID
	a.Logger.Info("updated last ping time of node for KIC", "node_id", a.NodeID)
	return nil
}

func (a *NodeAgent) clearOutdatedNodes() error {
	nodes, err := a.adminClient.ListNodes()
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	for _, node := range nodes.Items {
		if node.Type == NodeTypeIngressController && node.Hostname != a.Hostname {
			a.Logger.V(util.DebugLevel).Info("remove KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.adminClient.DeleteNode(node.ID)
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
	ingressControllerStatus := IngressControllerStateOperational

	updateNodeReq := &UpdateNodeRequest{
		Hostname: a.Hostname,
		Type:     NodeTypeIngressController,
		Version:  a.Version,
		LastPing: time.Now().Unix(),

		Status: string(ingressControllerStatus),
	}
	_, err = a.adminClient.UpdateNode(a.NodeID, updateNodeReq)
	if err != nil {
		a.Logger.Error(err, "failed to update node for KIC")
		return err
	}
	a.Logger.V(util.DebugLevel).Info("updated last ping time of node for KIC", "node_id", a.NodeID)
	return nil
}

func (a *NodeAgent) updateNodeLoop() {
	ticker := time.NewTicker(a.refreshInterval)
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
