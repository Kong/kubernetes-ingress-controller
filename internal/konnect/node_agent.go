package konnect

import (
	"context"
	"fmt"
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

// REVIEW: define the subscriber here, or internal/adminapi for common usage?
type ConfigStatusSubscriber interface {
	Subscribe() chan dataplane.ConfigStatus
}

type configStatusSubscriber struct {
	ch chan dataplane.ConfigStatus
}

var _ ConfigStatusSubscriber = &configStatusSubscriber{}

func (s *configStatusSubscriber) Subscribe() chan dataplane.ConfigStatus {
	return s.ch
}

func NewConfigStatusSubscriber(ch chan dataplane.ConfigStatus) *configStatusSubscriber {
	return &configStatusSubscriber{ch: ch}
}

type NodeAgent struct {
	NodeID   string
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient *NodeAPIClient
	refreshPeriod time.Duration

	configStatus           dataplane.ConfigStatus
	configStatusSubscriber ConfigStatusSubscriber
}

func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	logger logr.Logger,
	client *NodeAPIClient,
	configStatusSubscriber ConfigStatusSubscriber,
) *NodeAgent {
	if refreshPeriod < MinRefreshNodePeriod {
		refreshPeriod = MinRefreshNodePeriod
	}
	return &NodeAgent{
		Hostname: hostname,
		Version:  version,
		Logger: logger.
			WithName("konnect-node").WithValues("runtime_group_id", client.RuntimeGroupID),
		konnectClient:          client,
		refreshPeriod:          refreshPeriod,
		configStatus:           dataplane.ConfigStatusOK,
		configStatusSubscriber: configStatusSubscriber,
	}
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
			if node.Hostname != a.Hostname || time.Now().Sub(time.Unix(node.UpdatedAt, 0)) > NodeOutdateInterval {
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
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			a.Logger.Info("subscribe loop stopped", "message", err.Error())
			return
		case a.configStatus = <-a.configStatusSubscriber.Subscribe():
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
	switch a.configStatus {
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

func (a *NodeAgent) Run(ctx context.Context) {
	err := a.createNode()
	if err != nil {
		a.Logger.Error(err, "failed to create node, agent abort")
		return
	}
	go a.updateNodeLoop(ctx)
	go a.subscribeConfigStatus(ctx)
}
