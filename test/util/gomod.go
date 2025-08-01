package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
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

	// NOTE: When we rely on a pseudo-version (e.g. `v1.1.1-0.20250217181409-44e5ddce290d`),
	// we need to extract the commit hash from it to use it in the GitHub API.
	// When version is set to a tag (e.g. `v1.1.1`), we could use it against the /commits
	// endpoint as well but there's no need for that as we can use the tag directly.

	// The same logic as in scripts/generate-crd-kustomize.sh:11:18

	// If there are 2 or more hyphens, extract the part after the last hyphen as
	// that's a git commit hash (e.g. `v1.1.1-0.20250217181409-44e5ddce290d`).
	if hash, ok := GetHashFromPseudoVersion(version); ok {
		resp, err := retryablehttp.Get("https://api.github.com/repos/Kong/kubernetes-configuration/commits/" + hash)
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

// GetHashFromPseudoVersion extracts the commit hash from a pseudo version string.
// It returns the hash and a boolean indicating whether the provided string
// is a valid pseudo version.
func GetHashFromPseudoVersion(pseudoVersion string) (string, bool) {
	pseudoVersion = strings.TrimSpace(pseudoVersion)
	pseudoVersion = strings.TrimPrefix(pseudoVersion, "v") // Remove leading 'v' if present

	// The pseudo version is expected to be in the format:
	// v1.1.1-0.20250217181409-44e5ddce290d
	// We need to extract the last part after the last hyphen.
	parts := strings.Split(pseudoVersion, "-")
	if len(parts) <= 2 {
		return "", false
	}

	if _, err := semver.Parse(parts[0]); err != nil {
		return "", false
	}
	return parts[len(parts)-1], true
}
