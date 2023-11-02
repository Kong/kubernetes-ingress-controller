package adminapi_test

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

type mockClock struct {
	n time.Time
}

func newMockClock() *mockClock {
	return &mockClock{n: time.Now()}
}

func (m *mockClock) Now() time.Time {
	return m.n
}

func (m *mockClock) MoveBy(d time.Duration) {
	m.n = m.n.Add(d)
}

func TestKonnectBackoffStrategy(t *testing.T) {
	someTestHash := func(s string) []byte {
		h := sha256.Sum256([]byte(s))
		return h[:]
	}

	var (
		clock   = newMockClock()
		hashOne = someTestHash("1")
		hashTwo = someTestHash("2")
	)

	t.Run("on init allows any updates", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)

		// First try, no failure in the past, should always allow update, register failure.
		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow any update as a first one")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow any update as a first one")
		assert.Empty(t, whyNot)
	})

	t.Run("generic failure triggers backoff time requirement", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)

		// After failure, time moves by 1s (below initial backoff time), should not allow update.
		strategy.RegisterUpdateFailure(errors.New("error occurred"), hashOne)
		clock.MoveBy(time.Second)

		canUpdate, whyNot := strategy.CanUpdate(hashTwo)
		assert.False(t, canUpdate, "should not allow next update when last failed and backoff time wasn't satisfied")
		assert.Equal(t, "next attempt allowed in 2s", whyNot)

		// Time moves by 5s (enough for the backoff next try).
		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for different hash when enough time has passed")
		assert.Empty(t, whyNot)
	})

	t.Run("client error triggers faulty hash requirements", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusBadRequest, ""), hashOne)
		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "should not allow update for the same faulty hash")
		assert.Equal(t, "Config has to be changed: \"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b\" hash has already failed to be pushed with a client error", whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for another hash")
		assert.Empty(t, whyNot)
	})

	t.Run("success resets both hash and backoff requirements", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusBadRequest, ""), hashOne)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "should not allow next update when last failed and backoff time wasn't satisfied")
		assert.Equal(t, "Config has to be changed: \"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b\" hash has already failed to be pushed with a client error, next attempt allowed in 3s", whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.False(t, canUpdate, "should not allow next update when last failed and backoff time wasn't satisfied")
		assert.Equal(t, "next attempt allowed in 3s", whyNot)

		strategy.RegisterUpdateSuccess()

		canUpdate, whyNot = strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow any update after the last success")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow any update after the last success")
		assert.Empty(t, whyNot)
	})

	t.Run("server error does not trigger faulty hash requirement", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusInternalServerError, ""), hashOne)
		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow update for the same faulty hash as it was a server error")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for another hash")
		assert.Empty(t, whyNot)
	})

	t.Run("too many requests code with no details embedded", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)
		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusTooManyRequests, ""), hashOne)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "shouldn't allow update due to a standard backoff time")
		assert.Equal(t, "next attempt allowed in 3s", whyNot)

		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot = strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow update for the same hash after standard backoff time")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for another hash")
		assert.Empty(t, whyNot)
	})

	t.Run("too many requests code with details embedded", func(t *testing.T) {
		strategy := adminapi.NewKonnectBackoffStrategy(clock)
		tooManyRequestsAPIErr := kong.NewAPIError(http.StatusTooManyRequests, "")
		const retryAfter = time.Minute
		tooManyRequestsAPIErr.SetDetails(kong.ErrTooManyRequestsDetails{
			RetryAfter: retryAfter,
		})
		strategy.RegisterUpdateFailure(tooManyRequestsAPIErr, hashOne)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "shouldn't allow update due to the suggested retry-after backoff")
		assert.Equal(t, "next attempt allowed in 1m0s", whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.False(t, canUpdate, "shouldn't allow update due to the suggested retry-after backoff (different hash)")
		assert.Equal(t, "next attempt allowed in 1m0s", whyNot)

		clock.MoveBy(time.Minute + time.Second)
		canUpdate, whyNot = strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow update after the suggested retry-after backoff")
		assert.Empty(t, whyNot)
	})
}
