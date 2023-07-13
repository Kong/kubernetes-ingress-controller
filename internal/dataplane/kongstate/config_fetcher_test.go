package kongstate_test

import (
	"testing"

	"github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func TestKongRawStateToKongState(t *testing.T) {
	for _, tt := range []struct {
		name              string
		kongRawState      utils.KongRawState
		expectedKongState *kongstate.KongState
	}{
		{
			name: "sanitizes all services, routes, and upstreams and create a KongState out of a KongRawState",
			kongRawState: utils.KongRawState{
				Services: []*kong.Service{
					{
						Name:      kong.String("testService"),
						ID:        kong.String("abc"),
						CreatedAt: kong.Int(100),
					},
				},
				Routes: []*kong.Route{
					{
						Name:      kong.String("testRoute"),
						ID:        kong.String("def"),
						CreatedAt: kong.Int(101),
						Service: &kong.Service{
							ID: kong.String("abc"),
						},
					},
				},
				Upstreams: []*kong.Upstream{
					{
						Name: kong.String("testUpstream"),
						ID:   kong.String("ghi"),
					},
				},
				Targets: []*kong.Target{
					{
						ID:        kong.String("jkl"),
						CreatedAt: kong.Float64(102),
						Weight:    kong.Int(999),
						Upstream: &kong.Upstream{
							ID: kong.String("ghi"),
						},
					},
				},
			},
			expectedKongState: &kongstate.KongState{
				Services: []kongstate.Service{
					{
						Service: kong.Service{
							Name: kong.String("testService"),
						},
						Routes: []kongstate.Route{
							{
								Route: kong.Route{
									Name: kong.String("testRoute"),
								},
							},
						},
					},
				},
				Upstreams: []kongstate.Upstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("testUpstream"),
						},
						Targets: []kongstate.Target{
							{
								Target: kong.Target{
									Weight: kong.Int(999),
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			state := kongstate.KongRawStateToKongState(&tt.kongRawState)
			require.Equal(t, tt.expectedKongState, state)
		})
	}
}
