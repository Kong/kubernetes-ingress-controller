package manager

import "fmt"

// ValidatedValue implements `pflag.Value` interface. It can be used for hooking up arbitrary validation logic to any type.
// It should be passed to `pflag.FlagSet.Var()`.
type ValidatedValue[T any] struct {
	origin      string
	variable    *T
	constructor func(string) (T, error)
}

// NewValidatedValue creates a validated variable of type T. Constructor should validate the input and return an error
// in case of any failures. If validation passes, constructor should return a value that's to be set in the variable.
// The constructor accepts a flagValue that is raw input from user's command line (or an env variable that was bind to
// the flag, see: bindEnvVars).
func NewValidatedValue[T any](variable *T, constructor func(flagValue string) (T, error)) ValidatedValue[T] {
	return ValidatedValue[T]{
		constructor: constructor,
		variable:    variable,
	}
}

// NewValidatedValueWithDefault creates a validated variable of type T with a default value.
func NewValidatedValueWithDefault[T any](variable *T, constructor func(flagValue string) (T, error), value T) ValidatedValue[T] {
	*variable = value
	return ValidatedValue[T]{
		constructor: constructor,
		variable:    variable,
	}
}

func (v ValidatedValue[T]) String() string {
	return v.origin
}

func (v ValidatedValue[T]) Set(s string) error {
	value, err := v.constructor(s)
	if err != nil {
		return err
	}

	*v.variable = value
	return nil
}

func (v ValidatedValue[T]) Type() string {
	var t T
	return fmt.Sprintf("Validated%T", t)
}
