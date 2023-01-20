package manager

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

// FlagNamespacedName allows parsing command line flags straight to types.NamespacedName.
type FlagNamespacedName struct {
	NN types.NamespacedName
}

func (f *FlagNamespacedName) String() string {
	return f.NN.String()
}

func (f *FlagNamespacedName) Set(v string) error {
	s := strings.SplitN(v, "/", 3)
	if len(s) != 2 {
		return fmt.Errorf("namespaced name should be in the format: <namespace>/<name>")
	}
	f.NN = types.NamespacedName{
		Namespace: s[0],
		Name:      s[1],
	}
	return nil
}

func (f *FlagNamespacedName) Type() string {
	return "FlagNamespacedName"
}
