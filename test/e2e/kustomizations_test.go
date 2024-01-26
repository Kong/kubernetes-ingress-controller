//go:build e2e_tests || istio_tests || performance_tests

package e2e

import (
	"fmt"
	"io"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/kubectl"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

const (
	initRetryPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY
    value: "%s"
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRIES
    value: "%d"`

	livenessProbePatch = `- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/initialDelaySeconds
  value: %[1]d
- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/timeoutSeconds
  value: %[2]d
- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/failureThreshold
  value: %[3]d`

	kongRouterFlavorPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: KONG_ROUTER_FLAVOR
    value: "%s"`

	kongRouterFlavorPatchDelete = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: proxy-kong
  namespace: kong
spec:
  template:
    spec:
      containers:
      - name: proxy
        env:
        - name: KONG_ROUTER_FLAVOR
          $patch: delete`
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
						Name:      controllerDeploymentName,
						Namespace: "kong",
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

func patchKongRouterFlavorFn(flavor string) func(io.Reader) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "proxy-kong",
						Namespace: "kong",
					},
				},
				Patch: kongRouterFlavorPatchDelete,
			},
			{
				Patch: fmt.Sprintf(kongRouterFlavorPatch, flavor),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "proxy-kong",
						Namespace: "kong",
					},
				},
			},
		},
	}

	return func(baseManifestReader io.Reader) (io.Reader, error) {
		return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
	}
}

// patchLivenessProbes patches the given deployment's liveness probe, replacing the initial delay, period, and failure
// threshold.
func patchLivenessProbes(baseManifestReader io.Reader, deployment k8stypes.NamespacedName, failure int, initial, period time.Duration) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(livenessProbePatch, int(initial.Seconds()), int(period.Seconds()), failure),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      deployment.Name,
						Namespace: deployment.Namespace,
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}
