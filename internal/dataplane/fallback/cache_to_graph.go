package fallback

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// DefaultCacheGraphProvider is a default implementation of the CacheGraphProvider interface.
type DefaultCacheGraphProvider struct{}

func NewDefaultCacheGraphProvider() *DefaultCacheGraphProvider {
	return &DefaultCacheGraphProvider{}
}

// CacheToGraph creates a new ConfigGraph from the given cache stores. It adds all objects
// from the cache stores to the graph as vertices as well as edges between objects and their dependencies
// resolved by the ResolveDependencies function.
func (p DefaultCacheGraphProvider) CacheToGraph(c store.CacheStores) (*ConfigGraph, error) {
	g := NewConfigGraph()

	for _, s := range c.ListAllStores() {
		if s == nil {
			continue
		}
		for _, o := range s.List() {
			obj, ok := o.(client.Object)
			if !ok {
				// Should not happen since all objects in the cache are client.Objects, but better safe than sorry.
				return nil, fmt.Errorf("expected client.Object, got %T", o)
			}
			// Add the object to the graph. It can happen that the object is already in the graph (i.e. was already added
			// as a dependency of another object), in which case we ignore the error.
			if err := g.AddVertex(obj); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return nil, fmt.Errorf("failed to add %s to the graph: %w", GetObjectHash(obj), err)
			}

			deps, err := ResolveDependencies(c, obj)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve dependencies for %s: %w", GetObjectHash(obj), err)
			}
			// Add the object's dependencies to the graph.
			for _, dep := range deps {
				// Add the dependency to the graph in case it wasn't added before. If it was added before, we ignore the
				// error.
				if err := g.AddVertex(dep); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
					return nil, fmt.Errorf("failed to add %s to the graph: %w", GetObjectHash(obj), err)
				}

				// Add an edge from a dependency to the object. If the edge was already added before, we ignore the error.
				// It's on purpose that we add the edge from the dependency to the object, as it makes it easier to traverse
				// the graph from the object to its dependants once it is broken.
				if err := g.AddEdge(GetObjectHash(dep), GetObjectHash(obj)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge from %s to %s: %w", GetObjectHash(obj), GetObjectHash(dep), err)
				}
			}
		}
	}

	return g, nil
}
