package test

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

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

func getRepoRoot() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get repo root: runtime.Caller(0) failed")
	}
	d := filepath.Dir(path.Join(path.Dir(b), "../../")) // Number of ../ depends on the path of this file.
	return filepath.Abs(d)
}
