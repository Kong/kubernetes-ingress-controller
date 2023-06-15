package license_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/license"
)

type mockKonnectClientClient struct {
	konnectLicense mo.Option[license.KonnectLicense]
	err            error
	getCalls       []time.Time
	lock           sync.RWMutex
}

func newMockKonnectLicenseClient(license mo.Option[license.KonnectLicense]) *mockKonnectClientClient {
	return &mockKonnectClientClient{konnectLicense: license}
}

func (m *mockKonnectClientClient) Get(context.Context) (mo.Option[license.KonnectLicense], error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.getCalls = append(m.getCalls, time.Now())

	if m.err != nil {
		return mo.None[license.KonnectLicense](), m.err
	}
	return m.konnectLicense, nil
}

func (m *mockKonnectClientClient) ReturnError(err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.err = err
}

func (m *mockKonnectClientClient) ReturnSuccess(license mo.Option[license.KonnectLicense]) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.konnectLicense = license
	m.err = nil
}

func (m *mockKonnectClientClient) GetCalls() []time.Time {
	m.lock.RLock()
	defer m.lock.RUnlock()

	copied := make([]time.Time, len(m.getCalls))
	copy(copied, m.getCalls)
	return copied
}

func TestAgent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	expectedLicense := license.KonnectLicense{
		Payload:   "test-license",
		UpdatedAt: time.Now(),
	}

	expectLicenseToMatchEventually := func(t *testing.T, a *license.Agent, expectedPayload string) time.Time {
		var matchTime time.Time
		require.Eventually(t, func() bool {
			actualLicense, ok := a.GetLicense().Get()
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
		upstreamClient := newMockKonnectLicenseClient(mo.Some(expectedLicense))
		a := license.NewAgent(upstreamClient, logr.Discard())
		go func() {
			err := a.Start(ctx)
			require.NoError(t, err)
		}()
		expectLicenseToMatchEventually(t, a, expectedLicense.Payload)
	})

	t.Run("initial license retrieval fails and recovers", func(t *testing.T) {
		upstreamClient := newMockKonnectLicenseClient(mo.None[license.KonnectLicense]())

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
				return len(upstreamClient.GetCalls()) >= 1
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once")

			firstListCallTime := upstreamClient.GetCalls()[0]
			require.WithinDuration(t, startTime.Add(initialPollingPeriod), firstListCallTime, allowedDelta,
				"expected first call to List() to happen after the initial polling period")

			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) >= 2
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least twice")

			secondListCallTime := upstreamClient.GetCalls()[1]
			require.WithinDuration(t, firstListCallTime.Add(initialPollingPeriod), secondListCallTime, allowedDelta,
				"expected second call to List() to happen after the initial polling period as no license is retrieved yet")

			require.False(t, a.GetLicense().IsPresent(), "no license should be available due to an error in the upstream client")
		})

		t.Run("regular polling period is used after the initial license is retrieved", func(t *testing.T) {
			// Now return a valid response to ensure that the agent recovers.
			upstreamClient.ReturnSuccess(mo.Some(expectedLicense))
			expectLicenseToMatchEventually(t, a, expectedLicense.Payload)

			listCallsAfterMatchCount := len(upstreamClient.GetCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) > listCallsAfterMatchCount
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once after the license is retrieved")

			listCalls := upstreamClient.GetCalls()
			lastListCall := listCalls[len(listCalls)-1]
			lastButOneCall := listCalls[len(listCalls)-2]
			require.WithinDuration(t, lastButOneCall.Add(regularPollingPeriod), lastListCall, allowedDelta)
		})

		t.Run("after the license is retrieved, errors returned from upstream do not override the license", func(t *testing.T) {
			upstreamClient.ReturnError(errors.New("something went wrong on a backend"))

			// Wait for the call to happen.
			initialListCalls := len(upstreamClient.GetCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) > initialListCalls
			}, time.Second, time.Nanosecond)

			// The license should still be available.
			require.True(t, a.GetLicense().IsPresent(), "license should be available even if the upstream client returns an error")
		})

		t.Run("license is not updated when the upstream returns a license updated before the cached one", func(t *testing.T) {
			upstreamClient.ReturnSuccess(mo.Some(license.KonnectLicense{
				Payload:   "new-license",
				UpdatedAt: expectedLicense.UpdatedAt.Add(-time.Second),
			}))

			// Wait for the call to happen.
			initialListCalls := len(upstreamClient.GetCalls())
			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) > initialListCalls
			}, time.Second, time.Nanosecond)

			// The cached license should still be available.
			expectLicenseToMatchEventually(t, a, expectedLicense.Payload)
		})

		t.Run("license is updated when the upstream returns a license updated after the cached one", func(t *testing.T) {
			upstreamClient.ReturnSuccess(mo.Some(license.KonnectLicense{
				Payload:   "new-license",
				UpdatedAt: expectedLicense.UpdatedAt.Add(time.Second),
			}))
			expectLicenseToMatchEventually(t, a, "new-license")
		})
	})
}
