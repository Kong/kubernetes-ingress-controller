package manager

import (
	"github.com/spf13/pflag"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

// -----------------------------------------------------------------------------
// Controller Manager - flagSet
// -----------------------------------------------------------------------------

// flagSet extends flag.FlagSet with additional variable types.
type flagSet struct {
	pflag.FlagSet
}

// enablementStatusVar defines a flag of type EnablementStatus.
func (f *flagSet) enablementStatusVar(p *util.EnablementStatus, name string, value util.EnablementStatus, usage string) {
	*p = value
	f.Var(p, name, usage)
}
