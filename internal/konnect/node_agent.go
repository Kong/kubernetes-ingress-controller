package konnect

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	MinRefreshNodePeriod     = 30 * time.Second
	DefaultRefreshNodePeriod = 60 * time.Second
	NodeOutdateInterval      = 5 * time.Minute
)

type NodeAgent struct {
	NodeID   string
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient *NodeAPIClient
	refreshPeriod time.Duration

	configStatus           atomic.Uint32
	configStatusSubscriber dataplane.ConfigStatusSubscriber
}

func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	logger logr.Logger,
	client *NodeAPIClient,
	configStatusSubscriber dataplane.ConfigStatusSubscriber,
) *NodeAgent {
	if refreshPeriod < MinRefreshNodePeriod {
		refreshPeriod = MinRefreshNodePeriod
	}
	a := &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", client.RuntimeGroupID),
		konnectClient:          client,
		refreshPeriod:          refreshPeriod,
		configStatusSubscriber: configStatusSubscriber,
	}
	a.configStatus.Store(uint32(dataplane.ConfigStatusOK))
	return a
}

func (a *NodeAgent) Start(ctx context.Context) error {
	err := a.createNode()
	if err != nil {
		return fmt.Errorf("failed creating a node: %w", err)
	}
	go a.updateNodeLoop(ctx)
	go a.subscribeConfigStatus(ctx)

	// We're waiting here as that's the manager.Runnable interface requirement to block until the context is done.
	<-ctx.Done()
	return nil
}

// NeedLeaderElection implements LeaderElectionRunnable interface to ensure that the node agent is run only when
// the KIC instance is elected a leader.
func (a *NodeAgent) NeedLeaderElection() bool {
	return true
}

func (a *NodeAgent) createNode() error {
	err := a.clearOutdatedNodes()
	if err != nil {
		// still continue to update the current status if cleanup failed.
		a.Logger.Error(err, "failed to clear outdated nodes")
	}

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
		deleteNode := false
		if node.Type == NodeTypeIngressController {
			// nodes to remove:
			// (1) since only one KIC node is allowed in a runtime group, all the nodes with other hostnames are considered outdated.
			// (2) in some cases(kind/minikube restart), rebuilt pod uses the same name. So nodes updated for >5mins before should be deleted.
			if node.Hostname != a.Hostname || time.Since(time.Unix(node.UpdatedAt, 0)) > NodeOutdateInterval {
				deleteNode = true
			}
		}
		if deleteNode {
			a.Logger.V(util.DebugLevel).Info("remove outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.konnectClient.DeleteNode(node.ID)
			if err != nil {
				return fmt.Errorf("failed to delete node %s: %w", node.ID, err)
			}
		}
	}
	return nil
}

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

func (a *NodeAgent) updateNode() error {
	err := a.clearOutdatedNodes()
	if err != nil {
		// still continue to update the current status if cleanup failed.
		a.Logger.Error(err, "failed to clear outdated nodes")
	}

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
			err := a.updateNode()
			if err != nil {
				a.Logger.Error(err, "failed to update node", "node_id", a.NodeID)
			}
		}
	}
}
