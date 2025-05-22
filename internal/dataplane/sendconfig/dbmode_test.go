package sendconfig

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/state"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
)

func TestRefillPluginIDs(t *testing.T) {
	testCases := []struct {
		name             string
		currentState     *state.KongState
		targetState      *state.KongState
		expectedPluginID string
	}{
		{
			name: "plugin attached to the same service should be considered as the same plugin",
			currentState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-1"),
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
			}),
			targetState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-2"),
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
			}),
			expectedPluginID: "plugin-1",
		},
		{
			name: "plugin attached to different services should not be considered as the same plugin",
			currentState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-1"),
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
			}),
			targetState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-2"),
						ID:   kong.String("service-2"),
						Host: kong.String("kong.test"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-2"),
						Service: &kong.Service{
							Name: kong.String("service-2"),
							ID:   kong.String("service-2"),
						},
					},
				},
			}),
			expectedPluginID: "plugin-2",
		},
		{
			name: "plugin attached to the same comblination of route and consumer should be considered as the same plugin",
			currentState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Routes: []*kong.Route{
					{
						Name:  kong.String("route-1"),
						ID:    kong.String("route-1"),
						Paths: []*string{kong.String("/")},
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
				Consumers: []*kong.Consumer{
					{
						Username: kong.String("consumer-1"),
						ID:       kong.String("consumer-1"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-1"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer-1"),
						},
						Route: &kong.Route{
							ID: kong.String("route-1"),
						},
					},
				},
			}),
			targetState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Routes: []*kong.Route{
					{
						Name:  kong.String("route-1"),
						ID:    kong.String("route-1"),
						Paths: []*string{kong.String("/")},
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
				Consumers: []*kong.Consumer{
					{
						Username: kong.String("consumer-1"),
						ID:       kong.String("consumer-1"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-2"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer-1"),
						},
						Route: &kong.Route{
							ID: kong.String("route-1"),
						},
					},
				},
			}),
			expectedPluginID: "plugin-1",
		},
		{
			name: "plugin attached to the same route but different consumers should not be considered as the same plugin",
			currentState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Routes: []*kong.Route{
					{
						Name:  kong.String("route-1"),
						ID:    kong.String("route-1"),
						Paths: []*string{kong.String("/")},
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
				Consumers: []*kong.Consumer{
					{
						Username: kong.String("consumer-1"),
						ID:       kong.String("consumer-1"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-1"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer-1"),
						},
						Route: &kong.Route{
							ID: kong.String("route-1"),
						},
					},
				},
			}),
			targetState: mustNewKongStateFromRawState(t, &deckutils.KongRawState{
				Services: []*kong.Service{
					{
						Name: kong.String("service-1"),
						ID:   kong.String("service-1"),
						Host: kong.String("kong.test"),
					},
				},
				Routes: []*kong.Route{
					{
						Name:  kong.String("route-1"),
						ID:    kong.String("route-1"),
						Paths: []*string{kong.String("/")},
						Service: &kong.Service{
							Name: kong.String("service-1"),
							ID:   kong.String("service-1"),
						},
					},
				},
				Consumers: []*kong.Consumer{
					{
						Username: kong.String("consumer-2"),
						ID:       kong.String("consumer-2"),
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						ID:   kong.String("plugin-2"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer-2"),
						},
						Route: &kong.Route{
							ID: kong.String("route-1"),
						},
					},
				},
			}),
			expectedPluginID: "plugin-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &UpdateStrategyDBMode{
				logger: logr.Discard(),
			}
			err := s.refillPluginIDs(tc.currentState, tc.targetState)
			require.NoError(t, err)
			_, err = tc.targetState.Plugins.Get(tc.expectedPluginID)
			require.NoError(t, err)
		})
	}
}

func mustNewKongStateFromRawState(t *testing.T, rawState *deckutils.KongRawState) *state.KongState {
	t.Helper()

	kongState, err := state.Get(rawState)
	require.NoError(t, err, "failed to build Kong state from raw state")
	return kongState
}
