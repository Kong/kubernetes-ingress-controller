package flags_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/flags"
)

func TestValidatedValue(t *testing.T) {
	createFlagSet := func() *pflag.FlagSet { return pflag.NewFlagSet("", pflag.ContinueOnError) }

	t.Run("string", func(t *testing.T) {
		flagSet := createFlagSet()
		var validatedString string
		flagSet.Var(flags.NewValidatedValue(&validatedString, func(s string) (string, error) {
			if !strings.Contains(s, "magic-token") {
				return "", errors.New("no magic token passed")
			}
			return s, nil
		}), "validated-string", "")

		t.Run("invalid", func(t *testing.T) {
			err := flagSet.Parse([]string{
				"--validated-string", "invalid value",
			})
			require.Error(t, err)
		})

		t.Run("valid", func(t *testing.T) {
			err := flagSet.Parse([]string{
				"--validated-string", "magic-token",
			})
			require.NoError(t, err)
			require.Equal(t, "magic-token", validatedString)
		})
	})

	t.Run("struct", func(t *testing.T) {
		flagSet := createFlagSet()
		type customType struct {
			p1, p2 string
		}
		var customTypeVar customType
		flagSet.Var(flags.NewValidatedValue(&customTypeVar, func(s string) (customType, error) {
			parts := strings.Split(s, "/")
			if len(parts) != 2 {
				return customType{}, fmt.Errorf("expected '<string>/<string>' format, got: %q", s)
			}

			return customType{p1: parts[0], p2: parts[1]}, nil
		}), "custom-type", "")

		t.Run("valid", func(t *testing.T) {
			err := flagSet.Parse([]string{"--custom-type", "valid/value"})
			require.NoError(t, err)
			require.Equal(t, customType{p1: "valid", p2: "value"}, customTypeVar)
		})

		t.Run("invalid", func(t *testing.T) {
			err := flagSet.Parse([]string{"--custom-type", "invalid/format/"})
			require.ErrorContains(t, err, "expected '<string>/<string>'")
		})
	})

	t.Run("with default", func(t *testing.T) {
		flagSet := createFlagSet()

		var validatedString string
		flagSet.Var(flags.NewValidatedValue(&validatedString, func(s string) (string, error) {
			if !strings.Contains(s, "magic-token") {
				return "", errors.New("no magic token passed")
			}
			return s, nil
		}, flags.WithDefault("default-value")), "flag-with-default", "")

		t.Run("empty", func(t *testing.T) {
			err := flagSet.Parse(nil)
			require.NoError(t, err)
			v := validatedString
			require.Equal(t, "default-value", v)
		})

		t.Run("invalid", func(t *testing.T) {
			err := flagSet.Parse([]string{
				"--flag-with-default", "invalid value",
			})
			require.Error(t, err)
		})

		t.Run("valid", func(t *testing.T) {
			err := flagSet.Parse([]string{
				"--flag-with-default", "magic-token",
			})
			require.NoError(t, err)
			require.Equal(t, "magic-token", validatedString)
		})
	})
}

type customStringer struct{}

func (cs customStringer) String() string {
	return "custom-string-default"
}

func TestValidatedValue_WithDefault(t *testing.T) {
	createFlagSet := func() *pflag.FlagSet { return pflag.NewFlagSet("", pflag.ContinueOnError) }

	t.Run("default printed in usage for string flag", func(t *testing.T) {
		flagSet := createFlagSet()

		var validatedString string
		flagSet.Var(flags.NewValidatedValue(&validatedString, func(s string) (string, error) {
			return s, nil
		}, flags.WithDefault("default-value")), "flag-with-default", "")

		b := bytes.Buffer{}
		flagSet.SetOutput(&b)
		flagSet.PrintDefaults()
		require.Contains(t, b.String(), `(default "default-value")`)
	})

	t.Run("default printed in usage for fmt.Stringer flag", func(t *testing.T) {
		flagSet := createFlagSet()

		var cs customStringer
		flagSet.Var(flags.NewValidatedValue(&cs, func(_ string) (customStringer, error) {
			return customStringer{}, nil
		}, flags.WithDefault(customStringer{})), "flag-with-default", "")

		b := bytes.Buffer{}
		flagSet.SetOutput(&b)
		flagSet.PrintDefaults()
		require.Contains(t, b.String(), `(default "custom-string-default")`)
	})
}

func TestValidatedValue_Type(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var validatedString string
		vv := flags.NewValidatedValue(&validatedString, func(s string) (string, error) {
			return s, nil
		})
		require.Equal(t, "string", vv.Type())
	})

	t.Run("struct", func(t *testing.T) {
		type customType struct{}
		var customTypeVar customType
		vv := flags.NewValidatedValue(&customTypeVar, func(_ string) (customType, error) {
			return customType{}, nil
		})
		require.Equal(t, "flags_test.customType", vv.Type())
	})

	t.Run("overridden type name", func(t *testing.T) {
		type customType struct{}
		var customTypeVar customType
		vv := flags.NewValidatedValue(&customTypeVar, func(_ string) (customType, error) {
			return customType{}, nil
		}, flags.WithTypeNameOverride[customType]("custom-type-override"))
		require.Equal(t, "custom-type-override", vv.Type())
	})
}
