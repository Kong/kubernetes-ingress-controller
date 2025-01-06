package manager_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestAdminAPIClientFromServiceDiscovery(t *testing.T) {
	log := logr.Discard()
	adminAPISvcNN := k8stypes.NamespacedName{Name: "admin-api", Namespace: "kong"}
	kubeClient := fake.NewFakeClient()
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
		cancelContext  bool             // If true, will cancel the context after GetAdminAPIsForService is called 2 times.

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
			name: "two apis one error",
			discoveredAPIs: sets.New(
				someDiscoveredAPI("https://localhost:8444"),
				someDiscoveredAPI("https://localhost:8445"),
			),
			factoryErrs: map[string]error{
				"https://localhost:8445": genericErr,
			},
			expectedErr: genericErr,
		},
		{
			name:          "no admin apis waits forever",
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
			discoverer := mocks.NewAdminAPIDiscoverer(tc.discoveredAPIs, tc.discovererErr)
			factory := mocks.NewAdminAPIClientFactory(tc.factoryErrs)

			// If cancelContext is true, we will cancel the context after GetAdminAPIsForService is called >= 2 times.
			// This will mean that the retry loop is running. By cancelling the context we can ensure it will exit and
			// return the expected error.
			if tc.cancelContext {
				go func() {
					if assert.Eventually(t, func() bool {
						return discoverer.GetAdminAPIsForServiceCalledTimes() >= 2
					}, time.Second, time.Millisecond) {
						t.Log("cancelling context, GetAdminAPIsForService called >= 2 times")
						cancel()
					}
				}()
			} else {
				defer cancel()
			}

			retryEveryMs := retry.Delay(time.Millisecond) // For testing purposes, we want to retry as fast as possible.
			clients, err := manager.AdminAPIClientFromServiceDiscovery(ctx, log, adminAPISvcNN, kubeClient, discoverer, factory, retryEveryMs)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Len(t, clients, tc.expectedClientsCount)
		})
	}
}
