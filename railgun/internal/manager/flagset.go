package manager

import (
	"github.com/spf13/pflag"
)

// -----------------------------------------------------------------------------
// Controller Manager - flagSet
// -----------------------------------------------------------------------------

// flagSet extends flag.FlagSet with additional variable types.
type flagSet struct {
	pflag.FlagSet
}
