package sendconfig_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

type mockUpdateStrategy struct {
	wasUpdateCalled bool
	shouldSucceed   bool
}

func newMockUpdateStrategy(shouldSucceed bool) *mockUpdateStrategy {
	return &mockUpdateStrategy{shouldSucceed: shouldSucceed}
}

const mockUpdateReturnedConfigSize = 22

func (m *mockUpdateStrategy) Update(context.Context, sendconfig.ContentWithHash) (n int, err error) {
	m.wasUpdateCalled = true

	if !m.shouldSucceed {
		return 0, errors.New("update failure occurred")
	}

	return mockUpdateReturnedConfigSize, nil
}

func (m *mockUpdateStrategy) MetricsProtocol() metrics.Protocol {
	return "Mock"
}

func (m *mockUpdateStrategy) Type() string {
	return "Mock"
}

type mockBackoffStrategy struct {
	allowUpdate          bool
	wasSuccessRegistered bool
	wasFailureRegistered bool
}

func (m *mockBackoffStrategy) CanUpdate([]byte) (bool, string) {
	if m.allowUpdate {
		return true, ""
	}

	return false, "some reason"
}

func (m *mockBackoffStrategy) RegisterUpdateSuccess() {
	m.wasSuccessRegistered = true
}

func (m *mockBackoffStrategy) RegisterUpdateFailure(error, []byte) {
	m.wasFailureRegistered = true
}

func newMockBackoffStrategy(allowUpdate bool) *mockBackoffStrategy {
	return &mockBackoffStrategy{allowUpdate: allowUpdate}
}

func TestUpdateStrategyWithBackoff(t *testing.T) {
	ctx := context.Background()
	logger := zapr.NewLogger(zap.NewNop())

	testCases := []struct {
		name string

		updateShouldBeAllowed bool
		updateShouldSucceed   bool

		expectUpdateCalled      bool
		expectSuccessRegistered bool
		expectFailureRegistered bool
		expectError             error
	}{
		{
			name:                  "backoff allows update and it succeeds",
			updateShouldBeAllowed: true,
			updateShouldSucceed:   true,

			expectUpdateCalled:      true,
			expectSuccessRegistered: true,
		},
		{
			name:                  "backoff allows update and it fails",
			updateShouldBeAllowed: true,
			updateShouldSucceed:   false,

			expectUpdateCalled:      true,
			expectFailureRegistered: true,
			expectError:             errors.New("update failure occurred"),
		},
		{
			name:                  "backoff doesn't allow update, it doesn't happen and predefined error type is returned",
			updateShouldBeAllowed: false,

			expectUpdateCalled: false,
			expectError:        sendconfig.NewUpdateSkippedDueToBackoffStrategyError("some reason"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategy := newMockUpdateStrategy(tc.updateShouldSucceed)
			backoffStrategy := newMockBackoffStrategy(tc.updateShouldBeAllowed)

			decoratedStrategy := sendconfig.NewUpdateStrategyWithBackoff(updateStrategy, backoffStrategy, logger)
			size, err := decoratedStrategy.Update(ctx, sendconfig.ContentWithHash{})
			if tc.expectError != nil {
				require.Equal(t, tc.expectError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, mockUpdateReturnedConfigSize, size)
			}

			assert.Equal(t, tc.expectUpdateCalled, updateStrategy.wasUpdateCalled)
			assert.Equal(t, tc.expectSuccessRegistered, backoffStrategy.wasSuccessRegistered)
			assert.Equal(t, tc.expectFailureRegistered, backoffStrategy.wasFailureRegistered)
		})
	}
}

func TestUpdateSkippedDueToBackoffStrategyError(t *testing.T) {
	skippedErr := sendconfig.NewUpdateSkippedDueToBackoffStrategyError("reason")

	t.Run("errors.Is()", func(t *testing.T) {
		assert.False(t,
			errors.Is(skippedErr, deckerrors.ConfigConflictError{
				Err: sendconfig.NewUpdateSkippedDueToBackoffStrategyError("different reason"),
			}),
			"shouldn't panic when using errors.Is() with NewUpdateSkippedDueToBackoffStrategyError",
		)

		assert.False(t, errors.Is(skippedErr, errors.New("")),
			"empty error doesn't match",
		)
		assert.False(t, errors.Is(skippedErr, sendconfig.NewUpdateSkippedDueToBackoffStrategyError("different reason")),
			"error with different reason shouldn't match",
		)
		assert.True(t, errors.Is(skippedErr, skippedErr),
			"error with the same reason should match",
		)
	})

	t.Run("errors.As()", func(t *testing.T) {
		err := sendconfig.NewUpdateSkippedDueToBackoffStrategyError("reason")
		assert.True(t, errors.As(skippedErr, &err),
			"error with the same reason but wrapped should match",
		)
		err2 := sendconfig.NewUpdateSkippedDueToBackoffStrategyError("reason 2")
		assert.True(t, errors.As(skippedErr, &err2),
			"error with different reason should match",
		)
	})
}
