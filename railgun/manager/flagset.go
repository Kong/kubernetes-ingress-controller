package manager

import (
	"flag"

	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

// flagSet extends flag.FlagSet with additional variable types.
type flagSet struct {
	flag.FlagSet
}

func (f *flagSet) EnablementStatusVar(p *util.EnablementStatus, name string, value util.EnablementStatus, usage string) {
	*p = value
	f.Var(p, name, usage)
}
