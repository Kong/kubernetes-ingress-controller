package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

const (
	_kongRBACsKustomize        = "config/rbac/"
	_kongGatewayRBACsKustomize = "config/rbac/gateway"
	_kongCRDsRBACsKustomize    = "config/rbac/crds"
	_kongCRDsKustomize         = "config/crd/"
)

var (
	kongRBACsKustomize        string
	kongGatewayRBACsKustomize string
	kongCRDsRBACsKustomize    string

	kongCRDsKustomize string
)

// init initializes kustomize paths relative to the repo root directory so that
// variables containing kustomize files can be used from anywhere in the repository.
func init() {
	root, err := getRepoRoot()
	if err != nil {
		panic(err)
	}

	kongRBACsKustomize = filepath.Join(root, _kongRBACsKustomize)
	ensureDirExists(kongRBACsKustomize)

	kongGatewayRBACsKustomize = filepath.Join(root, _kongGatewayRBACsKustomize)
	ensureDirExists(kongGatewayRBACsKustomize)

	kongCRDsRBACsKustomize = filepath.Join(root, _kongCRDsRBACsKustomize)
	ensureDirExists(kongCRDsRBACsKustomize)

	kongCRDsKustomize = filepath.Join(root, _kongCRDsKustomize)
	ensureDirExists(kongCRDsKustomize)
}

func ensureDirExists(dir string) {
	fi, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("%s is not a directory", dir))
	}
}

func getRepoRoot() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get repo root: runtime.Caller(0) failed")
	}
	d := filepath.Dir(path.Join(path.Dir(b), "../../")) // Number of ../ depends on the path of this file.
	return filepath.Abs(d)
}
