package clients_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
)

const (
	testURL1 = "http://localhost:8001"
	testURL2 = "http://localhost:8002"
)

var testPodRef = k8stypes.NamespacedName{
	Namespace: "default",
	Name:      "mock",
}

type mockClientFactory struct {
	ready      map[string]bool // Maps address to readiness.
	callsCount map[string]int  // Maps address to number of CreateAdminAPIClient calls.
	t          *testing.T
}

func newMockClientFactory(t *testing.T, ready map[string]bool) mockClientFactory {
	return mockClientFactory{
		ready:      ready,
		callsCount: map[string]int{},
		t:          t,
	}
}

func (cf mockClientFactory) CreateAdminAPIClient(_ context.Context, adminAPI adminapi.DiscoveredAdminAPI) (*adminapi.Client, error) {
	address := adminAPI.Address

	cf.callsCount[address]++

	ready, ok := cf.ready[address]
	if !ok {
		cf.t.Errorf("unexpected client creation for %s", address)
	}
	if !ok || !ready {
		return nil, fmt.Errorf("client for %s is not ready", address)
	}

	return adminapi.NewTestClient(address)
}

type mockAlreadyCreatedClient struct {
	url     string
	isReady bool
}

func (m mockAlreadyCreatedClient) IsReady(context.Context) error {
	if !m.isReady {
		return errors.New("not ready")
	}
	return nil
}

func (m mockAlreadyCreatedClient) PodReference() (k8stypes.NamespacedName, bool) {
	return testPodRef, true
}

func (m mockAlreadyCreatedClient) BaseRootURL() string {
	return m.url
}

func TestDefaultReadinessChecker(t *testing.T) {
	testCases := []struct {
		name string

		alreadyCreatedClients   []clients.AlreadyCreatedClient
		pendingClients          []adminapi.DiscoveredAdminAPI
		pendingClientsReadiness map[string]bool

		expectedTurnedReady   []string
		expectedTurnedPending []string
	}{
		{
			name: "ready turning pending",
			alreadyCreatedClients: []clients.AlreadyCreatedClient{
				mockAlreadyCreatedClient{
					url:     testURL1,
					isReady: false,
				},
			},
			expectedTurnedPending: []string{testURL1},
		},
		{
			name: "pending turning ready",
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: testURL1,
					PodRef:  testPodRef,
				},
			},
			pendingClientsReadiness: map[string]bool{
				testURL1: true,
			},
			expectedTurnedReady: []string{testURL1},
		},
		{
			name: "ready turning pending, pending turning ready at once",
			alreadyCreatedClients: []clients.AlreadyCreatedClient{
				mockAlreadyCreatedClient{
					url:     testURL1,
					isReady: false,
				},
			},
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: testURL2,
					PodRef:  testPodRef,
				},
			},
			pendingClientsReadiness: map[string]bool{
				testURL2: true,
			},
			expectedTurnedReady:   []string{testURL2},
			expectedTurnedPending: []string{testURL1},
		},
		{
			name: "no changes",
			alreadyCreatedClients: []clients.AlreadyCreatedClient{
				mockAlreadyCreatedClient{
					url:     testURL1,
					isReady: true,
				},
			},
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: testURL2,
					PodRef:  testPodRef,
				},
			},
			pendingClientsReadiness: map[string]bool{
				testURL2: false,
			},
			expectedTurnedReady:   nil,
			expectedTurnedPending: nil,
		},
		{
			name:                  "no clients at all",
			expectedTurnedReady:   nil,
			expectedTurnedPending: nil,
		},
		{
			name: "multiple ready, one turning pending",
			alreadyCreatedClients: []clients.AlreadyCreatedClient{
				mockAlreadyCreatedClient{
					url:     testURL1,
					isReady: true,
				},
				mockAlreadyCreatedClient{
					url:     testURL2,
					isReady: false, // This one will turn pending.
				},
			},
			expectedTurnedReady: nil,
			expectedTurnedPending: []string{
				testURL2,
			},
		},
		{
			name: "multiple pending, one turning ready",
			pendingClients: []adminapi.DiscoveredAdminAPI{
				{
					Address: testURL1,
					PodRef:  testPodRef,
				},
				{
					Address: testURL2,
					PodRef:  testPodRef,
				},
			},
			pendingClientsReadiness: map[string]bool{
				testURL1: false,
				testURL2: true, // This one will turn ready.
			},
			expectedTurnedReady: []string{
				testURL2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			factory := newMockClientFactory(t, tc.pendingClientsReadiness)
			checker := clients.NewDefaultReadinessChecker(factory, logr.Discard())
			result := checker.CheckReadiness(context.Background(), tc.alreadyCreatedClients, tc.pendingClients)

			turnedPending := lo.Map(result.ClientsTurnedPending, func(c adminapi.DiscoveredAdminAPI, _ int) string { return c.Address })
			turnedReady := lo.Map(result.ClientsTurnedReady, func(c *adminapi.Client, _ int) string { return c.BaseRootURL() })

			require.ElementsMatch(t, tc.expectedTurnedReady, turnedReady)
			require.ElementsMatch(t, tc.expectedTurnedPending, turnedPending)

			// For every pending client turning ready we expect exactly one call to CreateAdminAPIClient.
			for _, url := range tc.pendingClients {
				require.Equal(t, 1, factory.callsCount[url.Address])
			}

			// For every already created client we expect NO calls to CreateAdminAPIClient.
			for _, url := range tc.alreadyCreatedClients {
				require.Zero(t, factory.callsCount[url.BaseRootURL()])
			}
		})
	}
}
