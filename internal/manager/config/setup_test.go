package config_test

import (
	"context"
	errors "errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type mockAdminAPIDiscoverer struct {
	apis sets.Set[adminapi.DiscoveredAdminAPI]
	err  error
}

func (m *mockAdminAPIDiscoverer) GetAdminAPIsForService(_ context.Context, _ client.Client, _ k8stypes.NamespacedName) (sets.Set[adminapi.DiscoveredAdminAPI], error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.apis, nil
}

type mockAdminAPIClientFactory struct {
	errs map[string]error // map from address to error
}

func newMockAdminAPIClientFactory(errs map[string]error) *mockAdminAPIClientFactory {
	if errs == nil {
		errs = make(map[string]error)
	}
	return &mockAdminAPIClientFactory{
		errs: errs,
	}
}

func (m *mockAdminAPIClientFactory) CreateAdminAPIClient(_ context.Context, api adminapi.DiscoveredAdminAPI) (*adminapi.Client, error) {
	err, ok := m.errs[api.Address]
	if !ok {
		return adminapi.NewTestClient(api.Address)
	}
	return nil, err
}

func TestAdminAPIClientFromServiceDiscovery(t *testing.T) {
	log := logr.Discard()
	adminAPISvcNN := k8stypes.NamespacedName{Name: "admin-api", Namespace: "kong"}
	kubeClient := fake.NewClientBuilder().Build()
	genericErr := errors.New("some generic error")
	someDiscoveredAPI := func(address string) adminapi.DiscoveredAdminAPI {
		return adminapi.DiscoveredAdminAPI{
			Address: address,
			PodRef: k8stypes.NamespacedName{
				Namespace: "kong",
				Name:      "pod",
			},
		}
	}
	testCases := []struct {
		name           string
		discoveredAPIs sets.Set[adminapi.DiscoveredAdminAPI]
		discovererErr  error
		factoryErrs    map[string]error // Map from address to error.
		cancelContext  bool             // If true, cancel the context after 100ms to not wait forever.

		expectedClientsCount int
		expectedErr          error
	}{
		{
			name: "no errors and one admin api",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
			),
			expectedClientsCount: 1,
		},
		{
			name: "no errors and two apis",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
				someDiscoveredAPI("https://localhost:8445"),
			),
			expectedClientsCount: 2,
		},
		{
			name: "two apis but one is not ready",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
				someDiscoveredAPI("https://localhost:8445"),
			),
			factoryErrs: map[string]error{
				"https://localhost:8445": adminapi.KongClientNotReadyError{},
			},
			expectedClientsCount: 1,
		},
		{
			name: "two apis both not ready waits forever",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
				someDiscoveredAPI("https://localhost:8445"),
			),
			factoryErrs: map[string]error{
				"https://localhost:8444": adminapi.KongClientNotReadyError{},
				"https://localhost:8445": adminapi.KongClientNotReadyError{},
			},
			cancelContext: true,
			expectedErr:   context.Canceled,
		},
		{
			name: "two apis one not ready one generic error aborts",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
				someDiscoveredAPI("https://localhost:8445"),
			),
			factoryErrs: map[string]error{
				"https://localhost:8444": adminapi.KongClientNotReadyError{},
				"https://localhost:8445": genericErr,
			},
			expectedErr: genericErr,
		},
		{
			name:          "no admin apis with no errors waits forever",
			cancelContext: true,
			expectedErr:   context.Canceled,
		},
		{
			name:          "any discoverer error aborts",
			discovererErr: genericErr,
			expectedErr:   genericErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			if tc.cancelContext {
				go func() {
					// We assume that 100ms is enough for the test to finish.
					<-time.After(100 * time.Millisecond)
					cancel()
				}()
			} else {
				defer cancel()
			}

			discoverer := &mockAdminAPIDiscoverer{apis: tc.discoveredAPIs, err: tc.discovererErr}
			factory := newMockAdminAPIClientFactory(tc.factoryErrs)

			clients, err := manager.AdminAPIClientFromServiceDiscovery(ctx, log, adminAPISvcNN, kubeClient, discoverer, factory)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Len(t, clients, tc.expectedClientsCount)
		})
	}

}
