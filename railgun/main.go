package main

import (
	"github.com/kong/kubernetes-ingress-controller/railgun/cmd/rootcmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

//go:generate go run github.com/kong/kubernetes-ingress-controller/railgun/hack/generators/controllers/networking

func main() {
	rootcmd.Execute(ctrl.SetupSignalHandler())
}
