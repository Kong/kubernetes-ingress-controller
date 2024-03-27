package deckerrors_test

import (
	"errors"
	"net/http"
	"testing"

	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckerrors"
)

func TestExtractAPIErrors(t *testing.T) {
	var (
		genericErr = errors.New("not an api error")
		apiErr     = kong.NewAPIError(http.StatusBadRequest, "api error")
	)

	testCases := []struct {
		name     string
		input    error
		expected []*kong.APIError
	}{
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "generic error",
			input:    genericErr,
			expected: nil,
		},
		{
			name:     "api error",
			input:    apiErr,
			expected: []*kong.APIError{apiErr},
		},
		{
			name:     "deck array of errors with no api error",
			input:    deckutils.ErrArray{Errors: []error{genericErr, genericErr}},
			expected: nil,
		},
		{
			name:     "deck array of errors with an api error among generic ones",
			input:    deckutils.ErrArray{Errors: []error{genericErr, apiErr, genericErr}},
			expected: []*kong.APIError{apiErr},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := deckerrors.ExtractAPIErrors(tc.input)
			require.Equal(t, tc.expected, out)
		})
	}
}
