package test

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
)

// DependencyModuleVersion returns the version of the dependency module based on the project's go.mod file.
func DependencyModuleVersion(dep string) (string, error) {
	var (
		goModPath = lo.Must(getRepoRoot()) + "/go.mod"
		goModData = lo.Must(os.ReadFile(goModPath))
		goMod     = lo.Must(modfile.Parse(goModPath, goModData, nil))
	)
	// If there's a replace directive, use the replace version.
	for _, r := range goMod.Replace {
		if r.Old.Path == dep {
			return r.New.Version, nil
		}
	}
	// Otherwise, use the require version.
	for _, r := range goMod.Require {
		if r.Mod.Path == dep {
			return r.Mod.Version, nil
		}
	}
	// If the module is not found, return an error.
	return "", fmt.Errorf("%s not found in go.mod's require nor replace section", dep)
}
