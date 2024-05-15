package fallback

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// ConfigGraph is a graph representation of the Kubernetes resources kept in the cache.
// Vertices are objects, edges are dependencies between objects (dependency -> dependant).
// It allows to quickly determine which objects are affected by traversing the graph from
// the affected object to its dependants.
//
// If you want to extend the graph with a new object type, you need to ensure ResolveDependencies
// function is implemented for that object type. If your object type has no dependencies, you can ignore it.
type ConfigGraph struct {
	graph graph.Graph[ObjectHash, client.Object]
}

// ObjectHash is a unique identifier for a given object that is used as a vertex key in the graph.
// It could consist of the object's UID only, but we also include the object's kind, namespace and name
// to make it easier to debug and understand the graph.
type ObjectHash struct {
	// UID is the unique identifier of the object.
	UID k8stypes.UID

	// Kind, Namespace and Name are the object's kind, namespace and name - included for debugging purposes.
	Kind      string
	Namespace string
	Name      string
}

// String returns a string representation of the ObjectHash. It intentionally does not include the UID
// as it is not human-readable and is not necessary for debugging purposes.
func (h ObjectHash) String() string {
	if h.Namespace == "" {
		return fmt.Sprintf("%s:%s", h.Kind, h.Name)
	}
	return fmt.Sprintf("%s:%s/%s", h.Kind, h.Namespace, h.Name)
}

// GetObjectHash is a function that returns a unique identifier for a given object that is used as a
// vertex key in the graph.
func GetObjectHash(obj client.Object) ObjectHash {
	return ObjectHash{
		UID:       obj.GetUID(),
		Kind:      obj.GetObjectKind().GroupVersionKind().Kind,
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

// NewConfigGraphFromCacheStores creates a new ConfigGraph from the given cache stores. It adds all objects
// from the cache stores to the graph as vertices as well as edges between objects and their dependencies
// resolved by the ResolveDependencies function.
func NewConfigGraphFromCacheStores(c store.CacheStores) (*ConfigGraph, error) {
	g := graph.New[ObjectHash, client.Object](GetObjectHash, graph.Directed())

	for _, s := range c.ListAllStores() {
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
				if err := g.AddEdge(GetObjectHash(dep), GetObjectHash(obj)); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge from %s to %s: %w", GetObjectHash(obj), GetObjectHash(dep), err)
				}
			}
		}
	}

	return &ConfigGraph{graph: g}, nil
}

// AdjacencyMap is a map of object hashes to their neighbours' hashes.
type AdjacencyMap map[ObjectHash][]ObjectHash

// AdjacencyMap returns a map of object hashes to their neighbours' hashes.
func (g *ConfigGraph) AdjacencyMap() (AdjacencyMap, error) {
	am, err := g.graph.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	m := make(map[ObjectHash][]ObjectHash)
	for v, neighbours := range am {
		m[v] = lo.Keys(neighbours)
	}
	return m, nil
}
