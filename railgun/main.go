package main

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/railgun/cmd/rootcmd"
)

//go:generate go run github.com/kong/kubernetes-ingress-controller/railgun/hack/generators/controllers/networking

func main() {
	rootcmd.Execute(ctrl.SetupSignalHandler())
}
