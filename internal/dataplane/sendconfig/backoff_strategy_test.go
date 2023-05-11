package sendconfig_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

type mockUpdateStrategy struct {
	wasUpdateCalled bool
	shouldSucceed   bool
}

func newMockUpdateStrategy(shouldSucceed bool) *mockUpdateStrategy {
	return &mockUpdateStrategy{shouldSucceed: shouldSucceed}
}

func (m *mockUpdateStrategy) Update(context.Context, sendconfig.ContentWithHash) (
	err error,
	resourceErrors []sendconfig.ResourceError,
	resourceErrorsParseErr error,
) {
	m.wasUpdateCalled = true

	if !m.shouldSucceed {
		return errors.New("update failure occurred"), nil, nil
	}

	return nil, nil, nil
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
	log := logrus.New()

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
			expectError:        sendconfig.NewErrUpdateSkippedDueToBackoffStrategy("some reason"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategy := newMockUpdateStrategy(tc.updateShouldSucceed)
			backoffStrategy := newMockBackoffStrategy(tc.updateShouldBeAllowed)

			decoratedStrategy := sendconfig.NewUpdateStrategyWithBackoff(updateStrategy, backoffStrategy, log)
			err, _, _ := decoratedStrategy.Update(ctx, sendconfig.ContentWithHash{})
			if tc.expectError != nil {
				require.Equal(t, tc.expectError, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, tc.expectUpdateCalled, updateStrategy.wasUpdateCalled)
			assert.Equal(t, tc.expectSuccessRegistered, backoffStrategy.wasSuccessRegistered)
			assert.Equal(t, tc.expectFailureRegistered, backoffStrategy.wasFailureRegistered)
		})
	}
}
