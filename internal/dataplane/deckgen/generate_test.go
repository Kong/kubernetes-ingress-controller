package deckgen_test

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

func TestToDeckContent(t *testing.T) {
	testCases := []struct {
		name     string
		params   deckgen.GenerateDeckContentParams
		input    *kongstate.KongState
		expected *file.Content
	}{
		{
			name:   "empty",
			params: deckgen.GenerateDeckContentParams{},
			input:  &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: versions.DeckFileFormatVersion,
			},
		},
		{
			name: "empty, generate stub entity",
			params: deckgen.GenerateDeckContentParams{
				AppendStubEntityWhenConfigEmpty: true,
			},
			input: &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: versions.DeckFileFormatVersion,
				Upstreams: []file.FUpstream{
					{
						Upstream: kong.Upstream{
							Name: lo.ToPtr(deckgen.StubUpstreamName),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := deckgen.ToDeckContent(context.Background(), zapr.NewLogger(zap.NewNop()), tc.input, tc.params)
			require.Equal(t, tc.expected, result)
		})
	}
}
