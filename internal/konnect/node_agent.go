package konnect

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	MinRefreshNodePeriod     = 30 * time.Second
	DefaultRefreshNodePeriod = 60 * time.Second
	NodeOutdateInterval      = 5 * time.Minute
)

type GatewayPod struct {
	namespace string
	name      string
	version   string
}

type GatewayPodGetter interface {
	GetGatewayPods() ([]GatewayPod, error)
}

type NodeAgent struct {
	Hostname string
	Version  string

	Logger logr.Logger

	konnectClient *NodeAPIClient
	refreshPeriod time.Duration

	configStatus           atomic.Uint32
	configStatusSubscriber dataplane.ConfigStatusSubscriber

	gatewayPodGetter GatewayPodGetter
}

func NewNodeAgent(
	hostname string,
	version string,
	refreshPeriod time.Duration,
	logger logr.Logger,
	client *NodeAPIClient,
	configStatusSubscriber dataplane.ConfigStatusSubscriber,
	gatewayPodGetter GatewayPodGetter,
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
		gatewayPodGetter:       gatewayPodGetter,
	}
	a.configStatus.Store(uint32(dataplane.ConfigStatusOK))
	return a
}

func (a *NodeAgent) Start(ctx context.Context) error {

	err := a.updateNodes()
	if err != nil {
		return fmt.Errorf("failed to run initial update of nodes, agent abort")
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

// sortNodesByLastPing sort nodes by descending order of last ping time
// so that nodes are sorted by the newest order.
func sortNodesByLastPing(nodes []*NodeItem) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].LastPing > nodes[j].LastPing
	})
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

func (a *NodeAgent) updateKICNode(existingNodes []*NodeItem) error {

	nodesWithSameName := []*NodeItem{}
	for _, node := range existingNodes {
		if node.Type != NodeTypeIngressController {
			continue
		}

		if node.Hostname == a.Hostname {
			nodesWithSameName = append(nodesWithSameName, node)
		} else {
			a.Logger.V(util.DebugLevel).Info("remove outdated KIC node", "node_id", node.ID, "hostname", node.Hostname)
			err := a.konnectClient.DeleteNode(node.ID)
			if err != nil {
				a.Logger.Error(err, "failed to delete KIC node", "node_id", node.ID, "hostname", node.Hostname)
				continue
			}
		}
	}

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

	if len(nodesWithSameName) == 0 {
		a.Logger.V(util.DebugLevel).Info("no nodes found for KIC pod, should create one", "hostname", a.Hostname)
		createNodeReq := &CreateNodeRequest{
			Hostname: a.Hostname,
			Version:  a.Version,
			Type:     NodeTypeIngressController,
			LastPing: time.Now().Unix(),
			Status:   string(ingressControllerStatus),
		}
		resp, err := a.konnectClient.CreateNode(createNodeReq)
		if err != nil {
			return fmt.Errorf("failed to create node, hostname %s: %w", a.Hostname, err)
		}
		a.Logger.Info("created node for KIC pod", "node_id", resp.Item.ID, "hostname", a.Hostname)
		return nil
	}

	latestNode := nodesWithSameName[0]
	updateNodeReq := &UpdateNodeRequest{
		Hostname: a.Hostname,
		Type:     NodeTypeIngressController,
		Version:  a.Version,
		LastPing: time.Now().Unix(),
		Status:   string(ingressControllerStatus),
	}
	_, err := a.konnectClient.UpdateNode(latestNode.ID, updateNodeReq)
	if err != nil {
		a.Logger.Error(err, "failed to update node for KIC")
		return err
	}
	a.Logger.V(util.DebugLevel).Info("updated last ping time of node for KIC", "node_id", latestNode.ID, "hostname", a.Hostname)

	// treat more nodes with the same name as outdated, and remove them.
	for i := 1; i < len(nodesWithSameName); i++ {
		node := nodesWithSameName[i]
		err := a.konnectClient.DeleteNode(node.ID)
		if err != nil {
			a.Logger.Error(err, "failed to delete outdated node", "node_id", node.ID, "hostname", node.Hostname)
			continue
		}
		a.Logger.V(util.DebugLevel).Info("remove outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
	}
	return nil
}

func (a *NodeAgent) updateGatewayNodes(existingNodes []*NodeItem) error {
	gatewayPods, err := a.gatewayPodGetter.GetGatewayPods()
	if err != nil {
		return fmt.Errorf("failed to get controlled kong gateway pods: %w", err)
	}
	gatewayPodMap := make(map[string]struct{})

	// TODO: final confirmation on node type used for controlled gateway nodes.
	nodeType := NodeTypeIngressProxy

	existingNodeMap := make(map[string][]*NodeItem)
	for _, node := range existingNodes {
		if node.Type == nodeType {
			existingNodeMap[node.Hostname] = append(existingNodeMap[node.Hostname], node)
		}
	}

	// TODO: generate a hostname if pod namespace and pod name is not available.

	for _, pod := range gatewayPods {
		podNN := pod.namespace + "/" + pod.name
		gatewayPodMap[podNN] = struct{}{}
		nodes, ok := existingNodeMap[podNN]
		if !ok || len(nodes) == 0 {
			// pod name not in existing nodes, should create a new node.
			createNodeReq := &CreateNodeRequest{
				Hostname: podNN,
				Version:  pod.version,
				Type:     nodeType,
				LastPing: time.Now().Unix(),
			}
			newNode, err := a.konnectClient.CreateNode(createNodeReq)
			if err != nil {
				a.Logger.Error(err, "failed to create node for pod", "pod_namespace", pod.namespace, "pod_name", pod.name)
			} else {
				a.Logger.Info("created kong gateway node for pod", "pod_namespace", pod.namespace, "pod_name", pod.name, "node_id", newNode.Item.ID)
			}
		} else {
			// sort the nodes by last ping, and only reserve the latest node.
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].LastPing > nodes[j].LastPing
			})
			updateNodeReq := &UpdateNodeRequest{
				Hostname: podNN,
				Version:  pod.version,
				Type:     nodeType,
				LastPing: time.Now().Unix(),
			}
			latestNode := nodes[0]
			_, err := a.konnectClient.UpdateNode(latestNode.ID, updateNodeReq)
			if err != nil {
				a.Logger.Error(err, "failed to update node for pod", "pod_namespace", pod.namespace, "pod_name", pod.name, "node_id", latestNode.ID)
			} else {
				a.Logger.Info("updated kong gateway node for pod", "pod_namespace", pod.namespace, "pod_name", pod.name, "node_id", latestNode.ID)
				// succeeded to update node, remove the other outdated nodes.
				for i := 1; i < len(nodes); i++ {
					node := nodes[i]
					err := a.konnectClient.DeleteNode(node.ID)
					if err != nil {
						a.Logger.Error(err, "failed to delete outdated node", "node_id", node.ID, "hostname", node.Hostname)
						continue
					}
					a.Logger.V(util.DebugLevel).Info("remove outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
				}
			}

		}
	}

	// delete nodes with no corresponding gateway pod.
	for hostname, nodes := range existingNodeMap {
		if _, ok := gatewayPodMap[hostname]; !ok {
			for _, node := range nodes {
				err := a.konnectClient.DeleteNode(node.ID)
				if err != nil {
					a.Logger.Error(err, "failed to delete outdated node", "node_id", node.ID, "hostname", node.Hostname)
					continue
				}
				a.Logger.V(util.DebugLevel).Info("remove outdated kong gateway node", "node_id", node.ID, "hostname", node.Hostname)
			}
		}
	}

	return nil
}

func (a *NodeAgent) updateNodes() error {
	existingNodes, err := a.konnectClient.ListAllNodes()
	if err != nil {
		return fmt.Errorf("failed to list existing nodes: %w", err)
	}

	err = a.updateKICNode(existingNodes)
	if err != nil {
		// REVIEW: not return here and continue to update kong gateway nodes?
		return fmt.Errorf("failed to update KIC node: %w", err)
	}

	err = a.updateGatewayNodes(existingNodes)
	if err != nil {
		return fmt.Errorf("failed to update controlled kong gateway nodes: %w", err)
	}
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
			err := a.updateNodes()
			if err != nil {
				a.Logger.Error(err, "failed to update nodes")
			}
		}
	}
}

type GatewayEndpointStore struct {
	lock                      sync.RWMutex
	logger                    logr.Logger
	gatewayEndpointsChan      chan []adminapi.DiscoveredAdminAPI
	gatewayEndpointVersionMap map[types.NamespacedName]string
	clientsProvider           dataplane.AdminAPIClientsProvider
}

func NewGatewayEndpointStore(
	ctx context.Context,
	logger logr.Logger,
	initGatewayEndpoints []adminapi.DiscoveredAdminAPI,
	gatewayEndpointsChan chan []adminapi.DiscoveredAdminAPI,
	clientsProvider dataplane.AdminAPIClientsProvider,
) *GatewayEndpointStore {
	gatewayEndpointVersionMap := make(map[types.NamespacedName]string)
	gatewayClients := clientsProvider.GatewayClients()

	// TODO: get the true address for each endpoint and get their versions.
	kongVersion := ""
	if len(gatewayClients) != 0 {
		v, err := gatewayClients[0].GetKongVersion(ctx)
		if err != nil {
			logger.Error(err, "failed to get kong version")
		} else {
			kongVersion = v
		}
	}
	for _, endpoint := range initGatewayEndpoints {
		logger.Info("init endpoint", "namespace", endpoint.PodRef.Namespace, "name", endpoint.PodRef.Name)
		gatewayEndpointVersionMap[endpoint.PodRef] = kongVersion
	}

	s := &GatewayEndpointStore{
		logger:                    logger.WithName("gateway_endpoint_store"),
		gatewayEndpointsChan:      gatewayEndpointsChan,
		gatewayEndpointVersionMap: gatewayEndpointVersionMap,
		clientsProvider:           clientsProvider,
	}

	go s.subscribeEndpointLoop(ctx)
	return s
}

func (s *GatewayEndpointStore) GetGatewayPods() ([]GatewayPod, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	gatewayPods := []GatewayPod{}
	for nsName, version := range s.gatewayEndpointVersionMap {
		gatewayPods = append(gatewayPods, GatewayPod{
			namespace: nsName.Namespace,
			name:      nsName.Name,
			version:   version,
		})
	}

	return gatewayPods, nil
}

func (s *GatewayEndpointStore) updateEndpoints(ctx context.Context, endpoints []adminapi.DiscoveredAdminAPI) {
	s.lock.Lock()
	defer s.lock.Unlock()
	gatewayClients := s.clientsProvider.GatewayClients()
	// TODO: get the true address for each endpoint and get their versions.
	kongVersion := ""
	if len(gatewayClients) != 0 {
		v, err := gatewayClients[0].GetKongVersion(ctx)
		if err != nil {
			s.logger.Error(err, "failed to get kong version")
		} else {
			kongVersion = v
		}
	}
	gatewayEndpointVersionMap := make(map[types.NamespacedName]string)
	for _, endpoint := range endpoints {
		s.logger.Info("updated endpoint", "namespace", endpoint.PodRef.Namespace, "name", endpoint.PodRef.Name)
		gatewayEndpointVersionMap[endpoint.PodRef] = kongVersion
	}
	s.gatewayEndpointVersionMap = gatewayEndpointVersionMap
}

func (s *GatewayEndpointStore) subscribeEndpointLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			s.logger.Info("update node loop stopped", "message", err.Error())
			return
		case endpoints := <-s.gatewayEndpointsChan:
			s.logger.V(util.DebugLevel).Info("update gateway endpoints")
			s.updateEndpoints(ctx, endpoints)
		}
	}
}

var _ GatewayPodGetter = &GatewayEndpointStore{}
