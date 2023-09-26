package graph

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/kong/deck/file"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Entity struct {
	Name string
	Type string
	Raw  any
}

type EntityHash string

type KongConfigGraph = graph.Graph[EntityHash, Entity]

func hashEntity(entity Entity) EntityHash {
	return EntityHash(entity.Type + ":" + entity.Name)
}

func RenderGraphSVG(g KongConfigGraph, outFilePath string) (string, error) {
	if outFilePath == "" {
		outFile, err := os.CreateTemp("", "*.svg")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
		defer outFile.Close()
		outFilePath = outFile.Name()
	}
	f, err := os.CreateTemp("", "*.dot")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	err = draw.DOT(g, f)
	if err != nil {
		return "", fmt.Errorf("failed to render dot file: %w", err)
	}

	if err = exec.Command("dot", "-Tsvg", "-o", outFilePath, f.Name()).Run(); err != nil {
		return "", fmt.Errorf("failed to render svg file: %w", err)
	}
	return outFilePath, nil
}

// FindConnectedComponents iterates over the graph vertices and runs a DFS on each vertex that has not been visited yet.
func FindConnectedComponents(g KongConfigGraph) ([]KongConfigGraph, error) {
	pm, err := g.PredecessorMap()
	if err != nil {
		return nil, err
	}

	var components []KongConfigGraph
	visited := sets.New[EntityHash]()
	for vertex := range pm {
		if visited.Has(vertex) {
			continue // it was already visited
		}
		component := graph.NewLike[EntityHash, Entity](g)
		if err := graph.DFS[EntityHash, Entity](g, vertex, func(visitedHash EntityHash) bool {
			visitedVertex, err := g.Vertex(visitedHash)
			if err != nil {
				return false // continue DFS, should never happen
			}
			if err := component.AddVertex(visitedVertex); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return false // continue DFS, should never happen
			}
			visited.Insert(visitedHash)
			return false // continue DFS
		}); err != nil {
			return nil, err
		}

		edges, err := g.Edges()
		if err != nil {
			return nil, err
		}

		// TODO: Might skip edges that were already added?
		for _, edge := range edges {
			_, sourceErr := component.Vertex(edge.Source)
			_, targetErr := component.Vertex(edge.Target)
			if sourceErr == nil && targetErr == nil {
				if err := component.AddEdge(edge.Source, edge.Target); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, err
				}
			}
		}

		components = append(components, component)
	}

	return components, nil
}

func BuildKongConfigGraph(config *file.Content) (KongConfigGraph, error) {
	g := graph.New(hashEntity)

	for _, caCert := range config.CACertificates {
		ecac := Entity{Name: *caCert.ID, Type: "ca-certificate", Raw: caCert.DeepCopy()}
		if err := g.AddVertex(ecac); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, err
		}
	}

	for _, service := range config.Services {
		es := Entity{Name: *service.Name, Type: "service"}
		if err := g.AddVertex(es); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, err
		}

		for _, route := range service.Routes {
			er := Entity{Name: *route.Name, Type: "route", Raw: route.DeepCopy()}
			if err := g.AddVertex(er); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(es), hashEntity(er)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}

		if service.ClientCertificate != nil {
			ecc := Entity{Name: *service.ClientCertificate.ID, Type: "certificate", Raw: service.ClientCertificate.DeepCopy()}
			if err := g.AddVertex(ecc); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(es), hashEntity(ecc)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}

		for _, caCert := range service.CACertificates {
			if err := g.AddEdge(hashEntity(es), hashEntity(Entity{Name: *caCert, Type: "ca-certificate"})); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
	}

	for _, upstream := range config.Upstreams {
		eu := Entity{Name: *upstream.Name, Type: "upstream"}
		if err := g.AddVertex(eu); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
			return nil, err
		}

		for _, target := range upstream.Targets {
			et := Entity{Name: *target.Target.Target, Type: "target"}
			if err := g.AddVertex(et); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(eu), hashEntity(et)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
	}

	for _, certificate := range config.Certificates {
		ec := Entity{Name: *certificate.ID, Type: "certificate"}
		if err := g.AddVertex(ec); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, err
		}
		for _, sni := range certificate.SNIs {
			esni := Entity{Name: *sni.Name, Type: "sni"}
			if err := g.AddVertex(esni); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(ec), hashEntity(esni)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
	}

	for _, consumer := range config.Consumers {
		if err := g.AddVertex(Entity{Name: *consumer.Username, Type: "consumer"}); err != nil {
			return nil, err
		}
		// TODO: handle consumer credentials
	}

	for _, plugin := range config.Plugins {
		// TODO: should we resolve edges for plugins that are not enabled?

		// TODO: should we resolve edges for plugins that refer other entities (e.g. mtls-auth -> ca_certificate)?

		// TODO: how to identify Plugins uniquely when no ID nor instance name is present? If we use Plugin.Name,
		// we will have just one vertex per plugin type, which could result in unwanted connections
		// (e.g. broken Service1 <-> Plugin <-> Service2 where Service1 and Service2 should not be connected).

		if plugin.InstanceName == nil {
			rel := util.Rel{}
			if plugin.Service != nil {
				rel.Service = *plugin.Service.ID
			}
			if plugin.Route != nil {
				rel.Route = *plugin.Route.ID
			}
			if plugin.Consumer != nil {
				rel.Consumer = *plugin.Consumer.Username
			}
			plugin.InstanceName = lo.ToPtr(kongstate.PluginInstanceName(*plugin.Name, sets.New[string](), rel))
		}
		ep := Entity{Name: *plugin.InstanceName, Type: "plugin"}
		if err := g.AddVertex(ep); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, err
		}

		if plugin.Service != nil {
			es := Entity{Name: *plugin.Service.ID, Type: "service"}
			if err := g.AddVertex(es); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(ep), hashEntity(es)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
		if plugin.Route != nil {
			er := Entity{Name: *plugin.Route.ID, Type: "route"}
			if err := g.AddVertex(er); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(ep), hashEntity(er)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
		if plugin.Consumer != nil {
			ec := Entity{Name: *plugin.Consumer.Username, Type: "consumer"}
			if err := g.AddVertex(ec); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(ep), hashEntity(ec)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}
	}

	return g, nil
}

// TODO: do we have to support full history or just the latest good config?
func BuildFallbackKongConfig(
	history []KongConfigGraph,
	currentConfig KongConfigGraph,
	entityErrors []sendconfig.FlatEntityError,
) (KongConfigGraph, error) {
	if len(history) == 0 {
		return nil, errors.New("history is empty")
	}
	if len(entityErrors) == 0 {
		return nil, errors.New("entityErrors is empty")
	}
	latestGoodConfig := history[len(history)-1]

	affectedEntities := lo.Map(entityErrors, func(ee sendconfig.FlatEntityError, _ int) EntityHash {
		return hashEntity(Entity{Name: ee.Name, Type: ee.Type})
	})

	currentConnectedComponents, err := FindConnectedComponents(currentConfig)
	if err != nil {
		return nil, fmt.Errorf("could not find connected components of the current config")
	}
	latestGoodConnectedComponents, err := FindConnectedComponents(latestGoodConfig)
	if err != nil {
		return nil, fmt.Errorf("could not find connected components of the latest good config")
	}

	fallbackConfig, err := currentConfig.Clone()
	if err != nil {
		return nil, fmt.Errorf("could not clone current config")
	}
	// We need to remove all connected components that contain affected entities.
	for _, affectedEntity := range affectedEntities {
		connectedComponent, err := findConnectedComponentContainingEntity(currentConnectedComponents, affectedEntity)
		if err != nil {
			return nil, fmt.Errorf("could not find connected component containing entity %s", affectedEntity)
		}

		if err := removeConnectedComponentFromGraph(fallbackConfig, connectedComponent); err != nil {
			return nil, fmt.Errorf("could not remove connected component from graph")
		}
	}

	// We need to add all connected components that contain affected entities from the latest good config.
	for _, affectedEntity := range affectedEntities {
		latestGoodComponent, err := findConnectedComponentContainingEntity(latestGoodConnectedComponents, affectedEntity)
		if err != nil {
			// TODO: If there's no connected component in the latest good config for the broken entity, we can skip it, right?
			continue
		}
		if err := addConnectedComponentToGraph(fallbackConfig, latestGoodComponent); err != nil {
			return nil, fmt.Errorf("could not add connected component to graph: %w", err)
		}
	}

	return fallbackConfig, nil
}

func addConnectedComponentToGraph(g KongConfigGraph, component KongConfigGraph) error {
	adjacencyMap, err := component.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map of connected component: %w", err)
	}

	for hash := range adjacencyMap {
		vertex, err := component.Vertex(hash)
		if err != nil {
			return fmt.Errorf("failed to get vertex %v: %w", hash, err)
		}
		_ = g.AddVertex(vertex)
	}

	edges, err := component.Edges()
	if err != nil {
		return fmt.Errorf("failed to get edges: %w", err)
	}
	for _, edge := range edges {
		_ = g.AddEdge(edge.Source, edge.Target)
	}

	return nil
}

func removeConnectedComponentFromGraph(g KongConfigGraph, component KongConfigGraph) error {
	adjacencyMap, err := component.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map of connected component")
	}
	for vertex, neighbours := range adjacencyMap {
		// First remove all edges from the vertex to its neighbours.
		for neighbour := range neighbours {
			_ = g.RemoveEdge(vertex, neighbour)
		}
		_ = g.RemoveVertex(vertex)
	}
	return nil
}

func findConnectedComponentContainingEntity(components []KongConfigGraph, entityHash EntityHash) (KongConfigGraph, error) {
	for _, component := range components {
		_, err := component.Vertex(entityHash)
		if err == nil {
			return component, nil
		}
	}

	return nil, fmt.Errorf("could not find connected component containing entity %s", entityHash)
}
