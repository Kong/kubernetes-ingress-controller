package fallback_test

import (
	"errors"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
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

func testService(t *testing.T, name string, modifiers ...func(s *corev1.Service)) *corev1.Service {
	svc := helpers.WithTypeMeta(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
	for _, mod := range modifiers {
		mod(svc)
	}
	return svc
}

func testSecret(t *testing.T, name string, modifiers ...func(s *corev1.Secret)) *corev1.Secret {
	s := helpers.WithTypeMeta(t, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
	for _, mod := range modifiers {
		mod(s)
	}
	return s
}

func testKongServiceFacade(t *testing.T, name string) *incubatorv1alpha1.KongServiceFacade {
	return helpers.WithTypeMeta(t, &incubatorv1alpha1.KongServiceFacade{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongPlugin(t *testing.T, name string, modifiers ...func(p *kongv1.KongPlugin)) *kongv1.KongPlugin {
	p := helpers.WithTypeMeta(t, &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
	for _, mod := range modifiers {
		mod(p)
	}
	return p
}

func testKongClusterPlugin(t *testing.T, name string) *kongv1.KongClusterPlugin {
	return helpers.WithTypeMeta(t, &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongUpstreamPolicy(t *testing.T, name string, modifiers ...func(kup *kongv1beta1.KongUpstreamPolicy)) *kongv1beta1.KongUpstreamPolicy {
	kup := helpers.WithTypeMeta(t, &kongv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
	for _, mod := range modifiers {
		mod(kup)
	}
	return kup
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

// getFromStore retrieves an object of type T from the given cache store.
func getFromStore[T client.Object](c cache.Store, obj client.Object) (client.Object, error) {
	o, exists, err := c.Get(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	if !exists {
		return nil, errors.New("object not found")
	}
	typedObject, ok := o.(T)
	if !ok {
		return nil, fmt.Errorf("expected object of type %T, got %T", typedObject, o)
	}
	return typedObject, nil
}
