package parser

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/kong/deck/file"
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
		outFilePath = path.Join(os.TempDir(), "kong-config-graph.svg")
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

	for _, service := range config.Services {
		es := Entity{Name: *service.Name, Type: "service"}
		if err := g.AddVertex(es); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return nil, err
		}

		for _, route := range service.Routes {
			er := Entity{Name: *route.Name, Type: "route"}
			if err := g.AddVertex(er); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(es), hashEntity(er)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}

		if service.ClientCertificate != nil {
			ecc := Entity{Name: *service.ClientCertificate.ID, Type: "certificate"}
			if err := g.AddVertex(ecc); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(es), hashEntity(ecc)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return nil, err
			}
		}

		for _, caCert := range service.CACertificates {
			ecac := Entity{Name: *caCert, Type: "ca-certificate"}
			if err := g.AddVertex(ecac); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, err
			}
			if err := g.AddEdge(hashEntity(es), hashEntity(ecac)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
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
