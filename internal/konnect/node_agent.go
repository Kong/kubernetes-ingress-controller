package konnect

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	MinRefreshNodePeriod     = 30 * time.Second
	DefaultRefreshNodePeriod = 60 * time.Second
)

type NodeAgent struct {
	NodeID   string
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient *Client
	refreshPeriod time.Duration

	hasTranslationFailureChan chan bool
	hasTranslationFailure     bool

	sendConfigErrorChan chan error
	sendCondifError     error
}

func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	hasTranslationFailureChan chan bool,
	sendConfigErrorChan chan error,
	logger logr.Logger,
	client *Client,
) *NodeAgent {
	if refreshPeriod < MinRefreshNodePeriod {
		refreshPeriod = MinRefreshNodePeriod
	}
	return &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", client.RuntimeGroupID),
		konnectClient: client,
		refreshPeriod: refreshPeriod,
	}
}

func (a *NodeAgent) createNode() error {
	// REVIEW: consider existing nodes in runtime group as outdated and delete them before creating?
	createNodeReq := &CreateNodeRequest{
		ID:       a.NodeID,
		Hostname: a.Hostname,
		Version:  a.Version,
		Type:     NodeTypeIngressController,
		LastPing: time.Now().Unix(),
	}
	resp, err := a.konnectClient.CreateNode(createNodeReq)
	if err != nil {
		return fmt.Errorf("failed to create node, hostname %s: %w", a.Hostname, err)
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
		// REVIEW: what should the condition be to delete the node in Konnect RG?
		// (1) Do we check the "last update" of the node, and only delete it when the last update is too old(say, 5 mins ago)?
		// (2) What if there is a node with the same name but not the same node exists?
		// for example, When KIC runs in minikube/kind env and whole cluster is stopped then started again.
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

func (a *NodeAgent) calculateStatus() IngressControllerState {
	if a.sendCondifError != nil {
		return IngressControllerStateInoperable
	}
	if a.hasTranslationFailure {
		return IngressControllerStatePartialConfigFail
	}
	return IngressControllerStateOperational
}

func (a *NodeAgent) updateNode() error {
	err := a.clearOutdatedNodes()
	if err != nil {
		// still continue to update the current status if cleanup failed.
		a.Logger.Error(err, "failed to clear outdated nodes")
	}

	ingressControllerStatus := a.calculateStatus()

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
	ticker := time.NewTicker(a.refreshPeriod)
	defer ticker.Stop()
	for range ticker.C {
		err := a.updateNode()
		if err != nil {
			a.Logger.Error(err, "failed to update node", "node_id", a.NodeID)
		}
	}
}

// receiveStatus receives the necessary information to set the status.
func (a *NodeAgent) receiveStatus() {
	for {
		select {
		case a.hasTranslationFailure = <-a.hasTranslationFailureChan:
		case a.sendCondifError = <-a.sendConfigErrorChan:
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
	go a.receiveStatus()
}
