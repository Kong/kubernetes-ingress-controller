package kongstate_test

import (
	"reflect"
	"testing"

	"github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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
						Name:      kong.String("service"),
						ID:        kong.String("service"),
						CreatedAt: kong.Int(100),
					},
				},
				Routes: []*kong.Route{
					{
						Name:      kong.String("route"),
						ID:        kong.String("route"),
						CreatedAt: kong.Int(101),
						Service: &kong.Service{
							ID: kong.String("service"),
						},
					},
				},
				Upstreams: []*kong.Upstream{
					{
						Name: kong.String("upstream"),
						ID:   kong.String("upstream"),
					},
				},
				Targets: []*kong.Target{
					{
						ID:        kong.String("target"),
						CreatedAt: kong.Float64(102),
						Weight:    kong.Int(999),
						Upstream: &kong.Upstream{
							ID: kong.String("upstream"),
						},
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("plugin1"),
						ID:   kong.String("plugin1"),
						Service: &kong.Service{
							ID: kong.String("service"),
						},
					},
					{
						Name: kong.String("plugin2"),
						ID:   kong.String("plugin2"),
						Route: &kong.Route{
							ID: kong.String("route"),
						},
					},
				},
			},
			expectedKongState: &kongstate.KongState{
				Services: []kongstate.Service{
					{
						Service: kong.Service{
							Name: kong.String("service"),
						},
						Plugins: []kong.Plugin{
							{
								Name: kong.String("plugin1"),
							},
						},
						Routes: []kongstate.Route{
							{
								Route: kong.Route{
									Name: kong.String("route"),
								},
								Plugins: []kong.Plugin{
									{
										Name: kong.String("plugin2"),
									},
								},
							},
						},
					},
				},
				Upstreams: []kongstate.Upstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("upstream"),
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

func TestKongStateToKongRawState_Ensure(t *testing.T) {
	kongRawStateFieldsKICDoesntSupport := []string{
		// These are fields that KIC explicitly doesn't support.
		"Vault",
	}
	allKongRawStateFields := func() []string {
		var fields []string
		typ := reflect.ValueOf(utils.KongRawState{}).Type()
		for i := 0; i < typ.NumField(); i++ {
			fields = append(fields, typ.Field(i).Name)
		}
		return fields
	}()

	testCases := []struct {
		testedFields []string
	}{
		{
			testedFields: []string{
				"Services",
				"Routes",
				"Upstreams",
				"Targets",
				"Plugins",
				"Certificates",
				"CACertificates",
			},
		},
	}

	// Kinda meta test - ensure we have testcases covering all fields in KongRawState.
	for _, field := range allKongRawStateFields {
		if lo.Contains(kongRawStateFieldsKICDoesntSupport, field) {
			t.Logf("skipping field %s - unsupported explicitly", field)
			continue
		}
		testCoveringFieldExists := lo.ContainsBy(testCases, func(tc struct{ testedFields []string }) bool {
			return lo.Contains(tc.testedFields, field)
		})
		assert.True(t, testCoveringFieldExists, "no test covering field %s", field)
	}

	// Run the tests.
	for _, tc := range testCases {
		// ...
		_ = tc
	}
}
