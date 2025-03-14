package sendconfig_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

type statusClientMock struct {
	expectedValue *kong.Status
	expectedError error
}

func (c statusClientMock) Status(context.Context) (*kong.Status, error) {
	return c.expectedValue, c.expectedError
}

func TestDefaultConfigurationChangeDetector_HasConfigurationChanged(t *testing.T) {
	ctx := t.Context()
	testSHAs := [][]byte{
		[]byte("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"),
		[]byte("82e35a63ceba37e9646434c5dd412ea577147f1e4a41ccde1614253187e3dbf9"),
	}

	createConfigContent := func() *file.Content {
		return &file.Content{
			FormatVersion: "1.1",
			Services: []file.FService{
				{
					Service: kong.Service{
						ID:   kong.String("id"),
						Name: kong.String("name"),
					},
				},
			},
		}
	}

	testCases := []struct {
		name           string
		oldSHA, newSHA []byte
		targetConfig   *file.Content
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
			targetConfig:   createConfigContent(),
			expectedResult: true,
		},
		{
			name:           "oldSHA == newSHA, but status signals crash",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			targetConfig:   createConfigContent(),
			statusSHA:      sendconfig.WellKnownInitialHash,
			expectedResult: true,
		},
		{
			name:   "oldSHA == newSHA, status signals init hash and we're trying to push empty config",
			oldSHA: testSHAs[0],
			newSHA: testSHAs[0],
			targetConfig: &file.Content{
				FormatVersion: "1.1",
			},
			statusSHA:      sendconfig.WellKnownInitialHash,
			expectedResult: false,
		},
		{
			name:           "oldSHA == newSHA and status signals same",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			targetConfig:   createConfigContent(),
			statusSHA:      string(testSHAs[0]),
			expectedResult: false,
		},
		{
			name:           "oldSHA == newSHA and status signals other",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			targetConfig:   createConfigContent(),
			statusSHA:      string(testSHAs[1]),
			expectedResult: false,
		},
		{
			name:         "oldSHA == newSHA, status returns error",
			oldSHA:       testSHAs[0],
			newSHA:       testSHAs[0],
			targetConfig: createConfigContent(),
			statusError:  errors.New("getting kong status failure"),
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			statusClient := statusClientMock{
				expectedValue: &kong.Status{
					ConfigurationHash: tc.statusSHA,
				},
				expectedError: tc.statusError,
			}
			detector := sendconfig.NewKongGatewayConfigurationChangeDetector(zapr.NewLogger(zap.NewNop()))

			result, err := detector.HasConfigurationChanged(ctx, tc.oldSHA, tc.newSHA, tc.targetConfig, statusClient)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestKonnectConfigurationChangeDetector(t *testing.T) {
	ctx := t.Context()
	testSHAs := [][]byte{
		[]byte("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"),
		[]byte("82e35a63ceba37e9646434c5dd412ea577147f1e4a41ccde1614253187e3dbf9"),
	}

	testCases := []struct {
		name           string
		oldSHA, newSHA []byte
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
			name:           "oldSHA == newSHA",
			oldSHA:         testSHAs[0],
			newSHA:         testSHAs[0],
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			detector := sendconfig.NewKonnectConfigurationChangeDetector()

			// Passing nil for content and status client explicitly, as they are not used in this detector.
			result, err := detector.HasConfigurationChanged(ctx, tc.oldSHA, tc.newSHA, nil, nil)
			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
