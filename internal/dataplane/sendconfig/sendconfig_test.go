package sendconfig

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
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

func TestPushFailureReason(t *testing.T) {
	apiConflictErr := kong.NewAPIError(http.StatusConflict, "conflict api error")
	networkErr := net.UnknownNetworkError("network error")
	genericError := errors.New("generic error")

	testCases := []struct {
		name           string
		err            error
		expectedReason string
	}{
		{
			name:           "generic_error",
			err:            genericError,
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "api_conflict_error",
			err:            apiConflictErr,
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "api_conflict_error_wrapped",
			err:            fmt.Errorf("wrapped conflict api err: %w", apiConflictErr),
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_empty",
			err:            deckConfigConflictError{},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_with_generic_error",
			err:            deckConfigConflictError{genericError},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_api_conflict_error",
			err:            deckutils.ErrArray{Errors: []error{apiConflictErr}},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "wrapped_deck_err_array_with_api_conflict_error",
			err:            fmt.Errorf("wrapped: %w", deckutils.ErrArray{Errors: []error{apiConflictErr}}),
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_generic_error",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "deck_err_array_empty",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "network_error",
			err:            networkErr,
			expectedReason: metrics.FailureReasonNetwork,
		},
		{
			name:           "network_error_wrapped_in_deck_config_conflict_error",
			err:            deckConfigConflictError{networkErr},
			expectedReason: metrics.FailureReasonNetwork,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reason := pushFailureReason(tc.err)
			require.Equal(t, tc.expectedReason, reason)
		})
	}
}
