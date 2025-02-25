package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
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

	var extractedVersion string
	// If there's a replace directive, use the replace version.
	for _, r := range goMod.Replace {
		if r.Old.Path == dep {
			extractedVersion = r.New.Version
		}
	}
	if extractedVersion == "" {
		// Otherwise, use the require version.
		for _, r := range goMod.Require {
			if r.Mod.Path == dep {
				extractedVersion = r.Mod.Version
			}
		}
	}
	if extractedVersion == "" {
		// If the module is not found, return an error.
		return "", fmt.Errorf("%s not found in go.mod's require nor replace section", dep)
	}
	return extractedVersion, nil
}

// DependencyModuleVersionGit returns the version of the dependency module based on the project's go.mod file
// and returns it in git compatible format, instead of go.mod one.
func DependencyModuleVersionGit(dep string) (string, error) {
	version, err := DependencyModuleVersion(dep)
	if err != nil {
		return "", err
	}

	// The same logic as in scripts/generate-crd-kustomize.sh:11:18
	if versionParts := strings.Split(version, "-"); len(versionParts) >= 2 {
		// If there are 2 or more hyphens, extract the part after the last hyphen as
		// that's a git commit hash (e.g. `v1.1.1-0.20250217181409-44e5ddce290d`).
		version = versionParts[len(versionParts)-1]

		resp, err := retryablehttp.Get("https://api.github.com/repos/Kong/kubernetes-configuration/commits/" + version)
		if err != nil {
			return "", fmt.Errorf("failed to fetch commit data: %w", err)
		}
		defer resp.Body.Close()

		response := struct {
			SHA string `json:"sha"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", fmt.Errorf("failed to decode response for commit details: %w", err)
		}
		version = response.SHA
	}

	return version, nil
}
