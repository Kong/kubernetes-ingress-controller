package main

import (
	"github.com/kong/kubernetes-ingress-controller/railgun/cmd/rootcmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	rootcmd.Execute(ctrl.SetupSignalHandler())
}
