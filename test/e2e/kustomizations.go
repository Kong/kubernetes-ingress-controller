//go:build e2e_tests || istio_tests
// +build e2e_tests istio_tests

package e2e

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/kubectl"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/api/types"
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

	addControllerEnvPatch = `- op: add
  path: "/spec/template/spec/containers/1/env/-"
  value:
    name: %s
    value: "%s"`
)

// patchControllerImage replaces the kong/kubernetes-ingress-controller image with the provided image and tag,
// and returns the modified manifest.
func patchControllerImage(baseManifestReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Images: []types.Image{
			{
				Name:    "kong/kubernetes-ingress-controller",
				NewName: image,
				NewTag:  tag,
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchKongImage replaces the kong and kong/kong-gateway images in a manifest with the provided image and tag,
// and returns the modified manifest.
func patchKongImage(baseManifestReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
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
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchControllerStartTimeout adds or updates the controller container CONTROLLER_KONG_ADMIN_INIT_RETRIES and
// CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY environment variables with the provided values, and returns the modified
// manifest.
func patchControllerStartTimeout(baseManifestReader io.Reader, tries int, delay time.Duration) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(initRetryPatch, delay.String(), tries),
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
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchLivenessProbes patches the given container's liveness probe, replacing the initial delay, period, and failure
// threshold.
func patchLivenessProbes(baseManifestReader io.Reader, container, failure int, initial, period time.Duration) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(livenessProbePatch, container, int(initial.Seconds()), int(period.Seconds()), failure),
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
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// addControllerEnv adds an environment variable to ingress-controller container.
func addControllerEnv(t *testing.T, baseManifestReader io.Reader, envName, value string) io.Reader {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(addControllerEnvPatch, envName, value),
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
	k, err := kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
	require.NoError(t, err)
	return k
}
