package deckerrors

import (
	"errors"

	"github.com/kong/go-database-reconciler/pkg/crud"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
)

// ExtractAPIErrors tries to extract kong.APIErrors from the generic error.
// It might be used when inspection of the error details is needed, e.g. its status code.
func ExtractAPIErrors(err error) []*kong.APIError {
	// It might be a single APIError.
	if apiErr, ok := castAsErr[*kong.APIError](err); ok {
		return []*kong.APIError{apiErr}
	}

	// It might be either a deckutils.ErrArray with APIErrors inside.
	var deckErrArray deckutils.ErrArray
	if errors.As(err, &deckErrArray) {
		return lo.FilterMap(deckErrArray.Errors, func(e error, _ int) (*kong.APIError, bool) {
			return castAsErr[*kong.APIError](e)
		})
	}

	return nil
}

func ExtractCRUDActionErrors(err error) []*crud.ActionError {
	// It might be a single crud.ActionError.
	if actionErr, ok := castAsErr[*crud.ActionError](err); ok {
		return []*crud.ActionError{actionErr}
	}

	// It might be either a deckutils.ErrArray with ActionErrors inside.
	var deckErrArray deckutils.ErrArray
	if errors.As(err, &deckErrArray) {
		return lo.FilterMap(deckErrArray.Errors, func(e error, _ int) (*crud.ActionError, bool) {
			return castAsErr[*crud.ActionError](e)
		})
	}

	return nil
}

func castAsErr[T error](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}
	return target, false
}
