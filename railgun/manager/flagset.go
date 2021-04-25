package manager

import (
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/spf13/pflag"
)

// flagSet extends flag.FlagSet with additional variable types.
type flagSet struct {
	pflag.FlagSet
}

// EnablementStatusVar defines a flag of type EnablementStatus.
func (f *flagSet) EnablementStatusVar(p *util.EnablementStatus, name string, value util.EnablementStatus, usage string) {
	*p = value
	f.Var(p, name, usage)
}
