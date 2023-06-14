package license_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	konnectLicense "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/license"
)

type mockKonnectClientClient struct {
	listResponse *konnectLicense.ListLicenseResponse
	err          error
	listCalls    []time.Time
	lock         sync.RWMutex
}

func newMockKonnectLicenseClient(listResponse *konnectLicense.ListLicenseResponse) *mockKonnectClientClient {
	return &mockKonnectClientClient{listResponse: listResponse}
}

func (m *mockKonnectClientClient) List(context.Context, int) (*konnectLicense.ListLicenseResponse, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.listCalls = append(m.listCalls, time.Now())

	if m.err != nil {
		return nil, m.err
	}
	return m.listResponse, nil
}

func (m *mockKonnectClientClient) ReturnError(err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.err = err
}

func (m *mockKonnectClientClient) ReturnSuccess(listResponse *konnectLicense.ListLicenseResponse) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.listResponse = listResponse
	m.err = nil
}

func (m *mockKonnectClientClient) ListCalls() []time.Time {
	m.lock.RLock()
	defer m.lock.RUnlock()

	copied := make([]time.Time, len(m.listCalls))
	copy(copied, m.listCalls)
	return copied
}

func TestAgent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedLicense := &konnectLicense.Item{
		License:   "test-license",
		UpdatedAt: 1234567890,
	}
	expectedListResponse := &konnectLicense.ListLicenseResponse{
		Items: []*konnectLicense.Item{
			expectedLicense,
		},
	}

	expectLicenseToMatchEventually := func(t *testing.T, a *license.Agent, expectedPayload string) time.Time {
		var matchTime time.Time
		require.Eventually(t, func() bool {
			actualLicense, ok := a.GetLicense()
			if !ok {
				t.Log("license not yet available")
				return false
			}
			if *actualLicense.Payload != expectedPayload {
				t.Logf("license mismatch: expected %q, got %q", expectedPayload, *actualLicense.Payload)
				return false
			}
			matchTime = time.Now()
			return true
		}, time.Second, time.Nanosecond)
		return matchTime
	}

	t.Run("initial license is retrieved", func(t *testing.T) {
		upstreamClient := newMockKonnectLicenseClient(expectedListResponse)
		a := license.NewAgent(upstreamClient, logr.Discard())
		go func() {
			err := a.Start(ctx)
			require.NoError(t, err)
		}()
		expectLicenseToMatchEventually(t, a, expectedLicense.License)
	})

	t.Run("initial license retrieval fails and recovers", func(t *testing.T) {
		upstreamClient := newMockKonnectLicenseClient(nil)

		// Return an error on the first call to List() to verify that the agent handles this correctly.
		upstreamClient.ReturnError(errors.New("something went wrong on a backend"))

		const (
			// Set the initial polling period to a very short duration to ensure that the agent retries quickly.
			initialPollingPeriod = time.Millisecond
			regularPollingPeriod = time.Millisecond * 5
			allowedDelta         = time.Millisecond
		)
		a := license.NewAgent(
			upstreamClient,
			logr.Discard(),
			license.WithInitialPollingPeriod(initialPollingPeriod),
			license.WithPollingPeriod(regularPollingPeriod),
		)

		startTime := time.Now()
		go func() {
			err := a.Start(ctx)
			require.NoError(t, err)
		}()

		t.Run("initial polling period is used when no license is retrieved", func(t *testing.T) {
			require.Eventually(t, func() bool {
				return len(upstreamClient.ListCalls()) >= 1
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once")

			firstListCallTime := upstreamClient.ListCalls()[0]
			require.WithinDuration(t, startTime.Add(initialPollingPeriod), firstListCallTime, allowedDelta,
				"expected first call to List() to happen after the initial polling period")

			require.Eventually(t, func() bool {
				return len(upstreamClient.ListCalls()) >= 2
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least twice")

			secondListCallTime := upstreamClient.ListCalls()[1]
			require.WithinDuration(t, firstListCallTime.Add(initialPollingPeriod), secondListCallTime, allowedDelta,
				"expected second call to List() to happen after the initial polling period as no license is retrieved yet")

			_, ok := a.GetLicense()
			require.False(t, ok, "no license should be available due to an error in the upstream client")
		})

		t.Run("regular polling period is used after the initial license is retrieved", func(t *testing.T) {
			// Now return a valid response to ensure that the agent recovers.
			upstreamClient.ReturnSuccess(expectedListResponse)
			expectLicenseToMatchEventually(t, a, expectedLicense.License)

			listCallsAfterMatchCount := len(upstreamClient.ListCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.ListCalls()) > listCallsAfterMatchCount
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once after the license is retrieved")

			listCalls := upstreamClient.ListCalls()
			lastListCall := listCalls[len(listCalls)-1]
			lastButOneCall := listCalls[len(listCalls)-2]
			require.WithinDuration(t, lastButOneCall.Add(regularPollingPeriod), lastListCall, allowedDelta)
		})

		t.Run("after the license is retrieved, errors returned from upstream do not override the license", func(t *testing.T) {
			upstreamClient.ReturnError(errors.New("something went wrong on a backend"))

			// Wait for the call to happen.
			initialListCalls := len(upstreamClient.ListCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.ListCalls()) > initialListCalls
			}, time.Second, time.Nanosecond)

			// The license should still be available.
			_, ok := a.GetLicense()
			require.True(t, ok, "license should be available even if the upstream client returns an error")
		})

		t.Run("license is not updated when the upstream returns a license updated before the cached one", func(t *testing.T) {
			upstreamClient.ReturnSuccess(&konnectLicense.ListLicenseResponse{
				Items: []*konnectLicense.Item{
					{
						License:   "new-license",
						UpdatedAt: expectedLicense.UpdatedAt - 1,
					},
				},
			})

			// Wait for the call to happen.
			initialListCalls := len(upstreamClient.ListCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.ListCalls()) > initialListCalls
			}, time.Second, time.Nanosecond)

			// The cached license should still be available.
			expectLicenseToMatchEventually(t, a, expectedLicense.License)
		})

		t.Run("license is updated when the upstream returns a license updated after the cached one", func(t *testing.T) {
			upstreamClient.ReturnSuccess(&konnectLicense.ListLicenseResponse{
				Items: []*konnectLicense.Item{
					{
						License:   "new-license",
						UpdatedAt: expectedLicense.UpdatedAt + 1,
					},
				},
			})
			expectLicenseToMatchEventually(t, a, "new-license")
		})
	})
}
