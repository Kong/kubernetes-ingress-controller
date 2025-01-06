package sendconfig

import (
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
)

// UpdateError wraps several pieces of error information relevant to a failed Kong update attempt.
type UpdateError struct {
	rawResponseBody  []byte
	configSize       mo.Option[int]
	resourceFailures []failures.ResourceFailure
	err              error
}

func NewUpdateErrorWithoutResponseBody(resourceFailures []failures.ResourceFailure, err error) UpdateError {
	return UpdateError{
		configSize:       mo.None[int](),
		resourceFailures: resourceFailures,
		err:              err,
	}
}

func NewUpdateErrorWithResponseBody(
	rawResponseBody []byte, configSize mo.Option[int], resourceFailures []failures.ResourceFailure, err error,
) UpdateError {
	return UpdateError{
		rawResponseBody:  rawResponseBody,
		configSize:       configSize,
		resourceFailures: resourceFailures,
		err:              err,
	}
}

// Error implements the Error interface. It returns the string value of the err field.
func (e UpdateError) Error() string {
	return e.err.Error()
}

// RawResponseBody returns the raw HTTP response body from Kong for the failed update if it was captured.
func (e UpdateError) RawResponseBody() []byte {
	return e.rawResponseBody
}

// ResourceFailures returns per-resource failures from a Kong configuration update attempt.
func (e UpdateError) ResourceFailures() []failures.ResourceFailure {
	return e.resourceFailures
}

// ConfigSize returns the size of the configuration that was attempted to be sent to Kong.
// When it's not applicable, returned option type contains mo.None.
func (e UpdateError) ConfigSize() mo.Option[int] {
	return e.configSize
}

func (e UpdateError) Unwrap() error {
	return e.err
}

// ResponseParsingError is an error type that is returned when the response from Kong is not in the expected format.
type ResponseParsingError struct {
	responseBody []byte
}

func NewResponseParsingError(responseBody []byte) ResponseParsingError {
	return ResponseParsingError{responseBody: responseBody}
}

// Error implements the Error interface.
func (e ResponseParsingError) Error() string {
	return "failed to parse Kong error response"
}

// ResponseBody returns the raw HTTP response body from Kong that could not be parsed.
func (e ResponseParsingError) ResponseBody() []byte {
	return e.responseBody
}
