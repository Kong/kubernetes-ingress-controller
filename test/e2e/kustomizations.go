package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

const (
	initRetryPatch = `- op: add
  path: /spec/template/spec/containers/1/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY
    value: "%s"
- op: add
  path: /spec/template/spec/containers/1/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRIES
    value: "%d"`

	livenessProbePatch = `- op: replace
  path: /spec/template/spec/containers/%[1]d/livenessProbe/initialDelaySeconds
  value: %[2]d
- op: replace
  path: /spec/template/spec/containers/%[1]d/livenessProbe/timeoutSeconds
  value: %[3]d
- op: replace
  path: /spec/template/spec/containers/%[1]d/livenessProbe/failureThreshold
  value: %[4]d`
)

// getKustomizedManifest takes a base manifest Reader and a Kustomization, and returns the manifest after applying the
// Kustomization. It overwrites the Kustomization's Bases; any existing Bases will be discarded.
func getKustomizedManifest(baseManifestReader io.Reader, kustomization types.Kustomization) (io.Reader, error) {
	workDir, err := os.MkdirTemp("", "kictest.")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workDir)
	orig, err := io.ReadAll(baseManifestReader)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(workDir, "base.yaml"), orig, 0o600)
	if err != nil {
		return nil, err
	}
	kustomization.Bases = []string{"base.yaml"}
	marshalled, err := yaml.Marshal(kustomization)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(workDir, "kustomization.yaml"), marshalled, 0o600)
	if err != nil {
		return nil, err
	}
	kustomized, err := runKustomize(workDir)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(kustomized), nil
}

// patchControllerImage replaces the kong/kubernetes-ingress-controller image with the provided image and tag,
// and returns the modified manifest.
func patchControllerImage(baseManifestReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Bases: []string{"base.yaml"},
		Images: []types.Image{
			{
				Name:    "kong/kubernetes-ingress-controller",
				NewName: image,
				NewTag:  tag,
			},
		},
	}
	return getKustomizedManifest(baseManifestReader, kustomization)
}

// patchKongImage replaces the kong and kong/kong-gateway images in a manifest with the provide image and tag,
// and returns the modified manifest.
func patchKongImage(baseManifestsReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Bases: []string{"base.yaml"},
		Images: []types.Image{
			{
				Name:    "kong/kong-gateway",
				NewName: image,
				NewTag:  tag,
			},
			{
				Name:    "kong",
				NewName: image,
				NewTag:  tag,
			},
		},
	}
	return getKustomizedManifest(baseManifestsReader, kustomization)
}

// patchControllerStartTimeout adds or updates the controller container CONTROLLER_KONG_ADMIN_INIT_RETRIES and
// CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY environment variables with the provided values, and returns the modified
// manifest.
func patchControllerStartTimeout(baseManifestReader io.Reader, tries int, delay string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Bases: []string{"base.yaml"},
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(initRetryPatch, delay, tries),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "ingress-kong",
						Namespace: "kong",
					},
				},
			},
		},
	}
	return getKustomizedManifest(baseManifestReader, kustomization)
}

// patchLivenessProbes patches the given container's liveness probe, replacing the initial delay, period, and failure
// threshold.
func patchLivenessProbes(baseManifestReader io.Reader, container, initial, period, failure int) (io.Reader, error) {
	kustomization := types.Kustomization{
		Bases: []string{"base.yaml"},
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(livenessProbePatch, container, initial, period, failure),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "ingress-kong",
						Namespace: "kong",
					},
				},
			},
		},
	}
	return getKustomizedManifest(baseManifestReader, kustomization)
}

// runKustomize runs kustomize on a path and returns the YAML output.
func runKustomize(path string) ([]byte, error) {
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	m, err := k.Run(filesys.MakeFsOnDisk(), path)
	if err != nil {
		return []byte{}, err
	}
	return m.AsYaml()
}
