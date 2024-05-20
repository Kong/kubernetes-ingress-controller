package fallback_test

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

const (
	testNamespace = "test-namespace"
)

func testIngressClass(t *testing.T, name string) *netv1.IngressClass {
	return helpers.WithTypeMeta(t, &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}

func testService(t *testing.T, name string) *corev1.Service {
	return helpers.WithTypeMeta(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongServiceFacade(t *testing.T, name string) *incubatorv1alpha1.KongServiceFacade {
	return helpers.WithTypeMeta(t, &incubatorv1alpha1.KongServiceFacade{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongPlugin(t *testing.T, name string) *kongv1.KongPlugin {
	return helpers.WithTypeMeta(t, &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongClusterPlugin(t *testing.T, name string) *kongv1.KongClusterPlugin {
	return helpers.WithTypeMeta(t, &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

// GraphBuilder is a helper to build a graph for testing.
type GraphBuilder struct {
	vertices []client.Object
	edges    map[fallback.ObjectHash][]fallback.ObjectHash
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		edges: make(map[fallback.ObjectHash][]fallback.ObjectHash),
	}
}

// WithVertices adds vertices to the graph.
func (b *GraphBuilder) WithVertices(objs ...client.Object) *GraphBuilder {
	b.vertices = append(b.vertices, objs...)
	return b
}

// WithEdge adds an edge between two vertices in the graph.
func (b *GraphBuilder) WithEdge(from, to client.Object) *GraphBuilder {
	fromHash := fallback.GetObjectHash(from)
	toHash := fallback.GetObjectHash(to)
	b.edges[fromHash] = append(b.edges[fromHash], toHash)
	return b
}

// Build builds the graph.
func (b *GraphBuilder) Build() (*fallback.ConfigGraph, error) {
	g := fallback.NewConfigGraph()

	for _, v := range b.vertices {
		if err := g.AddVertex(v); err != nil {
			return nil, fmt.Errorf("failed to add vertex %s to the graph: %w", fallback.GetObjectHash(v), err)
		}
	}

	for from, tos := range b.edges {
		for _, to := range tos {
			if err := g.AddEdge(from, to); err != nil {
				return nil, fmt.Errorf("failed to add edge from %s to %s: %w", from, to, err)
			}
		}
	}

	return g, nil
}

var _ client.Object = &MockObject{}

// MockObject is a mock object that implements the client.Object interface.
type MockObject struct {
	metav1.ObjectMeta
	metav1.TypeMeta
}

// NewMockObject creates a new mock object with the given name.
func NewMockObject(name string) *MockObject {
	return &MockObject{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// DeepCopyObject is required for runtime.Object interface.
func (m *MockObject) DeepCopyObject() runtime.Object {
	return m
}