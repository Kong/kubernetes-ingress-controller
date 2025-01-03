package flags

import "fmt"

type ValidatedValueOpt[T any] func(*ValidatedValue[T])

// WithDefault sets the default value for the validated variable.
func WithDefault[T any](defaultValue T) ValidatedValueOpt[T] {
	return func(v *ValidatedValue[T]) {
		*v.variable = defaultValue

		// Assign origin which is used in ValidatedValue[T]'s String() string
		// func so that we get a pretty printed default.
		v.origin = stringFromAny(defaultValue)
	}
}

func stringFromAny(s any) string {
	switch ss := s.(type) {
	case string:
		return ss
	case fmt.Stringer:
		return fmt.Sprintf("%q", ss.String())
	default:
		panic(fmt.Errorf("unknown type %T", s))
	}
}

// WithTypeNameOverride overrides the type name that's printed in the help message.
func WithTypeNameOverride[T any](typeName string) ValidatedValueOpt[T] {
	return func(v *ValidatedValue[T]) {
		v.typeName = typeName
	}
}

// ValidatedValue implements `pflag.Value` interface. It can be used for hooking up arbitrary validation logic to any type.
// It should be passed to `pflag.FlagSet.Var()`.
type ValidatedValue[T any] struct {
	origin      string
	variable    *T
	constructor func(string) (T, error)
	typeName    string
}

// NewValidatedValue creates a validated variable of type T. Constructor should validate the input and return an error
// in case of any failures. If validation passes, constructor should return a value that's to be set in the variable.
// The constructor accepts a flagValue that is raw input from user's command line (or an env variable that was bind to
// the flag, see: bindEnvVars).
// It accepts a variadic list of options that can be used e.g. to set the default value or override the type name.
func NewValidatedValue[T any](variable *T, constructor func(flagValue string) (T, error), opts ...ValidatedValueOpt[T]) ValidatedValue[T] {
	v := ValidatedValue[T]{
		constructor: constructor,
		variable:    variable,
	}
	for _, opt := range opts {
		opt(&v)
	}
	return v
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
	if v.typeName != "" {
		return v.typeName
	}

	var t T
	return fmt.Sprintf("%T", t)
}
