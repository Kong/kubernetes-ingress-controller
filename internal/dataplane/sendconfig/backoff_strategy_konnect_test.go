package sendconfig_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
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
	var (
		clock   = newMockClock()
		hashOne = []byte("1")
		hashTwo = []byte("2")
	)

	t.Run("on init allows any updates", func(t *testing.T) {
		strategy := sendconfig.NewKonnectBackoffStrategy(clock)

		// First try, no failure in the past, should always allow update, register failure.
		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow any update as a first one")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow any update as a first one")
		assert.Empty(t, whyNot)
	})

	t.Run("generic failure triggers backoff time requirement", func(t *testing.T) {
		strategy := sendconfig.NewKonnectBackoffStrategy(clock)

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
	})

	t.Run("client error triggers faulty hash requirements", func(t *testing.T) {
		strategy := sendconfig.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusBadRequest, ""), hashOne)
		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "should not allow update for the same faulty hash")
		assert.Equal(t, "configuration with \"1\" hash has already been attempted to be pushed and it resulted in a client error", whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for another hash")
		assert.Empty(t, whyNot)
	})

	t.Run("success resets both hash and backoff requirements", func(t *testing.T) {
		strategy := sendconfig.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusBadRequest, ""), hashOne)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.False(t, canUpdate, "should not allow next update when last failed and backoff time wasn't satisfied")
		assert.Equal(t, "configuration with \"1\" hash has already been attempted to be pushed and it resulted in a client error, next attempt allowed in 3s", whyNot)

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
		strategy := sendconfig.NewKonnectBackoffStrategy(clock)

		strategy.RegisterUpdateFailure(kong.NewAPIError(http.StatusInternalServerError, ""), hashOne)
		clock.MoveBy(time.Second * 5)

		canUpdate, whyNot := strategy.CanUpdate(hashOne)
		assert.True(t, canUpdate, "should allow update for the same faulty hash as it was a server error")
		assert.Empty(t, whyNot)

		canUpdate, whyNot = strategy.CanUpdate(hashTwo)
		assert.True(t, canUpdate, "should allow update for another hash")
		assert.Empty(t, whyNot)
	})
}
