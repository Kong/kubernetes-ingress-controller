package license_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/clock"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

type mockKonnectClientClient struct {
	konnectLicense mo.Option[license.KonnectLicense]
	err            error
	getCalls       []time.Time
	lock           sync.RWMutex
	clock          Clock
}

type Clock interface {
	Now() time.Time
}

func newMockKonnectLicenseClient(license mo.Option[license.KonnectLicense], clock Clock) *mockKonnectClientClient {
	return &mockKonnectClientClient{
		konnectLicense: license,
		clock:          clock,
	}
}

func (m *mockKonnectClientClient) Get(context.Context) (mo.Option[license.KonnectLicense], error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.getCalls = append(m.getCalls, m.clock.Now())

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
		}, time.Second, time.Millisecond)
		return matchTime
	}

	t.Run("initial license is retrieved", func(t *testing.T) {
		upstreamClient := newMockKonnectLicenseClient(mo.Some(expectedLicense), clock.System{})
		a := license.NewAgent(upstreamClient, logr.Discard())
		go a.Start(ctx) //nolint:errcheck
		expectLicenseToMatchEventually(t, a, expectedLicense.Payload)
	})

	t.Run("initial license retrieval fails and recovers", func(t *testing.T) {
		ticker := mocks.NewTicker()

		upstreamClient := newMockKonnectLicenseClient(mo.None[license.KonnectLicense](), ticker)

		// Return an error on the first call to List() to verify that the agent handles this correctly.
		upstreamClient.ReturnError(errors.New("something went wrong on a backend"))

		const (
			initialPollingPeriod = time.Minute * 3
			regularPollingPeriod = time.Minute * 20
			allowedDelta         = time.Second
		)

		a := license.NewAgent(
			upstreamClient,
			logr.Discard(),
			license.WithInitialPollingPeriod(initialPollingPeriod),
			license.WithPollingPeriod(regularPollingPeriod),
			license.WithTicker(ticker),
		)

		startTime := time.Now()
		go a.Start(ctx) //nolint:errcheck

		select {
		case <-a.Started():
		case <-time.After(time.Second):
			require.FailNow(t, "timed out waiting for agent to start")
		}

		t.Run("initial polling period is used when no license is retrieved", func(t *testing.T) {
			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) >= 1
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once")

			firstListCallTime := upstreamClient.GetCalls()[0]

			require.WithinDuration(t, startTime, firstListCallTime, allowedDelta,
				"expected first call to List() to happen immediately after starting the agent")

			// Initial polling period has passed...
			ticker.Add(initialPollingPeriod)

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

			// Regular polling period has passed...
			ticker.Add(regularPollingPeriod)

			expectLicenseToMatchEventually(t, a, expectedLicense.Payload)

			listCallsAfterMatchCount := len(upstreamClient.GetCalls())

			// Regular polling period has passed...
			ticker.Add(regularPollingPeriod)

			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) > listCallsAfterMatchCount
			}, time.Second, time.Millisecond, "expected upstream client to be called at least once after the license is retrieved")

			require.Eventually(t, func() bool {
				listCalls := upstreamClient.GetCalls()
				lastListCall := listCalls[len(listCalls)-1]
				lastButOneCall := listCalls[len(listCalls)-2]
				return lastListCall.Sub(lastButOneCall).Abs() <= allowedDelta
			}, time.Second, time.Millisecond)
		})

		t.Run("after the license is retrieved, errors returned from upstream do not override the license", func(t *testing.T) {
			upstreamClient.ReturnError(errors.New("something went wrong on a backend"))

			// Wait for the call to happen.
			initialListCalls := len(upstreamClient.GetCalls())

			// Regular polling period has passed...
			ticker.Add(regularPollingPeriod)

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
			// Regular polling period has passed...
			ticker.Add(regularPollingPeriod)

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

			// Regular polling period has passed...
			ticker.Add(regularPollingPeriod)

			expectLicenseToMatchEventually(t, a, "new-license")
		})
	})

	t.Run("initial license retrieval fails and recovers", func(t *testing.T) {
		ticker := mocks.NewTicker()

		upstreamClient := newMockKonnectLicenseClient(mo.None[license.KonnectLicense](), ticker)

		upstreamClient.ReturnSuccess(mo.Some(license.KonnectLicense{
			Payload:   fmt.Sprintf(testLicense, time.Now().AddDate(0, 0, -10).Format(time.DateOnly)),
			UpdatedAt: expectedLicense.UpdatedAt.Add(-50 * time.Second),
		}))

		const (
			initialPollingPeriod = time.Minute * 3
			regularPollingPeriod = time.Minute * 20
			allowedDelta         = time.Second
		)

		a := license.NewAgent(
			upstreamClient,
			logr.Discard(),
			license.WithInitialPollingPeriod(initialPollingPeriod),
			license.WithPollingPeriod(regularPollingPeriod),
			license.WithTicker(ticker),
		)

		startTime := time.Now()
		go a.Start(ctx) //nolint:errcheck

		select {
		case <-a.Started():
		case <-time.After(time.Second):
			require.FailNow(t, "timed out waiting for agent to start")
		}

		t.Run("initial polling period is used when license is expired", func(t *testing.T) {
			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) >= 1
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least once")

			firstListCallTime := upstreamClient.GetCalls()[0]

			require.WithinDuration(t, startTime, firstListCallTime, allowedDelta,
				"expected first call to List() to happen immediately after starting the agent")

			ticker.Add(initialPollingPeriod)

			require.Eventually(t, func() bool {
				return len(upstreamClient.GetCalls()) >= 2
			}, time.Second, time.Nanosecond, "expected upstream client to be called at least twice")

			secondListCallTime := upstreamClient.GetCalls()[1]
			require.WithinDuration(t, firstListCallTime.Add(initialPollingPeriod), secondListCallTime, allowedDelta,
				"expected second call to List() to happen after the initial polling period as current license is expired")

			require.True(t, a.GetLicense().IsPresent(), "license should be present")
		})

		t.Run("license is always updated when cached is expired", func(t *testing.T) {
			upstreamClient.ReturnSuccess(mo.Some(license.KonnectLicense{
				Payload: "new-license",
				// This is intentionally _older_ than the initial seed license. We want to disregard the normal updated_at
				// rules when our current license is expired.
				UpdatedAt: expectedLicense.UpdatedAt.Add(-100 * time.Second),
			}))

			ticker.Add(regularPollingPeriod)

			expectLicenseToMatchEventually(t, a, "new-license")
		})
	})
}

const testLicense = `{
    "license":{
      "payload":{
        "admin_seats":"1",
        "customer":"TESTTESTTEST",
        "dataplanes":"1",
        "license_creation_date":"2014-11-14",
        "license_expiration_date":"%s",
        "license_key":"TESTTESTTEST",
        "product_subscription":"Test",
        "support_plan":"None"
     },
     "signature":"fake",
     "version":"1"
    }
}`

func TestIsExpiredLicense(t *testing.T) {
	t.Parallel()

	require.True(t, license.IsExpiredLicense(fmt.Sprintf(testLicense, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))))

	require.False(t, license.IsExpiredLicense(fmt.Sprintf(testLicense, time.Now().AddDate(0, 0, 1).Format(time.DateOnly))))

	require.False(t, license.IsExpiredLicense(fmt.Sprintf(testLicense, "not a valid date")))

	require.False(t, license.IsExpiredLicense("this is not valid json"))

	require.False(t, license.IsExpiredLicense(`{"missing_expiration": 0}`))
}
