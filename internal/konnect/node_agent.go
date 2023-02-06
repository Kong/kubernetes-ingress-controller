package konnect

import (
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	defaultRefreshNodeInterval = 15 * time.Second
)

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

func (a *NodeAgent) updateNode() error {

	ingressControllerStatus := &IngressControllerStatus{
		State: IngressControllerStateOperational,
	}

	updateNodeReq := &UpdateNodeRequest{
		Hostname: a.Hostname,
		Type:     NodeTypeIngressController,
		Version:  a.Version,
		LastPing: time.Now().Unix(),

		IngressControllerStatus: ingressControllerStatus,
	}
	_, err := a.adminClient.UpdateNode(a.NodeID, updateNodeReq)
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
