package sendconfig

import (
	"context"
	"errors"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateReportingUtilities(t *testing.T) {
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
}

type konnectAwareClientMock struct {
	expected bool
}

func (c konnectAwareClientMock) IsKonnect() bool {
	return c.expected
}

type statusClientMock struct {
	expectedValue *kong.Status
	expectedError error
}

func (c statusClientMock) Status(context.Context) (*kong.Status, error) {
	return c.expectedValue, c.expectedError
}

func TestHasConfigurationChanged(t *testing.T) {
	ctx := context.Background()
	testSHAs := [][]byte{
		[]byte("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"),
		[]byte("82e35a63ceba37e9646434c5dd412ea577147f1e4a41ccde1614253187e3dbf9"),
	}

	testCases := []struct {
		name           string
		oldSHA, newSHA []byte
		statusSHA      string
		isKonnect      bool
		statusError    error

		expectedResult bool
		expectError    bool
	}{
		{
			name:           "oldSHA != newSHA",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[1],
			expectedResult: true,
		},
		{
			name:           "oldSHA == newSHA, but status signals crash",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			statusSHA:      wellKnownInitialHash,
			expectedResult: true,
		},
		{
			name:           "oldSHA == newSHA and status signals same",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			statusSHA:      string(testSHAs[0]),
			expectedResult: false,
		},
		{
			name:           "oldSHA == newSHA and status signals other",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			statusSHA:      string(testSHAs[1]),
			expectedResult: false,
		},
		{
			name:           "oldSHA == newSHA, status would signal crash, but it's konnect",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			statusSHA:      wellKnownInitialHash,
			isKonnect:      true,
			expectedResult: false,
		},
		{
			name:        "oldSHA == newSHA, status returns error",
			oldSHA:      testSHAs[0],
			newSHA:      testSHAs[0],
			statusError: errors.New("getting kong status failure"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			konnectAwareClient := konnectAwareClientMock{expected: tc.isKonnect}
			statusClient := statusClientMock{
				expectedValue: &kong.Status{
					ConfigurationHash: tc.statusSHA,
				},
				expectedError: tc.statusError,
			}

			result, err := hasConfigurationChanged(ctx, tc.oldSHA, tc.newSHA, konnectAwareClient, statusClient, logrus.New())
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
