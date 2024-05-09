package fallback

import (
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run ../../../hack/generators/config-graph

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
