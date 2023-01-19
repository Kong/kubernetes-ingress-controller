package manager_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestValidatedValue(t *testing.T) {
	flags := func() *pflag.FlagSet { return pflag.NewFlagSet("", pflag.ContinueOnError) }

	t.Run("string", func(t *testing.T) {
		flags := flags()
		var validatedString string
		flags.Var(manager.NewValidatedValue(&validatedString, func(s string) (string, error) {
			if !strings.Contains(s, "magic-token") {
				return "", errors.New("no magic token passed")
			}
			return s, nil
		}), "validated-string", "")

		t.Run("invalid", func(t *testing.T) {
			err := flags.Parse([]string{
				"--validated-string", "invalid value",
			})
			require.Error(t, err)
		})

		t.Run("valid", func(t *testing.T) {
			err := flags.Parse([]string{
				"--validated-string", "magic-token",
			})
			require.NoError(t, err)
			require.Equal(t, "magic-token", validatedString)
		})
	})

	t.Run("struct", func(t *testing.T) {
		flags := flags()
		type customType struct {
			p1, p2 string
		}
		var customTypeVar customType
		flags.Var(manager.NewValidatedValue(&customTypeVar, func(s string) (customType, error) {
			parts := strings.Split(s, "/")
			if len(parts) != 2 {
				return customType{}, fmt.Errorf("expected '<string>/<string>' format, got: %q", s)
			}

			return customType{p1: parts[0], p2: parts[1]}, nil
		}), "custom-type", "")

		t.Run("valid", func(t *testing.T) {
			err := flags.Parse([]string{"--custom-type", "valid/value"})
			require.NoError(t, err)
			require.Equal(t, customType{p1: "valid", p2: "value"}, customTypeVar)
		})

		t.Run("invalid", func(t *testing.T) {
			err := flags.Parse([]string{"--custom-type", "invalid/format/"})
			require.ErrorContains(t, err, "expected '<string>/<string>'")
		})
	})

	t.Run("with default", func(t *testing.T) {
		flags := flags()
		var validatedString string
		flags.Var(manager.NewValidatedValueWithDefault(&validatedString, func(s string) (string, error) {
			if !strings.Contains(s, "magic-token") {
				return "", errors.New("no magic token passed")
			}
			return s, nil
		}, "default-value"), "flag-with-default", "")

		t.Run("empty", func(t *testing.T) {
			err := flags.Parse(nil)
			require.NoError(t, err)
			v := validatedString
			require.Equal(t, "default-value", v)
		})

		t.Run("invalid", func(t *testing.T) {
			err := flags.Parse([]string{
				"--flag-with-default", "invalid value",
			})
			require.Error(t, err)
		})

		t.Run("valid", func(t *testing.T) {
			err := flags.Parse([]string{
				"--flag-with-default", "magic-token",
			})
			require.NoError(t, err)
			require.Equal(t, "magic-token", validatedString)
		})
	})
}
