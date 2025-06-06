package translator

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestTerminatingEndpointStickySessionSupport(t *testing.T) {
	t.Run("terminating endpoints should be included with weight 0", func(t *testing.T) {
		// Test that terminating endpoints are included with weight 0 to support sticky sessions
		endpoints := []util.Endpoint{
			{
				Address:     "10.0.0.1",
				Port:        "80",
				Terminating: false, // Normal endpoint
			},
			{
				Address:     "10.0.0.2",
				Port:        "80",
				Terminating: true, // Terminating endpoint
			},
		}

		targets := targetsForEndpoints(endpoints)
		require.Len(t, targets, 2)

		// Normal endpoint should not have weight set (will use default)
		require.Nil(t, targets[0].Weight)
		require.Equal(t, "10.0.0.1:80", *targets[0].Target.Target)

		// Terminating endpoint should have weight 0
		require.NotNil(t, targets[1].Weight)
		require.Equal(t, 0, *targets[1].Weight)
		require.Equal(t, "10.0.0.2:80", *targets[1].Target.Target)
	})

	t.Run("weight distribution should exclude terminating endpoints", func(t *testing.T) {
		// Test that backend weights are distributed only among non-terminating targets
		targets := []kongstate.Target{
			{
				Target: kong.Target{
					Target: kong.String("10.0.0.1:80"),
					Weight: nil, // Normal target
				},
			},
			{
				Target: kong.Target{
					Target: kong.String("10.0.0.2:80"),
					Weight: lo.ToPtr(0), // Terminating target
				},
			},
			{
				Target: kong.Target{
					Target: kong.String("10.0.0.3:80"),
					Weight: nil, // Normal target
				},
			},
		}

		// Simulate backend weight distribution with weight 100
		backendWeight := 100
		nonTerminatingTargets := 0
		for _, target := range targets {
			if target.Weight == nil || *target.Weight != 0 {
				nonTerminatingTargets++
			}
		}

		targetWeight := backendWeight / nonTerminatingTargets // Should be 50

		for i := range targets {
			// Don't override weight 0 for terminating targets
			if targets[i].Weight != nil && *targets[i].Weight == 0 {
				continue
			}
			targets[i].Weight = &targetWeight
		}

		// Verify results
		require.Equal(t, 50, *targets[0].Weight) // Normal target gets weight 50
		require.Equal(t, 0, *targets[1].Weight)  // Terminating target keeps weight 0
		require.Equal(t, 50, *targets[2].Weight) // Normal target gets weight 50
	})
}

func TestEndpointProcessingWithTerminatingCondition(t *testing.T) {
	t.Run("ready and terminating endpoints should be included", func(t *testing.T) {
		endpointSlice := &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-endpoints",
				Namespace: "default",
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
						Terminating: lo.ToPtr(true),
					},
				},
				{
					Addresses: []string{"10.0.0.3"},
					Conditions: discoveryv1.EndpointConditions{
						Ready:       lo.ToPtr(false),
						Terminating: lo.ToPtr(false),
					},
				},
			},
		}

		var processedEndpoints []util.Endpoint
		for _, endpoint := range endpointSlice.Endpoints {
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

		// Should have 2 endpoints: ready one and terminating one
		require.Len(t, processedEndpoints, 2)

		// First endpoint (ready)
		require.Equal(t, "10.0.0.1", processedEndpoints[0].Address)
		require.False(t, processedEndpoints[0].Terminating)

		// Second endpoint (terminating)
		require.Equal(t, "10.0.0.2", processedEndpoints[1].Address)
		require.True(t, processedEndpoints[1].Terminating)
	})
}

func TestStickySessionsTerminatingEndpointsDrainSupport(t *testing.T) {
	t.Run("drain support disabled - terminating endpoints should be excluded", func(t *testing.T) {
		// Test that when drain support is disabled, terminating endpoints are excluded
		endpointSlices := []*discoveryv1.EndpointSlice{
			{
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
							Terminating: lo.ToPtr(true),
						},
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name:     lo.ToPtr("http"),
						Port:     lo.ToPtr(int32(80)),
						Protocol: lo.ToPtr(corev1.ProtocolTCP),
					},
				},
			},
		}

		mockGetEndpointSlices := func(_, _ string) ([]*discoveryv1.EndpointSlice, error) {
			return endpointSlices, nil
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-service",
				Namespace: "default",
			},
		}

		port := &corev1.ServicePort{
			Name:     "http",
			Port:     80,
			Protocol: corev1.ProtocolTCP,
		}

		// Test with drain support disabled (false)
		processedEndpoints := getEndpoints(
			zapr.NewLogger(zap.NewNop()),
			service,
			port,
			corev1.ProtocolTCP,
			mockGetEndpointSlices,
			false, // isSvcUpstream
			"cluster.local",
			false, // enableStickySessionsTerminatingEndpoints = false
		)

		// Should only have the ready endpoint, not the terminating one
		require.Len(t, processedEndpoints, 1)
		require.Equal(t, "10.0.0.1", processedEndpoints[0].Address)
		require.False(t, processedEndpoints[0].Terminating)
	})

	t.Run("drain support enabled - terminating endpoints should be included", func(t *testing.T) {
		// Test that when drain support is enabled, terminating endpoints are included
		endpointSlices := []*discoveryv1.EndpointSlice{
			{
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
							Terminating: lo.ToPtr(true),
						},
					},
				},
				Ports: []discoveryv1.EndpointPort{
					{
						Name:     lo.ToPtr("http"),
						Port:     lo.ToPtr(int32(80)),
						Protocol: lo.ToPtr(corev1.ProtocolTCP),
					},
				},
			},
		}

		mockGetEndpointSlices := func(_, _ string) ([]*discoveryv1.EndpointSlice, error) {
			return endpointSlices, nil
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-service",
				Namespace: "default",
			},
		}

		port := &corev1.ServicePort{
			Name:     "http",
			Port:     80,
			Protocol: corev1.ProtocolTCP,
		}

		// Test with drain support enabled (true)
		processedEndpoints := getEndpoints(
			zapr.NewLogger(zap.NewNop()),
			service,
			port,
			corev1.ProtocolTCP,
			mockGetEndpointSlices,
			false, // isSvcUpstream
			"cluster.local",
			true, // enableStickySessionsTerminatingEndpoints = true
		)

		// Should have both endpoints
		require.Len(t, processedEndpoints, 2)

		// First endpoint (ready)
		require.Equal(t, "10.0.0.1", processedEndpoints[0].Address)
		require.False(t, processedEndpoints[0].Terminating)

		// Second endpoint (terminating, should be marked as terminating)
		require.Equal(t, "10.0.0.2", processedEndpoints[1].Address)
		require.True(t, processedEndpoints[1].Terminating)
	})
}
