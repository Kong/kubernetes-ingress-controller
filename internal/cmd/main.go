package main

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/internal/cmd/rootcmd"
)

//go:generate go run github.com/kong/kubernetes-ingress-controller/hack/generators/controllers/networking

func main() {
	rootcmd.Execute(ctrl.SetupSignalHandler())
}
