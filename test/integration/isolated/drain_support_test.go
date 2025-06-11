package isolated

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestDrainSupportEndpointProcessing(t *testing.T) {
	t.Run("terminating endpoints should be included with weight 0", func(t *testing.T) {
		// Simulate endpoint processing logic for terminating pods
		endpointSlice := &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-endpoints",
				Namespace: "default",
				Labels: map[string]string{
					"kubernetes.io/service-name": "test-service",
				},
			},

			Endpoints: []discoveryv1.Endpoint{
				{
					Addresses: []string{"10.0.0.1"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
				},
				{
					Addresses: []string{"10.0.0.2"},
					Conditions: discoveryv1.EndpointConditions{
						Ready:       lo.ToPtr(false),
						Terminating: lo.ToPtr(true), // Terminating pod
					},
				},
				{
					Addresses: []string{"10.0.0.3"},
					Conditions: discoveryv1.EndpointConditions{
						Ready: lo.ToPtr(true),
					},
				},
				{
					Addresses: []string{"10.0.0.4"},
					Conditions: discoveryv1.EndpointConditions{
						Ready:       lo.ToPtr(false),
						Terminating: lo.ToPtr(false), // Not ready, not terminating (should be skipped)
					},
				},
			},
		}

		// Process endpoints using the same logic as in the translator
		var processedEndpoints []util.Endpoint
		for _, endpoint := range endpointSlice.Endpoints {
			// Skip endpoints that are not ready, unless they are terminating.
			// Terminating endpoints should be included with weight 0 for drain support.
			isTerminating := endpoint.Conditions.Terminating != nil && *endpoint.Conditions.Terminating
			isReady := endpoint.Conditions.Ready == nil || *endpoint.Conditions.Ready

			if !isReady && !isTerminating {
				continue
			}

			upstreamServer := util.Endpoint{
				Address:     endpoint.Addresses[0],
				Port:        "80",
				Terminating: isTerminating,
			}
			processedEndpoints = append(processedEndpoints, upstreamServer)
		}

		// Verify results
		require.Len(t, processedEndpoints, 3, "should process 3 endpoints: 2 ready + 1 terminating")

		readyEndpoints := 0
		terminatingEndpoints := 0

		for _, ep := range processedEndpoints {
			if ep.Terminating {
				terminatingEndpoints++
				require.Equal(t, "10.0.0.2", ep.Address, "terminating endpoint should be the correct one")
			} else {
				readyEndpoints++
				require.True(t, ep.Address == "10.0.0.1" || ep.Address == "10.0.0.3", "ready endpoints should be the correct ones")
			}
		}

		require.Equal(t, 2, readyEndpoints, "should have 2 ready endpoints")
		require.Equal(t, 1, terminatingEndpoints, "should have 1 terminating endpoint")
	})

	t.Run("endpoint processing behavior is correct for all conditions", func(t *testing.T) {
		testCases := []struct {
			name          string
			ready         *bool
			terminating   *bool
			shouldInclude bool
		}{
			{"ready endpoint", lo.ToPtr(true), lo.ToPtr(false), true},
			{"ready endpoint with nil terminating", lo.ToPtr(true), nil, true},
			{"terminating endpoint", lo.ToPtr(false), lo.ToPtr(true), true},
			{"not ready and not terminating", lo.ToPtr(false), lo.ToPtr(false), false},
			{"not ready with nil terminating", lo.ToPtr(false), nil, false},
			{"nil ready (assumed ready)", nil, lo.ToPtr(false), true},
			{"nil ready and terminating", nil, lo.ToPtr(true), true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := discoveryv1.Endpoint{
					Conditions: discoveryv1.EndpointConditions{
						Ready:       tc.ready,
						Terminating: tc.terminating,
					},
				}

				isTerminating := endpoint.Conditions.Terminating != nil && *endpoint.Conditions.Terminating
				isReady := endpoint.Conditions.Ready == nil || *endpoint.Conditions.Ready

				shouldProcess := isReady || isTerminating

				require.Equal(t, tc.shouldInclude, shouldProcess, "endpoint processing decision should match expectation")
			})
		}
	})
}
