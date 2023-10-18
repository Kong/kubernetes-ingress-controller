package test

import "path/filepath"

var (
	kongRBACsKustomize        = "config/rbac/"
	kongGatewayRBACsKustomize = "config/rbac/gateway"
	kongCRDsRBACsKustomize    = "config/rbac/crds"

	kongCRDsKustomize = "config/crd/"
)

// init initializes kustomize paths relative to the repo root directory.
func init() {
	root, err := getRepoRoot()
	if err != nil {
		panic(err)
	}

	kongCRDsKustomize = filepath.Join(root, kongCRDsKustomize)

	kongRBACsKustomize = filepath.Join(root, kongRBACsKustomize)
	kongGatewayRBACsKustomize = filepath.Join(root, kongGatewayRBACsKustomize)
	kongCRDsRBACsKustomize = filepath.Join(root, kongCRDsRBACsKustomize)
}
