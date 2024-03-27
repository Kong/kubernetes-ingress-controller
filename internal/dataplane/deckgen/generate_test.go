package deckgen_test

import (
	"context"
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func TestToDeckContent(t *testing.T) {
	defaultTestParams := func() deckgen.GenerateDeckContentParams {
		return deckgen.GenerateDeckContentParams{
			FormatVersion: "3.0",
		}
	}
	modifiedDefaultTestParams := func(fn func(p *deckgen.GenerateDeckContentParams)) deckgen.GenerateDeckContentParams {
		p := defaultTestParams()
		fn(&p)
		return p
	}

	testCases := []struct {
		name     string
		params   deckgen.GenerateDeckContentParams
		input    *kongstate.KongState
		expected *file.Content
	}{
		{
			name:   "empty",
			params: defaultTestParams(),
			input:  &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: "3.0",
			},
		},
		{
			name: "empty, generate stub entity",
			params: modifiedDefaultTestParams(func(p *deckgen.GenerateDeckContentParams) {
				p.AppendStubEntityWhenConfigEmpty = true
			}),
			input: &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: "3.0",
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
			result := deckgen.ToDeckContent(context.Background(), logrus.New(), tc.input, tc.params)
			require.Equal(t, tc.expected, result)
		})
	}
}
