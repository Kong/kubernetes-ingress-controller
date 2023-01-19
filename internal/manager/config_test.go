package manager

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestValidatedVar(t *testing.T) {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)

	t.Run("string", func(t *testing.T) {
		validatedString := NewValidatedVar(func(s string) (string, error) {
			if !strings.Contains(s, "magic-token") {
				return "", errors.New("no magic token passed")
			}
			return s, nil
		})
		flags.Var(validatedString, "validated-string", "")

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
			v := validatedString.Get()
			require.Equal(t, "magic-token", v)
		})
	})

	t.Run("struct", func(t *testing.T) {
		type customType struct {
			p1, p2 string
		}
		validatedCustomType := NewValidatedVar(func(s string) (customType, error) {
			parts := strings.Split(s, "/")
			if len(parts) != 2 {
				return customType{}, fmt.Errorf("expected '<string>/<string>' format, got: %q", s)
			}

			return customType{p1: parts[0], p2: parts[1]}, nil
		})
		flags.Var(validatedCustomType, "custom-type", "")

		t.Run("valid", func(t *testing.T) {
			err := flags.Parse([]string{"--custom-type", "valid/value"})
			require.NoError(t, err)
			require.Equal(t, customType{p1: "valid", p2: "value"}, validatedCustomType.Get())
		})

		t.Run("invalid", func(t *testing.T) {
			err := flags.Parse([]string{"--custom-type", "invalid/format/"})
			require.ErrorContains(t, err, "expected '<string>/<string>'")
		})
	})
}
