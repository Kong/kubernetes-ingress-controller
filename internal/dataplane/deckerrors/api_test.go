package deckerrors_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kong/go-database-reconciler/pkg/crud"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
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
			expected: []*kong.APIError{},
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

func TestExtractCRUDActionErrors(t *testing.T) {
	var (
		genericErr = errors.New("not an api error")
		actionErr  = &crud.ActionError{
			OperationType: crud.Create,
			Kind:          crud.Kind("service"),
			Name:          "badservice",
			Err:           errors.New("something wrong"),
		}
	)

	testCases := []struct {
		name     string
		input    error
		expected []*crud.ActionError
	}{
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "single generic error",
			input:    genericErr,
			expected: nil,
		},
		{
			name:     "single action error",
			input:    actionErr,
			expected: []*crud.ActionError{actionErr},
		},
		{
			name:     "deck array of errors with no action errors",
			input:    deckutils.ErrArray{Errors: []error{genericErr, genericErr}},
			expected: []*crud.ActionError{},
		},
		{
			name:     "deck array of errors with an action error among generic ones",
			input:    deckutils.ErrArray{Errors: []error{genericErr, actionErr, genericErr}},
			expected: []*crud.ActionError{actionErr},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := deckerrors.ExtractCRUDActionErrors(tc.input)
			require.Equal(t, tc.expected, out)
		})
	}
}
