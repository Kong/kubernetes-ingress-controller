//go:build fips
// +build fips

package main

import (
	_ "crypto/tls/fipsonly"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
)

//go:generate go run github.com/kong/kubernetes-ingress-controller/v2/hack/generators/controllers/networking

func main() {
	rootcmd.Execute()
}
