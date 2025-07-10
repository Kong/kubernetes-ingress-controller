package deckerrors

import (
	"errors"
	"net/http"
	"slices"

	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
)

// ConfigConflictError is an error used to wrap deck config conflict errors
// returned from deck functions transforming KongRawState to KongState
// (e.g. state.Get, dump.Get).
type ConfigConflictError struct {
	Err error
}

func (e ConfigConflictError) Error() string {
	return e.Err.Error()
}

func (e ConfigConflictError) Is(err error) bool {
	return errors.Is(err, ConfigConflictError{})
}

func (e ConfigConflictError) Unwrap() error {
	return e.Err
}

func IsConflictErr(err error) bool {
	var apiErr *kong.APIError
	if errors.As(err, &apiErr) && apiErr.Code() == http.StatusConflict ||
		errors.Is(err, ConfigConflictError{}) {
		return true
	}

	var deckErrArray deckutils.ErrArray
	if errors.As(err, &deckErrArray) {
		if slices.ContainsFunc(deckErrArray.Errors, IsConflictErr) {
			return true
		}
	}

	return false
}
