package fallback_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
)

func TestConfigGraph_SubgraphObjects(t *testing.T) {
	var (
		A = NewMockObject("A")
		B = NewMockObject("B")
		C = NewMockObject("C")
		D = NewMockObject("D")
		E = NewMockObject("E")
		F = NewMockObject("F")
		G = NewMockObject("G")
		H = NewMockObject("H")
	)

	// Graph structure (edges define dependency -> dependant relationship):
	//     ┌───┐     ┌───┐   ┌───┐   ┌───┐
	//     │ A │     │ E │   │ F │   │ G │
	//     └─┬─┘     └─┬─┘   └─┬─┘   └─┬─┘
	//       │         │       │       │
	//   ┌───┴───┐     │       └───┬───┘
	//   │       │     │           │
	// ┌─▼─┐   ┌─▼─┐   │         ┌─▼─┐
	// │ B │   │ C │   │         │ H │
	// └───┘   └─┬─┘   │         └───┘
	//           │     │
	//           ├─────┘
	//           │
	//         ┌─▼─┐
	//         │ D │
	//         └───┘
	g, err := NewGraphBuilder().
		WithVertices(A, B, C, D, E).
		WithEdge(A, B).
		WithEdge(A, C).
		WithEdge(C, D).
		WithEdge(E, D).
		WithVertices(F, G, H).
		WithEdge(F, H).
		WithEdge(G, H).
		Build()
	require.NoError(t, err)

	objects, err := g.SubgraphObjects(fallback.GetObjectHash(A))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{A, B, C, D}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(B))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{B}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(C))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{C, D}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(D))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{D}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(E))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{E, D}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(F))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{F, H}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(G))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{G, H}, objects)

	objects, err = g.SubgraphObjects(fallback.GetObjectHash(H))
	require.NoError(t, err)
	require.ElementsMatch(t, []client.Object{H}, objects)
}
