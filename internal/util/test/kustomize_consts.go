package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
)

var (
	kongRBACsKustomize        = initKongRBACsKustomizePath()
	kongGatewayRBACsKustomize = initKongGatewayRBACsKustomizePath()
	kongCRDsRBACsKustomize    = initKongCRDsRBACsKustomizePath()
	kongCRDsKustomize         = initCRDsKustomizePath()
)

func initKongRBACsKustomizePath() string {
	dir := filepath.Join(lo.Must(getRepoRoot()), "config/rbac/")
	ensureDirExists(dir)
	return dir
}

func initKongGatewayRBACsKustomizePath() string {
	dir := filepath.Join(lo.Must(getRepoRoot()), "config/rbac/gateway")
	ensureDirExists(dir)
	return dir
}

func initKongCRDsRBACsKustomizePath() string {
	dir := filepath.Join(lo.Must(getRepoRoot()), "config/rbac/crds")
	ensureDirExists(dir)
	return dir
}

func initCRDsKustomizePath() string {
	dir := filepath.Join(lo.Must(getRepoRoot()), "config/crd/")
	ensureDirExists(dir)
	return dir
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
