package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
)

const (
	kubernetesConfigurationModulePath = "github.com/kong/kubernetes-configuration"
)

var (
	kongRBACsKustomize        = initKongRBACsKustomizePath()
	kongGatewayRBACsKustomize = initKongGatewayRBACsKustomizePath()
	kongCRDsRBACsKustomize    = initKongCRDsRBACsKustomizePath()

	kubernetesConfigurationModuleVersion = lo.Must(DependencyModuleVersion(kubernetesConfigurationModulePath))
	kongCRDsKustomize                    = initKongConfigurationCRDs()
	kongIncubatorCRDsKustomize           = initKongIncubatorCRDsKustomizePath()
)

func initKongIncubatorCRDsKustomizePath() string {
	return fmt.Sprintf("%s/config/crd/ingress-controller-incubator?ref=%s", kubernetesConfigurationModulePath, kubernetesConfigurationModuleVersion)
}

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

func initKongConfigurationCRDs() string {
	return fmt.Sprintf("%s/config/crd/ingress-controller?ref=%s", kubernetesConfigurationModulePath, kubernetesConfigurationModuleVersion)
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
