package deckerrors

import (
	"errors"

	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
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
		var apiErrs []*kong.APIError
		for _, err := range deckErrArray.Errors {
			if apiErr, ok := castAsErr[*kong.APIError](err); ok {
				apiErrs = append(apiErrs, apiErr)
			}
		}
		return apiErrs
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
